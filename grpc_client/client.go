package main

import (
	"context"
	"fmt"
	pb "grpcclient/proto/gen"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Did not connect: ", err)
	}
	defer conn.Close()
	client := pb.NewCalculateClient(conn)
	clientTwo := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := pb.AddRequest{
		A: 2,
		B: 5,
	}
	reqTwo := &pb.HelloRequest{
		Name: "Yonas",
	}
	resTwo, errTwo := clientTwo.Greet(ctx, reqTwo)
	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	if errTwo != nil {
		log.Fatalln("err: ", errTwo)
	}
	fmt.Println("-- response -- ", res.Sum)
	fmt.Println("---| response |---", resTwo.Message)
	state := conn.GetState()
	fmt.Println("--| state |--", state)
}
