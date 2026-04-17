package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLICommand[C any](use string) v2.Command[C, v2.NoFlags] {
	return v2.MustNewCommand[C, v2.NoFlags](use, noOpRunE[C])
}

// newTestCLICommandWithShort creates a leaf command with a short description.
func newTestCLICommandWithShort[C any](use, short string) v2.Command[C, v2.NoFlags] {
	return v2.MustNewCommand[C, v2.NoFlags](use, noOpRunE[C],
		v2.WithShort[C, v2.NoFlags](short),
	)
}

// newTestParentCommand creates a parent command with child subcommands.
func newTestParentCommand[C any](
	use, short, long string,
	children ...v2.Command[C, v2.NoFlags],
) v2.Command[C, v2.NoFlags] {
	return v2.MustNewParentCommand[C, v2.NoFlags](use, long, children,
		v2.WithShort[C, v2.NoFlags](short),
	)
}

func noOpRunE[C any](_ context.Context, _ *C, _ v2.NoFlags) error {
	return nil
}

// NoOpRunEWithFlags returns a no-op RunE function for commands with flags.
func NoOpRunEWithFlags[C, F any]() func(context.Context, *C, F) error {
	return func(_ context.Context, _ *C, _ F) error {
		return nil
	}
}

// RecordingHook returns a RunE function that records execution order.
func RecordingHook[C, F any](order *[]string, msg string) func(context.Context, *C, F) error {
	return func(_ context.Context, _ *C, _ F) error {
		*order = append(*order, msg)

		return nil
	}
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
