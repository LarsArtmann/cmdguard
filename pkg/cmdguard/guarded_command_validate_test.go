package cmdguard

import (
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

func TestGuardedCommand_validateCommand(t *testing.T) {
	t.Run("command with Run is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(*cobra.Command, []string) {},
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
			RunE: func(*cobra.Command, []string) error {
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
			Run: func(*cobra.Command, []string) {},
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
			Run: func(*cobra.Command, []string) {},
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
	tests := []struct {
		name    string
		cmdName string
	}{
		{"version command is added", "version"},
		{"validate command is added", "validate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New("testapp", "Test")

			if !hasSubcommand(g, tt.cmdName) {
				t.Errorf("%s command should be added by default", tt.cmdName)
			}
		})
	}
}

func TestGuardedCommand_ValidateCommandTree(t *testing.T) {
	t.Run("validates command tree successfully", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command with subcommand
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{
			Use: "child",
			Run: func(*cobra.Command, []string) {},
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
			Run: func(*cobra.Command, []string) {},
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
			Run: func(*cobra.Command, []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(*cobra.Command, []string) {},
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
			Run: func(*cobra.Command, []string) {},
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
			Run: func(*cobra.Command, []string) {},
		}
		child := &cobra.Command{
			Use: "child",
			Run: func(*cobra.Command, []string) {},
		}
		grandchild := &cobra.Command{
			Use: "grandchild",
			Run: func(*cobra.Command, []string) {},
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
	})

	t.Run("validate command succeeds with valid tree", func(t *testing.T) {
		g := New("testapp", "Test")

		// Add a valid command
		cmd := &cobra.Command{
			Use: "testcmd",
			Run: func(*cobra.Command, []string) {},
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
