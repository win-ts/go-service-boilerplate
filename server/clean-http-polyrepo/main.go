// Package main is the main entry point for servicename service
package main

import (
	"runtime"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/di"
	_ "github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/docs"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/repository"
)

// @title API Endpoints for servicename
// @version 1.0
// @description This is the API documentation for servicename

func init() {
	runtime.GOMAXPROCS(1)
}

// @BasePath /server
func main() {
	// Initiaize config
	cfg := config.New()

	// Initialize dependency injection
	di.New(di.Config{
		AppConfig:               cfg.AppConfig,
		ExampleRepositoryConfig: repository.ExampleRepositoryConfig{},
	})
}
