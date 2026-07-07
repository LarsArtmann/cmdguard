# Migration Guide: v2 → v3

> **TL;DR:** Change every `cmdguard/v2` import to `cmdguard/v3`, drop redundant
> type parameters from options, pass flags positionally to `NewCommand`, and import
> the relevant sub-module for spinner/telemetry/glamour/manpage/prompts features.

**Source version:** `v2.10.2` (the last stable v2)
**Target version:** `v3.0.0` (`github.com/larsartmann/cmdguard/v3`)
**Why:** v3 eliminates the type-parameter explosion, extracts heavy deps into
optional sub-modules, and fixes a semver violation (the breaking redesign was
briefly mis-tagged `v2.11.0`).

---

## 1. Update the module path

The import path gains a `/v3` suffix (Go major-version rule).

```diff
- v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
+ v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
```

Then update all qualified usages (`v2.NewCLI` → `v3.NewCLI`) and your `go.mod`:

```bash
go mod edit -require github.com/larsartmann/cmdguard/v3@v3.0.0
go mod tidy
```

Automated bulk replacement (from the migration itself):

```bash
find . -type f \( -name '*.go' -o -name 'go.mod' \) ! -path './.git/*' \
  -exec sed -i 's|github.com/larsartmann/cmdguard/v2|github.com/larsartmann/cmdguard/v3|g' {} +
find . -type f -name '*.go' ! -path './.git/*' \
  -exec sed -i 's/v2 "github.com\/larsartmann\/cmdguard\/v3/v3 "github.com\/larsartmann\/cmdguard\/v3/g' {} +
find . -type f -name '*.go' ! -path './.git/*' \
  -exec sed -i 's/\bv2\.\([A-Z]\)/v3.\1/g' {} +
```

---

## 2. Command API: non-generic options + positional flags

This is the biggest behavioral change. In v2 every option needed all type
parameters repeated. In v3 **metadata options are non-generic** and flags are
passed positionally to `NewCommand`, enabling full type inference.

### `NewCommand` — flags are now the 2nd positional argument

```go
// v2 — 2 explicit type params + WithFlags option
cmd, _ := v2.NewCommand[AppConfig, *ListFlags]("list", handler,
    v2.WithShort[AppConfig, *ListFlags]("List tasks"),
    v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
)

// v3 — flags passed positionally; options take zero type params
cmd, _ := v3.NewCommand("list", &ListFlags{}, handler,
    v3.WithShort("List tasks"),
)
```

**`WithFlags` is deleted** — pass flags as the second argument to `NewCommand`.

### Options that lost their type parameters

These are now plain non-generic functions (`func(...) CLIOption` /
`func(...) CommandOption`):

| v2 (with type params)          | v3 (non-generic)            |
| ------------------------------ | --------------------------- |
| `WithShort[T,F](...)`          | `WithShort(...)`            |
| `WithLong[T,F](...)`           | `WithLong(...)`             |
| `WithExample[T,F](...)`        | `WithExample(...)`          |
| `WithExactArgs[T,F](n)`        | `WithExactArgs(n)`          |
| `WithNoArgs[T,F]()`            | `WithNoArgs()`              |
| `WithCompletion[T,F](fn)`      | `WithCompletion(fn)`        |
| `WithStrictValidation[T]()`    | `WithStrictValidation()`    |
| `WithDraconianValidation[T]()` | `WithDraconianValidation()` |
| `WithConfigFile[T](paths...)`  | `WithConfigFile(paths...)`  |
| `WithAuditLog[T](plugin)`      | `WithAuditLog(plugin)`      |
| `WithPlugin[T](plugin)`        | `WithPlugin(plugin)`        |
| `WithSignalHandling[T]()`      | `WithSignalHandling()`      |
| `WithGracefulShutdown[T]()`    | `WithGracefulShutdown()`    |

…**and every other metadata option.** Drop the `[...]` everywhere.

### Options that are STILL generic (need `[T]`)

These remain generic because their callback signature needs the config type:

- `WithMiddleware[T](mw...)`
- `WithConfigValidation[T](fn)`
- `WithPostFlagParse[T](fns...)`
- `WithCleanup[T](fns...)`
- `WithSubcommands[T,F](cmds...)`
- `WithPreRunE[T,F](...)` / `WithPostRunE[T,F](...)`

> **Rule of thumb:** if the option takes a callback that receives `*T` (your
> config), it stays generic. Everything else is non-generic.

---

## 3. Sub-module adoption

Five features moved out of core into independently importable sub-modules.
**Core no longer depends on their libraries** — import a sub-module only when you
need it.

| Feature          | v2 (core)                   | v3 sub-module                                         |
| ---------------- | --------------------------- | ----------------------------------------------------- |
| Markdown help    | `v2.WithGlamourHelp[T]()`   | `glamour.WithHelp()`                                  |
| Terminal spinner | `v2.SpinnerMiddleware[T]`   | `spinner.Middleware[T](title)`                        |
| OpenTelemetry    | `v2.WithTelemetry[T]`       | `telemetry.WithTelemetry[T](tracer)`                  |
| Man pages        | `v2.GenerateManPageCommand` | `manpage.GenerateCommand[T](cli)`                     |
| huh prompts      | (built-in)                  | `prompts` package (`HuhRunner`) + `SetPromptRunner()` |

### Example: spinner + telemetry

```go
import (
    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
    "github.com/larsartmann/cmdguard/spinner"
    "github.com/larsartmann/cmdguard/telemetry"
)

cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    v3.WithMiddleware(
        spinner.Middleware[Config]("Loading..."),
        telemetry.Middleware[Config](tracer),
    ),
)
```

### Example: glamour markdown help

```go
import (
    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
    "github.com/larsartmann/cmdguard/glamour"
)

cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    glamour.WithHelp(),
)
```

### Example: interactive prompts (huh)

Core still defines the `prompt:"..."` tag, `WithPromptOnMissing()`, and the
`PromptRunner` interface. Wire the huh implementation from the sub-module:

```go
import (
    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
    "github.com/larsartmann/cmdguard/prompts"
)

prompts.Register() // sets the global PromptRunner to HuhRunner
```

---

## 4. Removed features

These were deleted entirely (not extracted):

| Removed                 | Reason                                      | Replacement            |
| ----------------------- | ------------------------------------------- | ---------------------- |
| `EditInEditor()`        | Marginal feature; not a CLI-library concern | Use `os/exec` directly |
| `Result[T]`             | Sum type; not a CLI concern                 | Return `error`         |
| `Validated[T]`          | Sum type; not a CLI concern                 | `errors.Join`          |
| `Ok[T]`/`Err[T]`        | Sum-type constructors                       | —                      |
| `Valid[T]`/`Invalid[T]` | Sum-type constructors                       | —                      |

If you used `Result[T]`/`Validated[T]`, replace with idiomatic Go error returns
or bring your own sum-type package.

---

## 5. The `v2.11.0` retraction

The v3 redesign was briefly mis-tagged `v2.11.0` on the `/v2` module path — a
semver violation (breaking change on a minor bump). That tag is **deleted and
retracted**.

- `v2.10.4` (on `release/v2.10`) contains a `retract v2.11.0` directive.
- If you reference `v2.11.0`, `go get` will warn you to upgrade.
- The v2 line is frozen at `v2.10.4`; all future work is on `/v3`.

If you accidentally pulled `v2.11.0`:

```bash
go get github.com/larsartmann/cmdguard/v3@v3.0.0
# update imports /v2 -> /v3, then:
go mod tidy
```

---

## 6. Quick checklist

- [ ] Replace all `cmdguard/v2` imports with `cmdguard/v3`
- [ ] Update import alias `v2` → `v3` and all `v2.X` usages → `v3.X`
- [ ] Pass flags positionally to `NewCommand` (delete `WithFlags`)
- [ ] Drop `[T,F]` / `[T]` from all metadata options
- [ ] Keep `[T]` only on callback-bearing options (`WithMiddleware`,
      `WithConfigValidation`, `WithPostFlagParse`, `WithCleanup`,
      `WithSubcommands`, `WithPreRunE`, `WithPostRunE`)
- [ ] Import the sub-module for spinner / telemetry / glamour / manpage / prompts
- [ ] Remove usages of `EditInEditor`, `Result[T]`, `Validated[T]`
- [ ] `go build ./... && go vet ./... && go test ./...`

---

## 7. Reference migration

See the `examples/taskctl/` directory — it was fully migrated to v3 and serves
as a complete, working reference for every pattern above.
