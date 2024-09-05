package main

import (
	"log"
	"net"
	"context"


	"github.com/hrmcardle0/envoy-grpc-trailer/pb"
	"google.golang.org/grpc"
	
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// Run simple gRPC server
func main() {
	log.Println("Running server")

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &GreeterServer{})

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}


}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("Received request with name %s\n", req.GetName())
	return &pb.Response{Message: "Hello " + req.GetName()}, nil
}