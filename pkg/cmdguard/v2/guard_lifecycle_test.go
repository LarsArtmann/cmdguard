package v2

import (
	"testing"
)

func TestCLI_Shutdown(t *testing.T) {
	t.Parallel()
	t.Run("shutdown succeeds", func(t *testing.T) {
		t.Parallel()
		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = cli.Shutdown(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCLI_HealthCheck(t *testing.T) {
	t.Parallel()
	t.Run("health check succeeds", func(t *testing.T) {
		t.Parallel()
		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = cli.HealthCheck()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
