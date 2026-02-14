// Package logging provides structured logging using slog with text handler.
package logging

import (
	"log/slog"
	"os"
)

// NewLogger creates a new slog.Logger with text handler.
func NewLogger(level string) *slog.Logger {
	logLevel := parseLevel(level)

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}

// parseLevel converts string level to slog.Level.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
