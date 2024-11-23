// Package main is the main entry point for service-name service
package main

import (
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"gitlab.com/greyxor/slogor"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/di"
	_ "github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/docs"
)

// @title API Endpoints for service-name
// @version 1.0

func init() {
	runtime.GOMAXPROCS(1)

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Panicf("failed to set timezone: %v", err)
	}
	time.Local = location
}

func main() {
	// Initialize logger
	env := os.Getenv("APP_ENV_STAGE")
	var logger *slog.Logger
	if env == "" || env == "LOCAL" {
		logger = slog.New(slogor.NewHandler(os.Stdout, slogor.SetTimeFormat(time.Stamp), slogor.ShowSource()))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}))
	}
	slog.SetDefault(logger)

	// Initiaize config
	cfg := config.New(env)

	// Initialize dependency injection
	di.New(cfg)
}
