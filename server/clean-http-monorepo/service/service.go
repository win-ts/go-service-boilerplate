// Package service provides the business logic service layer for the server
package service

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-monorepo/repository"
)

// Port represents the service layer functions
type Port interface {
	DoExample(pctx context.Context) (string, error)
}

type service struct {
	exampleRepository repository.ExampleRepository
}

// Dependencies represents the dependencies for the service
type Dependencies struct {
	ExampleRepository repository.ExampleRepository
}

// New creates a new service
func New(d Dependencies) Port {
	return &service{
		exampleRepository: d.ExampleRepository,
	}
}
