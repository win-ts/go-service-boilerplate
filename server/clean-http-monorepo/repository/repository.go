// Package repository provides the repository interfaces for the domain
package repository

import (
	"context"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock -source=./repository.go ExampleRepository

// ExampleRepository represents the repository layer functions of example repository
type ExampleRepository interface {
	DoExample(ctx context.Context) (string, error)
}
