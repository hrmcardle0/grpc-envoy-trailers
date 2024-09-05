package main

import (
	"context"
	"log"

	"github.com/hrmcardle0/envoy-grpc-trailer/pb"
	"google.golang.org/grpc"
)

var address = "grpc-server-service:80"

func main() {
	log.Println("Running client")

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	
	resp, err := client.SayHello(context.Background(), &pb.Request{Name: "world"})
	if err != nil {
		log.Fatalf("Failed to call SayHello: %v", err)
	}

	log.Printf("Response: %s", resp.GetMessage())
}
