// Package di provides dependency injection for the consumer
package di

import (
	"context"
	"log"
	"log/slog"

	"github.com/getsentry/sentry-go"

	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/config"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/handler"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/repository"
	"github.com/win-ts/go-service-boilerplate/consumer/kafka-consumer/service"
)

// New injects the dependencies for the consumer
func New(c *config.Config) {
	ctx := context.Background()

	// Sentry initialization
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.SentryConfig.SentryDSN,
		Debug:            true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		slog.Error("error - [di.initSentry] sentry initialization failed", slog.Any("error", err))
	}

	// Kafka Consumer initialization
	kafkaConsumer, err := newKafkaConsumer(kafkaConsumerOptions{
		username:          c.KafkaConsumerConfig.Username,
		password:          c.KafkaConsumerConfig.Password,
		sessionTimeout:    c.KafkaConsumerConfig.SessionTimeout,
		heartbeatInterval: c.KafkaConsumerConfig.HeartbeatInterval,
		bufferSize:        c.KafkaConsumerConfig.BufferSize,
		maxRetry:          c.KafkaConsumerConfig.MaxRetry,
		brokers:           c.KafkaConsumerConfig.Brokers,
		group:             c.KafkaConsumerConfig.Group,
	})
	if err != nil {
		log.Panicf("error - [di.New] unable to create Kafka consumer: %v", err)
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
	defer func() {
		if err := redisClient.client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close Redis connection", slog.Any("error", err))
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

	// Service initialization
	service := service.New(service.Dependencies{
		ExampleRepository:     exampleRepo,
		WiremockAPIRepository: wiremockAPIRepo,
		DatabaseRepository:    databaseRepo,
		CacheRepository:       cacheRepo,
	})

	// Handler initialization
	processorHandler := handler.NewProcessor(handler.ProcessorDependencies{
		Service: service,
	})

	// Consume Messages
	slog.Info("[di.New] consuming messages...")
	for {
		if err := kafkaConsumer.consumer.Consume(ctx, []string{c.KafkaConsumerConfig.Topic}, processorHandler); err != nil {
			log.Panicf("error - [di.New] unable to consume messages: %v", err)
		}
	}
}
