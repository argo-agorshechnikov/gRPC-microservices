package main

import (
	"context"
	"log"
	"net"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"

	"google.golang.org/grpc"
)

type server struct {
	userpb.UnimplementedUserServiceServer
}

func main() {
	ctx := context.Background()

	// Load config at app starting
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed load config (user): %v", err)
	}

	// Create repository, use config
	repo, err := repository.UserRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("failed connect db (user): %v", err)
	}
	defer repo.Pool.Close() // Close pool db connections

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen (user): %v", err)
	}
	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve (user): %v", err)
	}
}
