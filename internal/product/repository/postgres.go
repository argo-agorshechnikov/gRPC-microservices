package repository

import (
	"context"
	"fmt"
	"log"

	productpb "github.com/argo-agorshechnikov/gRPC-microservices/api/product-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	Pool *pgxpool.Pool
}

func (p *ProductRepository) CreateProduct(ctx context.Context, req *productpb.Product) (*productpb.Product, error) {

	var id int32
	err := p.Pool.QueryRow(
		ctx,
		"INSERT INTO products (productname, description, price) VALUES ($1, $2, $3) RETURNING id",
		req.ProductName, req.Description, req.Price).Scan(&id)
	if err != nil {
		return nil, err
	}

	product := &productpb.Product{
		Id:          id,
		ProductName: req.ProductName,
		Description: req.Description,
		Price:       req.Price,
	}
	return product, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, req *productpb.Product) (*productpb.Product, error) {

	query := `
		UPDATE products
		SET productname = $1, description = $2, price = $3
		WHERE id = $4
	`

	commTag, err := p.Pool.Exec(ctx, query,
		req.ProductName,
		req.Description,
		req.Price,
		req.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	if commTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no product found with id: %d", req.Id)
	}

	return req, nil

}

func (p *ProductRepository) ListProduct(ctx context.Context) ([]*productpb.Product, error) {

	rows, err := p.Pool.Query(ctx, "SELECT id, productname, description, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*productpb.Product, 0)
	for rows.Next() {
		var prod productpb.Product
		err := rows.Scan(&prod.Id, &prod.ProductName, &prod.Description, &prod.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &prod)
	}

	return products, nil
}

// func (p *ProductRepository) DeleteProduct(ctx context.Context, id int32) error {

// }

func CreateProductRepository(ctx context.Context, cfg *config.Config) (*ProductRepository, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed pool creating in product service: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed ping db in product service: %w", err)
	}

	log.Println("Successfully connected to db in product service")

	return &ProductRepository{Pool: pool}, nil
}
