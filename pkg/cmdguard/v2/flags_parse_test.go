package v2

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestFlagRegistry_ParseFlags(t *testing.T) {
	t.Run("parse string flag", func(t *testing.T) {
		type TestConfig struct {
			Name string `default:"default" flag:"name"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		// Set flag value
		if err := cmd.PersistentFlags().Set("name", "custom"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Name != "custom" {
			t.Errorf("expected Name 'custom', got %q", cfg.Name)
		}
	})

	t.Run("parse bool flag", func(t *testing.T) {
		type TestConfig struct {
			Verbose bool `default:"false" flag:"verbose"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("verbose", "true"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if !cfg.Verbose {
			t.Error("expected Verbose to be true")
		}
	})

	t.Run("parse int flag", func(t *testing.T) {
		type TestConfig struct {
			Count int `default:"0" flag:"count"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("count", "42"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Count != 42 {
			t.Errorf("expected Count 42, got %d", cfg.Count)
		}
	})

	t.Run("parse uint flag", func(t *testing.T) {
		type TestConfig struct {
			Workers uint `default:"0" flag:"workers"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("workers", "10"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Workers != 10 {
			t.Errorf("expected Workers 10, got %d", cfg.Workers)
		}
	})

	t.Run("parse float64 flag", func(t *testing.T) {
		type TestConfig struct {
			Rate float64 `default:"0.0" flag:"rate"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("rate", "3.14159"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		delta := 0.00001
		if cfg.Rate < 3.14159-delta || cfg.Rate > 3.14159+delta {
			t.Errorf("expected Rate ~3.14159, got %f", cfg.Rate)
		}
	})

	t.Run("parse Duration flag", func(t *testing.T) {
		type TestConfig struct {
			Timeout Duration `default:"1m" flag:"timeout"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("timeout", "5m30s"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := FromDuration(5*time.Minute + 30*time.Second)
		if cfg.Timeout != expected {
			t.Errorf("expected Timeout %v, got %v", expected, cfg.Timeout)
		}
	})

	t.Run("parse LogLevel flag valid", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("level", "debug"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Level.String() != "debug" {
			t.Errorf("expected Level 'debug', got %q", cfg.Level.String())
		}
	})

	t.Run("parse LogLevel flag invalid returns error", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("level", "invalid"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "level") {
			t.Errorf("expected error to contain 'level', got: %v", err)
		}
	})

	t.Run("parse Enum flag", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("mode", "prod"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Mode.String() != "prod" {
			t.Errorf("expected Mode 'prod', got %q", cfg.Mode.String())
		}
	})

	t.Run("parse invalid Enum returns error", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("mode", "invalid"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("parse LogFormat flag valid", func(t *testing.T) {
		type TestConfig struct {
			Format LogFormat `flag:"format"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("format", "json"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if cfg.Format.String() != "json" {
			t.Errorf("expected Format 'json', got %q", cfg.Format.String())
		}
	})

	t.Run("parse LogFormat flag invalid returns error", func(t *testing.T) {
		type TestConfig struct {
			Format LogFormat `flag:"format"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		if err := registry.RegisterFlags(cmd); err != nil {
			t.Fatalf("expected no error registering flags, got: %v", err)
		}

		if err := cmd.PersistentFlags().Set("format", "invalid"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "format") {
			t.Errorf("expected error to contain 'format', got: %v", err)
		}
	})
}

func TestFlagRegistry_FlagNotFound(t *testing.T) {
	t.Run("missing flag returns error", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Don't register flags on command
		cmd := &cobra.Command{Use: "test"}

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected error to contain 'not found', got: %v", err)
		}
	})
}
