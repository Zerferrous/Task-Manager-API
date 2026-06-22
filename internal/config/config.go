package config

import (
	"fmt"

	"github.com/Zerferrous/Task-Manager-API/internal/logger"
	"github.com/Zerferrous/Task-Manager-API/internal/transport/http/server"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPServer server.Config
	Logger     logger.Config
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		err := fmt.Errorf("failed to load config: %w", err)
		panic(err)
	}

	return config
}
