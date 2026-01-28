package integration

import (
	"log/slog"
	"net/http/httptest"
	"testing"
	"context"
	"time"
	"log"
	"os"
	"path/filepath"

	"session-zero-app/internal/ssr/config"
	"session-zero-app/internal/ssr/router"
	"session-zero-app/pkg/logger"
	"github.com/joho/godotenv"
	rbit "session-zero-app/pkg/messaging/rabbitmq"
)

func TestSetup(t *testing.T) error {
	t.Helper()
	root, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	envPath := filepath.Join(root, ".env.test")
	if err := godotenv.Load(envPath); err != nil {
					log.Printf("No .env.test found at %s, using environment variables", envPath)
	}
	err = logger.Init(logger.Config{
		Level:       slog.LevelDebug,
		Service:     "test-ssr",
		Environment: os.Getenv("ENV"),
		File:        "", // Don't write to file in tests
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	return nil
}

func SetupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolCfg := rbit.PoolConfig{
		URL: cfg.RabbitMQ.URL,
		MinConnections: cfg.RabbitMQ.MinConnections,
		MaxConnections: cfg.RabbitMQ.MaxConnections,
		ConnectionTimeout: cfg.RabbitMQ.Timeout,
	}
	pool, err := rbit.NewConnectionPool(&poolCfg)
	if err != nil {
		logger.Info("RabbitMQ Pool Error", "Pool Error:", err.Error())
	}
	defer pool.CloseAll()

	err = rbit.InitializeTopology(ctx,pool)
	if err != nil {
		logger.Info("RabbitMQ Pool Error", "Topology error:", err.Error())
	}

	handler := router.NewRouter(cfg, pool)

	server := httptest.NewServer(handler)

	t.Cleanup(func() {
		server.Close()
	})

	return server
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
