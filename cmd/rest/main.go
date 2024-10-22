package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_gw_template/protogen/golang/users"
)

func main() {
	const (
		addr             = "0.0.0.0:8080"
		orderServiceAddr = "localhost:50051"
	)

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := users.RegisterUsersHandlerFromEndpoint(context.Background(), mux, orderServiceAddr, opts); err != nil {
		log.Fatalf("failed to register the order server: %v", err)
	}

	// start listening to requests from the gateway server
	fmt.Println("API gateway server is running on " + addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}
}
