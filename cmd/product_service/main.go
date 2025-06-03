package main

import (
	"log"
	"net"

	prodpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"google.golang.org/grpc"
)

type server struct {
	prodpb.UnimplementedProductServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	s := grpc.NewServer()
	prodpb.RegisterProductServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed serve: %v", err)
	}
}
