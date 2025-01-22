// Package config provides configuration settings for the server
package config

import (
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var once sync.Once
var config *Config

// New loads the configuration from the .env file
func New() *Config {
	once.Do(func() {
		e := os.Getenv("APP_ENV_STAGE")
		if e == "" || e == "LOCAL" {
			if err := godotenv.Load(".env.generated"); err != nil {
				slog.Warn("[config.New] unable to load .env.generated file", slog.Any("error", err))
			}
		}

		cfg := &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Panicf("error - [config.New] unable to parse config: %v", err)
		}
		config = cfg
	})

	return config
}

// Config represents the configuration of the server
type Config struct {
	AppConfig         AppConfig
	LogConfig         LogConfig
	GRPCAuthConfig    GRPCAuthConfig
	SentryConfig      SentryConfig
	WiremockAPIConfig WiremockAPIConfig
}

// AppConfig represents the configuration of the application
type AppConfig struct {
	Name     string `env:"APP_NAME,notEmpty"`
	Port     string `env:"APP_PORT,notEmpty"`
	EnvStage string `env:"APP_ENV_STAGE,notEmpty"`
}

// LogConfig represents the configuration of the logger
type LogConfig struct {
	Level             string `env:"LOG_LEVEL,notEmpty"`
	MaskSensitiveData bool   `env:"LOG_MASK_SENSITIVE_DATA,notEmpty"`
}

// GRPCAuthConfig represents the configuration of the Auth GRPC client
type GRPCAuthConfig struct {
	URL string `env:"GRPC_AUTH_URL,notEmpty"`
}

// SentryConfig represents the configuration of Sentry.io
type SentryConfig struct {
	SentryDSN string `env:"SENTRY_DSN"`
}

// WiremockAPIConfig represents the configuration of the Wiremock API
type WiremockAPIConfig struct {
	BaseURL                  string        `env:"WIREMOCK_API_BASE_URL,notEmpty"`
	Path                     string        `env:"WIREMOCK_API_PATH,notEmpty"`
	MaxConns                 int           `env:"WIREMOCK_API_MAX_CONNS,notEmpty"`
	MaxRetry                 int           `env:"WIREMOCK_API_MAX_RETRY"`
	Timeout                  time.Duration `env:"WIREMOCK_API_TIMEOUT,notEmpty"`
	InsecureSkipVerify       bool          `env:"WIREMOCK_API_INSECURE_SKIP_VERIFY,notEmpty"`
	MaxTransactionsPerSecond int           `env:"WIREMOCK_API_MAX_TRANSACTIONS_PER_SECOND"`
}
