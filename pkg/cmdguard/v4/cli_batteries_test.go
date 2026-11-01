package v4

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestNewApp_RegistersCommands(t *testing.T) {
	t.Parallel()

	cli, err := NewApp(
		"app", "App", testConfig{},
		func(cli *CLI[testConfig]) error {
			cmd, cmdErr := NewCommand(
				"list",
				NoFlags{},
				func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
				WithShort("List things"),
			)
			if cmdErr != nil {
				return cmdErr
			}

			return AddCommand(cli, cmd)
		},
		WithFang(false),
	)
	testutil.AssertNoError(t, err)

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"list"}))
}

func TestNewApp_PropagatesBuildError(t *testing.T) {
	t.Parallel()

	buildErr := errors.New("build failed")

	cli, err := NewApp(
		"app", "App", testConfig{},
		func(_ *CLI[testConfig]) error { return buildErr },
		WithFang(false),
	)
	if cli != nil {
		t.Error("NewApp must discard a partially built CLI")
	}

	if !errors.Is(err, buildErr) {
		t.Fatalf("expected build error to be wrapped, got: %v", err)
	}
}

func TestNewApp_NilBuild(t *testing.T) {
	t.Parallel()

	cli, err := NewApp("app", "App", testConfig{}, nil, WithFang(false))
	testutil.AssertNoError(t, err)
	testutil.AssertNotNil(t, cli)
}

func TestWithBatteries_RegistersVersionAndDoctor(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithFang(false),
		WithCLIVersion("1.2.3"),
		WithBatteries(),
	)
	testutil.AssertNoError(t, err)

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"version"}))
	assertStringContains(t, out.String(), "app 1.2.3")

	names := commandNames(cli)
	for _, want := range []string{"version", "doctor", "describe"} {
		if !containsString(names, want) {
			t.Errorf("expected batteries to register %q, commands: %v", want, names)
		}
	}

	out.Reset()
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"doctor"}))
	assertStringContains(t, out.String(), "0 passed, 0 failed")
}

func TestWithBatteries_SkipsVersionWithoutVersion(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false), WithBatteries())
	testutil.AssertNoError(t, err)

	cli.RootCommand().SetOut(&bytes.Buffer{})

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"doctor"}))

	if containsString(commandNames(cli), "version") {
		t.Error("batteries must not register a version command when no version is set")
	}
}

func TestWithBatteries_EnablesStrictValidation(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false), WithBatteries())
	testutil.AssertNoError(t, err)

	bare, err := NewCommand(
		"bare",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, bare)
	if !errors.Is(err, ErrMissingShort) {
		t.Fatalf("expected ErrMissingShort under batteries, got: %v", err)
	}
}

func TestWithBatteries_ConsumerCommandsWin(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithFang(false),
		WithCLIVersion("1.2.3"),
		WithBatteries(),
	)
	testutil.AssertNoError(t, err)

	mine, err := NewCommand(
		"version",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error {
			_, writeErr := fmt.Fprintln(cli.RootCommand().OutOrStdout(), "mine")

			return writeErr
		},
		WithShort("Consumer version"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, mine))

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"version"}))
	testutil.AssertFieldEqString(t, out.String(), "mine\n", "consumer version output")
}

func TestWithBatteries_IdempotentAcrossExecutes(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithFang(false),
		WithCLIVersion("1.2.3"),
		WithBatteries(),
	)
	testutil.AssertNoError(t, err)

	cli.RootCommand().SetOut(&bytes.Buffer{})

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"doctor"}))
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"doctor"}))
}

// commandNames returns the names of the root command's direct subcommands.
func commandNames[T any](cli *CLI[T]) []string {
	cmds := cli.rootCmd.Commands()
	names := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		names = append(names, cmd.Name())
	}

	return names
}

// containsString reports whether values contains want.
func containsString(values []string, want string) bool {
	return slices.Contains(values, want)
}
