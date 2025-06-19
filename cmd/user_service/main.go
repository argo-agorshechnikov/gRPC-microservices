package main

import (
	"context"
	"log"
	"net"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/auth"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/kafka"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Load config at app starting
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed load config (user): %v", err)
	}

	// Create repository, use config
	repo, err := repository.CreateUserRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed connect db (user): %v", err)
	}
	defer repo.Pool.Close() // Close pool db connections

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen (user): %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor([]byte("secret_key"))),
	)

	userProducer := kafka.NewProducer("user-events")
	defer userProducer.Close()

	userService := service.NewUserService(repo, userProducer)

	userpb.RegisterUserServiceServer(s, userService)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve (user): %v", err)
	}
}
