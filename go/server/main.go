package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	personpb "example.com/grpc-person-crud/go/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// PersonServer implements our PersonService.
type PersonServer struct {
	personpb.UnimplementedPersonServiceServer

	mu       sync.Mutex
	persons  map[int32]*personpb.Person
	nextID   int32
}


// CreatePerson
func (s *PersonServer) CreatePerson(
	ctx context.Context,
	req *personpb.CreatePersonRequest,
) (*personpb.Person, error) {

	if req.Name == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"name is required",
		)
	}

	if req.Email == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"email is required",
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	person := &personpb.Person{
		Id:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	s.persons[s.nextID] = person

	fmt.Printf("Created person with ID %d\n", s.nextID)

	s.nextID++

	return person, nil
}


// GetPerson
func (s *PersonServer) GetPerson(
	ctx context.Context,
	req *personpb.GetPersonRequest,
) (*personpb.Person, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	person, exists := s.persons[req.Id]

	if !exists {
		return nil, status.Errorf(
			codes.NotFound,
			"person with ID %d not found",
			req.Id,
		)
	}

	fmt.Printf("Reading person with ID %d\n", req.Id)

	return person, nil
}


// UpdatePerson
func (s *PersonServer) UpdatePerson(
	ctx context.Context,
	req *personpb.UpdatePersonRequest,
) (*personpb.Person, error) {

	if req.Name == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"name is required",
		)
	}

	if req.Email == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"email is required",
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.persons[req.Id]

	if !exists {
		return nil, status.Errorf(
			codes.NotFound,
			"person with ID %d not found",
			req.Id,
		)
	}

	person := &personpb.Person{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	s.persons[req.Id] = person

	fmt.Printf("Updated person with ID %d\n", req.Id)

	return person, nil
}


// DeletePerson
func (s *PersonServer) DeletePerson(
	ctx context.Context,
	req *personpb.DeletePersonRequest,
) (*personpb.DeletePersonResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.persons[req.Id]

	if !exists {
		return nil, status.Errorf(
			codes.NotFound,
			"person with ID %d not found",
			req.Id,
		)
	}

	delete(s.persons, req.Id)

	fmt.Printf("Deleted person with ID %d\n", req.Id)

	return &personpb.DeletePersonResponse{
		Message: fmt.Sprintf(
			"Person %d deleted successfully",
			req.Id,
		),
	}, nil
}


func main() {

	listener, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		log.Fatalf(
			"Failed to listen: %v",
			err,
		)
	}

	grpcServer := grpc.NewServer()

	personServer := &PersonServer{
		persons: make(map[int32]*personpb.Person),
		nextID:  1,
	}

	personpb.RegisterPersonServiceServer(
		grpcServer,
		personServer,
	)

	fmt.Println(
		"Person gRPC Server started on port 50051",
	)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf(
			"Failed to serve: %v",
			err,
		)
	}
}