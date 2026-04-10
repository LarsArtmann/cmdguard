package config

import (
	"strings"
	"testing"
)

func FuzzValidate_LogLevel(f *testing.F) {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range validLevels {
		f.Add(level, true)
	}

	corpus := []struct {
		level       string
		expectValid bool
	}{
		{"", true},
		{" ", false},
		{"DEBUG", false},
		{"invalid", false},
		{"xyz", false},
		{"123", false},
		{"debug\x00info", false},
		{"🎉", false},
		{strings.Repeat("a", 10000), false},
		{"<script>alert('xss')</script>", false},
		{"'; DROP TABLE logs; --", false},
		{"debug\ninfo", false},
		{"debug\tinfo", false},
		{"${LOG_LEVEL}", false},
		{"$LOG_LEVEL", false},
		{"./debug", false},
		{"../debug", false},
	}

	for _, tt := range corpus {
		f.Add(tt.level, tt.expectValid)
	}

	f.Fuzz(func(t *testing.T, level string, expectValid bool) {
		cfg := &Config{LogLevel: level}
		err := cfg.Validate()

		if expectValid {
			if err != nil && !strings.Contains(err.Error(), "log level") {
				t.Errorf("unexpected error for level %q: %v", level, err)
			}
		} else {
			// For invalid levels, we expect an error about log level
			// unless it's empty (which is valid)
			if level != "" && !strings.Contains(level, "debug") &&
				!strings.Contains(level, "info") && !strings.Contains(level, "warn") &&
				!strings.Contains(level, "error") {
				if err == nil {
					t.Errorf("expected error for invalid level %q, got nil", level)
				}
			}
		}
	})
}

func FuzzValidate_LogFormat(f *testing.F) {
	validFormats := []string{"text", "json"}
	for _, format := range validFormats {
		f.Add(format, true)
	}

	corpus := []struct {
		format      string
		expectValid bool
	}{
		{"", true},
		{" ", false},
		{"xml", false},
		{"yaml", false},
		{"toml", false},
		{"text ", false},
		{"text/json", false},
		{"text+json", false},
		{"${FORMAT}", false},
		{"$FORMAT", false},
		{"./format", false},
		{"../format", false},
		{"http://format.json", false},
		{"text json", false},
		{"text\tjson", false},
	}

	for _, tt := range corpus {
		f.Add(tt.format, tt.expectValid)
	}

	f.Fuzz(func(t *testing.T, format string, expectValid bool) {
		cfg := &Config{LogFormat: format}
		err := cfg.Validate()

		if expectValid {
			if err != nil && !strings.Contains(err.Error(), "log format") {
				t.Errorf("unexpected error for format %q: %v", format, err)
			}
		} else {
			if format != "" && format != "text" && format != "json" {
				if err == nil {
					t.Errorf("expected error for invalid format %q, got nil", format)
				}
			}
		}
	})
}

func FuzzGetConfigFilePath(f *testing.F) {
	corpus := []string{
		"",
		"config.yaml",
		"/etc/cmdguard/config.yaml",
		"./config.yaml",
		"../config.yaml",
		strings.Repeat("../", 100) + ".yaml",
		"a/b/c/config.yaml",
		strings.Repeat("a", 1000) + ".yaml",
		"🎉.yaml",
		"<script>alert('xss')</script>.yaml",
		"${HOME}/config.yaml",
		"~/config.yaml",
		"config with spaces.yaml",
		"config\twith\ttabs.yaml",
	}

	for _, path := range corpus {
		f.Add(path)
	}

	f.Fuzz(func(t *testing.T, configFile string) {
		result := GetConfigFilePath(configFile)
		if configFile == "" {
			if result != "" {
				t.Errorf("GetConfigFilePath(%q) = %q, want empty", configFile, result)
			}
		} else {
			// Non-empty input should produce some result
			// We can't assert much about absolute paths without knowing the cwd
			_ = result
		}
	})
}

// isValidEnvValue checks if a string is valid for use as an environment variable value.
// Environment variables cannot contain null bytes or newlines on some platforms (e.g., Darwin).
func isValidEnvValue(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\x00', '\n', '\r':
			return false
		}
	}

	return true
}
