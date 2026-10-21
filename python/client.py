import grpc

import python.person_pb2 as person_pb2
import python.person_pb2_grpc as person_pb2_grpc


def create_person(stub):

    print("\n--- CREATE PERSON ---")

    name = input("Enter name: ")
    email = input("Enter email: ")
    phone = input("Enter phone: ")

    request = person_pb2.CreatePersonRequest(
        name=name,
        email=email,
        phone=phone
    )

    response = stub.CreatePerson(request)

    print("\nPerson created successfully!")

    print_person(response)


def get_person(stub):

    print("\n--- GET PERSON ---")

    person_id = int(
        input("Enter person ID: ")
    )

    request = person_pb2.GetPersonRequest(
        id=person_id
    )

    response = stub.GetPerson(request)

    print("\nPerson found:")

    print_person(response)


def update_person(stub):

    print("\n--- UPDATE PERSON ---")

    person_id = int(
        input("Enter person ID: ")
    )

    name = input("Enter new name: ")
    email = input("Enter new email: ")
    phone = input("Enter new phone: ")

    request = person_pb2.UpdatePersonRequest(
        id=person_id,
        name=name,
        email=email,
        phone=phone
    )

    response = stub.UpdatePerson(request)

    print("\nPerson updated successfully!")

    print_person(response)


def delete_person(stub):

    print("\n--- DELETE PERSON ---")

    person_id = int(
        input("Enter person ID: ")
    )

    request = person_pb2.DeletePersonRequest(
        id=person_id
    )

    response = stub.DeletePerson(request)

    print(response.message)


def print_person(person):

    print("---------------------")
    print("ID:", person.id)
    print("Name:", person.name)
    print("Email:", person.email)
    print("Phone:", person.phone)
    print("---------------------")


def main():

    channel = grpc.insecure_channel(
        "localhost:50051"
    )

    stub = person_pb2_grpc.PersonServiceStub(
        channel
    )

    while True:

        print("\n===== PERSON CRUD =====")
        print("1. Create Person")
        print("2. Get Person")
        print("3. Update Person")
        print("4. Delete Person")
        print("5. Exit")

        choice = input(
            "\nChoose an option: "
        )

        try:

            if choice == "1":
                create_person(stub)

            elif choice == "2":
                get_person(stub)

            elif choice == "3":
                update_person(stub)

            elif choice == "4":
                delete_person(stub)

            elif choice == "5":
                print("Goodbye!")
                break

            else:
                print(
                    "Invalid option. Please choose 1-5."
                )

        except ValueError:

            print(
                "Please enter a valid numeric ID."
            )

        except grpc.RpcError as error:

            print(
                "\ngRPC Error:",
                error.code().name
            )

            print(
                "Details:",
                error.details()
            )


if __name__ == "__main__":
    main()