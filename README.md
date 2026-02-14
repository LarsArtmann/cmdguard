# cmdguard

**A compile-time guard for Cobra CLI applications.**

Make it impossible to create broken CLI commands. Catch errors at construction time, not at runtime.

```go
// This will PANIC at init time - impossible to ignore
root := cmdguard.New("myapp", "My CLI")
root.AddCommand(&cobra.Command{
    Use: "broken",  // PANIC: no Run or RunE handler!
})

// This will PANIC - flag accessed before declaration
root.Flags().GetString("undeclared")  // PANIC: flag not registered
```

## The Problem

Cobra lets you create broken commands that fail at runtime:

```go
// This compiles fine, fails at runtime
root.AddCommand(&cobra.Command{Use: "sub"})  // No handler!

// User sees this error when they run the command:
// "Error: unknown command sub"
```

## The Solution

Guard your commands. Fail fast at construction time:

```go
// This PANICS immediately - impossible to ship broken code
root := cmdguard.New("myapp", "My CLI")
root.AddCommand(&cobra.Command{Use: "sub"})  // PANIC: command "sub" has no handler
```

## Installation

```bash
go get github.com/larsartmann/cmdguard
```

## Usage

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    // Create guarded root command
    root := cmdguard.New("myapp", "My application")

    // Add command - panics if invalid
    root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            // Safe: handler is required
        },
    })

    // Execute - safe to run
    root.Execute(context.Background())
}
```

## What Gets Guarded

| Violation | Result | When Caught |
|-----------|--------|-------------|
| Command without handler | **PANIC** | At `AddCommand()` |
| Duplicate command name | **PANIC** | At `AddCommand()` |
| Flag accessed before registration | **PANIC** | At `Flags().Get*()` |
| Duplicate flag name | **PANIC** | At flag registration |
| Conflicting aliases | **PANIC** | At alias registration |
| Subcommand without parent handler | **PANIC** | At validation |

## Philosophy

**Fail fast, fail loud.**

- **NOT a framework** - We don't manage your CLI
- **NOT a validator** - We don't return errors you might ignore
- **IS a guard** - We panic immediately on invalid construction

This follows the Go proverb: "Make the zero value useful." The zero value of a guarded command is a valid command.

## Why Panic?

Because broken commands are **programmer errors**, not runtime errors:

1. **Impossible to ignore** - Unlike returned errors
2. **Catches bugs early** - At construction, not execution
3. **Self-documenting** - Stack trace shows exactly where the bug is
4. **Zero runtime overhead** - No validation needed at execution

## Comparison

### Without cmdguard (framework approach)

```go
app := cmdguard.New()
app.Initialize()

// Returns error - might be ignored
if err := app.AddCommand(&cobra.Command{Use: "bad"}); err != nil {
    log.Println(err)  // Might be missed!
}

// Validate later - might be skipped
if err := app.Validate(); err != nil {
    panic(err)  // Too late!
}

app.Execute()
```

### With cmdguard (guard approach)

```go
root := cmdguard.New("myapp", "My app")

// PANIC immediately - impossible to ignore or skip
root.AddCommand(&cobra.Command{Use: "bad"})  // PANIC!

// Never reaches here if invalid
root.Execute()
```

## Advanced Usage

### Strict Mode

```go
root := cmdguard.New("myapp", "My app", cmdguard.WithStrictMode())

// Also guards against:
// - Flags without descriptions
// - Commands without usage text
// - Deprecated patterns
```

### Custom Guards

```go
root := cmdguard.New("myapp", "My app")
root.AddGuard(func(cmd *cobra.Command) error {
    if strings.Contains(cmd.Use, "-") {
        return fmt.Errorf("commands should not contain hyphens: %s", cmd.Use)
    }
    return nil
})
```

## Project Status

**⚠️ WORK IN PROGRESS**

This project is being transformed from a framework to a guard library. The public API will change significantly.

See [docs/planning/](docs/planning/) for the transformation plan.

## License

MIT
