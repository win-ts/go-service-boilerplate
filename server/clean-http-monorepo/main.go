// Package main is the main entry point for domain-service service
package main

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-monorepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-monorepo/di"
	_ "github.com/win-ts/go-service-boilerplate/server/clean-http-monorepo/docs"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-monorepo/repository"
)

// @title API Endpoints for domain-service
// @version 1.0
// @description This is the API documentation for domain-service

// @BasePath /server
func main() {
	env := os.Getenv("APP_ENV_STAGE")
	if env == "" {
		env = "local"
	}

	// Initiaize Config
	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Panicf("error - [main.loadConfig] error loading config: %v", err)
	}

	di.New(di.Config{
		AppConfig:               cfg.AppConfig,
		ExampleRepositoryConfig: repository.ExampleRepositoryConfig{},
	})
}
