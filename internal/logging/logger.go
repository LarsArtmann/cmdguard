// Package logging provides structured logging using slog with text or JSON handler.
package logging

import (
	"io"
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

// Level represents the log severity level.
type Level string

const (
	// LevelDebug enables all logs.
	LevelDebug Level = "debug"
	// LevelInfo enables info and above.
	LevelInfo Level = "info"
	// LevelWarn enables warn and above.
	LevelWarn Level = "warn"
	// LevelError enables only error logs.
	LevelError Level = "error"
)

// NewLogger creates a new slog.Logger with the specified format and level.
// Format can be "text" or "json". Defaults to text.
// Level can be "debug", "info", "warn", "error". Defaults to info.
// Writes to os.Stderr by default.
func NewLogger(format, level string) *slog.Logger {
	return NewLoggerWriter(format, level, os.Stderr)
}

// NewLoggerWriter creates a new slog.Logger writing to the given writer.
// Falls back to text format for unrecognized formats.
func NewLoggerWriter(format, level string, w io.Writer) *slog.Logger {
	logLevel := ParseLevel(level).SlogLevel()
	logFormat := ParseFormat(format)

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler

	switch logFormat {
	case FormatJSON:
		handler = slog.NewJSONHandler(w, opts)
	case FormatText:
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	return slog.New(handler)
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

// ValidLevel returns true if the level is valid.
func ValidLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}

// ParseLevel converts string level to Level type.
func ParseLevel(level string) Level {
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// SlogLevel converts Level to slog.Level.
func (l Level) SlogLevel() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// ParseFormat converts string format to Format type.
func ParseFormat(format string) Format {
	switch format {
	case "json":
		return FormatJSON
	case "text":
		return FormatText
	default:
		return FormatText
	}
}

// String returns the string representation of Format.
func (f Format) String() string {
	return string(f)
}

// String returns the string representation of Level.
func (l Level) String() string {
	return string(l)
}
