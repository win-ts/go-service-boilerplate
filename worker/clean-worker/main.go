// Package main is the main entry point for service-name service
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"gitlab.com/greyxor/slogor"

	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/config"
	"github.com/win-ts/go-service-boilerplate/worker/clean-worker/di"
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

	// Start worker
	{
		ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer stop()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					start := time.Now()

					cacheRes, err := s.DoSetGetCache(ctx)
					if err != nil {
						slog.Error("[main]: error - DoSetGetCache", slog.Any("error", err))
					}
					slog.Info("[main]: cache result", slog.Any("cacheRes", cacheRes))

					wiremockRes, err := s.DoWiremock(ctx)
					if err != nil {
						slog.Error("[main]: error - DoWiremock", slog.Any("error", err))
					}
					slog.Info("[main]: wiremock result", slog.Any("wiremockRes", wiremockRes))

					slog.Info(fmt.Sprintf("[main]: worker success ✅, waiting for %v to execute in next round", cfg.AppConfig.WorkerInterval),
						slog.Any("duration", time.Since(start)),
					)

					time.Sleep(cfg.AppConfig.WorkerInterval)
				}
			}
		}()

		wg.Wait()
		slog.Info("[main]: gracefully shutting down worker ...")
	}
}
