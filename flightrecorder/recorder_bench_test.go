package flightrecorder

import (
	"context"
	"runtime/trace"
	"testing"
	"time"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		_ = New(DefaultConfig())
	}
}

func BenchmarkMiddleware_Overhead(b *testing.B) {
	rec := New(Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Hour, // high threshold so no capture fires
		CaptureOnError: false,
	})

	if err := rec.Start(); err != nil {
		b.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	ctx := context.Background()
	cfg := &testConfig{}
	info := v4.CommandInfo{Name: "bench", Phase: v4.PhaseRun}
	middleware := Middleware[testConfig](rec)

	next := func() error { return nil }

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = middleware(ctx, cfg, info, next)
	}
}

func BenchmarkCapture(b *testing.B) {
	dir := b.TempDir()

	rec := New(Config{
		MinAge:    1 * time.Second,
		MaxBytes:  1 << 20,
		OutputDir: dir,
	})

	if err := rec.Start(); err != nil {
		b.Fatalf("Start failed: %v", err)
	}

	defer rec.Stop()

	trace.WithRegion(context.Background(), "benchmark-work", func() {
		time.Sleep(1 * time.Millisecond)
	})

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = rec.Capture(ctx, "bench", CaptureReasonSlow)
	}
}
