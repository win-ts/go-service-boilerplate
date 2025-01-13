// Package di provides dependency injection for the server
package di

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsentry/sentry-go"

	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/config"
	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/repository"
	"github.com/win-ts/go-service-boilerplate/consumer/amqt-consumer/service"
)

// New injects the dependencies for the server
func New(c *config.Config) {
	ctx := context.Background()

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

	// AMQP Consumer initialization
	amqpConsumer, err := newAMQPConsumer(amqpConsumerOptions{
		host:              c.AMQPConsumerConfig.Host,
		username:          c.AMQPConsumerConfig.Username,
		password:          c.AMQPConsumerConfig.Password,
		queueName:         c.AMQPConsumerConfig.QueueName,
		queueDurable:      c.AMQPConsumerConfig.QueueDurable,
		queueAutoDelete:   c.AMQPConsumerConfig.QueueAutoDelete,
		queueExclusive:    c.AMQPConsumerConfig.QueueExclusive,
		queueNoWait:       c.AMQPConsumerConfig.QueueNoWait,
		consumerName:      c.AMQPConsumerConfig.ConsumerName,
		consumerAutoAck:   c.AMQPConsumerConfig.ConsumerAutoAck,
		consumerExclusive: c.AMQPConsumerConfig.ConsumerExclusive,
		consumerNoLocal:   c.AMQPConsumerConfig.ConsumerNoLocal,
		consumerNoWait:    c.AMQPConsumerConfig.ConsumerNoWait,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to AMQP as consumer: %v", err)
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
		ExampleRepository:     exampleRepo,
		WiremockAPIRepository: wiremockAPIRepo,
		DatabaseRepository:    databaseRepo,
		CacheRepository:       cacheRepo,
	})

	_ = ctx
	_ = service

	// Start AMQP Consumer
	{
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Consume Messages
		go func() {
			for d := range amqpConsumer.deliveryChan {
				slog.Info("[main] received message", slog.Any("message", string(d.Body)))

				if err := d.Acknowledger.Ack(d.DeliveryTag, false); err != nil {
					slog.Error("[main] unable to acknowledge message", slog.Any("error", err))
				}
			}
		}()

		slog.Info("[main] waiting for messages...")

		<-sigChan
		slog.Info("[main]: gracefully shutting down consumer...")
	}
}
