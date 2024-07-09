# gRPC User Service

This repository contains a gRPC user service implemented in Go.

## Prerequisites
Before building and running the application, ensure you have the following installed on your system:
- Go 1.22.5 or higher
- Protobuf compiler (`protoc`)

## Install Dependencies and Generate Protobuf Code

To prepare your environment and generate necessary code, follow these steps:

```bash
# Clone the repository
git clone https://github.com/your-username/grpc-user-service.git
cd grpc-user-service

# Install dependencies
go mod download

# Generate protobuf code
protoc --go_out=. --go-grpc_out=. proto/user.proto

# Build the server
cd server
go build -o grpc-user-service .

# Run the server (replace with your actual command to run the server)
./grpc-user-service

# Run test
go test -v ./...

#Access gRPC Service
grpcurl -plaintext -proto proto/user.proto localhost:9879 user_service.UserService/GetUserById -d '{"id": 1}'
grpcurl -plaintext -proto proto/user.proto localhost:9879 user_service.UserService/GetUsersByIds -d '{"ids": [1, 2]}'
grpcurl -plaintext -proto proto/user.proto localhost:9879 user_service.UserService/SearchUsers -d '{"city": "LA"}'




