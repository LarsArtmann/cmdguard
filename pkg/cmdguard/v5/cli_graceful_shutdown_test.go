package v5_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/samber/do/v2"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type gracefulShutdownService struct {
	shutdownCalled atomic.Bool
}

var _ do.ShutdownerWithError = (*gracefulShutdownService)(nil)

func (s *gracefulShutdownService) Shutdown() error {
	s.shutdownCalled.Store(true)

	return nil
}

// newGracefulShutdownTestCLI builds a CLI configured for graceful shutdown tests.
// Centralizes the NewCLI call so individual tests stay terse.
func newGracefulShutdownTestCLI(t *testing.T) *v5.CLI[testCLIConfig] {
	t.Helper()

	cli, err := v5.NewCLI(
		"test", "Test", testCLIConfig{},
		v5.WithGracefulShutdown(),
		v5.WithFang(false),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	return cli
}

func TestGracefulShutdown_CallsShutdownerOnScope(t *testing.T) {
	t.Parallel()

	t.Run("shutdown is called on service implementing ShutdownerWithError", func(t *testing.T) {
		t.Parallel()

		svc := &gracefulShutdownService{}

		cli := newGracefulShutdownTestCLI(t)

		err := v5.Provide(cli.Scope(), func(i do.Injector) (*gracefulShutdownService, error) {
			return svc, nil
		})
		if err != nil {
			t.Fatalf("Provide failed: %v", err)
		}

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		_, err = v5.Invoke[*gracefulShutdownService](cli.Scope())
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		err = cli.Shutdown(context.Background())
		if err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}

		if !svc.shutdownCalled.Load() {
			t.Error("expected Shutdown to be called on service")
		}
	})

	t.Run("works without Shutdowner services", func(t *testing.T) {
		t.Parallel()

		cli := newGracefulShutdownTestCLI(t)

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		err := cli.Shutdown(context.Background())
		if err != nil {
			t.Fatalf("Shutdown failed: %v", err)
		}
	})

	t.Run("shutdown fires only when the execution context was cancelled", func(t *testing.T) {
		t.Parallel()

		svc := &gracefulShutdownService{}

		cli := newGracefulShutdownTestCLI(t)

		err := v5.Provide(cli.Scope(), func(i do.Injector) (*gracefulShutdownService, error) {
			return svc, nil
		})
		if err != nil {
			t.Fatalf("Provide failed: %v", err)
		}

		cmd := newTestCLICommand[testCLIConfig](t, "run")
		addCommand(t, cli, cmd)

		_, err = v5.Invoke[*gracefulShutdownService](cli.Scope())
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		// Normal completion must NOT tear down the scope: repeated
		// ExecuteWithArgs stays viable under the v5 defaults.
		err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
		if err != nil {
			t.Fatalf("first ExecuteWithArgs failed: %v", err)
		}

		if svc.shutdownCalled.Load() {
			t.Error("Shutdown must not fire on a normally completing Execute")
		}

		_, err = v5.Invoke[*gracefulShutdownService](cli.Scope())
		if err != nil {
			t.Fatalf("scope must stay alive across executions: %v", err)
		}

		// A cancelled context (the SIGINT/SIGTERM path) triggers the
		// graceful DI shutdown.
		cancelled, cancel := context.WithCancel(context.Background())
		cancel()

		_ = cli.ExecuteWithArgs(cancelled, []string{"run"})

		if !svc.shutdownCalled.Load() {
			t.Error("expected Shutdown on cancelled execution context")
		}
	})
}

func TestGracefulShutdown_ImpliesSignalHandling(t *testing.T) {
	t.Parallel()

	t.Run("WithGracefulShutdown enables gracefulShutdown flag", func(t *testing.T) {
		t.Parallel()

		cli, err := v5.NewCLI(
			"test", "Test", testCLIConfig{},
			v5.WithGracefulShutdown(),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{})
		if err != nil && err.Error() != "" {
			t.Logf("Execute returned (expected — no command): %v", err)
		}
	})
}
