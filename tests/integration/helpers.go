package integration

import (
	"log/slog"
	"net/http/httptest"
	"testing"

	"session-zero-app/internal/ssr/config"
	"session-zero-app/internal/ssr/router"
	"session-zero-app/pkg/logger"
)

func SetupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	err := logger.Init(logger.Config{
		Level:       slog.LevelDebug,
		Service:     "test-ssr",
		Environment: "test",
		File:        "", // Don't write to file in tests
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	cfg := &config.Config{}

	handler := router.NewRouter(cfg)

	server := httptest.NewServer(handler)

	t.Cleanup(func() {
		server.Close()
	})

	return server
}
