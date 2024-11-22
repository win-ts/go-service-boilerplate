// Package di provides dependency injection for the server
package di

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/handler"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/database"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/kafka"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/repository"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/service"
)

// New injects the dependencies for the server
func New(c *config.Config) {
	ctx := context.Background()

	// Sentry initialization
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.SentryConfig.SentryDSN,
		Debug:            true,
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
	mysqlDB, err := database.NewMySQL(database.MySQLOptions{
		Host:         c.MySQLConfig.Host,
		Username:     c.MySQLConfig.Username,
		Password:     c.MySQLConfig.Password,
		Database:     c.MySQLConfig.Database,
		Timeout:      c.MySQLConfig.Timeout,
		MaxIdleConns: c.MySQLConfig.MaxIdleConns,
		MaxOpenConns: c.MySQLConfig.MaxOpenConns,
		MaxLifetime:  c.MySQLConfig.MaxLifetime,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to connect to MySQL: %v", err)
	}
	defer func() {
		if err := mysqlDB.Client.Close(); err != nil {
			slog.Error("error - [main.New] unable to close MySQL connection", slog.Any("error", err))
		}
	}()

	// Kafka Producer initialization
	kafkaProducer, err := kafka.NewProducer(kafka.ProducerOptions{
		Username: c.KafkaProducerConfig.Username,
		Password: c.KafkaProducerConfig.Password,
		Brokers:  c.KafkaProducerConfig.Brokers,
		Timeout:  c.KafkaProducerConfig.Timeout,
		MaxRetry: c.KafkaProducerConfig.MaxRetry,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to create Kafka producer: %v", err)
	}
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
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
		Client: mysqlDB.Client,
	})

	kafkaProducerRepo := repository.NewKafkaProducerRepository(repository.KafkaProducerRepositoryConfig{
		TopicName: c.KafkaProducerConfig.Topic,
	}, repository.KafkaProducerRepositoryDependencies{
		Producer: kafkaProducer.Producer,
	})

	// Service initialization
	service := service.New(service.Dependencies{
		ExampleRepository:       exampleRepo,
		WiremockAPIRepository:   wiremockAPIRepo,
		DatabaseRepository:      databaseRepo,
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
