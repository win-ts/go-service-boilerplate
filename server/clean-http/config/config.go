// Package config provides configuration settings for the server
package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// Config represents the configuration of the server
type Config struct {
	AppConfig AppConfig
}

// AppConfig represents the configuration of the application
type AppConfig struct {
	Name     string
	Port     string
	EnvStage string
}

// LoadConfig loads the configuration from the .env file
func LoadConfig(env string) (Config, error) {
	if env == "" || env == "local" {
		if err := godotenv.Load(".env.local", ".env.secrets"); err != nil {
			return Config{}, err
		}
	}

	return Config{
		AppConfig: AppConfig{
			Name:     requiredEnv("APP_NAME"),
			Port:     requiredEnv("APP_PORT"),
			EnvStage: requiredEnv("APP_ENV_STAGE"),
		},
	}, nil
}

func requiredEnv(env string) string {
	val := os.Getenv(env)
	if val == "" {
		log.Panic("missing required environment variable: " + env)
	}
	return val
}
