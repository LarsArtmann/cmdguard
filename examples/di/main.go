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
	"fmt"
	"os"
	"time"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/samber/do/v2"
)

// Config is the application configuration.
type Config struct {
	LogLevel  v2.LogLevel `flag:"log-level" short:"l" default:"info" help:"Log level"`
	ServerURL string      `flag:"server-url" default:"http://localhost:8080" help:"Server URL"`
}

// DatabaseService simulates a database connection.
type DatabaseService struct {
	connected bool
	url       string
}

// Verify interface implementations at compile time.
var (
	_ do.Shutdowner               = (*DatabaseService)(nil)
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
func (d *DatabaseService) Shutdown(ctx context.Context) error {
	fmt.Println("Database: shutting down...")
	d.connected = false
	fmt.Println("Database: disconnected")
	return nil
}

// HealthCheck implements the HealthcheckerWithContext interface.
func (d *DatabaseService) HealthCheck(ctx context.Context) error {
	if !d.connected {
		return fmt.Errorf("database not connected")
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
	if !a.client.IsConnected() {
		return fmt.Errorf("database not available")
	}
	fmt.Println("API: calling database...")
	return nil
}

func main() {
	ctx := context.Background()

	// Create CLI with typed config
	root, err := v2.New[Config, v2.NoFlags]("di-app", "DI Example App", Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Register services in DI scope
	if err := v2.Provide(root.ScopeStruct(), NewDatabaseService); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering database: %v\n", err)
		os.Exit(1)
	}

	if err := v2.Provide(root.ScopeStruct(), NewAPIService); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering API: %v\n", err)
		os.Exit(1)
	}

	// Add health check command
	if err := root.AddCommand(v2.Command[Config, v2.NoFlags]{
		Use:   "check",
		Short: "Run health checks",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			fmt.Println("Running health checks...")

			if err := root.HealthCheckWithContext(ctx); err != nil {
				fmt.Printf("Health check FAILED: %v\n", err)
				return err
			}

			fmt.Println("All health checks PASSED!")
			return nil
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding check command: %v\n", err)
		os.Exit(1)
	}

	// Add API call command
	if err := root.AddCommand(v2.Command[Config, v2.NoFlags]{
		Use:   "call",
		Short: "Call the API",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			api, err := v2.Invoke[*APIService](root.ScopeStruct())
			if err != nil {
				return err
			}
			return api.Call(ctx)
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding call command: %v\n", err)
		os.Exit(1)
	}

	// Add shutdown command
	if err := root.AddCommand(v2.Command[Config, v2.NoFlags]{
		Use:   "shutdown",
		Short: "Shutdown services gracefully",
		RunE: func(ctx context.Context, cfg *Config, _ v2.NoFlags) error {
			fmt.Println("Shutting down services...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := root.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("Shutdown error: %v\n", err)
				return err
			}

			fmt.Println("Shutdown complete!")
			return nil
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error adding shutdown command: %v\n", err)
		os.Exit(1)
	}

	// Run health check before starting
	if err := root.HealthCheckWithContext(ctx); err != nil {
		fmt.Printf("Initial health check failed: %v\n", err)
		os.Exit(1)
	}

	// Execute
	if err := root.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown on exit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := root.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
		os.Exit(1)
	}
}
