package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "simplegrpcserver/proto/gen"
	chatpb "simplegrpcserver/proto/gen/chat"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculateServer
	pb.UnimplementedGreeterServer
	chatpb.UnimplementedChatServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	Sum := req.A + req.B
	return &pb.AddResponse{
		Sum: Sum,
	}, nil
}

func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.Name),
	}, nil
}

// Unary
func (s *server) SayHello(_ context.Context, msg *chatpb.Message) (*chatpb.Message, error) {
	return &chatpb.Message{
		Sender: "Server",
		Text:   "Hello, " + msg.Sender,
	}, nil
}

// Server streaming
func (s *server) ChatStream(req *chatpb.Message, stream chatpb.ChatService_ChatStreamServer) error {
	for i := 1; i <= 3; i++ {
		resp := &chatpb.Message{
			Sender: "Server",
			Text:   fmt.Sprintf("Message %d to %s", i, req.Sender),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

// Client streaming
func (s *server) SendStream(stream chatpb.ChatService_SendStreamServer) error {
	var count int
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&chatpb.Message{
				Sender: "Server",
				Text:   fmt.Sprintf("Received %d messages", count),
			})
		}
		if err != nil {
			return err
		}
		count++
		log.Printf("Received from %s: %s", msg.Sender, msg.Text)
	}
}

// Bidirectional streaming
func (s *server) FullChat(stream chatpb.ChatService_FullChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("%s says: %s", msg.Sender, msg.Text)
		resp := &chatpb.Message{
			Sender: "Server",
			Text:   "Got your message: " + msg.Text,
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
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
	pb.RegisterGreeterServer(grpcServer, &server{})
	chatpb.RegisterChatServiceServer(grpcServer, &server{})

	log.Printf("server running... on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("server failed to serve: %v", err)
	}
}
