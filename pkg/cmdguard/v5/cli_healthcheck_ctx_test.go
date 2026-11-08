package v5_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/samber/do/v2"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

// Compile-time proof that CLI satisfies samber/do's context-aware healthcheck
// interface (guards against silent signature drift; W4).
type healthConformanceConfig struct{}

var _ do.HealthcheckerWithContext = (*v5.CLI[healthConformanceConfig])(nil)

// ctxAwareService implements do.HealthcheckerWithContext and records the ctx
// state it was called with.
type ctxAwareService struct {
	sawCanceled atomic.Bool
	calls       atomic.Int64
}

var _ do.HealthcheckerWithContext = (*ctxAwareService)(nil)

func (s *ctxAwareService) HealthCheck(ctx context.Context) error {
	s.calls.Add(1)

	if ctx.Err() != nil {
		s.sawCanceled.Store(true)

		return ctx.Err()
	}

	return nil
}

func TestHealthCheck_HonorsContextCancellation(t *testing.T) {
	t.Parallel()

	svc := &ctxAwareService{}

	cli, err := v5.NewCLI("app", "App", healthConformanceConfig{}, v5.WithFang(false))
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	if err := v5.Provide(cli.Scope(), func(_ do.Injector) (*ctxAwareService, error) {
		return svc, nil
	}); err != nil {
		t.Fatalf("Provide: %v", err)
	}

	t.Run("live context passes", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := cli.HealthCheck(ctx); err != nil {
			t.Fatalf("HealthCheck with live ctx: %v", err)
		}

		if svc.sawCanceled.Load() {
			t.Error("service saw a cancelled context unexpectedly")
		}
	})

	t.Run("cancelled context surfaces the error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := cli.HealthCheck(ctx)
		if err == nil {
			t.Fatal("expected error from cancelled ctx health check")
		}

		if !errors.Is(err, context.Canceled) && svc.calls.Load() > 0 && !svc.sawCanceled.Load() {
			t.Errorf("expected cancellation to reach the service or the result, got: %v", err)
		}
	})

	t.Run("no-arg HealthCheck uses background context", func(t *testing.T) {
		t.Parallel()

		scope := cli.Scope()

		if err := scope.HealthCheck(); err != nil {
			t.Fatalf("Scope.HealthCheck background: %v", err)
		}

		results := scope.HealthCheckResults()
		if results == nil {
			t.Fatal("HealthCheckResults returned nil map")
		}

		if err, exists := results["*v5_test.ctxAwareService"]; exists && err != nil {
			t.Errorf("background health check must not report cancellation: %v", err)
		}
	})
}
