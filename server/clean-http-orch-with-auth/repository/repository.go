// Package repository provides the repository interfaces for the domain
package repository

import (
	"context"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/pkg/httpclient"
)

// AuthMiddlewareRepository represents the repository layer functions of auth middleware repository
type AuthMiddlewareRepository interface {
	VerifyToken(ctx context.Context, token string) error
}

// ExampleRepository represents the repository layer functions of example repository
type ExampleRepository interface {
	DoExample(ctx context.Context) (string, error)
}

// WiremockAPIRepository represents the repository layer functions of wiremock API repository
type WiremockAPIRepository interface {
	GetTest(ctx context.Context, h dto.WiremockGetTestHeader) (*httpclient.Response[dto.WiremockGetTestResponse], error)
}
