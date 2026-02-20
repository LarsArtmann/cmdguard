# cmdguard v2 Design: Type-Safe, DI-Powered CLI Framework

**Status:** DRAFT  
**Date:** 2026-02-15  
**Goal:** Rebuild cmdguard with single source of truth, no panics, and DI scopes

---

## Problem Statement

### Current Issues (v1)

| Issue               | Current                               | Impact                       |
| ------------------- | ------------------------------------- | ---------------------------- |
| **Panics**          | `AddCommand` panics on invalid        | Library crashes applications |
| **No Single Truth** | Flags defined separately from config  | Drift, duplication           |
| **No Type Safety**  | Flags return `string`, manual parsing | Runtime errors               |
| **No DI**           | Manual service wiring                 | Hard to test, hard to extend |
| **Static**          | Can't add plugin config dynamically   | Limits extensibility         |

### Design Goals

1. **Single Source of Truth** - Define config once, derive everything else
2. **No Panics** - All functions return errors
3. **Type Safety** - Strong types, branded IDs, make invalid states unrepresentable
4. **DI with Scopes** - samber/do/v2 for lifecycle management
5. **Dynamic Registration** - Plugins can register their own config/commands

---

## Architecture

### Scope Hierarchy

```
Root Scope (Application)
├── Config[T]           # Typed application config
├── Logger              # Structured logging
├── FlagProvider[T]     # Generates flags from config
├── CommandRegistry     # All commands
│
├── Plugin Scope (per-plugin)
│   ├── PluginConfig    # Plugin-specific config
│   └── PluginServices  # Plugin dependencies
│
└── Command Scope (per-execution)
    ├── CommandConfig   # Resolved flags for this command
    └── CommandServices # Command dependencies
```

### Core Types

```go
// Single source of truth - define your config struct
type AppConfig struct {
    LogLevel   config.Enum `flag:"log-level" short:"l" default:"info" values:"debug,info,warn,error"`
    LogFormat  config.Enum `flag:"log-format" default:"text" values:"text,json"`
    StrictMode bool        `flag:"strict" short:"s" default:"false"`
    PluginDir  string      `flag:"plugin-dir" default:"./plugins"`
    Timeout    duration    `flag:"timeout" default:"30s"`
}

// Branded types for type safety
type PluginName string // Can't mix with regular string

// Duration type with validation
type Duration time.Duration
```

---

## API Design

### 1. Application Bootstrap (Root Scope)

```go
package main

import (
    "context"

    "github.com/larsartmann/cmdguard/v2/pkg/cmdguard"
    "github.com/samber/do/v2"
)

// Define your config - SINGLE SOURCE OF TRUTH
type Config struct {
    LogLevel   cmdguard.Enum `flag:"log-level" short:"l" default:"info" values:"debug,info,warn,error"`
    LogFormat  cmdguard.Enum `flag:"log-format" default:"text" values:"text,json"`
    StrictMode bool          `flag:"strict" short:"s" default:"false"`
}

func main() {
    ctx := context.Background()

    // Create root scope with config
    root := cmdguard.New("myapp", "My application", Config{})

    // Register services in root scope
    do.Provide(root, NewLogger)
    do.Provide(root, NewDatabase)

    // Build command tree (returns errors, never panics)
    if err := root.AddCommand(greetCmd()); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Execute - injects typed config everywhere
    if err := root.Execute(ctx); err != nil {
        os.Exit(1)
    }
}
```

### 2. Command with Type-Safe Config Access

```go
func greetCmd() *cmdguard.Command[Config] {
    return cmdguard.Command[Config]{
        Use:   "greet",
        Short: "Greet someone",
        Flags: struct {
            Name    string `flag:"name" short:"n" default:"World"`
            Shout   bool   `flag:"shout" default:"false"`
        }{},
        RunE: func(ctx context.Context, cfg *Config, flags *GreetFlags) error {
            // cfg is typed! cfg.LogLevel is Enum, not string
            // flags is typed! flags.Name is string

            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }

            fmt.Println(msg)
            return nil
        },
    }
}
```

### 3. Plugin with Dynamic Config Registration

```go
// Plugin registers its own config scope
func RegisterPlugin(scope do.Injector) error {
    // Define plugin-specific config
    type PluginConfig struct {
        Enabled  bool   `flag:"enabled" default:"true"`
        Endpoint string `flag:"endpoint" default:"http://localhost:8080"`
    }

    // Register in plugin scope (child of root)
    pluginScope := do.Scope(scope, "my-plugin")

    do.Provide(pluginScope, func(i do.Injector) (PluginConfig, error) {
        // Auto-populated from flags with prefix: --my-plugin-enabled, --my-plugin-endpoint
        return cmdguard.ResolveConfig[PluginConfig](i, "my-plugin")
    })

    // Register plugin services
    do.Provide(pluginScope, NewPluginService)

    return nil
}
```

---

## Core Interfaces

### GuardedCommand (No Panics!)

```go
package cmdguard

import (
    "context"
    "github.com/samber/do/v2"
    "github.com/spf13/cobra"
)

// GuardedCommand provides type-safe CLI construction with DI
type GuardedCommand[T any] struct {
    name    string
    short   string
    config  T
    scope   do.Injector
    cmd     *cobra.Command
}

// New creates a new CLI application with typed config
func New[T any](name, short string, defaults T) *GuardedCommand[T]

// AddCommand adds a subcommand - returns error instead of panic
func (g *GuardedCommand[T]) AddCommand(cmd Command[T]) error

// AddCommandFunc adds a command using a constructor function
func (g *GuardedCommand[T]) AddCommandFunc(fn func() Command[T]) error

// Execute runs the CLI
func (g *GuardedCommand[T]) Execute(ctx context.Context) error

// ExecuteAndExit runs the CLI and exits with appropriate code
func (g *GuardedCommand[T]) ExecuteAndExit(ctx context.Context)

// Scope returns the DI scope for service registration
func (g *GuardedCommand[T]) Scope() do.Injector

// Config returns the resolved configuration
func (g *GuardedCommand[T]) Config() *T
```

### Command Definition

```go
// Command represents a type-safe command with typed flags
type Command[T any] struct {
    Use     string                              // Command name
    Short   string                              // Short description
    Long    string                              // Long description
    Aliases []string                            // Aliases
    Example string                              // Example usage

    // Flags is populated from struct tags
    Flags   any                                 // Struct with flag tags

    // RunE receives typed config and typed flags
    RunE    func(ctx context.Context, cfg *T, flags any) error

    // PreRunE for validation
    PreRunE func(ctx context.Context, cfg *T, flags any) error

    // Subcommands
    Commands []Command[T]
}
```

### Flag Provider (Auto-Generation)

```go
// FlagProvider generates Cobra flags from struct tags
type FlagProvider[T any] struct {
    config T
}

// RegisterFlags adds flags to a cobra command
func (p *FlagProvider[T]) RegisterFlags(cmd *cobra.Command) error

// ParseFlags populates config from parsed flags
func (p *FlagProvider[T]) ParseFlags(cmd *cobra.Command) (*T, error)

// Struct tag format:
// `flag:"name" short:"n" default:"value" help:"description" required:"true"`
```

---

## Error Handling (No Panics!)

```go
// All errors are typed and wrapped
var (
    ErrInvalidCommand    = errors.New("invalid command")
    ErrMissingHandler    = errors.New("command has no handler")
    ErrMissingName       = errors.New("command has no name")
    ErrFlagParseFailed   = errors.New("failed to parse flags")
    ErrConfigValidation  = errors.New("config validation failed")
)

// AddCommand returns error instead of panic
func (g *GuardedCommand[T]) AddCommand(cmd Command[T]) error {
    if cmd.Use == "" {
        return fmt.Errorf("%w: command has no name", ErrInvalidCommand)
    }

    if cmd.RunE == nil && len(cmd.Commands) == 0 {
        return fmt.Errorf("%w: %q has no RunE and no subcommands", ErrMissingHandler, cmd.Use)
    }

    // ... rest of implementation
    return nil
}
```

---

## DI Scope Usage

```go
// Register services in root scope
func main() {
    root := cmdguard.New("app", "desc", Config{})

    // Services available everywhere
    do.Provide(root.Scope(), NewLogger)
    do.Provide(root.Scope(), NewDatabase)
    do.Provide(root.Scope(), NewCache)

    // Command-specific scope
    root.AddCommand(cmdguard.Command[Config]{
        Use: "users",
        RunE: func(ctx context.Context, cfg *Config, flags any) error {
            // Access services via DI with proper error handling
            db, err := v2.Invoke[*Database](root.ScopeStruct())
            if err != nil {
                return v2.NewServiceError("*Database", err)
            }
            users, err := db.ListUsers(ctx)
            // ...
        },
    })
}
```

---

## Type Safety Examples

### Enum Type

```go
// Enum provides type-safe enum values
type Enum struct {
    value   string
    allowed []string
}

// Only allows valid values
func ParseEnum(value string, allowed []string) (Enum, error) {
    for _, a := range allowed {
        if value == a {
            return Enum{value: value, allowed: allowed}, nil
        }
    }
    return Enum{}, fmt.Errorf("invalid value %q, must be one of: %v", value, allowed)
}

// Usage in config
type Config struct {
    LogLevel Enum `flag:"log-level" values:"debug,info,warn,error"`
}
```

### Duration Type

```go
// Duration parses and validates time durations
type Duration struct {
    time.Duration
}

func ParseDuration(s string) (Duration, error) {
    d, err := time.ParseDuration(s)
    if err != nil {
        return Duration{}, fmt.Errorf("invalid duration: %w", err)
    }
    return Duration{Duration: d}, nil
}
```

---

## Migration Path (v1 → v2)

### Step 1: Add DI Scope Support

- Add samber/do/v2 dependency
- Wrap existing GuardedCommand in DI scope
- Keep existing API working

### Step 2: Add Typed Config

- Create `TypedGuardedCommand[T]`
- Implement flag auto-generation from struct tags
- Add `AddCommand` that returns error

### Step 3: Remove Panics

- Replace all `panic()` calls with error returns
- Update all `AddCommand` signatures
- Update documentation

### Step 4: Deprecate Old API

- Mark v1 API as deprecated
- Provide migration guide
- Eventually remove in v3

---

## File Structure

```
cmdguard/
├── pkg/cmdguard/
│   ├── command.go          # Command[T] definition
│   ├── guard.go            # GuardedCommand[T] implementation
│   ├── config.go           # Typed config with struct tags
│   ├── flags.go            # Flag auto-generation
│   ├── scope.go            # DI scope helpers
│   ├── types.go            # Enum, Duration, etc.
│   └── errors.go           # Typed errors
├── internal/
│   └── validation/         # Validation logic
└── examples/
    ├── basic/              # Simple usage
    ├── typed/              # Type-safe config
    └── plugin/             # Dynamic plugin registration
```

---

## Open Questions

1. **Flag prefix for plugins?** `--plugin-name-flag` or `--plugin.flag`?
2. **Config file support?** YAML/JSON merging with flags?
3. **Env var support?** `APP_LOG_LEVEL` auto-mapped to `--log-level`?
4. **Validation framework?** `validate:"required,min=1"` tags?
5. **Subcommand config inheritance?** How do child commands inherit parent config?

---

## References

- [samber/do/v2 scopes](https://github.com/samber/do)
- [Cobra flags](https://github.com/spf13/cobra)
- [struct tag parsing](https://pkg.go.dev/reflect#StructTag)
- [go-plugin-mvp](../go-plugin-mvp) - Plugin system reference
