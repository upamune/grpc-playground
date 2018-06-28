package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/upamune/grpc-playground/protobuf"
	"google.golang.org/grpc"
)

type GrpcPlaygroundService struct{}

func (s *GrpcPlaygroundService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	fmt.Println("in handler")
	return &pb.PingResponse{}, nil
}

func interceptor(char string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("interceptor%s\n", char)
		defer fmt.Printf("interceptor%s\n", char)
		return handler(ctx, req)
	}
}

func main() {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			interceptor("A"),
			interceptor("B"),
			interceptor("C"),
		),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGrpcPlaygroundServer(grpcServer, &GrpcPlaygroundService{})

	const port = ":8080"

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	go func() {
		if err := grpcServer.Serve(ln); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}()

	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	client := pb.NewGrpcPlaygroundClient(conn)
	_, err = client.Ping(context.TODO(), &pb.PingRequest{})
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	grpcServer.Stop()
}
