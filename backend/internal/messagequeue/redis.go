package messagequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// RedisClient represents a Redis client
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(config RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &RedisClient{client: client}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Publish publishes a message to a channel
func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.client.Publish(ctx, channel, payload).Err()
}

// Subscribe subscribes to a channel and returns a channel for messages
func (r *RedisClient) Subscribe(ctx context.Context, channel string) (<-chan *redis.Message, error) {
	pubsub := r.client.Subscribe(ctx, channel)

	// Test the connection
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to channel: %w", err)
	}

	return pubsub.Channel(), nil
}

// Set sets a key-value pair with optional expiration
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, payload, expiration).Err()
}

// Get gets a value by key
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete deletes a key
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// LPush pushes a value to the head of a list
func (r *RedisClient) LPush(ctx context.Context, key string, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.LPush(ctx, key, payload).Err()
}

// RPop pops a value from the tail of a list
func (r *RedisClient) RPop(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.RPop(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// BRPop blocks and pops a value from the tail of a list
func (r *RedisClient) BRPop(ctx context.Context, timeout time.Duration, key string, dest interface{}) error {
	result, err := r.client.BRPop(ctx, timeout, key).Result()
	if err != nil {
		return err
	}

	if len(result) < 2 {
		return fmt.Errorf("unexpected result length: %d", len(result))
	}

	return json.Unmarshal([]byte(result[1]), dest)
}
