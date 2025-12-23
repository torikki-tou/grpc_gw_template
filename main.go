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
	"strings"
	"time"

	"github.com/tmc/grpc-websocket-proxy/wsproxy"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	grpcAddr := "0.0.0.0:9090"

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptor), grpc.ChainStreamInterceptor(streamInterceptor))

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

	headerMatcher := func(key string) (string, bool) {
		print("header")
		switch strings.ToLower(key) {
		case "custom-header":
			return "custom-header", true
		default:
			return runtime.DefaultHeaderMatcher(key)
		}
	}

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithOutgoingHeaderMatcher(headerMatcher),
	)
	
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterPingServiceHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts); err != nil {
		log.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: wsproxy.WebsocketProxy(mux, wsproxy.WithForwardedHeaders(func(header string) bool { return true })),
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

func (h *handler) OneOf(ctx context.Context, req *pb.OneOfRequest) (*pb.OneOfResponse, error) {
	return &pb.OneOfResponse{
		Response: &pb.OneOfResponse_ResponseOne{
			ResponseOne: &pb.ResponseOne{Field: "test"},
		},
	}, nil
}

func (h *handler) Streaming(server grpc.BidiStreamingServer[pb.StreamingRequest, pb.StreamingResponse]) error {

	go func() {
		for {
			_, err := server.Recv()
			switch {
			case err == nil:
				continue
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

func unaryInterceptor(ctx context.Context, req any, serverInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata")
	}

	fmt.Printf("method: %s, custom unary header: %v\n", serverInfo.FullMethod, md.Get("custom-header"))

	return handler(ctx, req)
}

func streamInterceptor(srv any, ss grpc.ServerStream, serverInfo *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errors.New("no metadata")
	}

	fmt.Printf("method: %s, custom stream header: %v\n", serverInfo.FullMethod, md.Get("custom-header"))

	return handler(srv, &wrappedStream{ServerStream: ss})
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}
	r, _ := m.(*pb.StreamingRequest)
	fmt.Printf("RecvMsg: %v\n", r.Field)
	return nil
}

func (w *wrappedStream) SendMsg(m any) error {
	r, _ := m.(*pb.StreamingResponse)
	fmt.Printf("SendMsg: %v\n", r.Field)
	err := w.ServerStream.SendMsg(m)
	if err != nil {
		return err
	}
	return nil
}
