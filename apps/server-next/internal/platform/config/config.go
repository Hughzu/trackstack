package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Env      string
	Port     string
	LogLevel string
}

func Load() (Config, error) {
	cfg := Config{
		Env:      getEnv("APP_ENV", "local"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Port) == "" {
		return fmt.Errorf("Port must not be empty.")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
