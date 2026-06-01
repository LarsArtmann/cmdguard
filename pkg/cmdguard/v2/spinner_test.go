package v2

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestSpinnerMiddleware_SkipsNonTerminal(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	called := false

	mw := SpinnerMiddleware[testConfig]("Loading...")
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "test"},
		func() error {
			called = true

			return nil
		},
	)

	testutil.AssertNoError(t, err)

	if !called {
		t.Error("handler should be called when spinner skips")
	}
}

func TestSpinnerMiddleware_ReturnsError(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	expectedErr := errors.New("handler failed")

	mw := SpinnerMiddleware[testConfig]("Loading...")
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "test"},
		func() error {
			return expectedErr
		},
	)

	testutil.AssertErrorIs(t, err, expectedErr)
}

func TestSpinnerMiddleware_WritesToBuffer(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	buf := &bytes.Buffer{}
	cfg := SpinnerConfig{
		Title:    "Testing",
		Writer:   buf,
		Frames:   []string{".", "o", "O"},
		Interval: 10 * time.Millisecond,
	}

	s := newTextSpinner(cfg)
	s.Start()
	time.Sleep(25 * time.Millisecond)
	s.Stop()

	output := buf.String()

	if output == "" {
		t.Error("spinner should write to buffer")
	}

	if !strings.Contains(output, "Testing") {
		t.Errorf("output should contain title %q, got %q", "Testing", output)
	}

	frameCount := 0

	for _, frame := range cfg.Frames {
		if strings.Contains(output, frame) {
			frameCount++
		}
	}

	if frameCount == 0 {
		t.Errorf("output should contain at least one frame from %v, got %q", cfg.Frames, output)
	}

	if !strings.Contains(output, "\r") {
		t.Error("output should contain carriage return for line clearing")
	}

	clearSeq := "\r\033[K"
	if !strings.Contains(output, clearSeq) {
		t.Error("output should contain ANSI clear sequence from Stop()")
	}
}

func TestSpinnerMiddlewareWithConfig_SkipsNonTerminal(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	buf := &bytes.Buffer{}
	called := false

	cfg := SpinnerConfig{
		Title:    "Custom",
		Writer:   buf,
		Frames:   []string{"1", "2", "3"},
		Interval: 10 * time.Millisecond,
	}

	mw := SpinnerMiddlewareWithConfig[testConfig](cfg)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "test"},
		func() error {
			called = true

			return nil
		},
	)

	testutil.AssertNoError(t, err)

	if !called {
		t.Error("handler should be called when spinner skips non-terminal")
	}

	if buf.Len() > 0 {
		t.Error("nothing should be written to non-terminal buffer")
	}
}

func TestSpinnerMiddlewareWithConfig_ReturnsError(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	buf := &bytes.Buffer{}
	expectedErr := errors.New("handler failed")

	cfg := SpinnerConfig{
		Title:    "Custom",
		Writer:   buf,
		Frames:   []string{">"},
		Interval: 10 * time.Millisecond,
	}

	mw := SpinnerMiddlewareWithConfig[testConfig](cfg)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "test"},
		func() error { return expectedErr },
	)

	testutil.AssertErrorIs(t, err, expectedErr)
}

func TestDefaultSpinnerConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultSpinnerConfig("Test")

	if cfg.Title != "Test" {
		t.Errorf("expected title %q, got %q", "Test", cfg.Title)
	}

	if len(cfg.Frames) == 0 {
		t.Error("default frames should not be empty")
	}

	if cfg.Interval == 0 {
		t.Error("default interval should not be zero")
	}

	testutil.AssertNoError(t, cfg.Validate())
}

func TestSpinnerConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     SpinnerConfig
		wantErr error
	}{
		{
			name:    "nil writer",
			cfg:     SpinnerConfig{Writer: nil, Frames: defaultSpinnerFrames, Interval: 100 * time.Millisecond},
			wantErr: ErrNilValue,
		},
		{
			name:    "zero interval",
			cfg:     SpinnerConfig{Writer: os.Stderr, Frames: defaultSpinnerFrames, Interval: 0},
			wantErr: ErrValueTooSmall,
		},
		{
			name:    "negative interval",
			cfg:     SpinnerConfig{Writer: os.Stderr, Frames: defaultSpinnerFrames, Interval: -1 * time.Millisecond},
			wantErr: ErrValueTooSmall,
		},
		{
			name:    "empty frames",
			cfg:     SpinnerConfig{Writer: os.Stderr, Frames: []string{}, Interval: 100 * time.Millisecond},
			wantErr: ErrValueEmpty,
		},
		{
			name:    "valid config",
			cfg:     SpinnerConfig{Writer: os.Stderr, Frames: defaultSpinnerFrames, Interval: 100 * time.Millisecond},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErr == nil {
				testutil.AssertNoError(t, err)
			} else {
				testutil.AssertErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestSpinnerMiddlewareWithConfig_SkipsInvalidConfig(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	called := false

	cfg := SpinnerConfig{Writer: nil, Frames: nil, Interval: 0}
	mw := SpinnerMiddlewareWithConfig[testConfig](cfg)
	err := mw(
		context.Background(),
		&testConfig{},
		CommandInfo{Name: "test"},
		func() error {
			called = true

			return nil
		},
	)

	testutil.AssertNoError(t, err)

	if !called {
		t.Error("handler should be called when spinner skips invalid config")
	}
}
