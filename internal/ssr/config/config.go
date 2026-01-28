package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type RabbitMQConfig struct {
	URL string
	MinConnections int
	MaxConnections int
	Timeout time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SSR_PORT", "4000"),
			Host: getEnv("HOST", "localhost"),
		},
		RabbitMQ: RabbitMQConfig{
      URL:            getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
      MinConnections: getEnvInt("RABBITMQ_MIN_CONNECTIONS", 2),
      MaxConnections: getEnvInt("RABBITMQ_MAX_CONNECTIONS", 4),
      Timeout:        getEnvDuration("RABBITMQ_TIMEOUT", 5*time.Second),
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

func getEnvInt(key string, fallback int) int {
		if value := os.Getenv(key); value != "" {
				if intVal, err := strconv.Atoi(value); err == nil {
						return intVal
				}
		}
		return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
		if value := os.Getenv(key); value != "" {
				if duration, err := time.ParseDuration(value); err == nil {
						return duration
				}
		}
		return fallback
}
