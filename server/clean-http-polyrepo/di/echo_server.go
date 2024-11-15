package di

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func setupServer(ctx context.Context, e *echo.Echo, c Config) {
	// Request Timeout
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      30 * time.Second,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	// Body Limit
	e.Use(middleware.BodyLimit("10M"))

	// Recover
	e.Use(middleware.Recover())

	// Logger
	e.Use(middleware.Logger())

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

func gracefulShutdown(ctx context.Context, e *echo.Echo, c Config, quit <-chan os.Signal) {
	log.Infof("Starting server: %s", c.AppConfig.Name)
	<-quit
	log.Info("Shutting down server ...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Panicf("[main.echoServer] shutdown: %v", err)
	}
}
