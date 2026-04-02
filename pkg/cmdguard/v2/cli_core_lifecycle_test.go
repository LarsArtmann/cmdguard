package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestCLISetConfig(t *testing.T) {
	t.Run("updates config", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
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
	t.Run("shutdown succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
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
	t.Run("health check succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
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
	t.Run("health check with context succeeds", func(t *testing.T) {
		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheckWithContext(t.Context())
		if err != nil {
			t.Errorf("HealthCheckWithContext failed: %v", err)
		}
	})
}
