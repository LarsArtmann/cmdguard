package v2

import (
	"context"
	"testing"
	"time"

	"github.com/samber/do/v2"
)

func TestScope_Integration(t *testing.T) {
	t.Parallel()
	t.Run("full workflow with DI", func(t *testing.T) {
		t.Parallel()
		// Create root scope
		root := NewScope("app")

		// Register services
		type Config struct {
			Debug bool
		}
		mustProvideValue(t, root, Config{Debug: true})

		if err := Provide(root, func(i do.Injector) (string, error) {
			cfg, err := do.Invoke[Config](i)
			if err != nil {
				return "", err
			}

			if cfg.Debug {
				return "debug-mode", nil
			}

			return "production-mode", nil
		}); err != nil {
			t.Fatalf("expected no error providing service, got: %v", err)
		}

		// Verify services
		cfg, err := Invoke[Config](root)
		if err != nil {
			t.Fatalf("expected no error invoking config, got: %v", err)
		}

		if !cfg.Debug {
			t.Error("expected Debug to be true")
		}

		mode, err := Invoke[string](root)
		if err != nil {
			t.Fatalf("expected no error invoking mode, got: %v", err)
		}

		if mode != "debug-mode" {
			t.Errorf("expected mode to be 'debug-mode', got %q", mode)
		}

		// Health check
		if err := root.HealthCheck(); err != nil {
			t.Errorf("expected no error from health check, got: %v", err)
		}

		// Shutdown
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		if err := root.Shutdown(ctx); err != nil {
			t.Errorf("expected no error from shutdown, got: %v", err)
		}
	})

	t.Run("child scope can override parent services", func(t *testing.T) {
		t.Parallel()

		parent := NewScope("parent")
		mustProvideValue(t, parent, "parent-value")

		child := parent.Child("child")

		value, err := Invoke[string](child)
		if err != nil {
			t.Fatalf("expected no error invoking, got: %v", err)
		}

		if value != "parent-value" {
			t.Errorf("expected value 'parent-value', got %q", value)
		}
	})

	t.Run("Package function creates CLI with DI", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Name string
		}

		cli, err := Package("test-app", "Test Application", config{Name: "test"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cli == nil {
			t.Fatal("Package returned nil CLI")
		}
	})

	t.Run("Package function with options", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Version string
		}

		cli, err := Package(
			"test-app",
			"Test Application",
			config{Version: "1.0.0"},
			WithCLIVersion[config]("1.0.0"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cli == nil {
			t.Fatal("Package returned nil CLI")
		}
	})
}

func assertNotPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	fn()
}
