package config

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoading(t *testing.T) {
	resetConfig()

	// config
	t.Setenv("APP_NAME", "service-name")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_ENV_STAGE", "DEV")
	t.Setenv("WIREMOCK_API_PATH", "/api")
	t.Setenv("WIREMOCK_API_MAX_CONNS", "10")
	t.Setenv("WIREMOCK_API_MAX_RETRY", "3")
	t.Setenv("WIREMOCK_API_TIMEOUT", "5s")
	t.Setenv("WIREMOCK_API_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("WIREMOCK_API_MAX_TRANSACTIONS_PER_SECOND", "10")

	// secret
	t.Setenv("SENTRY_DSN", "http://sentry.io")
	t.Setenv("WIREMOCK_API_BASE_URL", "http://localhost:1324")

	conf := New()

	assert.Equal(t, "service-name", conf.AppConfig.Name)
}

func resetConfig() {
	config = nil
	once = sync.Once{}
}
