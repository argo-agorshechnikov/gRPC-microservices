package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	userpb "github.com/argo-agorshechnikov/gRPC-microservices/api/user-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserExists = errors.New("user already exists")

type UserRepository struct {
	Pool *pgxpool.Pool
}

func (r *UserRepository) CreateUser(ctx context.Context, user *userpb.User, passwordHash string) error {

	query := `
		INSERT INTO users (id, name, email, role, created_at, password_hash)
		VALUES ($1, $2, $3, $4, NOW(), $5)
	`

	_, err := r.Pool.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Email,
		user.Role.String(),
		passwordHash,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserExists
		}
		return fmt.Errorf("failed to insert user in db: %w", err)
	}

	return nil
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
