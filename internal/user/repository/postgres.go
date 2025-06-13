package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func CreateUserRepository(ctx context.Context, cfg *config.Config) (*UserRepository, error) {

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
		return nil, fmt.Errorf("failed create pool in user service: %w", err)
	}

	// Check db by ping
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed ping db in user service: %w", err)
	}

	log.Println("Successfully connected to db user service")

	return &UserRepository{Pool: pool}, nil
}
