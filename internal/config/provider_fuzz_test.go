package config

import (
	"path/filepath"
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

func FuzzLoad_EnvVarLevel(f *testing.F) {
	corpus := []string{
		"debug", "info", "warn", "error",
		"DEBUG", "Debug", "DEbUG",
		" ", "  ",
		"invalid", "xyz",
		"debug\n", "debug\x00",
		strings.Repeat("a", 1000),
		"🎉",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if !isValidEnvValue(value) {
			return
		}

		t.Setenv(
			"CMDGUARD_LOG_LEVEL",
			value,
		)

		cfg := Load()
		if cfg == nil {
			t.Fatalf("Load() returned nil")
		}

		_ = cfg.Validate()
	})
}

func FuzzLoad_EnvVarFormat(f *testing.F) {
	corpus := []string{
		"text", "json",
		"TEXT", "Text", "JSON", "Json",
		" ", "  ",
		"xml", "yaml",
		"json\n", "json\x00",
		strings.Repeat("a", 1000),
		"🎉",
	}
	for _, s := range corpus {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if !isValidEnvValue(value) {
			return
		}

		t.Setenv(
			"CMDGUARD_LOG_FORMAT",
			value,
		)

		cfg := Load()
		if cfg == nil {
			t.Fatalf("Load() returned nil")
		}

		_ = cfg.Validate()
	})
}

func FuzzLoad_EnvVarStrictMode(f *testing.F) {
	corpus := []string{
		"true", "false",
		"TRUE", "True", "1", "0",
		"yes", "no", "", " ",
		"true\n", "true\x00",
		"true ", " true",
		"=true", "${STRICT}",
	}
	for _, value := range corpus {
		f.Add(value)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if !isValidEnvValue(value) {
			return
		}

		t.Setenv(
			"CMDGUARD_STRICT_MODE",
			value,
		)

		cfg := Load()
		if cfg == nil {
			t.Fatalf("Load() returned nil")
		}

		// Only "true" (exact lowercase match) should set StrictMode to true
		expectedStrict := value == "true"
		if cfg.StrictMode != expectedStrict {
			t.Errorf(
				"Load().StrictMode = %v, want %v for input %q",
				cfg.StrictMode,
				expectedStrict,
				value,
			)
		}
	})
}

func TestValidate_EdgeCases(t *testing.T) {
	t.Run("concurrent validation should be safe", func(t *testing.T) {
		cfg := &Config{LogLevel: "debug"}
		done := make(chan bool)

		for range 100 {
			go func() {
				err := cfg.Validate()
				if err != nil {
					t.Errorf("concurrent validation failed: %v", err)
				}

				done <- true
			}()
		}

		for range 100 {
			<-done
		}
	})

	t.Run("both fields invalid returns first error", func(t *testing.T) {
		cfg := &Config{LogLevel: "invalid", LogFormat: "xml"}

		err := cfg.Validate()
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "invalid log level") {
			t.Errorf("error should contain 'invalid log level', got %q", err.Error())
		}
	})

	t.Run("null bytes in level", func(t *testing.T) {
		cfg := &Config{LogLevel: "debug\x00info"}

		err := cfg.Validate()
		if err == nil {
			t.Errorf("expected error for null bytes in level, got nil")
		}
	})

	t.Run("control characters in level", func(t *testing.T) {
		cfg := &Config{LogLevel: "de\x01bug"}

		err := cfg.Validate()
		if err == nil {
			t.Errorf("expected error for control characters in level, got nil")
		}
	})
}

func TestGetConfigFilePath_EdgeCases(t *testing.T) {
	t.Run("path with null bytes", func(t *testing.T) {
		result := GetConfigFilePath("config\x00.yaml")
		if result == "" {
			t.Errorf("GetConfigFilePath(%q) = empty, got %q", "config\x00.yaml", result)
		}
	})

	t.Run("path with newlines", func(t *testing.T) {
		result := GetConfigFilePath("config\n.yaml")
		if result == "" {
			t.Errorf("GetConfigFilePath(%q) = empty, got %q", "config\n.yaml", result)
		}
	})

	t.Run("very deep path traversal", func(t *testing.T) {
		deepPath := strings.Repeat("../", 1000) + "etc/passwd"

		result := GetConfigFilePath(deepPath)
		if result == "" {
			t.Errorf("GetConfigFilePath(%q) = empty, got %q", deepPath, result)
		}
	})

	t.Run("unicode in path", func(t *testing.T) {
		result := GetConfigFilePath("🎉-config.yaml")
		if result == "" {
			t.Errorf("GetConfigFilePath(%q) = empty, got %q", "🎉-config.yaml", result)
		}

		if !strings.Contains(result, "🎉") {
			t.Errorf("result should contain 🎉, got %q", result)
		}
	})
}

func testShellInjectionPayload(t *testing.T, payload string) {
	t.Helper()

	t.Setenv("CMDGUARD_LOG_LEVEL", payload)

	cfg := Load()
	if cfg == nil {
		t.Fatalf("Load() returned nil")
	}

	if cfg.LogLevel != payload {
		t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, payload)
	}

	err := cfg.Validate()
	if err == nil {
		t.Errorf("expected validation error for payload %q, got nil", payload)
	}
}

func TestLoad_EnvVarInjection(t *testing.T) {
	t.Run("shell injection attempt in level", func(t *testing.T) {
		testShellInjectionPayload(t, "$(whoami)")
	})

	t.Run("backtick injection attempt", func(t *testing.T) {
		testShellInjectionPayload(t, "`id`")
	})

	t.Run("pipe injection attempt", func(t *testing.T) {
		testShellInjectionPayload(t, "debug|cat /etc/passwd")
	})
}

func FuzzKoanfLoader_LoadEnv(f *testing.F) {
	corpus := []struct {
		level  string
		format string
	}{
		{"debug", "text"},
		{"info", "json"},
		{"warn", "text"},
		{"error", "json"},
		{"", ""},
		{"invalid", "xml"},
	}
	for _, tt := range corpus {
		f.Add(tt.level, tt.format)
	}

	f.Fuzz(func(t *testing.T, level, format string) {
		if level != "" && !isValidEnvValue(level) {
			return
		}

		if format != "" && !isValidEnvValue(format) {
			return
		}

		if level != "" {
			t.Setenv("CMDGUARD_LOG_LEVEL", level)
		}

		if format != "" {
			t.Setenv("CMDGUARD_LOG_FORMAT", format)
		}

		loader := NewLoader()

		err := loader.Load("")
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
	})
}

func FuzzKoanfLoader_FilePath(f *testing.F) {
	corpus := []string{
		"",
		"config.yaml",
		"/etc/app/config.yaml",
		"./config.yaml",
		"../config.yaml",
		filepath.Join(strings.Repeat("../", 50), "config.yaml"),
		"🎉.yaml",
		"<script>.yaml",
		strings.Repeat("a", 200) + ".yaml",
	}
	for _, path := range corpus {
		f.Add(path)
	}

	f.Fuzz(func(t *testing.T, configPath string) {
		loader := NewLoader()
		err := loader.Load(configPath)

		// Directories should return an error (can't read a directory as a file)
		if configPath == "." || configPath == ".." || configPath == "/" {
			if err == nil {
				t.Errorf("Load(%q) expected error for directory, got nil", configPath)
			}

			return
		}

		// Non-existent files are OK (config is optional)
		// But invalid paths might error depending on the OS
		// We just verify the function doesn't panic
		_ = err
	})
}
