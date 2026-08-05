# What This Project Is

cmdguard is a Go library that wraps the [Cobra](https://github.com/spf13/cobra) CLI framework with [fang](https://github.com/charmbracelet/fang) for styled error output, adding construction-time validation, type-safe flags, and dependency injection.

---

## Purpose

cmdguard ensures CLI commands are valid at construction time rather than failing silently at runtime. It provides a single API:

| API              | Use Case                         |
| ---------------- | -------------------------------- |
| **v4** (current) | Type-safe, DI-powered, no panics |

**Module path:** `github.com/larsartmann/cmdguard/v4`

---

## Core Features

### 1. Validation at Construction Time

Commands are validated when they're added to the CLI, not when executed. This catches misconfigurations during development.

```go
// Every constructor returns an error — no panics
cli, err := v4.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})
cmd, err := v4.NewCommand("greet", &GreetFlags{}, handler,
    v4.WithShort("Greet someone"),
)
err = v4.AddCommand(cli, cmd) // error if command is invalid
```

### 2. Type-Safe Configuration

Application configuration is typed via generics:

```go
type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v"`
    Output  string `flag:"output" short:"o" default:"text"`
}

cli, err := v4.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})
```

### 3. Type-Safe Flags

Define flags using struct tags, then pass the struct positionally to `NewCommand`:

```go
type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}

cmd, err := v4.NewCommand("greet", &GreetFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        // flags is fully typed - IDE autocomplete works
        fmt.Printf("Hello, %s!\n", flags.Name)
        return nil
    },
    v4.WithShort("Greet someone"),
)
```

### 4. Dependency Injection

Built-in DI through [samber/do/v2](https://github.com/samber/do):

```go
// Register services
v4.Provide(cli.Scope(), NewDatabaseService)
v4.ProvideValue(cli.Scope(), &Logger{Level: "info"})

// Invoke in command handlers
db, err := v4.Invoke[*DatabaseService](cli.Scope())
```

### 5. Lifecycle Management

- **PreRunE / PostRunE** hooks
- **HealthCheck()** for service health checks
- **Shutdown()** for graceful cleanup

---

## What You Get

- **Fail-fast validation** — errors caught at startup, not runtime
- **Type safety** — compile-time checking of config and flags
- **DI integration** — clean service management
- **Structured flag definition** — using struct tags
- **Flag typo suggestions** — Levenshtein distance for user-friendly errors

---

## Ideal For

- Building CLI tools where correctness matters
- Teams that want strict validation at development time
- Applications requiring dependency injection patterns
- Projects wanting typed flags with IDE support
