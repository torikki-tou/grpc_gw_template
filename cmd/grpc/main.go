package main

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"log"
	"net"

	"google.golang.org/grpc"
	"grpc_gw_template/internal/services"
	"grpc_gw_template/protogen/golang/users"
)

func main() {
	const addr = "localhost:50051"

	// create a TCP listener on the specified port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server instance
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()))

	// create a order service instance with a reference to the db
	orderService := services.NewUsersService()

	// register the order service with the grpc server
	users.RegisterUsersServer(server, orderService)

	// start listening to requests
	log.Printf("server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
