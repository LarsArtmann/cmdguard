# What This Project Is Not

cmdguard is a focused library, not a full-fledged CLI framework. Here's what it deliberately avoids:

---

## Not a Replacement for Cobra + Fang

cmdguard _wraps_ Cobra and _uses_ Fang, it doesn't replace them. Under the hood, it uses `cobra.Command` for command structure and `fang.Execute` for styled error output.

```go
// You still get a cobra.Command underneath
root := cli.RootCommand() // returns *cobra.Command
```

---

## Not a Standalone CLI Application

cmdguard is a **library**, not an executable. You use it to build your own CLI tools:

```go
// This is what YOU build with cmdguard
package main

import "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"

func main() {
    cli, _ := v2.New[AppConfig]("myapp", "My CLI", AppConfig{})
    cli.ExecuteAndExit(context.Background())
}
```

---

## Not Graceful Error Handling (v1)

The v1 API **panics** on invalid commands. This is intentional — it's "fail fast" design. If you need to handle errors gracefully, use v2:

```go
// v1: panics
root.AddCommand(invalidCmd) // PANIC!

// v2: returns error
err := cli.AddCommand(invalidCmd) // err != nil
```

---

## Not for Embedded Scenarios (v1)

If you're embedding a CLI in a larger application that cannot tolerate panics, v1 is not suitable. Use v2 instead, which returns errors rather than panicking.

---

## Not a Code Generator

cmdguard doesn't generate code. Flags are defined via struct tags at runtime:

```go
// This is NOT code generation — it's reflection at runtime
type MyFlags struct {
    Name string `flag:"name" short:"n"`
}
```

---

## Not a Complete CLI Framework

cmdguard doesn't provide:

- Built-in config file loading (use [knadh/koanf](https://github.com/knadh/koanf) separately)
- Built-in logging (use your preferred logger)
- Database ORMs
- Web servers
- Plugin systems (planned)

---

## Not Beginner-Friendly Defaults

cmdguard assumes you understand:

- Go generics
- Dependency injection concepts
- Cobra's command model

If you need simpler defaults, consider using Cobra directly.

---

## Not Multi-Language

cmdguard is Go-only. It doesn't generate CLIs for other languages.

---

## Not a Silver Bullet

It won't fix:

- Poorly designed CLI UX
- Missing input validation in your handlers
- Race conditions in your code
- Deployment issues

---

## Summary

| What cmdguard IS        | What cmdguard IS NOT         |
| ----------------------- | ---------------------------- |
| A library               | An executable                |
| Wraps Cobra + uses fang | Replaces Cobra               |
| Type-safe               | Code-generated               |
| Fail-fast (v1)          | Graceful error handling (v1) |
| DI-powered              | A complete framework         |
| Go only                 | Multi-language               |
