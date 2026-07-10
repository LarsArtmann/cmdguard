# AGENTS.md - cmdguard Contributor & AI Assistant Guide

> **Note:** This file serves as both a contributor guide and context for AI-assisted development. It documents architecture decisions, API reference, coding standards, and known gotchas.

**Last Updated:** 2026-07-10
**Project:** cmdguard - CLI Guard Library
**Go Version:** 1.26
**Status:** v3.0.0 - zero panics, 87.6% coverage, 0 lint issues, 0 race conditions

---

## Quick Start

```bash
# Enter dev shell (Go 1.26, gopls, golangci-lint)
nix develop

# Run tests (all packages, with race detection)
go test ./... -count=1 -timeout 120s -race

# Build all
go build ./...

# Lint (golangci-lint 2.x)
golangci-lint run ./...

# Format (nix + go via treefmt)
nix fmt

# Coverage
go test ./... -count=1 -timeout 120s -cover

# Check everything (format check)
nix flake check
```

**BuildFlow pre-commit hook:** runs golangci-lint + formatters on commit (auto-fixes applied automatically). Commits proceed normally — no `--no-verify` needed.

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with type-safe dependency injection.

| API | Package           | Use Case                         |
| --- | ----------------- | -------------------------------- |
| v3  | `pkg/cmdguard/v3` | Type-safe, DI-powered, no panics |

**Module path:** `github.com/larsartmann/cmdguard/v3`

**Current Status:** v3.0.0. 463 test functions (1288 runs incl. subtests), 26 benchmarks, 7 fuzz targets, 87.6% coverage, 0 build errors, 0 lint issues.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   ├── v3/                       # v3 API (recommended)
│   │   ├── cli.go                # CLI[T] struct, NewCLI, AddCommand, Execute
│   │   ├── cli_accessors.go      # CLI accessor methods (Config, Scope, etc.)
│   │   ├── cli_command.go        # Internal cobra wiring (cliToCobraCommand)
│   │   ├── cli_options.go        # CLI functional options (20 in this file; more in other files)
│   │   ├── cli_output.go         # Output format flag registration/parsing, dynamic help
│   │   ├── cli_errors_json.go    # Structured JSON error output for --output=json
│   │   ├── auditlog.go           # Audit-log export (ExportAuditLog, 11 formats), query helpers
│   │   ├── command.go            # Command[T,F] struct, constructors, options, Validate
│   │   ├── command_options.go    # CommandOption functions (WithShort, WithFlags, etc.)
│   │   ├── config.go             # Config type constraint
│   │   ├── config_file.go        # ConfigFileLoader, JSON loader, WithConfigFile
│   │   ├── config_parsing.go     # ParseFlagTags, DefaultValue (recurses into nested structs)
│   │   ├── config_setfield.go    # SetField for config structs
│   │   ├── configload/           # YAML/TOML/Auto/Koanf loaders (loader.go, koanf.go)
│   │   ├── docgen.go             # GenerateDocs (markdown command-tree docs)
│   │   ├── errors.go             # Error types (CommandError, FlagError, etc.) + sentinels
│   │   ├── errors_command.go     # Command-related sentinel errors
│   │   ├── errors_config.go      # Config-related sentinel errors
│   │   ├── errors_di.go          # DI-related sentinel errors
│   │   ├── errors_flags.go       # Flag-related sentinel errors
│   │   ├── flags.go              # FlagRegistry with struct tags
│   │   ├── flags_parse.go        # Flag parsing logic
│   │   ├── flags_suggest.go      # Typo suggestions (Levenshtein) for flags + commands
│   │   ├── flags_validate.go     # Flag validation
│   │   ├── completion.go         # Shell completion support
│   │   ├── doc.go                # Package documentation
│   │   ├── flag_helpers.go       # Flag type constraints, cloning, parsing helpers
│   │   ├── flow_context.go       # BranchingFlowContext for command path tracking
│   │   ├── flow_context_access.go # Flow context helpers (typed value access)
│   │   ├── prompts.go            # PromptRunner interface (huh/v2 impl in prompts/ sub-module)
│   │   ├── scope.go              # DI scope wrapping samber/do/v2
│   │   ├── type_handler.go       # Extensible type registry
│   │   ├── output.go             # Rich output (OutputTable, OutputResult, shape-aware errors)
│   │   ├── plugin.go             # Plugin system (Plugin interface, RegisterPlugin, WithPlugin)
│   │   ├── type_handler_kinds.go # Primitive kind handlers (string/bool/int/uint/float/slice)
│   │   ├── type_handler_intwidth.go # Narrow integer overflow validation (int8/16/32, uint8/16)
│   │   ├── type_handler_custom.go # Custom type handlers (Duration/Enum/URL/Email/Port)
│   │   ├── type_helpers.go       # Generic type helpers
│   │   ├── testutil/             # Test harness utilities
│   │   ├── version.go            # VersionCommand helper
│   │   ├── doctor.go             # DoctorCommand helper
│   │   ├── types_duration.go     # Duration type
│   │   ├── types_email.go        # Email type
│   │   ├── types_enum.go         # Enum[T] type
│   │   ├── types_filepath.go     # FilePath type
│   │   ├── types_hostport.go     # HostPort type
│   │   ├── types_log.go          # LogLevel type
│   │   ├── types_port.go         # Port type
│   │   └── types_url.go          # URL type
├── glamour/                      # SUB-MODULE: markdown help rendering (charm.land/glamour/v2)
├── manpage/                      # SUB-MODULE: man page generation (mango/roff)
├── prompts/                      # SUB-MODULE: huh/v2 interactive prompt runner
├── spinner/                      # SUB-MODULE: terminal spinner middleware (lipgloss/v2)
├── telemetry/                    # SUB-MODULE: OpenTelemetry middleware
├── pkg/testutil/
│   └── panic_test_helpers.go     # Shared test assertions
├── examples/
│   └── taskctl/                   # Single superb example: production task manager CLI
│       ├── main.go                # CLI construction, DI setup, all CLI options
│       ├── commands.go            # All command definitions with options
│       ├── types.go               # Config, flags, domain types, TaskStore service
│       ├── main_test.go           # Comprehensive integration tests (~66 tests)
│       └── README.md              # Feature matrix and usage guide
├── benchmarks/                   # Performance benchmarks
├── tests/integration/            # Integration tests
├── docs/                         # Documentation
├── AGENTS.md                     # This file (enduring context for AI sessions)
├── FEATURES.md                   # Feature inventory by status
├── TODO_LIST.md                  # Short/mid-term tasks
├── ROADMAP.md                    # Long-term direction and raw ideas
├── CHANGELOG.md                  # Change history per version
├── .golangci.yml                 # Lint configuration
├── go.work                       # Go workspace (6 modules: core + 5 sub-modules)
├── flake.nix                     # Nix dev shell, formatter, checks
├── flake.lock                    # Nix lock file
└── README.md                     # User documentation
```

### Package Guidelines

| Package           | Purpose       | Importable? | Coverage |
| ----------------- | ------------- | ----------- | -------- |
| `pkg/cmdguard/v3` | Type-safe API | Yes         | ~87.3%   |
| `pkg/testutil`    | Test helpers  | Yes         | —        |

---

## Key Dependencies

| Library                                     | Purpose              | Version |
| ------------------------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`                    | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`                   | Dependency injection | v2.0.0  |
| `github.com/spf13/pflag`                    | Flag parsing         | v1.0.10 |
| `charm.land/fang/v2`                        | Cobra styling        | v2.0.1  |
| `github.com/larsartmann/go-output`          | Rich output formats  | v0.30.1 |
| `github.com/larsartmann/samber-do-auditlog` | DI audit logging     | v0.4.0  |

### Optional Sub-Module Dependencies

Each sub-module is independently importable — core has **zero** dependencies on these libraries.

| Sub-module  | Library                          | Version | Purpose                 |
| ----------- | -------------------------------- | ------- | ----------------------- |
| `glamour`   | `charm.land/glamour/v2`          | v2.0.1  | Markdown help rendering |
| `prompts`   | `charm.land/huh/v2`              | v2.0.3  | Interactive prompts     |
| `spinner`   | `charm.land/lipgloss/v2`         | v2.0.5  | Terminal spinner        |
| `telemetry` | `go.opentelemetry.io/otel/trace` | v1.44.0 | OpenTelemetry spans     |
| `manpage`   | `muesli/mango` + `mango-cobra`   | v0.2.0  | Man page generation     |

---

## API Reference

See [docs/API.md](docs/API.md) for the full API reference (constructors, options, methods, DI, middleware, error handling, version/doctor commands).

Quick reference: `NewCLI[T]`, `NewCommand` (non-generic, flags passed positionally), `NewParentCommand[T]`, `AddCommand`, `Execute`. All functions return errors — zero panics. See [pkg.go.dev](https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3) for godoc.

---

## Coding Standards

### Go Conventions

- **Go 1.26** - Use modern Go features
- **gofumpt** formatting (via `golangci-lint fmt`)
- **Error handling** - Always check errors, wrap with `fmt.Errorf("context: %w", err)`
- **No panics** in v3 library code
- **Functional options** pattern for configuration
- **Constructor pattern** - All Command creation via `NewCommand`/`NewParentCommand`, struct fields unexported

### Testing

- `t.Parallel()` in every test function and subtest (paralleltest linter)
- `//nolint:paralleltest` for tests using `t.Setenv` or capturing `os.Stdout`
- `//nolint:fatcontext` at file level for test files with context in closures
- Table-driven tests: `tests := []struct{...}` pattern
- Two test packages: `v3` (internal, access private helpers) and `v3_test` (external)

### Test Commands

```bash
go test ./... -count=1 -timeout 120s -race     # All tests with race detection
go test ./... -count=1 -timeout 120s -cover     # Coverage report
golangci-lint run ./...                          # Lint (0 issues)
go build ./...                                   # Verify build
```

---

## Architecture Decisions

### Lint Strategy

**Goal:** 0 lint issues via real code fixes, not silencing linters.

**What was fixed (code refactored):**

- `wrapcheck` in `output.go` — external errors now wrapped via `wrapIfError` helper (nil-safe)
- `funlen` in `type_handler_kinds.go` — `registerKinds()` split into 7 focused helpers
- `funlen` in `type_handler_custom.go` — `registerCustomTypes()` split into `registerEnumTypes()` + `registerValueTypes()`
- `cyclop`/`funlen` in `cli.go` — `initialize()` split into `ensureScope()` + `setupPersistentPreRun()`
- `paralleltest` in `type_handler_test.go` — added `t.Parallel()` to all test functions and subtests
- `ireturn` — `TypeHandler` and `ConfigFileLoader` added to global ireturn allow list (legitimate interface returns)

**What remains as documented design decisions (`.golangci.yml` exclusions):**

- `gochecknoglobals` for `globalTypeRegistry`, `globalValidators`, `regexCache` — package-level registries are the COW pattern's foundation (ADR principle #11); injecting them would break the public `RegisterTypeHandler`/`RegisterValidator` API
- `gochecknoglobals` for `argsKey`/`configKey` in `cli_command.go` — context keys must be package-level (Go convention)
- `ireturn` for `do.Injector` returns in `scope.go`/`cli_accessors.go` — DI library interface, intentional
- `ireturn` for `koanf` interface in `configload/koanf.go` — factory pattern for config loaders

### v3 Design Principles

1. **Single type parameter on CLI only** — `CLI[T]` parameterizes on config. `CLIOption` and `CommandOption` are **non-generic** (`func(*spec)`); per-command flag types flow through `Command[T,F]`.
2. **Non-generic options** — Metadata options (`WithShort`, `WithLong`, `WithExample`, `WithAuditLog`, `WithConfigFile`, `WithStrictValidation`, `WithPlugin`, …) take **zero** type parameters. Type safety is preserved via generic constructors that _return_ non-generic options: `WithSubcommands[T,F](...)`, `WithPreRunE[T,F](...)`, `WithConfigValidation[T](...)`, `WithPostFlagParse[T](...)`, `WithCleanup[T](...)`. This eliminates the v2 "7 type params per command" explosion.
3. **No Panics** - All operations return errors
4. **DI-Powered** - samber/do/v2 for dependency injection
5. **Typed Flags** - Struct tags for flag definitions
6. **Standalone AddCommand** - Function (not method) to support per-command flag types
7. **Env tags** - `env:"VAR_NAME"` struct tag reads from environment
8. **Counting flags** - `count:"true"` tag enables -v/-vv/-vvv pattern
9. **Signal handling** — `WithSignalHandling()` cancels context on SIGINT/SIGTERM; `WithGracefulShutdown()` additionally triggers DI service shutdown (implies the former)
10. **Rich output** - OutputTable/OutputResult with 16 formats via go-output v0.30.1 registries
11. **Copy-on-write registries** — `FlagRegistry` shares global type/validator registries via copy-on-write; reads use the shared maps, writes trigger a lazy clone. `RegisterTypeHandler()`/`RegisterValidator()` write to global defaults (visible to instances that haven't cloned); `FlagRegistry.RegisterTypeHandler()`/`FlagRegistry.RegisterFlagValidator()` trigger COW clone and write to instance-local maps
12. **Typo suggestions** - `SuggestFlag`/`SuggestCommand` with Levenshtein
13. **Full sentinel coverage** - All 40+ errors identifiable via `errors.Is()`
14. **Generic helpers** - `textMarshal[T]`/`textUnmarshal[T]`, `renderAndWrite`/`marshalAndWrite`, `branchWithCtx`
15. **Modular sub-modules** — 5 optional importable sub-modules (`glamour`, `manpage`, `prompts`, `spinner`, `telemetry`) isolate heavy dependencies; core stays lean (13 direct deps). Extension hooks: `WithHelpTransform[T]()` (markdown rendering injection point), `PromptRunner` interface + `SetPromptRunner()` (prompt injection point). Import a sub-module only when you need its feature.
16. **Audit log integration** — `WithAuditLog(plugin)` wires `samber-do-auditlog` into the DI injector; `cli.AuditLog()`/`cli.AuditLogReport()` for programmatic access; `AuditLogServiceByName`/`AuditLogFailedServices` query helpers; `ExportAuditLog[T]` supports 11 formats (html, json, ndjson, csv, tsv, mermaid, dot, d2, plantuml, tree, htmltree). No built-in subcommand — consumers export via their own flag/env pattern (e.g. `DO_AUDITLOG_ENABLED` + `AUDIT_LOG_FORMAT`)
17. **Plugin system** — `Plugin` interface bundles custom type handlers + validators; `RegisterPlugin()` applies globally, `WithPlugin()` / `FlagRegistry.RegisterPlugin()` apply per-instance
18. **Nested config structs** — `ParseFlagTags` recurses into nested structs; `FieldTag.Index` tracks the reflect path for flattened flag registration
19. **Docs generation** — `cli.GenerateDocs(w)` writes markdown documentation for the full command tree to any `io.Writer`
20. **Go workspace** — `go.work` spans 6 modules (core + 5 sub-modules) for unified local builds; `go build ./...` compiles all modules

### Key Gotchas

#### Testing & Build

- `t.Setenv` + `t.Parallel()` = panic — use `//nolint:paralleltest`
- `NoFlags` is a distinct named type (`type NoFlags struct{}`, not an alias) — use `(NoFlags{})` with parens for comparisons
- **Nested modules** — `go build ./...` from the repo root does NOT descend into the 5 sub-module directories (each has its own `go.mod`). Build/test them individually: `for m in glamour manpage prompts spinner telemetry; do (cd pkg/cmdguard/$m && go build ./... && go test ./...); done`
- `flake.nix` provides devShell + formatter + format check only (no `buildGoModule` or vet checks)

#### Cobra Behavior

- `PostRunE` is NOT called when `RunE` errors (Cobra semantics)
- `AddCommand` calls `cmd.Validate()` as defense-in-depth even though constructors already validate
- `CommandInfo.FullPath` is set via `cobra.CommandPath()` inside the handler closure, NOT at registration — empty in unit tests unless run through cobra execution

#### Flag System

- **Counting flags** — must use `int` type with `count:"true"` tag; don't reuse flag names from root config
- **Copy-on-write registries** — `FlagRegistry` shares global `typeRegistry`/`validatorRegistry` lazily via `share()` and clones only on first write via `register()`. Package-level `RegisterTypeHandler()`/`RegisterValidator()` write to global defaults (visible to instances that haven't cloned); `FlagRegistry.RegisterTypeHandler()`/`RegisterFlagValidator()` trigger the lazy clone and write to instance-local maps. Reduces NewCLI by ~48%.
- **Direct go-output usage** — Users import `output.FormatTable`, `output.FormatJSON`, etc. directly from `github.com/larsartmann/go-output`; cmdguard only re-exports the `OutputFormat = output.Format` type alias. No `ParseOutputFormat`, `SupportedFormats`, `IsFormatSupported`, or `Format*` constant re-exports.
- **Regex validation cache** — `validateRegex` caches compiled patterns in `sync.Map` (concurrency-safe; tests run in parallel)
- **Integer overflow** — `int8`/`int16`/`int32`/`uint8`/`uint16` flag values are range-checked at parse time → `ErrIntegerOverflow`
- **Iterator methods (`iter.Seq`)** — `TagsSeq()`, `FlagNamesSeq()`, `PathSeq()`, `ChildrenSeq()` are zero-allocation alternatives; the slice-returning methods return defensive copies

#### Config Files

- **Precedence** — explicit flag → `env:"VAR"` (with optional prefix) → config file → default value
- `WithConfigFile(paths...)` loads config BEFORE flag registration; config values become the new tag defaults, so flags/env still override them
- Paths support `$ENV` and `~` expansion; missing files are silently skipped
- If the config struct has a `config` flag, its value overrides the default search paths
- **Nested structs supported** — `ParseFlagTags` recurses into nested structs; `FieldTag.Index` tracks the reflect path
- **`configload.Auto()`** — tries YAML → TOML → JSON sequentially (NOT file-extension based); since JSON is valid YAML, JSON data is handled by the YAML parser first. Use `LoaderForPath()` for precise extension-based detection
- **configload internals** — all loaders use `genericLoader` with a pluggable `unmarshalFunc`; TOML import aliased as `toml` to avoid conflict with the local `cmdguard` import alias. `KoanfLoader()` handles nested config objects (e.g. `{"db":{"host":"x"}}` → `--db-host`)
- `os.UserHomeDir()` is cached via `sync.OnceValue` (`cachedHomeDir`) to eliminate redundant syscalls across multiple `~/` paths

#### Output & Styling

- **16 output formats** via go-output `v0.30.1` registries — `RenderTableData` (all 16) and `RenderAnyData` (JSON/YAML/TOML) via thread-safe `formatRegistry[T]`. `OutputTable()` uses `AddRowChecked()` for fail-fast row validation. `--output` flag help is auto-generated from `RegisteredTableDataFormats()`.
- **go-output sub-modules** — `markdown/` and `tree/` are standalone sub-modules (like `d2/`, `table/`, etc.); `output.go` imports them explicitly so `FormatMarkdown`/`FormatTree` stay available. All go-output modules are pinned at v0.30.1. The `enum` and `envdetect` sub-modules were absorbed into go-output core.
- **fang styling** — styled output by default; `--no-color` persistent flag is registered by default and sets `NO_COLOR=1` for fang; `NO_COLOR` env var also respected automatically via fang's colorprofile. `cli.NoColor()` returns true if either is set.
- **Help rendering hook** — `WithHelpTransform[T](fn)` is the core extension point for transforming command help text. The `glamour` sub-module provides a ready-made markdown transformer (see [Sub-Modules](#sub-modules) below); it is NOT imported by core.

#### Dependency Injection

- **Zero panics** — every function returns errors; no `Must*` variants exist
- **`WithGracefulShutdown()`** enables DI service shutdown on SIGINT/SIGTERM (implies `WithSignalHandling`). Services implementing `do.ShutdownerWithError` shut down in reverse invocation order. `WithSignalHandling` only cancels context and does NOT trigger DI shutdown.
- **Override + CloneScope** — `Override[T](scope, provider)` / `OverrideValue[T](scope, value)` replace services for testing. `CloneScope(scope)` copies registrations without invoked state. Pattern: clone → override → invoke.
- **`NewScopeWithOpts(name, opts)`** — creates scope with `do.InjectorOpts` (custom logging, lifecycle hooks, health check timeouts). `WithDILogging[T](logf)` is the CLI convenience option.
- **`Package[T]`** takes a pre-existing `*Scope` as the first parameter (callers create it with `NewScope(name)` first) — separates scope construction from CLI construction.
- **Shutdown error chain** — `Shutdown()` wraps with `ErrServiceConstruction` once; `ShutdownAll` collects without additional wrapping.
- **`buildInjectorOpts`** merges `diLogf` and `auditLog` into a single `*do.InjectorOpts`; returns nil (default injector) when neither is configured.

#### Validation

- **Strict** — `WithStrictValidation()` requires `WithShort` on all commands; enforced at `AddCommand` time
- **Draconian** — `WithDraconianValidation()` is strict + requires `WithExample` on leaf commands; parent commands are exempt
- **Config validation** — `WithConfigValidation[T](fn)` runs after root flag parsing but before any command handler; blocks execution on error
- **Args** — `WithExactArgs`/`WithMinimumArgs`/`WithMaximumArgs`/`WithRangeArgs` use cobra's built-in validators (run during execution, not registration); invalid args (negative n, min > max) return errors from `NewCommand`/`NewParentCommand`
- **Error aggregation** — `ValidateConfig` uses `errors.Join(append([]error{ErrConfigValidation}, errs...)...)` so individual errors are reachable via `errors.Is`
- **Sentinel wrapping** — all errors use `fmt.Errorf("%w: ...", sentinel)` for `errors.Is()` chainability
- **`ErrLogLevel`/`ErrLogFormat` chain** — `ErrLogLevel → EnumError → ErrInvalidEnum`
- **`errors.AsType` (Go 1.26)** — use `errors.AsType[*T]` instead of `errors.As(err, &v)` (consistent with `cli.go`/`output.go`)
- **Deduplicated validators** — `validateEmail`/`validateURL` delegate to `ParseEmail()`/`ParseURL()`

#### Error Handling & Exit Codes

- `ExecuteAndExit` checks for `ExitCoder`; `NewExitError(code, err)` returns `(*ExitError, error)` and validates the 0–255 range
- **Error/exit contract** — cmdguard owns error display: the error is printed exactly once (fang when enabled, cobra when disabled). `SilenceUsage` is **true by default** (kills the #1 cobra footgun: usage-on-error); use `WithoutSilenceUsage()` to re-enable usage-on-error. The error returned by `Execute` is for exit-code mapping only — consumers must NOT re-print it (that double-prints). `ExecuteAndExit` is the blessed entry point; `ExitCode(err) int` (public) supports the post-execution-work case (flush/audit/teardown before exit). `ExitCode(nil)==0`.
- **Cobra escape hatch** — `ConfigFromContext[T](ctx)` retrieves resolved config from any cobra command context (stored by PersistentPreRunE). `WithPostFlagParse[T](fns...)` registers hooks that run after flag parsing + config validation but before command handlers (replaces manual PersistentPreRunE wrapping). `RegisterLocalCommandFlags(cmd)` registers the root's `local:"true"` flags on a subcommand.
- **`WithCleanup[T]`** — registers hooks that fire after EVERY command's `RunE`, including when `RunE` errors (Cobra's `PostRunE`/`PersistentPostRunE` do NOT fire on `RunE` error). Implemented as a tree-walk at `Execute` time (`applyCleanupHooks` in `cli.go`) that wraps each command's `RunE` — so it covers BOTH cmdguard-managed `Command[T,F]` and raw cobra subcommands (escape hatch). Hook signature `func(cmd, cfg, runErr) error`; the original `runErr` is never swallowed (cleanup errors joined via `errors.Join`). Idempotent (`cleanupWired` guard) so calling `Execute` twice doesn't double-wrap. No-op when no hooks registered.
- **Scoped flags** — `local:"true"` tag: registered on owning command only, NOT inherited by subcommands. `hidden:"true"` tag: excluded from `--help` but fully functional. Both parsed via `parseBoolTags`.
- `NewScopeFromInjector` returns `(*Scope, error)` — nil injector returns error
- `SuggestFlag` returns `(string, bool)`

#### Middleware

- **Core middleware chain** — `Middleware[T]` (middleware.go) is the generic chain type. Wire middleware via `WithMiddleware[T](mw...)`. The `spinner` and `telemetry` middleware implementations live in their sub-modules (see below).

#### Sub-Modules (glamour / manpage / prompts / spinner / telemetry)

- **Import path** — each is `github.com/larsartmann/cmdguard/<name>`; import only what you need. Core has zero deps on these.
- **Directory layout is load-bearing** — each sub-module lives at the **repo root** (`<name>/`), NOT under `pkg/cmdguard/`. Go resolves a module path by finding `go.mod` at the matching directory in the repo: `github.com/larsartmann/cmdguard/telemetry` requires `telemetry/go.mod` at the repo root. The root `go.mod` `replace` directives only work locally (in the workspace); they are **ignored by downstream consumers**. Moving a sub-module under `pkg/` breaks external `go get` silently (builds still pass via workspace `replace`).
- **glamour** — provides a markdown help transformer for the `WithHelpTransform[T]()` hook; uses `RenderWithEnvironmentConfig` (checks `GLAMOUR_STYLE` env var, defaults to `"dark"`). The string `"auto"` is NOT a valid glamour theme — pass `""` for env-based detection.
- **spinner** — `SpinnerMiddleware[T]` auto-skips when `os.Stderr` is not a terminal; override with `SpinnerConfig{Writer: ...}`.
- **telemetry** — `TelemetryMiddleware[T]` starts a span per command but cannot propagate the new context to the handler (`next func() error` signature); child spans must use the original context.
- **prompts** — provides the `huh/v2` implementation of the core `PromptRunner` interface; wire via `SetPromptRunner()`.
- **manpage** — `manpage.GenerateCommand` produces roff man pages via `muesli/mango-cobra`.

#### Audit Log

- `WithAuditLog(plugin)` wires `samber-do-auditlog` hooks into the injector via `buildInjectorOpts()`. `cli.AuditLog()` returns the plugin; `cli.AuditLogReport()` returns a snapshot. `AuditLogServiceByName`/`AuditLogFailedServices` query the report.
- `ExportAuditLog[T]` + `AuditLogExportConfig` write to file or `io.Writer` in **11 formats** (html, json, ndjson, csv, tsv, mermaid, dot, d2, plantuml, tree, htmltree). `ParseAuditLogFormat` validates input. No built-in `audit-log` subcommand — consumers implement their own export via flags/env.
- `samber-do-auditlog` is consumed from the Go module proxy (`v0.4.0`). The sibling repo at `../samber-do-auditlog` is for local dev only — a `replace` directive works for local builds but is **ignored by downstream consumers** (replace directives in a library's go.mod only affect the module's own build/CI).

#### Fang Integration (ADR-001)

- `WithCLIVersion` auto-pipes to `fang.WithVersion`; `WithCLICommit` auto-pipes to `fang.WithCommit`. Do NOT combine `WithFangOptions(fang.WithVersion(...))` with `WithCLIVersion` (duplicate opts).
- `fang.WithNotifySignal` is intentionally skipped — cmdguard's `WithSignalHandling`/`WithGracefulShutdown` provides DI-aware signal handling that fang cannot. See `docs/adr/001-fang-integration-strategy.md`.

#### Custom Types & Commands

- **Prompt tag** — `prompt:"Question?"` + `WithPromptOnMissing` enables interactive prompting for missing flags via the `PromptRunner` interface. The `huh/v2` implementation (bool → confirm, enum → select, else → input) lives in the `prompts` sub-module; without it, prompting returns an error.
- **envPrefix** — `WithEnvPrefix` sets the prefix on root AND command-level flags
- **GoDuration handler** — `RegisterGoDurationHandler()` validates the default at registration time (error for non-empty invalid defaults; empty defaults allowed as zero value)
- **DoctorCommand** — calls `HealthCheckResultsWithContext(ctx)` returning `map[string]error`; DI services with `do.HealthcheckerWithContext` are included automatically; custom checks via `WithDoctorCheck` run after DI checks

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
