// Package config provides configuration management for cmdguard.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// Config holds the application configuration.
type Config struct {
	StrictMode bool   `koanf:"strict_mode"`
	ConfigFile string `koanf:"config_file"`
	LogLevel   string `koanf:"log_level"`
	LogFormat  string `koanf:"log_format"`
}

// Load loads configuration from environment variables.
// Returns config with defaults if no environment variables are set.
func Load() *Config {
	cfg := &Config{
		LogLevel:   "info",
		LogFormat:  "text",
		StrictMode: false,
	}

	// Load from environment variables
	if level := os.Getenv("CMDGUARD_LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}
	if format := os.Getenv("CMDGUARD_LOG_FORMAT"); format != "" {
		cfg.LogFormat = format
	}
	if os.Getenv("CMDGUARD_STRICT_MODE") == "true" {
		cfg.StrictMode = true
	}

	return cfg
}

// Validate performs validation on the configuration.
func (c *Config) Validate() error {
	if c.LogLevel != "" {
		validLevels := []string{"debug", "info", "warn", "error"}
		if !slices.Contains(validLevels, c.LogLevel) {
			return fmt.Errorf("invalid log level %q, must be one of: debug, info, warn, error", c.LogLevel)
		}
	}
	if c.LogFormat != "" {
		validFormats := []string{"text", "json"}
		if !slices.Contains(validFormats, c.LogFormat) {
			return fmt.Errorf("invalid log format %q, must be one of: text, json", c.LogFormat)
		}
	}
	return nil
}

// GetConfigFilePath returns the absolute path to the config file.
func GetConfigFilePath(configFile string) string {
	if configFile == "" {
		return ""
	}
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		return configFile
	}
	return absPath
}
