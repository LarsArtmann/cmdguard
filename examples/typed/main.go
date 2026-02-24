// Typed example demonstrating the v2 API with type-safe commands,
// typed flags, DI integration, and lifecycle hooks.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
// This is shared across all commands and can be populated
// from flags, config files, or environment variables.
type AppConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	Output  string `flag:"output" short:"o" default:"text" help:"Output format (text, json)"`
	APIURL  string `flag:"api-url" default:"https://api.example.com" help:"API server URL"`
}

// GreetFlags defines flags for the greet command.
type GreetFlags struct {
	Name   string `flag:"name" short:"n" default:"World" help:"Name to greet"`
	Shout  bool   `flag:"shout" short:"s" default:"false" help:"Print greeting in uppercase"`
	Count  int    `flag:"count" short:"c" default:"1" help:"Number of times to greet"`
	Prefix string `flag:"prefix" default:"Hello" help:"Greeting prefix"`
	Suffix string `flag:"suffix" default:"!" help:"Greeting suffix"`
}

// Logger is a service registered via DI.
type Logger struct {
	verbose bool
}

func (l *Logger) Log(msg string) {
	if l.verbose {
		fmt.Printf("[LOG] %s\n", msg)
	}
}

// HealthCheck implements do.HealthcheckerWithContext for lifecycle demonstration.
func (l *Logger) HealthCheck(ctx context.Context) error {
	if l.verbose {
		fmt.Println("[LOG] Health check passed")
	}
	return nil
}

// Database is a service registered via DI.
type Database struct {
	connectionString string
}

// Shutdown implements do.Shutdowner for lifecycle demonstration.
func (d *Database) Shutdown() error {
	fmt.Printf("[DB] Closing connection to %s\n", d.connectionString)
	return nil
}

func main() {
	fmt.Println("=== Typed CLI Example (v2 API) ===")
	fmt.Println()

	// Create the CLI with typed config
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
		os.Exit(1)
	}

	// Set version
	cli.SetVersion("1.0.0")

	// Add global flags (available to all commands)
	cli.AddGlobalBoolFlag("debug", "d", false, "Enable debug mode")

	// Register services in the DI scope
	registerServices(cli.ScopeStruct())

	// Add commands
	if err := addCommands(cli); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add commands: %v\n", err)
		os.Exit(1)
	}

	// Run health checks
	if err := cli.HealthCheck(); err != nil {
		fmt.Fprintf(os.Stderr, "Health check failed: %v\n", err)
	}

	// Execute the CLI
	ctx := context.Background()
	if err := cli.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cli.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
	}
}

func registerServices(scope *v2.Scope) {
	// Register config first so providers can depend on it
	if err := v2.ProvideValue(scope, AppConfig{Verbose: true}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register config: %v\n", err)
	}

	// Register a logger service - gets config via DI, not closure capture
	if err := v2.Provide(scope, func(i do.Injector) (*Logger, error) {
		cfg, err := v2.Invoke[*AppConfig](scope)
		if err != nil {
			return nil, v2.NewServiceError("*AppConfig", err)
		}
		return &Logger{verbose: cfg.Verbose}, nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register logger: %v\n", err)
	}

	// Register a database service (simulated)
	if err := v2.ProvideValue(scope, &Database{connectionString: "postgres://localhost:5432/mydb"}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register database: %v\n", err)
	}
}

func addCommands(cli *v2.GuardedCommand[AppConfig, v2.NoFlags]) error {
	// Greet command with typed flags
	greetCmd := v2.Command[AppConfig, *GreetFlags]{
		Use:     "greet [message]",
		Short:   "Greet someone",
		Long:    "Greet someone with a customizable message.",
		Example: "myapp greet --name Alice --shout --count 3",
		Flags:   &GreetFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			if cfg.Verbose {
				fmt.Println("Preparing to greet...")
			}
			if flags.Count < 1 {
				return fmt.Errorf("count must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			logger, err := v2.Invoke[*Logger](cli.ScopeStruct())
			if err != nil {
				return v2.NewServiceError("*Logger", err)
			}

			for i := 0; i < flags.Count; i++ {
				msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
				if flags.Shout {
					msg = stringsToUpper(msg)
				}
				fmt.Println(msg)
				logger.Log(fmt.Sprintf("Greeted %s (iteration %d)", flags.Name, i+1))
			}
			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			if cfg.Verbose {
				fmt.Println("Greeting complete!")
			}
			return nil
		},
	}
	if err := v2.AddAnyCommand(cli, greetCmd); err != nil {
		return fmt.Errorf("failed to add greet command: %w", err)
	}

	// Version command
	versionCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "version",
		Short: "Print version information",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("myapp version 1.0.0")
			fmt.Println("Built with cmdguard v2")
			return nil
		},
	}
	if err := cli.AddCommand(versionCmd); err != nil {
		return fmt.Errorf("failed to add version command: %w", err)
	}

	// Config command that uses the app config
	configCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "config",
		Short: "Show current configuration",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("Current configuration:")
			fmt.Printf("  Verbose: %v\n", cfg.Verbose)
			fmt.Printf("  Output:  %s\n", cfg.Output)
			fmt.Printf("  API URL: %s\n", cfg.APIURL)
			return nil
		},
	}
	if err := cli.AddCommand(configCmd); err != nil {
		return fmt.Errorf("failed to add config command: %w", err)
	}

	// Parent command with subcommands
	dbCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "db",
		Short: "Database operations",
		Commands: []v2.Command[AppConfig, v2.NoFlags]{
			{
				Use:   "status",
				Short: "Check database connection",
				RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
					db, err := v2.Invoke[*Database](cli.ScopeStruct())
					if err != nil {
						return v2.NewServiceError("*Database", err)
					}
					fmt.Printf("Database: %s\n", db.connectionString)
					fmt.Println("Status: Connected (simulated)")
					return nil
				},
			},
			{
				Use:   "migrate",
				Short: "Run database migrations",
				RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
					fmt.Println("Running migrations...")
					fmt.Println("Migration complete!")
					return nil
				},
			},
		},
	}
	if err := cli.AddCommand(dbCmd); err != nil {
		return fmt.Errorf("failed to add db command: %w", err)
	}

	// Hidden command (won't show in help)
	hiddenCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:    "secret",
		Short:  "Secret command",
		Hidden: true,
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("You found the secret command!")
			return nil
		},
	}
	if err := cli.AddCommand(hiddenCmd); err != nil {
		return fmt.Errorf("failed to add hidden command: %w", err)
	}

	// Deprecated command
	deprecatedCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:        "oldcmd",
		Short:      "Old command (deprecated)",
		Deprecated: "Use 'greet' instead",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("This command is deprecated. Use 'greet' instead.")
			return nil
		},
	}
	if err := cli.AddCommand(deprecatedCmd); err != nil {
		return fmt.Errorf("failed to add deprecated command: %w", err)
	}

	return nil
}

func stringsToUpper(s string) string {
	result := make([]byte, len(s))
	for i, c := range []byte(s) {
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}
