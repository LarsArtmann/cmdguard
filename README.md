# cmdguard

**A Go library for building validated Cobra CLI applications with panic-at-construction-time guards.**

This library wraps Cobra with validation that panics at construction time, ensuring invalid commands are caught immediately at startup rather than failing silently at runtime.

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    // Single-step initialization
    root := cmdguard.New("myapp", "My CLI application")

    // Add commands - panics if invalid (no handler)
    root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, World!")
        },
    })

    // Execute
    root.ExecuteAndExit(context.Background())
}
```

## How It Works

cmdguard validates commands at construction time:

- **Panic on invalid commands** - Commands without handlers (`Run` or `RunE`) cause immediate panic
- **Panic on empty names** - Commands must have a valid `Use` field
- **Strict mode** - Optional enforcement of `RunE` (error-returning handlers)

This "fail fast" approach catches configuration errors during development, not production.

## API Reference

### GuardedCommand

```go
// Create root command
root := cmdguard.New(name, shortDescription)

// Add subcommands (panics if invalid)
root.AddCommand(cmd *cobra.Command)

// Add nested subcommands (panics if child is invalid)
root.AddSubcommand(parent, child *cobra.Command)

// Execution
err := root.Execute(ctx)           // Returns error
root.ExecuteAndExit(ctx)           // Calls os.Exit(1) on error

// Access underlying Cobra command for advanced customization
cmd := root.Command()

// Configuration access
cfg := root.Config()
strict := root.IsStrictMode()
```

### Built-in Commands

Every GuardedCommand includes:

- `version` - Prints version information
- `validate` - Validates the entire command tree
- `help` - Cobra's built-in help

### Global Flags

- `--config, -c` - Config file path
- `--log-level, -l` - Log level (debug, info, warn, error)
- `--strict, -s` - Enable strict mode

## Examples

### Nested Commands

```go
func main() {
    root := cmdguard.New("myapp", "My CLI")

    // Parent command (intermediate - no handler needed when it has children)
    db := &cobra.Command{
        Use:   "db",
        Short: "Database operations",
    }

    // Leaf commands must have handlers
    db.AddCommand(&cobra.Command{
        Use:   "migrate",
        Short: "Run migrations",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runMigrations()
        },
    })

    root.AddCommand(db)
    root.ExecuteAndExit(context.Background())
}
```

### Strict Mode

Strict mode requires all handlers to be `RunE` (returning error):

```go
func main() {
    // Enable via environment
    os.Setenv("CMDGUARD_STRICT_MODE", "true")

    root := cmdguard.New("myapp", "My CLI")

    // This works in strict mode
    root.AddCommand(&cobra.Command{
        Use:   "check",
        Short: "Run checks",
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    })

    // This would panic in strict mode (Run instead of RunE)
    // root.AddCommand(&cobra.Command{
    //     Use: "bad",
    //     Run: func(cmd *cobra.Command, args []string) {},
    // })

    root.ExecuteAndExit(context.Background())
}
```

### Custom Flags

```go
func main() {
    root := cmdguard.New("myapp", "My CLI")

    cmd := &cobra.Command{
        Use:   "greet",
        Short: "Greet someone",
        RunE: func(cmd *cobra.Command, args []string) error {
            name, _ := cmd.Flags().GetString("name")
            fmt.Printf("Hello, %s!\n", name)
            return nil
        },
    }

    cmd.Flags().StringP("name", "n", "World", "Name to greet")
    root.AddCommand(cmd)

    root.ExecuteAndExit(context.Background())
}
```

## Configuration

Configuration via environment variables:

| Variable | Values | Default |
|----------|--------|---------|
| `CMDGUARD_LOG_LEVEL` | debug, info, warn, error | info |
| `CMDGUARD_LOG_FORMAT` | text, json | text |
| `CMDGUARD_STRICT_MODE` | true, false | false |

### JSON Logging

For machine-parseable logs, set the format to JSON:

```bash
export CMDGUARD_LOG_FORMAT=json
myapp command
# Output: {"time":"2026-02-14T10:00:00Z","level":"INFO","msg":"message"}
```

## Architecture

```
pkg/cmdguard/
└── guarded_command.go    # Public API (GuardedCommand)

internal/
├── config/               # Configuration loading
└── logging/              # slog integration
```

### Dependencies

| Library | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/fang` | Beautiful CLI output |
| `log/slog` | Structured logging (stdlib) |

## Philosophy

**Why panics?**

Go lacks compile-time macros. The closest equivalent to "fail at compile time" is "fail at init time". By panicking on invalid commands during construction:

1. Errors are caught during development, not production
2. Invalid states are impossible to represent at runtime
3. The API is simple - no error handling boilerplate

**When to use cmdguard:**

- You want guaranteed-valid CLI configurations
- You prefer "crash early" over "handle errors later"
- You're building CLIs where panics are acceptable (most CLIs)

**When NOT to use cmdguard:**

- You need to handle configuration errors gracefully
- You're embedding CLI in a larger application that can't panic

## Project Status

**v0.1.0 released.** Core API is stable. JSON logging added in v0.2.0.

## License

MIT
