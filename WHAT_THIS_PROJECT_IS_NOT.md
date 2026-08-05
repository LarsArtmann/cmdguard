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

import "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"

func main() {
    cli, _ := v4.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})
    cli.ExecuteAndExit(context.Background())
}
```

---

## Not for Embedded Scenarios (legacy v1)

The legacy v1 API **panicked** on invalid commands. The v4 API returns errors instead:

```go
// v4: returns error
err := v4.AddCommand(cli, invalidCmd) // err != nil
```

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

- Advanced config file loading beyond JSON/YAML/TOML auto-detection (`WithConfigFile` covers the three common formats via KoanfLoader; for exotic formats use `WithConfigFileLoader` with a custom [ConfigFileLoader](https://pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4#ConfigFileLoader) or [knadh/koanf](https://github.com/knadh/koanf) directly)
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
| Error-returning (v4)    | A complete framework         |
| DI-powered              | Multi-language               |
| Go only                 | A silver bullet              |
