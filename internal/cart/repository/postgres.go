package repository

import (
	"context"
	"fmt"
	"log"

	cartpb "github.com/argo-agorshechnikov/gRPC-microservices/api/cart-service"
	"github.com/argo-agorshechnikov/gRPC-microservices/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	Pool *pgxpool.Pool
}

func (r *CartRepository) GetCart(ctx context.Context, user_id int32) ([]*cartpb.CartItem, error) {

	query := `
		SELECT product_id, quantity FROM cart_items WHERE user_id=$1 
	`

	rows, err := r.Pool.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*cartpb.CartItem

	for rows.Next() {
		var product_id int32
		var quantity int32

		if err := rows.Scan(&product_id, &quantity); err != nil {
			return nil, err
		}

		items = append(items, &cartpb.CartItem{
			ProductId: product_id,
			Quantity:  quantity,
		})
	}

	return items, nil
}

func (r *CartRepository) AddItem(ctx context.Context, user_id, product_id, quantity int32) (string, error) {

	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + $3
	`

	_, err := r.Pool.Query(ctx, query, user_id, product_id, quantity)
	if err != nil {
		return "", err
	}

	return "success add item in cart", nil
}

func (r *CartRepository) RemoveItem(ctx context.Context, user_id, product_id int32) (string, error) {
	query := `
		DELETE FROM cart_items WHERE user_id=$1 AND product_id=$2
	`
	_, err := r.Pool.Exec(ctx, query, user_id, product_id)
	if err != nil {
		return "", nil
	}

	return "success remove item from cart", nil
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
