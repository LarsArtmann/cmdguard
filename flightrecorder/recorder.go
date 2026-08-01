// Package flightrecorder provides Go runtime execution trace recording for
// cmdguard CLIs via the flight recorder pattern (Go 1.25+ runtime/trace).
//
// The flight recorder continuously buffers the last few seconds of execution
// trace in memory. When something goes wrong — a slow command, an error, or a
// panic — you can snapshot the buffer to capture exactly the problematic window
// of time. The resulting .trace file can be analyzed with `go tool trace`.
//
// This is an optional module — it has zero external dependencies (uses only the
// Go standard library runtime/trace package). Import it only when you need
// execution trace diagnostics.
//
// # Quick start
//
// Wire it as a single CLIOption and every command that exceeds the slow
// threshold or returns an error will automatically produce a .trace file:
//
//	import (
//	    v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
//	    "github.com/larsartmann/cmdguard/flightrecorder"
//	)
//
//	cli, _ := v4.NewCLI[Config]("app", "My app", Config{},
//	    flightrecorder.WithFlightRecorder[Config](flightrecorder.Config{
//	        CaptureOnSlow:    true,
//	        SlowThreshold:    200 * time.Millisecond,
//	        CaptureOnError:   true,
//	    }),
//	)
//
// # Manual control
//
// For advanced use cases (custom capture triggers, shared recorder across
// multiple CLIs), use the Recorder and Middleware types directly:
//
//	rec := flightrecorder.New(flightrecorder.DefaultConfig())
//	if err := rec.Start(); err != nil {
//	    log.Printf("flight recorder: %v", err)
//	}
//	defer rec.Stop()
//
//	cli, _ := v4.NewCLI[Config]("app", "My app", Config{},
//	    v4.WithMiddleware(flightrecorder.Middleware[Config](rec)),
//	)
//
// # Analyzing traces
//
// After a snapshot is captured, analyze it with:
//
//	go tool trace /path/to/snapshot.trace
//
// This launches a local web server with an interactive trace viewer.
package flightrecorder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/trace"
	"sync"
	"time"
)

// Default values for Config fields. See DefaultConfig.
const (
	defaultMinAge        = 5 * time.Second
	defaultMaxBytes      = uint64(10 << 20) // 10 MiB — covers ~1s of busy-service trace data
	defaultSlowThreshold = 100 * time.Millisecond
	defaultPrefix        = "cmdguard"
	dirPermission        = 0o755
)

// Sentinel errors for Recorder operations.
var (
	// ErrAlreadyStarted is returned when Start is called on an already-running recorder.
	ErrAlreadyStarted = errors.New("flightrecorder: recorder already started")

	// ErrNotEnabled is returned when an operation requires an active recorder but none is running.
	ErrNotEnabled = errors.New("flightrecorder: recorder not enabled")
)

// CaptureReason describes why a snapshot was taken.
type CaptureReason string

const (
	// CaptureReasonSlow indicates the command exceeded the configured slow threshold.
	CaptureReasonSlow CaptureReason = "slow"

	// CaptureReasonError indicates the command returned an error.
	CaptureReasonError CaptureReason = "error"
)

// Config configures the flight recorder behavior.
type Config struct {
	// MinAge is the duration for which trace data is reliably retained in the
	// in-memory buffer. Older data may be evicted. Suggested: ~2x the time
	// window of the event you are debugging (e.g. if debugging a 5-second
	// timeout, set to 10 seconds).
	//
	// Default: 5 * time.Second
	MinAge time.Duration

	// MaxBytes limits the in-memory buffer size to prevent unbounded memory
	// growth. Expect roughly 10 MB/s of trace data for a busy service.
	//
	// Default: 10 MiB (10 << 20)
	MaxBytes uint64

	// CaptureOnSlow controls whether a snapshot is automatically captured when
	// a command's run phase exceeds SlowThreshold.
	//
	// Default: true
	CaptureOnSlow bool

	// SlowThreshold is the duration above which a command is considered slow
	// and triggers a snapshot (when CaptureOnSlow is true).
	//
	// Default: 100 * time.Millisecond
	SlowThreshold time.Duration

	// CaptureOnError controls whether a snapshot is automatically captured when
	// a command returns a non-nil error from its run phase.
	//
	// Default: false
	CaptureOnError bool

	// OutputDir is the directory where snapshot .trace files are written.
	// The directory is created if it does not exist. If empty, uses os.TempDir().
	//
	// Default: "" (resolves to os.TempDir())
	OutputDir string

	// FilenamePrefix is prepended to snapshot filenames. Snapshot files are
	// named: {prefix}-{command}-{reason}-{timestamp}.trace
	//
	// Default: "cmdguard"
	FilenamePrefix string

	// Log receives diagnostic messages (e.g. snapshot captured, start failed).
	// If nil, messages are written to os.Stderr.
	//
	// Default: nil (os.Stderr)
	Log func(format string, args ...any)
}

// DefaultConfig returns a Config with sensible defaults for typical CLI usage.
func DefaultConfig() Config {
	return Config{
		MinAge:         defaultMinAge,
		MaxBytes:       defaultMaxBytes,
		CaptureOnSlow:  true,
		SlowThreshold:  defaultSlowThreshold,
		CaptureOnError: false,
		OutputDir:      "",
		FilenamePrefix: defaultPrefix,
		Log:            nil,
	}
}

// withDefaults applies default values to zero-valued fields.
func (c Config) withDefaults() Config {
	result := c

	if result.MinAge == 0 {
		result.MinAge = defaultMinAge
	}

	if result.MaxBytes == 0 {
		result.MaxBytes = defaultMaxBytes
	}

	if result.SlowThreshold == 0 {
		result.SlowThreshold = defaultSlowThreshold
	}

	if result.FilenamePrefix == "" {
		result.FilenamePrefix = defaultPrefix
	}

	return result
}

// Recorder wraps runtime/trace.FlightRecorder with snapshot management.
// It is safe for concurrent use.
//
// At most one flight recorder may be active per process at any given time
// (runtime/trace limitation). If you create multiple Recorders, only the first
// Start will succeed; subsequent starts return an error.
type Recorder struct {
	fr  *trace.FlightRecorder
	cfg Config

	mu       sync.Mutex
	started  bool
	inflight sync.WaitGroup // tracks in-flight WriteTo/Capture operations
}

// New creates a Recorder from the given Config. The recorder is NOT started —
// call Start to begin recording. Zero-valued Config fields are replaced with
// defaults.
func New(cfg Config) *Recorder {
	return &Recorder{
		fr: trace.NewFlightRecorder(trace.FlightRecorderConfig{
			MinAge:   cfg.MinAge,
			MaxBytes: cfg.MaxBytes,
		}),
		cfg:      cfg.withDefaults(),
		mu:       sync.Mutex{},
		started:  false,
		inflight: sync.WaitGroup{},
	}
}

// Config returns the effective configuration (with defaults applied).
func (rec *Recorder) Config() Config {
	return rec.cfg
}

// Start begins recording execution traces into the in-memory buffer.
// Returns an error if the recorder is already started or if the runtime
// rejects the start (e.g. another flight recorder is already active).
func (rec *Recorder) Start() error {
	rec.mu.Lock()
	defer rec.mu.Unlock()

	if rec.started {
		return ErrAlreadyStarted
	}

	if err := rec.fr.Start(); err != nil {
		return fmt.Errorf("flightrecorder: starting flight recorder: %w", err)
	}

	rec.started = true

	return nil
}

// Stop stops recording. Safe to call multiple times. After stopping, the
// recorder can be restarted with Start.
//
// Stop waits for any in-flight snapshot capture (WriteTo/Capture) to complete
// before stopping the underlying recorder, preventing races between snapshot
// writes and recorder shutdown.
func (rec *Recorder) Stop() {
	rec.mu.Lock()

	if !rec.started {
		rec.mu.Unlock()

		return
	}

	rec.started = false
	rec.mu.Unlock()

	// Wait for in-flight WriteTo/Capture operations to finish before
	// calling fr.Stop(), which races with fr.WriteTo in the runtime.
	rec.inflight.Wait()
	rec.fr.Stop()
}

// Enabled reports whether the flight recorder is actively recording.
func (rec *Recorder) Enabled() bool {
	rec.mu.Lock()
	defer rec.mu.Unlock()

	return rec.started
}

// WriteTo writes the current contents of the trace buffer to writer.
// Returns the number of bytes written. If the recorder is not enabled, returns
// an error.
func (rec *Recorder) WriteTo(writer io.Writer) (int64, error) {
	rec.mu.Lock()

	if !rec.started {
		rec.mu.Unlock()

		return 0, ErrNotEnabled
	}

	rec.inflight.Add(1)

	defer rec.inflight.Done()

	rec.mu.Unlock()

	written, err := rec.fr.WriteTo(writer)
	if err != nil {
		return written, fmt.Errorf("flightrecorder: writing trace snapshot: %w", err)
	}

	return written, nil
}

// Capture writes a trace snapshot to a file in OutputDir.
// The filename includes the command name, capture reason, and timestamp.
// Returns the path to the written file.
//
// If the recorder is not enabled, Capture returns an error without writing.
func (rec *Recorder) Capture(ctx context.Context, commandName string, reason CaptureReason) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("flightrecorder: capture cancelled: %w", err)
	}

	if !rec.Enabled() {
		return "", ErrNotEnabled
	}

	path, err := rec.buildSnapshotPath(commandName, reason)
	if err != nil {
		return "", err
	}

	written, err := rec.writeSnapshot(path)
	if err != nil {
		return "", err
	}

	rec.logf("flightrecorder: captured %s snapshot (%d bytes) to %s", reason, written, path)

	return path, nil
}

// CaptureToWriter writes a trace snapshot to the given writer instead of a file.
// This decouples snapshot writing from the filesystem, useful for piping to
// stdout, network storage, or testing. Returns the number of bytes written.
//
// If the recorder is not enabled, CaptureToWriter returns an error without writing.
func (rec *Recorder) CaptureToWriter(
	ctx context.Context,
	writer io.Writer,
	commandName string,
	reason CaptureReason,
) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("flightrecorder: capture cancelled: %w", err)
	}

	if !rec.Enabled() {
		return 0, ErrNotEnabled
	}

	written, err := rec.WriteTo(writer)
	if err != nil {
		return written, err
	}

	rec.logf("flightrecorder: captured %s snapshot (%d bytes) to writer for %s", reason, written, commandName)

	return written, nil
}

// buildSnapshotPath resolves the output directory (creating it if needed) and
// constructs the snapshot filename from the command name and capture reason.
func (rec *Recorder) buildSnapshotPath(commandName string, reason CaptureReason) (string, error) {
	dir := rec.cfg.OutputDir
	if dir == "" {
		dir = os.TempDir()
	}

	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return "", fmt.Errorf("flightrecorder: creating output directory %q: %w", dir, err)
	}

	safeName := sanitizeFilename(commandName)
	if safeName == "" {
		safeName = "command"
	}

	filename := fmt.Sprintf(
		"%s-%s-%s-%s.trace",
		rec.cfg.FilenamePrefix,
		safeName,
		reason,
		time.Now().Format("20060102-150405.000000000"),
	)

	return filepath.Join(dir, filename), nil
}

// writeSnapshot creates the file at path, writes the trace buffer to it,
// and removes the partial file if the write fails.
func (rec *Recorder) writeSnapshot(path string) (int64, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("flightrecorder: creating snapshot file %q: %w", path, err)
	}

	defer func() { _ = file.Close() }()

	written, err := rec.WriteTo(file)
	if err != nil {
		_ = os.Remove(path) // best-effort cleanup of partial file

		return written, err
	}

	return written, nil
}

// logf writes a diagnostic message via the configured Log function,
// defaulting to os.Stderr.
func (rec *Recorder) logf(format string, args ...any) {
	if rec.cfg.Log != nil {
		rec.cfg.Log(format, args...)

		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// sanitizeFilename replaces characters that are unsafe in filenames with hyphens.
func sanitizeFilename(input string) string {
	var builder []rune

	for _, char := range input {
		switch {
		case char >= 'a' && char <= 'z',
			char >= 'A' && char <= 'Z',
			char >= '0' && char <= '9',
			char == '-', char == '_', char == '.':
			builder = append(builder, char)
		default:
			builder = append(builder, '-')
		}
	}

	return string(builder)
}
