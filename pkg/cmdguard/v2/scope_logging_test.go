package v2_test

import (
	"sync"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func TestWithDILogging(t *testing.T) {
	t.Parallel()

	t.Run("captures DI log output", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		var logs []string

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test", testCLIConfig{},
			v2.WithDILogging[testCLIConfig](func(format string, args ...any) {
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

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
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

		scope := v2.NewScopeWithOpts("test", nil)
		if scope == nil {
			t.Fatal("NewScopeWithOpts returned nil")
		}

		if scope.Name() != "test" {
			t.Errorf("Name = %q, want %q", scope.Name(), "test")
		}
	})
}
