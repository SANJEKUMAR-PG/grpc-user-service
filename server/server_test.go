package main

import (
	"context"
	"io"
	"log"
	"testing"

	pb "github.com/SANJEKUMAR-PG/grpc-user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var expectedUsers = []*pb.User{
	{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
	{Id: 2, Fname: "Sanjay", City: "India", Phone: 6363543646, Height: 5.8, Married: false},
}

func newServer(t *testing.T, register func(server *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	server := grpc.NewServer()
	t.Cleanup(func() {
		server.Stop()
	})

	register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("server.Serve %v", err)
		}
	}()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("localhost:9879", opts...)
	t.Cleanup(func() {
		conn.Close()
	})

	if err != nil {
		t.Fatalf("grpc.NewClient %v", err)
	}
	return conn
}

func TestUserServiceServer_GetUserById(t *testing.T) {
	service := UserServiceServer{}
	conn := newServer(t, func(server *grpc.Server) {
		pb.RegisterUserServiceServer(server, &service)
	})

	client := pb.NewUserServiceClient(conn)
	resp, err := client.GetUserById(context.Background(), &pb.UserIdRequest{Id: 1})
	if err != nil {
		t.Fatalf("client.GetUserById %v", err)
	}
	if resp.User.Id != 1 && resp.User.Fname != "Steve" {
		t.Fatalf("Unexpected values %v", resp.User)
	}
}

func TestUserServiceServer_GetUsersByIds(t *testing.T) {
	service := &UserServiceServer{}
	conn := newServer(t, func(server *grpc.Server) {
		pb.RegisterUserServiceServer(server, service)
	})
	t.Cleanup(func() {
		conn.Close()
	})

	client := pb.NewUserServiceClient(conn)
	req := &pb.UsersIdsRequest{Ids: []int32{1, 2}}
	stream, err := client.GetUsersByIds(context.Background(), req)
	if err != nil {
		t.Fatalf("client.GetUsersByIds: %v", err)
	}
	var receivedUsers []*pb.User
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream.Recv: %v", err)
		}

		receivedUsers = append(receivedUsers, resp.User)
	}
	if len(receivedUsers) != len(expectedUsers) {
		t.Fatalf("Expected %d users, got %d", len(expectedUsers), len(receivedUsers))
	}
	for i, expected := range expectedUsers {
		if receivedUsers[i].Id != expected.Id || receivedUsers[i].Fname != expected.Fname || receivedUsers[i].City != expected.City {
			t.Errorf("Unexpected user fields: %+v, expected: %+v", receivedUsers[i], expected)
		}
	}
}

func TestUserServiceServer_SearchUsers(t *testing.T) {
	service := &UserServiceServer{}
	conn := newServer(t, func(server *grpc.Server) {
		pb.RegisterUserServiceServer(server, service)
	})
	t.Cleanup(func() {
		conn.Close()
	})
	client := pb.NewUserServiceClient(conn)
	req := &pb.SearchRequest{City: "LA"}
	stream, err := client.SearchUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("client.SearchUsers: %v", err)
	}
	var receivedUsers []*pb.User
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream.Recv: %v", err)
		} 

		receivedUsers = append(receivedUsers, resp.User)
	}
	if len(receivedUsers) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(receivedUsers))
	}
	expectedUser := expectedUsers[0]
	if receivedUsers[0].Fname != expectedUser.Fname || receivedUsers[0].City != expectedUser.City || receivedUsers[0].Married != expectedUser.Married {
		t.Errorf("Unexpected user fields: %+v, expected: %+v", receivedUsers[0], expectedUser)
	}
}
