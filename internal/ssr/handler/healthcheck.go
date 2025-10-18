package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// HandleHealth returns a simple health check endpoint
func HealthCheck() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "ssr-service",
		}

		json.NewEncoder(w).Encode(response)
	})
}
