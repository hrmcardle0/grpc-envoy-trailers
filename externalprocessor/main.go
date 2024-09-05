package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

// Server implements the ExternalProcessor service
type server struct {
	extprocv3.UnimplementedExternalProcessorServer
}

// Get value of a header based on it's key
func GetHeaderValue(headers []*core.HeaderValue, key string) (string, bool) {
	result := ""
	ok := false
	for _, headerValue := range headers {
		if headerValue.Key == key {
			result = headerValue.Value
			ok = true
			break
		}
	}
	return result, ok
}

// Implement our Process method that envoy envokes
func (s *server) Process(stream extprocv3.ExternalProcessor_ProcessServer) error {
	ctx := stream.Context()
	fmt.Println("Processing stream...")

	for {
		log.Println("Waiting for message...")
		select {
			case <-ctx.Done():
				log.Println("Context cancelled")
				return ctx.Err()
			default:
		}
		// Receive messages from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream")
			return nil
		}

		if err != nil {
			if status.Code(err) == codes.Canceled {
				log.Println("Stream cancelled")
				return nil
			}
			log.Printf("Error receiving message: %v", err)
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		// Do some work
		fmt.Println("Message Recievved")

		resp := &extprocv3.ProcessingResponse{}
		// Switch based on request type
		switch value := req.Request.(type) {
			case *extprocv3.ProcessingRequest_RequestHeaders:
				httpMethod, _ := GetHeaderValue(value.RequestHeaders.Headers.Headers, ":method")
				requestPath, _ := GetHeaderValue(value.RequestHeaders.Headers.Headers, ":path")
				// print out all headers
				for _, header := range value.RequestHeaders.Headers.Headers {
					log.Printf("%s\n", fmt.Sprintf("REQUEST Header: %s: %s", header.Key, header.Value))
				}
				log.Printf("%s\n", fmt.Sprintf("Handle (REQ_HEAD): downstream -> ext_proc -> upstream, Method:%s, Path:%s", httpMethod, requestPath))
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_RequestHeaders{},
				}

			case *extprocv3.ProcessingRequest_RequestTrailers:
				// print out all trailers
				for _, trailer := range value.RequestTrailers.Trailers.Headers{
					log.Printf("%s\n", fmt.Sprintf("REQUEST Trailer: %s: %s", trailer.Key, trailer.Value))
				}
				log.Printf("%s\n", "Handle (REQ_TRAILERS): downstream -> ext_proc -> upstream")
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_RequestTrailers{},
				}

			case *extprocv3.ProcessingRequest_RequestBody:
				log.Printf("%s\n", "Handle (REQ_BODY): downstream -> ext_proc -> upstream")
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_RequestBody{},
				}

			case *extprocv3.ProcessingRequest_ResponseHeaders:
				status, _ := GetHeaderValue(value.ResponseHeaders.Headers.Headers, ":status")
				// print out all headers
				for _, header := range value.ResponseHeaders.Headers.Headers {
					log.Printf("%s\n", fmt.Sprintf("RESPONSE Header: %s: %s", header.Key, header.Value))
				}
				log.Printf("%s\n", fmt.Sprintf("Handle (REQ_HEAD): upstream -> ext_proc -> downstream, status:%v", status))
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_ResponseHeaders{},
				}

			case *extprocv3.ProcessingRequest_ResponseTrailers:
				// print out all trailers
				for _, trailer := range value.ResponseTrailers.Trailers.Headers{
					log.Printf("%s\n", fmt.Sprintf("RESPONSE Trailer: %s: %s", trailer.Key, trailer.Value))
				}
				log.Printf("%s\n", "Handle (REQ_TRAILERS): upstream -> ext_proc -> downstream")
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_ResponseTrailers{},
				}

			case *extprocv3.ProcessingRequest_ResponseBody:
				log.Printf("%s\n", "Handle (REQ_BODY): upstream -> ext_proc -> downstream")
				resp = &extprocv3.ProcessingResponse{
					Response: &extprocv3.ProcessingResponse_ResponseBody{},
				}

			default:
				log.Printf("%s\n", fmt.Sprintf("Unknown Request type %v\n", value))
		
		}
		
		// Send response
		if err := stream.Send(resp); err != nil {
			log.Printf("Error sending response: %v", err)
		}
		log.Println("Response sent")
	}
}

func main() {
	// Set up a listener on port 50051 for the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register the external processor service
	extprocv3.RegisterExternalProcessorServer(s, &server{})
	log.Println("Starting External Processor Service on :50051")

	// Start the gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
