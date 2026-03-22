// Integration test for typed example
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// newGreetCmd creates a new greet command instance for testing.
// This helper ensures consistency across tests and avoids code duplication.
func newGreetCmd() v2.Command[AppConfig, *GreetFlags] {
	return v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			if flags.Shout {
				msg = strings.ToUpper(msg)
			}

			for range flags.Count {
				fmt.Println(msg)
			}

			return nil
		},
	}
}

// captureOutput captures stdout during the execution of f and returns it as a string.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, r) // Error intentionally ignored in test helper
	os.Stdout = old

	return buf.String()
}

func TestTypedExample_CreateCLI(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cli == nil {
		t.Fatal("cli is nil")
	}

	cfg := cli.Config()
	if cfg == nil {
		t.Fatal("config is nil")
	}
}

func TestTypedExample_VersionCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cli.SetVersion("1.0.0")

	versionCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "version",
		Short: "Print version information",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("myapp version 1.0.0")

			return nil
		},
	}

	err = cli.AddCommand(versionCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"version"})
		_ = cli.Execute(context.Background())
	})

	if !strings.Contains(output, "myapp version 1.0.0") {
		t.Errorf("output should contain %q, got %q", "myapp version 1.0.0", output)
	}
}

func TestTypedExample_GreetCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd := newGreetCmd()

	err = v2.AddAnyCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test basic greeting
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet"})
		_ = cli.Execute(context.Background())
	})
	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("output should contain %q, got %q", "Hello, World!", output)
	}
}

func TestTypedExample_GreetCommandWithFlags(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd := newGreetCmd()

	err = v2.AddAnyCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test with name flag
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--name", "Alice"})
		_ = cli.Execute(context.Background())
	})
	if !strings.Contains(output, "Hello, Alice!") {
		t.Errorf("output should contain %q, got %q", "Hello, Alice!", output)
	}

	// Test with shout flag
	output = captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--name", "Bob", "--shout"})
		_ = cli.Execute(context.Background())
	})
	if !strings.Contains(output, "HELLO, BOB!") {
		t.Errorf("output should contain %q, got %q", "HELLO, BOB!", output)
	}

	// Test with count flag - recreate CLI to avoid flag pollution
	cli, _ = v2.New[AppConfig, v2.NoFlags](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	greetCmd = v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			for range flags.Count {
				fmt.Println(msg)
			}

			return nil
		},
	}
	_ = v2.AddAnyCommand(cli, greetCmd)

	output = captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--count", "3"})
		_ = cli.Execute(context.Background())
	})
	// With count=3, we should see "Hello, World!" three times
	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("output should contain %q, got %q", "Hello, World!", output)
	}
	// Count the occurrences
	occurrences := 0

	for i := 0; i <= len(output)-len("Hello, World!"); i++ {
		if output[i:i+len("Hello, World!")] == "Hello, World!" {
			occurrences++
		}
	}

	if occurrences != 3 {
		t.Errorf("expected 3 occurrences of 'Hello, World!', got %d", occurrences)
	}
}

func TestTypedExample_ConfigCommand(t *testing.T) {
	// Use default config values for this test
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{
		Verbose: false,
		Output:  "text",
		APIURL:  "https://api.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "config",
		Short: "Show current configuration",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Printf("Verbose: %v\n", cfg.Verbose)
			fmt.Printf("Output: %s\n", cfg.Output)
			fmt.Printf("API URL: %s\n", cfg.APIURL)

			return nil
		},
	}

	err = cli.AddCommand(configCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"config"})
		_ = cli.Execute(context.Background())
	})

	// Verify the default config values are displayed
	if !strings.Contains(output, "Verbose: false") {
		t.Errorf("output should contain %q, got %q", "Verbose: false", output)
	}

	if !strings.Contains(output, "Output: text") {
		t.Errorf("output should contain %q, got %q", "Output: text", output)
	}

	if !strings.Contains(output, "API URL: https://api.example.com") {
		t.Errorf("output should contain %q, got %q", "API URL: https://api.example.com", output)
	}
}

func TestTypedExample_DIRegistration(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scope := cli.ScopeStruct()

	// Register config
	err = v2.ProvideValue(scope, AppConfig{Verbose: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Register logger
	err = v2.Provide(scope, func(i do.Injector) (*Logger, error) {
		cfg, err := v2.Invoke[*AppConfig](scope)
		if err != nil {
			return nil, err
		}

		return &Logger{verbose: cfg.Verbose}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we can invoke the logger
	logger, err := v2.Invoke[*Logger](scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if logger == nil {
		t.Fatal("logger is nil")
	}

	if !logger.verbose {
		t.Error("logger.verbose should be true")
	}
}

func TestTypedExample_DatabaseService(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scope := cli.ScopeStruct()

	// Register database
	err = v2.ProvideValue(scope, &Database{connectionString: "postgres://test:5432/db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test command that uses database
	dbCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "db-status",
		Short: "Check database status",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			db, err := v2.Invoke[*Database](scope)
			if err != nil {
				return err
			}

			fmt.Printf("Database: %s\n", db.connectionString)

			return nil
		},
	}

	err = cli.AddCommand(dbCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"db-status"})
		_ = cli.Execute(context.Background())
	})

	if !strings.Contains(output, "postgres://test:5432/db") {
		t.Errorf("output should contain %q, got %q", "postgres://test:5432/db", output)
	}
}

func TestTypedExample_PreRunEValidation(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd := v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			if flags.Count < 1 {
				return errors.New("count should be at least 1")
			}

			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			fmt.Println("Greeting executed")

			return nil
		},
	}

	err = v2.AddAnyCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test with invalid count
	cli.RootCommand().SetArgs([]string{"greet", "--count", "0"})

	err = cli.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error for count < 1")
	}

	if !strings.Contains(err.Error(), "count should be at least 1") {
		t.Errorf("error should contain %q, got %q", "count should be at least 1", err.Error())
	}

	// Reset and test with valid count
	cli, _ = v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	greetCmd.RunE = func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
		fmt.Println("Greeting executed")

		return nil
	}

	err = v2.AddAnyCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--count", "1"})
		_ = cli.Execute(context.Background())
	})
	if !strings.Contains(output, "Greeting executed") {
		t.Errorf("output should contain %q, got %q", "Greeting executed", output)
	}
}
