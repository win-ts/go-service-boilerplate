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
		Debug:              true,
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
