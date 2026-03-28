# Quickstart Guide

**Learn cmdguard in 5 minutes.**

cmdguard is a Go library for building validated CLI applications with type-safe flags and dependency injection.

---

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

---

## Your First CLI (2 minutes)

### 1. Create main.go

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct{}

func main() {
    cli, err := v2.New[AppConfig, v2.NoFlags]("hello", "A simple CLI", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
        Use:   "greet",
        Short: "Greet someone",
        RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            fmt.Println("Hello, World!")
            return nil
        },
    })

    cli.ExecuteAndExit(context.Background())
}
```

### 2. Run It

```bash
go run main.go greet
# Output: Hello, World!

go run main.go --help
# Shows usage help
```

---

## Adding Flags (1 minute)

### 1. Define Flag Struct

```go
type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}
```

### 2. Use in Command

```go
cli.AddCommand(v2.Command[AppConfig, *GreetFlags]{
    Use:   "greet",
    Short: "Greet someone",
    Flags: &GreetFlags{},
    RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        msg := fmt.Sprintf("Hello, %s!", flags.Name)
        if flags.Shout {
            msg = fmt.Sprintf("%s!!!", msg)
        }
        fmt.Println(msg)
        return nil
    },
})
```

### 3. Run It

```bash
go run main.go greet --name=Alice
# Output: Hello, Alice!

go run main.go greet -n Bob -s
# Output: Hello, Bob!!!

go run main.go greet --help
# Shows flag help
```

---

## Subcommands (1 minute)

```go
root := v2.Command[AppConfig, v2.NoFlags]{
    Use:   "db",
    Short: "Database operations",
    Commands: []v2.Command[AppConfig, v2.NoFlags]{
        {
            Use:   "migrate",
            Short: "Run migrations",
            RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
                fmt.Println("Running migrations...")
                return nil
            },
        },
        {
            Use:   "rollback",
            Short: "Rollback last migration",
            RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
                fmt.Println("Rolling back...")
                return nil
            },
        },
    },
}

cli.AddCommand(root)
```

```bash
go run main.go db migrate
go run main.go db rollback
```

---

## Dependency Injection (1 minute)

### 1. Create Service

```go
type Logger struct {
    Level string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.Level, msg)
}
```

### 2. Register Service

```go
cli, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})

// Register service
v2.ProvideValue(cli.Scope(), &Logger{Level: "info"})
```

### 3. Use in Handler

```go
cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
    Use:   "log",
    Short: "Log a message",
    RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        // Get service from DI
        logger, err := v2.Invoke[*Logger](cli.Scope())
        if err != nil {
            return err
        }
        logger.Log("Hello from DI!")
        return nil
    },
})
```

---

## Key Concepts

### Type Parameters

| Parameter | Description     | Example                     |
| --------- | --------------- | --------------------------- |
| `T`       | App config type | `AppConfig`                 |
| `F`       | Flags type      | `*GreetFlags`, `v2.NoFlags` |

### NoFlags

Use `v2.NoFlags` for commands without flags:

```go
v2.Command[AppConfig, v2.NoFlags]{
    Use:   "version",
    RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        fmt.Println("v1.0.0")
        return nil
    },
}
```

### Error Handling

All v2 functions return errors—no panics:

```go
cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})
if err != nil {
    // Handle error
    return err
}
```

---

## Next Steps

- **[API Reference](README.md#v2-api-reference)** - Full API documentation
- **[Migration Guide](MIGRATION_v1_v2.md)** - Migrate from v1
- **[examples/typed](examples/typed/)** - Complete examples
- **[examples/di](examples/di/)** - DI examples
