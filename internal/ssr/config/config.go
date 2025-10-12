package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port string
	Host string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SSR_PORT", "4000"),
			Host: getEnv("HOST", "localhost"),
		},
	}

	return cfg, nil
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
