// Guarded example demonstrating panic behavior on invalid commands.
// This example intentionally shows what happens when you try to add invalid commands.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
)

func main() {
	fmt.Println("=== Guarded CLI Example ===")
	fmt.Println()
	fmt.Println("This example demonstrates cmdguard's compile-time validation.")
	fmt.Println("Invalid commands will cause a panic at startup.")
	fmt.Println()

	// You can enable strict mode via environment variable
	// os.Setenv("CMDGUARD_STRICT_MODE", "true")

	root := cmdguard.New("guarded", "Demonstrates guard validation")

	// Valid: Command with Run handler
	root.AddCommand(&cobra.Command{
		Use:   "valid-run",
		Short: "This command is valid (has Run)",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	// Valid: Command with RunE handler
	root.AddCommand(&cobra.Command{
		Use:   "valid-rune",
		Short: "This command is valid (has RunE)",
		RunE:  func(cmd *cobra.Command, args []string) error { return nil },
	})

	// Valid: Parent command with children
	parent := &cobra.Command{
		Use:   "parent",
		Short: "Parent command (valid because it has children)",
	}
	parent.AddCommand(&cobra.Command{
		Use:   "child",
		Short: "Child command",
		Run:   func(cmd *cobra.Command, args []string) {},
	})
	root.AddCommand(parent)

	// UNCOMMENT TO SEE PANIC:
	// The following would panic because the command has no handler:
	//
	// root.AddCommand(&cobra.Command{
	//     Use:   "invalid",
	//     Short: "This command has no handler!",
	// })

	// In strict mode, the following would also panic:
	// (Run is not allowed, only RunE)
	//
	// os.Setenv("CMDGUARD_STRICT_MODE", "true")
	// root.AddCommand(&cobra.Command{
	//     Use:  "not-strict",
	//     Run:  func(cmd *cobra.Command, args []string) {},
	// })

	// Show current mode
	if root.IsStrictMode() {
		fmt.Println("✓ Strict mode: ENABLED (only RunE handlers allowed)")
	} else {
		fmt.Println("✓ Strict mode: disabled (Run or RunE allowed)")
	}
	fmt.Println()

	// Run the validate command to verify everything is working
	root.Command().SetArgs([]string{"validate"})
	if err := root.Execute(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
