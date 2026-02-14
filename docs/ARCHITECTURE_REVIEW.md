# cmdguard Architecture Review

> **Review Date:** 2026-02-14
> **Reviewer:** Crush (AI Assistant)
> **Library:** github.com/larsartmann/cmdguard
> **Version:** 0.1.0 (pre-release)

---

## Executive Summary

**cmdguard** is a Go library for building validated Cobra CLI applications with a "guard" approach - panicking at construction time if commands are invalid. It enforces correctness at initialization rather than runtime.

|| Library | Purpose |
|---------|---------|
| `spf13/cobra` | CLI command framework |
| `charmbracelet/fang` | Styled error output and execution |

### Overall Assessment

|| Criterion | Score | Notes |
|-----------|-------|-------|
| **Architecture** | 9/10 | Clean, simple, focused |
| **Code Quality** | 9/10 | Well-structured, idiomatic Go |
| **Test Coverage** | 9/10 | 91%+ coverage on main packages |
| **Documentation** | 8/10 | Good inline docs, principles doc |
| **Production Readiness** | 8/10 | Ready for early adopters |

---

## Module Structure

```
cmdguard/
├── pkg/cmdguard/           # Public API - GuardedCommand
│   ├── guarded_command.go  # Main entry point and validation
│   └── guarded_command_test.go
├── internal/
│   ├── config/             # Configuration management
│   │   ├── provider.go     # Environment-based config loading
│   │   └── provider_test.go
│   └── logging/            # Structured logging
│       ├── logger.go       # slog wrapper with format/level
│       └── logger_test.go
├── examples/               # Usage examples
│   ├── basic/main.go
│   ├── guarded/main.go
│   └── advanced/main.go
└── tests/integration/      # Integration tests
```

---

## Core Components

### 1. GuardedCommand (`pkg/cmdguard`)

The `GuardedCommand` struct wraps `cobra.Command` with compile-time validation:

```go
type GuardedCommand struct {
    cmd        *cobra.Command
    cfg        *config.Config
    logger     *slog.Logger
    validated  bool
    strictMode bool
}
```

**Lifecycle:**

```
New() → AddCommand() → Execute()
  ↓          ↓             ↓
Create   Validate+Panic   Run CLI
```

**Key Principle:** Invalid commands cause immediate panic at construction time, not runtime errors when invoked.

**Public API Surface:**

|| Method | Purpose | Panics? |
|--------|---------|---------|
| `New(name, short)` | Constructor | No |
| `AddCommand(cmd)` | Add subcommand | Yes, if invalid |
| `AddSubcommand(parent, child)` | Add nested subcommand | Yes, if invalid |
| `Execute(ctx)` | Run command | No |
| `ExecuteAndExit(ctx)` | Run and exit | No |
| `Command()` | Get underlying cobra.Command | No |
| `Config()` | Get configuration | No |
| `IsStrictMode()` | Check strict mode | No |
| `Version()` | Get version string | No |

---

### 2. Configuration (`internal/config`)

Simple environment-based configuration:

```go
type Config struct {
    StrictMode bool   // Enable strict RunE requirement
    ConfigFile string // Path to config file
    LogLevel   string // debug, info, warn, error
    LogFormat  string // text, json
}
```

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `CMDGUARD_LOG_LEVEL` | `info` | Log level |
| `CMDGUARD_LOG_FORMAT` | `text` | Output format |
| `CMDGUARD_STRICT_MODE` | `false` | Enable strict mode |

---

### 3. Logging (`internal/logging`)

Slog wrapper with configurable format and level:

```go
func NewLogger(format, level string) *slog.Logger
func ValidFormat(format string) bool
```

---

## Validation Rules

### Command Validation

A command is valid if:

1. **Has a name** - `cmd.Name() != ""`
2. **Has handler OR subcommands:**
   - Leaf commands: must have `Run` or `RunE`
   - Parent commands: can omit handler if they have subcommands
3. **Strict mode:** Requires `RunE` (returns error) instead of `Run`

**Example - Valid:**

```go
// Leaf command with handler
&cobra.Command{
    Use: "deploy",
    Run: func(cmd *cobra.Command, args []string) { ... },
}

// Parent with subcommands
parent := &cobra.Command{Use: "server"}
parent.AddCommand(&cobra.Command{Use: "start", Run: ...})
```

**Example - Invalid (will panic):**

```go
// No handler, no subcommands
&cobra.Command{Use: "broken"}  // PANIC!
```

---

## Built-in Commands

Every GuardedCommand includes:

| Command | Purpose |
|---------|---------|
| `version` | Print version information |
| `validate` | Validate entire command tree |

---

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `""` | Config file path |
| `--log-level` | `-l` | `info` | Log level (debug/info/warn/error) |
| `--strict` | `-s` | `false` | Enable strict validation |

---

## Test Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `pkg/cmdguard` | 91.0% | ✅ Excellent |
| `internal/config` | 82.6% | ✅ Good |
| `internal/logging` | 100% | ✅ Complete |

---

## Usage Example

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    // Create root - single step initialization
    root := cmdguard.New("myapp", "My application")

    // Add commands - panics if invalid
    root.AddCommand(&cobra.Command{
        Use:   "greet",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            println("Hello!")
        },
    })

    // Execute
    root.Execute(context.Background())
}
```

---

## Design Decisions

### Why Panic Instead of Error?

Go lacks compile-time macros. Panicking at construction time is the closest to "fail fast":
- Errors caught immediately at startup
- No silent failures at runtime
- Forces developers to fix issues before deployment

### Why No DI Framework?

The Guard API is intentionally simple:
- No external DI container needed
- Configuration loaded inline
- Services created at construction time
- Minimal dependencies

### Why Fang Instead of Raw Cobra?

`charmbracelet/fang` provides:
- Beautiful error styling
- Consistent output formatting
- Better UX out of the box

---

## Production Readiness Checklist

### Completed ✅

- [x] Clean public API
- [x] Comprehensive test coverage (90%+)
- [x] Built-in validation commands
- [x] Environment configuration
- [x] Structured logging
- [x] Example applications

### Recommended Enhancements

- [ ] Shell completion generation
- [ ] Custom validation hooks
- [ ] Metrics/telemetry integration
- [ ] Config file support (YAML/TOML)

---

## Conclusion

cmdguard implements a clean "guard" approach to CLI construction:

- ✅ Simple, focused API
- ✅ Fail-fast validation
- ✅ Excellent test coverage
- ✅ Minimal dependencies
- ✅ Good documentation

**Verdict:** 8/10 production readiness. Ready for early adopters.

---

*Review generated by Crush AI Assistant*
*Assisted-by: Crush via crush@charm.land*
