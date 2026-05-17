package v2

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

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

		var exitCoder ExitCoder
		if !errors.As(wrapped, &exitCoder) {
			t.Fatal("expected errors.As to find ExitCoder")
		}

		testutil.AssertFieldEq(t, exitCoder.ExitCode(), 5, "ExitCode()")
	})
}

func TestExecuteAndExit_ExitCodes(t *testing.T) {
	t.Parallel()

	t.Run("default exit code is 1", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("fail",
			func(_ context.Context, _ *testConfig, _ NoFlags) error {
				return errors.New("always fails")
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		// Can't test os.Exit directly, but test the ExitCoder logic
		err = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
		testutil.AssertExpectedError(t, err)

		var exitCoder ExitCoder
		if errors.As(err, &exitCoder) {
			t.Error("plain error should not implement ExitCoder")
		}
	})

	t.Run("ExitError propagates exit code", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("fail-custom",
			func(_ context.Context, _ *testConfig, _ NoFlags) error {
				exitErr, _ := NewExitError(42, errors.New("custom failure"))

				return exitErr
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"fail-custom"})
		testutil.AssertExpectedError(t, err)

		var exitCoder ExitCoder
		if !errors.As(err, &exitCoder) {
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
		cli, err := NewCLI[config]("test", "Test", config{},
			WithConfigValidation[config](func(cfg *config) error {
				validated = true

				return nil
			}),
		)
		testutil.AssertNoError(t, err)

		executed := false
		cmd, err := NewCommand[config, NoFlags]("run",
			func(_ context.Context, _ *config, _ NoFlags) error {
				executed = true

				return nil
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

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

		cli, err := NewCLI[config]("test", "Test", config{},
			WithConfigValidation[config](func(cfg *config) error {
				if cfg.Port < 1 {
					return fmt.Errorf("port must be > 0, got %d", cfg.Port)
				}

				return nil
			}),
		)
		testutil.AssertNoError(t, err)

		executed := false
		cmd, err := NewCommand[config, NoFlags]("run",
			func(_ context.Context, _ *config, _ NoFlags) error {
				executed = true

				return nil
			},
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

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

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{},
			WithStrictValidation[testConfig](),
		)
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("noshort",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		)
		testutil.AssertNoError(t, err)

		err = AddCommand(cli, cmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort, got %v", err)
		}
	})

	t.Run("command with short passes in strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{},
			WithStrictValidation[testConfig](),
		)
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("good",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithShort[testConfig, NoFlags]("A good command"),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("command without short passes without strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("noshort",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("subcommand without short fails in strict mode", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{},
			WithStrictValidation[testConfig](),
		)
		testutil.AssertNoError(t, err)

		childCmd, err := NewCommand[testConfig, NoFlags]("child-no-short",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		)
		testutil.AssertNoError(t, err)

		parentCmd, err := NewParentCommand[testConfig, NoFlags]("parent",
			"Parent description",
			[]Command[testConfig, NoFlags]{childCmd},
			WithShort[testConfig, NoFlags]("Parent"),
		)
		testutil.AssertNoError(t, err)

		err = AddCommand(cli, parentCmd)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrMissingShort) {
			t.Errorf("expected ErrMissingShort, got %v", err)
		}
	})
}

func TestCommand_ValidateStrict(t *testing.T) {
	t.Parallel()

	t.Run("strict mode rejects missing short on leaf", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			use:  "test",
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
			use:  "test",
			runE: noOpHandler(),
		}

		err := cmd.Validate()
		testutil.AssertNoError(t, err)
	})

	t.Run("strict passes with short", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			use:   "test",
			short: "Test command",
			runE:  noOpHandler(),
		}

		err := cmd.ValidateStrict()
		testutil.AssertNoError(t, err)
	})
}

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	t.Run("creates version command with version set", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("myapp", "Test", testConfig{},
			WithCLIVersion[testConfig]("1.2.3"),
		)
		testutil.AssertNoError(t, err)

		cmd, err := VersionCommand[testConfig](cli)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cmd.Use(), "version", "Use()")
		testutil.AssertFieldEqString(t, cmd.Short(), "Print version information", "Short()")
	})

	t.Run("fails without version set", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("myapp", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		_, err = VersionCommand[testConfig](cli)
		testutil.AssertExpectedError(t, err)
	})

	t.Run("executes and prints version", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("myapp", "Test", testConfig{},
			WithCLIVersion[testConfig]("2.0.0"),
			WithFang[testConfig](false),
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
		if !strings.Contains(output, "myapp 2.0.0") {
			t.Errorf("expected output to contain 'myapp 2.0.0', got %q", output)
		}
	})

	t.Run("MustVersionCommand panics without version", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("myapp", "Test", testConfig{})
		testutil.AssertNoError(t, err)

		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic")
			}
		}()

		MustVersionCommand[testConfig](cli)
	})
}

func TestWithExactArgs(t *testing.T) {
	t.Parallel()

	t.Run("exact args - correct count passes", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		var received []string
		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(ctx context.Context, _ *testConfig, _ NoFlags) error {
				received = ArgsFromContext(ctx)

				return nil
			},
			WithExactArgs[testConfig, NoFlags](2),
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

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithExactArgs[testConfig, NoFlags](2),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "only-one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts 2 arg(s)", "received 1")
	})
}

func TestWithMinimumArgs(t *testing.T) {
	t.Parallel()

	t.Run("minimum args - enough args passes", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithMinimumArgs[testConfig, NoFlags](1),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "arg1", "arg2"})
		testutil.AssertNoError(t, err)
	})

	t.Run("minimum args - too few fails", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithMinimumArgs[testConfig, NoFlags](2),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "requires at least 2 arg(s)")
	})
}

func TestWithMaximumArgs(t *testing.T) {
	t.Parallel()

	t.Run("maximum args - within limit passes", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithMaximumArgs[testConfig, NoFlags](2),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd"})
		testutil.AssertNoError(t, err)
	})

	t.Run("maximum args - exceeds limit fails", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithMaximumArgs[testConfig, NoFlags](1),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "a", "b"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts at most 1 arg(s)")
	})
}

func TestWithNoArgs(t *testing.T) {
	t.Parallel()

	t.Run("no args - zero args passes", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithNoArgs[testConfig, NoFlags](),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd"})
		testutil.AssertNoError(t, err)
	})

	t.Run("no args - any args fails", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithNoArgs[testConfig, NoFlags](),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "--", "unexpected"})
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

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithRangeArgs[testConfig, NoFlags](1, 3),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "a", "b"})
		testutil.AssertNoError(t, err)
	})

	t.Run("below range fails", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testConfig]("test", "Test", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand[testConfig, NoFlags]("cmd",
			func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
			WithRangeArgs[testConfig, NoFlags](2, 4),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		err = cli.ExecuteWithArgs(context.Background(), []string{"cmd", "only-one"})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "accepts between 2 and 4 arg(s)")
	})
}
