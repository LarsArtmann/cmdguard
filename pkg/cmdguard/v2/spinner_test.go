package v2

import (
	"bytes"
	"context"
	"errors"
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

	if buf.Len() == 0 {
		t.Error("spinner should write to buffer")
	}
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
}
