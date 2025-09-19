package main

import (
	"context"
	"log"
	"net"

	pb "simplegrpcserver/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		log.Fatal(err)
	}
	cred, error := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if error != nil {
		log.Fatal("could not load tls keys")
	}

	grpcServer := grpc.NewServer(grpc.Creds(cred))
	pb.RegisterCalculateServer(grpcServer, &server{})

	log.Printf("server running... on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
