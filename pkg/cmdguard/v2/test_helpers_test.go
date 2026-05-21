package v2

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func assertDurationField(t *testing.T, d Duration, expected time.Duration) {
	t.Helper()

	testutil.AssertEqual(t, d.Duration(), expected)
}

func assertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()
	testutil.AssertErrorContains(t, err, substrings...)
}

func assertStderrContains(t *testing.T, stderr string, substrings ...string) {
	t.Helper()
	testutil.AssertStderrContains(t, stderr, substrings...)
}

func assertStringContains(t *testing.T, s string, substrings ...string) {
	t.Helper()

	for _, sub := range substrings {
		if !strings.Contains(s, sub) {
			t.Errorf("expected %q to contain %q", s, sub)
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
		use: use,
		runE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
			return runErr
		},
	}
}

func registerAndSetFlag[T any](
	t *testing.T,
	registry *FlagRegistry,
	cmd *cobra.Command,
	cfg *T,
	flagName, flagValue string,
) {
	t.Helper()

	registerAndParseFlags(t, registry, cmd, flagName, flagValue)

	err := registry.ParseFlags(cmd, cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func registerAndParseFlags(
	t *testing.T,
	registry *FlagRegistry,
	cmd *cobra.Command,
	flagName, flagValue string,
) {
	t.Helper()

	err := registry.RegisterFlags(cmd)
	if err != nil {
		t.Fatalf("expected no error registering flags, got: %v", err)
	}

	err = cmd.Flags().Set(flagName, flagValue)
	if err != nil {
		t.Fatalf("expected no error setting flag, got: %v", err)
	}
}

func noOpRunE[T, F any](_ context.Context, _ *T, _ F) error {
	return nil
}

func noOpHandler() func(context.Context, *testConfig, NoFlags) error {
	return noOpRunE[testConfig, NoFlags]
}

func noOpHandlerForTestAppConfig() func(context.Context, *testAppConfig, NoFlags) error {
	return noOpRunE[testAppConfig, NoFlags]
}

// makeHookRunE creates a RunE function that records execution order.
func makeHookRunE(
	order *[]string,
	msg string,
) func(context.Context, *testAppConfig, NoFlags) error {
	return func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
		*order = append(*order, msg)

		return nil
	}
}

func assertEnumString(t *testing.T, got, want, fieldName string) {
	t.Helper()

	testutil.AssertFieldEqString(t, got, want, "unmarshaled "+fieldName)
}

func setFlag(t *testing.T, cmd *cobra.Command, name, value string) {
	t.Helper()

	err := cmd.Flags().Set(name, value)
	if err != nil {
		t.Fatalf("expected no error setting flag %q to %q, got: %v", name, value, err)
	}
}

func addCommand[T, F any](t *testing.T, cli *CLI[T], cmd Command[T, F]) {
	t.Helper()

	err := AddCommand(cli, cmd)
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}
}

func assertFlowValue(t *testing.T, bfc *BranchingFlowContext, key, expected any, msg string) {
	t.Helper()

	if bfc.Value(key) != expected {
		t.Error(msg)
	}
}

func mustProvideValue[T any](t *testing.T, scope *Scope, value T) {
	t.Helper()

	if err := ProvideValue(scope, value); err != nil {
		t.Fatalf("expected no error providing value, got: %v", err)
	}
}

func mustInvoke[T any](t *testing.T, scope *Scope) T {
	t.Helper()

	value, err := Invoke[T](scope)
	if err != nil {
		t.Fatalf("expected no error invoking value, got: %v", err)
	}

	return value
}

func mustInvokeNamed[T any](t *testing.T, scope *Scope, name string) T {
	t.Helper()

	value, err := InvokeNamed[T](scope, name)
	if err != nil {
		t.Fatalf("expected no error invoking named value %q, got: %v", name, err)
	}

	return value
}

func assertValidatorError(t *testing.T, name string, validator func(string) error, param string) {
	t.Helper()

	err := validator(param)
	if err == nil {
		t.Fatalf("expected error for %s with param %q, got nil", name, param)
	}

	assertErrorContains(t, err, "invalid")
}

func addShortCommandToStrictCLI(
	t *testing.T,
	opts ...CLIOption[testConfig],
) {
	t.Helper()

	cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, opts...)
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand[testConfig, NoFlags](
		"good",
		noOpHandler(),
		WithShort[testConfig, NoFlags]("A good command"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))
}
