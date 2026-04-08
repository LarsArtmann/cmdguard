package logging

import (
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			logger := NewLogger(tt.format, tt.level)
			if logger == nil {
				t.Error("NewLogger() returned nil, expected non-nil logger")
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			result := ParseFormat(tt.format)
			if result != tt.expectedFormat {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.format, result, tt.expectedFormat)
			}
		})
	}
}

func TestValidFormat(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := ValidFormat(tt.format)
			if result != tt.isValid {
				t.Errorf("ValidFormat(%q) = %v, want %v", tt.format, result, tt.isValid)
			}
		})
	}
}

func TestValidLevel(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := ValidLevel(tt.level)
			if result != tt.isValid {
				t.Errorf("ValidLevel(%q) = %v, want %v", tt.level, result, tt.isValid)
			}
		})
	}
}
