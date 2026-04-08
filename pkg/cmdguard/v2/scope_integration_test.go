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
		if err := ProvideValue(root, Config{Debug: true}); err != nil {
			t.Fatalf("expected no error providing value, got: %v", err)
		}

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
		assertChildInheritsParent(t)
	})

	t.Run("Package function creates CLI with DI", func(t *testing.T) {
		t.Parallel()
		type config struct {
			Name string
		}

		pkg := Package("test-app", "Test Application", config{Name: "test"})

		// Package should return a valid function
		if pkg == nil {
			t.Fatal("Package returned nil function")
		}

		// Should not panic when called
		assertNotPanic(t, func() {
			pkg(nil)
		})
	})

	t.Run("Package function with options", func(t *testing.T) {
		t.Parallel()
		type config struct {
			Version string
		}

		pkg := Package(
			"test-app",
			"Test Application",
			config{Version: "1.0.0"},
			WithCLIVersion[config]("1.0.0"),
		)

		assertNotPanic(t, func() {
			pkg(nil)
		})
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
