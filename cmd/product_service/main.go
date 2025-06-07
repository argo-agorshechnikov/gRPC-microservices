package main

import (
	"context"
	"log"
	"net"

	prodpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/product/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"google.golang.org/grpc"
)

type server struct {
	prodpb.UnimplementedProductServiceServer
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed load config (product): %v", err)
	}

	repo, err := repository.ProductRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed repo created (product): %v", err)
	}
	defer repo.Pool.Close()

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
