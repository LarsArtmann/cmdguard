package v3_test

import (
	"sync"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestWithDILogging(t *testing.T) {
	t.Parallel()

	t.Run("captures DI log output", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		var logs []string

		cli, err := v3.NewCLI(
			"test", "Test", testCLIConfig{},
			v3.WithDILogging(func(format string, args ...any) {
				mu.Lock()
				logs = append(logs, format)
				mu.Unlock()
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{})
		if err == nil {
			t.Log("Execute completed")
		}

		mu.Lock()
		captured := len(logs)
		mu.Unlock()

		if captured == 0 {
			t.Error("expected DI log output to be captured")
		}
	})

	t.Run("no logging when option not set", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{})
		if err == nil {
			t.Log("Execute completed without logging")
		}
	})
}

func TestNewScopeWithOpts(t *testing.T) {
	t.Parallel()

	t.Run("creates scope with custom opts", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScopeWithOpts("test", nil)
		if scope == nil {
			t.Fatal("NewScopeWithOpts returned nil")
		}

		if scope.Name() != "test" {
			t.Errorf("Name = %q, want %q", scope.Name(), "test")
		}
	})
}
