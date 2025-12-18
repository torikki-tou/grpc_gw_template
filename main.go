package main

import (
	pb "buftest/gen/go"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/tmc/grpc-websocket-proxy/wsproxy"

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
		Handler: wsproxy.WebsocketProxy(mux),
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

func (h *handler) OneOf(ctx context.Context, req *pb.Empty) (*pb.OneOfResponse, error) {
	return &pb.OneOfResponse{
		Response: &pb.OneOfResponse_ResponseOne{
			ResponseOne: &pb.ResponseOne{Field: "test"},
		},
	}, nil
}

func (h *handler) Streaming(server grpc.BidiStreamingServer[pb.StreamingRequest, pb.StreamingResponse]) error {

	go func() {
		for {
			message, err := server.Recv()
			switch {
			case err == nil:
				fmt.Println(message.Field)
			case errors.Is(err, io.EOF):
				return
			default:
				log.Print(err)
			}
		}
	}()

	for range 100 {
		err := server.Send(&pb.StreamingResponse{
			Field: "test",
		})
		if err != nil {
			log.Print(err)
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
