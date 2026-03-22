package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
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
			if logger == nil {
				t.Error("NewLogger() returned nil, expected non-nil logger")
			}
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
			if result != tt.expectedLevel {
				t.Errorf(
					"ParseLevel(%q).SlogLevel() = %v, want %v",
					tt.level,
					result,
					tt.expectedLevel,
				)
			}
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
			if result != tt.expectedFormat {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.format, result, tt.expectedFormat)
			}
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
			if result != tt.isValid {
				t.Errorf("ValidFormat(%q) = %v, want %v", tt.format, result, tt.isValid)
			}
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
			if result != tt.isValid {
				t.Errorf("ValidLevel(%q) = %v, want %v", tt.level, result, tt.isValid)
			}
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
			if result != tt.expectedLevel {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.level, result, tt.expectedLevel)
			}
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
			if result != tt.expectedLevel {
				t.Errorf("Level(%q).SlogLevel() = %v, want %v", tt.level, result, tt.expectedLevel)
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	if got := LevelDebug.String(); got != "debug" {
		t.Errorf("LevelDebug.String() = %q, want %q", got, "debug")
	}

	if got := LevelInfo.String(); got != "info" {
		t.Errorf("LevelInfo.String() = %q, want %q", got, "info")
	}

	if got := LevelWarn.String(); got != "warn" {
		t.Errorf("LevelWarn.String() = %q, want %q", got, "warn")
	}

	if got := LevelError.String(); got != "error" {
		t.Errorf("LevelError.String() = %q, want %q", got, "error")
	}
}

func TestFormat_String(t *testing.T) {
	if got := FormatText.String(); got != "text" {
		t.Errorf("FormatText.String() = %q, want %q", got, "text")
	}

	if got := FormatJSON.String(); got != "json" {
		t.Errorf("FormatJSON.String() = %q, want %q", got, "json")
	}
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
	if !strings.Contains(output, "debug message") {
		t.Error("output missing 'debug message'")
	}

	if !strings.Contains(output, "info message") {
		t.Error("output missing 'info message'")
	}

	if !strings.Contains(output, "warn message") {
		t.Error("output missing 'warn message'")
	}

	if !strings.Contains(output, "error message") {
		t.Error("output missing 'error message'")
	}
}

func TestLoggerOutput_JSON(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	logger.Info("json test", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "json test") {
		t.Error("output missing 'json test'")
	}

	if !strings.Contains(output, "key") {
		t.Error("output missing 'key'")
	}

	if !strings.Contains(output, "value") {
		t.Error("output missing 'value'")
	}
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
			if !strings.Contains(output, tt.shouldContain) {
				t.Errorf("output should contain %q, but doesn't", tt.shouldContain)
			}

			if strings.Contains(output, tt.shouldNotContain) {
				t.Errorf("output should not contain %q, but does", tt.shouldNotContain)
			}
		})
	}
}
