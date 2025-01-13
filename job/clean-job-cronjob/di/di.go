// Package di provides dependency injection for the server
package di

import (
	"database/sql"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"

	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/config"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/repository"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/service"
)

// New injects the dependencies for the server
func New(c *config.Config) (service.Port, *sql.DB, *redis.Client) {
	// Sentry initialization
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:                c.SentryConfig.SentryDSN,
		Debug:              false,
		EnableTracing:      true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
	}); err != nil {
		slog.Error("error - [main.New] sentry initialization failed", slog.Any("error", err))
	}

	// HTTP Client initialization
	httpClientWiremock := httpclient.NewHTTPClient(httpclient.Options{
		MaxConns:                 c.WiremockAPIConfig.MaxConns,
		MaxRetry:                 c.WiremockAPIConfig.MaxRetry,
		Timeout:                  c.WiremockAPIConfig.Timeout,
		InsecureSkipVerify:       c.WiremockAPIConfig.InsecureSkipVerify,
		MaxTransactionsPerSecond: c.WiremockAPIConfig.MaxTransactionsPerSecond,
	})

	// PostgreSQL initialization
	postgresDB, err := newPostgreSQL(postgreSQLOptions{
		host:         c.PostgreSQLConfig.Host,
		username:     c.PostgreSQLConfig.Username,
		password:     c.PostgreSQLConfig.Password,
		database:     c.PostgreSQLConfig.Database,
		timeout:      c.PostgreSQLConfig.Timeout,
		maxIdleConns: c.PostgreSQLConfig.MaxIdleConns,
		maxOpenConns: c.PostgreSQLConfig.MaxOpenConns,
		maxLifetime:  c.PostgreSQLConfig.MaxLifetime,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to PostgreSQL: %v", err)
	}

	// Redis initialization
	redisClient, err := newRedis(redisOptions{
		host:     c.RedisConfig.Host,
		password: c.RedisConfig.Password,
		timeout:  c.RedisConfig.Timeout,
		maxRetry: c.RedisConfig.MaxRetry,
		poolSize: c.RedisConfig.PoolSize,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to Redis: %v", err)
	}

	// Repository initialization
	exampleRepo := repository.NewExampleRepository(repository.ExampleRepositoryConfig{})

	wiremockAPIRepo := repository.NewWiremockAPIRepository(repository.WiremockAPIRepositoryConfig{
		BaseURL: c.WiremockAPIConfig.BaseURL,
		Path:    c.WiremockAPIConfig.Path,
	}, repository.WiremockAPIRepositoryDependencies{
		Client: httpClientWiremock,
	})

	databaseRepo := repository.NewDatabaseRepository(repository.DatabaseRepositoryConfig{
		Database: c.PostgreSQLConfig.Database,
	}, repository.DatabaseRepositoryDependencies{
		Client: postgresDB.client,
	})

	cacheRepo := repository.NewCacheRepository(repository.CacheRepositoryConfig{}, repository.CacheRepositoryDependencies{
		Client: redisClient.client,
	})

	// Service initialization
	return service.New(service.Dependencies{
		ExampleRepository:     exampleRepo,
		WiremockAPIRepository: wiremockAPIRepo,
		DatabaseRepository:    databaseRepo,
		CacheRepository:       cacheRepo,
	}), postgresDB.client, redisClient.client
}
