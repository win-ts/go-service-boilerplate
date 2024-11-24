// Package main is the main entry point for service-name service
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"gitlab.com/greyxor/slogor"

	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/config"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/di"
)

func init() {
	runtime.GOMAXPROCS(1)

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Panicf("failed to set timezone: %v", err)
	}
	time.Local = location
}

func main() {
	ctx := context.Background()
	now := time.Now()

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
	s, mysqlClient, redisClient := di.New(cfg)
	defer func() {
		if err := mysqlClient.Close(); err != nil {
			slog.Error("error - [main.New] unable to close MySQL connection", slog.Any("error", err))
		}
		if err := redisClient.Close(); err != nil {
			slog.Error("error - [main.New] unable to close Redis connection", slog.Any("error", err))
		}
	}()

	// Start job
	{
		slog.Info("starting job")

		cacheRes, err := s.DoSetGetCache(ctx)
		if err != nil {
			slog.Error("error - [main] unable to set/get cache", slog.Any("error", err))
		}

		slog.Info("DoSetGetCache",
			slog.Any("result", cacheRes),
		)

		wiremockRes, err := s.DoWiremock(ctx)
		if err != nil {
			slog.Error("error - [main] unable to call Wiremock API", slog.Any("error", err))
		}

		slog.Info("DoWiremock",
			slog.Any("result", wiremockRes),
		)

		slog.Info("job completed", slog.Any("duration", time.Since(now)))
	}
}
