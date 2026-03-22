package cmdguard

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// hasSubcommand checks if a command with the given name exists in the GuardedCommand's subcommands.
func hasSubcommand(g *GuardedCommand, name string) bool {
	for _, cmd := range g.cmd.Commands() {
		if cmd.Name() == name {
			return true
		}
	}

	return false
}

func TestNew(t *testing.T) {
	t.Run("creates GuardedCommand with defaults", func(t *testing.T) {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")

		g := New("testapp", "Test application")

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}
		if g.cmd == nil {
			t.Error("expected non-nil cmd")
		}
		if g.cfg == nil {
			t.Error("expected non-nil cfg")
		}
		if g.cmd.Use != "testapp" {
			t.Errorf("cmd.Use = %q, want %q", g.cmd.Use, "testapp")
		}
		if g.cmd.Short != "Test application" {
			t.Errorf("cmd.Short = %q, want %q", g.cmd.Short, "Test application")
		}
	})

	t.Run("loads config from environment", func(t *testing.T) {
		t.Setenv("CMDGUARD_LOG_LEVEL", "debug")
		t.Setenv("CMDGUARD_STRICT_MODE", "true")

		g := New("testapp", "Test")

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}
		if g.cfg.LogLevel != "debug" {
			t.Errorf("cfg.LogLevel = %q, want %q", g.cfg.LogLevel, "debug")
		}
		if !g.cfg.StrictMode {
			t.Error("cfg.StrictMode = false, want true")
		}
		if !g.strictMode {
			t.Error("strictMode = false, want true")
		}
	})
}

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Run("accepts valid command with Run", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddCommand(cmd)
		}()

		if didPanic {
			t.Error("AddCommand should not panic for valid command")
		}
	})

	t.Run("accepts valid command with RunE", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddCommand(cmd)
		}()

		if didPanic {
			t.Error("AddCommand should not panic for valid command")
		}
	})

	t.Run("panics on command without handler", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "invalid",
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddCommand(cmd)
		}()

		if !didPanic {
			t.Error("AddCommand should panic for command without handler")
		}
	})

	t.Run("panics on command without name", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddCommand(cmd)
		}()

		if !didPanic {
			t.Error("AddCommand should panic for command without name")
		}
	})

	t.Run("panics after Execute called", func(t *testing.T) {
		g := New("testapp", "Test")
		g.validated = true // Simulate post-execute state

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddCommand(cmd)
		}()

		if !didPanic {
			t.Error("AddCommand should panic after Execute called")
		}
	})
}

func TestGuardedCommand_AddSubcommand(t *testing.T) {
	t.Run("adds subcommand to parent", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddSubcommand(parent, child)
		}()

		if didPanic {
			t.Error("AddSubcommand should not panic for valid child")
		}

		found := false
		for _, c := range parent.Commands() {
			if c.Name() == "child" {
				found = true
				break
			}
		}

		if !found {
			t.Error("child command should be added to parent")
		}
	})

	t.Run("panics on invalid child", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "invalid-child",
			// No Run or RunE
		}

		didPanic := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
			}()
			g.AddSubcommand(parent, child)
		}()

		if !didPanic {
			t.Error("AddSubcommand should panic for invalid child")
		}
	})
}

func TestGuardedCommand_Execute(t *testing.T) {
	t.Run("executes command successfully", func(t *testing.T) {
		g := New("testapp", "Test")
		g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return nil
		}

		err := g.Execute(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !g.validated {
			t.Error("validated = false, want true")
		}
	})
}

func TestGuardedCommand_Accessors(t *testing.T) {
	t.Run("Command returns underlying cobra command", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := g.Command()

		if cmd == nil {
			t.Fatal("expected non-nil command")
		}
		if cmd.Use != "testapp" {
			t.Errorf("cmd.Use = %q, want %q", cmd.Use, "testapp")
		}
	})

	t.Run("Config returns config", func(t *testing.T) {
		g := New("testapp", "Test")

		cfg := g.Config()

		if cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if cfg.LogLevel != "info" {
			t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, "info")
		}
	})

	t.Run("IsStrictMode returns correct value", func(t *testing.T) {
		g := New("testapp", "Test")
		if g.IsStrictMode() {
			t.Error("IsStrictMode() = true, want false")
		}

		t.Setenv("CMDGUARD_STRICT_MODE", "true")

		g2 := New("testapp2", "Test2")
		if !g2.IsStrictMode() {
			t.Error("IsStrictMode() = false, want true")
		}
	})
}

func TestGuardedCommand_validateCommand(t *testing.T) {
	t.Run("command with Run is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		err := g.validateCommand(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("command with RunE is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		err := g.validateCommand(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("command without name is invalid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{}

		err := g.validateCommand(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no name") {
			t.Errorf("error should contain 'no name', got: %v", err)
		}
	})

	t.Run("command without handler is invalid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "invalid",
		}

		err := g.validateCommand(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no handler") {
			t.Errorf("error should contain 'no handler', got: %v", err)
		}
	})

	t.Run("parent command with subcommands does not need handler", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)

		err := g.validateCommand(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("strict mode requires RunE", func(t *testing.T) {
		t.Setenv("CMDGUARD_STRICT_MODE", "true")

		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		err := g.validateCommand(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "strict mode requires RunE") {
			t.Errorf("error should contain 'strict mode requires RunE', got: %v", err)
		}
	})
}

func TestGuardedCommand_DefaultCommands(t *testing.T) {
	t.Run("version command is added", func(t *testing.T) {
		g := New("testapp", "Test")

		if !hasSubcommand(g, "version") {
			t.Error("version command should be added by default")
		}
	})

	t.Run("validate command is added", func(t *testing.T) {
		g := New("testapp", "Test")

		if !hasSubcommand(g, "validate") {
			t.Error("validate command should be added by default")
		}
	})
}

func TestVersion(t *testing.T) {
	t.Run("returns version string", func(t *testing.T) {
		v := Version()
		if v == "" {
			t.Error("version should not be empty")
		}
		if v != "dev" {
			t.Errorf("Version() = %q, want %q", v, "dev")
		}
	})
}

func TestGuardedCommand_ValidateCommandTree(t *testing.T) {
	t.Run("validates command tree successfully", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command with subcommand
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)
		g.AddCommand(parent)

		err := g.validateCommandTree()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid nested command", func(t *testing.T) {
		g := New("testapp", "Test")

		// Create a command tree with an invalid nested subcommand
		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		// Manually add an invalid child (bypassing guard validation)
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)
		g.cmd.AddCommand(parent)

		err := g.validateCommandTree()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid") {
			t.Errorf("error should contain 'invalid', got: %v", err)
		}
	})
}

func TestGuardedCommand_ValidateSubcommands(t *testing.T) {
	t.Run("validates valid subcommands", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid subcommand", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		// Manually add invalid child (bypassing guard)
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)

		err := g.validateSubcommands(parent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "parent invalid") {
			t.Errorf("error should contain 'parent invalid', got: %v", err)
		}
	})

	t.Run("validates nested subcommands", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		grandchild := &cobra.Command{
			Use: "grandchild",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		child.AddCommand(grandchild)
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGuardedCommand_DefaultCommands_Execution(t *testing.T) {
	t.Run("version command executes", func(t *testing.T) {
		g := New("testapp", "Test")

		// Find version command
		var versionCmd *cobra.Command
		for _, cmd := range g.cmd.Commands() {
			if cmd.Name() == "version" {
				versionCmd = cmd
				break
			}
		}

		if versionCmd == nil {
			t.Fatal("version command not found")
		}

		// Execute it
		versionCmd.Run(versionCmd, []string{})
		// If we get here without panic, it worked
	})

	t.Run("validate command succeeds with valid tree", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command
		cmd := &cobra.Command{
			Use: "testcmd",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		g.AddCommand(cmd)

		// Find validate command
		var validateCmd *cobra.Command
		for _, c := range g.cmd.Commands() {
			if c.Name() == "validate" {
				validateCmd = c
				break
			}
		}

		if validateCmd == nil {
			t.Fatal("validate command not found")
		}

		// Execute it
		err := validateCmd.RunE(validateCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("validate command fails with invalid tree", func(t *testing.T) {
		g := New("testapp", "Test")

		// Manually add invalid command (bypassing guard)
		invalidCmd := &cobra.Command{Use: "invalid"} // No handler
		g.cmd.AddCommand(invalidCmd)

		// Find validate command
		var validateCmd *cobra.Command
		for _, c := range g.cmd.Commands() {
			if c.Name() == "validate" {
				validateCmd = c
				break
			}
		}

		if validateCmd == nil {
			t.Fatal("validate command not found")
		}

		// Execute it - should error
		err := validateCmd.RunE(validateCmd, []string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("error should contain 'validation failed', got: %v", err)
		}
	})
}
