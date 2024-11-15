// Package di provides dependency injection for the server
package di

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/handler"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/repository"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/service"
)

// Config represents the configuration of the service
type Config struct {
	AppConfig               config.AppConfig
	SentryConfig            config.SentryConfig
	ExampleRepositoryConfig repository.ExampleRepositoryConfig
}

// New injects the dependencies for the server
func New(c Config) {
	ctx := context.Background()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:   c.SentryConfig.SentryDSN,
		Debug: true,
	}); err != nil {
		log.Errorf("error - [main.initServer] sentry initialization failed: %v", err)
	}

	e := echo.New()
	setupServer(ctx, e, c)

	exampleRepo := repository.NewExampleRepository(c.ExampleRepositoryConfig)

	service := service.New(service.Dependencies{
		ExampleRepository: exampleRepo,
	})

	handler.New(e, handler.Dependencies{
		Service: service,
	})

	// HTTP Listening
	if err := e.Start(":" + c.AppConfig.Port); err != nil && err != http.ErrServerClosed {
		log.Panicf("error - [main.New] unable to start server: %v", err)
	}
}
