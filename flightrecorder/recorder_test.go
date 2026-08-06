package flightrecorder

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/trace"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultConfig()

	if cfg.MinAge != defaultMinAge {
		t.Errorf("MinAge = %v, want %v", cfg.MinAge, defaultMinAge)
	}

	if cfg.MaxBytes != defaultMaxBytes {
		t.Errorf("MaxBytes = %v, want %v", cfg.MaxBytes, defaultMaxBytes)
	}

	if !cfg.CaptureOnSlow {
		t.Error("CaptureOnSlow should be true by default")
	}

	if cfg.CaptureOnError {
		t.Error("CaptureOnError should be false by default")
	}

	if cfg.SlowThreshold != defaultSlowThreshold {
		t.Errorf("SlowThreshold = %v, want %v", cfg.SlowThreshold, defaultSlowThreshold)
	}

	if cfg.FilenamePrefix != defaultPrefix {
		t.Errorf("FilenamePrefix = %q, want %q", cfg.FilenamePrefix, defaultPrefix)
	}
}

func TestConfig_WithDefaults(t *testing.T) {
	t.Parallel()

	t.Run("empty fields_get_defaults", func(t *testing.T) {
		t.Parallel()

		cfg := Config{}.withDefaults()

		if cfg.MinAge != defaultMinAge {
			t.Errorf("MinAge = %v, want %v", cfg.MinAge, defaultMinAge)
		}

		if cfg.MaxBytes != defaultMaxBytes {
			t.Errorf("MaxBytes = %v, want %v", cfg.MaxBytes, defaultMaxBytes)
		}

		if cfg.SlowThreshold != defaultSlowThreshold {
			t.Errorf("SlowThreshold = %v, want %v", cfg.SlowThreshold, defaultSlowThreshold)
		}

		if cfg.FilenamePrefix != defaultPrefix {
			t.Errorf("FilenamePrefix = %q, want %q", cfg.FilenamePrefix, defaultPrefix)
		}
	})

	t.Run("non-zero_fields_preserved", func(t *testing.T) {
		t.Parallel()

		cfg := Config{
			MinAge:         42 * time.Second,
			MaxBytes:       999,
			SlowThreshold:  7 * time.Second,
			FilenamePrefix: "myapp",
		}.withDefaults()

		if cfg.MinAge != 42*time.Second {
			t.Errorf("MinAge = %v, want 42s", cfg.MinAge)
		}

		if cfg.MaxBytes != 999 {
			t.Errorf("MaxBytes = %v, want 999", cfg.MaxBytes)
		}

		if cfg.SlowThreshold != 7*time.Second {
			t.Errorf("SlowThreshold = %v, want 7s", cfg.SlowThreshold)
		}

		if cfg.FilenamePrefix != "myapp" {
			t.Errorf("FilenamePrefix = %q, want %q", cfg.FilenamePrefix, "myapp")
		}
	})
}

func TestNew_AppliesDefaults(t *testing.T) {
	t.Parallel()

	rec := New(Config{})

	if rec.Config().MinAge != defaultMinAge {
		t.Errorf("MinAge = %v, want %v", rec.Config().MinAge, defaultMinAge)
	}

	if rec.Config().FilenamePrefix != defaultPrefix {
		t.Errorf("FilenamePrefix = %q, want %q", rec.Config().FilenamePrefix, defaultPrefix)
	}
}

func TestNew_PreservesExplicitConfig(t *testing.T) {
	t.Parallel()

	rec := New(Config{
		MinAge:         10 * time.Second,
		FilenamePrefix: "custom",
		CaptureOnError: true,
	})

	if rec.Config().MinAge != 10*time.Second {
		t.Errorf("MinAge = %v, want 10s", rec.Config().MinAge)
	}

	if rec.Config().FilenamePrefix != "custom" {
		t.Errorf("FilenamePrefix = %q, want %q", rec.Config().FilenamePrefix, "custom")
	}

	if !rec.Config().CaptureOnError {
		t.Error("CaptureOnError should be true")
	}
}

func TestRecorder_StartStopLifecycle(t *testing.T) {
	rec := New(DefaultConfig())

	if rec.Enabled() {
		t.Fatal("new recorder should not be enabled")
	}

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if !rec.Enabled() {
		t.Fatal("recorder should be enabled after Start")
	}

	// Double-start should error.
	if err := rec.Start(); err == nil {
		t.Fatal("double Start should return error")
	}

	rec.Stop()

	if rec.Enabled() {
		t.Fatal("recorder should not be enabled after Stop")
	}

	// Double-stop should be safe.
	rec.Stop()
}

func TestRecorder_CanRestartAfterStop(t *testing.T) {
	rec := New(DefaultConfig())

	if err := rec.Start(); err != nil {
		t.Fatalf("first Start failed: %v", err)
	}

	rec.Stop()

	if err := rec.Start(); err != nil {
		t.Fatalf("second Start failed: %v", err)
	}

	rec.Stop()
}

func TestRecorder_WriteTo_DisabledReturnsError(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())

	var buf bytes.Buffer

	_, err := rec.WriteTo(&buf)
	if err == nil {
		t.Fatal("WriteTo on disabled recorder should return error")
	}
}

func TestRecorder_WriteTo_EnabledWritesData(t *testing.T) {
	rec := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	// Generate some trace data by starting a region.
	ctx, task := trace.NewTask(context.Background(), "test-task")

	trace.WithRegion(ctx, "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	task.End()

	var buf bytes.Buffer

	written, err := rec.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	// We just verify we got *some* data.
	if written == 0 {
		t.Error("WriteTo wrote 0 bytes, expected trace data")
	}

	if buf.Len() == 0 {
		t.Error("buffer is empty after WriteTo")
	}
}

func TestRecorder_Capture_DisabledReturnsError(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())

	path, err := rec.Capture(context.Background(), "test", CaptureReasonSlow)
	if err == nil {
		t.Fatalf("Capture on disabled recorder should return error, got path: %s", path)
	}
}

func TestRecorder_Capture_WritesFileWithExpectedName(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		FilenamePrefix: "testprefix",
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	// Generate some trace data.
	ctx, task := trace.NewTask(context.Background(), "capture-test")

	trace.WithRegion(ctx, "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	task.End()

	path, err := rec.Capture(context.Background(), "deploy", CaptureReasonSlow)
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	if path == "" {
		t.Fatal("Capture returned empty path")
	}

	// Verify file exists.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}

	if info.Size() == 0 {
		t.Error("snapshot file is empty")
	}

	// Verify filename contains expected components.
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "testprefix-deploy-slow-") {
		t.Errorf("filename %q does not match expected pattern", base)
	}

	if !strings.HasSuffix(base, ".trace") {
		t.Errorf("filename %q should end with .trace", base)
	}
}

func TestRecorder_Capture_UsesFullPathWhenAvailable(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	// Generate minimal trace data.
	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	// Simulate a command path with spaces.
	path, err := rec.Capture(context.Background(), "myapp deploy prod", CaptureReasonError)
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	base := filepath.Base(path)

	// Spaces should be sanitized to hyphens.
	if !strings.Contains(base, "myapp-deploy-prod") {
		t.Errorf("filename %q should contain sanitized path", base)
	}
}

func TestRecorder_Capture_CancelledContext(t *testing.T) {
	rec := New(DefaultConfig())

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := rec.Capture(ctx, "test", CaptureReasonSlow)
	if err == nil {
		t.Fatal("Capture with cancelled context should return error")
	}
}

func TestRecorder_Capture_CreatesOutputDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deep")

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	// Generate trace data.
	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	path, err := rec.Capture(context.Background(), "test", CaptureReasonSlow)
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("snapshot file not found after Capture: %v", err)
	}
}

func TestRecorder_Capture_UsesTempDirWhenOutputDirEmpty(t *testing.T) {
	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: "",
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	path, err := rec.Capture(context.Background(), "test", CaptureReasonSlow)
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	// File should be in os.TempDir().
	if !strings.HasPrefix(path, os.TempDir()) {
		t.Errorf("snapshot path %q should be in temp dir %q", path, os.TempDir())
	}

	_ = os.Remove(path)
}

func TestRecorder_Log_UsesCustomLogFunc(t *testing.T) {
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

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	_, _ = rec.Capture(context.Background(), "test", CaptureReasonSlow)

	mutex.Lock()
	defer mutex.Unlock()

	foundCaptureLog := false

	for _, msg := range logged {
		if strings.Contains(msg, "captured") {
			foundCaptureLog = true

			break
		}
	}

	if !foundCaptureLog {
		t.Errorf("expected a 'captured' log message, got: %v", logged)
	}
}

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{input: "deploy", want: "deploy"},
		{input: "deploy prod", want: "deploy-prod"},
		{input: "app/deploy", want: "app-deploy"},
		{input: "myapp deploy prod", want: "myapp-deploy-prod"},
		{input: "command_1.0", want: "command_1.0"},
		{input: "", want: ""},
		{input: "UPPER", want: "UPPER"},
		{input: "café", want: "caf-"},
	}

	for _, table := range tests {
		t.Run(table.input, func(t *testing.T) {
			t.Parallel()

			got := sanitizeFilename(table.input)
			if got != table.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", table.input, got, table.want)
			}
		})
	}
}

func TestCaptureReason_Constants(t *testing.T) {
	t.Parallel()

	if CaptureReasonSlow != "slow" {
		t.Errorf("CaptureReasonSlow = %q, want %q", CaptureReasonSlow, "slow")
	}

	if CaptureReasonError != "error" {
		t.Errorf("CaptureReasonError = %q, want %q", CaptureReasonError, "error")
	}
}

func TestRecorder_ErrorsAreSentinels(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())

	_, err := rec.WriteTo(&bytes.Buffer{})

	if !errors.Is(err, ErrNotEnabled) {
		t.Errorf("expected ErrNotEnabled, got: %v", err)
	}

	_, captureErr := rec.Capture(context.Background(), "test", CaptureReasonSlow)
	if !errors.Is(captureErr, ErrNotEnabled) {
		t.Errorf("expected ErrNotEnabled from Capture, got: %v", captureErr)
	}
}

func TestRecorder_DoubleStartReturnsErrAlreadyStarted(t *testing.T) {
	rec := New(DefaultConfig())

	if err := rec.Start(); err != nil {
		t.Fatalf("first Start failed: %v", err)
	}

	defer rec.Stop()

	err := rec.Start()
	if !errors.Is(err, ErrAlreadyStarted) {
		t.Errorf("double Start should return ErrAlreadyStarted, got: %v", err)
	}
}

func TestRecorder_WriteTo_AfterStopReturnsErrNotEnabled(t *testing.T) {
	rec := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	rec.Stop()

	var buf bytes.Buffer

	_, err := rec.WriteTo(&buf)
	if !errors.Is(err, ErrNotEnabled) {
		t.Errorf("WriteTo after Stop should return ErrNotEnabled, got: %v", err)
	}
}

func TestRecorder_Capture_EmptyCommandNameFallback(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	path, err := rec.Capture(context.Background(), "", CaptureReasonSlow)
	if err != nil {
		t.Fatalf("Capture with empty name failed: %v", err)
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "command-slow-") {
		t.Errorf("filename %q should contain 'command-slow-' (empty name fallback)", base)
	}
}

func TestRecorder_Capture_ConcurrentCallsSafe(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	const goroutines = 10

	var waitGroup sync.WaitGroup

	for range goroutines {
		waitGroup.Go(func() {
			_, _ = rec.Capture(context.Background(), "concurrent", CaptureReasonSlow)
		})
	}

	waitGroup.Wait()

	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Error("expected at least one snapshot file from concurrent captures")
	}
}

func TestRecorder_CaptureToWriter_WritesData(t *testing.T) {
	rec := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	var buf bytes.Buffer

	written, err := rec.CaptureToWriter(context.Background(), &buf, "test", CaptureReasonError)
	if err != nil {
		t.Fatalf("CaptureToWriter failed: %v", err)
	}

	if written == 0 {
		t.Error("CaptureToWriter wrote 0 bytes, expected trace data")
	}

	if buf.Len() == 0 {
		t.Error("buffer is empty after CaptureToWriter")
	}
}

func TestRecorder_CaptureToWriter_DisabledReturnsError(t *testing.T) {
	t.Parallel()

	rec := New(DefaultConfig())

	var buf bytes.Buffer

	_, err := rec.CaptureToWriter(context.Background(), &buf, "test", CaptureReasonSlow)
	if !errors.Is(err, ErrNotEnabled) {
		t.Errorf("CaptureToWriter on disabled recorder should return ErrNotEnabled, got: %v", err)
	}
}

func TestRecorder_CaptureToWriter_CancelledContext(t *testing.T) {
	rec := New(Config{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var buf bytes.Buffer

	_, err := rec.CaptureToWriter(ctx, &buf, "test", CaptureReasonSlow)
	if err == nil {
		t.Fatal("CaptureToWriter with cancelled context should return error")
	}
}

// TestTraceSnapshot_IsParseableByGoToolTrace validates that snapshots produced
// by the flight recorder can be parsed by `go tool trace`. This makes the
// manual M10 validation repeatable in CI.
func TestTraceSnapshot_IsParseableByGoToolTrace(t *testing.T) {
	dir := t.TempDir()

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "validation-work", func() {
		time.Sleep(2 * time.Millisecond)
	})

	path, err := rec.Capture(context.Background(), "validation", CaptureReasonSlow)
	if err != nil {
		t.Fatalf("Capture failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}

	if info.Size() < 100 {
		t.Fatalf("snapshot file too small (%d bytes), expected trace data", info.Size())
	}

	// Run `go tool trace` with a short timeout. If parsing fails, the tool
	// exits quickly with a non-zero status. If parsing succeeds, it starts a
	// web server and blocks — the timeout kills it, which we treat as success.
	// We discard output to prevent pipe-buffer deadlocks.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("failed to open /dev/null: %v", err)
	}
	defer devNull.Close()

	cmd := exec.CommandContext(ctx, "go", "tool", "trace", path)
	cmd.Stdout = devNull
	cmd.Stderr = devNull

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start go tool trace: %v", err)
	}

	err = cmd.Wait()

	// DeadlineExceeded means the tool was still running after 5 seconds
	// (parsing succeeded, web server started). Any other error means either
	// the tool failed during parsing or was killed by a signal.
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("go tool trace failed during parsing (exit before timeout): %v", err)
	}
}
