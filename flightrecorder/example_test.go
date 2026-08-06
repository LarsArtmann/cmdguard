package flightrecorder_test

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/larsartmann/cmdguard/flightrecorder"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func ExampleDefaultConfig() {
	cfg := flightrecorder.DefaultConfig()

	fmt.Println("capture on slow:", cfg.CaptureOnSlow)
	fmt.Println("slow threshold:", cfg.SlowThreshold)
	fmt.Println("capture on error:", cfg.CaptureOnError)
	fmt.Println("filename prefix:", cfg.FilenamePrefix)

	// Output:
	// capture on slow: true
	// slow threshold: 100ms
	// capture on error: false
	// filename prefix: cmdguard
}

func ExampleNew() {
	rec := flightrecorder.New(flightrecorder.Config{
		CaptureOnSlow:  true,
		SlowThreshold:  500 * time.Millisecond,
		CaptureOnError: true,
		OutputDir:      "/tmp/myapp-traces",
	})

	fmt.Println("enabled before start:", rec.Enabled())

	if err := rec.Start(); err != nil {
		fmt.Println("failed to start:", err)

		return
	}

	defer rec.Stop()

	fmt.Println("enabled after start:", rec.Enabled())

	// Output:
	// enabled before start: false
	// enabled after start: true
}

func ExampleCaptureReason() {
	fmt.Println(flightrecorder.CaptureReasonSlow)
	fmt.Println(flightrecorder.CaptureReasonError)

	// Output:
	// slow
	// error
}

func ExampleRecorder_CaptureToWriter() {
	rec := flightrecorder.New(flightrecorder.Config{
		CaptureOnSlow: true,
		SlowThreshold: 1 * time.Millisecond,
		OutputDir:     "", // use os.TempDir()
	})

	if err := rec.Start(); err != nil {
		fmt.Println("start error:", err)

		return
	}
	defer rec.Stop()

	written, err := rec.CaptureToWriter(
		context.Background(),
		io.Discard,
		"my-command",
		flightrecorder.CaptureReasonSlow,
	)
	if err != nil {
		fmt.Println("capture error:", err)

		return
	}

	fmt.Println("captured bytes:", written > 0)

	// Output:
	// captured bytes: true
}

func ExampleWithFlightRecorderRecorder() {
	rec := flightrecorder.New(flightrecorder.DefaultConfig())
	defer rec.Stop()

	type appConfig struct {
		Debug bool `flag:"debug" default:"false" help:"Enable debug mode"`
	}

	// WithFlightRecorderRecorder injects a pre-built Recorder into the CLI
	// middleware chain. Use this when you need direct access to the Recorder
	// (e.g. for programmatic CaptureToWriter calls or custom lifecycle control).
	cli, err := v4.NewCLI[appConfig](
		"myapp",
		"My Application",
		appConfig{},
		flightrecorder.WithFlightRecorderRecorder[appConfig](rec),
	)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("cli created:", cli.Name())

	// Output:
	// cli created: myapp
}
