// Package config provides Koanf-based configuration management for cmdguard.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/samber/do/v2"
	"github.com/spf13/cobra"
)

// Config holds the application configuration.
type Config struct {
	StrictMode bool   `koanf:"strict_mode"`
	ConfigFile string `koanf:"config_file"`
	LogLevel   string `koanf:"log_level"`
}

// NewConfig creates a new Config instance loaded from files, env vars, and flags.
// This is registered as an eager service to ensure config is loaded immediately.
func NewConfig(i do.Injector) (*Config, error) {
	k := koanf.New(".")

	// Try to load from default config file first
	_ = k.Load(file.Provider("config.yaml"), yaml.Parser())

	// Load from environment variables with CMDGUARD_ prefix
	_ = k.Load(env.Provider("CMDGUARD_", ".", nil), nil)

	// Check for custom config file via env
	if configFile := os.Getenv("CMDGUARD_CONFIG_FILE"); configFile != "" {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("failed to load config file %q: %w", configFile, err)
		}
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return &cfg, nil
}

// NewConfigWithCommand creates a Config that also loads from cobra flags.
// This should be called after cobra flags are parsed.
func NewConfigWithCommand(i do.Injector, cmd *cobra.Command) (*Config, error) {
	k := koanf.New(".")

	// Try to load from default config file first
	_ = k.Load(file.Provider("config.yaml"), yaml.Parser())

	// Load from environment variables
	_ = k.Load(env.Provider("CMDGUARD_", ".", nil), nil)

	// Load from cobra flags if available
	if cmd != nil {
		// Get config file from flag if provided
		configFile, _ := cmd.Flags().GetString("config")
		if configFile != "" {
			if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
				return nil, fmt.Errorf("failed to load config file %q: %w", configFile, err)
			}
		}

		// Load remaining flags
		if err := k.Load(posflagProvider(cmd, ".", k), nil); err != nil {
			return nil, fmt.Errorf("failed to load flags: %w", err)
		}
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return &cfg, nil
}

// posflagProvider adapts cobra pflag to koanf provider.
func posflagProvider(cmd *cobra.Command, delim string, k *koanf.Koanf) koanf.Provider {
	return posflag.Provider(cmd.Flags(), delim, k)
}

// Shutdown implements the Shutdowner interface for graceful cleanup.
func (c *Config) Shutdown() error {
	// No cleanup needed for config
	return nil
}

// Validate performs validation on the configuration.
func (c *Config) Validate() error {
	if c.LogLevel != "" {
		validLevels := []string{"debug", "info", "warn", "error"}
		found := false
		for _, level := range validLevels {
			if c.LogLevel == level {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid log level %q, must be one of: debug, info, warn, error", c.LogLevel)
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
