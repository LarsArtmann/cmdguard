# cmdguard

**A Go library for building validated Cobra CLI applications with panic-at-construction-time guards and type-safe flags.**

This library wraps Cobra with validation that panics at construction time, ensuring invalid commands are caught immediately at startup rather than failing silently at runtime.

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

## Quick Start

### v2 API (Recommended) - Type-Safe Commands

The v2 API provides full type safety for both configuration and command-specific flags:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is your application-level configuration
type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
    Output  string `flag:"output" short:"o" default:"text" help:"Output format"`
}

// GreetFlags defines command-specific flags (fully typed!)
type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Uppercase output"`
}

func main() {
    // Create CLI with typed config - T=AppConfig, F=NoFlags (root has no flags)
    cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My CLI application", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
        os.Exit(1)
    }

    // Add a command with typed flags - T=AppConfig, F=*GreetFlags
    greetCmd := v2.Command[AppConfig, *GreetFlags]{
        Use:   "greet",
        Short: "Greet someone",
        Flags: &GreetFlags{},
        RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
    }

    // Use AddAnyCommand when command flags type differs from root flags type
    if err := v2.AddAnyCommand(cli, greetCmd); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to add command: %v\n", err)
        os.Exit(1)
    }

    // Execute
    cli.ExecuteAndExit(context.Background())
}
```

### v1 API - Simple Cobra Wrapper

For simpler use cases or migration from raw Cobra:

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    root := cmdguard.New("myapp", "My CLI application")

    root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, World!")
        },
    })

    root.ExecuteAndExit(context.Background())
}
```

## API Comparison

| Feature           | v1 API | v2 API                                   |
| ----------------- | ------ | ---------------------------------------- |
| Type-safe config  | No     | Yes (type parameter T)                   |
| Type-safe flags   | No     | Yes (type parameter F)                   |
| DI integration    | No     | Yes (samber/do/v2)                       |
| Flag tags         | No     | Yes (`flag`, `short`, `default`, `help`) |
| Lifecycle hooks   | No     | Yes (PreRunE, PostRunE)                  |
| Health checks     | No     | Yes                                      |
| Graceful shutdown | No     | Yes                                      |

## v2 API Reference

### GuardedCommand[T, F]

The main CLI type with two type parameters:

- `T` - Application-level config type
- `F` - Root command flags type (use `v2.NoFlags` if none)

```go
// Create root command
cli, err := v2.New[T, F](name, shortDescription, defaultConfig)

// Add subcommands with same F type
cli.AddCommand(cmd Command[T, F]) error

// Add subcommands with different F type
v2.AddAnyCommand[T, F, F2](cli, cmd Command[T, F2]) error

// Execution
err := cli.Execute(ctx)           // Returns error
cli.ExecuteAndExit(ctx)           // Calls os.Exit(1) on error

// DI and lifecycle
scope := cli.Scope()              // do.Injector for service registration
err := cli.HealthCheck()          // Run health checks
err := cli.Shutdown(ctx)          // Graceful shutdown

// Configuration access
cfg := cli.Config()               // *T - typed config
cli.SetConfig(cfg)                // Update config

// Advanced
cmd := cli.RootCommand()          // Underlying cobra.Command
cli.AddGlobalFlag(name, short, default, help)
cli.AddGlobalBoolFlag(name, short, default, help)
```

### Command[T, F]

Type-safe command definition:

```go
type Command[T any, F any] struct {
    Use        string               // Command name/usage
    Short      string               // Short description
    Long       string               // Long description
    Aliases    []string             // Alternative names
    Example    string               // Example usage
    Flags      F                    // Typed flags struct
    RunE       func(ctx, cfg, flags) error
    PreRunE    func(ctx, cfg, flags) error
    PostRunE   func(ctx, cfg, flags) error
    Commands   []Command[T, F]      // Subcommands
    Hidden     bool                 // Hide from help
    Deprecated string               // Deprecation message
}
```

### Flag Tags

Define flags using struct tags:

```go
type MyFlags struct {
    // Required: flag name
    Name string `flag:"name"`

    // Optional: short flag (-n)
    Name string `flag:"name" short:"n"`

    // Optional: default value
    Name string `flag:"name" default:"World"`

    // Optional: help text
    Name string `flag:"name" help:"Name to greet"`

    // All together
    Name string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Count int    `flag:"count" short:"c" default:"1" help:"Number of times"`
    Verbose bool  `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
}
```

### NoFlags

Use `v2.NoFlags` for commands without flags:

```go
versionCmd := v2.Command[AppConfig, v2.NoFlags]{
    Use:   "version",
    Short: "Print version",
    RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        fmt.Println("v1.0.0")
        return nil
    },
}
```

### AddAnyCommand

When a command's flags type differs from the CLI root:

```go
// CLI root has NoFlags
cli, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "...", AppConfig{})

// Greet command has *GreetFlags - different type!
greetCmd := v2.Command[AppConfig, *GreetFlags]{
    Use:   "greet",
    Flags: &GreetFlags{},
    RunE:  ...,
}

// Use AddAnyCommand (standalone function, not method)
v2.AddAnyCommand(cli, greetCmd)
```

### DI Integration

Register and invoke services:

```go
// Register services
v2.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{...}, nil
})

v2.ProvideValue(scope, config)

// Invoke services in handlers (with proper error handling)
RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
    db, err := v2.Invoke[*Database](cli.ScopeStruct())
    if err != nil {
        return v2.NewServiceError("*Database", err)
    }
    // use db...
},
```

### Lifecycle Hooks

```go
cmd := v2.Command[AppConfig, *Flags]{
    Use: "example",
    PreRunE: func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        // Validation before main handler
        if flags.Count < 1 {
            return fmt.Errorf("count must be at least 1")
        }
        return nil
    },
    RunE: func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        // Main command logic
        return nil
    },
    PostRunE: func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        // Cleanup after main handler (called even on error)
        return nil
    },
}
```

### Functional Options

```go
cmd, err := v2.NewCommand[AppConfig, NoFlags]("version",
    v2.WithShort("Print version"),
    v2.WithRunE(func(ctx context.Context, cfg *AppConfig, flags NoFlags) error {
        fmt.Println("v1.0.0")
        return nil
    }),
)
```

## Examples

### Nested Commands

```go
func main() {
    cli, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "My CLI", AppConfig{})

    dbCmd := v2.Command[AppConfig, v2.NoFlags]{
        Use:   "db",
        Short: "Database operations",
        Commands: []v2.Command[AppConfig, v2.NoFlags]{
            {
                Use:   "migrate",
                Short: "Run migrations",
                RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
                    return runMigrations()
                },
            },
            {
                Use:   "rollback",
                Short: "Rollback migrations",
                RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
                    return rollbackMigrations()
                },
            },
        },
    }

    cli.AddCommand(dbCmd)
    cli.ExecuteAndExit(context.Background())
}
```

### Mixed Flag Types

```go
func main() {
    cli, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "My CLI", AppConfig{})

    // Command with GreetFlags
    v2.AddAnyCommand(cli, v2.Command[AppConfig, *GreetFlags]{
        Use:   "greet",
        Flags: &GreetFlags{},
        RunE:  greetHandler,
    })

    // Command with ConfigFlags (different type!)
    v2.AddAnyCommand(cli, v2.Command[AppConfig, *ConfigFlags]{
        Use:   "config",
        Flags: &ConfigFlags{},
        RunE:  configHandler,
    })

    // Command with no flags
    cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
        Use:  "version",
        RunE: versionHandler,
    })
}
```

## How It Works

cmdguard validates commands at construction time:

- **Panic on invalid commands** - Commands without handlers (`Run` or `RunE`) cause immediate panic
- **Panic on empty names** - Commands must have a valid `Use` field
- **Type-safe flags** - v2 API ensures flags are properly typed

This "fail fast" approach catches configuration errors during development, not production.

## Philosophy

**Why panics?**

Go lacks compile-time macros. The closest equivalent to "fail at compile time" is "fail at init time". By panicking on invalid commands during construction:

1. Errors are caught during development, not production
2. Invalid states are impossible to represent at runtime
3. The API is simple - no error handling boilerplate

**When to use cmdguard:**

- You want guaranteed-valid CLI configurations
- You prefer "crash early" over "handle errors later"
- You want type-safe flags (v2)
- You're building CLIs where panics are acceptable (most CLIs)

**When NOT to use cmdguard:**

- You need to handle configuration errors gracefully
- You're embedding CLI in a larger application that can't panic

## Project Status

- **v1** - Stable, minimal Cobra wrapper
- **v2** - Stable, full type-safe API with DI integration

## License

MIT
