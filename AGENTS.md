# AGENTS.md - cmdguard Project Guide

**Last Updated:** 2026-02-28
**Project:** cmdguard - CLI Guard Library
**Go Version:** 1.26.1
**Status:** v2.0.0 COMPLETE - Production Ready

---

## Quick Start

```bash
# Run tests
go test ./...
go test -v -cover ./...        # Verbose with coverage
go test -race ./...            # Race detection

# Build examples
cd examples/typed && go build -o myapp .
```

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with type-safe dependency injection.

**Two APIs Available:**

| API                  | Package           | Use Case                         |
| -------------------- | ----------------- | -------------------------------- |
| **v2** (Recommended) | `pkg/cmdguard/v2` | Type-safe, DI-powered, no panics |
| v1 (Legacy)          | `pkg/cmdguard`    | Simple, panic-at-construction    |

**Current Status:** v2.0.0 Complete. All packages tested with 90%+ coverage.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   ├── v2/                      # v2 API (recommended)
│   │   ├── errors.go            # Typed errors
│   │   ├── types.go             # Common types (Enum, Duration, LogLevel)
│   │   ├── config.go            # Configuration merging/validation
│   │   ├── flags.go             # FlagRegistry with struct tags
│   │   ├── flags_parse.go       # Flag parsing
│   │   ├── flags_suggest.go     # Typo suggestions
│   │   ├── scope.go             # DI scope with samber/do/v2
│   │   ├── command.go           # Command[T] definition
│   │   ├── guard.go             # GuardedCommand[T]
│   │   ├── guard_command.go     # Command conversion
│   │   ├── guard_exec.go        # Execution logic
│   │   └── guard_flags.go       # Flag setup
│   └── guarded_command.go       # v1 API
├── internal/
│   ├── config/                  # Configuration (95.7% coverage)
│   └── logging/                 # Structured logging (100% coverage)
├── examples/
│   ├── basic/                   # v1 API demo
│   └── typed/                   # v2 API demo with DI
├── docs/
│   ├── architecture.d2          # D2 diagram source
│   └── CLI_DESIGN_PRINCIPLES.md
├── AGENTS.md                    # This file
├── FEATURES.md                  # Feature status
├── TODO_LIST.md                 # Remaining tasks
├── .golangci.yml                # Lint configuration
└── README.md                    # User documentation
```

### Package Guidelines

| Package            | Purpose           | Importable? | Coverage |
| ------------------ | ----------------- | ----------- | -------- |
| `pkg/cmdguard/v2`  | v2 Type-safe API  | Yes         | 90.6%    |
| `pkg/cmdguard`     | v1 Guard API      | Yes         | 94.3%    |
| `internal/config`  | Configuration     | No          | 95.7%    |
| `internal/logging` | Logging utilities | No          | 100%     |

---

## Key Dependencies

| Library                         | Purpose              | Version |
| ------------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`        | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`       | Dependency injection | v2.0.0  |
| `github.com/charmbracelet/fang` | Cobra styling        | v2.0.1  |
| `github.com/knadh/koanf/v2`     | Configuration        | v2.3.3  |
| `github.com/onsi/ginkgo/v2`     | BDD testing          | v2.28.1 |
| `github.com/onsi/gomega`        | Test matchers        | v1.39.1 |

---

## v2 API (Recommended)

### Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// Define your config - single source of truth
type AppConfig struct {
    LogLevel  v2.LogLevel `flag:"log-level" short:"l" default:"info" help:"Log level"`
    LogFormat string      `flag:"log-format" default:"text" help:"Log format"`
}

func main() {
    ctx := context.Background()

    // Create CLI with typed config
    root, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My application", AppConfig{})
    if err != nil {
        panic(err)
    }

    // Add command
    err = root.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
        Use:   "hello",
        Short: "Say hello",
        RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            fmt.Printf("Hello! Log level: %s\n", cfg.LogLevel)
            return nil
        },
    })
    if err != nil {
        panic(err)
    }

    // Execute
    if err := root.Execute(ctx); err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Dependency Injection with samber/do/v2

```go
package main

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
    "github.com/samber/do/v2"
)

type DatabaseService struct {
    db *sql.DB
}

// Compile-time interface verification
var _ do.Shutdowner = (*DatabaseService)(nil)
var _ do.HealthcheckerWithContext = (*DatabaseService)(nil)

func NewDatabaseService(i do.Injector) (*DatabaseService, error) {
    // Use MustInvoke for required dependencies
    cfg := do.MustInvoke[*AppConfig](i)

    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }

    return &DatabaseService{db: db}, nil
}

func (d *DatabaseService) Shutdown(ctx context.Context) error {
    return d.db.Close()
}

func (d *DatabaseService) HealthCheck(ctx context.Context) error {
    return d.db.PingContext(ctx)
}

func main() {
    root, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})

    // Register services in DI scope
    v2.Provide(root.ScopeStruct(), NewDatabaseService)

    // Add command that uses DI
    root.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
        Use:   "check",
        Short: "Health check",
        RunE: func(ctx context.Context, cfg *AppConfig, _ v2.NoFlags) error {
            // Get service from DI
            db, err := v2.Invoke[*DatabaseService](root.ScopeStruct())
            if err != nil {
                return err
            }

            // Use service
            return db.db.PingContext(ctx)
        },
    })

    // Health check before starting
    if err := root.HealthCheckWithContext(context.Background()); err != nil {
        panic(err)
    }

    // Execute
    root.ExecuteAndExit(context.Background())

    // Shutdown on exit
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    root.Shutdown(ctx)
}
```

### DI Helper Functions

```go
// Provide - Register a service provider
v2.Provide(scope, NewDatabaseService)

// ProvideNamed - Register named service (multiple implementations)
v2.ProvideNamed[Cache](scope, "redis", NewRedisCache)
v2.ProvideNamed[Cache](scope, "memory", NewMemoryCache)

// ProvideValue - Register existing value
v2.ProvideValue(scope, &MyService{})

// Invoke - Get service (returns error)
svc, err := v2.Invoke[*DatabaseService](scope)

// MustInvoke - Get service (panics on error, for constructors)
svc := v2.MustInvoke[*DatabaseService](scope)

// InvokeNamed - Get named service
redis, err := v2.InvokeNamed[Cache](scope, "redis")

// MustInvokeNamed - Get named service (panics on error)
redis := v2.MustInvokeNamed[Cache](scope, "redis")

// HealthCheck - Run health checks
err := scope.HealthCheck()

// HealthCheckWithContext - Run health checks with context
err := scope.HealthCheckWithContext(ctx)

// Shutdown - Graceful shutdown
err := scope.Shutdown(ctx)

// Child - Create child scope
child := scope.Child("plugin-scope")
```

### BranchingFlowContext

`BranchingFlowContext` tracks the command execution path and allows context values to flow through the command tree. It's automatically created when `Execute` is called and accessible via `GetBranchingFlowContext`.

```go
// BranchingFlowContext is automatically integrated with CLI[T]
// Access it after Execute
root, _ := v2.NewCLI[AppConfig]("myapp", "My app", AppConfig{})

// Access flow context in command handlers
root.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
    Use:   "check",
    Short: "Check flow context",
    RunE: func(ctx context.Context, cfg *AppConfig, _ v2.NoFlags) error {
        // Get branching flow context from context
        bfc, ok := v2.GetBranchingFlowContext(ctx)
        if !ok {
            return errors.New("no flow context")
        }

        // Access path and values
        fmt.Println("Path:", bfc.PathString())
        fmt.Println("Depth:", bfc.Depth())

        return nil
    },
})

root.ExecuteAndExit(context.Background())

// After Execute, access via CLI method
fc := root.FlowContext()
```

**Key Functions:**

| Function | Purpose |
|----------|---------|
| `NewBranchingFlowContext(ctx)` | Create root context |
| `GetBranchingFlowContext(ctx)` | Get from context (returns ok) |
| `RequireBranchingFlowContext(ctx)` | Get or panic |
| `WithBranchingFlowContext(ctx, bfc)` | Wrap context with flow context |
| `bfc.Branch(name)` | Create child context for subcommand |
| `bfc.PathString()` | Get dot-separated path (e.g., "app.subcmd") |
| `bfc.SetValue(key, val)` | Set value (propagates to children) |
| `bfc.GetValue(key)` | Get value (looks up hierarchy) |
| `bfc.Cancel()` | Cancel this and all children |

**Use Cases:**
- Track which commands are executing (audit logging)
- Propagate request-scoped values down command tree
- Cancel all child operations on error
- Debug command execution paths

### Command with Custom Flags

Supported flag types: `string`, `bool`, `int`, `uint`, `float64`, `[]string`, `Duration`, `Enum`, `LogLevel`, `LogFormat`

```go
type GreetFlags struct {
    Name    string   `flag:"name"    short:"n" default:"World" help:"Name to greet"`
    Count   uint     `flag:"count"   short:"c" default:"1"    help:"Number of greetings"`
    Shout   bool     `flag:"shout"   default:"false"         help:"Shout the greeting"`
}

root.AddCommand(v2.Command[AppConfig, GreetFlags]{
    Use:   "greet",
    Short: "Greet someone",
    Flags: GreetFlags{}, // Provide defaults
    RunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        for i := uint(0); i < flags.Count; i++ {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
        }
        return nil
    },
})
```

### Subcommands

```go
parent := v2.Command[AppConfig, v2.NoFlags]{
    Use:   "user",
    Short: "User management",
    Commands: []v2.Command[AppConfig, v2.NoFlags]{
        {
            Use:   "list",
            Short: "List users",
            RunE:  listUsersHandler,
        },
        {
            Use:   "create",
            Short: "Create user",
            RunE:  createUserHandler,
        },
    },
}

root.AddCommand(parent)
```

### Error Handling

```go
// All v2 functions return errors
root, err := v2.New[Config, NoFlags]("app", "My app", Config{})
if err != nil {
    // Handle initialization error
}

// Check specific errors with errors.Is
if errors.Is(err, v2.ErrInvalidCommand) {
    // Handle invalid command
}

if errors.Is(err, v2.ErrDuplicateCommand) {
    // Handle duplicate command name
}

// Available sentinel errors:
// - ErrInvalidCommand
// - ErrMissingHandler
// - ErrMissingName
// - ErrDuplicateCommand
// - ErrFlagParseFailed
// - ErrConfigValidation
// - ErrInvalidScope
// - ErrServiceNotFound
// - ErrServiceConstruction
// - ErrServiceRegistration
// - ErrInvalidEnum
// - ErrInvalidDuration
// - ErrInvalidFlagType
// - ErrConfigNil
// - ErrConfigNotPointer
// - ErrFlagNotFound
// - ErrRequiredFlag
```

---

## v1 API (Legacy)

### Basic Usage

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    root := cmdguard.New("myapp", "My application")

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

### Panic Conditions (Intentional)

The v1 Guard API panics at construction time if:

1. Command has no `Run`/`RunE` and no subcommands
2. Strict mode requires `RunE` but only `Run` provided
3. Command has no name

---

## Coding Standards

### Go Conventions

- **Go 1.26.1** - Use modern Go features
- **gofumpt** formatting preferred
- **Error handling** - Always check errors, wrap with context
- **Interface naming** - `-er` suffix (e.g., `Validator`)
- **Constructor naming** - `New` + type name (e.g., `NewLogger`)
- **File size** - Max 250 lines (split if larger)
- **Function size** - Max 30 lines (extract if larger)

### Testing

Use Ginkgo/Gomega for BDD-style tests:

```go
package v2_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

var _ = Describe("Command", func() {
    Describe("Validate", func() {
        It("returns error for empty Use field", func() {
            cmd := v2.Command[struct{}, struct{}]{}
            err := cmd.Validate()
            Expect(err).To(HaveOccurred())
            Expect(errors.Is(err, v2.ErrInvalidCommand)).To(BeTrue())
        })
    })
})
```

### Test Commands

```bash
# Run all tests with coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Race detection
go test -race ./...

# Specific package
go test -v ./pkg/cmdguard/v2/...
```

---

## Configuration

Uses [knadh/koanf](https://github.com/knadh/koanf) for configuration management per library policy.

### Configuration Sources (Priority Order)

1. **Environment variables** (highest priority)
2. **Config file** (YAML)
3. **Default values** (lowest priority)

### Environment Variables

Prefix: `CMDGUARD_`

- `CMDGUARD_LOG_LEVEL` - debug, info, warn, error (default: info)
- `CMDGUARD_LOG_FORMAT` - text, json (default: text)
- `CMDGUARD_STRICT_MODE` - true/false (default: false)

### Config File (YAML)

```yaml
log_level: debug
log_format: json
strict_mode: true
```

### Programmatic Usage

```go
package main

import (
    "github.com/larsartmann/cmdguard/internal/config"
)

func main() {
    // Create loader
    loader := config.NewLoader()

    // Load configuration (optional config file path)
    err := loader.Load("config.yaml")
    if err != nil {
        panic(err)
    }

    // Get values
    logLevel := loader.GetString("log_level")
    strictMode := loader.GetBool("strict_mode")

    // Or unmarshal to struct
    type Config struct {
        LogLevel   string `koanf:"log_level"`
        LogFormat  string `koanf:"log_format"`
        StrictMode bool   `koanf:"strict_mode"`
    }

    var cfg Config
    err = loader.Unmarshal(&cfg)
}
```

### Priority Example

```yaml
# config.yaml
log_level: warn
```

```bash
# Environment overrides file
export CMDGUARD_LOG_LEVEL=debug

# Final value: debug (env wins)
```

---

## Architecture Decisions

### v2 Design Principles

1. **Type Safety** - Generic type parameters for config and flags
2. **No Panics** - All operations return errors
3. **DI-Powered** - samber/do/v2 for dependency injection
4. **Typed Flags** - Struct tags for flag definitions
5. **Composable** - Commands are values, easy to test

### Why samber/do/v2?

- Clean API with no global state
- Scope support for nested DI containers
- Lifecycle hooks (Shutdowner, Healthchecker)
- Compile-time dependency checking
- Context-aware operations

### Error Handling Strategy

- Sentinel errors for `errors.Is()` checking
- Wrapped errors with context
- No panics in library code
- Rich error types (CommandError, FlagError, ServiceError)

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
