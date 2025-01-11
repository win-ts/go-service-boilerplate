// Package di provides dependency injection for the server
package di

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/handler"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/pkg/httpclient"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/repository"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/service"
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

	// GRPC Auth Client initialization
	grpcClient, err := newGRPCClient(grpcClientOptions{
		url: c.GRPCAuthConfig.URL,
	})
	if err != nil {
		log.Panicf("error - [main.New] unable to create grpc client: %v", err)
	}

	// HTTP Client initialization
	httpClientWiremock := httpclient.NewHTTPClient(httpclient.Options{
		MaxConns:                 c.WiremockAPIConfig.MaxConns,
		MaxRetry:                 c.WiremockAPIConfig.MaxRetry,
		Timeout:                  c.WiremockAPIConfig.Timeout,
		InsecureSkipVerify:       c.WiremockAPIConfig.InsecureSkipVerify,
		MaxTransactionsPerSecond: c.WiremockAPIConfig.MaxTransactionsPerSecond,
	})

	// Repository initialization
	authRepo := repository.NewAuthMiddlewareRepository(repository.AuthMiddlewareRepositoryConfig{}, repository.AuthMiddlewareRepositoryDependencies{
		GrpcClient: grpcClient.client,
	})

	exampleRepo := repository.NewExampleRepository(repository.ExampleRepositoryConfig{})

	wiremockAPIRepo := repository.NewWiremockAPIRepository(repository.WiremockAPIRepositoryConfig{
		BaseURL: c.WiremockAPIConfig.BaseURL,
		Path:    c.WiremockAPIConfig.Path,
	}, repository.WiremockAPIRepositoryDependencies{
		Client: httpClientWiremock,
	})

	// Service initialization
	service := service.New(service.Dependencies{
		AuthMiddlewareRepository: authRepo,
		ExampleRepository:        exampleRepo,
		WiremockAPIRepository:    wiremockAPIRepo,
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
