package service

import (
	"context"
	"log"

	productpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/product/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) ListProduct(ctx context.Context, req *productpb.ListProductRequest) (*productpb.ListProductResponse, error) {
	products, err := s.repo.ListProduct(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	return &productpb.ListProductResponse{
		Products: products,
	}, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.Product, error) {

	role, ok := ctx.Value(auth.RoleKey).(string)
	log.Printf("role received: %s", role)
	if !ok || role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "(service) only admin can create products")
	}

	product := &productpb.Product{
		ProductName: req.ProductName,
		Description: req.Description,
		Price:       req.Price,
	}

	createdProduct, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.Product, error) {

	role, ok := ctx.Value(auth.RoleKey).(string)
	if !ok || role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "(service) only admin can update product!")
	}

	product := &productpb.Product{
		Id:          req.Id,
		ProductName: req.ProductName,
		Description: req.Description,
		Price:       req.Price,
	}

	updatedProduct, err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, status.Error(codes.Internal, "(service) failed to update product")
	}

	return updatedProduct, nil
}
