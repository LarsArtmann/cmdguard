package spinner

import (
	"bytes"
	"context"
	"testing"
	"time"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type testConfig struct {
	Name string `flag:"name" short:"n" default:"" help:"Your name"`
}

func TestMiddleware_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	middleware := Middleware[testConfig]("Loading...")
	if middleware == nil {
		t.Fatal("Middleware returned nil")
	}
}

func TestMiddlewareWithConfig_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	middleware := MiddlewareWithConfig[testConfig](Config{
		Title:    "Working",
		Frames:   []string{"|", "/", "-", "\\"},
		Interval: 50 * time.Millisecond,
		Writer:   &bytes.Buffer{},
	})
	if middleware == nil {
		t.Fatal("MiddlewareWithConfig returned nil")
	}
}

func TestMiddleware_NilTracerSkips(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := &testConfig{}
	info := v3.CommandInfo{Name: "test", Phase: v3.PhaseRun}

	called := false
	middleware := MiddlewareWithConfig[testConfig](Config{
		Title:    "Loading",
		Frames:   []string{"-"},
		Interval: 10 * time.Millisecond,
		Writer:   &bytes.Buffer{},
	})

	err := middleware(ctx, cfg, info, func() error {
		called = true

		return nil
	})
	if err != nil {
		t.Fatalf("middleware returned error: %v", err)
	}

	if !called {
		t.Fatal("next function was not called")
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	t.Parallel()

	c := Config{Title: "test"}
	if c.Title != "test" {
		t.Fatalf("expected Title='test', got %q", c.Title)
	}
}
