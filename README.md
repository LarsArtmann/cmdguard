# cmdguard

**A Go library for building validated Cobra CLI applications with compile-time guards.**

This library wraps Cobra with panic-at-construction-time validation, ensuring invalid commands are caught immediately at startup rather than at runtime.

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func main() {
    // Single-step initialization - panics on invalid
    root := cmdguard.New("myapp", "My application description")

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

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    // Create guarded root command
    root := cmdguard.New("myapp", "My CLI application")

    // Add custom commands - panics if command has no handler
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

### Nested Commands

```go
func main() {
    root := cmdguard.New("myapp", "My CLI")

    // Parent command (intermediate)
    parent := &cobra.Command{
        Use:   "db",
        Short: "Database operations",
    }

    // Child commands
    parent.AddCommand(&cobra.Command{
        Use:   "migrate",
        Short: "Run migrations",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runMigrations()
        },
    })

    // Add nested structure - AddSubcommand validates children
    root.AddCommand(parent)
    root.ExecuteAndExit(context.Background())
}
```

### Strict Mode

```go
func main() {
    // Enable strict mode via environment
    os.Setenv("CMDGUARD_STRICT_MODE", "true")

    root := cmdguard.New("myapp", "My CLI")

    // In strict mode, this panics - RunE required instead of Run
    root.AddCommand(&cobra.Command{
        Use:   "strict",
        Short: "Must use RunE in strict mode",
        RunE: func(cmd *cobra.Command, args []string) error {
            return nil
        },
    })

    root.ExecuteAndExit(context.Background())
}
```

## Architecture

The framework provides:

- **Dependency Injection**: Built on `samber/do/v2` for service management
- **Configuration Management**: Integrated with `koanf` for config loading
- **Command Registry**: Centralized command management and validation
- **Lifecycle Management**: Structured init → validate → execute → shutdown flow

### Components

| Component | Purpose |
|-----------|---------|
| `Application` | Main entry point, orchestrates lifecycle |
| `Module` | DI container and service provider |
| `Registry` | Command registration and management |
| `Validator` | Runtime validation of command tree |
| `Config` | Configuration management |

## API Reference

### Application

```go
// Create new application
app := cmdguard.New()

// Initialize services
err := app.Initialize()
err := app.InitializeWithOptions(opts...)

// Validation
err := app.Validate()
app.MustValidate()  // Panics on error

// Execution
err := app.Execute(ctx)
app.ExecuteAndExit(ctx)  // Calls os.Exit

// Access components
root := app.Root()           // *cobra.Command
registry := app.Registry()   // *commands.Registry
config := app.Config()       // *config.Config
validator := app.Validator() // *validation.Validator
injector := app.Injector()   // do.Injector

// Lifecycle
err := app.Shutdown()
err := app.HealthCheck()
```

### Options

```go
// Add command during initialization
cmdguard.WithCommand(cmd *cobra.Command)

// Add validation hook
cmdguard.WithValidationHook(hook func() error)
```

## Validation

The framework validates:

- All commands have handlers (Run or RunE)
- Flags are properly bound
- No duplicate command names
- No conflicting aliases

```go
// Run validation
if err := app.Validate(); err != nil {
    log.Fatal("Validation failed:", err)
}

// Or panic on validation failure
app.MustValidate()
```

## Configuration

Configuration is loaded from:
1. Default values
2. Config file (YAML)
3. Environment variables
4. Command-line flags

```go
// Access config
config := app.Config()
if config.StrictMode {
    // Additional validations enabled
}
```

## Project Status

**⚠️ WORK IN PROGRESS**

The API is evolving. Expect breaking changes in future releases.

## License

MIT
