// Package main demonstrates cmdguard v2 with dependency injection.
//
// This example shows:
// - Creating services with samber/do/v2
// - Implementing Shutdowner for cleanup
// - Implementing HealthcheckerWithContext for health checks
// - Graceful shutdown handling
//
// Usage:
//
//	go run examples/di/main.go
//	go run examples/di/main.go check
//	go run examples/di/main.db shutdown
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// Config is the application configuration.
type Config struct {
	LogLevel  v2.LogLevel `default:"info"                  flag:"log-level"  help:"Log level"  short:"l"`
	ServerURL string      `default:"http://localhost:8080" flag:"server-url" help:"Server URL"`
}

// DatabaseService simulates a database connection.
type DatabaseService struct {
	connected bool
	url       string
}

// Verify interface implementations at compile time.
var (
	_ do.ShutdownerWithError      = (*DatabaseService)(nil)
	_ do.HealthcheckerWithContext = (*DatabaseService)(nil)
)

// NewDatabaseService creates a new database service.
func NewDatabaseService(i do.Injector) (*DatabaseService, error) {
	cfg := do.MustInvoke[*Config](i)

	return &DatabaseService{
		connected: true,
		url:       cfg.ServerURL,
	}, nil
}

// Shutdown implements the Shutdowner interface.
func (d *DatabaseService) Shutdown() error {
	fmt.Println("Database: shutting down...")

	d.connected = false

	fmt.Println("Database: disconnected")

	return nil
}

// HealthCheck implements the HealthcheckerWithContext interface.
func (d *DatabaseService) HealthCheck(ctx context.Context) error {
	if !d.connected {
		return errors.New("database not connected")
	}

	return nil
}

// IsConnected returns the connection status.
func (d *DatabaseService) IsConnected() bool {
	return d.connected
}

// APIService simulates an HTTP API client.
type APIService struct {
	client *DatabaseService
}

// NewAPIService creates a new API service.
func NewAPIService(i do.Injector) (*APIService, error) {
	db := do.MustInvoke[*DatabaseService](i)

	return &APIService{client: db}, nil
}

// Call simulates an API call.
func (a *APIService) Call(ctx context.Context) error {
	_ = ctx // context required by interface but not used in this example

	if !a.client.IsConnected() {
		return errors.New("database not available")
	}

	fmt.Println("API: calling database...")

	return nil
}

// execute runs the CLI and exits on error.
func execute(ctx context.Context, cli *v2.CLI[Config]) {
	if err := cli.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// fatal prints the error to stderr and exits with code 1.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	root, err := v2.NewCLI[Config]("di-app", "DI Example App", Config{})
	if err != nil {
		fatal("Error: %v\n", err)
	}

	// Register services in DI scope
	if err := v2.Provide(root.Scope(), NewDatabaseService); err != nil {
		fatal("Error registering database: %v\n", err)
	}

	if err := v2.Provide(root.Scope(), NewAPIService); err != nil {
		fatal("Error registering API: %v\n", err)
	}

	// Add health check command
	if err := v2.AddCommand(root, v2.Command[Config, v2.NoFlags]{
		Use:   "check",
		Short: "Run health checks",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			fmt.Println("Running health checks...")

			err := root.Scope().HealthCheckWithContext(ctx)
			if err != nil {
				fmt.Printf("Health check FAILED: %v\n", err)

				return fmt.Errorf("health check failed: %w", err)
			}

		fmt.Println("All health checks PASSED!")

			return nil
		},
	}); err != nil {
		fatal("Error adding check command: %v\n", err)
	}

	// Add API call command
	if err := v2.AddCommand(root, v2.Command[Config, v2.NoFlags]{
		Use:   "call",
		Short: "Call the API",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			api, err := v2.Invoke[*APIService](root.Scope())
			if err != nil {
				return err
			}

			return api.Call(ctx)
		},
	}); err != nil {
		fatal("Error adding call command: %v\n", err)
	}

	// Add shutdown command
	if err := v2.AddCommand(root, v2.Command[Config, v2.NoFlags]{
		Use:   "shutdown",
		Short: "Shutdown services gracefully",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			fmt.Println("Shutting down services...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := root.Shutdown(shutdownCtx)
			if err != nil {
				fmt.Printf("Shutdown error: %v\n", err)

				return fmt.Errorf("shutdown failed: %w", err)
			}

			fmt.Println("Shutdown complete!")

			return nil
		},
	}); err != nil {
		fatal("Error adding shutdown command: %v\n", err)
	}

	// Run health check before starting
	if err := root.Scope().HealthCheckWithContext(ctx); err != nil {
		fatal("Initial health check failed: %v\n", err)
	}

	// Execute
	execute(ctx, root)

	// Graceful shutdown on exit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := root.Shutdown(shutdownCtx); err != nil {
		fatal("Shutdown error: %v\n", err)
	}
}
