# cmdguard

**A Go library for building correct Cobra CLI applications.**

Import this library to add compile-time and construction-time safeguards to your CLI commands. It validates commands and flags as you build them, catching errors before they reach users.

```go
import "github.com/larsartmann/cmdguard"

func main() {
    // Guarded command - validates at construction
    root, err := cmdguard.New("myapp", "My CLI")
    if err != nil {
        log.Fatal(err)
    }
    
    // This returns an error - must handle it
    if err := root.AddCommand(&cobra.Command{Use: "broken"}); err != nil {
        log.Fatal(err) // "command 'broken' has no handler"
    }
}
```

## The Problem

Cobra lets you create broken commands that fail at runtime:

```go
// This compiles fine, fails at runtime
root.AddCommand(&cobra.Command{Use: "sub"})  // No handler!

// User sees this error when they run the command:
// "Error: unknown command sub"
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
    "log"
    "github.com/larsartmann/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    // Create guarded root command
    root, err := cmdguard.New("myapp", "My application")
    if err != nil {
        log.Fatal(err)
    }

    // Add command - returns error if invalid
    if err := root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            // Safe: handler is required
        },
    }); err != nil {
        log.Fatal(err)
    }

    // Execute - safe to run
    root.Execute(context.Background())
}
```

## What Gets Guarded

| Violation | Result | When Caught |
|-----------|--------|-------------|
| Command without handler | **error** | At `AddCommand()` |
| Duplicate command name | **error** | At `AddCommand()` |
| Flag accessed before registration | **error** | At `Flags().Get*()` |
| Duplicate flag name | **error** | At flag registration |
| Conflicting aliases | **error** | At alias registration |
| Subcommand without parent handler | **error** | At validation |

## Philosophy

**Fail fast with clear errors.**

- **NOT a framework** - We don't manage your CLI
- **IS a library** - You use our types to build your CLI
- **IS a guard** - We return errors on invalid construction

This follows the Go proverb: "Make the zero value useful." The zero value of a guarded command is a valid command.

## Why Return Errors?

Because broken commands are **programmer errors**, but we respect Go conventions:

1. **Explicit error handling** - Following Go idioms
2. **Catches bugs early** - At construction, not execution
3. **Clear error messages** - Tell you exactly what's wrong
4. **Testable** - You can test error conditions

## Comparison

### Without cmdguard

```go
root := &cobra.Command{Use: "myapp"}

// No validation - compiles fine
root.AddCommand(&cobra.Command{Use: "bad"})  // No handler!

// Error only shows when user runs it
root.Execute()  // User sees: "Error: unknown command"
```

### With cmdguard

```go
root, _ := cmdguard.New("myapp", "My app")

// Returns error - must handle it
if err := root.AddCommand(&cobra.Command{Use: "bad"}); err != nil {
    log.Fatal(err)  // "command 'bad' has no handler" - caught early!
}

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
