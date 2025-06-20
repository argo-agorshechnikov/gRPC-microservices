package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Rdb *redis.Client
}

// Created and returned connection to Redis

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed connection to Redis: %v", err)
	}

	log.Println("Successfully redis connection!")
	return &RedisClient{Rdb: rdb}
}
