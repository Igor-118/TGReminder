package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	LogLevel slog.Level
}

func NewLogger(cfg *Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
