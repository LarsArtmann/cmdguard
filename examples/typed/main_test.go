// Integration test for typed example
package main

import (
	"bytes"
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

func TestTypedExample_CreateCLI(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	require.NoError(t, err)
	assert.NotNil(t, cli)

	cfg := cli.Config()
	assert.NotNil(t, cfg)
}

func TestTypedExample_VersionCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	require.NoError(t, err)

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
	require.NoError(t, err)

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"version"})
		_ = cli.Execute(context.Background())
	})

	assert.Contains(t, output, "myapp version 1.0.0")
}

func TestTypedExample_GreetCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{Verbose: false})
	require.NoError(t, err)

	greetCmd := v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			if flags.Shout {
				msg = stringsToUpper(msg)
			}
			for i := 0; i < flags.Count; i++ {
				fmt.Println(msg)
			}
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, greetCmd)
	require.NoError(t, err)

	// Test basic greeting
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet"})
		_ = cli.Execute(context.Background())
	})
	assert.Contains(t, output, "Hello, World!")
}

func TestTypedExample_GreetCommandWithFlags(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{Verbose: false})
	require.NoError(t, err)

	greetCmd := v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			if flags.Shout {
				msg = stringsToUpper(msg)
			}
			for i := 0; i < flags.Count; i++ {
				fmt.Println(msg)
			}
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, greetCmd)
	require.NoError(t, err)

	// Test with name flag
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--name", "Alice"})
		_ = cli.Execute(context.Background())
	})
	assert.Contains(t, output, "Hello, Alice!")

	// Test with shout flag
	output = captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--name", "Bob", "--shout"})
		_ = cli.Execute(context.Background())
	})
	assert.Contains(t, output, "HELLO, BOB!")

	// Test with count flag
	output = captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--count", "3"})
		_ = cli.Execute(context.Background())
	})
	count := bytes.Count([]byte(output), []byte("Hello, World!"))
	assert.Equal(t, 3, count)
}

func TestTypedExample_ConfigCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{
		Verbose: true,
		Output:  "json",
		APIURL:  "https://api.test.com",
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"config"})
		_ = cli.Execute(context.Background())
	})

	assert.Contains(t, output, "Verbose: true")
	assert.Contains(t, output, "Output: json")
	assert.Contains(t, output, "API URL: https://api.test.com")
}

func TestTypedExample_DIRegistration(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{Verbose: true})
	require.NoError(t, err)

	scope := cli.ScopeStruct()

	// Register config
	err = v2.ProvideValue(scope, AppConfig{Verbose: true})
	require.NoError(t, err)

	// Register logger
	err = v2.Provide(scope, func(i do.Injector) (*Logger, error) {
		cfg, err := v2.Invoke[*AppConfig](scope)
		if err != nil {
			return nil, err
		}
		return &Logger{verbose: cfg.Verbose}, nil
	})
	require.NoError(t, err)

	// Verify we can invoke the logger
	logger, err := v2.Invoke[*Logger](scope)
	require.NoError(t, err)
	assert.NotNil(t, logger)
	assert.True(t, logger.verbose)
}

func TestTypedExample_DatabaseService(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	require.NoError(t, err)

	scope := cli.ScopeStruct()

	// Register database
	err = v2.ProvideValue(scope, &Database{connectionString: "postgres://test:5432/db"})
	require.NoError(t, err)

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
	require.NoError(t, err)

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"db-status"})
		_ = cli.Execute(context.Background())
	})

	assert.Contains(t, output, "postgres://test:5432/db")
}

func TestTypedExample_PreRunEValidation(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	require.NoError(t, err)

	greetCmd := v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			if flags.Count < 1 {
				return fmt.Errorf("count must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			fmt.Println("Greeting executed")
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, greetCmd)
	require.NoError(t, err)

	// Test with invalid count
	cli.RootCommand().SetArgs([]string{"greet", "--count", "0"})
	err = cli.Execute(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "count must be at least 1")

	// Reset and test with valid count
	cli, _ = v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	greetCmd.RunE = func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
		fmt.Println("Greeting executed")
		return nil
	}
	err = v2.AddAnyCommand(cli, greetCmd)
	require.NoError(t, err)

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet", "--count", "1"})
		_ = cli.Execute(context.Background())
	})
	assert.Contains(t, output, "Greeting executed")
}
