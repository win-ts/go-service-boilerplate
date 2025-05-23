// Package di provides dependency injection for the server
package di

import (
	"database/sql"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"

	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/config"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/repository"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/service"
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

	// MySQL initialization
	mysqlDB, err := newMySQL(mySQLOptions{
		host:         c.MySQLConfig.Host,
		username:     c.MySQLConfig.Username,
		password:     c.MySQLConfig.Password,
		database:     c.MySQLConfig.Database,
		timeout:      c.MySQLConfig.Timeout,
		maxIdleConns: c.MySQLConfig.MaxIdleConns,
		maxOpenConns: c.MySQLConfig.MaxOpenConns,
		maxLifetime:  c.MySQLConfig.MaxLifetime,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to MySQL: %v", err)
	}
	defer func() {
		if err := mysqlDB.client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close MySQL connection", slog.Any("error", err))
		}
	}()

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
		Database: c.MySQLConfig.Database,
	}, repository.DatabaseRepositoryDependencies{
		Client: mysqlDB.client,
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
	}), mysqlDB.client, redisClient.client
}
