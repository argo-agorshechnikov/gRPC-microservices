package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/user/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/kafka"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var jwtKey = []byte("secret_key")

type UserRegisteredEvent struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserLoginEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo     *repository.UserRepository
	producer *kafka.Producer
}

func NewUserService(repo *repository.UserRepository, producer *kafka.Producer) *UserService {
	return &UserService{
		repo:     repo,
		producer: producer,
	}
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

	events := UserLoginEvent{
		Email: req.Email,
		Token: tokenString,
	}

	eventsBytes, err := json.Marshal(events)
	if err != nil {
		log.Printf("failed marshal user login events: %v", err)
	} else {
		err = s.producer.SendMessage([]byte(fmt.Sprintf("%s", req.Email)), eventsBytes)

		if err != nil {
			log.Printf("failed to send message user login event to kafka: %v", err)
		}
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
		Name:      req.Name,
		Email:     req.Email,
		Role:      userpb.Role_USER,
		CreatedAt: timestamppb.New(time.Now()),
	}

	log.Println(user)
	// Save in DB
	err = s.repo.CreateUser(ctx, user, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		log.Printf("create user err: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	event := UserRegisteredEvent{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.String(),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal user registered event: %v", err)
	} else {
		err = s.producer.SendMessage([]byte(fmt.Sprintf("%d", user.Id)), eventBytes)

		if err != nil {
			log.Printf("failed to send user registered event to kafka: %v", err)
		}
	}

	return &userpb.RegisterUserResponse{User: user}, nil
}
