// Package logging provides structured logging using slog with text or JSON handler.
package logging

import (
	"log/slog"
	"os"
)

// Format represents the log output format.
type Format string

const (
	// FormatText uses human-readable text output.
	FormatText Format = "text"
	// FormatJSON uses JSON output for machine parsing.
	FormatJSON Format = "json"
)

// NewLogger creates a new slog.Logger with the specified format and level.
// Format can be "text" or "json". Defaults to text.
// Level can be "debug", "info", "warn", "error". Defaults to info.
func NewLogger(format, level string) *slog.Logger {
	logLevel := parseLevel(level)
	logFormat := parseFormat(format)

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	switch logFormat {
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

// parseFormat converts string format to Format type.
func parseFormat(format string) Format {
	switch format {
	case "json":
		return FormatJSON
	case "text":
		return FormatText
	default:
		return FormatText
	}
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

// ValidFormat returns true if the format is valid.
func ValidFormat(format string) bool {
	switch format {
	case "text", "json":
		return true
	default:
		return false
	}
}
