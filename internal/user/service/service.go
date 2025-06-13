package service

import (
	"context"
	"errors"
	"time"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {

	// Data validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "required: username, password and email!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed hash password")
	}

	// Create user
	user := &userpb.User{
		Id:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Role:      userpb.Role_USER,
		CreatedAt: timestamppb.New(time.Now()),
	}

	// Save in DB
	err = s.repo.CreateUser(ctx, user, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &userpb.RegisterUserResponse{User: user}, nil
}
