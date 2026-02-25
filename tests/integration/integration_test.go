// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func TestGuardedCommand_FullLifecycle(t *testing.T) {
	// Create a new guarded command
	root := cmdguard.New("testapp", "Test application")
	require.NotNil(t, root)

	// Add a simple command
	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(os.Stdout)
			cmd.Println("Hello, World!")
		},
	})

	// Verify we can access the underlying command
	cmd := root.Command()
	require.NotNil(t, cmd)
	assert.Equal(t, "testapp", cmd.Name())
}

func TestGuardedCommand_PanicOnInvalidCommand(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")

	// Adding a command without a handler should panic
	assert.Panics(t, func() {
		root.AddCommand(&cobra.Command{
			Use:   "invalid",
			Short: "This has no handler",
		})
	}, "Adding command without handler should panic")
}

func TestGuardedCommand_PanicOnEmptyName(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")

	// Adding a command with empty name should panic
	assert.Panics(t, func() {
		root.AddCommand(&cobra.Command{
			Short: "No name here",
			Run:   func(cmd *cobra.Command, args []string) {},
		})
	}, "Adding command without name should panic")
}

func TestGuardedCommand_ParentWithChildren(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")

	// Parent command without handler is valid when it has children
	parent := &cobra.Command{
		Use:   "parent",
		Short: "Parent command",
	}

	// Child command with handler
	child := &cobra.Command{
		Use:   "child",
		Short: "Child command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	// This should not panic because parent has children
	parent.AddCommand(child)
	root.AddCommand(parent)

	// Verify structure (built-in: validate, version + parent)
	cmd := root.Command()
	assert.Len(t, cmd.Commands(), 3)
}

func TestGuardedCommand_StrictMode(t *testing.T) {
	// Set strict mode via environment (errors ignored for test setup)
	_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")

	defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()

	root := cmdguard.New("testapp", "Test application")
	assert.True(t, root.IsStrictMode(), "Should be in strict mode")

	// RunE should work in strict mode
	root.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Run checks",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	})

	// Run (not RunE) should panic in strict mode
	assert.Panics(t, func() {
		root.AddCommand(&cobra.Command{
			Use: "bad",
			Run: func(cmd *cobra.Command, args []string) {},
		})
	}, "Run without RunE should panic in strict mode")
}

func TestGuardedCommand_ConfigAccess(t *testing.T) {
	// Set custom log level (errors ignored for test setup)
	_ = os.Setenv("CMDGUARD_LOG_LEVEL", "debug")

	defer func() { _ = os.Unsetenv("CMDGUARD_LOG_LEVEL") }()

	root := cmdguard.New("testapp", "Test application")
	cfg := root.Config()

	require.NotNil(t, cfg)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestGuardedCommand_AddSubcommand(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")

	parent := &cobra.Command{
		Use:   "db",
		Short: "Database commands",
	}

	child := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	}

	// AddSubcommand should work
	root.AddSubcommand(parent, child)
	root.AddCommand(parent)

	// Verify
	cmd := root.Command()
	dbCmd, _, err := cmd.Find([]string{"db"})
	require.NoError(t, err)
	assert.Equal(t, "db", dbCmd.Name())
}

func TestGuardedCommand_BuiltInCommands(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")
	cmd := root.Command()

	// Check that built-in commands exist
	commands := cmd.Commands()
	require.GreaterOrEqual(
		t,
		len(commands),
		2,
		"Should have at least version and validate commands",
	)

	// Find version command
	var hasVersion bool

	for _, c := range commands {
		if c.Name() == "version" {
			hasVersion = true

			break
		}
	}

	assert.True(t, hasVersion, "Should have version command")
}

func TestGuardedCommand_ExecuteWithContext(t *testing.T) {
	root := cmdguard.New("testapp", "Test application")

	// Add a command that writes to output
	var output bytes.Buffer

	root.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(&output)
			cmd.Println("test output")
		},
	})

	// Set args to run our test command
	root.Command().SetArgs([]string{"test"})

	// Execute should work with context
	ctx := context.Background()
	err := root.Execute(ctx)

	// Note: Execute may fail since we're not running from main
	// but it shouldn't panic
	assert.NoError(t, err)
}
