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

	// secret

	conf := New()

	assert.Equal(t, "service-name", conf.AppConfig.Name)
}

func resetConfig() {
	config = nil
	once = sync.Once{}
}
