// Integration test for basic example
package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func TestBasicExample_HelloCommand(t *testing.T) {
	t.Parallel()
	root := cmdguard.New("basic", "A basic CLI example")

	var output bytes.Buffer

	root.AddCommand(newCommand("hello", "Say hello", func(_ *cobra.Command, _ []string) {
		output.WriteString("Hello, World!\n")
	}))

	root.AddCommand(newCommand("goodbye", "Say goodbye", func(_ *cobra.Command, _ []string) {
		output.WriteString("Goodbye, World!\n")
	}))

	// Test hello command
	root.Command().SetArgs([]string{"hello"})

	err := root.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.String() != "Hello, World!\n" {
		t.Errorf("output = %q, want %q", output.String(), "Hello, World!\n")
	}

	// Reset output for next test
	output.Reset()

	// Test goodbye command
	root.Command().SetArgs([]string{"goodbye"})

	err = root.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.String() != "Goodbye, World!\n" {
		t.Errorf("output = %q, want %q", output.String(), "Goodbye, World!\n")
	}
}

func TestBasicExample_RootHasSubcommands(t *testing.T) {
	t.Parallel()
	root := cmdguard.New("basic", "A basic CLI example")

	emptyRun := func(_ *cobra.Command, _ []string) {}
	root.AddCommand(newCommand("hello", "Say hello", emptyRun))
	root.AddCommand(newCommand("goodbye", "Say goodbye", emptyRun))

	cmd := root.Command()
	if cmd.Use != "basic" {
		t.Errorf("cmd.Use = %q, want %q", cmd.Use, "basic")
	}

	if cmd.Short != "A basic CLI example" {
		t.Errorf("cmd.Short = %q, want %q", cmd.Short, "A basic CLI example")
	}

	// GuardedCommand adds built-in commands (completion, help, validate, version + our 2 custom commands)
	if len(cmd.Commands()) < 2 {
		t.Errorf("len(cmd.Commands()) = %d, want at least 2", len(cmd.Commands()))
	}
}

func TestBasicExample_HelpOutput(t *testing.T) {
	t.Parallel()
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping help output test in CI")
	}

	root := cmdguard.New("basic", "A basic CLI example")

	emptyRun := func(_ *cobra.Command, _ []string) {}
	root.AddCommand(newCommand("hello", "Say hello", emptyRun))
	root.AddCommand(newCommand("goodbye", "Say goodbye", emptyRun))

	// Set args to run help
	root.Command().SetArgs([]string{"--help"})

	// Execute should not panic, we just want to verify it doesn't crash
	_ = root.Execute(context.Background())
}
