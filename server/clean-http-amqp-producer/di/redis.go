package di

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

type redisOptions struct {
	host     string
	password string
	timeout  time.Duration
	maxRetry int
	poolSize int
}

func newRedis(opts redisOptions) (*redisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        opts.host,
		Password:    opts.password,
		ReadTimeout: opts.timeout,
		MaxRetries:  opts.maxRetry,
		PoolSize:    opts.poolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	slog.Info("[di.newRedis] redis database connected",
		slog.String("host", opts.host),
	)

	return &redisClient{rdb}, nil
}
