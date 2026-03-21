# What This Project Is

cmdguard is a Go library that wraps the [Cobra](https://github.com/spf13/cobra) CLI framework with [fang](https://github.com/charmbracelet/fang) for styled error output, adding construction-time validation, type-safe flags, and dependency injection.

---

## Purpose

cmdguard ensures CLI commands are valid at construction time rather than failing silently at runtime. It provides two APIs:

| API                  | Use Case                                  |
| -------------------- | ----------------------------------------- |
| **v2** (recommended) | Type-safe, DI-powered, no panics          |
| **v1** (legacy)      | Simple wrapper with panic-at-construction |

---

## Core Features

### 1. Validation at Construction Time

Commands are validated when they're added to the CLI, not when executed. This catches misconfigurations during development.

```go
// v1 API: panics immediately if command is invalid
root.AddCommand(&cobra.Command{Use: "hello"}) // PANIC: no Run handler

// v2 API: returns error instead
err := cli.AddCommand(cmd) // error: missing handler
```

### 2. Type-Safe Configuration (v2)

Application configuration is typed via generics:

```go
type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v"`
    Output  string `flag:"output" short:"o" default:"text"`
}

cli, err := v2.New[AppConfig, NoFlags]("myapp", "My CLI", AppConfig{})
```

### 3. Type-Safe Flags (v2)

Define flags using struct tags:

```go
type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}

cmd := v2.Command[AppConfig, *GreetFlags]{
    Use:   "greet",
    Flags: &GreetFlags{},
    RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        // flags is fully typed - IDE autocomplete works
        fmt.Printf("Hello, %s!\n", flags.Name)
        return nil
    },
}
```

### 4. Dependency Injection (v2)

Built-in DI through [samber/do/v2](https://github.com/samber/do):

```go
// Register services
v2.Provide(scope, NewDatabaseService)
v2.ProvideValue(scope, &Logger{Level: "info"})

// Invoke in command handlers
db, err := v2.Invoke[*DatabaseService](scope)
```

### 5. Lifecycle Management (v2)

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
