// Integration test for typed example
package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// greetRunE returns a RunE function that prints "Greeting executed".
func greetRunE() func(context.Context, *AppConfig, *GreetFlags) error {
	return func(_ context.Context, _ *AppConfig, _ *GreetFlags) error {
		fmt.Println("Greeting executed")

		return nil
	}
}

// greetValidationPreRunE returns a PreRunE that validates count >= 1.
func greetValidationPreRunE() func(context.Context, *AppConfig, *GreetFlags) error {
	return func(_ context.Context, _ *AppConfig, flags *GreetFlags) error {
		if flags.Count < 1 {
			return errors.New("count should be at least 1")
		}

		return nil
	}
}

// newGreetCommandWithValidation creates a greet command with PreRunE validation.
func newGreetCommandWithValidation() (v2.Command[AppConfig, *GreetFlags], error) {
	return v2.NewCommand[AppConfig, *GreetFlags](
		"greet",
		greetRunE(),
		v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
		v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
		v2.WithPreRunE[AppConfig, *GreetFlags](greetValidationPreRunE()),
	)
}

//nolint:paralleltest // captures os.Stdout, not safe for parallel execution
func TestTypedExample_GreetCommandWithFlags(t *testing.T) {
	cli, err := v2.NewCLI[AppConfig](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd := newGreetCmd()

	err = v2.AddCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test with name flag
	output := runCLIWithArgs(cli, "greet", "--name", "Alice")
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
	cli, _ = v2.NewCLI[AppConfig](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	greetCmd = v2.MustNewCommand[AppConfig, *GreetFlags](
		"greet",
		func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			for range flags.Count {
				fmt.Println(msg)
			}

			return nil
		},
		v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
		v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
	)
	_ = v2.AddCommand(cli, greetCmd)

	output = runCLIWithArgs(cli, "greet", "--count", "3")
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

//nolint:paralleltest // captures os.Stdout, not safe for parallel execution
func TestTypedExample_ConfigCommand(t *testing.T) {
	// Use default config values for this test
	cli, err := v2.NewCLI[AppConfig]("myapp", "A typed CLI application", AppConfig{
		Verbose: false,
		Output:  "text",
		APIURL:  "https://api.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"config",
		func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Printf("Verbose: %v\n", cfg.Verbose)
			fmt.Printf("Output: %s\n", cfg.Output)
			fmt.Printf("API URL: %s\n", cfg.APIURL)

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Show current configuration"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, configCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := runCLIWithArgs(cli, "config")

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

//nolint:paralleltest // captures os.Stdout, not safe for parallel execution
func TestTypedExample_DIRegistration(t *testing.T) {
	cli, err := v2.NewCLI[AppConfig](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: true},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scope := cli.Scope()

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

//nolint:paralleltest // captures os.Stdout, not safe for parallel execution
func TestTypedExample_DatabaseService(t *testing.T) {
	cli, err := v2.NewCLI[AppConfig]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	scope := cli.Scope()

	// Register database
	err = v2.ProvideValue(scope, &Database{connectionString: "postgres://test:5432/db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test command that uses database
	dbCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"db-status",
		func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			db, err := v2.Invoke[*Database](scope)
			if err != nil {
				return err
			}

			fmt.Printf("Database: %s\n", db.connectionString)

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Check database status"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, dbCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := runCLIWithArgs(cli, "db-status")

	if !strings.Contains(output, "postgres://test:5432/db") {
		t.Errorf("output should contain %q, got %q", "postgres://test:5432/db", output)
	}
}

//nolint:paralleltest // captures os.Stdout, not safe for parallel execution
func TestTypedExample_PreRunEValidation(t *testing.T) {
	// Test 1: Invalid count should return error
	t.Run("invalid count", func(t *testing.T) {
		cli, err := v2.NewCLI[AppConfig]("myapp", "A typed CLI application", AppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		greetCmd, err := newGreetCommandWithValidation()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = v2.AddCommand(cli, greetCmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cli.RootCommand().SetArgs([]string{"greet", "--count", "0"})

		err = cli.Execute(context.Background())
		if err == nil {
			t.Fatal("expected error for count < 1")
		}

		if !strings.Contains(err.Error(), "count should be at least 1") {
			t.Errorf("error should contain %q, got %q", "count should be at least 1", err.Error())
		}
	})

	// Test 2: Valid count should succeed
	t.Run("valid count", func(t *testing.T) {
		cli, err := v2.NewCLI[AppConfig]("myapp", "A typed CLI application", AppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		greetCmd, err := newGreetCommandWithValidation()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = v2.AddCommand(cli, greetCmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := runCLIWithArgs(cli, "greet", "--count", "1")
		if !strings.Contains(output, "Greeting executed") {
			t.Errorf("output should contain %q, got %q", "Greeting executed", output)
		}
	})
}
