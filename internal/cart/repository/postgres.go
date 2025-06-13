package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	Pool *pgxpool.Pool
}

func CreateCartRepository(ctx context.Context, cfg *config.Config) (*CartRepository, error) {

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
		return nil, fmt.Errorf("failed pool created in cart service: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed ping db in cart service: %w", err)
	}

	log.Println("Successfully connected to db in cart service")

	return &CartRepository{Pool: pool}, nil
}
