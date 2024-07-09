FROM golang:1.22.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o grpc-user-service ./server/server.go

EXPOSE 9879

CMD ["./grpc-user-service"]