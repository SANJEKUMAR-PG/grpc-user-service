 
# gRPC User Service with Search

This project implements a gRPC service for managing user details with search functionality. It allows fetching user details by ID, retrieving users by a list of IDs, and searching for users based on specific criteria.

## Prerequisites

Before running this application, ensure you have the following installed:
- Go 1.22.5 or later
- Docker (optional, for containerization)

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/SANJEKUMAR-PG/grpc-user-service.git
   cd grpc-user-service


go mod download

docker build -t grpc-user-service .

docker run -p 9879:9879 grpc-user-service

go run server/server.go

go test -v ./...
