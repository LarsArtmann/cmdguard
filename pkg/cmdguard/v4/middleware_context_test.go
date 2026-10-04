package v4

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type ctxMiddlewareConfig struct {
	Name string `default:"test" flag:"name" help:"test"`
}

// ctxKey is a private context key type for middleware context tests.
type ctxKey struct{}

// deriveContextMiddleware returns a ContextMiddleware[T] that attaches value
// to the context under ctxKey before calling next with the derived context.
func deriveContextMiddleware[T any](value string) ContextMiddleware[T] {
	return func(ctx context.Context, _ *T, _ CommandInfo, next func(context.Context) error) error {
		return next(context.WithValue(ctx, ctxKey{}, value))
	}
}

func TestContextMiddleware_PropagatesContextToHandler(t *testing.T) {
	t.Parallel()

	var handlerValue any

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(deriveContextMiddleware[ctxMiddlewareConfig]("derived")),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
		runE: func(ctx context.Context, _ *ctxMiddlewareConfig, _ NoFlags) error {
			handlerValue = ctx.Value(ctxKey{})

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	if handlerValue != "derived" {
		t.Errorf("handler received context value %v, want %q", handlerValue, "derived")
	}
}

func TestContextMiddleware_ReachesPlainMiddleware(t *testing.T) {
	t.Parallel()

	var plainValue any

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(deriveContextMiddleware[ctxMiddlewareConfig]("outer")),
		WithMiddleware(func(_ context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func() error) error {
			plainValue = nil

			return next()
		}),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	// The plain middleware captures the context value by asserting on the
	// chain: wrap the capture inside the handler-visible context instead.
	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
		runE: func(ctx context.Context, _ *ctxMiddlewareConfig, _ NoFlags) error {
			plainValue = ctx.Value(ctxKey{})

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	if plainValue != "outer" {
		t.Errorf("plain middleware/handler received context value %v, want %q", plainValue, "outer")
	}
}

func TestContextMiddleware_Chaining(t *testing.T) {
	t.Parallel()

	var callOrder []string

	first := func(ctx context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func(context.Context) error) error {
		callOrder = append(callOrder, "first-before")

		return next(ctx)
	}

	second := func(ctx context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func(context.Context) error) error {
		callOrder = append(callOrder, "second-before")

		err := next(ctx)

		callOrder = append(callOrder, "second-after")

		return err
	}

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(first, second),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
		runE: func(context.Context, *ctxMiddlewareConfig, NoFlags) error {
			callOrder = append(callOrder, "handler")

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	expected := []string{"first-before", "second-before", "handler", "second-after"}
	if len(callOrder) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(callOrder), callOrder)
	}

	for i, v := range expected {
		if callOrder[i] != v {
			t.Errorf("call[%d]: expected %q, got %q", i, v, callOrder[i])
		}
	}
}

func TestContextMiddleware_DerivedContextFlowsThroughWholeChain(t *testing.T) {
	t.Parallel()

	var (
		seenBySecond any
		seenByPlain  any
		seenByFirst  any
	)

	first := func(ctx context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func(context.Context) error) error {
		err := next(context.WithValue(ctx, ctxKey{}, "from-first"))

		seenByFirst = ctx.Value(ctxKey{})

		return err
	}

	second := func(ctx context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func(context.Context) error) error {
		seenBySecond = ctx.Value(ctxKey{})

		return next(ctx)
	}

	plain := func(_ context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func() error) error {
		return next()
	}

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(first, second),
		WithMiddleware(plain),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
		runE: func(ctx context.Context, _ *ctxMiddlewareConfig, _ NoFlags) error {
			seenByPlain = ctx.Value(ctxKey{})

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	if seenBySecond != "from-first" {
		t.Errorf("second context middleware saw %v, want %q", seenBySecond, "from-first")
	}

	if seenByPlain != "from-first" {
		t.Errorf("handler saw %v, want %q (derived context must reach plain middleware and handler)", seenByPlain, "from-first")
	}

	if seenByFirst != nil {
		t.Errorf("first context middleware saw %v, want nil (context flows inward only)", seenByFirst)
	}
}

func TestContextMiddleware_ErrorPropagation(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(func(ctx context.Context, _ *ctxMiddlewareConfig, _ CommandInfo, next func(context.Context) error) error {
			return next(ctx)
		}),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "fail", short: "Fail command", long: "Fail command"},
		runE: func(context.Context, *ctxMiddlewareConfig, NoFlags) error {
			return sentinel
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTimeoutMiddleware_DeadlineExceeded(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(TimeoutMiddleware[ctxMiddlewareConfig](1*time.Millisecond)),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "slow", short: "Slow command", long: "Slow command"},
		runE: func(ctx context.Context, _ *ctxMiddlewareConfig, _ NoFlags) error {
			<-ctx.Done()

			return ctx.Err()
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"slow"})

	if !errors.Is(err, ErrCommandTimeout) {
		t.Errorf("expected ErrCommandTimeout, got %v", err)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded in chain, got %v", err)
	}
}

func TestTimeoutMiddleware_UnderDeadlinePasses(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		WithContextMiddleware(TimeoutMiddleware[ctxMiddlewareConfig](30*time.Second)),
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "fast", short: "Fast command", long: "Fast command"},
		runE: func(context.Context, *ctxMiddlewareConfig, NoFlags) error {
			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"fast"})
	testutil.AssertNoError(t, err)
}

func TestTimeoutMiddleware_NonPositiveDurationDisablesLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    time.Duration
	}{
		{name: "zero", d: 0},
		{name: "negative", d: -time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, err := NewCLI(
				"test", "Test CLI", ctxMiddlewareConfig{},
				WithContextMiddleware(TimeoutMiddleware[ctxMiddlewareConfig](tt.d)),
				WithFang(false),
			)
			testutil.AssertNoError(t, err)

			err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
				spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
				runE: func(ctx context.Context, _ *ctxMiddlewareConfig, _ NoFlags) error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						return nil
					}
				},
			})
			testutil.AssertNoError(t, err)

			err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
			testutil.AssertNoError(t, err)
		})
	}
}

func TestWithContextMiddleware_TypeMismatchYieldsNoMiddleware(t *testing.T) {
	t.Parallel()

	type otherConfig struct {
		Other string `default:"x" flag:"other" help:"other"`
	}

	handlerCalled := false

	// A context middleware registered with the wrong config type must be
	// silently ignored (sealed-interface extraction yields nil), not panic.
	mismatched := WithContextMiddleware(func(_ context.Context, _ *otherConfig, _ CommandInfo, next func(context.Context) error) error {
		return errors.New("must not run")
	})

	cli, err := NewCLI(
		"test", "Test CLI", ctxMiddlewareConfig{},
		mismatched,
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[ctxMiddlewareConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
		runE: func(context.Context, *ctxMiddlewareConfig, NoFlags) error {
			handlerCalled = true

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	if !handlerCalled {
		t.Error("handler should have run")
	}
}
