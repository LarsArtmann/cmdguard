package v2

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFlagRegistry_ParseFlags_Advanced(t *testing.T) {
	t.Parallel()
	t.Run("parse LogLevel flag valid", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Level LogLevel `flag:"level"`
		}

		cfg := &TestConfig{}

		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		registerAndSetFlag(t, registry, cmd, cfg, "level", "debug")

		assertEnumString(t, cfg.Level.String(), "debug", "Level")
	})

	t.Run("parse LogLevel flag invalid returns error", func(t *testing.T) {
		t.Parallel()

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

		setFlag(t, cmd, "level", "invalid")

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "level")
	})

	t.Run("parse invalid Enum returns error", func(t *testing.T) {
		t.Parallel()

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

		setFlag(t, cmd, "mode", "invalid")

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("parse LogFormat flag valid", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Format LogFormat `flag:"format"`
		}

		cfg := &TestConfig{}

		registry, err := NewFlagRegistry(*cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}
		registerAndSetFlag(t, registry, cmd, cfg, "format", "json")

		assertEnumString(t, cfg.Format.String(), "json", "Format")
	})

	t.Run("parse LogFormat flag invalid returns error", func(t *testing.T) {
		t.Parallel()

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

		setFlag(t, cmd, "format", "invalid")

		err = registry.ParseFlags(cmd, cfg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "format")
	})
}

func TestFlagRegistry_FlagNotFound(t *testing.T) {
	t.Parallel()
	t.Run("missing flag returns error", func(t *testing.T) {
		t.Parallel()

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
