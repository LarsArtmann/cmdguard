package v2

import (
	"context"
	"fmt"
	"runtime/debug"
	"slices"
	"time"
)

// Middleware wraps command execution for cross-cutting concerns like logging,
// metrics, tracing, and error recovery.
//
// The next function calls the next middleware in the chain (or the final handler).
// Middleware must call next() exactly once. To short-circuit, return before calling next.
//
// Example — timing middleware:
//
//	TimingMiddleware := func(ctx context.Context, cfg *MyConfig, info CommandInfo, next func() error) error {
//	    start := time.Now()
//	    err := next()
//	    fmt.Fprintf(os.Stderr, "%s took %v\n", info.Name, time.Since(start))
//	    return err
//	}
type Middleware[T any] func(ctx context.Context, cfg *T, info CommandInfo, next func() error) error

// Phase identifies the execution stage of a command handler.
type Phase string

const (
	// PhaseRun is the main command execution phase.
	PhaseRun Phase = "run"
	// PhasePreRun is the pre-validation hook phase.
	PhasePreRun Phase = "pre-run"
	// PhasePostRun is the post-success cleanup phase.
	PhasePostRun Phase = "post-run"
)

// CommandInfo provides command metadata to middleware.
type CommandInfo struct {
	Name    string
	Phase   Phase
	HasRunE bool
}

// buildChain builds a single handler from a slice of middleware.
// Middleware are applied in order: first middleware wraps the second, etc.
func buildChain[T any](
	ctx context.Context,
	cfg *T,
	info CommandInfo,
	middlewares []Middleware[T],
	final func() error,
) func() error {
	for _, v := range slices.Backward(middlewares) {
		mw := v
		prev := final
		final = func() error {
			return mw(ctx, cfg, info, prev)
		}
	}

	return final
}

// TimingMiddleware returns a middleware that logs command execution duration.
// The log function receives the command name and duration.
func TimingMiddleware[T any](log func(commandName string, d time.Duration)) Middleware[T] {
	return func(_ context.Context, _ *T, info CommandInfo, next func() error) error {
		start := time.Now()
		err := next()

		log(info.Name, time.Since(start))

		return err
	}
}

// RecoveryMiddleware returns a middleware that recovers from panics in command handlers,
// converting them to errors with stack traces.
func RecoveryMiddleware[T any]() Middleware[T] {
	return func(_ context.Context, _ *T, info CommandInfo, next func() error) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf(
					"%w: panic in command %q: %v\n%s",
					ErrCommandPanic,
					info.Name,
					r,
					debug.Stack(),
				)
			}
		}()

		return next()
	}
}
