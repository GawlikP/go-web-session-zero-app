package router

import (
	"net/http"
	apiHandler "session-zero-app/internal/ssr/api/handler"
	"session-zero-app/internal/ssr/config"
	"session-zero-app/internal/ssr/handler"
	"session-zero-app/pkg/middleware"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	staticHandler := handler.StaticHandler("web/static", true)
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))
	mux.Handle("GET /health", handler.HealthCheck())
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r) // Return 404
	})
	mux.HandleFunc("GET /{$}", handler.HandleLoginGet())
	mux.HandleFunc("POST /api/login", apiHandler.HandleLoginPost())
	handler := middleware.Chain(
		mux,
		middleware.LoggingMiddleware,
	)
	return handler
}
