// Package di provides dependency injection for the server
package di

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/handler"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/repository"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-kafka-producer/service"
)

// New injects the dependencies for the server
func New(c *config.Config) {
	ctx := context.Background()

	// Sentry initialization
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.SentryConfig.SentryDSN,
		Debug:            true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		slog.Error("error - [main.initServer] sentry initialization failed", slog.Any("error", err))
	}

	// Echo server initialization
	e := echo.New()
	setupServer(ctx, e, c)

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
	defer func() {
		if err := redisClient.client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close Redis connection", slog.Any("error", err))
		}
	}()

	// Kafka Producer initialization
	kafkaProducer, err := newKafkaProducer(kafkaProducerOptions{
		username: c.KafkaProducerConfig.Username,
		password: c.KafkaProducerConfig.Password,
		brokers:  c.KafkaProducerConfig.Brokers,
		timeout:  c.KafkaProducerConfig.Timeout,
		maxRetry: c.KafkaProducerConfig.MaxRetry,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to create Kafka producer: %v", err)
	}
	defer func() {
		if err := kafkaProducer.producer.Close(); err != nil {
			slog.Error("error - [main.New] unable to close Kafka producer", slog.Any("error", err))
		}
	}()

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

	kafkaProducerRepo := repository.NewKafkaProducerRepository(repository.KafkaProducerRepositoryConfig{
		TopicName: c.KafkaProducerConfig.Topic,
	}, repository.KafkaProducerRepositoryDependencies{
		Producer: kafkaProducer.producer,
	})

	// Service initialization
	service := service.New(service.Dependencies{
		ExampleRepository:       exampleRepo,
		WiremockAPIRepository:   wiremockAPIRepo,
		DatabaseRepository:      databaseRepo,
		CacheRepository:         cacheRepo,
		KafkaProducerRepository: kafkaProducerRepo,
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
