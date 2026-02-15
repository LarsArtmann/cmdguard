// Typed example demonstrating the v2 API with type-safe commands,
// typed flags, DI integration, and lifecycle hooks.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/samber/do/v2"
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
	Name    string `flag:"name" short:"n" default:"World" help:"Name to greet"`
	Shout   bool   `flag:"shout" short:"s" default:"false" help:"Print greeting in uppercase"`
	Count   int    `flag:"count" short:"c" default:"1" help:"Number of times to greet"`
	Prefix  string `flag:"prefix" default:"Hello" help:"Greeting prefix"`
	Suffix  string `flag:"suffix" default:"!" help:"Greeting suffix"`
}

// Database is a service registered via DI.
type Database string

// Logger is a service registered via DI.
type Logger struct {
	verbose bool
}

func (l *Logger) Log(msg string) {
	if l.verbose {
		fmt.Printf("[LOG] %s\n", msg)
	}
}

func main() {
	fmt.Println("=== Typed CLI Example (v2 API) ===")
	fmt.Println()

	// Create the CLI with typed config
	cli, err := v2.New[AppConfig]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
		os.Exit(1)
	}

	// Set version
	cli.SetVersion("1.0.0")

	// Add global flags (available to all commands)
	cli.AddGlobalBoolFlag("debug", "d", false, "Enable debug mode")

	// Register services in the DI scope
	registerServices(cli.ScopeStruct(), AppConfig{Verbose: true})

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

func registerServices(scope *v2.Scope, cfg AppConfig) {
	// Register a logger service
	if err := v2.Provide(scope, func(i do.Injector) (*Logger, error) {
		return &Logger{verbose: cfg.Verbose}, nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register logger: %v\n", err)
	}

	// Register a database service (simulated)
	if err := v2.ProvideValue(scope, Database("postgres://localhost:5432/mydb")); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register database: %v\n", err)
	}
}

func addCommands(cli *v2.GuardedCommand[AppConfig]) error {
	// Greet command with typed flags
	greetCmd := v2.Command[AppConfig]{
		Use:     "greet [message]",
		Short:   "Greet someone",
		Long:    "Greet someone with a customizable message.",
		Example: "myapp greet --name Alice --shout --count 3",
		Flags:   &GreetFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
			if cfg.Verbose {
				fmt.Println("Preparing to greet...")
			}
			greetFlags := flags.(*GreetFlags)
			if greetFlags.Count < 1 {
				return fmt.Errorf("count must be at least 1")
			}
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
			greetFlags := flags.(*GreetFlags)
			logger := do.MustInvoke[*Logger](cli.Scope())

			for i := 0; i < greetFlags.Count; i++ {
				msg := fmt.Sprintf("%s, %s%s", greetFlags.Prefix, greetFlags.Name, greetFlags.Suffix)
				if greetFlags.Shout {
					msg = stringsToUpper(msg)
				}
				fmt.Println(msg)
				logger.Log(fmt.Sprintf("Greeted %s (iteration %d)", greetFlags.Name, i+1))
			}
			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
			if cfg.Verbose {
				fmt.Println("Greeting complete!")
			}
			return nil
		},
	}
	if err := cli.AddCommand(greetCmd); err != nil {
		return fmt.Errorf("failed to add greet command: %w", err)
	}

	// Version command
	versionCmd := v2.Command[AppConfig]{
		Use:   "version",
		Short: "Print version information",
		RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
			fmt.Println("myapp version 1.0.0")
			fmt.Println("Built with cmdguard v2")
			return nil
		},
	}
	if err := cli.AddCommand(versionCmd); err != nil {
		return fmt.Errorf("failed to add version command: %w", err)
	}

	// Config command that uses the app config
	configCmd := v2.Command[AppConfig]{
		Use:   "config",
		Short: "Show current configuration",
		RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
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
	dbCmd := v2.Command[AppConfig]{
		Use:   "db",
		Short: "Database operations",
		Commands: []v2.Command[AppConfig]{
			{
				Use:   "status",
				Short: "Check database connection",
				RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
					db := do.MustInvoke[Database](cli.Scope())
					fmt.Printf("Database: %s\n", db)
					fmt.Println("Status: Connected (simulated)")
					return nil
				},
			},
			{
				Use:   "migrate",
				Short: "Run database migrations",
				RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
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
	hiddenCmd := v2.Command[AppConfig]{
		Use:    "secret",
		Short:  "Secret command",
		Hidden: true,
		RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
			fmt.Println("You found the secret command!")
			return nil
		},
	}
	if err := cli.AddCommand(hiddenCmd); err != nil {
		return fmt.Errorf("failed to add hidden command: %w", err)
	}

	// Deprecated command
	deprecatedCmd := v2.Command[AppConfig]{
		Use:        "oldcmd",
		Short:      "Old command (deprecated)",
		Deprecated: "Use 'greet' instead",
		RunE: func(ctx context.Context, cfg *AppConfig, flags any) error {
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
