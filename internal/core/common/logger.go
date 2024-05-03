package common

import (
	"log/slog"
	"os"
)

func NewLogger(logLevel slog.Level) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))
	return logger
}
