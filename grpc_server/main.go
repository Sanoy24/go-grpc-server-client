package main

import (
	"context"
	"log"
	"net"

	pb "simplegrpcserver/proto/gen"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculateServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	Sum := req.A + req.B
	return &pb.AddResponse{
		Sum: Sum,
	}, nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err) // Changed log.Fatal to log.Fatalf for better formatting
	}

	// --- REMOVE TLS/Certificate Setup ---
	// The following lines are removed:
	// cred, error := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	// if error != nil {
	//     log.Fatal("could not load tls keys")
	// }

	// Initialize the gRPC server without credentials (insecure)
	// The grpc.Creds(cred) option is removed.
	grpcServer := grpc.NewServer() // No arguments needed for an insecure server

	pb.RegisterCalculateServer(grpcServer, &server{})

	log.Printf("server running... on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}
