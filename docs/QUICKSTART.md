# Quickstart Guide

**Learn cmdguard in 5 minutes.**

cmdguard is a Go library for building validated CLI applications with type-safe flags, dependency injection, and 12 output formats.

---

## Installation

```bash
go get github.com/larsartmann/cmdguard/v3
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

    "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type AppConfig struct{}

func main() {
    cli, err := v3.NewCLI[AppConfig]("hello", "A simple CLI", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cmd, err := v3.NewCommand[AppConfig, v3.NoFlags]("greet",
        func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
            fmt.Println("Hello, World!")
            return nil
        },
        v3.WithShort[AppConfig, v3.NoFlags]("Greet someone"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if err := v3.AddCommand(cli, cmd); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cli.ExecuteAndExit(context.Background())
}
```

### 2. Run It

```bash
go run main.go greet
# Output: Hello, World!

go run main.go --help
# Shows styled usage help
```

---

## Adding Flags (1 minute)

### 1. Define Flag Struct

```go
type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}
```

### 2. Use in Command

```go
cmd, err := v3.NewCommand[AppConfig, *GreetFlags]("greet",
    func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        msg := fmt.Sprintf("Hello, %s!", flags.Name)
        if flags.Shout {
            msg = strings.ToUpper(msg)
        }
        fmt.Println(msg)
        return nil
    },
    v3.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v3.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v3.AddCommand(cli, cmd)
```

### 3. Run It

```bash
go run main.go greet --name=Alice
# Output: Hello, Alice!

go run main.go greet -n Bob -s
# Output: HELLO, BOB!

go run main.go greet --help
# Shows flag help
```

---

## Subcommands (1 minute)

```go
migrateCmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("migrate",
    func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
        fmt.Println("Running migrations...")
        return nil
    },
    v3.WithShort[AppConfig, v3.NoFlags]("Run migrations"),
)

rollbackCmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("rollback",
    func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
        fmt.Println("Rolling back...")
        return nil
    },
    v3.WithShort[AppConfig, v3.NoFlags]("Rollback last migration"),
)

dbCmd, _ := v3.NewParentCommand[AppConfig, v3.NoFlags]("db",
    "Database operations",
    []v3.Command[AppConfig, v3.NoFlags]{migrateCmd, rollbackCmd},
    v3.WithShort[AppConfig, v3.NoFlags]("Database operations"),
)
v3.AddCommand(cli, dbCmd)
```

```bash
go run main.go db migrate
go run main.go db rollback
```

---

## Environment Variables

Flags can read from environment variables with the `env` tag:

```go
type DBFlags struct {
    Host     string `flag:"host"     env:"DB_HOST"     default:"localhost" help:"Database host"`
    Port     int    `flag:"port"     env:"DB_PORT"     default:"5432"      help:"Database port"`
    Password string `flag:"password" env:"DB_PASSWORD" default:""          help:"Database password"`
}
```

Priority: explicit flag > env var > default.

Add a prefix with `WithEnvPrefix`:

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v3.WithEnvPrefix[AppConfig]("MYAPP_"),
)
// Now reads MYAPP_DB_HOST, MYAPP_DB_PORT, etc.
```

---

## Counting Flags

Use `count:"true"` for `-v`/`-vv`/`-vvv` verbosity:

```go
type Flags struct {
    Verbose int `flag:"verbose" short:"v" help:"Verbosity level" count:"true"`
}
```

```bash
myapp run          # Verbose = 0
myapp run -v       # Verbose = 1
myapp run -vvv     # Verbose = 3
```

---

## Dependency Injection (1 minute)

### 1. Register Services

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{})

v3.ProvideValue(cli.Scope(), &Logger{Level: "info"})
```

### 2. Use in Handler

```go
cmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("log",
    func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
        logger, err := v3.Invoke[*Logger](cli.Scope())
        if err != nil {
            return err
        }
        logger.Log("Hello from DI!")
        return nil
    },
    v3.WithShort[AppConfig, v3.NoFlags]("Log a message"),
)
v3.AddCommand(cli, cmd)
```

---

## Signal Handling

One-line opt-in for graceful shutdown:

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v3.WithSignalHandling[AppConfig](),
)
// Ctrl+C cancels the context in handlers
```

---

## Rich Output (16 Formats)

```go
output "github.com/larsartmann/go-output"

data := v3.DefaultOutputConfig()
data.Format = output.FormatJSON

// Or parse from string:
format, _ := output.ParseFormat("yaml")

// Render table data
headers := []string{"Name", "Age"}
rows := [][]string{{"Alice", "30"}, {"Bob", "25"}}
v3.OutputTable(output.FormatTable, headers, rows)
```

Supported formats: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.

---

## Key Concepts

### Constructor Pattern

All commands are created via constructors, not struct literals:

```go
// Leaf command
cmd, err := v3.NewCommand[T, F]("name", handler, opts...)

// Parent command
parent, err := v3.NewParentCommand[T, F]("name", "desc", subcommands, opts...)

// Add to CLI
v3.AddCommand(cli, cmd)
```

### Type Parameters

| Parameter | Description     | Example                     |
| --------- | --------------- | --------------------------- |
| `T`       | App config type | `AppConfig`                 |
| `F`       | Flags type      | `*GreetFlags`, `v3.NoFlags` |

### NoFlags

Use `v3.NoFlags` for commands without flags:

```go
cmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("version",
    func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
        fmt.Println("v2.2.0")
        return nil
    },
)
```

### Error Handling

All v2 functions return errors — no panics in library code:

```go
cli, err := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
if err != nil {
    return err
}
```

---

## Next Steps

- **[API Reference](../README.md)** - Full API documentation
- **[examples/](../examples/)** - Complete working examples
- **[Feature Status](../FEATURES.md)** - All features and their status
