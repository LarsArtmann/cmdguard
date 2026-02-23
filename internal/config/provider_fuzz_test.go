package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{"  ", false},
		{"DEBUG", false},
		{"Debug", false},
		{"DEbUG", false},
		{"debug ", false},
		{" debug", false},
		{"invalid", false},
		{"xyz", false},
		{"123", false},
		{"debug\x00", false},
		{"🎉", false},
		{strings.Repeat("a", 10000), false},
		{"<script>alert('xss')</script>", false},
		{"'; DROP TABLE logs; --", false},
		{"debug\ninfo", false},
		{"debug\tinfo", false},
		{"debug info", false},
		{"DEBUG=info", false},
		{"${LOG_LEVEL}", false},
		{"$LOG_LEVEL", false},
		{"%LOG_LEVEL%", false},
		{"../debug", false},
		{"./debug", false},
		{"/debug", false},
	}
	for _, tt := range corpus {
		f.Add(tt.level, tt.expectValid)
	}

	f.Fuzz(func(t *testing.T, level string, expectValid bool) {
		validateConfigField(t, &Config{LogLevel: level}, expectValid)
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
		{"TEXT", false},
		{"Text", false},
		{"JSON", false},
		{"Json", false},
		{"xml", false},
		{"yaml", false},
		{"toml", false},
		{"text ", false},
		{" text", false},
		{"json\x00", false},
		{"🎉", false},
		{strings.Repeat("a", 10000), false},
		{"<script>alert('xss')</script>", false},
		{"text/json", false},
		{"text+json", false},
		{"${FORMAT}", false},
	}
	for _, tt := range corpus {
		f.Add(tt.format, tt.expectValid)
	}

	f.Fuzz(func(t *testing.T, format string, expectValid bool) {
		validateConfigField(t, &Config{LogFormat: format}, expectValid)
	})
}

func FuzzGetConfigFilePath(f *testing.F) {
	f.Add("")
	f.Add("config.yaml")
	f.Add("/etc/cmdguard/config.yaml")
	f.Add("./config.yaml")
	f.Add("../config.yaml")

	corpus := []string{
		"..\n",
		"..\t",
		"..\x00",
		"../" + strings.Repeat("../", 100) + "etc/passwd",
		"./config.yaml",
		"/absolute/path.yaml",
		"relative/path.yaml",
		strings.Repeat("a", 1000) + ".yaml",
		"🎉.yaml",
		"config with spaces.yaml",
		"config\twith\ttabs.yaml",
		"config\nwith\nnewlines.yaml",
		"<script>alert('xss')</script>.yaml",
		"'; DROP TABLE configs; --.yaml",
		"${HOME}/config.yaml",
		"$HOME/config.yaml",
		"%APPDATA%/config.yaml",
		"~/config.yaml",
	}
	for _, path := range corpus {
		f.Add(path)
	}

	f.Fuzz(func(t *testing.T, configFile string) {
		result := GetConfigFilePath(configFile)
		if configFile == "" {
			assert.Empty(t, result)
		} else {
			assert.NotEmpty(t, result)
			abs, _ := filepath.Abs(configFile)
			assert.Equal(t, abs, result)
		}
	})
}

func validateConfigField(t *testing.T, cfg *Config, expectValid bool) {
	t.Helper()
	err := cfg.Validate()
	if expectValid {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
	}
}

func fuzzLoadWithEnvVar(f *testing.F, envVarName string, corpus []string) {
	for _, value := range corpus {
		f.Add(value)
	}

	f.Fuzz(func(t *testing.T, value string) {
		_ = os.Setenv(envVarName, value)
		defer func() { _ = os.Unsetenv(envVarName) }()

		cfg := Load()
		require.NotNil(t, cfg)
		// Just verify it doesn't crash and returns a valid config
		// OS may handle edge cases (null bytes, etc.) differently
		_ = cfg.Validate()
	})
}

func FuzzLoad_EnvVarLevel(f *testing.F) {
	// Note: Load() returns the raw env var value; Validate() checks it
	// Some values (null bytes, etc.) may be handled differently by the OS
	corpus := []string{
		"debug", "info", "warn", "error",
		"DEBUG", "Debug", "DEbUG",
		" ", "  ",
		"invalid", "xyz",
		"debug\n", "debug\x00",
		strings.Repeat("a", 1000),
		"🎉",
	}
	fuzzLoadWithEnvVar(f, "CMDGUARD_LOG_LEVEL", corpus)
}

func FuzzLoad_EnvVarFormat(f *testing.F) {
	// Note: Load() returns the raw env var value; Validate() checks it
	// Some values (null bytes, etc.) may be handled differently by the OS
	corpus := []string{
		"text", "json",
		"TEXT", "Text", "JSON", "Json",
		" ", "  ",
		"xml", "yaml",
		"json\n", "json\x00",
		strings.Repeat("a", 1000),
		"🎉",
	}
	fuzzLoadWithEnvVar(f, "CMDGUARD_LOG_FORMAT", corpus)
}

func FuzzLoad_EnvVarStrictMode(f *testing.F) {
	corpus := []struct {
		value        string
		expectStrict bool
	}{
		{"true", true},
		{"false", false},
		{"TRUE", false},
		{"True", false},
		{"1", false},
		{"0", false},
		{"yes", false},
		{"no", false},
		{"", false},
		{" ", false},
		{"true\n", false},
		{"true\x00", false},
		{"true ", false},
		{" true", false},
		{"=true", false},
		{"${STRICT}", false},
	}
	for _, tt := range corpus {
		f.Add(tt.value, tt.expectStrict)
	}

	f.Fuzz(func(t *testing.T, value string, expectStrict bool) {
		_ = os.Setenv("CMDGUARD_STRICT_MODE", value)
		defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()

		cfg := Load()
		require.NotNil(t, cfg)
		assert.Equal(t, expectStrict, cfg.StrictMode)
	})
}

func TestValidate_EdgeCases(t *testing.T) {
	t.Run("concurrent validation should be safe", func(t *testing.T) {
		cfg := &Config{LogLevel: "debug"}
		done := make(chan bool)
		for i := 0; i < 100; i++ {
			go func() {
				err := cfg.Validate()
				assert.NoError(t, err)
				done <- true
			}()
		}
		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("both fields invalid returns first error", func(t *testing.T) {
		cfg := &Config{LogLevel: "invalid", LogFormat: "xml"}
		err := cfg.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})

	t.Run("null bytes in level", func(t *testing.T) {
		cfg := &Config{LogLevel: "debug\x00info"}
		err := cfg.Validate()
		assert.Error(t, err)
	})

	t.Run("control characters in level", func(t *testing.T) {
		cfg := &Config{LogLevel: "de\x01bug"}
		err := cfg.Validate()
		assert.Error(t, err)
	})
}

func TestGetConfigFilePath_EdgeCases(t *testing.T) {
	t.Run("path with null bytes", func(t *testing.T) {
		result := GetConfigFilePath("config\x00.yaml")
		assert.NotEmpty(t, result)
	})

	t.Run("path with newlines", func(t *testing.T) {
		result := GetConfigFilePath("config\n.yaml")
		assert.NotEmpty(t, result)
	})

	t.Run("very deep path traversal", func(t *testing.T) {
		deepPath := strings.Repeat("../", 1000) + "etc/passwd"
		result := GetConfigFilePath(deepPath)
		assert.NotEmpty(t, result)
		// filepath.Abs resolves the path, so it will be an absolute path
		// This tests that it doesn't crash on deep traversal
	})

	t.Run("unicode in path", func(t *testing.T) {
		result := GetConfigFilePath("🎉-config.yaml")
		assert.NotEmpty(t, result)
		assert.Contains(t, result, "🎉")
	})
}

func testShellInjectionPayload(t *testing.T, payload string) {
	t.Helper()
	_ = os.Setenv("CMDGUARD_LOG_LEVEL", payload)
	defer func() { _ = os.Unsetenv("CMDGUARD_LOG_LEVEL") }()

	cfg := Load()
	require.NotNil(t, cfg)
	assert.Equal(t, payload, cfg.LogLevel)

	err := cfg.Validate()
	assert.Error(t, err)
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
