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
		name   string
		format string
		level  string
	}{
		{
			name:   "text format with debug level",
			format: "text",
			level:  "debug",
		},
		{
			name:   "json format with info level",
			format: "json",
			level:  "info",
		},
		{
			name:   "text format with warn level",
			format: "text",
			level:  "warn",
		},
		{
			name:   "json format with error level",
			format: "json",
			level:  "error",
		},
		{
			name:   "unknown format defaults to text",
			format: "unknown",
			level:  "info",
		},
		{
			name:   "empty format defaults to text",
			format: "",
			level:  "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.format, tt.level)
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
			result := ParseLevel(tt.level).SlogLevel()
			assert.Equal(t, tt.expectedLevel, result)
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name           string
		format         string
		expectedFormat Format
	}{
		{
			name:           "text format",
			format:         "text",
			expectedFormat: FormatText,
		},
		{
			name:           "json format",
			format:         "json",
			expectedFormat: FormatJSON,
		},
		{
			name:           "unknown defaults to text",
			format:         "unknown",
			expectedFormat: FormatText,
		},
		{
			name:           "empty defaults to text",
			format:         "",
			expectedFormat: FormatText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFormat(tt.format)
			assert.Equal(t, tt.expectedFormat, result)
		})
	}
}

func TestValidFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		isValid bool
	}{
		{
			name:    "text is valid",
			format:  "text",
			isValid: true,
		},
		{
			name:    "json is valid",
			format:  "json",
			isValid: true,
		},
		{
			name:    "unknown is invalid",
			format:  "unknown",
			isValid: false,
		},
		{
			name:    "empty is invalid",
			format:  "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidFormat(tt.format)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestValidLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		isValid bool
	}{
		{name: "debug is valid", level: "debug", isValid: true},
		{name: "info is valid", level: "info", isValid: true},
		{name: "warn is valid", level: "warn", isValid: true},
		{name: "error is valid", level: "error", isValid: true},
		{name: "unknown is invalid", level: "unknown", isValid: false},
		{name: "empty is invalid", level: "", isValid: false},
		{name: "uppercase is invalid", level: "DEBUG", isValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidLevel(tt.level)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestParseLevel_Type(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel Level
	}{
		{name: "debug", level: "debug", expectedLevel: LevelDebug},
		{name: "info", level: "info", expectedLevel: LevelInfo},
		{name: "warn", level: "warn", expectedLevel: LevelWarn},
		{name: "error", level: "error", expectedLevel: LevelError},
		{name: "unknown defaults to info", level: "foobar", expectedLevel: LevelInfo},
		{name: "empty defaults to info", level: "", expectedLevel: LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLevel(tt.level)
			assert.Equal(t, tt.expectedLevel, result)
		})
	}
}

func TestLevel_SlogLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         Level
		expectedLevel slog.Level
	}{
		{name: "debug", level: LevelDebug, expectedLevel: slog.LevelDebug},
		{name: "info", level: LevelInfo, expectedLevel: slog.LevelInfo},
		{name: "warn", level: LevelWarn, expectedLevel: slog.LevelWarn},
		{name: "error", level: LevelError, expectedLevel: slog.LevelError},
		{name: "unknown defaults to info", level: Level("unknown"), expectedLevel: slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.SlogLevel()
			assert.Equal(t, tt.expectedLevel, result)
		})
	}
}

func TestLevel_String(t *testing.T) {
	assert.Equal(t, "debug", LevelDebug.String())
	assert.Equal(t, "info", LevelInfo.String())
	assert.Equal(t, "warn", LevelWarn.String())
	assert.Equal(t, "error", LevelError.String())
}

func TestFormat_String(t *testing.T) {
	assert.Equal(t, "text", FormatText.String())
	assert.Equal(t, "json", FormatJSON.String())
}

func TestLoggerOutput_Text(t *testing.T) {
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

func TestLoggerOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	logger.Info("json test", "key", "value")

	output := buf.String()
	require.Contains(t, output, "json test")
	require.Contains(t, output, "key")
	require.Contains(t, output, "value")
}

func TestLoggerLevelFiltering(t *testing.T) {
	tests := []struct {
		name             string
		logLevel         slog.Level
		logFunc          func(*slog.Logger)
		shouldContain    string
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
