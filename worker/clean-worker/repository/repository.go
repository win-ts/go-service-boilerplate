// Package repository provides the repository interfaces for the domain
package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/dto"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/pkg/httpclient"
)

// ExampleRepository represents the repository layer functions of example repository
type ExampleRepository interface {
	DoExample(ctx context.Context) (string, error)
}

// WiremockAPIRepository represents the repository layer functions of wiremock API repository
type WiremockAPIRepository interface {
	GetTest(ctx context.Context, h dto.WiremockGetTestHeader) (*httpclient.Response[dto.WiremockGetTestResponse], error)
}

// DatabaseRepository represents the repository layer functions of database repository
type DatabaseRepository interface {
	QueryTest() (*[]dto.TestEntity, error)
}

// CacheRepository represents the repository layer functions of cache repository
type CacheRepository interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
}
