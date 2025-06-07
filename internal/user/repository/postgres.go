package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, cfg *config.Config) (*PostgresRepository, error) {

	// string connection
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// pool db connections
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed create pool: %w", err)
	}

	// Check db by ping
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed ping db: %w", err)
	}

	log.Println("Successfully connected to db")

	return &PostgresRepository{Pool: pool}, nil
}
