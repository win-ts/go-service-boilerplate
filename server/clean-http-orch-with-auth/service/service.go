// Package service provides the business logic service layer for the server
package service

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/repository"
)

// Port represents the service layer functions
type Port interface {
	Auth(c echo.Context, token string) (echo.Context, error)

	DoExample(ctx context.Context) (string, error)
	DoWiremock(ctx context.Context) (*dto.WiremockGetTestResponse, error)
}

type service struct {
	authMiddlewareRepository repository.AuthMiddlewareRepository
	exampleRepository        repository.ExampleRepository
	wiremockAPIRepository    repository.WiremockAPIRepository
}

// Dependencies represents the dependencies for the service
type Dependencies struct {
	AuthMiddlewareRepository repository.AuthMiddlewareRepository
	ExampleRepository        repository.ExampleRepository
	WiremockAPIRepository    repository.WiremockAPIRepository
}

// New creates a new service
func New(d Dependencies) Port {
	return &service{
		authMiddlewareRepository: d.AuthMiddlewareRepository,
		exampleRepository:        d.ExampleRepository,
		wiremockAPIRepository:    d.WiremockAPIRepository,
	}
}
