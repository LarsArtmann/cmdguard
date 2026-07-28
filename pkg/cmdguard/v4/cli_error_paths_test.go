package v4_test

import (
	"context"
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestInitializeErrorPaths(t *testing.T) {
	t.Parallel()
	t.Run("NewCLI with custom scope works", func(t *testing.T) {
		t.Parallel()

		scope := v4.NewScope("custom")

		cli, err := v4.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v4.WithCLIScope(scope),
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

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type ptrFlags struct {
			Name string `default:"test" flag:"name" help:"Name"`
		}

		cmd, err := v4.NewCommand(
			"ptrcmd",
			&ptrFlags{},
			func(_ context.Context, _ *testCLIConfig, _ *ptrFlags) error {
				return nil
			},
			v4.WithShort("Command with pointer flags"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with nil pointer flags failed: %v", err)
		}
	})

	t.Run("AddCommand with nested subcommands", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		childCmd := newTestCLICommandWithShort[testCLIConfig](t, "child", "Child command")
		parentCmd := newTestParentCommand(
			t,
			"parent",
			"Parent command",
			"Parent command long description",
			childCmd,
		)

		err = v4.AddCommand(cli, parentCmd)
		if err != nil {
			t.Fatalf("AddCommand with nested subcommands failed: %v", err)
		}
	})

	t.Run("AddCommand with command missing RunE and subcommands", func(t *testing.T) {
		t.Parallel()

		_, err := v4.NewCommand[testCLIConfig](
			"norun", v4.NoFlags{}, nil,
			v4.WithShort("Command without RunE"),
		)
		if err == nil {
			t.Fatal("expected error for command without RunE and no subcommands")
		}
	})

	t.Run("AddCommand with empty Use field", func(t *testing.T) {
		t.Parallel()

		_, err := v4.NewCommand(
			"",
			v4.NoFlags{},
			noOpRunE[testCLIConfig],
		)
		if err == nil {
			t.Fatal("expected error for empty Use field")
		}
	})

	t.Run("AddCommand with flags containing multiple fields", func(t *testing.T) {
		t.Parallel()

		cli, err := v4.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type multiFlags struct {
			Name    string `default:"World" flag:"name"    help:"Name"`
			Count   int    `default:"1"     flag:"count"   help:"Count"`
			Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
		}

		cmd, err := v4.NewCommand(
			"multi",
			multiFlags{},
			func(_ context.Context, _ *testCLIConfig, _ multiFlags) error {
				return nil
			},
			v4.WithShort("Multi-flag command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v4.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with multiple flags failed: %v", err)
		}
	})
}
