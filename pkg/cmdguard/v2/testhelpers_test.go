package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func newTestCLICommand[C any](t *testing.T, use string) v2.Command[C, v2.NoFlags] {
	t.Helper()

	cmd, err := v2.NewCommand[C, v2.NoFlags](use, noOpRunE[C])
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func newTestCLICommandWithShort[C any](t *testing.T, use, short string) v2.Command[C, v2.NoFlags] {
	t.Helper()

	cmd, err := v2.NewCommand[C, v2.NoFlags](
		use, noOpRunE[C],
		v2.WithShort[C, v2.NoFlags](short),
	)
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func newTestParentCommand[C any](
	t *testing.T,
	use, short, long string,
	children ...v2.Command[C, v2.NoFlags],
) v2.Command[C, v2.NoFlags] {
	t.Helper()

	cmd, err := v2.NewParentCommand[C, v2.NoFlags](
		use, long, children,
		v2.WithShort[C, v2.NoFlags](short),
	)
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func noOpRunE[C any](_ context.Context, _ *C, _ v2.NoFlags) error {
	return nil
}

func NoOpRunEWithFlags[C, F any]() func(context.Context, *C, F) error {
	return func(_ context.Context, _ *C, _ F) error {
		return nil
	}
}

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

func testHostPortPortInt(t *testing.T, hp v2.HostPort, expected int) {
	t.Helper()

	if hp.Port().Int() != expected {
		t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), expected)
	}
}
