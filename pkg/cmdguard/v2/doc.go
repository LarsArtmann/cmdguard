// Package v2 provides a type-safe, dependency-injection-powered CLI framework built on Cobra.
//
// cmdguard eliminates runtime flag lookups by using struct tags for flag definitions,
// provides compile-time type safety for commands and flags, and integrates samber/do/v2
// for dependency injection — all without panics in library code.
//
// Version: 2.7.0-dev
//
// # Quick Start
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//
//	    v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
//	)
//
//	type AppConfig struct {
//	    Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
//	}
//
//	func main() {
//	    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("hello",
//	        func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
//	            fmt.Println("Hello, World!")
//	            return nil
//	        },
//	        v2.WithShort[AppConfig, v2.NoFlags]("Say hello"),
//	    )
//	    if err != nil {
//	        panic(err)
//	    }
//
//	    if err := v2.AddCommand(cli, cmd); err != nil {
//	        panic(err)
//	    }
//
//	    if err := cli.Execute(context.Background()); err != nil {
//	        fmt.Println("Error:", err)
//	    }
//	}
//
// # Type-Safe Flags
//
// Define flags as struct fields with tags instead of string lookups:
//
//	type DeployFlags struct {
//	    Env     string        `flag:"env"     short:"e" required:"true" help:"Target environment"`
//	    DryRun  bool          `flag:"dry-run" short:"d" default:"false" help:"Simulate deployment"`
//	    Timeout time.Duration `flag:"timeout" short:"t" default:"5m"    help:"Deployment timeout"`
//	}
//
// Pass the flags struct to NewCommand with WithFlags:
//
//	cmd, err := v2.NewCommand[AppConfig, *DeployFlags]("deploy", handler,
//	    v2.WithShort[AppConfig, *DeployFlags]("Deploy the application"),
//	    v2.WithFlags[AppConfig, *DeployFlags](&DeployFlags{}),
//	)
//
// Supported tags: flag, short, default, help, env, required, count.
//
// # Dependency Injection
//
// Register services on the CLI scope and invoke them in handlers:
//
//	scope := cli.Scope()
//
//	v2.Provide(scope, func(i do.Injector) (*Database, error) {
//	    return &Database{DSN: "postgres://..."}, nil
//	})
//
//	// In handler:
//	db, err := v2.Invoke[*Database](scope)
//
// Services can implement HealthCheck and Shutdown for lifecycle management.
//
// For testing, clone the scope and override services with mocks:
//
//	cloned := v2.CloneScope(scope)
//	v2.OverrideValue(cloned, &MockDatabase{})
//	mockDB, _ := v2.Invoke[*Database](cloned)
//
// # Command Options
//
// Commands are created via constructors and configured with functional options:
//
//	WithShort[T, F](short)        // Short description (required with StrictValidation)
//	WithLong[T, F](long)          // Long description
//	WithExample[T, F](example)    // Example usage
//	WithFlags[T, F](flags)        // Typed flags struct
//	WithPreRunE[T, F](fn)         // Pre-validation hook
//	WithPostRunE[T, F](fn)        // Post-success cleanup
//	WithExactArgs[T, F](n)        // Require exactly n positional args
//	WithCompletion[T, F](fn)      // Dynamic shell completion
//	WithHidden[T, F](bool)        // Hide from help
//	WithDeprecated[T, F](msg)     // Deprecation message
//
// Parent commands (commands with subcommands) use NewParentCommand:
//
//	parent, err := v2.NewParentCommand[AppConfig, v2.NoFlags]("user", "User management",
//	    []v2.Command[AppConfig, v2.NoFlags]{listCmd, createCmd},
//	    v2.WithShort[AppConfig, v2.NoFlags]("User management"),
//	)
//
// # CLI Options
//
// Configure the root CLI with options passed to NewCLI:
//
//	WithCLIVersion[T](v)              // Version string (auto-pipes to fang)
//	WithCLICommit[T](commit)          // Git commit hash (auto-pipes to fang)
//	WithEnvPrefix[T](prefix)          // Prefix for env var lookups
//	WithSignalHandling[T]()           // Cancel context on SIGINT/SIGTERM
//	WithGracefulShutdown[T]()         // Graceful DI shutdown on SIGINT/SIGTERM
//	WithDILogging[T](logf)            // Internal DI container logging
//	WithMiddleware[T](mw...)          // Wrap all command handlers
//	WithStrictValidation[T]()         // Require WithShort on all commands
//	WithConfigValidation[T](fn)       // Validate config after flag parsing
//	WithConfigFile[T](paths...)       // Load JSON config before flags
//	WithFang[T](bool)                 // Styled help output
//	WithFangOptions[T](opts...)       // Custom fang options
//	WithFangErrorHandler[T](handler)  // Custom fang error handler
//	WithFangColorScheme[T](cs)        // Custom fang color scheme
//
// # Error Handling
//
// All v2 constructors return errors. Functions never panic — every function
// returns errors. Sentinel errors support errors.Is() for identification:
//
//	errors.Is(err, v2.ErrInvalidCommand)
//	errors.Is(err, v2.ErrMissingHandler)
//	errors.Is(err, v2.ErrDuplicateCommand)
//
// Rich error types add context:
//
//	v2.NewCommandError(name, err)
//	v2.NewFlagError(name, err)
//	v2.NewFlagErrorWithSuggestion(name, err, suggestion)
//	v2.NewExitError(code, err)    // Custom exit code
//
// Check for custom exit codes with the ExitCoder interface:
//
//	if exitCoder, ok := errors.AsType[v2.ExitCoder](err); ok {
//	    fmt.Println("Exit code:", exitCoder.ExitCode())
//	}
//
// # Middleware
//
// Wrap all command handlers with middleware:
//
//	cli, err := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
//	    v2.WithMiddleware[AppConfig](
//	        v2.TimingMiddleware[AppConfig](),
//	        v2.RecoveryMiddleware[AppConfig](),
//	    ),
//	)
//
// Write custom middleware by implementing the MiddlewareFunc type.
//
// # Output Formats
//
// Render structured data in 16 formats using go-output directly:
//
//	v2.OutputTable(output.FormatJSON, headers, rows)
//	v2.OutputTable(output.FormatCSV, headers, rows)
//	v2.OutputTable(output.FormatYAML, headers, rows)
//
// Available: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree,
// mermaid, dot, jsonl, asciidoc, toml, plantuml.
//
// WithCLIVersion and WithCLICommit automatically pipe version/commit info to fang.
// Do NOT also pass these via WithFangOptions or you will get duplicates.
//
// # BranchingFlowContext
//
// Track command execution paths and share values across the hierarchy:
//
//	bfc, ok := v2.GetBranchingFlowContext(ctx)
//	if ok {
//	    fmt.Println("Path:", bfc.PathString())
//	    bfc.SetValue("key", "value")
//	}
//
// # Editor Support
//
// Open the user's $EDITOR for interactive input:
//
//	edited, err := v2.EditInEditor(ctx, "# Edit here\n")
//
// # Version Command
//
// Add a built-in version command:
//
//	cmd, err := v2.VersionCommand[AppConfig](cli)
//	if err != nil {
//	    panic(err)
//	}
//	v2.AddCommand(cli, cmd)
//
// # Custom Types
//
// Register custom flag types with full parse and validate support:
//
//	v2.RegisterTypeHandler(reflect.TypeFor[MyType](), v2.TypeHandlerFunc{
//	    ParseFunc:   func(value string, _ v2.FlagTag) (any, error) { return MyType{value}, nil },
//	    DefaultFunc: func(_ v2.FlagTag) any { return MyType{} },
//	})
//
// Built-in custom types: Duration, Enum[T], LogLevel, URL, Email, Port, FilePath, HostPort.
//
// # Further Reading
//
// See the examples/ directory for working demonstrations of each feature.
// Visit https://pkg.go.dev/github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2 for the full API reference.
package v2
