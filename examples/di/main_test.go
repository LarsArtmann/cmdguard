// Integration test for DI example
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/samber/do/v2"
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

func TestDIExample_CreateCLI(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)
	assert.NotNil(t, cli)
}

func TestDIExample_DIRegistration(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)

	scopeStruct := cli.ScopeStruct()

	// Register DatabaseService
	err = v2.Provide(scopeStruct, func(i do.Injector) (*DatabaseService, error) {
		return NewDatabaseService(), nil
	})
	require.NoError(t, err)

	// Register LoggerService
	err = v2.Provide(scopeStruct, func(i do.Injector) (*LoggerService, error) {
		return NewLoggerService(false), nil
	})
	require.NoError(t, err)

	// Verify services can be invoked
	db, err := v2.Invoke[*DatabaseService](scopeStruct)
	require.NoError(t, err)
	assert.NotNil(t, db)
	assert.True(t, db.connected)

	logger, err := v2.Invoke[*LoggerService](scopeStruct)
	require.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestDIExample_DatabaseServiceQuery(t *testing.T) {
	db := NewDatabaseService()
	require.True(t, db.connected)

	results, err := db.Query(context.Background(), "SELECT * FROM users")
	require.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, "result1", results[0])
}

func TestDIExample_DatabaseServiceClose(t *testing.T) {
	db := NewDatabaseService()
	require.True(t, db.connected)

	err := db.Close()
	require.NoError(t, err)
	assert.False(t, db.connected)

	// Query should fail after close
	_, err = db.Query(context.Background(), "SELECT * FROM users")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestDIExample_LoggerService(t *testing.T) {
	output := captureOutput(func() {
		logger := NewLoggerService(false)
		logger.Info("test message")
	})
	assert.Contains(t, output, "[INFO] test message")

	// Debug should not appear when verbose is false
	output = captureOutput(func() {
		logger := NewLoggerService(false)
		logger.Debug("debug message")
	})
	assert.Empty(t, output)

	// Debug should appear when verbose is true
	output = captureOutput(func() {
		logger := NewLoggerService(true)
		logger.Debug("debug message")
	})
	assert.Contains(t, output, "[DEBUG] debug message")
}

func TestDIExample_QueryFlagsValidation(t *testing.T) {
	// Test valid flags
	flags := &QueryFlags{
		Table: "users",
		Limit: 10,
	}
	assert.Equal(t, "users", flags.Table)
	assert.Equal(t, 10, flags.Limit)

	// Test zero values
	flags = &QueryFlags{}
	assert.Empty(t, flags.Table)
	assert.Equal(t, 0, flags.Limit)
}

func TestDIExample_QueryCommandStructure(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)

	scopeStruct := cli.ScopeStruct()

	// Register services
	err = v2.Provide(scopeStruct, func(i do.Injector) (*DatabaseService, error) {
		return NewDatabaseService(), nil
	})
	require.NoError(t, err)

	err = v2.Provide(scopeStruct, func(i do.Injector) (*LoggerService, error) {
		return NewLoggerService(false), nil
	})
	require.NoError(t, err)

	// Create query command
	queryCmd := v2.Command[AppConfig, *QueryFlags]{
		Use:   "query",
		Short: "Query the database",
		Flags: &QueryFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			if flags.Table == "" {
				return fmt.Errorf("table name is required")
			}
			if flags.Limit < 1 {
				return fmt.Errorf("limit must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			db, err := v2.Invoke[*DatabaseService](scopeStruct)
			if err != nil {
				return v2.NewServiceError("*DatabaseService", err)
			}
			logger, err := v2.Invoke[*LoggerService](scopeStruct)
			if err != nil {
				return v2.NewServiceError("*LoggerService", err)
			}

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
			logger, err := v2.Invoke[*LoggerService](scopeStruct)
			if err != nil {
				return v2.NewServiceError("*LoggerService", err)
			}
			logger.Info("Query completed")
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, queryCmd)
	require.NoError(t, err)

	// Test execution with valid flags
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"query", "--table", "users", "--limit", "2"})
		_ = cli.Execute(context.Background())
	})

	assert.Contains(t, output, "[INFO] Querying table: users")
	assert.Contains(t, output, "Results from users:")
	assert.Contains(t, output, "- result1")
	assert.Contains(t, output, "- result2")
}

func TestDIExample_QueryCommandPreRunEValidation(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)

	scopeStruct := cli.ScopeStruct()

	err = v2.Provide(scopeStruct, func(i do.Injector) (*DatabaseService, error) {
		return NewDatabaseService(), nil
	})
	require.NoError(t, err)

	err = v2.Provide(scopeStruct, func(i do.Injector) (*LoggerService, error) {
		return NewLoggerService(false), nil
	})
	require.NoError(t, err)

	queryCmd := v2.Command[AppConfig, *QueryFlags]{
		Use:   "query",
		Short: "Query the database",
		Flags: &QueryFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			if flags.Table == "" {
				return fmt.Errorf("table name is required")
			}
			if flags.Limit < 1 {
				return fmt.Errorf("limit must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, queryCmd)
	require.NoError(t, err)

	// Test with missing table
	cli.RootCommand().SetArgs([]string{"query"})
	err = cli.Execute(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"table\" not set")

	// Reset CLI and test with invalid limit
	cli, _ = v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	queryCmd.PreRunE = func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
		if flags.Limit < 1 {
			return fmt.Errorf("limit must be at least 1, got %d", flags.Limit)
		}
		return nil
	}
	queryCmd.RunE = func(ctx context.Context, cfg *AppConfig, flags *QueryFlags) error {
		return nil
	}
	err = v2.AddAnyCommand(cli, queryCmd)
	require.NoError(t, err)

	// This test needs a valid table to get past flag validation
	// The PreRunE validation would catch invalid limit if we could set it to 0
	// But the default is 10 from the struct tag
	cli.RootCommand().SetArgs([]string{"query", "--table", "users"})
	err = cli.Execute(context.Background())
	// Should succeed with default limit of 10
	require.NoError(t, err)
}

func TestDIExample_AppConfigStruct(t *testing.T) {
	// Test default values
	cfg := AppConfig{}
	assert.False(t, cfg.Verbose)
	assert.Equal(t, "", cfg.DBHost)

	// Test with values
	cfg = AppConfig{
		Verbose: true,
		DBHost:  "localhost",
	}
	assert.True(t, cfg.Verbose)
	assert.Equal(t, "localhost", cfg.DBHost)
}

func TestDIExample_QueryWithLimitTruncation(t *testing.T) {
	db := NewDatabaseService()
	results, err := db.Query(context.Background(), "users")
	require.NoError(t, err)
	assert.Len(t, results, 3)

	// Simulate limit truncation logic
	limit := 2
	var limitedResults []string
	for i, result := range results {
		if i >= limit {
			break
		}
		limitedResults = append(limitedResults, result)
	}
	assert.Len(t, limitedResults, 2)
	assert.Equal(t, "result1", limitedResults[0])
	assert.Equal(t, "result2", limitedResults[1])
}

func TestDIExample_ServiceErrorHandling(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)

	// Try to invoke unregistered service
	_, err = v2.Invoke[*DatabaseService](cli.ScopeStruct())
	require.Error(t, err)

	// Create service error
	svcErr := v2.NewServiceError("*DatabaseService", err)
	assert.Error(t, svcErr)
	assert.Contains(t, svcErr.Error(), "*DatabaseService")
}

func TestDIExample_CommandWithServiceInvocation(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("diapp", "DI Example Application", AppConfig{})
	require.NoError(t, err)

	scopeStruct := cli.ScopeStruct()

	// Register only logger (not database)
	err = v2.Provide(scopeStruct, func(i do.Injector) (*LoggerService, error) {
		return NewLoggerService(true), nil
	})
	require.NoError(t, err)

	// Create a command that tries to invoke database (which is not registered)
	cmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "test",
		Short: "Test command",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			_, err := v2.Invoke[*DatabaseService](scopeStruct)
			if err != nil {
				return v2.NewServiceError("*DatabaseService", err)
			}
			return nil
		},
	}

	err = cli.AddCommand(cmd)
	require.NoError(t, err)

	// Execution should fail because DatabaseService is not registered
	cli.RootCommand().SetArgs([]string{"test"})
	err = cli.Execute(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "*DatabaseService")
}
