package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzParseLevel(f *testing.F) {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range validLevels {
		f.Add(level)
	}

	corpus := []string{
		"", " ", "  ", "\t", "\n",
		"DEBUG", "Debug", "DEbUG",
		"debug ", " debug", " debug ",
		"debug\x00", "debug\n", "debug\t",
		"invalid", "xyz", "123",
		strings.Repeat("a", 1000),
		strings.Repeat("debug", 100),
		"🎉", "-debug-", "debug_info",
		"<script>alert('xss')</script>",
		"'; DROP TABLE logs; --",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, level string) {
		result := ParseLevel(level)
		slogLevel := result.SlogLevel()
		str := result.String()

		assert.Contains(t, []Level{LevelDebug, LevelInfo, LevelWarn, LevelError}, result)
		assert.Contains(t, []string{"debug", "info", "warn", "error"}, str)
		assert.Contains(t, []int{-4, 0, 4, 8}, int(slogLevel))
	})
}

func FuzzParseFormat(f *testing.F) {
	validFormats := []string{"text", "json"}
	for _, format := range validFormats {
		f.Add(format)
	}

	corpus := []string{
		"", " ", "  ", "\t", "\n",
		"TEXT", "Text", "TeXt",
		"JSON", "Json", "JsOn",
		"text ", " text", " text ",
		"xml", "yaml", "toml",
		strings.Repeat("a", 1000),
		"🎉", "-text-", "json_xml",
		"<script>alert('xss')</script>",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, format string) {
		result := ParseFormat(format)
		str := result.String()

		assert.Contains(t, []Format{FormatText, FormatJSON}, result)
		assert.Contains(t, []string{"text", "json"}, str)
	})
}

func FuzzValidLevel(f *testing.F) {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range validLevels {
		f.Add(level, true)
	}

	corpus := []string{
		"", " ", "DEBUG", "Debug",
		"invalid", "xyz", "123",
		strings.Repeat("a", 1000),
		"🎉", "debug\x00",
	}
	for _, s := range corpus {
		f.Add(s, false)
	}

	f.Fuzz(func(t *testing.T, level string, expectValid bool) {
		result := ValidLevel(level)
		assert.Equal(t, expectValid, result)
	})
}

func FuzzValidFormat(f *testing.F) {
	validFormats := []string{"text", "json"}
	for _, format := range validFormats {
		f.Add(format, true)
	}

	corpus := []string{
		"", " ", "TEXT", "JSON",
		"xml", "yaml", "invalid",
		strings.Repeat("a", 1000),
		"🎉", "json\x00",
	}
	for _, s := range corpus {
		f.Add(s, false)
	}

	f.Fuzz(func(t *testing.T, format string, expectValid bool) {
		result := ValidFormat(format)
		assert.Equal(t, expectValid, result)
	})
}

func TestNewLogger_JSONOutputIsValid(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	logger.Info("test message", "key", "value", "count", 42)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 1, "should produce exactly one JSON line")

	var parsed map[string]any
	err := json.Unmarshal([]byte(lines[0]), &parsed)
	require.NoError(t, err, "output should be valid JSON")

	assert.Equal(t, "test message", parsed["msg"])
	assert.Equal(t, "value", parsed["key"])
	assert.Equal(t, float64(42), parsed["count"])
	assert.Contains(t, parsed, "time")
	assert.Contains(t, parsed, "level")
}

func TestNewLogger_TextOutputFormat(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	logger.Info("test message", "key", "value")

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key=value")
	assert.Contains(t, output, "level=INFO")
}

func TestLevel_CaseSensitivity(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelInfo},
		{"Debug", LevelInfo},
		{"DEbUG", LevelInfo},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"Info", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelInfo},
		{"error", LevelError},
		{"ERROR", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			assert.Equal(t, tt.expected, result, "ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}

func TestFormat_CaseSensitivity(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{"Text", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatText},
		{"Json", FormatText},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseFormat(tt.input)
			assert.Equal(t, tt.expected, result, "ParseFormat(%q) = %v, want %v", tt.input, result, tt.expected)
		})
	}
}

func TestLevel_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{" debug", LevelInfo},
		{"debug ", LevelInfo},
		{" debug ", LevelInfo},
		{"\tdebug", LevelInfo},
		{"debug\n", LevelInfo},
		{"", LevelInfo},
		{"   ", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidLevel_EdgeCases(t *testing.T) {
	assert.False(t, ValidLevel(" debug"), "leading space should be invalid")
	assert.False(t, ValidLevel("debug "), "trailing space should be invalid")
	assert.False(t, ValidLevel("DEBUG"), "uppercase should be invalid")
	assert.False(t, ValidLevel(""), "empty should be invalid")
	assert.False(t, ValidLevel("   "), "whitespace only should be invalid")
	assert.False(t, ValidLevel("debug\x00"), "null byte should be invalid")
	assert.False(t, ValidLevel("🎉"), "emoji should be invalid")
	assert.False(t, ValidLevel(strings.Repeat("a", 10000)), "very long string should be invalid")
}

func TestValidFormat_EdgeCases(t *testing.T) {
	assert.False(t, ValidFormat(" json"), "leading space should be invalid")
	assert.False(t, ValidFormat("json "), "trailing space should be invalid")
	assert.False(t, ValidFormat("JSON"), "uppercase should be invalid")
	assert.False(t, ValidFormat(""), "empty should be invalid")
	assert.False(t, ValidFormat("   "), "whitespace only should be invalid")
	assert.False(t, ValidFormat(strings.Repeat("a", 10000)), "very long string should be invalid")
}
