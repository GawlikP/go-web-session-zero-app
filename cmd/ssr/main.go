package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"session-zero-app/internal/ssr/config"
	"session-zero-app/internal/ssr/router"
	"session-zero-app/pkg/logger"
	rbit "session-zero-app/pkg/messaging/rabbitmq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	err = logger.Init(logger.Config{
		Service:     "SSR",
		Environment: os.Getenv("ENV"),
		Level:       slog.LevelInfo,
		File:        "./logs/ssr.log",
	})
	if err != nil {
		log.Fatalf("Logger error: %v", err)
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
		log.Fatalf("Server error: %v", err)
	}
	defer pool.CloseAll()

	err = rbit.InitializeTopology(ctx,pool)
	if err != nil {
		logger.Info("RabbitMQ Pool Error", "Topology error:", err.Error())
		log.Fatalf("Server error: %v", err)
	}
	rt := router.NewRouter(cfg, pool)

	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      rt,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("SSR service starting on %s", "Adress", cfg.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
