package flightrecorder_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/flightrecorder"
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type integrationConfig struct {
	Debug bool `flag:"debug" default:"false" help:"Debug mode"`
}

// TestIntegration_CaptureOnCommandError is the critical integration test: it
// wires WithFlightRecorderRecorder through a real CLI[T].ExecuteWithArgs flow,
// runs a command that returns an error, and asserts that a .trace snapshot file
// is produced in the output directory.
//
// This validates the entire chain: CLIOption → WithMiddleware → middleware
// chain execution → cobra context propagation → lazy Start → evaluateCapture
// → async capture goroutine → file write.
func TestIntegration_CaptureOnCommandError(t *testing.T) {
	dir := t.TempDir()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnError: true,
	})

	defer rec.Stop()

	cli, err := v5.NewCLI[integrationConfig](
		"testapp", "Test app", integrationConfig{},
		flightrecorder.WithFlightRecorderRecorder[integrationConfig](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	cmd, err := v5.NewCommand(
		"fail",
		v5.NoFlags{},
		func(_ context.Context, _ *integrationConfig, _ v5.NoFlags) error {
			return errors.New("boom")
		},
		v5.WithShort("Always fails"),
	)
	if err != nil {
		t.Fatalf("NewCommand: %v", err)
	}

	if err := v5.AddCommand(cli, cmd); err != nil {
		t.Fatalf("AddCommand: %v", err)
	}

	// Execute the failing command — the error return is expected.
	_ = cli.ExecuteWithArgs(context.Background(), []string{"fail"})

	path := waitTraceFile(t, dir)
	if path == "" {
		t.Fatal("expected a trace snapshot file after command error")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat snapshot: %v", err)
	}

	if info.Size() == 0 {
		t.Error("snapshot file is empty")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "-error-") {
		t.Errorf("filename %q should contain '-error-' reason", base)
	}
}

// TestIntegration_CaptureOnSlow verifies that a slow command triggers a snapshot
// through the full CLI.Execute() flow.
func TestIntegration_CaptureOnSlow(t *testing.T) {
	dir := t.TempDir()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:        1 * time.Second,
		MaxBytes:      1 << 20,
		OutputDir:     dir,
		CaptureOnSlow: true,
		SlowThreshold: 1 * time.Millisecond,
	})

	defer rec.Stop()

	cli, err := v5.NewCLI[integrationConfig](
		"testapp", "Test app", integrationConfig{},
		flightrecorder.WithFlightRecorderRecorder[integrationConfig](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	cmd, err := v5.NewCommand(
		"slow",
		v5.NoFlags{},
		func(_ context.Context, _ *integrationConfig, _ v5.NoFlags) error {
			time.Sleep(5 * time.Millisecond)

			return nil
		},
		v5.WithShort("Slow command"),
	)
	if err != nil {
		t.Fatalf("NewCommand: %v", err)
	}

	if err := v5.AddCommand(cli, cmd); err != nil {
		t.Fatalf("AddCommand: %v", err)
	}

	if err := cli.ExecuteWithArgs(context.Background(), []string{"slow"}); err != nil {
		t.Fatalf("ExecuteWithArgs: %v", err)
	}

	path := waitTraceFile(t, dir)
	if path == "" {
		t.Fatal("expected a trace snapshot file after slow command")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "-slow-") {
		t.Errorf("filename %q should contain '-slow-' reason", base)
	}
}

// TestIntegration_NoCaptureOnSuccess verifies that a fast, successful command
// does NOT produce a snapshot file.
func TestIntegration_NoCaptureOnSuccess(t *testing.T) {
	dir := t.TempDir()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Hour, // effectively never slow
		CaptureOnError: true,
	})

	defer rec.Stop()

	cli, err := v5.NewCLI[integrationConfig](
		"testapp", "Test app", integrationConfig{},
		flightrecorder.WithFlightRecorderRecorder[integrationConfig](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	cmd, err := v5.NewCommand(
		"ok",
		v5.NoFlags{},
		func(_ context.Context, _ *integrationConfig, _ v5.NoFlags) error {
			return nil
		},
		v5.WithShort("Succeeds fast"),
	)
	if err != nil {
		t.Fatalf("NewCommand: %v", err)
	}

	if err := v5.AddCommand(cli, cmd); err != nil {
		t.Fatalf("AddCommand: %v", err)
	}

	if err := cli.ExecuteWithArgs(context.Background(), []string{"ok"}); err != nil {
		t.Fatalf("ExecuteWithArgs: %v", err)
	}

	time.Sleep(100 * time.Millisecond) // give async goroutine time to NOT write

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected no snapshot files, got %d", len(entries))
	}
}

// waitTraceFile polls directory for up to 500ms, returning the first file path
// that has non-zero content (ensuring the async write has completed).
func waitTraceFile(t *testing.T, dir string) string {
	t.Helper()

	deadline := time.Now().Add(500 * time.Millisecond)

	for time.Now().Before(deadline) {
		entries, err := os.ReadDir(dir)
		if err == nil && len(entries) > 0 {
			path := filepath.Join(dir, entries[0].Name())

			info, statErr := os.Stat(path)
			if statErr == nil && info.Size() > 0 {
				return path
			}
		}

		time.Sleep(5 * time.Millisecond)
	}

	return ""
}
