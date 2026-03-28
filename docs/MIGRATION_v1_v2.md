# Migration Guide: v1 to v2

**From:** cmdguard v1 (`pkg/cmdguard`)  
**To:** cmdguard v2 (`pkg/cmdguard/v2`)

---

## Why Migrate?

| Aspect                     | v1     | v2                     |
| -------------------------- | ------ | ---------------------- |
| Panics on invalid commands | ✅ Yes | ❌ No (returns errors) |
| Type-safe config           | ❌ No  | ✅ Yes                 |
| Type-safe flags            | ❌ No  | ✅ Yes                 |
| DI integration             | ❌ No  | ✅ Yes                 |
| Flag struct tags           | ❌ No  | ✅ Yes                 |

**Recommendation:** Use v2 for all new projects. Migrate existing projects when convenient.

---

## Quick Comparison

### v1: Simple but Panics

```go
root := cmdguard.New("myapp", "My app")

root.AddCommand(&cobra.Command{
    Use:   "hello",
    Short: "Say hello",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello!")
    },
})

root.ExecuteAndExit(context.Background())
// ⚠️ Panics if command has no Run/RunE
```

### v2: Type-Safe and Never Panics

```go
type AppConfig struct{}

cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})
if err != nil {
    // Handle error gracefully
    return err
}

cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
    Use:   "hello",
    Short: "Say hello",
    RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        fmt.Println("Hello!")
        return nil
    },
})

cli.ExecuteAndExit(context.Background())
// ✅ Never panics, all errors handled
```

---

## Migration Steps

### 1. Update Imports

```go
// Before
import "github.com/larsartmann/cmdguard/pkg/cmdguard"

// After
import "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
```

### 2. Define Your Config Type

```go
// v1: No config type
root := cmdguard.New("myapp", "My app")

// v2: Define config struct
type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
    LogLevel string `flag:"log-level" default:"info" help:"Log level"`
}
```

### 3. Create CLI with Config

```go
// v1
root := cmdguard.New("myapp", "My app")

// v2
cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})
if err != nil {
    return err
}
```

### 4. Update Command Definitions

```go
// v1
root.AddCommand(&cobra.Command{
    Use:   "hello",
    Short: "Say hello",
    Run: func(cmd *cobra.Command, args []string) {
        name, _ := cmd.Flags().GetString("name")
        fmt.Printf("Hello, %s!\n", name)
    },
})
cmd.Flags().StringP("name", "n", "World", "Name to greet")

// v2: Define flag struct
type HelloFlags struct {
    Name string `flag:"name" short:"n" default:"World" help:"Name to greet"`
}

cli.AddCommand(v2.Command[AppConfig, *HelloFlags]{
    Use:   "hello",
    Short: "Say hello",
    Flags: &HelloFlags{},
    RunE: func(ctx context.Context, cfg *AppConfig, flags *HelloFlags) error {
        fmt.Printf("Hello, %s!\n", flags.Name)
        return nil
    },
})
```

### 5. Use Functional Options for Convenience

```go
// v2 supports builder pattern
cmd := v2.Command[AppConfig, v2.NoFlags]{
    Use:   "version",
    Short: "Print version",
    RunE:  versionHandler,
}

// Or use functional options
cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("version",
    v2.WithShort("Print version"),
    v2.WithRunE(versionHandler),
)
```

### 6. Add DI Services (v2 Only)

```go
// v2: Register services in scope
v2.ProvideValue(cli.Scope(), &Logger{})

// Get services in handlers
db, err := v2.Invoke[*Database](cli.Scope())
```

### 7. Update Execution

```go
// v1
root.ExecuteAndExit(context.Background())

// v2
cli.ExecuteAndExit(context.Background())
```

---

## Flag Tag Reference

| Tag        | Description   | Example                |
| ---------- | ------------- | ---------------------- |
| `flag`     | Flag name     | `flag:"name"`          |
| `short`    | Short flag    | `short:"n"`            |
| `default`  | Default value | `default:"World"`      |
| `help`     | Help text     | `help:"Name to greet"` |
| `required` | Required flag | `required:"true"`      |

---

## Error Handling

```go
// v1: Panics on invalid
root.AddCommand(&cobra.Command{
    Use:   "broken",
    Short: "This will panic",  // No Run or RunE!
})
// ❌ Panics at runtime!

// v2: Returns errors
cmd := v2.Command[AppConfig, v2.NoFlags]{
    Use:   "fixed",
    Short: "This returns error",
    // No RunE!
}
err := cli.AddCommand(cmd)
if err != nil {
    // ✅ Gracefully handle
    return fmt.Errorf("adding command: %w", err)
}
```

---

## Common Patterns

### Nested Commands

```go
cli.AddCommand(v2.Command[AppConfig, v2.NoFlags]{
    Use:   "db",
    Short: "Database operations",
    Commands: []v2.Command[AppConfig, v2.NoFlags]{
        {Use: "migrate", Short: "Run migrations", RunE: migrateHandler},
        {Use: "rollback", Short: "Rollback migrations", RunE: rollbackHandler},
    },
})
```

### Mixed Flag Types

```go
// Root has NoFlags
cli, _ := v2.New[AppConfig, v2.NoFlags]("myapp", "My app", AppConfig{})

// Commands can have different flag types
v2.AddAnyCommand(cli, v2.Command[AppConfig, *GreetFlags]{
    Use:   "greet",
    Flags: &GreetFlags{},
    RunE:  greetHandler,
})

v2.AddAnyCommand(cli, v2.Command[AppConfig, *ConfigFlags]{
    Use:   "config",
    Flags: &ConfigFlags{},
    RunE:  configHandler,
})
```

---

## Feature Comparison

| Feature           | v1  | v2  |
| ----------------- | --- | --- |
| Panic on invalid  | ✅  | ❌  |
| Error returns     | ❌  | ✅  |
| Typed config      | ❌  | ✅  |
| Typed flags       | ❌  | ✅  |
| Struct tag flags  | ❌  | ✅  |
| Flag suggestions  | ❌  | ✅  |
| DI integration    | ❌  | ✅  |
| Health checks     | ❌  | ✅  |
| Graceful shutdown | ❌  | ✅  |

---

## Need Help?

- See [README.md](README.md) for full API reference
- Check [examples/basic](examples/basic/) for simple usage
- Check [examples/typed](examples/typed/) for v2 patterns
- Check [examples/di](examples/di/) for DI examples
