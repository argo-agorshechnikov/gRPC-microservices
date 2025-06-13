package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	Pool *pgxpool.Pool
}

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
