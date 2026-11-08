package flightrecorder

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

// TestMiddleware_ErrorCaptureSurvivesExecuteReturn pins the WithoutCancel
// contract (2026-11-08 fix). Under v5's default signal handling, the run
// context is cancelled as Execute returns; the error-triggered snapshot is
// post-hoc telemetry and must still land. Before the fix, the background
// capture raced Execute's return and wrote nothing (context cancelled
// mid-write).
func TestMiddleware_ErrorCaptureSurvivesExecuteReturn(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		CaptureOnError: true,
		OutputDir:      dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	defer rec.Stop() // waits for in-flight captures before tearing down

	cli, err := v5.NewCLI("snap", "Snapshot app", struct{}{},
		v5.WithoutSignalHandling(), // deterministic: no real signals; we cancel manually
		WithFlightRecorderRecorder[struct{}](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "boom", v5.NoFlags{}, func(ctx context.Context, cfgPtr *struct{}, flags v5.NoFlags) error {
		return boomError{}
	}, v5.WithShort("Boom"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	// The v5 signal path: the run context is LIVE while the command executes
	// and cancelled as Execute returns (NotifyContext teardown). That is
	// exactly the race the WithoutCancel fix addresses: the capture goroutine
	// starts during the run and must survive Execute's teardown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = cli.ExecuteWithArgs(ctx, []string{"boom"})
	cancel() // teardown fires before the async capture may have finished

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

	t.Fatal("no .trace snapshot after erroring run on cancelled context: capture did not survive Execute return")
}

type boomError struct{}

func (boomError) Error() string { return "boom" }
