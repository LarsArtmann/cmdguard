// Example: Middleware Patterns with PreRunE and PostRunE
//
// This example demonstrates how to use lifecycle hooks for:
// - Validation before command execution (PreRunE)
// - Cleanup after command execution (PostRunE)
// - Composable middleware chains
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application configuration
type AppConfig struct {
	Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
}

// ProcessingFlags for the process command
type ProcessingFlags struct {
	Input  string `flag:"input" short:"i" required:"true" help:"Input file"`
	Output string `flag:"output" short:"o" default:"output.txt" help:"Output file"`
	DryRun bool   `flag:"dry-run" default:"false" help:"Simulate processing"`
}

// Timer is a simple timing middleware
type Timer struct {
	start time.Time
	label string
}

func NewTimer(label string) *Timer {
	return &Timer{label: label}
}

func (t *Timer) Start() {
	t.start = time.Now()
	fmt.Printf("[%s] Starting...\n", t.label)
}

func (t *Timer) Stop() {
	duration := time.Since(t.start)
	fmt.Printf("[%s] Completed in %v\n", t.label, duration)
}

// Validator is a validation middleware
func validateFlags(flags *ProcessingFlags) error {
	if flags.Input == "" {
		return fmt.Errorf("input file is required")
	}
	if flags.Output == "" {
		return fmt.Errorf("output file is required")
	}
	if flags.Input == flags.Output {
		return fmt.Errorf("input and output cannot be the same file")
	}
	return nil
}

// Authenticator is an auth middleware simulation
func authenticate(ctx context.Context) error {
	// In real app, check auth tokens, etc.
	fmt.Println("[Auth] Checking permissions...")
	return nil
}

// Logger is a logging middleware
func logStart(flags *ProcessingFlags) {
	fmt.Printf("[Log] Processing %s -> %s (dry-run: %v)\n",
		flags.Input, flags.Output, flags.DryRun)
}

func logComplete(flags *ProcessingFlags) {
	fmt.Printf("[Log] Processing complete for %s\n", flags.Output)
}

// Cleanup is a cleanup middleware
func cleanup() error {
	fmt.Println("[Cleanup] Releasing resources...")
	return nil
}

// Metrics records execution metrics
func recordMetrics(duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	fmt.Printf("[Metrics] Duration: %v, Status: %s\n", duration, status)
}

func main() {
	// Create CLI
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
		os.Exit(1)
	}

	// Create timer middleware
	timer := NewTimer("process")

	// Command with full middleware chain
	processCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "process",
		Short: "Process a file with middleware",
		Long:  "Demonstrates PreRunE and PostRunE lifecycle hooks",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			// Chain of PreRunE middleware

			// 1. Authentication
			if err := authenticate(ctx); err != nil {
				return fmt.Errorf("auth failed: %w", err)
			}

			// 2. Validation
			if err := validateFlags(flags); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// 3. Start timer
			timer.Start()

			// 4. Logging
			logStart(flags)

			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			// Main processing logic
			if flags.DryRun {
				fmt.Println("[Dry Run] Would process file:")
				fmt.Printf("  Input:  %s\n", flags.Input)
				fmt.Printf("  Output: %s\n", flags.Output)
				return nil
			}

			// Simulate processing
			fmt.Println("[Processing] Reading input...")
			fmt.Println("[Processing] Transforming data...")
			fmt.Println("[Processing] Writing output...")

			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			// Chain of PostRunE middleware
			// Note: PostRunE runs even if RunE fails

			// 1. Stop timer and record metrics
			duration := time.Since(timer.start)
			timer.Stop()

			// 2. Logging
			logComplete(flags)

			// 3. Cleanup
			if err := cleanup(); err != nil {
				fmt.Fprintf(os.Stderr, "[Cleanup Error] %v\n", err)
			}

			// 4. Metrics
			recordMetrics(duration, true)

			return nil
		},
	}

	if err := v2.AddAnyCommand(cli, processCmd); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add process command: %v\n", err)
		os.Exit(1)
	}

	// Simple command demonstrating just PreRunE
	validateCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "validate",
		Short: "Validate configuration without processing",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			fmt.Println("[Validate] Running validation...")
			if err := validateFlags(flags); err != nil {
				return err
			}
			fmt.Println("[Validate] ✓ Configuration is valid")
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			fmt.Println("[Validate] Validation complete, exiting without processing")
			return nil
		},
	}

	if err := v2.AddAnyCommand(cli, validateCmd); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add validate command: %v\n", err)
		os.Exit(1)
	}

	// Execute
	cli.ExecuteAndExit(context.Background())
}
