// Package config provides configuration management using koanf.
//
// This implementation uses knadh/koanf for configuration loading,
// supporting multiple sources: defaults, config files, and environment variables.
//
// Configuration priority (highest to lowest):
// 1. Environment variables (CMDGUARD_*)
// 2. Config file (YAML)
// 3. Default values
package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Loader handles configuration loading from multiple sources.
type Loader struct {
	k *koanf.Koanf
}

// NewLoader creates a new configuration loader.
func NewLoader() *Loader {
	return &Loader{
		k: koanf.New("."),
	}
}

// Load loads configuration from all sources with the following priority:
// 1. Environment variables (highest priority)
// 2. Config file (if provided)
// 3. Default values (lowest priority)
func (l *Loader) Load(configPath string) error {
	// 1. Load defaults
	if err := l.loadDefaults(); err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	// 2. Load config file (if path provided)
	if configPath != "" {
		if err := l.loadFile(configPath); err != nil {
			// Config file is optional, don't fail if not found
			if !strings.Contains(err.Error(), "no such file") {
				return fmt.Errorf("failed to load config file: %w", err)
			}
		}
	}

	// 3. Load environment variables (highest priority)
	if err := l.loadEnv(); err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	return nil
}

// loadDefaults loads default configuration values.
func (l *Loader) loadDefaults() error {
	return l.k.Load(confmap.Provider(map[string]any{
		"strict_mode": false,
		"log_level":   "info",
		"log_format":  "text",
	}, "."), nil)
}

// loadFile loads configuration from a YAML file.
func (l *Loader) loadFile(path string) error {
	return l.k.Load(file.Provider(path), yaml.Parser())
}

// loadEnv loads configuration from environment variables.
// Only variables with CMDGUARD_ prefix are loaded.
// Example: CMDGUARD_LOG_LEVEL=debug becomes log_level=debug
func (l *Loader) loadEnv() error {
	return l.k.Load(env.Provider(".", env.Opt{
		Prefix: "CMDGUARD_",
		TransformFunc: func(key, value string) (string, any) {
			// Transform: CMDGUARD_LOG_LEVEL -> log_level
			key = strings.ToLower(strings.TrimPrefix(key, "CMDGUARD_"))
			return key, value
		},
	}), nil)
}

// Unmarshal unmarshals the configuration into the provided struct.
// The struct should have `koanf` tags for field mapping.
func (l *Loader) Unmarshal(dest any) error {
	return l.k.Unmarshal("", dest)
}

// UnmarshalWithConf unmarshals with custom configuration.
func (l *Loader) UnmarshalWithConf(path string, dest any, conf koanf.UnmarshalConf) error {
	return l.k.UnmarshalWithConf(path, dest, conf)
}

// GetString returns a string value from the configuration.
func (l *Loader) GetString(key string) string {
	return l.k.String(key)
}

// GetBool returns a boolean value from the configuration.
func (l *Loader) GetBool(key string) bool {
	return l.k.Bool(key)
}

// GetInt returns an integer value from the configuration.
func (l *Loader) GetInt(key string) int {
	return l.k.Int(key)
}

// Print prints the configuration for debugging.
func (l *Loader) Print() {
	l.k.Print()
}
