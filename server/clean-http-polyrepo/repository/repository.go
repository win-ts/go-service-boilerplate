// Package repository provides the repository interfaces for the domain
package repository

import (
	"context"
)

// ExampleRepository represents the repository layer functions of example repository
type ExampleRepository interface {
	DoExample(ctx context.Context) (string, error)
}
