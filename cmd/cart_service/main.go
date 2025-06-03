package main

import (
	"log"
	"net"

	cartpb "github.com/argo-agorshechnikov/gRPC-microservices/api/cart-service"
	"google.golang.org/grpc"
)

type server struct {
	cartpb.UnimplementedCartServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}

	s := grpc.NewServer()
	cartpb.RegisterCartServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failder serve: %v", err)
	}
}
