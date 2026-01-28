package rabbitmq

import "testing"

func TestPoolConfig_Validate(t *testing.T) {
		t.Run("valid config", func(t *testing.T) {
			config := PoolConfig{
				URL: "amqp://admin:admin@localhost:5672/",
				ConnectionTimeout: 1,
				MinConnections: 1,
				MaxConnections: 2,
			}
			err := config.Validate()
			if err != nil {
				t.Errorf("Expected no error for valid config, got: %v", err)
			}
		})
		t.Run("URL should not be empty", func(t *testing.T) {
			config := PoolConfig{
				URL: "",
				ConnectionTimeout: 1,
				MinConnections: 1,
				MaxConnections: 2,
			}
			err := config.Validate()
			if err == nil {
				t.Error("Expected empty URL to be invalid")
			}
		})
		t.Run("MinConnections should be higher than 1", func(t *testing.T) {
			config := PoolConfig{
				URL: "amqp://admin:admin@localhost:5672/",
				ConnectionTimeout: 1,
				MinConnections: 0,
				MaxConnections: 2,
			}
			err := config.Validate()
			if err == nil {
				t.Error("Expected MinConnections to be invalid when less than 0")
			}
		})
		t.Run("MaxConnections should be lower than 100", func(t *testing.T) {
			config := PoolConfig{
				URL: "amqp://admin:admin@localhost:5672/",
				ConnectionTimeout: 1,
				MinConnections: 1,
				MaxConnections: 100,
			}
			err := config.Validate()
			if err == nil {
				t.Error("Expected MaxConnections to be invalid when less than 0")
			}
		})
		t.Run("MaxConnections should be higher than Min", func(t *testing.T) {
			config := PoolConfig{
				URL: "amqp://admin:admin@localhost:5672/",
				ConnectionTimeout: 1,
				MinConnections: 3,
				MaxConnections: 2,
			}
			err := config.Validate()
			if err == nil {
				t.Error("Expected MaxConnections to be invalid when less than MinConnections")
			}
		})
}

func TestPoolConfig_Default(t *testing.T) {
	t.Run("Default config should be valid", func(t *testing.T) {
		config := DefaultConfig("amqp://admin:admin@localhost:5672/")
		err := config.Validate()
		if err != nil {
			t.Errorf("Expected no error for default config, got: %v", err)
		}
	})
}
