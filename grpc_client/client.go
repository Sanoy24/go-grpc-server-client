package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mainapipb "grpcclient/proto/gen"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	cert, err := os.ReadFile("cert.pem")
	if err != nil {
		log.Fatal("failed to read the cert.pem: ", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatal("failed to append cert to pool")
	}
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	cred := credentials.NewTLS(tlsConfig)
	if err != nil {
		log.Fatal("err: ", err)
	}
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(cred))
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
