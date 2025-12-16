package main

import (
	pb "buftest/gen/go"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcAddr := "0.0.0.0:9090"

	grpcServer := grpc.NewServer()

	pb.RegisterPingServiceServer(grpcServer, &handler{})

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	mux := runtime.NewServeMux()

	pb.RegisterPingServiceHandlerFromEndpoint(context.Background(), mux, grpcAddr, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})

	gwServer := &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: mux,
	}

	fmt.Println("Starting server at 0.0.0.0:8090")
	log.Fatal(gwServer.ListenAndServe())
}

type handler struct {
	pb.UnimplementedPingServiceServer
}

func (h *handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Field: req.Field}, nil
}

func (h *handler) HeavyPing(ctx context.Context, req *pb.HeavyPingRequest) (*pb.HeavyPingResponse, error) {
	return &pb.HeavyPingResponse{List: req.List}, nil
}
