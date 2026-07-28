package v4_test

import (
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestCLISetConfig(t *testing.T) {
	t.Parallel()
	t.Run("updates config", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cli.SetConfig(testCLIConfig{Verbose: true, Level: "debug"})

		cfg := cli.Config()
		if cfg == nil {
			t.Fatal("Config() returned nil")
		}

		if !cfg.Verbose {
			t.Error("Verbose not updated")
		}

		if cfg.Level != "debug" {
			t.Errorf("Level = %q, want %q", cfg.Level, "debug")
		}
	})
}

func TestCLIShutdown(t *testing.T) {
	t.Parallel()
	t.Run("shutdown succeeds", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.Shutdown(t.Context())
		if err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}
	})
}

func TestCLIHealthCheck(t *testing.T) {
	t.Parallel()
	t.Run("health check succeeds", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheck()
		if err != nil {
			t.Errorf("HealthCheck failed: %v", err)
		}
	})
}

func TestCLIHealthCheckWithContext(t *testing.T) {
	t.Parallel()
	t.Run("health check with context succeeds", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheckWithContext(t.Context())
		if err != nil {
			t.Errorf("HealthCheckWithContext failed: %v", err)
		}
	})
}
