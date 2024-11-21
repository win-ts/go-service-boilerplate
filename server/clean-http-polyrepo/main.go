// Package main is the main entry point for service-name service
package main

import (
	"runtime"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/config"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/di"
	_ "github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/docs"
)

// @title API Endpoints for service-name
// @version 1.0
// @description This is the API documentation for service-name

func init() {
	runtime.GOMAXPROCS(1)
}

// @BasePath /server
func main() {
	// Initiaize config
	cfg := config.New()

	// Initialize dependency injection
	di.New(cfg)
}
