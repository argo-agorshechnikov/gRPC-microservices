package service

import (
	"context"
	"encoding/json"
	"log"

	productpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/internal/product/repository"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/auth"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/kafka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductListEvent struct {
	Id          int32   `json:"id"`
	ProductName string  `json:"productName"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	repo     *repository.ProductRepository
	producer *kafka.Producer
}

func NewProductService(repo *repository.ProductRepository, producer *kafka.Producer) *ProductService {
	return &ProductService{
		repo:     repo,
		producer: producer,
	}
}

func (s *ProductService) ListProduct(ctx context.Context, req *productpb.ListProductRequest) (*productpb.ListProductResponse, error) {
	products, err := s.repo.ListProduct(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	for _, product := range products {
		event := ProductListEvent{
			Id:          product.Id,
			ProductName: product.ProductName,
			Description: product.Description,
			Price:       product.Price,
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed marshal product list event: %v", err)
		} else {
			err = s.producer.SendMessage([]byte(event.ProductName), eventBytes)

			if err != nil {
				log.Printf("failed send message product list event to kafka: %v", err)
			}
		}

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

	event := ProductListEvent{
		ProductName: createdProduct.ProductName,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed marshal product create event: %v", err)
	} else {
		err = s.producer.SendMessage([]byte(event.ProductName), eventBytes)
		if err != nil {
			log.Printf("failed send message product create event to kafka:%v", err)
		}
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

func (s *ProductService) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*emptypb.Empty, error) {

	role, ok := ctx.Value(auth.RoleKey).(string)
	if !ok || role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "(service) only admin can delete product!")
	}

	err := s.repo.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "(service) failed to delete product")
	}

	return &emptypb.Empty{}, nil
}
