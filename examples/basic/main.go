// Basic example demonstrating simple cmdguard usage.
package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
)

// newCommand creates a simple cobra command with the given name and run function.
func newCommand(use, short string, run func(*cobra.Command, []string)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run:   run,
	}
}

func main() {
	// Create guarded CLI
	root := cmdguard.New("basic", "A basic CLI example")

	// Add commands
	root.AddCommand(newCommand("hello", "Say hello", func(_ *cobra.Command, _ []string) {
		fmt.Println("Hello, World!")
	}))
	root.AddCommand(newCommand("goodbye", "Say goodbye", func(_ *cobra.Command, _ []string) {
		fmt.Println("Goodbye, World!")
	}))

	// Execute
	root.ExecuteAndExit(context.Background())
}
