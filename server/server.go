package main

import (
	"context"
	"log"
	"net"

	pb "github.com/SANJEKUMAR-PG/grpc-user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var users = []*pb.User{
	{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
	{Id: 2, Fname: "Sanjay", City: "India", Phone: 6363543646, Height: 5.8, Married: false},
}

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *pb.UserIdRequest) (*pb.UserResponse, error) {
	for _, user := range users {
		if user.Id == req.Id {
			return &pb.UserResponse{User: user}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "User with ID %d not found", req.Id)
}

func (s *UserServiceServer) GetUsersByIds(req *pb.UsersIdsRequest, stream pb.UserService_GetUsersByIdsServer) error {
	for _, id := range req.Ids {
		for _, user := range users {
			if user.Id == id {
				if err := stream.Send(&pb.UserResponse{User: user}); err != nil {
					return err // Return error if sending fails
				}
			}
		}
	}

	return nil
}

func (s *UserServiceServer) SearchUsers(req *pb.SearchRequest, stream pb.UserService_SearchUsersServer) error {
	for _, user := range users {
		if (req.City != "" && user.City != req.City) ||
			(req.Phone != 0 && user.Phone != req.Phone) ||
			(req.Married && !user.Married) {
			continue
		}
		if err := stream.Send(&pb.UserResponse{User: user}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServiceServer{})
	reflection.Register(grpcServer)
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
