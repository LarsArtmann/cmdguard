package v3

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

// newTestCLIWithNoOpCmd creates a CLI with a no-op "cmd" command using the given options.
func newTestCLIWithNoOpCmd(t *testing.T, opts ...CommandOption) *CLI[testConfig] {
	t.Helper()

	cli, err := NewCLI("test", "Test", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"cmd",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		opts...,
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	return cli
}

func TestExitError(t *testing.T) {
	t.Parallel()

	t.Run("Error returns wrapped message", func(t *testing.T) {
		t.Parallel()

		err, _ := NewExitError(2, errors.New("something failed"))
		testutil.AssertFieldEqString(t, err.Error(), "something failed", "Error()")
	})

	t.Run("ExitCode returns code", func(t *testing.T) {
		t.Parallel()

		err, _ := NewExitError(42, errors.New("custom"))
		testutil.AssertFieldEq(t, err.ExitCode(), 42, "ExitCode()")
	})

	t.Run("Unwrap returns inner error", func(t *testing.T) {
		t.Parallel()

		inner := errors.New("inner")
		err, _ := NewExitError(1, inner)
		testutil.AssertEqual(t, err.Unwrap(), inner)
	})

	t.Run("implements ExitCoder", func(t *testing.T) {
		t.Parallel()

		exitErr, errCheck := NewExitError(3, errors.New("fail"))
		if errCheck != nil {
			t.Fatal(errCheck)
		}

		_, ok := any(exitErr).(ExitCoder)
		if !ok {
			t.Fatal("expected ExitError to implement ExitCoder")
		}
	})

	t.Run("errors.As detects ExitCoder", func(t *testing.T) {
		t.Parallel()

		err, _ := NewExitError(5, errors.New("wrapped"))
		wrapped := fmt.Errorf("outer: %w", err)

		exitCoder, ok := errors.AsType[ExitCoder](wrapped)
		if !ok {
			t.Fatal("expected errors.AsType to find ExitCoder")
		}

		testutil.AssertFieldEq(t, exitCoder.ExitCode(), 5, "ExitCode()")
	})
}

func TestExecuteAndExit_ExitCodes(t *testing.T) {
	t.Parallel()

	t.Run("default exit code is 1", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand(
			"fail",
			NoFlags{},
			func(_ context.Context, _ *testConfig, _ NoFlags) error {
				return errors.New("always fails")
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		// Can't test os.Exit directly, but test the ExitCoder logic
		err = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
		testutil.AssertExpectedError(t, err)

		if _, ok := errors.AsType[ExitCoder](err); ok {
			t.Error("plain error should not implement ExitCoder")
		}
	})

	t.Run("ExitError propagates exit code", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand(
			"fail-custom",
			NoFlags{},
			func(_ context.Context, _ *testConfig, _ NoFlags) error {
				exitErr, _ := NewExitError(42, errors.New("custom failure"))

				return exitErr
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"fail-custom"})
		testutil.AssertExpectedError(t, err)

		exitCoder, ok := errors.AsType[ExitCoder](err)
		if !ok {
			t.Fatal("expected error to implement ExitCoder")
		}

		testutil.AssertFieldEq(t, exitCoder.ExitCode(), 42, "ExitCode()")
	})
}

func TestWithConfigValidation(t *testing.T) {
	t.Parallel()

	t.Run("validation passes - command runs", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Name string `flag:"name" default:"world" help:"Name"`
		}

		var validated bool
		cli, err := NewCLI(
			"test", "Test", config{},
			WithConfigValidation(func(cfg *config) error {
				validated = true

				return nil
			}),
		)
		testutil.AssertNoError(t, err)

		executed := false
		runFlagCommand(t, cli, &executed)

		err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, validated, true)
		testutil.AssertEqual(t, executed, true)
	})

	t.Run("validation fails - command does not run", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Port int `flag:"port" default:"0" help:"Port"`
		}

		cli, err := NewCLI(
			"test", "Test", config{},
			WithConfigValidation(func(cfg *config) error {
				if cfg.Port < 1 {
					return fmt.Errorf("port must be > 0, got %d", cfg.Port)
				}

				return nil
			}),
		)
		testutil.AssertNoError(t, err)

		executed := false
		runFlagCommand(t, cli, &executed)

		err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "config validation failed", "port must be > 0")
		testutil.AssertEqual(t, executed, false)
	})
}

func TestWithStrictValidation(t *testing.T) {
	t.Parallel()

	t.Run("command without short fails in strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithStrictValidation(),
		)
		testutil.AssertNoError(t, err)

		cmd := noShortCommand(t)

		err = AddCommand(cli, cmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort, got %v", err)
		}
	})

	t.Run("command with short passes in strict mode", assertShortCommandAcceptedOnStrictCLI)

	t.Run("command without short passes without strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd := noShortCommand(t)
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("subcommand without short fails in strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithStrictValidation(),
		)
		testutil.AssertNoError(t, err)

		childCmd, err := NewCommand(
			"child-no-short",
			NoFlags{},
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		)
		testutil.AssertNoError(t, err)

		parentCmd, err := NewParentCommand[testConfig](
			"parent",
			"Parent description",
			NoFlags{},
			WithSubcommands(childCmd),
			WithShort("Parent"),
		)
		testutil.AssertNoError(t, err)

		err = AddCommand(cli, parentCmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort, got %v", err)
		}
	})
}

func TestWithDraconianValidation(t *testing.T) {
	t.Parallel()

	t.Run("leaf without example fails in draconian mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithDraconianValidation(),
		)
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand(
			"noexample",
			NoFlags{},
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithShort("Has short but no example"),
		)
		testutil.AssertNoError(t, err)

		err = AddCommand(cli, cmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingExample) {
			t.Errorf("expected ErrMissingExample, got %v", err)
		}
	})

	t.Run("leaf with example passes in draconian mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithDraconianValidation(),
		)
		testutil.AssertNoError(t, err)

		cmd := goodCommand(t, "good", "Good command", "test good")
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("parent command without example passes in draconian mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithDraconianValidation(),
		)
		testutil.AssertNoError(t, err)

		childCmd := goodCommand(t, "child", "Child", "test parent child")

		parentCmd, err := NewParentCommand[testConfig](
			"parent",
			"Parent description", NoFlags{},
			WithSubcommands(childCmd),
			WithShort("Parent"),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, parentCmd))
	})

	t.Run("draconian also enforces short description", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"test", "Test", testConfig{},
			WithDraconianValidation(),
		)
		testutil.AssertNoError(t, err)

		cmd := noShortCommand(t)

		err = AddCommand(cli, cmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort (draconian includes strict), got %v", err)
		}
	})

	t.Run("leaf without example passes in strict mode", assertShortCommandAcceptedOnStrictCLI)
}

func TestCommand_ValidateStrict(t *testing.T) {
	t.Parallel()

	t.Run("strict mode rejects missing short on leaf", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "test"},
			runE: noOpHandler(),
		}

		err := cmd.ValidateStrict()
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort, got %v", err)
		}
	})

	t.Run("non-strict allows missing short on leaf", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "test"},
			runE: noOpHandler(),
		}

		err := cmd.Validate()
		testutil.AssertNoError(t, err)
	})

	t.Run("strict passes with short", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "test", short: "Test command"},
			runE: noOpHandler(),
		}

		err := cmd.ValidateStrict()
		testutil.AssertNoError(t, err)
	})
}

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	t.Run("creates version command with version set", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"myapp", "Test", testConfig{},
			WithCLIVersion("1.2.3"),
		)
		testutil.AssertNoError(t, err)

		cmd, err := VersionCommand(cli)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cmd.Use(), "version", "Use()")
		testutil.AssertFieldEqString(t, cmd.Short(), "Print version information", "Short()")
	})

	t.Run("fails without version set", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("myapp", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		_, err = VersionCommand[testConfig](cli)
		testutil.AssertExpectedError(t, err)
	})

	t.Run("executes and prints version", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI(
			"myapp", "Test", testConfig{},
			WithCLIVersion("2.0.0"),
			WithFang(false),
		)
		testutil.AssertNoError(t, err)

		cmd, err := VersionCommand[testConfig](cli)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var buf strings.Builder
		cli.rootCmd.SetOut(&buf)

		err = cli.ExecuteWithArgs(context.Background(), []string{"version"})
		testutil.AssertNoError(t, err)

		output := buf.String()
		testutil.AssertOutputContains(t, output, "myapp 2.0.0")
	})

	t.Run("VersionCommand returns error without version", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("myapp", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		_, err = VersionCommand[testConfig](cli)
		if err == nil {
			t.Fatal("expected error when no version set")
		}
	})
}

func TestWithExactArgs(t *testing.T) {
	t.Parallel()

	t.Run("exact args - correct count passes", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("test", "Test", testConfig{}, WithFang(false))
		testutil.AssertNoError(t, err)

		var received []string
		cmd, err := NewCommand(
			"cmd",
			NoFlags{},
			func(ctx context.Context, _ *testConfig, _ NoFlags) error {
				received = ArgsFromContext(ctx)

				return nil
			},
			WithExactArgs(2),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "arg1", "arg2"})
		testutil.AssertNoError(t, err)
		if len(received) != 2 {
			t.Errorf("expected 2 args, got %d", len(received))
		}
	})

	t.Run("exact args - wrong count fails", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithExactArgs(2))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "only-one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts 2 arg(s)", "received 1")
	})
}

func TestWithMinimumArgs(t *testing.T) {
	t.Parallel()

	t.Run("minimum args - enough args passes", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithMinimumArgs(1))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "arg1", "arg2"})
		testutil.AssertNoError(t, err)
	})

	t.Run("minimum args - too few fails", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithMinimumArgs(2))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "requires at least 2 arg(s)")
	})
}

func TestWithMaximumArgs(t *testing.T) {
	t.Parallel()

	t.Run("maximum args - within limit passes", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithMaximumArgs(2))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd"})
		testutil.AssertNoError(t, err)
	})

	t.Run("maximum args - exceeds limit fails", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithMaximumArgs(1))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "a", "b"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts at most 1 arg(s)")
	})
}

func TestWithNoArgs(t *testing.T) {
	t.Parallel()

	t.Run("no args - zero args passes", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithNoArgs())

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd"})
		testutil.AssertNoError(t, err)
	})

	t.Run("no args - any args fails", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithNoArgs())

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "--", "unexpected"})
		testutil.AssertExpectedError(t, err)

		// Cobra treats unknown positional args as unknown commands before
		// the Args validator runs when the command has no subcommands.
		// Either "unknown command" or "accepts no args" is acceptable.
		errMsg := err.Error()
		if !strings.Contains(errMsg, "accepts no args") &&
			!strings.Contains(errMsg, "unknown command") {
			t.Errorf("error should contain 'accepts no args' or 'unknown command', got %q", errMsg)
		}
	})
}

func TestWithRangeArgs(t *testing.T) {
	t.Parallel()

	t.Run("within range passes", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithRangeArgs(1, 3))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "a", "b"})
		testutil.AssertNoError(t, err)
	})

	t.Run("below range fails", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLIWithNoOpCmd(t, WithRangeArgs(2, 4))

		err := cli.ExecuteWithArgs(context.Background(), []string{"cmd", "only-one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts between 2 and 4 arg(s)")
	})
}
