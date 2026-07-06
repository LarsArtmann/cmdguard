package v2

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
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

// containsAll reports whether s contains every one of the given substrings.
// It is the conjunction of strings.Contains across all substrings, written as
// a function so callers can use it in a single if-condition.
func containsAll(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if !strings.Contains(s, sub) {
			return false
		}
	}

	return true
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
		spec: commandSpec{use: use},
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

// registerFlags registers all flags in registry on cmd and fails the test on error.
// Centralizes the RegisterFlags + t.Fatalf pattern for tests that don't set a flag value.
func registerFlags(t *testing.T, registry *FlagRegistry, cmd *cobra.Command) {
	t.Helper()

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("expected no error registering flags, got: %v", err)
	}
}

func noOpHandler() func(context.Context, *testConfig, NoFlags) error {
	return testutil.NoOpRunE[testConfig, NoFlags]
}

func noOpHandlerForTestAppConfig() func(context.Context, *testAppConfig, NoFlags) error {
	return testutil.NoOpRunE[testAppConfig, NoFlags]
}

// noOpRunE wraps testutil.NoOpRunE as a regular function for use in struct fields
// where generic type parameters would be unwieldy.
func noOpRunE[T, F any](ctx context.Context, cfg *T, flags F) error {
	return testutil.NoOpRunE[T, F](ctx, cfg, flags)
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

// addGroupedCommand builds a no-op Command[T, NoFlags] with use/short/long/group
// and registers it on the CLI. Used by command-grouping tests to reduce struct-literal boilerplate.
func addGroupedCommand[T any](t *testing.T, cli *CLI[T], use, short, group string) {
	t.Helper()

	err := AddCommand(cli, Command[T, NoFlags]{
		spec: commandSpec{use: use, short: short, long: short, group: group},
		runE: testutil.NoOpRunE[T, NoFlags],
	})
	if err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}
}

// recordHandlerCall returns a `next` callback that sets *called to true and returns nil.
// Used by middleware tests to assert the handler was invoked.
func recordHandlerCall(called *bool) func() error {
	return func() error {
		*called = true

		return nil
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
	opts ...CLIOption,
) {
	t.Helper()

	cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, opts...)
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"good",
		NoFlags{},
		noOpHandler(),
		WithShort("A good command"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))
}

// assertShortCommandAcceptedOnStrictCLI is a t.Run body that asserts a command
// with a short description is accepted under strict validation. It exists so
// tests that only need this assertion can use it without rewriting the
// t.Run/t.Parallel/helper boilerplate per subtest.
func assertShortCommandAcceptedOnStrictCLI(t *testing.T) {
	t.Parallel()
	addShortCommandToStrictCLI(t, WithStrictValidation())
}

// noShortCommand builds a minimal test command that omits a short description.
// Used to exercise strict/draconian validation rules that require WithShort.
func noShortCommand(t *testing.T) Command[testConfig, NoFlags] {
	t.Helper()

	cmd, err := NewCommand(
		"noshort",
		NoFlags{},
		noOpHandler(),
	)
	testutil.AssertNoError(t, err)

	return cmd
}

// goodCommand builds a test command that is fully described (short + example).
// Used to verify that strict and draconian validation accept a complete command.
func goodCommand(t *testing.T, use, short, example string) Command[testConfig, NoFlags] {
	t.Helper()

	cmd, err := NewCommand(
		use,
		NoFlags{},
		noOpHandler(),
		WithShort(short),
		WithExample(example),
	)
	testutil.AssertNoError(t, err)

	return cmd
}

// runFlagCommand adds a "run" command that flips *executed when invoked.
// Used by config-validation tests to assert whether the handler ran.
func runFlagCommand[T any](t *testing.T, cli *CLI[T], executed *bool) {
	t.Helper()

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, _ *T, _ NoFlags) error {
			*executed = true

			return nil
		},
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))
}
