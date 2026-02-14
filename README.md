# cmdguard

**A Go framework for building validated Cobra CLI applications.**

This framework provides lifecycle management, dependency injection, and validation for Cobra-based CLI applications. It wraps Cobra with structured initialization and runtime validation.

```go
package main

import (
    "context"
    "log"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func main() {
    // Create application
    app := cmdguard.New()

    // Initialize services and DI container
    if err := app.Initialize(); err != nil {
        log.Fatal(err)
    }

    // Validate command tree before execution
    if err := app.Validate(); err != nil {
        log.Fatal(err)
    }

    // Execute and exit
    app.ExecuteAndExit(context.Background())
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
    "log"

    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    app := cmdguard.New()

    if err := app.Initialize(); err != nil {
        log.Fatal(err)
    }

    // Add custom commands
    app.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, World!")
        },
    })

    if err := app.Validate(); err != nil {
        log.Fatal(err)
    }

    app.ExecuteAndExit(context.Background())
}
```

### With Options

```go
package main

import (
    "fmt"
    "log"

    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    app := cmdguard.New()

    if err := app.InitializeWithOptions(
        cmdguard.WithCommand(&cobra.Command{
            Use:   "version",
            Short: "Print version",
            Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("v1.0.0")
            },
        }),
    ); err != nil {
        log.Fatal(err)
    }
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
