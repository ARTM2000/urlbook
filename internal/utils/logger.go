package utils

import (
	"log/slog"
	"os"
)

func NewLogger(logLevel slog.Level) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	return logger
}
