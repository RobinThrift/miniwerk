package app

import (
	"strings"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Addr    string `env:"MINIWERK_ADDRESS" envDefault:":8080"`
	BaseURL string `env:"MINIWERK_BASE_URL"`

	LogLevel  string `env:"MINIWERK_LOG_LEVEL" envDevault:"info"`
	LogFormat string `env:"MINIWERK_LOG_FORMAT" envDevault:"json"`
}

func NewConfigFromEnv() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	if config.BaseURL == "" {
		config.BaseURL = "localhost" + config.Addr
	}

	if !strings.HasPrefix(config.BaseURL, "http") {
		config.BaseURL = "http://" + config.BaseURL
	}

	return &config, nil
}
