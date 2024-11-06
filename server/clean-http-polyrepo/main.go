// Package main is the main entry point for servicename service
package main

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/di"
	_ "github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/docs"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/repository"
)

// @title API Endpoints for servicename
// @version 1.0
// @description This is the API documentation for servicename

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
