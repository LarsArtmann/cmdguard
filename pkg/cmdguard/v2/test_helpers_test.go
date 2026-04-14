package v2

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	testutil.AssertPanics(t, fn)
}

func assertDurationField(t *testing.T, d Duration, expected time.Duration) {
	t.Helper()

	if d.Duration() != expected {
		t.Errorf("expected %v, got %v", expected, d.Duration())
	}
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
		Use: use,
		RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
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

func noOpHandler() func(context.Context, *testConfig, NoFlags) error {
	return func(_ context.Context, _ *testConfig, _ NoFlags) error {
		return nil
	}
}

func noOpHandlerForTestAppConfig() func(context.Context, *testAppConfig, NoFlags) error {
	return func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
		return nil
	}
}

func assertEnumString(t *testing.T, got, want, fieldName string) {
	t.Helper()

	if got != want {
		t.Errorf("unmarshaled %s = %q, want %q", fieldName, got, want)
	}
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
