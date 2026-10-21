package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	personpb "example.com/grpc-person-crud/go/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)


var reader = bufio.NewReader(os.Stdin)


func readInput(prompt string) string {

	fmt.Print(prompt)

	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input)
}


func readID(prompt string) (int32, error) {

	input := readInput(prompt)

	id, err := strconv.Atoi(input)

	if err != nil {
		return 0, err
	}

	return int32(id), nil
}


func printPerson(person *personpb.Person) {

	fmt.Println("---------------------")
	fmt.Println("ID:", person.Id)
	fmt.Println("Name:", person.Name)
	fmt.Println("Email:", person.Email)
	fmt.Println("Phone:", person.Phone)
	fmt.Println("---------------------")
}


func createPerson(
	client personpb.PersonServiceClient,
) error {

	fmt.Println("\n--- CREATE PERSON ---")

	name := readInput("Enter name: ")
	email := readInput("Enter email: ")
	phone := readInput("Enter phone: ")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	person, err := client.CreatePerson(
		ctx,
		&personpb.CreatePersonRequest{
			Name:  name,
			Email: email,
			Phone: phone,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("\nPerson created successfully!")

	printPerson(person)

	return nil
}


func getPerson(
	client personpb.PersonServiceClient,
) error {

	fmt.Println("\n--- GET PERSON ---")

	id, err := readID("Enter person ID: ")

	if err != nil {
		fmt.Println("Please enter a valid numeric ID.")
		return nil
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	person, err := client.GetPerson(
		ctx,
		&personpb.GetPersonRequest{
			Id: id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("\nPerson found:")

	printPerson(person)

	return nil
}


func updatePerson(
	client personpb.PersonServiceClient,
) error {

	fmt.Println("\n--- UPDATE PERSON ---")

	id, err := readID("Enter person ID: ")

	if err != nil {
		fmt.Println("Please enter a valid numeric ID.")
		return nil
	}

	name := readInput("Enter new name: ")
	email := readInput("Enter new email: ")
	phone := readInput("Enter new phone: ")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	person, err := client.UpdatePerson(
		ctx,
		&personpb.UpdatePersonRequest{
			Id:    id,
			Name:  name,
			Email: email,
			Phone: phone,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("\nPerson updated successfully!")

	printPerson(person)

	return nil
}


func deletePerson(
	client personpb.PersonServiceClient,
) error {

	fmt.Println("\n--- DELETE PERSON ---")

	id, err := readID("Enter person ID: ")

	if err != nil {
		fmt.Println("Please enter a valid numeric ID.")
		return nil
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	response, err := client.DeletePerson(
		ctx,
		&personpb.DeletePersonRequest{
			Id: id,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println(response.Message)

	return nil
}


func handleError(err error) {

	if err == nil {
		return
	}

	s, ok := status.FromError(err)

	if ok {
		fmt.Println("\ngRPC Error:", s.Code())
		fmt.Println("Details:", s.Message())
	} else {
		fmt.Println("\nError:", err)
	}
}


func main() {

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatalf(
			"Failed to create connection: %v",
			err,
		)
	}

	defer conn.Close()

	client := personpb.NewPersonServiceClient(conn)

	for {

		fmt.Println("\n===== PERSON CRUD =====")
		fmt.Println("1. Create Person")
		fmt.Println("2. Get Person")
		fmt.Println("3. Update Person")
		fmt.Println("4. Delete Person")
		fmt.Println("5. Exit")

		choice := readInput(
			"\nChoose an option: ",
		)

		switch choice {

		case "1":
			handleError(
				createPerson(client),
			)

		case "2":
			handleError(
				getPerson(client),
			)

		case "3":
			handleError(
				updatePerson(client),
			)

		case "4":
			handleError(
				deletePerson(client),
			)

		case "5":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println(
				"Invalid option. Please choose 1-5.",
			)
		}
	}
}