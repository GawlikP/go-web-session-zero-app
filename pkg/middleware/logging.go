package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"session-zero-app/pkg/logger"
	"time"
)

var logFilter = DefaultLogFilter()

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := generateRequestID()
		logger := logger.Get().With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		ctx := context.WithValue(r.Context(), "logger", logger)
		r = r.WithContext(ctx)

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		wrapped.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(wrapped, r)
		if logFilter.ShouldSkip(r.URL.Path, wrapped.statusCode) {
			return
		}
		duration := time.Since(start)
		logger.Info("Request completed",
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

func generateRequestID() string {
	b := make([]byte, 16) // 128 bits
	rand.Read(b)
	return hex.EncodeToString(b)
}
