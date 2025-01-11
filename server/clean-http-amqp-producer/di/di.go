// Package di provides dependency injection for the server
package di

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/handler"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/repository"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/service"
)

// New injects the dependencies for the server
func New(c *config.Config) {
	ctx := context.Background()

	// Sentry initialization
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:                c.SentryConfig.SentryDSN,
		Debug:              true,
		EnableTracing:      true,
		TracesSampleRate:   1.0,
		ProfilesSampleRate: 1.0,
	}); err != nil {
		slog.Error("error - [main.New] sentry initialization failed", slog.Any("error", err))
	}

	// Echo server initialization
	e := echo.New()
	setupServer(ctx, e, c)

	// AMQP Producer initialization
	amqpProducer, err := newAMQPProducer(amqpProducerOptions{
		host:            c.AMQPProducerConfig.Host,
		username:        c.AMQPProducerConfig.Username,
		password:        c.AMQPProducerConfig.Password,
		queueName:       c.AMQPProducerConfig.QueueName,
		queueDurable:    c.AMQPProducerConfig.QueueDurable,
		queueAutoDelete: c.AMQPProducerConfig.QueueAutoDelete,
		queueExclusive:  c.AMQPProducerConfig.QueueExclusive,
		queueNoWait:     c.AMQPProducerConfig.QueueNoWait,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to AMQP Producer: %v", err)
	}
	defer func() {
		if err := amqpProducer.channel.Close(); err != nil {
			slog.Error("error - [main.New] unable to close AMQP Producer connection", slog.Any("error", err))
		}
	}()

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
	defer func() {
		if err := postgresDB.client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close PostgreSQL connection", slog.Any("error", err))
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
	defer func() {
		if err := redisClient.client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close Redis connection", slog.Any("error", err))
		}
	}()

	// Repository initialization
	amqpProducerRepo := repository.NewAMQPProducerRepository(repository.AMQPProducerRepositoryConfig{
		QueueName: c.AMQPProducerConfig.QueueName,
	}, repository.AMQPProducerRepositoryDependencies{
		Channel: amqpProducer.channel,
	})

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
	service := service.New(service.Dependencies{
		AMQPProducerRepository: amqpProducerRepo,
		ExampleRepository:      exampleRepo,
		WiremockAPIRepository:  wiremockAPIRepo,
		DatabaseRepository:     databaseRepo,
		CacheRepository:        cacheRepo,
	})

	// Handler initialization
	handler.New(e, handler.Dependencies{
		Service: service,
	})

	// HTTP Listening
	if err := e.Start(":" + c.AppConfig.Port); err != nil && err != http.ErrServerClosed {
		log.Panicf("error - [main.New] unable to start server: %v", err)
	}
}
