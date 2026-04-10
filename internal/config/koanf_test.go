package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKoanfLoader_LoadDefaults(t *testing.T) {
	t.Parallel()

	loader := NewLoader()

	err := loader.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := loader.GetString("log_level"); got != "info" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "info")
	}

	if got := loader.GetString("log_format"); got != "text" {
		t.Errorf("GetString(log_format) = %q, want %q", got, "text")
	}

	if loader.GetBool("strict_mode") {
		t.Errorf("GetBool(strict_mode) = true, want false")
	}
}

func TestKoanfLoader_LoadEnv(t *testing.T) {
	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")
	t.Setenv("CMDGUARD_LOG_FORMAT", "json")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()

	err := loader.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := loader.GetString("log_level"); got != "debug" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "debug")
	}

	if got := loader.GetString("log_format"); got != "json" {
		t.Errorf("GetString(log_format) = %q, want %q", got, "json")
	}

	if !loader.GetBool("strict_mode") {
		t.Errorf("GetBool(strict_mode) = false, want true")
	}
}

func TestKoanfLoader_LoadFile(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
log_level: warn
log_format: json
strict_mode: true
`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loader := NewLoader()

	err = loader.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := loader.GetString("log_level"); got != "warn" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "warn")
	}

	if got := loader.GetString("log_format"); got != "json" {
		t.Errorf("GetString(log_format) = %q, want %q", got, "json")
	}

	if !loader.GetBool("strict_mode") {
		t.Errorf("GetBool(strict_mode) = false, want true")
	}
}

func TestKoanfLoader_EnvOverridesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
log_level: warn
log_format: json
`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()

	err = loader.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := loader.GetString("log_level"); got != "debug" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "debug")
	}

	if got := loader.GetString("log_format"); got != "json" {
		t.Errorf("GetString(log_format) = %q, want %q", got, "json")
	}

	if !loader.GetBool("strict_mode") {
		t.Errorf("GetBool(strict_mode) = false, want true")
	}
}

func TestKoanfLoader_Unmarshal(t *testing.T) {
	t.Setenv("CMDGUARD_LOG_LEVEL", "error")
	t.Setenv("CMDGUARD_LOG_FORMAT", "json")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()

	err := loader.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type TestConfig struct {
		LogLevel   string `koanf:"log_level"`
		LogFormat  string `koanf:"log_format"`
		StrictMode bool   `koanf:"strict_mode"`
	}

	var cfg TestConfig

	err = loader.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LogLevel != "error" {
		t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, "error")
	}

	if cfg.LogFormat != "json" {
		t.Errorf("cfg.LogFormat = %q, want %q", cfg.LogFormat, "json")
	}

	if !cfg.StrictMode {
		t.Errorf("cfg.StrictMode = false, want true")
	}
}

func TestKoanfLoader_MissingFile(t *testing.T) {
	t.Parallel()

	loader := NewLoader()
	err := loader.Load("/nonexistent/path/config.yaml")
	// Should not error (file is optional)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have defaults
	if got := loader.GetString("log_level"); got != "info" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "info")
	}
}

func TestKoanfLoader_Priority(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `log_level: warn`

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")

	loader := NewLoader()

	err = loader.Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be debug (env wins)
	if got := loader.GetString("log_level"); got != "debug" {
		t.Errorf("GetString(log_level) = %q, want %q", got, "debug")
	}
}
