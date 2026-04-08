package v2

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestFlagRegistry_ParseFlags(t *testing.T) {
	t.Parallel()
	t.Run("parse string flag", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
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

	t.Run("parse uint64 flag", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			MaxBytes uint64 `default:"0" flag:"max-bytes"`
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

		if err := cmd.PersistentFlags().Set("max-bytes", "18446744073709551615"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.MaxBytes != 18446744073709551615 {
			t.Errorf("expected MaxBytes 18446744073709551615, got %d", cfg.MaxBytes)
		}
	})

	t.Run("parse float64 flag", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			Alpha float64 `default:"1.0" flag:"alpha"`
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

		if err := cmd.PersistentFlags().Set("alpha", "0.5"); err != nil {
			t.Fatalf("expected no error setting flag, got: %v", err)
		}

		err = registry.ParseFlags(cmd, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Alpha != 0.5 {
			t.Errorf("expected Alpha %f, got %f", 0.5, cfg.Alpha)
		}
	})

	t.Run("parse Duration flag", func(t *testing.T) {
		t.Parallel()
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
}
