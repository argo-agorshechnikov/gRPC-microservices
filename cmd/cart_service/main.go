package main

import (
	"context"
	"log"
	"net"

	cartpb "github.com/argo-agorshechnikov/gRPC-microservices/api/cart-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/cart/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"google.golang.org/grpc"
)

type server struct {
	cartpb.UnimplementedCartServiceServer
}

func main() {

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed load config (cart): %v", err)
	}

	repo, err := repository.CartRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed created repo (cart): %v", err)
	}
	defer repo.Pool.Close()

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
