package flightrecorder

import (
	"context"
	"time"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// Middleware returns a cmdguard middleware that manages the flight recorder
// lifecycle and automatically captures trace snapshots when configured
// conditions are met.
//
// The recorder is lazily started on the first command invocation. Snapshots
// are only captured during the run phase to avoid duplicate captures from
// pre-run/post-run phases.
//
// The middleware never blocks command execution on recorder operations: if
// Start fails (e.g. another recorder is already active), the command proceeds
// without tracing. Snapshot capture runs in a background goroutine so the
// command response is not delayed.
//
// The returned error from next is always passed through unchanged — the
// middleware only adds observability, it never changes program behavior.
func Middleware[T any](rec *Recorder) v4.Middleware[T] {
	return func(ctx context.Context, _ *T, info v4.CommandInfo, next func() error) error {
		if !rec.Enabled() {
			if startErr := rec.Start(); startErr != nil {
				rec.logf("flightrecorder: failed to start: %v", startErr)

				return next()
			}
		}

		// Only evaluate capture conditions during the run phase.
		// Pre-run and post-run phases should not trigger captures.
		if info.Phase != v4.PhaseRun {
			return next()
		}

		start := time.Now()
		runErr := next()
		elapsed := time.Since(start)

		evaluateCapture(ctx, rec, info, elapsed, runErr)

		return runErr
	}
}

// evaluateCapture checks capture conditions and triggers a snapshot if warranted.
// When both CaptureOnSlow and CaptureOnError are enabled and a command is both
// slow AND errors, error takes precedence over slow — the error reason is always
// the primary concern for debugging.
// The snapshot write runs in a background goroutine to avoid blocking.
func evaluateCapture(ctx context.Context, rec *Recorder, info v4.CommandInfo, elapsed time.Duration, runErr error) {
	cfg := rec.Config()

	var reason CaptureReason

	shouldCapture := false

	if cfg.CaptureOnSlow && elapsed > cfg.SlowThreshold {
		reason = CaptureReasonSlow
		shouldCapture = true
	}

	// Error takes precedence over slow when both conditions are met:
	// the error is always the primary concern for debugging.
	if cfg.CaptureOnError && runErr != nil {
		reason = CaptureReasonError
		shouldCapture = true
	}

	if !shouldCapture || !rec.Enabled() {
		return
	}

	name := info.Name
	if info.FullPath != "" {
		name = info.FullPath
	}

	finalName := name
	finalReason := reason

	go func() {
		if _, captureErr := rec.Capture(ctx, finalName, finalReason); captureErr != nil {
			rec.logf("flightrecorder: failed to capture %s snapshot: %v", finalReason, captureErr)
		}
	}()
}

// WithFlightRecorder is a convenience CLIOption that creates a Recorder from
// the given Config and registers it as middleware via v4.WithMiddleware.
//
// The recorder is lazily started on first command execution. For explicit
// lifecycle control (Start/Stop), use New + Middleware instead.
//
// Example:
//
//	cli, _ := v4.NewCLI[Config]("app", "My app", Config{},
//	    flightrecorder.WithFlightRecorder[Config](flightrecorder.Config{
//	        CaptureOnSlow:  true,
//	        SlowThreshold:  200 * time.Millisecond,
//	        CaptureOnError: true,
//	    }),
//	)
func WithFlightRecorder[T any](cfg Config) v4.CLIOption {
	rec := New(cfg)

	return v4.WithMiddleware(Middleware[T](rec))
}

// WithFlightRecorderRecorder is a CLIOption that registers a pre-created Recorder
// as middleware. Unlike WithFlightRecorder (which creates a new Recorder from Config),
// this variant gives the caller explicit lifecycle control — useful when multiple
// CLIs share a single recorder or when Start/Stop must be called at specific times.
//
// Example:
//
//	rec := flightrecorder.New(flightrecorder.DefaultConfig())
//	defer rec.Stop()
//
//	cli, _ := v4.NewCLI[Config]("app", "My app", Config{},
//	    flightrecorder.WithFlightRecorderRecorder[Config](rec),
//	)
func WithFlightRecorderRecorder[T any](rec *Recorder) v4.CLIOption {
	return v4.WithMiddleware(Middleware[T](rec))
}
