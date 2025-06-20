package main

import (
	"context"
	"log"
	"net"

	cartpb "github.com/argo-agorshechnikov/gRPC-microservices/api/cart-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/cart/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/cart/service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/auth"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	redispkg "github.com/argo-agorshechnikov/gRPC-microservices/pkg/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed load config (cart): %v", err)
	}

	repo, err := repository.CreateCartRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed created repo (cart): %v", err)
	}
	defer repo.Pool.Close()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor([]byte("secret_key"))),
	)

	redisClient := redispkg.NewRedisClient("redis:6379", "", 0)

	cartService := service.NewCartService(repo, redisClient)

	cartpb.RegisterCartServiceServer(s, cartService)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failder serve: %v", err)
	}
}
