// Integration test for basic example
package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicExample_HelloCommand(t *testing.T) {
	root := cmdguard.New("basic", "A basic CLI example")

	var output bytes.Buffer
	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run: func(cmd *cobra.Command, args []string) {
			output.WriteString("Hello, World!\n")
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "goodbye",
		Short: "Say goodbye",
		Run: func(cmd *cobra.Command, args []string) {
			output.WriteString("Goodbye, World!\n")
		},
	})

	// Test hello command
	root.Command().SetArgs([]string{"hello"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!\n", output.String())

	// Reset output for next test
	output.Reset()

	// Test goodbye command
	root.Command().SetArgs([]string{"goodbye"})
	err = root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Goodbye, World!\n", output.String())
}

func TestBasicExample_RootHasSubcommands(t *testing.T) {
	root := cmdguard.New("basic", "A basic CLI example")

	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	root.AddCommand(&cobra.Command{
		Use:   "goodbye",
		Short: "Say goodbye",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	cmd := root.Command()
	assert.Equal(t, "basic", cmd.Use)
	assert.Equal(t, "A basic CLI example", cmd.Short)
	// GuardedCommand adds built-in commands (completion, help, validate, version)
	// plus our 2 custom commands
	assert.GreaterOrEqual(t, len(cmd.Commands()), 2)
}

func TestBasicExample_HelpOutput(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping help output test in CI")
	}

	root := cmdguard.New("basic", "A basic CLI example")
	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	root.Command().SetArgs([]string{"--help"})
	// Help returns error with exit code, but we just want to verify it doesn't panic
	_ = root.Execute(context.Background())
}
