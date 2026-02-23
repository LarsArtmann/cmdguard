// Integration test for advanced example
package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvancedExample_DatabaseMigrate(t *testing.T) {
	root := cmdguard.New("advanced", "An advanced CLI example with nested commands")

	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}

	var output bytes.Buffer
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, _ := cmd.Flags().GetString("version")
			if version != "" {
				output.WriteString(fmt.Sprintf("Migrating to version: %s\n", version))
			} else {
				output.WriteString("Running all pending migrations...\n")
			}
			return nil
		},
	}
	migrateCmd.Flags().StringP("version", "v", "", "Target migration version")

	dbCmd.AddCommand(migrateCmd)
	root.AddCommand(dbCmd)

	// Test migrate without version
	root.Command().SetArgs([]string{"db", "migrate"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Running all pending migrations...\n", output.String())

	// Reset and test with version
	output.Reset()
	root.Command().SetArgs([]string{"db", "migrate", "--version", "1.2.3"})
	err = root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Migrating to version: 1.2.3\n", output.String())
}

func TestAdvancedExample_DatabaseSeed(t *testing.T) {
	root := cmdguard.New("advanced", "An advanced CLI example with nested commands")

	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}

	var output bytes.Buffer
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			count, _ := cmd.Flags().GetInt("count")
			output.WriteString(fmt.Sprintf("Seeding database with %d records...\n", count))
			return nil
		},
	}
	seedCmd.Flags().IntP("count", "n", 10, "Number of records to seed")

	dbCmd.AddCommand(seedCmd)
	root.AddCommand(dbCmd)

	// Test seed with default count
	root.Command().SetArgs([]string{"db", "seed"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Seeding database with 10 records...\n", output.String())

	// Reset and test with custom count
	output.Reset()
	root.Command().SetArgs([]string{"db", "seed", "--count", "25"})
	err = root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Seeding database with 25 records...\n", output.String())
}

func TestAdvancedExample_UserCommands(t *testing.T) {
	root := cmdguard.New("advanced", "An advanced CLI example with nested commands")

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management",
	}

	var output bytes.Buffer
	userCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			output.WriteString("Listing users...\n")
		},
	})

	userCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("username required")
			}
			output.WriteString(fmt.Sprintf("Creating user: %s\n", args[0]))
			return nil
		},
	})

	root.AddCommand(userCmd)

	// Test list
	root.Command().SetArgs([]string{"user", "list"})
	err := root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Listing users...\n", output.String())

	// Test create with argument
	output.Reset()
	root.Command().SetArgs([]string{"user", "create", "john"})
	err = root.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "Creating user: john\n", output.String())

	// Test create without argument (should error)
	output.Reset()
	root.Command().SetArgs([]string{"user", "create"})
	err = root.Execute(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "username required")
}

func TestAdvancedExample_CommandStructure(t *testing.T) {
	root := cmdguard.New("advanced", "An advanced CLI example with nested commands")

	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database operations",
	}
	dbCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management",
	}
	userCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run:   func(cmd *cobra.Command, args []string) {},
	})

	root.AddCommand(dbCmd)
	root.AddCommand(userCmd)

	cmd := root.Command()
	assert.Len(t, cmd.Commands(), 2)

	dbSubcommands := dbCmd.Commands()
	assert.Len(t, dbSubcommands, 1)
	assert.Equal(t, "migrate", dbSubcommands[0].Use)

	userSubcommands := userCmd.Commands()
	assert.Len(t, userSubcommands, 1)
	assert.Equal(t, "list", userSubcommands[0].Use)
}
