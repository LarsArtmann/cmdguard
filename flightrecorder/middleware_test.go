package flightrecorder

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type testConfig struct {
	Name string `flag:"name" short:"n" default:"" help:"Your name"`
}

func TestMiddleware_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())
	middleware := Middleware[testConfig](rec)
	if middleware == nil {
		t.Fatal("Middleware returned nil")
	}
}

func TestMiddleware_CallsNextAndPassesError(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	info := v4.CommandInfo{Name: "test", Phase: v4.PhaseRun}

	expectedErr := errors.New("command failed")
	called := false
	middleware := Middleware[testConfig](rec)

	err := middleware(ctx, cfg, info, func() error {
		called = true

		return expectedErr
	})

	if !called {
		t.Fatal("next function was not called")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("middleware should pass through error, got: %v", err)
	}
}

func TestMiddleware_NonRunPhaseDoesNotCapture(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Nanosecond,
		CaptureOnError: true,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	// Pre-run phase: should NOT capture even though slow/error conditions are met.
	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "cmd", Phase: v4.PhasePreRun}, func() error {
		return errors.New("error")
	})

	// Post-run phase: should NOT capture.
	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "cmd", Phase: v4.PhasePostRun}, func() error {
		return errors.New("error")
	})

	// Give goroutines time to write.
	time.Sleep(50 * time.Millisecond)

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("non-run phases should not capture snapshots, got %d files", len(entries))
	}
}

func TestMiddleware_CapturesOnSlow(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:        1 * time.Second,
		MaxBytes:      1 << 20,
		OutputDir:     dir,
		CaptureOnSlow: true,
		SlowThreshold: 1 * time.Millisecond,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "slowcmd", Phase: v4.PhaseRun}, func() error {
		time.Sleep(5 * time.Millisecond) // exceeds 1ms threshold

		return nil
	})

	// Wait for async capture goroutine.
	path := waitForFile(t, dir)
	if path == "" {
		t.Fatal("expected a snapshot file to be captured")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "slowcmd-slow-") {
		t.Errorf("filename %q should contain 'slowcmd-slow-'", base)
	}
}

func TestMiddleware_CapturesOnError(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  false, // disable slow capture to isolate error capture
		CaptureOnError: true,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	cmdErr := errors.New("boom")
	err := middleware(ctx, cfg, v4.CommandInfo{Name: "errcmd", Phase: v4.PhaseRun}, func() error {
		return cmdErr
	})

	if !errors.Is(err, cmdErr) {
		t.Fatalf("middleware should pass through error, got: %v", err)
	}

	path := waitForFile(t, dir)
	if path == "" {
		t.Fatal("expected a snapshot file to be captured on error")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "errcmd-error-") {
		t.Errorf("filename %q should contain 'errcmd-error-'", base)
	}
}

func TestMiddleware_NoCaptureWhenConditionsNotMet(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Hour, // very high threshold
		CaptureOnError: false,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "fastcmd", Phase: v4.PhaseRun}, func() error {
		return nil // fast, no error
	})

	time.Sleep(50 * time.Millisecond)

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("should not capture when conditions not met, got %d files", len(entries))
	}
}

func TestMiddleware_UsesFullPathWhenAvailable(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:        1 * time.Second,
		MaxBytes:      1 << 20,
		OutputDir:     dir,
		CaptureOnSlow: true,
		SlowThreshold: 1 * time.Millisecond,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	_ = middleware(ctx, cfg, v4.CommandInfo{
		Name:     "deploy",
		FullPath: "myapp deploy",
		Phase:    v4.PhaseRun,
	}, func() error {
		time.Sleep(5 * time.Millisecond)

		return nil
	})

	path := waitForFile(t, dir)
	if path == "" {
		t.Fatal("expected a snapshot file")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "myapp-deploy-slow-") {
		t.Errorf("filename %q should use full path with sanitized name", base)
	}
}

func TestMiddleware_LazilyStartsRecorder(t *testing.T) {
	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		CaptureOnSlow:  false,
		CaptureOnError: false,
	})

	if rec.Enabled() {
		t.Fatal("recorder should not be started before middleware runs")
	}

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "test", Phase: v4.PhaseRun}, func() error {
		return nil
	})

	if !rec.Enabled() {
		t.Fatal("recorder should be started after middleware runs")
	}

	rec.Stop()
}

func TestWithFlightRecorder_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	opt := WithFlightRecorder[testConfig](DefaultConfig())
	if opt == nil {
		t.Fatal("WithFlightRecorder returned nil")
	}
}

func TestWithFlightRecorderRecorder_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())
	opt := WithFlightRecorderRecorder[testConfig](rec)
	if opt == nil {
		t.Fatal("WithFlightRecorderRecorder returned nil")
	}
}

func TestWithFlightRecorder_WithCustomConfig(t *testing.T) {
	t.Parallel()

	opt := WithFlightRecorder[testConfig](Config{
		MinAge:         10 * time.Second,
		CaptureOnError: true,
		OutputDir:      "/tmp/test-fr",
	})
	if opt == nil {
		t.Fatal("WithFlightRecorder returned nil for custom config")
	}
}

// waitForFile polls directory for up to 500ms, returning the first file path found.
func waitForFile(t *testing.T, directory string) string {
	t.Helper()

	deadline := time.Now().Add(500 * time.Millisecond)

	for time.Now().Before(deadline) {
		entries, err := os.ReadDir(directory)
		if err == nil && len(entries) > 0 {
			return filepath.Join(directory, entries[0].Name())
		}

		time.Sleep(5 * time.Millisecond)
	}

	return ""
}

func TestMiddleware_ErrorTakesPrecedenceOverSlow(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Nanosecond,
		CaptureOnError: true,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	// Command is both slow (exceeds 1ns threshold) AND errors.
	// The snapshot reason should be "error" (takes precedence over "slow").
	_ = middleware(ctx, cfg, v4.CommandInfo{Name: "both", Phase: v4.PhaseRun}, func() error {
		return errors.New("failed")
	})

	path := waitForFile(t, dir)
	if path == "" {
		t.Fatal("expected a snapshot file")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "both-error-") {
		t.Errorf("filename %q should use 'error' reason (takes precedence over 'slow')", base)
	}
}

// TestMiddleware_StartFailure_LogsAndContinues verifies that when the recorder
// fails to start (because another flight recorder is already active), the
// middleware logs the failure and still calls next(). This covers the
// previously-untested Start() error branch (middleware.go lines 28-32).
func TestMiddleware_StartFailure_LogsAndContinues(t *testing.T) {
	// Occupy the process-wide trace singleton so the second recorder cannot start.
	blocker := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	if err := blocker.Start(); err != nil {
		t.Fatalf("blocker Start failed: %v", err)
	}

	defer blocker.Stop()

	var logged []string

	var mutex sync.Mutex

	rec := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
		Log: func(format string, _ ...any) {
			mutex.Lock()
			logged = append(logged, format)
			mutex.Unlock()
		},
	})

	ctx := context.Background()
	cfg := &testConfig{}
	middleware := Middleware[testConfig](rec)

	nextCalled := false

	err := middleware(ctx, cfg, v4.CommandInfo{Name: "test", Phase: v4.PhaseRun}, func() error {
		nextCalled = true

		return nil
	})

	if !nextCalled {
		t.Fatal("middleware should call next() even when Start fails")
	}

	if err != nil {
		t.Fatalf("middleware should return nil (next returned nil), got: %v", err)
	}

	mutex.Lock()
	defer mutex.Unlock()

	foundStartFailure := false

	for _, msg := range logged {
		if strings.Contains(msg, "failed to start") {
			foundStartFailure = true

			break
		}
	}

	if !foundStartFailure {
		t.Errorf("expected a 'failed to start' log message, got: %v", logged)
	}
}
