from concurrent import futures

import grpc

import python.person_pb2 as person_pb2
import python.person_pb2_grpc as person_pb2_grpc


# In-memory database
persons = {}

# Used to automatically generate IDs
next_id = 1


class PersonService(person_pb2_grpc.PersonServiceServicer):

    def CreatePerson(self, request, context):
        global next_id

        # Basic validation
        if not request.name:
            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Name is required"
            )

        if not request.email:
            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Email is required"
            )

        # Create the Person object
        person = person_pb2.Person(
            id=next_id,
            name=request.name,
            email=request.email,
            phone=request.phone
        )

        # Store the person in memory
        persons[next_id] = person

        print(f"Created person with ID {next_id}")

        next_id += 1

        return person


    def GetPerson(self, request, context):

        person_id = request.id

        if person_id not in persons:
            context.abort(
                grpc.StatusCode.NOT_FOUND,
                f"Person with ID {person_id} not found"
            )

        print(f"Reading person with ID {person_id}")

        return persons[person_id]


    def UpdatePerson(self, request, context):

        person_id = request.id

        if person_id not in persons:
            context.abort(
                grpc.StatusCode.NOT_FOUND,
                f"Person with ID {person_id} not found"
            )

        if not request.name:
            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Name is required"
            )

        if not request.email:
            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Email is required"
            )

        updated_person = person_pb2.Person(
            id=person_id,
            name=request.name,
            email=request.email,
            phone=request.phone
        )

        persons[person_id] = updated_person

        print(f"Updated person with ID {person_id}")

        return updated_person


    def DeletePerson(self, request, context):

        person_id = request.id

        if person_id not in persons:
            context.abort(
                grpc.StatusCode.NOT_FOUND,
                f"Person with ID {person_id} not found"
            )

        del persons[person_id]

        print(f"Deleted person with ID {person_id}")

        return person_pb2.DeletePersonResponse(
            message=f"Person {person_id} deleted successfully"
        )


def serve():

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10)
    )

    person_pb2_grpc.add_PersonServiceServicer_to_server(
        PersonService(),
        server
    )

    server.add_insecure_port(
        "localhost:50051"
    )

    server.start()

    print("Person gRPC Server started on port 50051")

    server.wait_for_termination()


if __name__ == "__main__":
    serve()