package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "hello/helloworld"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func printSayHello(c pb.GreeterClient, HelloRequest *pb.HelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, HelloRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func printAllHello(c pb.GreeterClient, HelloRequest *pb.HelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.AllHello(ctx, HelloRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
}

func runInputHello(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.InputHello(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// var input string
	for i := 0; i < 3; i++ {
		// fmt.Scan(&input)
		// if input == "exit" {
		// 	break
		// }
		if err := stream.Send(&pb.HelloRequest{Name: "namae"}); err != nil {
			log.Fatalf("this is: %v", err)
		}
	}
	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	printSayHello(c, &pb.HelloRequest{Name: name})
	printAllHello(c, &pb.HelloRequest{Name: name})
	runInputHello(c)
}
