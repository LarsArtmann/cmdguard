# cmdguard

[![CI](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml)

**A Go library for building validated CLI applications with type-safe flags and dependency injection.**

cmdguard wraps [Cobra](https://github.com/spf13/cobra) with validation and provides two APIs:

| API                  | Package           | Description                                |
| -------------------- | ----------------- | ------------------------------------------ |
| **v2** (Recommended) | `pkg/cmdguard/v2` | Type-safe, DI-powered, returns errors      |
| v1                   | `pkg/cmdguard`    | Simple wrapper, panics on invalid commands |

v2 is the recommended API—it never panics and provides full type safety for configuration and command flags.

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

cmdguard v2 provides built-in dependency injection through [samber/do/v2](https://github.com/samber/do), enabling clean service management and lifecycle handling.

#### Scope Hierarchy

Each `GuardedCommand` has a root scope that can create child scopes:

```go
// Get the root scope
scope := cli.ScopeStruct()

// Create child scopes for isolation
childScope := scope.Child("worker")

// Access scope hierarchy
path := scope.Path()     // Returns ["myapp"]
isRoot := scope.IsRoot() // Returns true for root scope
```

#### Registering Services

Register services using constructors or values:

```go
// Register with constructor (lazy initialization)
err := v2.Provide(scope, func(i do.Injector) (*Database, error) {
    // Access config through injector
    cfg, err := v2.Invoke[*AppConfig](scope)
    if err != nil {
        return nil, err
    }
    return &Database{DSN: cfg.DSN}, nil
})

// Register pre-constructed value
err = v2.ProvideValue(scope, &Logger{Level: "info"})
```

#### Invoking Services

Retrieve services in command handlers:

```go
RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
    db, err := v2.Invoke[*Database](cli.ScopeStruct())
    if err != nil {
        return v2.NewServiceError("*Database", err)
    }
    // Use db...
},
```

#### Scoped Providers

Create providers that operate within specific scopes:

```go
// Register a scoped provider
err := v2.ScopedProvider(scope, "worker", func(i do.Injector) (*Worker, error) {
    return &Worker{}, nil
})

// Invoke within that scope
worker, err := v2.Invoke[*Worker](scope.Child("worker"))
```

#### Lifecycle Management

Services can implement lifecycle hooks:

```go
// Implement do.HealthcheckerWithContext for health checks
type Database struct{}

func (d *Database) HealthCheck(ctx context.Context) error {
    return d.Ping(ctx)
}

// Implement do.Shutdowner for graceful shutdown
type Server struct{}

func (s *Server) Shutdown() error {
    return s.server.Close()
}

// Run health checks
err := cli.HealthCheck()

// Graceful shutdown
err := cli.Shutdown(ctx)
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

- **v2 API (recommended)** returns errors on invalid commands—no panics
- **v1 API** panics on invalid commands for fail-fast behavior
- **Type-safe flags** ensure flags are properly typed at compile time
- **DI integration** (v2) enables clean service management and lifecycle handling

This approach catches configuration errors during development, not production.

## Philosophy

**Why two APIs?**

- **v1** panics on invalid commands—simple, fail-fast approach
- **v2** returns errors—flexible, production-friendly, type-safe

**When to use v2 (recommended):**

- You want type-safe flags with struct tags
- You need dependency injection for services
- You prefer graceful error handling over panics
- You're building production CLIs

**When to use v1:**

- You want a simple Cobra wrapper
- You prefer "crash early" over error handling
- You're building quick scripts or prototypes

## Project Status

| API    | Status | Description                            |
| ------ | ------ | -------------------------------------- |
| **v2** | Stable | Full type-safe API with DI integration |
| v1     | Stable | Simple Cobra wrapper with panics       |

Both APIs are production-ready. Use v2 for new projects.

## Documentation

- [Quick Start Guide](docs/QUICKSTART.md) - Get started with cmdguard v2 in 5 minutes
- [Migration Guide v1 to v2](docs/MIGRATION_v1_v2.md) - Migrating from v1 to v2 API
- [FEATURES.md](docs/FEATURES.md) - Feature status and roadmap
- [CLI Design Principles](docs/CLI_DESIGN_PRINCIPLES.md) - Design guidelines

## License

MIT
