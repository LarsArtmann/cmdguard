// Example: Dependency Injection Patterns
//
// This example demonstrates how to use cmdguard v2's DI integration
// with samber/do/v2 for service management.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/samber/do/v2"
)

// DatabaseService is a mock database service
type DatabaseService struct {
	connected bool
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{connected: true}
}

func (db *DatabaseService) Query(ctx context.Context, query string) ([]string, error) {
	if !db.connected {
		return nil, fmt.Errorf("database not connected")
	}
	// Mock implementation
	return []string{"result1", "result2", "result3"}, nil
}

func (db *DatabaseService) Close() error {
	db.connected = false
	return nil
}

// LoggerService is a mock logger service
type LoggerService struct {
	verbose bool
}

func NewLoggerService(verbose bool) *LoggerService {
	return &LoggerService{verbose: verbose}
}

func (l *LoggerService) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (l *LoggerService) Debug(msg string) {
	if l.verbose {
		fmt.Printf("[DEBUG] %s\n", msg)
	}
}

// AppConfig is the application configuration
type AppConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose logging"`
	DBHost  string `flag:"db-host" default:"localhost" help:"Database host"`
}

type QueryFlags struct {
	Table string `flag:"table" short:"t" required:"true" help:"Table to query"`
	Limit int    `flag:"limit" short:"l" default:"10" help:"Result limit"`
}

func main() {
	// Create CLI with typed config
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
		os.Exit(1)
	}

	// Get the DI scope struct for service registration
	scopeStruct := cli.ScopeStruct()

	// Register services in DI container
	// These services will be available to all commands
	v2.Provide(scopeStruct, func(i do.Injector) (*DatabaseService, error) {
		// In real app, use config to configure connection
		return NewDatabaseService(), nil
	})

	v2.Provide(scopeStruct, func(i do.Injector) (*LoggerService, error) {
		// Access config through injector if needed
		return NewLoggerService(false), nil
	})

	// Add commands
	queryCmd := v2.Command[AppConfig, *QueryFlags]{
		Use:   "query",
		Short: "Query the database",
		Flags: &QueryFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			// Validation before execution
			if flags.Table == "" {
				return fmt.Errorf("table name is required")
			}
			if flags.Limit < 1 {
				return fmt.Errorf("limit must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			// Invoke services from DI container
			db := do.MustInvoke[*DatabaseService](cli.Scope())
			logger := do.MustInvoke[*LoggerService](cli.Scope())

			logger.Info(fmt.Sprintf("Querying table: %s (limit: %d)", flags.Table, flags.Limit))

			results, err := db.Query(ctx, flags.Table)
			if err != nil {
				return fmt.Errorf("query failed: %w", err)
			}

			fmt.Printf("Results from %s:\n", flags.Table)
			for i, result := range results {
				if i >= flags.Limit {
					break
				}
				fmt.Printf("  - %s\n", result)
			}

			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			// Cleanup after execution
			logger := do.MustInvoke[*LoggerService](cli.Scope())
			logger.Info("Query completed")
			return nil
		},
	}

	if err := v2.AddAnyCommand(cli, queryCmd); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add query command: %v\n", err)
		os.Exit(1)
	}

	// Execute
	ctx := context.Background()
	defer cli.Shutdown(ctx) // Cleanup services on exit
	cli.ExecuteAndExit(ctx)
}
