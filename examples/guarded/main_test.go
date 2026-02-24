// Integration test for guarded example
package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestGuardedExample_ValidRunCommand(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	var output bytes.Buffer
	root.AddCommand(&cobra.Command{
		Use:   "valid-run",
		Short: "This command is valid (has Run)",
		Run: func(cmd *cobra.Command, args []string) {
			output.WriteString("Run command executed\n")
		},
	})

	root.Command().SetArgs([]string{"valid-run"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Run command executed\n", output.String())
}

func TestGuardedExample_ValidRunECommand(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	var output bytes.Buffer
	root.AddCommand(&cobra.Command{
		Use:   "valid-rune",
		Short: "This command is valid (has RunE)",
		RunE: func(cmd *cobra.Command, args []string) error {
			output.WriteString("RunE command executed\n")
			return nil
		},
	})

	root.Command().SetArgs([]string{"valid-rune"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "RunE command executed\n", output.String())
}

func TestGuardedExample_ValidParentWithChildren(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	var output bytes.Buffer
	parent := &cobra.Command{
		Use:   "parent",
		Short: "Parent command (valid because it has children)",
	}
	parent.AddCommand(&cobra.Command{
		Use:   "child",
		Short: "Child command",
		Run: func(cmd *cobra.Command, args []string) {
			output.WriteString("Child command executed\n")
		},
	})
	root.AddCommand(parent)

	root.Command().SetArgs([]string{"parent", "child"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Child command executed\n", output.String())
}

func TestGuardedExample_StrictModeDisabled(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	// By default, strict mode is disabled
	assert.False(t, root.IsStrictMode())
}

func TestGuardedExample_StrictModeEnabled(t *testing.T) {
	// Set strict mode via environment variable
	t.Setenv("CMDGUARD_STRICT_MODE", "true")

	root := cmdguard.New("guarded", "Demonstrates guard validation")
	assert.True(t, root.IsStrictMode())
}

func TestGuardedExample_MultipleValidCommands(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	root.AddCommand(&cobra.Command{
		Use:   "cmd1",
		Short: "First command",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	root.AddCommand(&cobra.Command{
		Use:   "cmd2",
		Short: "Second command",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	})

	root.AddCommand(&cobra.Command{
		Use:   "cmd3",
		Short: "Third command",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	cmd := root.Command()
	// GuardedCommand adds built-in commands, so we just verify our commands were added
	assert.GreaterOrEqual(t, len(cmd.Commands()), 3)
}

func TestGuardedExample_PanicOnInvalidCommand(t *testing.T) {
	// This test documents that invalid commands panic
	// We cannot test this directly as it would crash the test
	// But we document the expected behavior

	assert.Panics(t, func() {
		root := cmdguard.New("guarded", "Demonstrates guard validation")
		// This will panic because the command has no handler
		root.AddCommand(&cobra.Command{
			Use:   "invalid",
			Short: "This command has no handler - will panic",
		})
	})
}

func TestGuardedExample_PanicOnEmptyName(t *testing.T) {
	assert.Panics(t, func() {
		root := cmdguard.New("guarded", "Demonstrates guard validation")
		root.AddCommand(&cobra.Command{
			Use:   "",
			Short: "Empty name - will panic",
			Run:   func(cmd *cobra.Command, args []string) {},
		})
	})
}

func TestGuardedExample_CommandAccess(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	cmd := root.Command()
	assert.Equal(t, "guarded", cmd.Use)
	assert.Equal(t, "Demonstrates guard validation", cmd.Short)
}

func TestGuardedExample_OutputCapture(t *testing.T) {
	root := cmdguard.New("guarded", "Demonstrates guard validation")

	var capturedOutput string
	root.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			capturedOutput = "Test output captured"
		},
	})

	root.Command().SetArgs([]string{"test"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Test output captured", capturedOutput)
}
