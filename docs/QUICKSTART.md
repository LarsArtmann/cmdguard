# Quickstart Guide

**Learn cmdguard in 5 minutes.**

cmdguard is a Go library for building validated CLI applications with type-safe flags, dependency injection, and 16 output formats.

---

## Installation

```bash
go get github.com/larsartmann/cmdguard/v4
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

    "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type AppConfig struct{}

func main() {
    cli, err := v4.NewCLI[AppConfig]("hello", "A simple CLI", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cmd, err := v4.NewCommand("greet", v4.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
            fmt.Println("Hello, World!")
            return nil
        },
        v4.WithShort("Greet someone"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if err := v4.AddCommand(cli, cmd); err != nil {
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
cmd, err := v4.NewCommand("greet", &GreetFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        msg := fmt.Sprintf("Hello, %s!", flags.Name)
        if flags.Shout {
            msg = strings.ToUpper(msg)
        }
        fmt.Println(msg)
        return nil
    },
    v4.WithShort("Greet someone"),
)
v4.AddCommand(cli, cmd)
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
migrateCmd, _ := v4.NewCommand("migrate", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
        fmt.Println("Running migrations...")
        return nil
    },
    v4.WithShort("Run migrations"),
)

rollbackCmd, _ := v4.NewCommand("rollback", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
        fmt.Println("Rolling back...")
        return nil
    },
    v4.WithShort("Rollback last migration"),
)

dbCmd, _ := v4.NewParentCommand[AppConfig]("db", "Database operations", v4.NoFlags{},
    v4.WithSubcommands(migrateCmd, rollbackCmd),
    v4.WithShort("Database operations"),
)
v4.AddCommand(cli, dbCmd)
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
cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v4.WithEnvPrefix("MYAPP_"),
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
cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{})

v4.ProvideValue(cli.Scope(), &Logger{Level: "info"})
```

### 2. Use in Handler

```go
cmd, _ := v4.NewCommand("log", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
        logger, err := v4.Invoke[*Logger](cli.Scope())
        if err != nil {
            return err
        }
        logger.Log("Hello from DI!")
        return nil
    },
    v4.WithShort("Log a message"),
)
v4.AddCommand(cli, cmd)
```

---

## Signal Handling

One-line opt-in for graceful shutdown:

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v4.WithSignalHandling(),
)
// Ctrl+C cancels the context in handlers
```

---

## Rich Output (16 Formats)

```go
import "github.com/larsartmann/go-output"

// Use go-output directly for format parsing
format := output.FormatJSON

// Render table data
headers := []string{"Name", "Age"}
rows := [][]string{{"Alice", "30"}, {"Bob", "25"}}
v4.OutputTable(format, headers, rows)
```

Supported formats: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.

---

## Key Concepts

### Constructor Pattern

All commands are created via constructors, not struct literals:

```go
// Leaf command — flags passed positionally, options are non-generic
cmd, err := v4.NewCommand("name", &Flags{}, handler, opts...)

// Parent command — subcommands added via WithSubcommands option
parent, err := v4.NewParentCommand[T]("name", "long desc", v4.NoFlags{},
    v4.WithSubcommands(childCmd1, childCmd2),
)

// Add to CLI
v4.AddCommand(cli, cmd)
```

### Type Parameters

| Parameter | Description     | Example                     |
| --------- | --------------- | --------------------------- |
| `T`       | App config type | `AppConfig`                 |
| `F`       | Flags type      | `*GreetFlags`, `v4.NoFlags` |

Type parameters are inferred from the flags argument — you rarely need to specify them explicitly.

### NoFlags

Use `v4.NoFlags` for commands without flags:

```go
cmd, _ := v4.NewCommand("version", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
        fmt.Println("v4.0.0")
        return nil
    },
)
```

### Error Handling

All functions return errors — no panics in library code:

```go
cli, err := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
if err != nil {
    return err
}
```

---

## Next Steps

- **[API Reference](../README.md)** - Full API documentation
- **[examples/](../examples/)** - Complete working examples
- **[Feature Status](../FEATURES.md)** - All features and their status
