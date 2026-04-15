// Typed example demonstrating the v2 API with type-safe commands,
// typed flags, DI integration, and lifecycle hooks.
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/do/v2"

	examplesinternal "github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
// This is shared across all commands and can be populated
// from flags, config files, or environment variables.
type AppConfig struct {
	Verbose bool   `default:"false"                   flag:"verbose" help:"Enable verbose output"      short:"v"`
	Output  string `default:"text"                    flag:"output"  help:"Output format (text, json)" short:"o"`
	APIURL  string `default:"https://api.example.com" flag:"api-url" help:"API server URL"`
}

// GreetFlags defines flags for the greet command.
type GreetFlags struct {
	Name   string `default:"World" flag:"name"   help:"Name to greet"               short:"n"`
	Shout  bool   `default:"false" flag:"shout"  help:"Print greeting in uppercase" short:"s"`
	Count  int    `default:"1"     flag:"count"  help:"Number of times to greet"    short:"c"`
	Prefix string `default:"Hello" flag:"prefix" help:"Greeting prefix"`
	Suffix string `default:"!"     flag:"suffix" help:"Greeting suffix"`
}

// StatsFlags defines flags demonstrating uint and float32 support.
type StatsFlags struct {
	MaxRetries uint    `default:"3"   flag:"max-retries" help:"Maximum retry attempts"`
	Threshold  float32 `default:"0.5" flag:"threshold"   help:"Threshold value (0.0-1.0)"`
	PageSize   uint    `default:"10"  flag:"page-size"   help:"Items per page"`
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
	_ = ctx // context required by interface but not used in this example

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

// printRunE creates a RunE function that prints messages and returns nil.
func printRunE(messages ...string) func(context.Context, *AppConfig, v2.NoFlags) error {
	return func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
		for _, msg := range messages {
			fmt.Println(msg)
		}

		return nil
	}
}

func main() {
	fmt.Println("=== Typed CLI Example (v2 API) ===")
	fmt.Println()

	cli, err := v2.NewCLI[AppConfig]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		examplesinternal.Fatalf("Failed to create CLI: %v\n", err)
	}

	// Set version
	cli.SetVersion("1.0.0")

	// Add global flags (available to all commands)
	cli.AddGlobalBoolFlag("debug", "d", false, "Enable debug mode")

	// Register services in the DI scope
	registerServices(cli.Scope())

	// Add commands
	if err := addCommands(cli); err != nil {
		examplesinternal.Fatalf("Failed to add commands: %v\n", err)
	}

	// Run health checks
	if err := cli.HealthCheck(); err != nil {
		examplesinternal.Fatalf("Health check failed: %v\n", err)
	}

	// Execute the CLI
	ctx := context.Background()
	if err := cli.Execute(ctx); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Shutdown(shutdownCtx); err != nil {
		examplesinternal.Fatalf("Shutdown error: %v\n", err)
	}
}

func registerServices(scope *v2.Scope) {
	// Register config first so providers can depend on it
	err := v2.ProvideValue(scope, AppConfig{Verbose: true})
	if err != nil {
		examplesinternal.Fatalf("Failed to register config: %v\n", err)
	}

	// Register a logger service - gets config via DI, not closure capture
	err = v2.Provide(scope, func(i do.Injector) (*Logger, error) {
		cfg, err := v2.Invoke[*AppConfig](scope)
		if err != nil {
			return nil, fmt.Errorf("invoking *AppConfig in scope %p: %w", scope, err)
		}

		return &Logger{verbose: cfg.Verbose}, nil
	})
	if err != nil {
		examplesinternal.Fatalf("Failed to register logger: %v\n", err)
	}

	// Register a database service (simulated)
	err = v2.ProvideValue(scope, &Database{connectionString: "postgres://localhost:5432/mydb"})
	if err != nil {
		examplesinternal.Fatalf("Failed to register database: %v\n", err)
	}
}

func addCommands(cli *v2.CLI[AppConfig]) error {
	// Greet command with typed flags
	// Using CLI[T] API - each command can have its own flags type
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
				return fmt.Errorf("count must be at least 1 (got %d)", flags.Count)
			}

			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			logger, err := v2.Invoke[*Logger](cli.Scope())
			if err != nil {
				return fmt.Errorf("invoking *Logger in scope %p: %w", cli.Scope(), err)
			}

			for iteration := range flags.Count {
				msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
				if flags.Shout {
					msg = strings.ToUpper(msg)
				}

				fmt.Println(msg)
				logger.Log(fmt.Sprintf("Greeted %s (iteration %d)", flags.Name, iteration+1))
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

	// CLI[T] uses AddCommand function (not method) - supports any flags type per command
	err := v2.AddCommand(cli, greetCmd)
	if err != nil {
		return fmt.Errorf("failed to add greet command: %w", err)
	}

	// Version command
	versionCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "version",
		Short: "Print version information",
		RunE:  printRunE("myapp version 1.0.0", "Built with cmdguard v2"),
	}

	err = v2.AddCommand(cli, versionCmd)
	if err != nil {
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

	err = v2.AddCommand(cli, configCmd)
	if err != nil {
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
					db, err := v2.Invoke[*Database](cli.Scope())
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
				RunE:  printRunE("Running migrations...", "Migration complete!"),
			},
		},
	}

	err = v2.AddCommand(cli, dbCmd)
	if err != nil {
		return fmt.Errorf("failed to add db command: %w", err)
	}

	// Hidden command (won't show in help)
	hiddenCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:    "secret",
		Short:  "Secret command",
		Hidden: true,
		RunE:   printRunE("You found the secret command!"),
	}

	err = v2.AddCommand(cli, hiddenCmd)
	if err != nil {
		return fmt.Errorf("failed to add hidden command: %w", err)
	}

	// Deprecated command
	deprecatedCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:        "oldcmd",
		Short:      "Old command (deprecated)",
		Deprecated: "Use 'greet' instead",
		RunE:       printRunE("This command is deprecated. Use 'greet' instead."),
	}

	err = v2.AddCommand(cli, deprecatedCmd)
	if err != nil {
		return fmt.Errorf("failed to add deprecated command: %w", err)
	}

	// Stats command demonstrating uint and float32 flag support
	statsCmd := v2.Command[AppConfig, *StatsFlags]{
		Use:     "stats",
		Short:   "Display statistics configuration",
		Example: "myapp stats --max-retries 5 --threshold 0.75 --page-size 20",
		Flags:   &StatsFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *StatsFlags) error {
			fmt.Println("Statistics Configuration:")
			fmt.Printf("  Max Retries: %d\n", flags.MaxRetries)
			fmt.Printf("  Threshold:   %.2f\n", flags.Threshold)
			fmt.Printf("  Page Size:   %d\n", flags.PageSize)

			return nil
		},
	}

	err = v2.AddCommand(cli, statsCmd)
	if err != nil {
		return fmt.Errorf("failed to add stats command: %w", err)
	}

	return nil
}
