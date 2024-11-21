// Package repository provides the repository interfaces for the domain
package repository

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/pkg/httpclient"
)

// ExampleRepository represents the repository layer functions of example repository
type ExampleRepository interface {
	DoExample(ctx context.Context) (string, error)
}

type WiremockAPIRepository interface {
	GetTest(ctx context.Context, h dto.WiremockGetTestHeader) (*httpclient.Response[dto.WiremockGetTestResponse], error)
}
