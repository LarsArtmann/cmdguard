# Migrating from v1 to v2

This guide helps you migrate from the v1 API (`pkg/cmdguard`) to the v2 API (`pkg/cmdguard/v2`).

## Quick Reference

| Feature            | v1                    | v2                     |
| ------------------ | --------------------- | ---------------------- |
| **Command Type**   | `*cobra.Command`      | `Command[T, F]`        |
| **Flags**          | Cobra flags API       | Struct tags            |
| **Config**         | Fixed `config.Config` | Custom type `T`        |
| **Error Handling** | Panics on invalid     | Returns errors         |
| **DI**             | None                  | Built-in `Scope`       |
| **Type Safety**    | Runtime only          | Compile-time + runtime |

## Key Changes

### 1. Import Path

```go
// v1
import "github.com/larsartmann/cmdguard/pkg/cmdguard"

// v2
import v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
```

### 2. CLI Creation

**v1:**

```go
root := cmdguard.New("myapp", "My application")
```

**v2:**

```go
type AppConfig struct {
    LogLevel string `flag:"log-level" short:"l" default:"info" help:"Log level"`
}

defaults := AppConfig{LogLevel: "info"}
root, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My application", defaults)
if err != nil {
    log.Fatal(err)
}
```

### 3. Adding Commands

**v1:**

```go
root.AddCommand(&cobra.Command{
    Use:   "greet",
    Short: "Greet someone",
    Run: func(cmd *cobra.Command, args []string) {
        name, _ := cmd.Flags().GetString("name")
        fmt.Printf("Hello, %s!\n", name)
    },
})
```

**v2:**

```go
type GreetFlags struct {
    Name string `flag:"name" short:"n" default:"World" help:"Name to greet"`
}

greetCmd := v2.Command[AppConfig, GreetFlags]{
    Use:   "greet",
    Short: "Greet someone",
    Flags: &GreetFlags{},
    RunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        fmt.Printf("Hello, %s!\n", flags.Name)
        return nil
    },
}

if err := root.AddCommand(greetCmd); err != nil {
    log.Fatal(err)
}
```

### 4. Flag Definitions

**v1 (Cobra flags):**

```go
cmd := &cobra.Command{Use: "greet"}
cmd.Flags().StringP("name", "n", "World", "Name to greet")
cmd.Flags().BoolP("shout", "s", false, "Shout the greeting")
```

**v2 (Struct tags):**

```go
type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}
```

### 5. Error Handling

**v1 (Panics):**

```go
// This panics if command is invalid
root.AddCommand(&cobra.Command{Use: "invalid"}) // No handler!
```

**v2 (Returns errors):**

```go
// This returns an error
err := root.AddCommand(v2.Command[AppConfig, v2.NoFlags]{Use: "invalid"})
if err != nil {
    // Handle error gracefully
    log.Printf("invalid command: %v", err)
}
```

### 6. Command Handlers

**v1:**

```go
Run: func(cmd *cobra.Command, args []string) {
    // Access flags via cmd.Flags().GetString()
}

RunE: func(cmd *cobra.Command, args []string) error {
    // Return error instead of panic
}
```

**v2:**

```go
RunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
    // ctx: context for cancellation
    // cfg: typed application config
    // flags: typed command flags (already parsed!)
    return nil
}
```

### 7. Lifecycle Hooks

**v1:**

```go
cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
    return nil
}
cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
    return nil
}
```

**v2:**

```go
v2.Command[AppConfig, GreetFlags]{
    PreRunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        return nil
    },
    RunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        return nil
    },
    PostRunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        return nil
    },
}
```

### 8. Subcommands

**v1:**

```go
parent := &cobra.Command{Use: "parent"}
child := &cobra.Command{Use: "child", Run: func(...){}}
parent.AddCommand(child)
root.AddCommand(parent)
```

**v2:**

```go
parentCmd := v2.Command[AppConfig, v2.NoFlags]{
    Use: "parent",
    Commands: []v2.Command[AppConfig, v2.NoFlags]{
        {Use: "child", RunE: func(...){...}},
    },
}
root.AddCommand(parentCmd)
```

### 9. Dependency Injection

**v1:** Not available

**v2:**

```go
// Register services
scope := root.Scope()
v2.ProvideValue(scope, &MyService{})

// Invoke in handler
RunE: func(ctx context.Context, cfg *AppConfig, flags Flags) error {
    svc, err := v2.Invoke[MyService](scope)
    if err != nil {
        return err
    }
    return svc.DoSomething(ctx)
}
```

### 10. Execution

**v1:**

```go
root.Execute(context.Background())
// or
root.ExecuteAndExit(context.Background())
```

**v2:**

```go
root.Execute(context.Background())
// or
root.ExecuteAndExit(context.Background())

// Don't forget cleanup
defer root.Shutdown(context.Background())
```

## Mixed Flag Types

If subcommands need different flag types than the root, use `AddAnyCommand`:

```go
type RootFlags struct {
    Verbose bool `flag:"verbose" short:"v" default:"false"`
}

type GreetFlags struct {
    Name string `flag:"name" short:"n" default:"World"`
}

root, _ := v2.New[AppConfig, RootFlags]("myapp", "...", defaults)

// GreetFlags differs from RootFlags
greetCmd := v2.Command[AppConfig, GreetFlags]{
    Use:   "greet",
    Flags: &GreetFlags{},
    RunE:  func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error { ... },
}

// Use AddAnyCommand for different flag types
v2.AddAnyCommand[AppConfig, RootFlags, GreetFlags](root, greetCmd)
```

## Functional Options Pattern

v2 supports functional options for command configuration:

```go
cmd, err := v2.NewCommand[AppConfig, GreetFlags]("greet",
    v2.WithShort("Greet someone"),
    v2.WithLong("Send a greeting to the specified person"),
    v2.WithExample("myapp greet --name Alice"),
    v2.WithFlags(&GreetFlags{}),
    v2.WithRunE(func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
        fmt.Printf("Hello, %s!\n", flags.Name)
        return nil
    }),
    v2.WithAliases("hello", "hi"),
)
```

## Common Migration Issues

### Issue: Command has no handler

**v1 behavior:** Panics at `AddCommand()`

**v2 behavior:** Returns error from `AddCommand()`

**Fix:** Ensure every leaf command has a `RunE` handler:

```go
// This is invalid - no RunE
v2.Command[AppConfig, v2.NoFlags]{Use: "greet"}

// This is valid - has RunE
v2.Command[AppConfig, v2.NoFlags]{
    Use: "greet",
    RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        return nil
    },
}
```

### Issue: Flag not found

**v1:** Access via `cmd.Flags().GetString("name")`

**v2:** Access directly from typed `flags` parameter

**Fix:** Define flags struct and use it in handler:

```go
type Flags struct {
    Name string `flag:"name"` // Required tag
}

RunE: func(ctx context.Context, cfg *AppConfig, flags Flags) error {
    fmt.Println(flags.Name) // Direct access
}
```

### Issue: Config access

**v1:** `root.Config()` returns fixed `*config.Config`

**v2:** Config is your custom type `*T`, passed to handlers

**Fix:** Define your config type and access via handler parameter:

```go
type AppConfig struct {
    DatabaseURL string `flag:"db-url" default:"localhost:5432"`
}

RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    // ...
}
```

## Complete Migration Example

### Before (v1)

```go
package main

import (
    "context"
    "fmt"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    root := cmdguard.New("myapp", "My application")

    greetCmd := &cobra.Command{
        Use:   "greet [name]",
        Short: "Greet someone",
        Run: func(cmd *cobra.Command, args []string) {
            name, _ := cmd.Flags().GetString("name")
            shout, _ := cmd.Flags().GetBool("shout")

            msg := fmt.Sprintf("Hello, %s!", name)
            if shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
        },
    }
    greetCmd.Flags().StringP("name", "n", "World", "Name to greet")
    greetCmd.Flags().BoolP("shout", "s", false, "Shout the greeting")

    root.AddCommand(greetCmd)
    root.ExecuteAndExit(context.Background())
}
```

### After (v2)

```go
package main

import (
    "context"
    "fmt"
    "strings"

    v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct {
    LogLLevel string `flag:"log-level" short:"l" default:"info" help:"Log level"`
}

type GreetFlags struct {
    Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}

func main() {
    defaults := AppConfig{LogLevel: "info"}

    root, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My application", defaults)
    if err != nil {
        panic(err)
    }

    greetCmd := v2.Command[AppConfig, GreetFlags]{
        Use:   "greet [name]",
        Short: "Greet someone",
        Flags: &GreetFlags{},
        RunE: func(ctx context.Context, cfg *AppConfig, flags GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
    }

    if err := root.AddCommand(greetCmd); err != nil {
        panic(err)
    }

    root.ExecuteAndExit(context.Background())
}
```

## When to Stay on v1

Stay on v1 if:

- You need direct cobra command access
- You prefer panics over error returns
- You don't need type-safe flags
- You're using the built-in `config.Config` struct

Migrate to v2 if:

- You want compile-time type safety
- You prefer struct tags for flags
- You need custom config types
- You want integrated dependency injection
- You prefer error returns over panics
