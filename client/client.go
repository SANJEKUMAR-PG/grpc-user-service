package main

import (
	"context"
	"io"
	"log"

	pb "github.com/SANJEKUMAR-PG/grpc-user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("localhost:9879", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	getUserById(client)
	GetUsersByIds(client)
	SearchUsers(client)
}

func getUserById(client pb.UserServiceClient) {
	req := pb.UserIdRequest{Id: 1}
	resp, err := client.GetUserById(context.Background(), &req)
	if err != nil {
		log.Fatalf("error getting user by Id: %v", err)
	}
	log.Printf("User by Id: %+v", resp.User)
}

func GetUsersByIds(client pb.UserServiceClient) {
	req := pb.UsersIdsRequest{Ids: []int32{1, 2}}
	stream, err := client.GetUsersByIds(context.Background(), &req)
	if err != nil {
		log.Fatalf("error getting users by Ids: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error recieving user: %v", err)
		}
		log.Printf("User by Ids: %+v", resp.User)
	}
}

func SearchUsers(client pb.UserServiceClient) {
	req := pb.SearchRequest{City: "LA"}
	stream, err := client.SearchUsers(context.Background(), &req)
	if err != nil {
		log.Fatalf("error searching users: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving user: %v", err)
		}
		log.Fatalf("User found: %+v", resp.User)
	}
}
