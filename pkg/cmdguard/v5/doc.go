// Package v5 provides a type-safe, dependency-injection-powered CLI framework built on Cobra.
//
// cmdguard eliminates runtime flag lookups by using struct tags for flag definitions,
// provides compile-time type safety for commands and flags, and integrates samber/do/v2
// for dependency injection — all without panics in library code.
//
// # Quick Start
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "log"
//
//	    "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
//	)
//
//	type AppConfig struct {
//	    Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
//	}
//
//	func main() {
//	    cli, err := v5.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    cmd, err := v5.NewCommand("hello", v5.NoFlags{},
//	        func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
//	            fmt.Println("Hello, World!")
//	            return nil
//	        },
//	        v5.WithShort("Say hello"),
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    if err := v5.AddCommand(cli, cmd); err != nil {
//	        log.Fatal(err)
//	    }
//
//	    if err := cli.Execute(context.Background()); err != nil {
//	        fmt.Println("Error:", err)
//	    }
//	}
//
// # Type-Safe Flags
//
// Define flags as struct fields with tags, then pass the struct positionally
// to NewCommand. Type parameters T and F are inferred automatically — no
// explicit type parameters needed on options:
//
//	type DeployFlags struct {
//	    Env     string        `flag:"env"     short:"e" required:"true" help:"Target environment"`
//	    DryRun  bool          `flag:"dry-run" short:"d" default:"false" help:"Simulate deployment"`
//	    Timeout time.Duration `flag:"timeout" short:"t" default:"5m"    help:"Deployment timeout"`
//	}
//
//	cmd, err := v5.NewCommand("deploy", &DeployFlags{}, handler,
//	    v5.WithShort("Deploy the application"),
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
//	v5.Provide(scope, func(i do.Injector) (*Database, error) {
//	    return &Database{DSN: "postgres://..."}, nil
//	})
//
//	// In handler:
//	db, err := v5.Invoke[*Database](scope)
//
// Services can implement HealthCheck(ctx) (the do.HealthcheckerWithContext
// interface, wired into DoctorCommand) and Shutdown for lifecycle management.
//
// For testing, clone the scope and override services with mocks:
//
//	cloned := v5.CloneScope(scope)
//	v5.OverrideValue(cloned, &MockDatabase{})
//	mockDB, _ := v5.Invoke[*Database](cloned)
//
// # Command Options
//
// All metadata options are non-generic — no type parameters needed.
// Only lifecycle hooks (WithPreRunE, WithPostRunE, WithSubcommands)
// are generic functions that return a non-generic CommandOption:
//
//	WithShort(short)              // Short description (required with StrictValidation)
//	WithLong(long)                // Long description
//	WithExample(example)          // Example usage
//	WithPreRunE[T,F](fn)          // Pre-validation hook
//	WithPostRunE[T,F](fn)         // Post-success cleanup
//	WithExactArgs(n)              // Require exactly n positional args
//	WithCompletion(fn)            // Dynamic shell completion
//	WithHidden(bool)              // Hide from help
//	WithDeprecated(msg)           // Deprecation message
//
// Parent commands use NewParentCommand with WithSubcommands:
//
//	parent, err := v5.NewParentCommand[AppConfig]("user", "User management", v5.NoFlags{},
//	    v5.WithSubcommands(listCmd, createCmd),
//	    v5.WithShort("User management"),
//	)
//
// # Optional Sub-Modules
//
// Heavy dependencies are extracted into optional importable modules to keep
// the core dependency tree minimal:
//
//	github.com/larsartmann/cmdguard/glamour   — Markdown help rendering (glamour/v2)
//	github.com/larsartmann/cmdguard/prompts   — Interactive prompts (huh/v2)
//	github.com/larsartmann/cmdguard/telemetry — OpenTelemetry spans
//
// # Error Handling
//
// All constructors return errors. Functions never panic — every function
// returns errors. Sentinel errors support errors.Is() for identification:
//
//	errors.Is(err, v5.ErrInvalidCommand)
//	errors.Is(err, v5.ErrMissingHandler)
//	errors.Is(err, v5.ErrDuplicateCommand)
//
// Rich error types add context:
//
//	v5.NewCommandError(name, err)
//	v5.NewFlagError(name, err)
//	v5.NewFlagErrorWithSuggestion(name, err, suggestion)
//	v5.NewExitError(code, err)    // Custom exit code
//
// Check for custom exit codes with the ExitCoder interface:
//
//	if exitCoder, ok := errors.AsType[v5.ExitCoder](err); ok {
//	    fmt.Println("Exit code:", exitCoder.ExitCode())
//	}
//
// # Middleware
//
// Wrap all command handlers with middleware:
//
//	cli, err := v5.NewCLI[AppConfig]("myapp", "...", AppConfig{},
//	    v5.WithMiddleware[AppConfig](
//	        v5.TimingMiddleware[AppConfig](),
//	        v5.RecoveryMiddleware[AppConfig](),
//	    ),
//	)
//
// # Version and Doctor Commands
//
// Add built-in helper commands:
//
//	cmd, err := v5.VersionCommand[AppConfig](cli)
//	docCmd, err := v5.DoctorCommand[AppConfig](cli)
//
// # Custom Types
//
// Register custom flag types with full parse and validate support:
//
//	v5.RegisterTypeCodec(reflect.TypeFor[MyType](), v5.TypeCodecFunc{
//	    ParseFunc:   func(value string, _ v5.FlagTag) (any, error) { return MyType{value}, nil },
//	    DefaultFunc: func(_ v5.FlagTag) any { return MyType{} },
//	})
//
// Built-in custom types: Duration, Enum, LogLevel, URL, Email, Port, FilePath, HostPort.
//
// # Further Reading
//
// See the examples/ directory for working demonstrations of each feature.
// Visit https://pkg.go.dev/github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5 for the full API reference.
package v5
