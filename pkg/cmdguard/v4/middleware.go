package v4

import (
	"context"
	"errors"
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

// ContextMiddleware is a context-aware variant of [Middleware]. The next
// function accepts a context, so middleware that derives a new context
// (timeouts, cancellation, tracing scopes) can propagate it to every inner
// middleware and to the final handler — something Middleware cannot do,
// because its next closure is bound to the invocation context.
//
// Context middleware run OUTSIDE plain middleware: they wrap the whole chain,
// and the context they pass to next reaches every plain middleware and the
// command handler.
//
// Example — timeout middleware:
//
//	TimeoutMiddleware := func(ctx context.Context, cfg *MyConfig, info CommandInfo, next func(context.Context) error) error {
//	    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//	    defer cancel()
//	    return next(ctx)
//	}
type ContextMiddleware[T any] func(ctx context.Context, cfg *T, info CommandInfo, next func(context.Context) error) error

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
// TODO(v5): "Info" suffix is vague; consider CommandMetadata (see naming-review 2026-07-18).
type CommandInfo struct {
	Name     string
	FullPath string
	Phase    Phase
	HasRunE  bool
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

// buildContextChain is the [ContextMiddleware] counterpart of buildChain:
// the context each middleware passes to next is threaded through the whole
// chain down to the final handler.
func buildContextChain[T any](
	ctx context.Context,
	cfg *T,
	info CommandInfo,
	middlewares []ContextMiddleware[T],
	final func(context.Context) error,
) func(context.Context) error {
	for _, v := range slices.Backward(middlewares) {
		mw := v
		prev := final
		final = func(c context.Context) error {
			return mw(c, cfg, info, prev)
		}
	}

	return final
}

// TimingMiddleware returns a middleware that logs command execution duration.
// The log function receives the command name, duration, and any error from execution.
func TimingMiddleware[T any](log func(commandName string, d time.Duration, err error)) Middleware[T] {
	return func(_ context.Context, _ *T, info CommandInfo, next func() error) error {
		start := time.Now()
		err := next()

		log(info.Name, time.Since(start), err)

		return err
	}
}

// TimeoutMiddleware returns a context middleware that bounds command
// execution to d: handlers receive the timeout-derived context and are
// expected to honor its cancellation. When the handler returns an error
// matching context.DeadlineExceeded, the returned error matches both
// [ErrCommandTimeout] and context.DeadlineExceeded via errors.Is.
//
// A non-positive d disables the limit (next is called unchanged). Note that a
// handler which ignores its context cannot be interrupted — Go cannot kill a
// goroutine — so the middleware guarantees the derived context and the error
// mapping, not forcible cancellation.
func TimeoutMiddleware[T any](d time.Duration) ContextMiddleware[T] {
	return func(ctx context.Context, _ *T, info CommandInfo, next func(context.Context) error) error {
		if d <= 0 {
			return next(ctx)
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		err := next(timeoutCtx)
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf(
				"%w: command %q exceeded %s: %w",
				ErrCommandTimeout,
				info.Name,
				d,
				context.DeadlineExceeded,
			)
		}

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
