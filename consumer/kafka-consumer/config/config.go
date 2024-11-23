// Package config provides configuration settings for the consumer
package config

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var once sync.Once
var config *Config

// New loads the configuration from the .env file
func New(e string) *Config {
	once.Do(func() {
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
	SentryConfig        SentryConfig
	WiremockAPIConfig   WiremockAPIConfig
	MySQLConfig         MySQLConfig
	RedisConfig         RedisConfig
	KafkaConsumerConfig KafkaConsumerConfig
}

// AppConfig represents the configuration of the application
type AppConfig struct {
	Name     string `env:"APP_NAME,notEmpty"`
	EnvStage string `env:"APP_ENV_STAGE,notEmpty"`
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

// MySQLConfig represents the configuration of the MySQL database
type MySQLConfig struct {
	Host         string        `env:"MYSQL_HOST,notEmpty"`
	Username     string        `env:"MYSQL_USERNAME,notEmpty"`
	Password     string        `env:"MYSQL_PASSWORD,notEmpty"`
	Database     string        `env:"MYSQL_DATABASE,notEmpty"`
	Timeout      time.Duration `env:"MYSQL_TIMEOUT,notEmpty"`
	MaxIdleConns int           `env:"MYSQL_MAX_IDLE_CONNS,notEmpty"`
	MaxOpenConns int           `env:"MYSQL_MAX_OPEN_CONNS,notEmpty"`
	MaxLifetime  time.Duration `env:"MYSQL_MAX_LIFETIME,notEmpty"`
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
