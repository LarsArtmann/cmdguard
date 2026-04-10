package v2

import (
	"context"
	"strings"
	"testing"
	"time"
)

func assertDurationField(t *testing.T, d Duration, expected time.Duration) {
	t.Helper()
	if d.Duration() != expected {
		t.Errorf("expected %v, got %v", expected, d.Duration())
	}
}

func assertStderrContains(t *testing.T, stderr string, substrings ...string) {
	t.Helper()
	for _, s := range substrings {
		if !strings.Contains(strings.ToLower(stderr), strings.ToLower(s)) {
			t.Errorf("stderr should contain %q, got %q", s, stderr)
		}
	}
}

type testAppConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Output  string `default:"-"     flag:"output"  help:"Output file"           short:"o"`
}

func newTestCmd(use string, err ...error) Command[testAppConfig, NoFlags] {
	var runErr error
	if len(err) > 0 {
		runErr = err[0]
	}

	return Command[testAppConfig, NoFlags]{
		Use: use,
		RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
			return runErr
		},
	}
}
