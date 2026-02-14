# cmdguard

A CLI validation library that ensures every command and flag is actually implemented. Combines **fang** (Cobra styling), **koanf** (configuration), and **samber/do/v2** (DI) with compile-time and runtime validation.

## Features

- **Command Validation**: Ensures all declared commands have handlers
- **Flag Validation**: Verifies all flags are properly bound
- **Health Checks**: Built-in health checking for all services
- **Graceful Shutdown**: Proper cleanup with context support
- **Dependency Injection**: Uses samber/do/v2 for clean DI
- **Styled Output**: Uses Charm's fang for beautiful CLI output
- **Configuration**: Koanf-based config with env/file/flag support

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

## Quick Start

### Using the Public API

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
)

func main() {
    app := cmdguard.New()

    if err := app.Initialize(); err != nil {
        panic(err)
    }

    // Validate before running
    if err := app.Validate(); err != nil {
        panic(err)
    }

    app.ExecuteAndExit(context.Background())
}
```

### Manual Setup

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/internal/commands"
    "github.com/larsartmann/cmdguard/internal/config"
    "github.com/larsartmann/cmdguard/internal/di"
    "github.com/larsartmann/cmdguard/internal/validation"
    "github.com/samber/do/v2"
)

func main() {
    // Create DI module
    module := di.NewModule()

    // Register services
    module.ProvideServices()

    // Get services
    cfg := module.MustInvokeConfig()
    registry := module.MustInvokeRegistry()
    validator := module.MustInvokeValidator()

    // Link and validate
    registry.SetValidator(validator)

    // Run
    ctx := context.Background()
    registry.ExecuteAndExit(ctx)
}
```

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    fang     │────▶│    cobra    │◀────│   koanf     │
│   (style)   │     │   (commands)│     │  (config)   │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   cmdguard  │
                    │ (validator) │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ samber/do   │
                    │    (DI)     │
                    └─────────────┘
```

## DI Scope Hierarchy

```
RootScope (application lifecycle)
├── "config" scope (koanf, env vars, flags)
│   └── Config service (eager - load immediately)
├── "commands" scope (cobra command tree)
│   ├── CommandRegistry service
│   ├── Validator service
│   └── Individual command scopes (per subcommand)
│       └── Command-specific services
└── "output" scope (fang styling, logging)
    └── Theme service
```

## Validation Levels

| Level        | Target            | Validation           |
| ------------ | ----------------- | -------------------- |
| Compile-time | Command structs   | Interface compliance |
| Startup      | Flag definitions  | Schema matching      |
| Runtime      | Command execution | Handler exists       |
| Runtime      | Flag access       | Flag registered      |

## Commands

- `validate` - Run validation on the command tree
- `version` - Print version information
- `example` - Example command with flags

## Configuration

Configuration is loaded from (in order of precedence):

1. Command-line flags
2. Environment variables (prefix: `CMDGUARD_`)
3. Config file (default: `config.yaml`)

### Config File Example

```yaml
strict_mode: true
log_level: debug
config_file: "custom-config.yaml"
```

### Environment Variables

```bash
export CMDGUARD_STRICT_MODE=true
export CMDGUARD_LOG_LEVEL=debug
export CMDGUARD_CONFIG_FILE=/path/to/config.yaml
```

## Project Structure

```
cmdguard/
├── cmd/
│   └── cmdguard/
│       └── main.go          # Entry point
├── internal/
│   ├── config/
│   │   └── provider.go      # Koanf integration
│   ├── di/
│   │   └── module.go        # DI bindings
│   ├── validation/
│   │   ├── registry.go      # Command/flag tracking
│   │   └── validator.go     # Validation logic
│   └── commands/
│       └── root.go          # Cobra command tree
├── pkg/
│   └── cmdguard/
│       └── public_api.go    # Public interface
├── go.mod
└── README.md
```

## Key Patterns

### Eager Config (load at startup)

```go
do.ProvideEager(injector, config.NewConfig)
```

### Lazy Command Registration

```go
do.Provide(injector, commands.NewRegistry)
```

### Transient Validators (per-validation instance)

```go
do.ProvideTransient(injector, validation.NewFlagValidator)
```

### Child Scopes for Commands

```go
func RegisterSubcommand(parent do.Injector, name string) do.Injector {
    child := parent.Scope(name)
    do.Provide(child, NewSubcommandHandler)
    return child
}
```

## License

MIT
