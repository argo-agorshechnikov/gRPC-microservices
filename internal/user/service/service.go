package service

import (
	"context"
	"errors"
	"time"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var jwtKey = []byte("secret_key")

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {

	// Data validation
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "required: email and password")
	}

	// Get user from DB by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Check password with passwordHash from db
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Create data kit for token with role, user_id, exp
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"role":    user.Role.String(),

		// time to token live
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failde to generate token")
	}

	return &userpb.LoginResponse{
		Token: tokenString,
	}, nil

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
