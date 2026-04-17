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

The v2 API uses `CLI[T]` with a single type parameter for your config type. Each command can have its own flags type via `Command[T, F]`.

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
    // Create CLI with typed config
    cli, err := v2.NewCLI[AppConfig]("myapp", "My CLI application", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
        os.Exit(1)
    }

    // Add a command with typed flags (constructor validates at creation)
    greetCmd, err := v2.NewCommand[AppConfig, *GreetFlags]("greet",
        func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
        v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
        v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create command: %v\n", err)
        os.Exit(1)
    }

    err = v2.AddCommand(cli, greetCmd)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to add command: %v\n", err)
        os.Exit(1)
    }

    // Execute
    if err := cli.Execute(context.Background()); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
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

### CLI[T]

The main CLI type with a single type parameter `T` for your application config.

```go
// Create root CLI
cli, err := v2.NewCLI[T](name, shortDescription, defaultConfig)

// With options
cli, err := v2.NewCLI[T](name, short, defaults,
    v2.WithCLIVersion[T]("1.0.0"),
    v2.WithCLILong[T]("A longer description..."),
    v2.WithSilenceErrors[T](),
    v2.WithSilenceUsage[T](),
    v2.WithColor[T](false),  // disable fang styling
)

// Add subcommands (standalone function — each command has its own flags type)
err = v2.AddCommand(cli, cmd)

// Execution
err := cli.Execute(ctx)           // Returns error
cli.ExecuteAndExit(ctx)           // Calls os.Exit(1) on error

// DI and lifecycle
scope := cli.Scope()              // *Scope for service registration
err := cli.HealthCheck()          // Run health checks
err := cli.Shutdown(ctx)          // Graceful shutdown

// Configuration access
cfg := cli.Config()               // *T - typed config
cli.SetConfig(cfg)                // Update config

// Advanced
cmd := cli.RootCommand()          // Underlying cobra.Command
cli.AddGlobalFlag(name, short, default, help)
cli.AddGlobalBoolFlag(name, short, default, help)
cli.SetVersion("1.0.0")
```

### Command[T, F]

Type-safe command definition created via constructors:

```go
// Leaf command — requires use string and handler
func NewCommand[T, F any](use string, runE func(ctx context.Context, cfg *T, flags F) error, opts ...CommandOption[T, F]) (Command[T, F], error)

// Parent command — requires use, long description, and subcommands
func NewParentCommand[T, F any](use string, long string, subcommands []Command[T, F], opts ...CommandOption[T, F]) (Command[T, F], error)

// Panic variants for compile-time-known configuration
func MustNewCommand[T, F any](...) Command[T, F]
func MustNewParentCommand[T, F any](...) Command[T, F]
```

Command options:

| Option                           | Purpose              |
| -------------------------------- | -------------------- |
| `WithShort[T, F](short)`         | Short description    |
| `WithLong[T, F](long)`           | Long description     |
| `WithAliases[T, F](aliases...)`  | Alternative names    |
| `WithExample[T, F](example)`     | Example usage        |
| `WithFlags[T, F](flags)`         | Typed flags struct   |
| `WithPreRunE[T, F](preRunE)`     | Pre-validation hook  |
| `WithPostRunE[T, F](postRunE)`   | Post-success cleanup |
| `WithSubcommands[T, F](cmds...)` | Child commands       |
| `WithHidden[T, F](hidden)`       | Hide from help       |
| `WithDeprecated[T, F](msg)`      | Deprecation message  |
| `WithGroupID[T, F](group)`       | Help group name      |

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
    Name    string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Count   int    `flag:"count" short:"c" default:"1" help:"Number of times"`
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
}
```

Supported types: `string`, `bool`, `int`, `uint`, `float64`, `[]string`, `time.Duration`.

### NoFlags

Use `v2.NoFlags` for commands without flags:

```go
cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("version",
    func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        fmt.Println("v1.0.0")
        return nil
    },
    v2.WithShort[AppConfig, v2.NoFlags]("Print version"),
)
```

### Mixing Flag Types

`AddCommand` is a standalone function, not a method, because each command can have a different flags type:

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{})

// Command with GreetFlags
greetCmd, _ := v2.NewCommand[AppConfig, *GreetFlags]("greet", greetHandler,
    v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v2.AddCommand(cli, greetCmd)

// Command with ConfigFlags (different type!)
configCmd, _ := v2.NewCommand[AppConfig, *ConfigFlags]("config", configHandler,
    v2.WithFlags[AppConfig, *ConfigFlags](&ConfigFlags{}),
)
v2.AddCommand(cli, configCmd)

// Command with no flags
versionCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("version", versionHandler)
v2.AddCommand(cli, versionCmd)
```

### DI Integration

cmdguard v2 provides built-in dependency injection through [samber/do/v2](https://github.com/samber/do), enabling clean service management and lifecycle handling.

#### Registering Services

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{})
scope := cli.Scope()

// Register with constructor (lazy initialization)
err := v2.Provide(scope, func(i do.Injector) (*Database, error) {
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

```go
RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
    db, err := v2.Invoke[*Database](cli.Scope())
    if err != nil {
        return v2.NewServiceError("*Database", err)
    }
    // Use db...
    return nil
},
```

#### Lifecycle Management

Services can implement lifecycle hooks:

```go
// Health checks — implement do.HealthcheckerWithContext
func (d *Database) HealthCheck(ctx context.Context) error {
    return d.Ping(ctx)
}

// Graceful shutdown — implement do.Shutdowner
func (s *Server) Shutdown() error {
    return s.server.Close()
}

// Run health checks
err := cli.HealthCheck()

// Graceful shutdown
err := cli.Shutdown(ctx)
```

#### Scope Hierarchy

```go
scope := cli.Scope()

// Create child scopes for isolation
childScope := scope.Child("worker")

// Access scope hierarchy
path := scope.Path()     // Returns ["myapp"]
isRoot := scope.IsRoot() // Returns true for root scope
```

### Lifecycle Hooks

```go
cmd, err := v2.NewCommand[AppConfig, *Flags]("example", runHandler,
    v2.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        // Validation before main handler
        if flags.Count < 1 {
            return fmt.Errorf("count must be at least 1")
        }
        return nil
    }),
    v2.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        // Cleanup after main handler (only called on success)
        return nil
    }),
)
```

### Functional Options

Commands can be built with functional options:

```go
// Using NewCommand with functional options
cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("version",
    func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        fmt.Println("v1.0.0")
        return nil
    },
    v2.WithShort[AppConfig, v2.NoFlags]("Print version"),
)

// Using MustNewCommand for compile-time-known config (panics on error)
cmd := v2.MustNewCommand[AppConfig, v2.NoFlags]("version", versionHandler,
    v2.WithShort[AppConfig, v2.NoFlags]("Print version"),
)
```

## Examples

### Nested Commands

```go
func main() {
    cli, _ := v2.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})

    migrateCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("migrate",
        func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            return runMigrations()
        },
        v2.WithShort[AppConfig, v2.NoFlags]("Run migrations"),
    )

    rollbackCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("rollback",
        func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            return rollbackMigrations()
        },
        v2.WithShort[AppConfig, v2.NoFlags]("Rollback migrations"),
    )

    dbCmd, _ := v2.NewParentCommand[AppConfig, v2.NoFlags]("db",
        "Database operations",
        []v2.Command[AppConfig, v2.NoFlags]{migrateCmd, rollbackCmd},
        v2.WithShort[AppConfig, v2.NoFlags]("Database operations"),
    )

    v2.AddCommand(cli, dbCmd)
    cli.ExecuteAndExit(context.Background())
}
```

### Full DI Example

See [`examples/typed/main.go`](examples/typed/main.go) for a complete example with DI, lifecycle hooks, typed flags, and nested commands.

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
- [CLI Design Principles](docs/CLI_DESIGN_PRINCIPLES.md) - Design guidelines

## License

MIT
