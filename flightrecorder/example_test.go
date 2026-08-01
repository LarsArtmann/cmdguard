package flightrecorder_test

import (
	"fmt"
	"time"

	"github.com/larsartmann/cmdguard/flightrecorder"
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
