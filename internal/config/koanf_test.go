package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKoanfLoader_LoadDefaults(t *testing.T) {
	loader := NewLoader()
	err := loader.Load("")
	require.NoError(t, err)

	// Verify defaults
	assert.Equal(t, "info", loader.GetString("log_level"))
	assert.Equal(t, "text", loader.GetString("log_format"))
	assert.Equal(t, false, loader.GetBool("strict_mode"))
}

func TestKoanfLoader_LoadEnv(t *testing.T) {
	// Set environment variables
	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")
	t.Setenv("CMDGUARD_LOG_FORMAT", "json")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()
	err := loader.Load("")
	require.NoError(t, err)

	// Verify env overrides defaults
	assert.Equal(t, "debug", loader.GetString("log_level"))
	assert.Equal(t, "json", loader.GetString("log_format"))
	assert.Equal(t, true, loader.GetBool("strict_mode"))
}

func TestKoanfLoader_LoadFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
log_level: warn
log_format: json
strict_mode: true
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	loader := NewLoader()
	err = loader.Load(configPath)
	require.NoError(t, err)

	// Verify file overrides defaults
	assert.Equal(t, "warn", loader.GetString("log_level"))
	assert.Equal(t, "json", loader.GetString("log_format"))
	assert.Equal(t, true, loader.GetBool("strict_mode"))
}

func TestKoanfLoader_EnvOverridesFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
log_level: warn
log_format: json
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	// Set environment variables
	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()
	err = loader.Load(configPath)
	require.NoError(t, err)

	// Verify env overrides file
	assert.Equal(t, "debug", loader.GetString("log_level"))
	assert.Equal(t, "json", loader.GetString("log_format")) // From file
	assert.Equal(t, true, loader.GetBool("strict_mode"))    // From env
}

func TestKoanfLoader_Unmarshal(t *testing.T) {
	t.Setenv("CMDGUARD_LOG_LEVEL", "error")
	t.Setenv("CMDGUARD_LOG_FORMAT", "json")
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	loader := NewLoader()
	err := loader.Load("")
	require.NoError(t, err)

	type TestConfig struct {
		LogLevel   string `koanf:"log_level"`
		LogFormat  string `koanf:"log_format"`
		StrictMode bool   `koanf:"strict_mode"`
	}

	var cfg TestConfig
	err = loader.Unmarshal(&cfg)
	require.NoError(t, err)

	assert.Equal(t, "error", cfg.LogLevel)
	assert.Equal(t, "json", cfg.LogFormat)
	assert.Equal(t, true, cfg.StrictMode)
}

func TestKoanfLoader_MissingFile(t *testing.T) {
	// Try to load non-existent file
	loader := NewLoader()
	err := loader.Load("/nonexistent/path/config.yaml")

	// Should not error (file is optional)
	require.NoError(t, err)

	// Should have defaults
	assert.Equal(t, "info", loader.GetString("log_level"))
}

func TestKoanfLoader_Priority(t *testing.T) {
	// Priority test: env > file > defaults
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// File sets log_level to warn
	configContent := `log_level: warn`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	// Env sets log_level to debug
	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")

	loader := NewLoader()
	err = loader.Load(configPath)
	require.NoError(t, err)

	// Should be debug (env wins)
	assert.Equal(t, "debug", loader.GetString("log_level"))
}
