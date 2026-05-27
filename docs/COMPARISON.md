# cmdguard vs. Alternative Go CLI Frameworks

> An honest comparison for developers evaluating CLI frameworks.
> **Last Updated:** 2026-05-27

---

## Summary

| Framework | Stars | Cobra Integration | Struct Tags | Type Safety | DI | Typo Suggestions |
| --- | --- | --- | --- | --- | --- | --- |
| **cmdguard** | — | ✅ Native | ✅ | ✅ Compile-time | ✅ samber/do/v2 | ✅ |
| **Kong** | ~2.4k | ❌ Separate framework | ✅ | ✅ Mapper interface | ❌ | ❌ |
| **sflags** | ~167 | ✅ Native | ✅ | ❌ Runtime | ❌ | ❌ |
| **go-flags** | ~2.6k | ❌ Separate framework | ✅ | ❌ Runtime | ❌ | ❌ |
| **structcli** | ~8 | ✅ Native | ✅ | ❌ Runtime | ❌ | ❌ |
| **urfave/cli** | ~22k | ❌ Separate framework | ❌ | ❌ Runtime | ❌ | ❌ |

---

## Detailed Comparison

### cmdguard

**Best for:** Teams already using Cobra who want type safety, DI, and better UX without switching frameworks.

- ✅ **Native Cobra integration** — wraps Cobra, doesn't replace it
- ✅ **Compile-time type safety** — flags are typed structs, not string lookups
- ✅ **Built-in DI** — samber/do/v2 with `Provide`/`Invoke` and lifecycle hooks
- ✅ **Typo suggestions** — "did you mean?" for flags and subcommands
- ✅ **Constructor validation** — invalid commands caught at registration time
- ✅ **Zero panics** — all APIs return errors
- ✅ **Rich output** — 12+ formats (JSON, CSV, YAML, table, etc.)
- ✅ **Config file support** — JSON/YAML/TOML with flag/env override
- ⚠️ **Newer project** — smaller community than Kong or urfave/cli
- ⚠️ **Cobra required** — not a standalone framework

### Kong

**Best for:** Greenfield projects that want a declarative, self-contained CLI framework.

- ✅ **Mature ecosystem** — ~2.4k stars, extensive documentation
- ✅ **Type-safe mappers** — custom types via mapper interfaces
- ✅ **Plugin system** — extensible architecture
- ✅ **Validation** — built-in validators and custom validation tags
- ❌ **No Cobra integration** — separate framework, can't reuse existing Cobra commands
- ❌ **No built-in DI** — bring your own dependency injection
- ❌ **No typo suggestions** — standard "unknown flag" errors

### sflags

**Best for:** Simple Cobra apps that just want struct-tag flag generation.

- ✅ **Cobra native** — generates pflag/cobra flags from structs
- ✅ **Multi-backend** — supports urfave/cli, kingpin, not just Cobra
- ✅ **Lightweight** — minimal API surface
- ❌ **No type safety** — runtime reflection only
- ❌ **No DI** — no dependency injection support
- ❌ **No typo suggestions** — plain Cobra error handling
- ❌ **No custom types** — basic types only

### go-flags

**Best for:** Standalone CLI apps that don't need Cobra.

- ✅ **Extensive tag support** — rich struct tag options
- ✅ **Mature** — ~2.6k stars, stable API
- ✅ **Subcommands** — built-in subcommand support
- ❌ **No Cobra integration** — separate framework
- ❌ **Legacy API** — older design patterns
- ❌ **No DI** — no dependency injection

### urfave/cli

**Best for:** Simple CLIs that don't need subcommands or complex flag hierarchies.

- ✅ **Popular** — ~22k stars, widely used
- ✅ **Simple API** — easy to get started
- ✅ **Middleware** — `cli.Before`/`cli.After` hooks
- ❌ **No Cobra integration** — separate framework
- ❌ **No struct tags** — flags defined imperatively
- ❌ **Stringly-typed** — `ctx.String("name")` lookups
- ❌ **No DI** — no built-in dependency injection

---

## Feature Matrix

| Feature | cmdguard | Kong | sflags | go-flags | urfave/cli |
| --- | --- | --- | --- | --- | --- |
| **Struct tags for flags** | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Cobra integration** | ✅ Native | ❌ | ✅ Native | ❌ | ❌ |
| **Compile-time type safety** | ✅ Generics | ⚠️ Mappers | ❌ | ❌ | ❌ |
| **Required validation** | ✅ `required:"true"` | ✅ | ✅ | ✅ | ✅ |
| **Typo suggestions** | ✅ Levenshtein | ❌ | ❌ | ❌ | ❌ |
| **Custom types** | ✅ Enum, Duration, etc. | ✅ Mappers | ❌ | ⚠️ Limited | ❌ |
| **Dependency injection** | ✅ samber/do/v2 | ❌ | ❌ | ❌ | ❌ |
| **Environment variables** | ✅ `env:"VAR"` | ✅ | ✅ | ✅ | ✅ |
| **Config files** | ✅ JSON/YAML/TOML | ✅ | ❌ | ❌ | ❌ |
| **Shell completion** | ✅ | ✅ | ❌ | ❌ | ✅ |
| **Middleware** | ✅ | ✅ | ❌ | ❌ | ✅ |
| **Man page generation** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Zero panics** | ✅ | ✅ | N/A | N/A | N/A |
| **Counting flags** | ✅ `-v`/`-vv`/`-vvv` | ❌ | ❌ | ❌ | ❌ |
| **Positional args validation** | ✅ `WithExactArgs`, etc. | ✅ | ❌ | ✅ | ✅ |

---

## API Comparison: Flag Definition

### cmdguard

```go
type Flags struct {
    Name   string `flag:"name"   short:"n" default:"World" help:"Name to greet"`
    Count  int    `flag:"count"  short:"c" default:"1"     help:"Number of greetings"`
    Verbose bool  `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
}

cmd, err := v2.NewCommand[Config, *Flags]("greet", handler,
    v2.WithFlags[Config, *Flags](&Flags{}),
)
```

### Kong

```go
type CLI struct {
    Greet struct {
        Name    string `help:"Name to greet" default:"World"`
        Count   int    `help:"Number of greetings" default:"1"`
        Verbose bool   `help:"Verbose output"`
    } `cmd:"" help:"Greet someone"`
}

var cli CLI
ctx := kong.Parse(&cli)
```

### sflags

```go
type Flags struct {
    Name    string `flag:"name"   short:"n" default:"World" desc:"Name to greet"`
    Count   int    `flag:"count"  short:"c" default:"1"     desc:"Number of greetings"`
    Verbose bool   `flag:"verbose" short:"v" default:"false" desc:"Verbose output"`
}

flags := &Flags{}
cmd := &cobra.Command{Use: "greet"}
sflags.ParseTo(cmd.Flags(), flags)
```

### urfave/cli

```go
app := &cli.App{
    Commands: []*cli.Command{{
        Name:  "greet",
        Flags: []cli.Flag{
            &cli.StringFlag{Name: "name", Aliases: []string{"n"}, Value: "World"},
            &cli.IntFlag{Name: "count", Aliases: []string{"c"}, Value: 1},
            &cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
        },
        Action: func(ctx *cli.Context) error {
            name := ctx.String("name")  // string lookup
            return nil
        },
    }},
}
```

---

## When to Choose Each

**Choose cmdguard when:**
- You're already using Cobra and want incremental improvements
- You want compile-time type safety for flags
- You need built-in dependency injection
- You want typo suggestions and better error messages for users
- You're building a complex CLI with config files, middleware, and lifecycle management

**Choose Kong when:**
- You're starting a new project and don't mind adopting a full framework
- You want a declarative, tag-driven API without Cobra
- You need a plugin system for extensibility

**Choose sflags when:**
- You just want struct-tag flag generation for Cobra
- You don't need DI, custom types, or typo suggestions
- You want to stay close to raw Cobra

**Choose urfave/cli when:**
- You're building a simple CLI without many subcommands
- You prefer an imperative API over struct tags
- You value a large community and extensive examples

---

## Migration Paths

| From | To cmdguard | Effort |
| --- | --- | --- |
| Plain Cobra | Wrap root CLI, migrate commands one at a time | Low |
| sflags | Replace flag definitions with cmdguard structs | Low |
| Kong | Rewrite — different architecture | High |
| urfave/cli | Rewrite — different architecture | High |

See [MIGRATION_FROM_COBRA.md](MIGRATION_FROM_COBRA.md) for a step-by-step guide.
