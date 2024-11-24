// Package service provides the business logic service layer for the server
package service

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-amqp-producer/repository"
)

// Port represents the service layer functions
type Port interface {
	DoExample(ctx context.Context) (string, error)
	DoWiremock(ctx context.Context) (*dto.WiremockGetTestResponse, error)
	DoDBTest() (*[]dto.TestModel, error)
	DoSetGetCache(ctx context.Context) (*dto.TestModel, error)
	DoProduce(ctx context.Context, event dto.Event) error
}

type service struct {
	amqpProducerRepository repository.AMQPProducerRepository
	exampleRepository      repository.ExampleRepository
	wiremockAPIRepository  repository.WiremockAPIRepository
	databaseRepository     repository.DatabaseRepository
	cacheRepository        repository.CacheRepository
}

// Dependencies represents the dependencies for the service
type Dependencies struct {
	AMQPProducerRepository repository.AMQPProducerRepository
	ExampleRepository      repository.ExampleRepository
	WiremockAPIRepository  repository.WiremockAPIRepository
	DatabaseRepository     repository.DatabaseRepository
	CacheRepository        repository.CacheRepository
}

// New creates a new service
func New(d Dependencies) Port {
	return &service{
		amqpProducerRepository: d.AMQPProducerRepository,
		exampleRepository:      d.ExampleRepository,
		wiremockAPIRepository:  d.WiremockAPIRepository,
		databaseRepository:     d.DatabaseRepository,
		cacheRepository:        d.CacheRepository,
	}
}
