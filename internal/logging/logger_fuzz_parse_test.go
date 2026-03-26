package logging

import (
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
		"'; DROP TABLE logs; --",
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
	corpus := []string{
		"", " ", "DEBUG", "Debug",
		"invalid", "xyz", "123",
		strings.Repeat("a", 1000),
		"🎉", "debug\x00",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	validSet := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	f.Fuzz(func(t *testing.T, level string) {
		result := ValidLevel(level)
		expected := validSet[level]
		if result != expected {
			t.Errorf("ValidLevel(%q) = %v, want %v", level, result, expected)
		}
	})
}

func FuzzValidFormat(f *testing.F) {
	validFormats := []string{"text", "json"}
	for _, format := range validFormats {
		f.Add(format)
	}

	corpus := []string{
		"", " ", "TEXT", "JSON",
		"xml", "yaml", "invalid",
		strings.Repeat("a", 1000),
		"🎉", "json\x00",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	validSet := map[string]bool{
		"text": true,
		"json": true,
	}

	f.Fuzz(func(t *testing.T, format string) {
		result := ValidFormat(format)
		expected := validSet[format]
		if result != expected {
			t.Errorf("ValidFormat(%q) = %v, want %v", format, result, expected)
		}
	})
}
