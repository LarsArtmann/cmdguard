// Basic example demonstrating simple cmdguard usage.
package main

import (
	"context"
	"fmt"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
)

func main() {
	// Create guarded CLI
	root := cmdguard.New("basic", "A basic CLI example")

	// Add hello command
	root.AddCommand(&cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	})

	// Add goodbye command
	root.AddCommand(&cobra.Command{
		Use:   "goodbye",
		Short: "Say goodbye",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Goodbye, World!")
		},
	})

	// Execute
	root.ExecuteAndExit(context.Background())
}
