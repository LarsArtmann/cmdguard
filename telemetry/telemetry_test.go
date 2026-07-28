package telemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace/noop"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type testConfig struct {
	Name string `flag:"name" short:"n" default:"" help:"Your name"`
}

func TestMiddleware_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	mw := Middleware[testConfig](noop.NewTracerProvider().Tracer("test"))
	if mw == nil {
		t.Fatal("Middleware returned nil")
	}
}

func TestMiddleware_NilTracerCallsNext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := &testConfig{}
	info := v4.CommandInfo{Name: "test", Phase: v4.PhaseRun}

	called := false
	mw := Middleware[testConfig](nil)

	err := mw(ctx, cfg, info, func() error {
		called = true

		return nil
	})
	if err != nil {
		t.Fatalf("middleware returned error: %v", err)
	}

	if !called {
		t.Fatal("next function was not called with nil tracer")
	}
}

func TestMiddleware_NoopTracerCallsNext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := &testConfig{}
	info := v4.CommandInfo{Name: "test", Phase: v4.PhaseRun}

	called := false
	tracer := noop.NewTracerProvider().Tracer("test")
	mw := Middleware[testConfig](tracer)

	err := mw(ctx, cfg, info, func() error {
		called = true

		return nil
	})
	if err != nil {
		t.Fatalf("middleware returned error: %v", err)
	}

	if !called {
		t.Fatal("next function was not called with noop tracer")
	}
}

func TestWithTelemetry_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	tracer := noop.NewTracerProvider().Tracer("test")
	opt := WithTelemetry[testConfig](tracer)
	if opt == nil {
		t.Fatal("WithTelemetry returned nil")
	}
}
