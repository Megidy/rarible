package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL,required"`

	HttpServerPort string `env:"HTTP_SERVER_PORT,required"`

	RaribleApiKey string `env:"RARIBLE_API_KEY,required"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
