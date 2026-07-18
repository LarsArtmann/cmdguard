package v3_test

import (
	"context"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestInitializeErrorPaths(t *testing.T) {
	t.Parallel()
	t.Run("NewCLI with custom scope works", func(t *testing.T) {
		t.Parallel()

		scope := v3.NewScope("custom")

		cli, err := v3.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v3.WithCLIScope(scope),
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

		cli, err := v3.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type ptrFlags struct {
			Name string `default:"test" flag:"name" help:"Name"`
		}

		cmd, err := v3.NewCommand(
			"ptrcmd",
			&ptrFlags{},
			func(_ context.Context, _ *testCLIConfig, _ *ptrFlags) error {
				return nil
			},
			v3.WithShort("Command with pointer flags"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with nil pointer flags failed: %v", err)
		}
	})

	t.Run("AddCommand with nested subcommands", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test CLI", testCLIConfig{})
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

		err = v3.AddCommand(cli, parentCmd)
		if err != nil {
			t.Fatalf("AddCommand with nested subcommands failed: %v", err)
		}
	})

	t.Run("AddCommand with command missing RunE and subcommands", func(t *testing.T) {
		t.Parallel()

		_, err := v3.NewCommand[testCLIConfig](
			"norun", v3.NoFlags{}, nil,
			v3.WithShort("Command without RunE"),
		)
		if err == nil {
			t.Fatal("expected error for command without RunE and no subcommands")
		}
	})

	t.Run("AddCommand with empty Use field", func(t *testing.T) {
		t.Parallel()

		_, err := v3.NewCommand(
			"",
			v3.NoFlags{},
			noOpRunE[testCLIConfig],
		)
		if err == nil {
			t.Fatal("expected error for empty Use field")
		}
	})

	t.Run("AddCommand with flags containing multiple fields", func(t *testing.T) {
		t.Parallel()

		cli, err := v3.NewCLI("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		type multiFlags struct {
			Name    string `default:"World" flag:"name"    help:"Name"`
			Count   int    `default:"1"     flag:"count"   help:"Count"`
			Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
		}

		cmd, err := v3.NewCommand(
			"multi",
			multiFlags{},
			func(_ context.Context, _ *testCLIConfig, _ multiFlags) error {
				return nil
			},
			v3.WithShort("Multi-flag command"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		err = v3.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with multiple flags failed: %v", err)
		}
	})
}
