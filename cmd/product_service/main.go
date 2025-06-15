package main

import (
	"context"
	"log"
	"net"

	prodpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/product/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/product/service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/auth"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	repo, err := repository.CreateProductRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed repo created (product): %v", err)
	}
	defer repo.Pool.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor([]byte("secret_key"))),
	)
	productService := service.NewProductService(repo)
	prodpb.RegisterProductServiceServer(s, productService)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed serve (product): %v", err)
	}
}
