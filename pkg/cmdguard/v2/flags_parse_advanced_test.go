package v2

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFlagRegistry_ParseFlags_Advanced(t *testing.T) {
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
