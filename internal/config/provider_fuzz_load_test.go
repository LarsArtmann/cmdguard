package config

import (
	"path/filepath"
	"strings"
	"testing"
)

// fuzzEnvVarLoad tests that Load() handles environment variable fuzzing safely.
// It sets the given env var to the fuzzed value and verifies Load() doesn't panic.
func fuzzEnvVarLoad(t *testing.T, envVar, value string) {
	t.Helper()

	if !isValidEnvValue(value) {
		return
	}

	t.Setenv(envVar, value)

	cfg := Load()
	if cfg == nil {
		t.Fatalf("Load() returned nil")
	}

	_ = cfg.Validate()
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
		fuzzEnvVarLoad(t, "CMDGUARD_LOG_LEVEL", value)
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
		fuzzEnvVarLoad(t, "CMDGUARD_LOG_FORMAT", value)
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

	invalidLevelTests := []struct {
		name  string
		level string
	}{
		{"null bytes in level", "debug\x00info"},
		{"control characters in level", "de\x01bug"},
	}

	for _, tt := range invalidLevelTests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{LogLevel: tt.level}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestGetConfigFilePath_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantEmpty bool
		check     func(t *testing.T, result string)
	}{
		{"path with null bytes", "config\x00.yaml", false, nil},
		{"path with newlines", "config\n.yaml", false, nil},
		{"very deep path traversal", strings.Repeat("../", 1000) + "etc/passwd", false, nil},
		{"unicode in path", "🎉-config.yaml", false, func(t *testing.T, result string) {
			if !strings.Contains(result, "🎉") {
				t.Errorf("result should contain 🎉, got %q", result)
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConfigFilePath(tt.input)
			if tt.wantEmpty && result != "" {
				t.Errorf("GetConfigFilePath(%q) = %q, want empty", tt.input, result)
			} else if !tt.wantEmpty && result == "" {
				t.Errorf("GetConfigFilePath(%q) = empty, want non-empty", tt.input)
			}
			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
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
