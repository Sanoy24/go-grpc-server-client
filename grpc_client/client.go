package main

import (
	"context"
	"fmt"
	mainapipb "grpcclient/proto/gen"
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
	client := mainapipb.NewCalculateClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := mainapipb.AddRequest{
		A: 2,
		B: 5,
	}
	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatalln("err: ", err)
	}
	fmt.Println("-- response -- ", res)
}
