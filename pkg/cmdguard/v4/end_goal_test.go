package v4

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// TestEndGoal_MinimalSuperbCLI is the acceptance test for the project's end
// goal: a developer who writes only the domain handler still gets a superb
// CLI, and every mistake is caught with an actionable message.
func TestEndGoal_MinimalSuperbCLI(t *testing.T) {
	t.Parallel()

	type appConfig struct {
		Token string `flag:"token" required:"true" help:"API token" env:"ENDGOAL_TOKEN"`
	}

	type runFlags struct {
		Name string `flag:"name" short:"n" required:"true" help:"Name" validate:"min=3"`
	}

	cli, err := NewApp(
		"app", "A superb CLI", appConfig{},
		func(cli *CLI[appConfig]) error {
			cmd, cmdErr := NewCommand(
				"run",
				runFlags{},
				func(_ context.Context, _ *appConfig, _ runFlags) error { return nil },
				WithShort("Run the thing"),
				WithExample("app run --token x --name abc"),
			)
			if cmdErr != nil {
				return cmdErr
			}

			return AddCommands(cli, cmd)
		},
		WithFang(false),
		WithCLIVersion("9.9.9"),
		WithBatteries(),
	)
	testutil.AssertNoError(t, err)

	var out, errOut bytes.Buffer

	cli.RootCommand().SetOut(&out)
	cli.RootCommand().SetErr(&errOut)

	// Discovery works before the caller knows any business flags.
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))

	var desc CLIDescription
	testutil.AssertNoError(t, json.Unmarshal(out.Bytes(), &desc))

	testutil.AssertFieldEqString(t, desc.SchemaVersion, "1", "schema version")

	if !strings.Contains(strings.Join(commandNames(cli), ","), "run") {
		t.Errorf("describe must document the run command, got: %v", commandNames(cli))
	}

	out.Reset()
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"version"}))
	assertStringContains(t, out.String(), "app 9.9.9")

	out.Reset()
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"doctor"}))
	assertStringContains(t, out.String(), "0 passed, 0 failed")

	// Missing root required flag fails with a hint.
	errOut.Reset()

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	if !errors.Is(err, ErrRequiredFlag) {
		t.Fatalf("expected ErrRequiredFlag for the missing token, got: %v", err)
	}

	assertStringContains(t, errOut.String(), "Hint: "+ErrorHint(ErrRequiredFlag))

	// Missing command required flag fails with a hint.
	errOut.Reset()

	err = cli.ExecuteWithArgs(context.Background(), []string{"run", "--token", "x"})
	if !errors.Is(err, ErrRequiredFlag) {
		t.Fatalf("expected ErrRequiredFlag for the missing name, got: %v", err)
	}

	assertStringContains(t, errOut.String(), `"name"`)

	// A too-short string respects the min-length rule.
	errOut.Reset()

	err = cli.ExecuteWithArgs(context.Background(), []string{"run", "--token", "x", "--name", "ab"})
	if !errors.Is(err, ErrValueTooShort) {
		t.Fatalf("expected ErrValueTooShort for a 2-character name, got: %v", err)
	}

	assertStringContains(t, errOut.String(), "Hint: "+ErrorHint(ErrValueTooShort))

	// Valid input runs cleanly.
	testutil.AssertNoError(t, cli.ExecuteWithArgs(
		context.Background(),
		[]string{"run", "--token", "x", "--name", "abc"},
	))
}
