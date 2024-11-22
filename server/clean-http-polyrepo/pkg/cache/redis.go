// Package cache provides the cache repository for the server
package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisOptions represents the configuration of the redis client
type RedisOptions struct {
	Host     string
	Password string
	Timeout  time.Duration
	MaxRetry int
	PoolSize int
}

// RedisClient represents the redis client
type RedisClient struct {
	Client *redis.Client
}

// NewRedis creates a new redis client
func NewRedis(opts RedisOptions) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        opts.Host,
		Password:    opts.Password,
		ReadTimeout: opts.Timeout,
		MaxRetries:  opts.MaxRetry,
		PoolSize:    opts.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	slog.Info("[pkg.NewRedis] redis database connected",
		slog.String("host", opts.Host),
	)

	return &RedisClient{rdb}, nil
}

// Ping pings the redis client
func (r *RedisClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := r.Client.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

// Close closes the redis client
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
