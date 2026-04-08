package v2

import (
	"testing"
)

func TestGuardedCommand_Shutdown(t *testing.T) {
	t.Parallel()
	t.Run("shutdown succeeds", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = g.Shutdown(t.Context())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGuardedCommand_HealthCheck(t *testing.T) {
	t.Parallel()
	t.Run("health check succeeds", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = g.HealthCheck()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
