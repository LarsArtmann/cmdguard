package v2_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// Example_basic demonstrates creating a minimal CLI application.
func Example_basic() {
	// Define your application config
	type AppConfig struct {
		Verbose bool `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	}

	// Create the CLI
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A simple CLI application", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	// Add a command without flags
	err = cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
		Use:   "hello",
		Short: "Say hello",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			if cfg.Verbose {
				fmt.Println("Verbose mode enabled")
			}

			fmt.Println("Hello, World!")

			return nil
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	// Execute with arguments (useful for testing)
	_ = cli.ExecuteWithArgs(context.Background(), []string{"hello"})
	// Output:
	// Hello, World!
}

// Example_withFlags demonstrates commands with typed flags.
func Example_withFlags() {
	type AppConfig struct {
		Debug bool `default:"false" flag:"debug" help:"Enable debug mode" short:"d"`
	}

	// Define command-specific flags
	type GreetFlags struct {
		Name  string `default:"World" flag:"name"  help:"Name to greet"       short:"n"`
		Shout bool   `default:"false" flag:"shout" help:"Shout the greeting"  short:"s"`
		Count int    `default:"1"     flag:"count" help:"Number of greetings" short:"c"`
	}

	cli, err := v2.New[AppConfig, *GreetFlags]("greeter", "A greeting CLI", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	err = cli.AddCommand(v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{}, // Initialize with defaults
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			for range flags.Count {
				msg := fmt.Sprintf("Hello, %s!", flags.Name)
				if flags.Shout {
					msg = strings.ToUpper(msg)
				}

				fmt.Println(msg)
			}

			return nil
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	_ = cli.ExecuteWithArgs(
		context.Background(),
		[]string{"greet", "--name", "Alice", "--count", "2"},
	)
	// Output:
	// Hello, Alice!
	// Hello, Alice!
}

// Example_withSubcommands demonstrates nested command hierarchies.
func Example_withSubcommands() {
	type AppConfig struct {
		LogLevel v2.LogLevel `default:"info" flag:"log-level" help:"Log level"`
	}

	cli, err := v2.New[AppConfig, v2.NoFlags]("git", "Git version control", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	// Add nested commands
	err = cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
		Use:   "remote",
		Short: "Manage remote repositories",
		Commands: []v2.Command[AppConfig, v2.NoFlags]{
			{
				Use:   "list",
				Short: "List remotes",
				RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
					fmt.Println("origin")
					fmt.Println("upstream")

					return nil
				},
			},
			{
				Use:   "add",
				Short: "Add a remote",
				RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
					fmt.Println("Remote added")

					return nil
				},
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	_ = cli.ExecuteWithArgs(context.Background(), []string{"remote", "list"})
	// Output:
	// origin
	// upstream
}

// Example_withEnum demonstrates using enum types for restricted values.
func Example_withEnum() {
	type AppConfig struct {
		Environment v2.Enum `default:"dev" flag:"env" help:"Target environment" values:"dev,staging,prod"`
	}

	type DeployFlags struct {
		Version string `flag:"version" help:"Version to deploy" required:"true" short:"v"`
	}

	// Parse enum value (error handling omitted for brevity in this example)
	env, _ := v2.ParseEnum("dev", []string{"dev", "staging", "prod"})

	cli, err := v2.New[AppConfig, *DeployFlags]("deployer", "Deployment CLI", AppConfig{
		Environment: env,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	err = cli.AddCommand(v2.Command[AppConfig, *DeployFlags]{
		Use:   "deploy",
		Short: "Deploy to environment",
		Flags: &DeployFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *DeployFlags) error {
			fmt.Printf("Deploying version %s to %s\n", flags.Version, cfg.Environment.String())

			return nil
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	_ = cli.ExecuteWithArgs(context.Background(), []string{"deploy", "--version", "1.2.3"})
	// Output:
	// Deploying version 1.2.3 to dev
}

// Example_withPreRunE demonstrates validation using PreRunE.
func Example_withPreRunE() {
	type AppConfig struct{}

	type CreateUserFlags struct {
		Email    string `flag:"email"    help:"User email"             required:"true"`
		Password string `flag:"password" help:"User password"          required:"true"`
		Admin    bool   `flag:"admin"    help:"Grant admin privileges"                 default:"false"`
	}

	cli, err := v2.New[AppConfig, *CreateUserFlags]("userctl", "User management CLI", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	err = cli.AddCommand(v2.Command[AppConfig, *CreateUserFlags]{
		Use:   "create",
		Short: "Create a new user",
		Flags: &CreateUserFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *CreateUserFlags) error {
			// Custom validation
			if len(flags.Password) < 8 {
				return errors.New("password must be at least 8 characters")
			}

			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *CreateUserFlags) error {
			role := "user"
			if flags.Admin {
				role = "admin"
			}

			fmt.Printf("Created %s: %s\n", role, flags.Email)

			return nil
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	_ = cli.ExecuteWithArgs(context.Background(), []string{
		"create",
		"--email", "user@example.com",
		"--password", "securepass123",
		"--admin",
	})
	// Output:
	// Created admin: user@example.com
}

// Example_withFunctionalOptions demonstrates using functional options.
func Example_withFunctionalOptions() {
	type AppConfig struct{}

	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My application", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	// Create command using functional options
	cmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"version",
		v2.WithShort[AppConfig, v2.NoFlags]("Show version info"),
		v2.WithRunE(func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("myapp version 1.0.0")

			return nil
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	if err := cli.AddCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return
	}

	_ = cli.ExecuteWithArgs(context.Background(), []string{"version"})
	// Output:
	// myapp version 1.0.0
}
