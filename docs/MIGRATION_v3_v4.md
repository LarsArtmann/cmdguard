# Migration Guide: v3 → v4

> **TL;DR:** Change every `cmdguard/v3` import to `cmdguard/v4`, replace all
> `configload.NewKoanfLoader` / `NewJSONLoader` calls with `WithConfigFile`,
> and update import aliases from `v3` to `v4`. No API signatures changed —
> v4 is a module-path + config-loading consolidation release.

**Source version:** `v3.1.0` (the last stable v3)
**Target version:** `v4.0.0` (`github.com/larsartmann/cmdguard/v4`)
**Why:** v4 consolidates config loading into a single `KoanfLoader` (deleting
the `configload` sub-package and `NewJSONLoader`), upgrades go-output to v0.37.0,
and bumps the module path to the next major version.

---

## 1. Update the module path

The import path gains a `/v4` suffix (Go major-version rule).

```diff
- v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
+ v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
```

Then update all qualified usages (`v3.NewCLI` → `v4.NewCLI`) and your `go.mod`:

```bash
go mod edit -require github.com/larsartmann/cmdguard/v4@v4.0.0
go mod tidy
```

Automated bulk replacement:

```bash
find . -type f \( -name '*.go' -o -name 'go.mod' \) ! -path './.git/*' \
  -exec sed -i 's|github.com/larsartmann/cmdguard/v3|github.com/larsartmann/cmdguard/v4|g' {} +
find . -type f -name '*.go' ! -path './.git/*' \
  -exec sed -i 's/v3 "github.com\/larsartmann\/cmdguard\/v4/v4 "github.com\/larsartmann\/cmdguard\/v4/g' {} +
find . -type f -name '*.go' ! -path './.git/*' \
  -exec sed -i 's/\bv3\.\([A-Z]\)/v4.\1/g' {} +
```

---

## 2. Config loading: `configload` deleted, use `WithConfigFile`

This is the only **behavioral** change. In v3, config file loading required the
`configload` sub-package and explicit loader construction. In v4, all config
loading is consolidated into a single `KoanfLoader` accessed via the
`WithConfigFile` CLI option.

### Before (v3)

```go
import (
    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
    "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3/configload"
)

cli, _ := v3.NewCLI[AppConfig](
    v3.WithConfigFileLoader(configload.NewKoanfLoader("config.yaml")),
)
```

Or with the JSON loader:

```go
cli, _ := v3.NewCLI[AppConfig](
    v3.WithConfigFileLoader(v3.NewJSONLoader()),
)
```

### After (v4)

```go
import (
    v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

cli, _ := v4.NewCLI[AppConfig](
    v4.WithConfigFile("config.yaml"),
)
```

`WithConfigFile` accepts one or more paths and auto-detects JSON, YAML, and TOML
by file extension. Missing files are silently skipped. Paths support `$ENV` and
`~` expansion.

### What was deleted

| v3 API                                               | v4 Replacement                                  |
| ---------------------------------------------------- | ----------------------------------------------- |
| `configload.NewKoanfLoader(path)`                    | `WithConfigFile(path)`                          |
| `configload.NewKoanfLoaderFromBytes(data, ext)`      | Load via file path instead                      |
| `configload.NewJSONLoader()` / `NewJSONLoader()`     | `WithConfigFile(path)` (auto-detects JSON)      |
| `configload` sub-package                             | Deleted — `KoanfLoader` lives in the main package |

### Config precedence

Precedence is unchanged: explicit flag → `env:"VAR"` tag → config file → default value.

---

## 3. No other API changes

The following are **unchanged** from v3 to v4:

- `NewCLI[T]`, `NewCommand`, `NewParentCommand[T]`, `AddCommand`, `Execute` signatures
- All `CLIOption` and `CommandOption` functions
- `Command[T, F]` struct shape
- Middleware chain (`Middleware[T]`, `WithMiddleware[T]`)
- DI scope (`samber/do/v2` integration)
- Sub-module imports (`glamour`, `prompts`, `spinner`, `telemetry`)
- Output formats (16 formats via go-output)
- Error types and sentinels
- Plugin system

### Dependency upgrades

| Dependency                     | v3       | v4       |
| ------------------------------ | -------- | -------- |
| `github.com/larsartmann/go-output` | v0.35.0  | v0.37.0  |
| `samber-do-auditlog`           | v0.7.0   | v0.8.1   |

These are transitive — `go mod tidy` after updating the module path will pull
the correct versions.

---

## 4. New sub-module: `flightrecorder`

v4.0.0 introduces the `flightrecorder` sub-module (added during the v4 line, not
a v3→v4 breaking change):

```go
import (
    "github.com/larsartmann/cmdguard/flightrecorder"
)

cli, _ := v4.NewCLI[AppConfig](
    flightrecorder.WithFlightRecorder[AppConfig](flightrecorder.Config{
        Enabled:         true,
        CaptureOnSlow:   true,
        SlowThreshold:   5 * time.Second,
        CaptureOnError:  true,
    }),
)
```

This is optional — only import it if you need execution trace recording.

---

## 5. Verification checklist

After migrating:

```bash
go build ./...      # Compiles
go test ./...       # Tests pass
go mod tidy         # go.mod / go.sum clean
```

- [ ] All imports changed from `cmdguard/v3` to `cmdguard/v4`
- [ ] No references to `configload` sub-package
- [ ] No references to `NewJSONLoader`
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
