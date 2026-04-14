package cmdguard

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

// findCommand iterates over subcommands and returns the first one matching the predicate.
func findCommand(g *GuardedCommand, name string) *cobra.Command {
	for _, cmd := range g.cmd.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}

	return nil
}

// hasSubcommand checks if a command with the given name exists in the GuardedCommand's subcommands.
func hasSubcommand(g *GuardedCommand, name string) bool {
	return findCommand(g, name) != nil
}

// findSubcommand finds and returns a subcommand with the given name, or nil if not found.
func findSubcommand(g *GuardedCommand, name string) *cobra.Command {
	return findCommand(g, name)
}

//nolint:paralleltest // uses t.Setenv
func TestGuardedCommand_validateCommand(t *testing.T) {
	t.Run("command with Run is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := newCobraCommand("sub")

		err := g.validateCommand(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("command with RunE is valid", func(t *testing.T) {
		g := New("testapp", "Test")

		cmd := newTestCommand("sub")

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

		testutil.AssertErrorContains(t, err, "no name")
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

		testutil.AssertErrorContains(t, err, "no handler")
	})

	t.Run("parent command with subcommands does not need handler", func(t *testing.T) {
		g := New("testapp", "Test")

		parent := &cobra.Command{Use: "parent"}
		child := newCobraCommand("child")
		parent.AddCommand(child)

		err := g.validateCommand(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("strict mode requires RunE", func(t *testing.T) {
		t.Setenv("CMDGUARD_STRICT_MODE", "true")

		g := New("testapp", "Test")

		cmd := newCobraCommand("sub")

		err := g.validateCommand(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		testutil.AssertErrorContains(t, err, "strict mode requires RunE")
	})
}

func TestGuardedCommand_DefaultCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cmdName string
	}{
		{"version command is added", "version"},
		{"validate command is added", "validate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := New("testapp", "Test")

			if !hasSubcommand(g, tt.cmdName) {
				t.Errorf("%s command should be added by default", tt.cmdName)
			}
		})
	}
}

func TestGuardedCommand_ValidateCommandTree(t *testing.T) {
	t.Parallel()
	t.Run("validates command tree successfully", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		// Add a valid command with subcommand
		parent := &cobra.Command{Use: "parent"}
		child := newCobraCommand("child")
		parent.AddCommand(child)
		g.AddCommand(parent)

		err := g.validateCommandTree()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid nested command", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		// Create a command tree with an invalid nested subcommand
		parent := newCobraCommand("parent")
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)
		g.cmd.AddCommand(parent)

		err := g.validateCommandTree()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		testutil.AssertErrorContains(t, err, "invalid")
	})
}

func TestGuardedCommand_ValidateSubcommands(t *testing.T) {
	t.Parallel()
	t.Run("validates valid subcommands", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		parent := newCobraCommand("parent")
		child := newCobraCommand("child")
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns error for invalid subcommand", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		parent := newCobraCommand("parent")
		// Manually add invalid child (bypassing guard)
		invalidChild := &cobra.Command{Use: "invalid"} // No handler
		parent.AddCommand(invalidChild)

		err := g.validateSubcommands(parent)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		testutil.AssertErrorContains(t, err, "parent invalid")
	})

	t.Run("validates nested subcommands", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		parent := newCobraCommand("parent")
		child := newCobraCommand("child")
		grandchild := newCobraCommand("grandchild")
		child.AddCommand(grandchild)
		parent.AddCommand(child)

		err := g.validateSubcommands(parent)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGuardedCommand_DefaultCommands_Execution(t *testing.T) {
	t.Parallel()
	t.Run("version command executes", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		g := New("testapp", "Test")

		// Add a valid command
		cmd := newCobraCommand("testcmd")
		g.AddCommand(cmd)

		validateCmd := findSubcommand(g, "validate")
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
		t.Parallel()

		g := New("testapp", "Test")

		// Manually add invalid command (bypassing guard)
		invalidCmd := &cobra.Command{Use: "invalid"} // No handler
		g.cmd.AddCommand(invalidCmd)

		validateCmd := findSubcommand(g, "validate")
		if validateCmd == nil {
			t.Fatal("validate command not found")
		}

		// Execute it - should error
		err := validateCmd.RunE(validateCmd, []string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		testutil.AssertErrorContains(t, err, "validation failed")
	})
}
