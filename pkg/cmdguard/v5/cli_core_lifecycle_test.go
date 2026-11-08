package v5_test

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestCLIShutdown(t *testing.T) {
	t.Parallel()
	t.Run("shutdown succeeds", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI("test", "Test", testCLIConfig{})
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

		cli, err := v5.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheck(t.Context())
		if err != nil {
			t.Errorf("HealthCheck failed: %v", err)
		}
	})
}

func TestCLIHealthCheckWithContext(t *testing.T) {
	t.Parallel()
	t.Run("health check with context succeeds", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.HealthCheckWithContext(t.Context())
		if err != nil {
			t.Errorf("HealthCheckWithContext failed: %v", err)
		}
	})
}
