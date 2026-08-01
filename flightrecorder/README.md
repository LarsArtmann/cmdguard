# flightrecorder

Go runtime execution trace flight recorder middleware for [cmdguard](https://github.com/LarsArtmann/cmdguard).

Wraps Go 1.25+ [`runtime/trace.FlightRecorder`](https://pkg.go.dev/runtime/trace#FlightRecorder) to continuously buffer execution traces in memory and automatically capture `.trace` snapshots when commands are slow or error. **Zero external dependencies** — uses only the Go standard library.

## Installation

```bash
go get github.com/larsartmann/cmdguard/flightrecorder
```

## Quick Start

```go
package main

import (
    "time"

    "github.com/larsartmann/cmdguard/flightrecorder"
    v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type Config struct{}

func main() {
    cli, _ := v4.NewCLI[Config]("myapp", "My CLI", Config{},
        flightrecorder.WithFlightRecorder[Config](flightrecorder.Config{
            CaptureOnSlow:  true,
            SlowThreshold:  200 * time.Millisecond,
            CaptureOnError: true,
            OutputDir:      "/tmp/myapp-traces",
        }),
    )
    _ = cli.Execute()
}
```

Snapshots are written as `{prefix}-{command}-{reason}-{timestamp}.trace` files in `OutputDir` (defaults to `os.TempDir()`). Analyze them with:

```bash
go tool trace /tmp/myapp-traces/cmdguard-deploy-error-20260801-154032.000000000.trace
```

## Manual Lifecycle Control

For explicit Start/Stop control (e.g. multiple CLIs sharing one recorder):

```go
rec := flightrecorder.New(flightrecorder.DefaultConfig())

if err := rec.Start(); err != nil {
    log.Fatal(err)
}
defer rec.Stop()

// Capture a snapshot at any time:
path, err := rec.Capture(ctx, "manual-snapshot", flightrecorder.CaptureReasonSlow)

// Or write to an arbitrary io.Writer:
n, err := rec.CaptureToWriter(ctx, os.Stdout, "deploy", flightrecorder.CaptureReasonError)
```

## Configuration

| Field            | Type            | Default        | Description                                                     |
| ---------------- | --------------- | -------------- | --------------------------------------------------------------- |
| `MinAge`         | `time.Duration` | `5s`           | Minimum time trace data is retained in the in-memory buffer     |
| `MaxBytes`       | `uint64`        | `10 MiB`       | Maximum in-memory buffer size                                   |
| `CaptureOnSlow`  | `bool`          | `true`         | Auto-capture when a command exceeds `SlowThreshold`             |
| `SlowThreshold`  | `time.Duration` | `100ms`        | Duration above which a command triggers a slow snapshot         |
| `CaptureOnError` | `bool`          | `false`        | Auto-capture when a command returns an error                    |
| `OutputDir`      | `string`        | `""` (tmpdir)  | Directory for `.trace` snapshot files (created if missing)      |
| `FilenamePrefix` | `string`        | `"cmdguard"`   | Prepended to snapshot filenames                                 |
| `Log`            | `func(...)`     | `nil` (stderr) | Diagnostic log function (snapshot captured, start failed, etc.) |

## Constraints

- **Process-wide singleton**: at most one flight recorder may be active per process (`runtime/trace` limitation). If you create multiple `Recorder` instances, only the first `Start()` succeeds; subsequent starts return `ErrAlreadyStarted`.
- **Non-blocking**: the middleware never blocks command execution on recorder operations. Snapshot capture runs in a background goroutine.
- **Error precedence**: when both `CaptureOnSlow` and `CaptureOnError` are enabled and a command is both slow AND errors, the error reason takes precedence.
