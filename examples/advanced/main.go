// Advanced example demonstrating nested commands and flags.
package main

import (
	"context"
	"fmt"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
)

func main() {
	root := cmdguard.New("advanced", "An advanced CLI example with nested commands")

	// Database command group
	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}

	// db migrate
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, _ := cmd.Flags().GetString("version")
			if version != "" {
				fmt.Printf("Migrating to version: %s\n", version)
			} else {
				fmt.Println("Running all pending migrations...")
			}
			return nil
		},
	}
	migrateCmd.Flags().StringP("version", "v", "", "Target migration version")

	// db seed
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, _ := cmd.Flags().GetInt("count")
			fmt.Printf("Seeding database with %d records...\n", count)
			return nil
		},
	}
	seedCmd.Flags().IntP("count", "n", 10, "Number of records to seed")

	// Add subcommands
	dbCmd.AddCommand(migrateCmd)
	dbCmd.AddCommand(seedCmd)
	root.AddCommand(dbCmd)

	// User command group
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management",
	}

	userCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing users...")
		},
	})

	userCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("username required")
			}
			fmt.Printf("Creating user: %s\n", args[0])
			return nil
		},
	})

	root.AddCommand(userCmd)

	// Execute
	root.ExecuteAndExit(context.Background())
}
