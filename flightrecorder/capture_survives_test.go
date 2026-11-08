package flightrecorder

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

// R5.4: pins the WithoutCancel contract (2026-11-08 fix). Under v5's default
// signal handling, the run context is cancelled when Execute returns; the
// error-triggered snapshot is post-hoc telemetry and must still land.
// Before the fix, the background capture raced Execute's return and wrote
// nothing (context cancelled mid-write).

func TestMiddleware_ErrorCaptureSurvivesExecuteReturn(t *testing.T) { //nolint:paralleltest // process-wide singleton recorder
	dir := t.TempDir()

	cfg := Config{
		CaptureOnError: true,
		OutputDir:      dir,
	}

	cli, err := v5.NewCLI("snap", "Snapshot app", struct{}{},
		v5.WithoutSignalHandling(), // deterministic: no real signals; we cancel manually
		WithFlightRecorder[v5.CLI[struct{}]](cfg),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "boom", v5.NoFlags{}, func(ctx context.Context, cfgPtr *struct{}, flags v5.NoFlags) error {
		return errBoom{}
	}, v5.WithShort("Boom"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Cancelled parent context simulates the v5 signal path: NotifyContext
	// wraps Execute and cancels at return, so the middleware sees a ctx that
	// is (or is about to be) cancelled while the capture goroutine runs.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_ = cli.ExecuteWithArgs(ctx, []string{"boom"})

	// The capture goroutine is asynchronous: poll briefly for the snapshot.
	deadline := time.Now().Add(3 * time.Second)

	for time.Now().Before(deadline) {
		entries, readErr := os.ReadDir(dir)
		if readErr == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".trace") {
					return // snapshot landed despite the cancelled context
				}
			}
		}

		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("no .trace snapshot captured after an erroring run on a cancelled context — post-hoc capture did not survive Execute return")
}

type errBoom struct{}

func (errBoom) Error() string { return "boom" }

