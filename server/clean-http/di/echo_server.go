package di

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func setupServer(pctx context.Context, e *echo.Echo, c Config) {
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

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	e.Use(middleware.Logger())
	go gracefulShutdown(pctx, e, c, quit)
}

func gracefulShutdown(pctx context.Context, e *echo.Echo, c Config, quit <-chan os.Signal) {
	log.Infof("Starting server: %s", c.AppConfig.Name)
	<-quit
	log.Info("Shutting down server ...")

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Panicf("[main.echoServer] shutdown: %v", err)
	}
}
