// Package v3 provides a type-safe, dependency-injection-powered CLI framework built on Cobra.
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
//	    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
//	)
//
//	type AppConfig struct {
//	    Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
//	}
//
//	func main() {
//	    cli, err := v3.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    cmd, err := v3.NewCommand("hello", v3.NoFlags{},
//	        func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
//	            fmt.Println("Hello, World!")
//	            return nil
//	        },
//	        v3.WithShort("Say hello"),
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    if err := v3.AddCommand(cli, cmd); err != nil {
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
//	cmd, err := v3.NewCommand("deploy", &DeployFlags{}, handler,
//	    v3.WithShort("Deploy the application"),
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
//	v3.Provide(scope, func(i do.Injector) (*Database, error) {
//	    return &Database{DSN: "postgres://..."}, nil
//	})
//
//	// In handler:
//	db, err := v3.Invoke[*Database](scope)
//
// Services can implement HealthCheck and Shutdown for lifecycle management.
//
// For testing, clone the scope and override services with mocks:
//
//	cloned := v3.CloneScope(scope)
//	v3.OverrideValue(cloned, &MockDatabase{})
//	mockDB, _ := v3.Invoke[*Database](cloned)
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
//	parent, err := v3.NewParentCommand[AppConfig]("user", "User management", v3.NoFlags{},
//	    v3.WithSubcommands(listCmd, createCmd),
//	    v3.WithShort("User management"),
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
//	github.com/larsartmann/cmdguard/manpage   — Man page generation (mango/roff)
//
// # Error Handling
//
// All v2 constructors return errors. Functions never panic — every function
// returns errors. Sentinel errors support errors.Is() for identification:
//
//	errors.Is(err, v3.ErrInvalidCommand)
//	errors.Is(err, v3.ErrMissingHandler)
//	errors.Is(err, v3.ErrDuplicateCommand)
//
// Rich error types add context:
//
//	v3.NewCommandError(name, err)
//	v3.NewFlagError(name, err)
//	v3.NewFlagErrorWithSuggestion(name, err, suggestion)
//	v3.NewExitError(code, err)    // Custom exit code
//
// Check for custom exit codes with the ExitCoder interface:
//
//	if exitCoder, ok := errors.AsType[v3.ExitCoder](err); ok {
//	    fmt.Println("Exit code:", exitCoder.ExitCode())
//	}
//
// # Middleware
//
// Wrap all command handlers with middleware:
//
//	cli, err := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
//	    v3.WithMiddleware[AppConfig](
//	        v3.TimingMiddleware[AppConfig](),
//	        v3.RecoveryMiddleware[AppConfig](),
//	    ),
//	)
//
// # Version and Doctor Commands
//
// Add built-in helper commands:
//
//	cmd, err := v3.VersionCommand[AppConfig](cli)
//	docCmd, err := v3.DoctorCommand[AppConfig](cli)
//
// # Custom Types
//
// Register custom flag types with full parse and validate support:
//
//	v3.RegisterTypeHandler(reflect.TypeFor[MyType](), v3.TypeHandlerFunc{
//	    ParseFunc:   func(value string, _ v3.FlagTag) (any, error) { return MyType{value}, nil },
//	    DefaultFunc: func(_ v3.FlagTag) any { return MyType{} },
//	})
//
// Built-in custom types: Duration, Enum, LogLevel, URL, Email, Port, FilePath, HostPort.
//
// # Further Reading
//
// See the examples/ directory for working demonstrations of each feature.
// Visit https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3 for the full API reference.
package v3
