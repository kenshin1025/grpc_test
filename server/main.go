package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "hello/helloworld"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) AllHello(in *pb.HelloRequest, stream pb.Greeter_AllHelloServer) error {
	for i := 0; i < 4; i++ {
		if err := stream.Send(&pb.HelloReply{Message: fmt.Sprint(i) + "Hello " + in.GetName()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) InputHello(stream pb.Greeter_InputHelloServer) error {
	var names string
	for {
		HelloRequest, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "Hello " + names})
		}
		if err != nil {
			return err
		}
		names = names + HelloRequest.GetName()
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	fmt.Println("server start")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
