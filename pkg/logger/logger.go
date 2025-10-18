package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var globalLogger *slog.Logger

type Config struct {
	Service     string
	Environment string
	Level       slog.Level
	Format      string
	File        string
}

func Init(cfg Config) error {
	writer, err := createWriters(cfg)
	if err != nil {
		return err
	}
	handler := createHandler(writer, cfg)
	logger := createLogger(handler, cfg)
	globalLogger = logger
	slog.SetDefault(globalLogger)
	return nil
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

func Get() *slog.Logger {
	return globalLogger
}

func createWriters(cfg Config) (io.Writer, error) {
	writers := []io.Writer{}

	writers = append(writers, os.Stdout)
	if cfg.File != "" {
		file, err := openLogFile(cfg.File)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}
	return io.MultiWriter(writers...), nil
}

func createHandler(writer io.Writer, cfg Config) slog.Handler {
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: true,
	}
	if cfg.Format == "json" {
		return slog.NewJSONHandler(writer, opts)
	}
	return slog.NewTextHandler(writer, opts)
}

func createLogger(handler slog.Handler, cfg Config) *slog.Logger {
	return slog.New(handler).With(
		"service", cfg.Service,
		"environment", cfg.Environment,
	)
}

func openLogFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return file, nil
}
