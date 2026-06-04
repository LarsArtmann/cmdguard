package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func TestInitializeErrorPaths(t *testing.T) {
	t.Parallel()
	t.Run("NewCLI with custom scope works", func(t *testing.T) {
		t.Parallel()

		scope := v2.NewScope("custom")

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test CLI", testCLIConfig{},
			v2.WithCLIScope[testCLIConfig](scope),
		)
		if err != nil {
			t.Fatalf("NewCLI with custom scope failed: %v", err)
		}

		if cli.Scope() != scope {
			t.Error("expected scope to be the custom scope")
		}
	})

	t.Run("AddCommand with nil pointer flags creates prototype", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type ptrFlags struct {
			Name string `default:"test" flag:"name" help:"Name"`
		}

		cmd, err := v2.NewCommand[testCLIConfig, *ptrFlags](
			"ptrcmd",
			func(_ context.Context, _ *testCLIConfig, _ *ptrFlags) error {
				return nil
			},
			v2.WithShort[testCLIConfig, *ptrFlags]("Command with pointer flags"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with nil pointer flags failed: %v", err)
		}
	})

	t.Run("AddCommand with nested subcommands", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		childCmd := newTestCLICommandWithShort[testCLIConfig]("child", "Child command")
		parentCmd := newTestParentCommand[testCLIConfig](
			"parent",
			"Parent command",
			"Parent command long description",
			childCmd,
		)

		err = v2.AddCommand(cli, parentCmd)
		if err != nil {
			t.Fatalf("AddCommand with nested subcommands failed: %v", err)
		}
	})

	t.Run("AddCommand with command missing RunE and subcommands", func(t *testing.T) {
		t.Parallel()

		_, err := v2.NewCommand[testCLIConfig, v2.NoFlags](
			"norun", nil,
			v2.WithShort[testCLIConfig, v2.NoFlags]("Command without RunE"),
		)
		if err == nil {
			t.Fatal("expected error for command without RunE and no subcommands")
		}
	})

	t.Run("AddCommand with empty Use field", func(t *testing.T) {
		t.Parallel()

		_, err := v2.NewCommand[testCLIConfig, v2.NoFlags](
			"",
			noOpRunE[testCLIConfig],
		)
		if err == nil {
			t.Fatal("expected error for empty Use field")
		}
	})

	t.Run("AddCommand with flags containing multiple fields", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type multiFlags struct {
			Name    string `default:"World" flag:"name"    help:"Name"`
			Count   int    `default:"1"     flag:"count"   help:"Count"`
			Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
		}

		cmd, err := v2.NewCommand[testCLIConfig, multiFlags](
			"multi",
			func(_ context.Context, _ *testCLIConfig, _ multiFlags) error {
				return nil
			},
			v2.WithShort[testCLIConfig, multiFlags]("Multi-flag command"),
			v2.WithFlags[testCLIConfig, multiFlags](multiFlags{}),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with multiple flags failed: %v", err)
		}
	})
}
