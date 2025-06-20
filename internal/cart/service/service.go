package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cartpb "github.com/argo-agorshechnikov/gRPC-microservices/api/cart-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/cart/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService struct {
	cartpb.UnimplementedCartServiceServer
	repo        *repository.CartRepository
	redisClient *redis.RedisClient
}

func NewCartService(repo *repository.CartRepository, redisClient *redis.RedisClient) *CartService {
	return &CartService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *CartService) GetCart(ctx context.Context, req *cartpb.GetCartRequest) (*cartpb.GetCartResponse, error) {

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cacheKey := fmt.Sprintf("cart:get:%d", req.UserId)

	cached, err := s.redisClient.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var items []*cartpb.CartItem
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			return &cartpb.GetCartResponse{Items: items}, nil
		}
	}

	// Get data from db
	items, err := s.repo.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	data, err := json.Marshal(items)
	if err == nil {
		s.redisClient.Rdb.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return &cartpb.GetCartResponse{
		Items: items,
	}, nil
}

func (s *CartService) AddItem(ctx context.Context, req *cartpb.AddItemRequest) (*cartpb.AddItemResponse, error) {

	if req.UserId == 0 || req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user and product id is required")
	}

	msg, err := s.repo.AddItem(ctx, req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item in cart: %v", err)
	}

	return &cartpb.AddItemResponse{
		Message: msg,
	}, nil

}

func (s *CartService) RemoveItem(ctx context.Context, req *cartpb.RemoveItemRequest) (*cartpb.RemoveItemResponse, error) {

	if req.UserId == 0 || req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user and product_id is required")
	}

	msg, err := s.repo.RemoveItem(ctx, req.UserId, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove item from cart: %v", err)
	}

	return &cartpb.RemoveItemResponse{
		Message: msg,
	}, nil
}
