package infrastructure

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisFileStorage struct {
	redisClient *redis.Client
}

func NewRedisFileStorage(redisAddr string) (*RedisFileStorage, error) {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisFileStorage{
		redisClient: client,
	}, nil
}

func (r *RedisFileStorage) Save(key string, data []byte) error {
	// Save the file data to Redis
	err := r.redisClient.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save file to Redis: %v", err)
	}

	return nil
}

func (r *RedisFileStorage) Get(key string) ([]byte, error) {
	// Save the file data to Redis
	res := r.redisClient.Get(context.Background(), key)
	if res.Err() != nil {
		return nil, fmt.Errorf("failed to get file from Redis: %v", res.Err())
	}
	return res.Bytes()
}
