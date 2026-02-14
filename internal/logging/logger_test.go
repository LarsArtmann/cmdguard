package logging

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel slog.Level
	}{
		{
			name:          "debug level",
			level:         "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info level",
			level:         "info",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn level",
			level:         "warn",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error level",
			level:         "error",
			expectedLevel: slog.LevelError,
		},
		{
			name:          "unknown level defaults to info",
			level:         "unknown",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "empty level defaults to info",
			level:         "",
			expectedLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.level)
			assert.NotNil(t, logger)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel slog.Level
	}{
		{
			name:          "debug",
			level:         "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info",
			level:         "info",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn",
			level:         "warn",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error",
			level:         "error",
			expectedLevel: slog.LevelError,
		},
		{
			name:          "unknown defaults to info",
			level:         "foobar",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "empty defaults to info",
			level:         "",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "uppercase defaults to info",
			level:         "DEBUG",
			expectedLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLevel(tt.level)
			assert.Equal(t, tt.expectedLevel, result)
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	require.Contains(t, output, "debug message")
	require.Contains(t, output, "info message")
	require.Contains(t, output, "warn message")
	require.Contains(t, output, "error message")
}

func TestLoggerLevelFiltering(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       slog.Level
		logFunc        func(*slog.Logger)
		shouldContain  string
		shouldNotContain string
	}{
		{
			name:     "warn level filters debug and info",
			logLevel: slog.LevelWarn,
			logFunc: func(l *slog.Logger) {
				l.Debug("debug msg")
				l.Info("info msg")
				l.Warn("warn msg")
			},
			shouldContain:    "warn msg",
			shouldNotContain: "debug msg",
		},
		{
			name:     "error level filters warn",
			logLevel: slog.LevelError,
			logFunc: func(l *slog.Logger) {
				l.Warn("warn msg")
				l.Error("error msg")
			},
			shouldContain:    "error msg",
			shouldNotContain: "warn msg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
				Level: tt.logLevel,
			})
			logger := slog.New(handler)

			tt.logFunc(logger)

			output := buf.String()
			assert.Contains(t, output, tt.shouldContain)
			assert.NotContains(t, output, tt.shouldNotContain)
		})
	}
}
