package integration

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

const cmdRun = "run"

type lifecycleConfig struct {
	Verbose bool `flag:"verbose" short:"v" default:"false" help:"Verbose"`
}

func newLifecycleCmd(t *testing.T, use, short string) v2.Command[lifecycleConfig, v2.NoFlags] {
	t.Helper()

	cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
		use,
		func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
			return nil
		},
		v2.WithShort[lifecycleConfig, v2.NoFlags](short),
	)
	if err != nil {
		t.Fatalf("NewCommand %s: %v", use, err)
	}

	return cmd
}

func newLifecycleParentCmd(
	t *testing.T,
	child v2.Command[lifecycleConfig, v2.NoFlags],
	short string,
) v2.Command[lifecycleConfig, v2.NoFlags] {
	t.Helper()

	parent, err := v2.NewParentCommand[lifecycleConfig, v2.NoFlags](
		"parent", "Parent description",
		[]v2.Command[lifecycleConfig, v2.NoFlags]{child},
		v2.WithShort[lifecycleConfig, v2.NoFlags](short),
	)
	if err != nil {
		t.Fatalf("NewParentCommand: %v", err)
	}

	return parent
}

// newLifecyclePostRunFlag returns a WithPostRunE option that sets *flag to true.
// Used to verify post-run is (or is not) called from lifecycle tests.
func newLifecyclePostRunFlag(flag *bool) v2.CommandOption[lifecycleConfig, v2.NoFlags] {
	return v2.WithPostRunE[lifecycleConfig, v2.NoFlags](
		func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
			*flag = true

			return nil
		},
	)
}

// recordLifecycleStep returns a RunE/PreRunE/PostRunE handler that appends step
// to *order, then returns nil. Used by lifecycle tests to keep handler bodies terse.
func recordLifecycleStep(order *[]string, step string) func(
	_ context.Context, _ *lifecycleConfig, _ v2.NoFlags,
) error {
	return func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
		*order = append(*order, step)

		return nil
	}
}

func TestCLI_Lifecycle_PreRunAndPostRun(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a CLI with pre-run and post-run hooks, When command executes successfully, "+
			"Then all hooks and handler run in order",
		func(t *testing.T) {
			t.Parallel()

			var order []string

			cli, err := v2.NewCLI[lifecycleConfig](
				"lifecycle", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				cmdRun,
				recordLifecycleStep(&order, cmdRun),
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Run"),
				v2.WithPreRunE[lifecycleConfig, v2.NoFlags](
					recordLifecycleStep(&order, "pre-run"),
				),
				v2.WithPostRunE[lifecycleConfig, v2.NoFlags](
					recordLifecycleStep(&order, "post-run"),
				),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			if err := cli.ExecuteWithArgs(context.Background(), []string{cmdRun}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			expected := []string{"pre-run", cmdRun, "post-run"}
			if len(order) != len(expected) {
				t.Fatalf("execution order: got %v, want %v", order, expected)
			}

			for i, step := range expected {
				if order[i] != step {
					t.Errorf("order[%d]: got %q, want %q", i, order[i], step)
				}
			}
		},
	)

	t.Run(
		"Given a CLI with post-run hook, When command fails, Then post-run is NOT called",
		func(t *testing.T) {
			t.Parallel()

			var postRunCalled bool

			cli, err := v2.NewCLI[lifecycleConfig](
				"lifecycle", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"fail",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					return errors.New("command failed")
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Fail"),
				newLifecyclePostRunFlag(&postRunCalled),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
			if err == nil {
				t.Fatal("expected error")
			}

			if postRunCalled {
				t.Error("post-run should NOT be called when RunE fails")
			}
		},
	)

	t.Run(
		"Given a CLI with pre-run hook, When pre-run fails, "+
			"Then handler and post-run are NOT called",
		func(t *testing.T) {
			t.Parallel()

			var handlerCalled, postRunCalled bool

			cli, err := v2.NewCLI[lifecycleConfig](
				"lifecycle", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"prefail",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					handlerCalled = true

					return nil
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Pre-fail"),
				v2.WithPreRunE[lifecycleConfig, v2.NoFlags](
					func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
						return errors.New("pre-run rejected")
					},
				),
				newLifecyclePostRunFlag(&postRunCalled),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{"prefail"})
			if err == nil {
				t.Fatal("expected error from pre-run")
			}

			if handlerCalled {
				t.Error("handler should NOT be called when pre-run fails")
			}

			if postRunCalled {
				t.Error("post-run should NOT be called when pre-run fails")
			}
		},
	)
}

func TestCLI_Middleware_Chain(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a CLI with multiple middleware, When command runs, "+
			"Then middleware executes in order wrapping handler",
		func(t *testing.T) {
			t.Parallel()

			var order []string

			trackingMW := func(name string) v2.Middleware[lifecycleConfig] {
				return func(
					_ context.Context, _ *lifecycleConfig, _ v2.CommandInfo, next func() error,
				) error {
					order = append(order, name+"-before")
					err := next()
					order = append(order, name+"-after")

					return err
				}
			}

			cli, err := v2.NewCLI[lifecycleConfig](
				"mw", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithMiddleware[lifecycleConfig](trackingMW("mw1"), trackingMW("mw2")),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"run",
				recordLifecycleStep(&order, "handler"),
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Run"),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			if err := cli.ExecuteWithArgs(context.Background(), []string{"run"}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
			if len(order) != len(expected) {
				t.Fatalf("middleware order: got %v, want %v", order, expected)
			}

			for i, step := range expected {
				if order[i] != step {
					t.Errorf("order[%d]: got %q, want %q", i, order[i], step)
				}
			}
		},
	)

	t.Run(
		"Given a CLI with recovery middleware, When handler panics, "+
			"Then panic is caught and converted to error",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"recover", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithMiddleware[lifecycleConfig](v2.RecoveryMiddleware[lifecycleConfig]()),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"panic",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					panic("something went wrong")
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Panic"),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{"panic"})
			if err == nil {
				t.Fatal("expected error from panic recovery")
			}

			if !errors.Is(err, v2.ErrCommandPanic) {
				t.Errorf("expected ErrCommandPanic, got: %v", err)
			}

			if !strings.Contains(err.Error(), "something went wrong") {
				t.Errorf("error should contain panic message, got: %v", err)
			}
		},
	)

	t.Run(
		"Given a CLI with timing middleware, When command runs, Then timing callback is invoked",
		func(t *testing.T) {
			t.Parallel()

			var timedCommand string
			var gotDuration bool

			timingMW := v2.TimingMiddleware[lifecycleConfig](
				func(name string, d time.Duration, err error) {
					timedCommand = name
					gotDuration = d > 0
				},
			)

			cli, err := v2.NewCLI[lifecycleConfig](
				"timing", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithMiddleware[lifecycleConfig](timingMW),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd := newLifecycleCmd(t, "timed", "Timed")

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			if err := cli.ExecuteWithArgs(context.Background(), []string{"timed"}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			if timedCommand != "timed" {
				t.Errorf("timedCommand = %q, want %q", timedCommand, "timed")
			}

			if !gotDuration {
				t.Error("timing callback was not invoked")
			}
		},
	)
}

func TestCLI_DependencyInjection_Scope(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a CLI with DI services, When handler invokes a service, "+
			"Then the service is resolved correctly",
		func(t *testing.T) {
			t.Parallel()

			type Database struct {
				DSN string
			}

			cli, err := v2.NewCLI[lifecycleConfig](
				"di", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			err = v2.Provide(cli.Scope(), func(_ do.Injector) (*Database, error) {
				return &Database{DSN: "postgres://localhost:5432"}, nil
			})
			if err != nil {
				t.Fatalf("Provide: %v", err)
			}

			var resolvedDSN string

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"query",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					db, err := v2.Invoke[*Database](cli.Scope())
					if err != nil {
						return err
					}

					resolvedDSN = db.DSN

					return nil
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Query"),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			if err := cli.ExecuteWithArgs(context.Background(), []string{"query"}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			if resolvedDSN != "postgres://localhost:5432" {
				t.Errorf("resolvedDSN = %q, want %q", resolvedDSN, "postgres://localhost:5432")
			}
		},
	)

	t.Run(
		"Given a CLI with DI child scope, When service is provided in child, "+
			"Then child resolves its own service",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"di-child", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			childScope := cli.Scope().Child("plugin")

			err = v2.ProvideValue(childScope, "child-service-value")
			if err != nil {
				t.Fatalf("ProvideValue: %v", err)
			}

			val, err := v2.Invoke[string](childScope)
			if err != nil {
				t.Fatalf("Invoke from child: %v", err)
			}

			if val != "child-service-value" {
				t.Errorf("val = %q, want %q", val, "child-service-value")
			}
		},
	)

	t.Run(
		"Given a CLI scope, When providing nil scope, Then error is returned",
		func(t *testing.T) {
			t.Parallel()

			err := v2.Provide[string](nil, func(_ do.Injector) (string, error) {
				return "", nil
			})
			if err == nil {
				t.Fatal("expected error for nil scope")
			}

			if !errors.Is(err, v2.ErrInvalidScope) {
				t.Errorf("expected ErrInvalidScope, got: %v", err)
			}
		},
	)
}

func TestCLI_ErrorChains(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a command that returns ExitError, When executed, "+
			"Then errors.As finds ExitCoder",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"exit", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"die",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					exitErr, _ := v2.NewExitError(42, errors.New("permission denied"))

					return exitErr
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Die"),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{"die"})
			if err == nil {
				t.Fatal("expected error")
			}

			exitCoder, ok := errors.AsType[v2.ExitCoder](err)
			if !ok {
				t.Fatal("expected errors.AsType to find ExitCoder")
			}

			if exitCoder.ExitCode() != 42 {
				t.Errorf("exit code = %d, want 42", exitCoder.ExitCode())
			}
		},
	)

	t.Run(
		"Given a command error, When unwrapped, Then inner error is accessible",
		func(t *testing.T) {
			t.Parallel()

			inner := errors.New("root cause")
			cmdErr := v2.NewCommandError("mycommand", inner)

			if cmdErr.Error() == "" {
				t.Error("CommandError.Error() should not be empty")
			}

			if !strings.Contains(cmdErr.Error(), "mycommand") {
				t.Error("CommandError should contain command name")
			}

			if !strings.Contains(cmdErr.Error(), "root cause") {
				t.Error("CommandError should contain inner error message")
			}

			unwrapped := cmdErr.Unwrap()
			if unwrapped != inner {
				t.Error("Unwrap should return inner error")
			}
		},
	)

	t.Run(
		"Given a flag error with suggestion, When Error() is called, "+
			"Then suggestion appears in message",
		func(t *testing.T) {
			t.Parallel()

			flagErr := v2.NewFlagErrorWithSuggestion("nme", errors.New("unknown flag"), "name")
			msg := flagErr.Error()

			if !strings.Contains(msg, "nme") {
				t.Error("should contain flag name")
			}

			if !strings.Contains(msg, "name") {
				t.Error("should contain suggestion")
			}

			if !strings.Contains(msg, "did you mean") {
				t.Error("should contain suggestion hint")
			}
		},
	)
}

func TestCLI_ConfigValidation_Integration(t *testing.T) {
	t.Parallel()

	type serverConfig struct {
		Name string `flag:"name" default:"" help:"Server name"`
	}

	newValidatedServerCLI := func(t *testing.T) *v2.CLI[serverConfig] {
		t.Helper()

		cli, err := v2.NewCLI[serverConfig](
			"server", "Test", serverConfig{},
			v2.WithFang[serverConfig](false),
			v2.WithConfigValidation[serverConfig](func(cfg *serverConfig) error {
				if cfg.Name == "" {
					return errors.New("name is required")
				}

				return nil
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI: %v", err)
		}

		return cli
	}

	addStartCmd := func(t *testing.T, cli *v2.CLI[serverConfig], handler func() error) {
		t.Helper()

		cmd, err := v2.NewCommand[serverConfig, v2.NoFlags](
			"start",
			func(_ context.Context, _ *serverConfig, _ v2.NoFlags) error {
				return handler()
			},
			v2.WithShort[serverConfig, v2.NoFlags]("Start"),
		)
		if err != nil {
			t.Fatalf("NewCommand: %v", err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatalf("AddCommand: %v", err)
		}
	}

	t.Run(
		"Given CLI with config validation requiring non-empty name, "+
			"When name is empty (default), Then execution is rejected",
		func(t *testing.T) {
			t.Parallel()

			cli := newValidatedServerCLI(t)
			addStartCmd(t, cli, func() error { return nil })

			err := cli.ExecuteWithArgs(context.Background(), []string{"start"})
			if err == nil {
				t.Fatal("expected validation error")
			}

			if !strings.Contains(err.Error(), "config validation failed") {
				t.Errorf("error should mention config validation, got: %v", err)
			}

			if !strings.Contains(err.Error(), "name is required") {
				t.Errorf("error should contain validation message, got: %v", err)
			}
		},
	)

	t.Run(
		"Given CLI with config validation, When valid config is provided via flags, "+
			"Then execution succeeds",
		func(t *testing.T) {
			t.Parallel()

			var handlerCalled bool

			cli := newValidatedServerCLI(t)
			addStartCmd(t, cli, func() error {
				handlerCalled = true

				return nil
			})

			err := cli.ExecuteWithArgs(context.Background(), []string{"start", "--name=production"})
			if err != nil {
				t.Fatalf("expected success, got: %v", err)
			}

			if !handlerCalled {
				t.Error("handler should have been called")
			}
		},
	)
}

func TestCLI_StrictMode_Integration(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given CLI in strict mode, When command tree is fully described, "+
			"Then all commands pass validation",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"strict", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithStrictValidation[lifecycleConfig](),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			child := newLifecycleCmd(t, "child", "Child command")

			parent := newLifecycleParentCmd(t, child, "Parent command")

			if err := v2.AddCommand(cli, parent); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{"parent", "child"})
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
		},
	)

	t.Run(
		"Given CLI in strict mode, When any command lacks short description, "+
			"Then AddCommand rejects it",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"strict", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithStrictValidation[lifecycleConfig](),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			child, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"child",
				func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					return nil
				},
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			parent := newLifecycleParentCmd(t, child, "Parent")

			err = v2.AddCommand(cli, parent)
			if err == nil {
				t.Fatal("expected strict mode to reject command without short")
			}

			if !errors.Is(err, v2.ErrMissingShort) {
				t.Errorf("expected ErrMissingShort, got: %v", err)
			}
		},
	)
}

func TestCLI_VersionCommand_Integration(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given CLI with version and version command, When 'version' is executed, "+
			"Then output contains app name and version",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"myapp", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
				v2.WithCLIVersion[lifecycleConfig]("3.14.0"),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			vCmd, err := v2.VersionCommand[lifecycleConfig](cli)
			if err != nil {
				t.Fatalf("VersionCommand: %v", err)
			}

			if err := v2.AddCommand(cli, vCmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			var buf strings.Builder
			cli.RootCommand().SetOut(&buf)

			if err := cli.ExecuteWithArgs(context.Background(), []string{"version"}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, "myapp 3.14.0") {
				t.Errorf("output should contain 'myapp 3.14.0', got: %q", output)
			}
		},
	)
}

func TestCLI_FlowContext_Integration(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a running CLI, When handler accesses flow context, "+
			"Then branching flow context is present",
		func(t *testing.T) {
			t.Parallel()

			cli, err := v2.NewCLI[lifecycleConfig](
				"flow", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			var flowCtxPresent bool

			cmd, err := v2.NewCommand[lifecycleConfig, v2.NoFlags](
				"run",
				func(ctx context.Context, _ *lifecycleConfig, _ v2.NoFlags) error {
					_, ok := v2.GetBranchingFlowContext(ctx)
					flowCtxPresent = ok

					return nil
				},
				v2.WithShort[lifecycleConfig, v2.NoFlags]("Run"),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			if err := cli.ExecuteWithArgs(context.Background(), []string{"run"}); err != nil {
				t.Fatalf("Execute: %v", err)
			}

			if !flowCtxPresent {
				t.Error("flow context should be present in handler")
			}
		},
	)

	t.Run(
		"Given a manually branched flow context, When PathString is called, "+
			"Then path reflects branch structure",
		func(t *testing.T) {
			t.Parallel()

			root := v2.NewBranchingFlowContext(context.Background())
			child, cancel := root.Branch("parent")
			defer cancel()

			grandchild, cancel2 := child.Branch("leaf")
			defer cancel2()

			path := grandchild.PathString()
			if !strings.Contains(path, "parent") || !strings.Contains(path, "leaf") {
				t.Errorf("path should contain 'parent' and 'leaf', got: %q", path)
			}
		},
	)
}

func TestCLI_FlagTypes_Integration(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a command with typed flags, When flags are parsed, "+
			"Then all types are correctly deserialized",
		func(t *testing.T) {
			t.Parallel()

			type allTypesFlags struct {
				Name    string  `flag:"name"    default:"test"  help:"Name"`
				Count   int     `flag:"count"   default:"5"     help:"Count"`
				Ratio   float64 `flag:"ratio"   default:"0.5"   help:"Ratio"`
				Enabled bool    `flag:"enabled" default:"false" help:"Enabled"`
			}

			cli, err := v2.NewCLI[lifecycleConfig](
				"flags", "Test", lifecycleConfig{},
				v2.WithFang[lifecycleConfig](false),
			)
			if err != nil {
				t.Fatalf("NewCLI: %v", err)
			}

			var parsedFlags *allTypesFlags

			cmd, err := v2.NewCommand[lifecycleConfig, *allTypesFlags](
				"check",
				func(_ context.Context, _ *lifecycleConfig, flags *allTypesFlags) error {
					parsedFlags = flags

					return nil
				},
				v2.WithShort[lifecycleConfig, *allTypesFlags]("Check"),
				v2.WithFlags[lifecycleConfig, *allTypesFlags](&allTypesFlags{}),
			)
			if err != nil {
				t.Fatalf("NewCommand: %v", err)
			}

			if err := v2.AddCommand(cli, cmd); err != nil {
				t.Fatalf("AddCommand: %v", err)
			}

			err = cli.ExecuteWithArgs(context.Background(), []string{
				"check",
				"--name=hello",
				"--count=42",
				"--ratio=3.14",
				"--enabled=true",
			})
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}

			if parsedFlags.Name != "hello" {
				t.Errorf("Name = %q, want %q", parsedFlags.Name, "hello")
			}

			if parsedFlags.Count != 42 {
				t.Errorf("Count = %d, want %d", parsedFlags.Count, 42)
			}

			if parsedFlags.Ratio != 3.14 {
				t.Errorf("Ratio = %f, want %f", parsedFlags.Ratio, 3.14)
			}

			if !parsedFlags.Enabled {
				t.Error("Enabled should be true")
			}
		},
	)
}
