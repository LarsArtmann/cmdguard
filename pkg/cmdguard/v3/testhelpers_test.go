package v3_test

import (
	"context"
	"testing"

	auditlog "github.com/larsartmann/samber-do-auditlog"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func newTestCLICommand[C any](t *testing.T, use string) v3.Command[C, v3.NoFlags] {
	t.Helper()

	cmd, err := v3.NewCommand(use, v3.NoFlags{}, noOpRunE[C])
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func newTestCLICommandWithShort[C any](t *testing.T, use, short string) v3.Command[C, v3.NoFlags] {
	t.Helper()

	cmd, err := v3.NewCommand(
		use, v3.NoFlags{}, noOpRunE[C],
		v3.WithShort(short),
	)
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func newTestParentCommand[C any](
	t *testing.T,
	use, short, long string,
	children ...v3.Command[C, v3.NoFlags],
) v3.Command[C, v3.NoFlags] {
	t.Helper()

	cmd, err := v3.NewParentCommand[C](
		use, long, v3.NoFlags{},
		v3.WithSubcommands(children...),
		v3.WithShort(short),
	)
	if err != nil {
		t.Fatal(err)
	}

	return cmd
}

func noOpRunE[C any](_ context.Context, _ *C, _ v3.NoFlags) error {
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

func testHostPortPortInt(t *testing.T, hp v3.HostPort, expected int) {
	t.Helper()

	if hp.Port().Int() != expected {
		t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), expected)
	}
}

// addCommand registers a Command on a CLI and fails the test on error.
// Delegates to testutil.AddCommand to keep the canonical helper in one place.
func addCommand[T, F any](t *testing.T, cli *v3.CLI[T], cmd v3.Command[T, F]) {
	t.Helper()

	if err := v3.AddCommand(cli, cmd); err != nil {
		t.Fatalf("AddCommand: %v", err)
	}
}

// newTestCLI builds a CLI[testCLIConfig] with no options. Centralizes the
// trivial "test", "Test", testCLIConfig{} + t.Fatalf pattern used across
// cli_auditlog_test.go.
func newTestCLI(t *testing.T) *v3.CLI[testCLIConfig] {
	t.Helper()

	cli, err := v3.NewCLI[testCLIConfig]("test", "Test", testCLIConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	return cli
}

// newTestCLIWithAuditLog builds a CLI[testCLIConfig] with WithAuditLog(plugin).
// Centralizes the NewCLI + WithAuditLog + t.Fatalf pattern used across
// cli_auditlog_test.go.
func newTestCLIWithAuditLog(
	t *testing.T,
	plugin *auditlog.Plugin,
) *v3.CLI[testCLIConfig] {
	t.Helper()

	cli, err := v3.NewCLI[testCLIConfig](
		"test", "Test", testCLIConfig{},
		v3.WithAuditLog(plugin),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	return cli
}
