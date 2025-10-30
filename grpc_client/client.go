package main

import (
	"context"
	"fmt"
	pb "grpcclient/proto/gen"
	chatpb "grpcclient/proto/gen/chat"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculateClient(conn)
	clientTwo := pb.NewGreeterClient(conn)
	chatClient := chatpb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := &pb.AddRequest{A: 2, B: 5}
	reqTwo := &pb.HelloRequest{Name: "Yonas"}

	// --- Unary examples ---
	res, err := client.Add(ctx, req)
	if err != nil {
		log.Fatalf("Add error: %v", err)
	}

	resTwo, err := clientTwo.Greet(ctx, reqTwo)
	if err != nil {
		log.Fatalf("Greet error: %v", err)
	}

	fmt.Println("-- response -- ", res.Sum)
	fmt.Println("---| response |---", resTwo.Message)

	// --- 1️⃣ Unary call ---
	resp, err := chatClient.SayHello(context.Background(), &chatpb.Message{Sender: "Yonas"})
	if err != nil {
		log.Fatalf("SayHello error: %v", err)
	}
	log.Println("Unary:", resp.Text)

	// --- 2️⃣ Server streaming ---
	stream, err := chatClient.ChatStream(context.Background(), &chatpb.Message{Sender: "Yonas"})
	if err != nil {
		log.Fatalf("ChatStream error: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Server stream error: %v", err)
			break
		}
		log.Println("Server stream:", msg.Text)
	}

	// --- 3️⃣ Client streaming ---
	cs, err := chatClient.SendStream(context.Background())
	if err != nil {
		log.Fatalf("SendStream error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := cs.Send(&chatpb.Message{Sender: "Yonas", Text: "Ping"}); err != nil {
			log.Printf("Send error: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	reply, err := cs.CloseAndRecv()
	if err != nil {
		log.Printf("CloseAndRecv error: %v", err)
	}
	log.Println("Client stream:", reply.GetText())

	// --- 4️⃣ Bidirectional streaming ---
	bs, err := chatClient.FullChat(context.Background())
	if err != nil {
		log.Fatalf("FullChat error: %v", err)
	}

	go func() {
		for i := 1; i <= 3; i++ {
			bs.Send(&chatpb.Message{Sender: "Yonas", Text: fmt.Sprintf("Message %d", i)})
			time.Sleep(time.Second)
		}
		bs.CloseSend()
	}()
	for {
		msg, err := bs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("BiDi stream error: %v", err)
			break
		}
		log.Println("BiDi stream:", msg.Text)
	}

	state := conn.GetState()
	fmt.Println("--| connection state |--", state)
}
