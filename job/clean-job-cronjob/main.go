// Package main is the main entry point for service-name service
package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
	_ "time/tzdata"

	"gitlab.com/greyxor/slogor"

	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/config"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/di"
	"github.com/win-ts/go-service-boilerplate/job/clean-job-cronjob/middleware"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {
	ctx := context.Background()
	start := time.Now()

	// Initiaize config
	cfg := config.New()

	// Initialize logger
	var logger *slog.Logger
	if cfg.LogConfig.Level == "DEBUG" {
		logger = slog.New(slogor.NewHandler(os.Stdout, slogor.SetTimeFormat(time.Stamp), slogor.ShowSource()))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if cfg.LogConfig.MaskSensitiveData {
					value := middleware.MaskSensitiveData(a.Key, a.Value.Any())
					return slog.Any(a.Key, value)
				}

				return a
			},
		}))
	}
	slog.SetDefault(logger)

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
		slog.Info("[main]: starting job")

		cacheRes, err := s.DoSetGetCache(ctx)
		if err != nil {
			slog.Error("error - [main] unable to set/get cache", slog.Any("error", err))
		}

		slog.Info("[main]: DoSetGetCache",
			slog.Any("result", cacheRes),
		)

		wiremockRes, err := s.DoWiremock(ctx)
		if err != nil {
			slog.Error("error - [main] unable to call Wiremock API", slog.Any("error", err))
		}

		slog.Info("[main]: DoWiremock",
			slog.Any("result", wiremockRes),
		)

		slog.Info("[main]: job completed ✅", slog.Any("duration", time.Since(start)))
	}
}
