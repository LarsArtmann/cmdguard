// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func TestGuardedCommand_FullLifecycle(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")
	if root == nil {
		t.Fatal("root is nil")
	}

	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.SetOut(os.Stdout)
			cmd.Println("Hello, World!")
		},
	})

	cmd := root.Command()
	if cmd == nil {
		t.Fatal("cmd is nil")
	}

	if cmd.Name() != "testapp" {
		t.Errorf("cmd.Name() = %q, want %q", cmd.Name(), "testapp")
	}
}

func TestGuardedCommand_PanicOnInvalidCommand(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")

	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		root.AddCommand(&cobra.Command{
			Use:   "invalid",
			Short: "This has no handler",
		})
	}()

	if !didPanic {
		t.Errorf("expected panic for command without handler")
	}
}

func TestGuardedCommand_PanicOnEmptyName(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")

	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		root.AddCommand(&cobra.Command{
			Short: "No name here",
			Run:   func(_ *cobra.Command, _ []string) {},
		})
	}()

	if !didPanic {
		t.Errorf("expected panic for command without name")
	}
}

func TestGuardedCommand_ParentWithChildren(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")

	parent := &cobra.Command{
		Use:   "parent",
		Short: "Parent command",
	}

	child := &cobra.Command{
		Use:   "child",
		Short: "Child command",
		Run:   func(_ *cobra.Command, _ []string) {},
	}

	parent.AddCommand(child)
	root.AddCommand(parent)

	cmd := root.Command()
	if len(cmd.Commands()) < 3 {
		t.Errorf("len(cmd.Commands()) = %d, want at least 3", len(cmd.Commands()))
	}
}

func TestGuardedCommand_StrictMode(t *testing.T) {
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	root := cmdguard.New("testapp", "Test application")
	if !root.IsStrictMode() {
		t.Error("Should be in strict mode")
	}

	root.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Run checks",
		RunE:  func(_ *cobra.Command, _ []string) error { return nil },
	})

	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		root.AddCommand(&cobra.Command{
			Use: "bad",
			Run: func(*cobra.Command, []string) {},
		})
	}()

	if !didPanic {
		t.Errorf("expected panic for Run without RunE in strict mode")
	}
}

func TestGuardedCommand_ConfigAccess(t *testing.T) {
	t.Setenv("CMDGUARD_LOG_LEVEL", "debug")

	root := cmdguard.New("testapp", "Test application")
	cfg := root.Config()

	if cfg == nil {
		t.Fatal("config is nil")
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}

func TestGuardedCommand_AddSubcommand(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")

	parent := &cobra.Command{
		Use:   "db",
		Short: "Database commands",
	}

	child := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE:  func(*cobra.Command, []string) error { return nil },
	}

	root.AddSubcommand(parent, child)
	root.AddCommand(parent)

	cmd := root.Command()

	dbCmd, _, err := cmd.Find([]string{"db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dbCmd.Name() != "db" {
		t.Errorf("dbCmd.Name() = %q, want %q", dbCmd.Name(), "db")
	}
}

func TestGuardedCommand_BuiltInCommands(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")
	cmd := root.Command()

	commands := cmd.Commands()
	if len(commands) < 2 {
		t.Errorf("len(commands) = %d, want at least 2", len(commands))
	}

	var hasVersion bool

	for _, c := range commands {
		if c.Name() == "version" {
			hasVersion = true

			break
		}
	}

	if !hasVersion {
		t.Error("Should have version command")
	}
}

func TestGuardedCommand_ExecuteWithContext(t *testing.T) {
	t.Parallel()

	root := cmdguard.New("testapp", "Test application")

	var output bytes.Buffer

	root.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.SetOut(&output)
			cmd.Println("test output")
		},
	})

	root.Command().SetArgs([]string{"test"})

	ctx := t.Context()

	err := root.Execute(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
