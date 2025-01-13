package di

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/middleware"
)

func setupServer(ctx context.Context, e *echo.Echo, c *config.Config) {
	// Request Timeout
	e.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Body Limit
	e.Use(echoMiddleware.BodyLimit("10M"))

	// Recover
	e.Use(echoMiddleware.Recover())

	// Logger
	e.Use(middleware.IncomingLogTrace())

	// Sentry
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
				hub.Scope().SetTag("service-name", c.AppConfig.Name)
			}
			return next(ctx)
		}
	})

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go gracefulShutdown(ctx, e, c, quit)
}

func gracefulShutdown(ctx context.Context, e *echo.Echo, c *config.Config, quit <-chan os.Signal) {
	slog.Info("Starting server...",
		slog.String("name", c.AppConfig.Name),
	)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Panicf("[main.echoServer] shutdown: %v", err)
	}
}
