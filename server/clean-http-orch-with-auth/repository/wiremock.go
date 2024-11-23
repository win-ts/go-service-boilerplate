package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/dto"
	"github.com/win-ts/go-service-boilerplate/server/clean-http-orch-with-auth/pkg/httpclient"
)

type wiremockAPIRepository struct {
	baseURL string
	path    string
	client  *http.Client
}

// WiremockAPIRepositoryConfig represents the configuration for wiremock API repository
type WiremockAPIRepositoryConfig struct {
	BaseURL string
	Path    string
}

// WiremockAPIRepositoryDependencies represents the dependencies for wiremock API repository
type WiremockAPIRepositoryDependencies struct {
	Client *http.Client
}

// NewWiremockAPIRepository creates a new wiremock API repository
func NewWiremockAPIRepository(c WiremockAPIRepositoryConfig, d WiremockAPIRepositoryDependencies) WiremockAPIRepository {
	return &wiremockAPIRepository{
		baseURL: c.BaseURL,
		path:    c.Path,
		client:  d.Client,
	}
}

// GetTest returns the test response from wiremock API
func (r *wiremockAPIRepository) GetTest(ctx context.Context, h dto.WiremockGetTestHeader) (*httpclient.Response[dto.WiremockGetTestResponse], error) {
	url := fmt.Sprintf("%s%s", r.baseURL, r.path)
	return httpclient.Get[dto.WiremockGetTestResponse](ctx, r.client, url, h.ToMap())
}
