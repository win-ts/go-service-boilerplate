// Package service provides the business logic service layer for the server
package service

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/repository"
)

// Port represents the service layer functions
type Port interface {
	DoExample(ctx context.Context) (string, error)
	DoWiremock(ctx context.Context) (*dto.WiremockGetTestResponse, error)
	DoDBTest() (*[]dto.TestModel, error)
}

type service struct {
	exampleRepository     repository.ExampleRepository
	wiremockAPIRepository repository.WiremockAPIRepository
	databaseRepository    repository.DatabaseRepository
}

// Dependencies represents the dependencies for the service
type Dependencies struct {
	ExampleRepository     repository.ExampleRepository
	WiremockAPIRepository repository.WiremockAPIRepository
	DatabaseRepository    repository.DatabaseRepository
}

// New creates a new service
func New(d Dependencies) Port {
	return &service{
		exampleRepository:     d.ExampleRepository,
		wiremockAPIRepository: d.WiremockAPIRepository,
		databaseRepository:    d.DatabaseRepository,
	}
}
