package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"slices"
	"strings"
	"testing"
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

		validLevels := []Level{LevelDebug, LevelInfo, LevelWarn, LevelError}
		found := slices.Contains(validLevels, result)

		if !found {
			t.Errorf("ParseLevel(%q) = %v, expected valid level", level, result)
		}

		validStrings := []string{"debug", "info", "warn", "error"}
		found = slices.Contains(validStrings, str)

		if !found {
			t.Errorf("ParseLevel(%q).String() = %q, expected valid string", level, str)
		}

		validInts := []int{-4, 0, 4, 8}
		found = slices.Contains(validInts, int(slogLevel))

		if !found {
			t.Errorf("ParseLevel(%q).SlogLevel() = %d, expected valid level", level, int(slogLevel))
		}
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

		validFormats := []Format{FormatText, FormatJSON}
		found := slices.Contains(validFormats, result)

		if !found {
			t.Errorf("ParseFormat(%q) = %v, expected valid format", format, result)
		}

		validStrings := []string{"text", "json"}
		found = slices.Contains(validStrings, str)

		if !found {
			t.Errorf("ParseFormat(%q).String() = %q, expected valid string", format, str)
		}
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
		if result != expectValid {
			t.Errorf("ValidLevel(%q) = %v, want %v", level, result, expectValid)
		}
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
		if result != expectValid {
			t.Errorf("ValidFormat(%q) = %v, want %v", format, result, expectValid)
		}
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
	if len(lines) != 1 {
		t.Fatalf("should produce exactly one JSON line, got %d lines", len(lines))
	}

	var parsed map[string]any

	err := json.Unmarshal([]byte(lines[0]), &parsed)
	if err != nil {
		t.Fatalf("output should be valid JSON: %v", err)
	}

	if parsed["msg"] != "test message" {
		t.Errorf("parsed[\"msg\"] = %v, want %q", parsed["msg"], "test message")
	}

	if parsed["key"] != "value" {
		t.Errorf("parsed[\"key\"] = %v, want %q", parsed["key"], "value")
	}

	if parsed["count"] != float64(42) {
		t.Errorf("parsed[\"count\"] = %v, want %v", parsed["count"], float64(42))
	}

	if _, ok := parsed["time"]; !ok {
		t.Error("parsed output missing 'time' field")
	}

	if _, ok := parsed["level"]; !ok {
		t.Error("parsed output missing 'level' field")
	}
}

func TestNewLogger_TextOutputFormat(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	logger.Info("test message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("output missing 'test message'")
	}

	if !strings.Contains(output, "key=value") {
		t.Error("output missing 'key=value'")
	}

	if !strings.Contains(output, "level=INFO") {
		t.Error("output missing 'level=INFO'")
	}
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
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
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
			if result != tt.expected {
				t.Errorf("ParseFormat(%q) = %v, want %v", tt.input, result, tt.expected)
			}
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
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidLevel_EdgeCases(t *testing.T) {
	if ValidLevel(" debug") {
		t.Error("leading space should be invalid")
	}

	if ValidLevel("debug ") {
		t.Error("trailing space should be invalid")
	}

	if ValidLevel("DEBUG") {
		t.Error("uppercase should be invalid")
	}

	if ValidLevel("") {
		t.Error("empty should be invalid")
	}

	if ValidLevel("   ") {
		t.Error("whitespace only should be invalid")
	}

	if ValidLevel("debug\x00") {
		t.Error("null byte should be invalid")
	}

	if ValidLevel("🎉") {
		t.Error("emoji should be invalid")
	}

	if ValidLevel(strings.Repeat("a", 10000)) {
		t.Error("very long string should be invalid")
	}
}

func TestValidFormat_EdgeCases(t *testing.T) {
	if ValidFormat(" json") {
		t.Error("leading space should be invalid")
	}

	if ValidFormat("json ") {
		t.Error("trailing space should be invalid")
	}

	if ValidFormat("JSON") {
		t.Error("uppercase should be invalid")
	}

	if ValidFormat("") {
		t.Error("empty should be invalid")
	}

	if ValidFormat("   ") {
		t.Error("whitespace only should be invalid")
	}

	if ValidFormat(strings.Repeat("a", 10000)) {
		t.Error("very long string should be invalid")
	}
}
