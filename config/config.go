package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Logging is the logging configuration.
	Logging struct {
		// LogLevel is the log level to use.
		Level string `envconfig:"LOG_LEVEL" default:"debug"`
		// LogFormat is the log format to use.
		Format string `envconfig:"LOG_FORMAT" default:"console"`
	}

	Connection struct {
		ServerHost string `envconfig:"SERVER_HOST" default:"localhost"`
		ServerPort int    `envconfig:"SERVER_PORT" default:"8080"`

		DB_DRIVER string `envconfig:"DB_DRIVER" default:"postgres"`
		DB_USER string `envconfig:"DB_USER" default:"root"`
		DB_PASSWORD string `envconfig:"DB_PASSWORD" default:"postgres"`
		DB_PORT string `envconfig:"DB_PORT" default:"5432"`
		DB_TABLE string `envconfig:"DB_TABLE" default:"soundproof_db"`
	}
}

// NewConfig returns a new Config instance, populated with environment variables and defaults.
func NewConfig() (*Config, error) {
	// Create a new Config instance
	cfg := &Config{}

	// Populate the Config instance with environment variables
	err := envconfig.Process("", cfg)
	// Return an error if the environment variables could not be processed
	if err != nil {
		return nil, fmt.Errorf("processing environment variables: %w", err)
	}

	// Return the Config instance and no error
	return cfg, nil
}
