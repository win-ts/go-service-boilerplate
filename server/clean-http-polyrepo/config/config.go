// Package config provides configuration settings for the server
package config

import (
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

var once sync.Once
var config *Config

// New loads the configuration from the .env file
func New() *Config {
	once.Do(func() {
		e := os.Getenv("APP_ENV_STAGE")
		if e == "" {
			if err := godotenv.Load(".env.local"); err != nil {
				log.Warnf("Warning - [config.New] unable to load .env.local file: %v", err)
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
	AppConfig AppConfig
}

// AppConfig represents the configuration of the application
type AppConfig struct {
	Name     string `env:"APP_NAME,notEmpty"`
	Port     string `env:"APP_PORT,notEmpty"`
	EnvStage string `env:"APP_ENV_STAGE,notEmpty"`
}
