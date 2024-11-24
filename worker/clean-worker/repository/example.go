package repository

import (
	"context"
)

type exampleRepository struct {
	config ExampleRepositoryConfig
}

// ExampleRepositoryConfig represents the configuration for example repository
type ExampleRepositoryConfig struct {
}

// NewExampleRepository creates a new example repository
func NewExampleRepository(c ExampleRepositoryConfig) ExampleRepository {
	return &exampleRepository{
		config: c,
	}
}

// DoExample returns example string
func (r *exampleRepository) DoExample(ctx context.Context) (string, error) {
	_ = ctx
	return "example", nil
}
