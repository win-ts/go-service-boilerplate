// Package config provides configuration settings for the consumer
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

// Config represents the configuration of the consumer
type Config struct {
	AppConfig           AppConfig
	LogConfig           LogConfig
	SentryConfig        SentryConfig
	WiremockAPIConfig   WiremockAPIConfig
	PostgreSQLConfig    PostgreSQLConfig
	RedisConfig         RedisConfig
	KafkaConsumerConfig KafkaConsumerConfig
}

// AppConfig represents the configuration of the application
type AppConfig struct {
	Name     string `env:"APP_NAME,notEmpty"`
	EnvStage string `env:"APP_ENV_STAGE,notEmpty"`
}

// LogConfig represents the configuration of the logger
type LogConfig struct {
	Level             string `env:"LOG_LEVEL,notEmpty"`
	MaskSensitiveData bool   `env:"LOG_MASK_SENSITIVE_DATA,notEmpty"`
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
	MaxRetry                 int           `env:"WIREMOCK_API_MAX_RETRY,notEmpty"`
	Timeout                  time.Duration `env:"WIREMOCK_API_TIMEOUT,notEmpty"`
	InsecureSkipVerify       bool          `env:"WIREMOCK_API_INSECURE_SKIP_VERIFY,notEmpty"`
	MaxTransactionsPerSecond int           `env:"WIREMOCK_API_MAX_TRANSACTIONS_PER_SECOND"`
}

// PostgreSQLConfig represents the configuration of the PostgreSQL database
type PostgreSQLConfig struct {
	Host         string        `env:"POSTGRES_HOST,notEmpty"`
	Username     string        `env:"POSTGRES_USERNAME,notEmpty"`
	Password     string        `env:"POSTGRES_PASSWORD,notEmpty"`
	Database     string        `env:"POSTGRES_DATABASE,notEmpty"`
	Timeout      string        `env:"POSTGRES_TIMEOUT,notEmpty"`
	MaxIdleConns int           `env:"POSTGRES_MAX_IDLE_CONNS,notEmpty"`
	MaxOpenConns int           `env:"POSTGRES_MAX_OPEN_CONNS,notEmpty"`
	MaxLifetime  time.Duration `env:"POSTGRES_MAX_LIFETIME,notEmpty"`
}

// RedisConfig represents the configuration of the Redis cache
type RedisConfig struct {
	Host     string        `env:"REDIS_HOST,notEmpty"`
	Password string        `env:"REDIS_PASSWORD,notEmpty"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT,notEmpty"`
	MaxRetry int           `env:"REDIS_MAX_RETRY,notEmpty"`
	PoolSize int           `env:"REDIS_POOL_SIZE,notEmpty"`
}

// KafkaConsumerConfig represents the configuration of the Kafka consumer
type KafkaConsumerConfig struct {
	Username          string        `env:"KAFKA_USERNAME,notEmpty"`
	Password          string        `env:"KAFKA_PASSWORD,notEmpty"`
	SessionTimeout    time.Duration `env:"KAFKA_SESSION_TIMEOUT"`
	HeartbeatInterval time.Duration `env:"KAFKA_HEARTBEAT_INTERVAL"`
	BufferSize        int           `env:"KAFKA_CONSUMER_BUFFER_SIZE"`
	MaxRetry          int           `env:"KAFKA_CONSUMER_MAX_RETRY"`
	Brokers           []string      `env:"KAFKA_BROKERS,notEmpty"`
	Group             string        `env:"KAFKA_CONSUMER_GROUP,notEmpty"`
	Topic             string        `env:"KAFKA_CONSUMER_TOPIC,notEmpty"`
}
