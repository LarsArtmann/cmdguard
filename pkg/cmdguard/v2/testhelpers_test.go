package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLICommand[C any](use string) v2.Command[C, v2.NoFlags] {
	return v2.Command[C, v2.NoFlags]{
		Use:  use,
		RunE: noOpRunE[C],
	}
}

func noOpRunE[C any](_ context.Context, _ *C, _ v2.NoFlags) error {
	return nil
}

func testParseError[T any](t *testing.T, parseFn func() (T, error), typeName string) {
	t.Helper()

	_, err := parseFn()
	if err == nil {
		t.Fatalf("expected error for %s", typeName)
	}
}

func testMustParsePanics[T any](t *testing.T, mustFn func(string) T, typeName string) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid %s", typeName)
		}
	}()

	_ = mustFn("invalid")
}

func testHostPortPortInt(t *testing.T, hp v2.HostPort, expected int) {
	t.Helper()

	if hp.Port().Int() != expected {
		t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), expected)
	}
}
