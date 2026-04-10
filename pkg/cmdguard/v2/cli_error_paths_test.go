package v2_test

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
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

		cmd := v2.Command[testCLIConfig, *ptrFlags]{
			Use:   "ptrcmd",
			Short: "Command with pointer flags",
			Flags: nil,
			RunE: func(_ context.Context, _ *testCLIConfig, _ *ptrFlags) error {
				return nil
			},
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

		childCmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "child",
			Short: "Child command",
			RunE: func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				return nil
			},
		}

		parentCmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "parent",
			Short: "Parent command",
			Commands: []v2.Command[testCLIConfig, v2.NoFlags]{
				childCmd,
			},
		}

		err = v2.AddCommand(cli, parentCmd)
		if err != nil {
			t.Fatalf("AddCommand with nested subcommands failed: %v", err)
		}
	})

	t.Run("AddCommand with command missing RunE and subcommands", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use:   "norun",
			Short: "Command without RunE",
		}

		err = v2.AddCommand(cli, cmd)
		if err == nil {
			t.Fatal("expected error for command without RunE and no subcommands")
		}
	})

	t.Run("AddCommand with empty Use field", func(t *testing.T) {
		t.Parallel()

		cli, err := v2.NewCLI[testCLIConfig]("test", "Test CLI", testCLIConfig{})
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd := v2.Command[testCLIConfig, v2.NoFlags]{
			Use: "",
			RunE: func(_ context.Context, _ *testCLIConfig, _ v2.NoFlags) error {
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
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

		cmd := v2.Command[testCLIConfig, multiFlags]{
			Use:   "multi",
			Short: "Multi-flag command",
			Flags: multiFlags{},
			RunE: func(_ context.Context, _ *testCLIConfig, _ multiFlags) error {
				return nil
			},
		}

		err = v2.AddCommand(cli, cmd)
		if err != nil {
			t.Fatalf("AddCommand with multiple flags failed: %v", err)
		}
	})
}
