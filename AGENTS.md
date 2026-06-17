# AGENTS.md - cmdguard Contributor & AI Assistant Guide

> **Note:** This file serves as both a contributor guide and context for AI-assisted development. It documents architecture decisions, API reference, coding standards, and known gotchas.

**Last Updated:** 2026-06-17
**Project:** cmdguard - CLI Guard Library
**Go Version:** 1.26
**Status:** v2.7.0 - zero panics, 85.6% coverage, 0 lint issues, 0 race conditions

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

**Important:** `git commit --no-verify` is required (pre-commit hooks have pre-existing errors).

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with type-safe dependency injection.

| API | Package           | Use Case                         |
| --- | ----------------- | -------------------------------- |
| v2  | `pkg/cmdguard/v2` | Type-safe, DI-powered, no panics |

**Current Status:** v2.7.0. 396+ tests passing, 85.6% coverage, 0 build errors.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   ├── v2/                       # v2 API (recommended)
│   │   ├── cli.go                # CLI[T] struct, NewCLI, AddCommand, Execute
│   │   ├── cli_accessors.go      # CLI accessor methods (Config, Scope, etc.)
│   │   ├── cli_command.go        # Internal cobra wiring (cliToCobraCommand)
│   │   ├── cli_options.go        # CLI functional options (WithCLIVersion, etc.)
│   │   ├── auditlog.go            # AuditLogCommand, WithAuditLog convenience helpers
│   │   ├── cli_output.go          # Output format flag registration and parsing, dynamic help from registry
│   │   ├── command.go            # Command[T,F] struct, constructors, options, Validate
│   │   ├── command_options.go    # All 19 CommandOption functions (WithShort, WithFlags, etc.)
│   │   ├── command_suggest.go    # (removed — consolidated into flags_suggest.go)
│   │   ├── config.go             # Config type constraint
│   │   ├── config_file.go        # ConfigFileLoader, JSON loader, WithConfigFile
│   │   ├── config_parsing.go     # ParseFlagTags, DefaultValue
│   │   ├── config_setfield.go    # SetField for config structs
│   │   ├── configload/           # Optional YAML/TOML loaders
│   │   ├── counting_flag.go      # (removed — logic in type_handler_kinds.go)
│   │   ├── editor.go             # EditInEditor ($EDITOR support)
│   │   ├── errors.go             # Error types (CommandError, FlagError, etc.) and remaining sentinels
│   │   ├── errors_command.go     # Command-related sentinel errors
│   │   ├── errors_config.go      # Config-related sentinel errors
│   │   ├── errors_di.go          # DI-related sentinel errors
│   │   ├── errors_flags.go       # Flag-related sentinel errors
│   │   ├── errors_audit.go       # Audit-log-related sentinel errors
│   │   ├── flags.go              # FlagRegistry with struct tags
│   │   ├── flags_parse.go        # Flag parsing logic
│   │   ├── flags_suggest.go      # Typo suggestions (Levenshtein)
│   │   ├── flags_validate.go     # Flag validation
│   │   ├── completion.go         # Shell completion support
│   │   ├── doc.go                # Package documentation
│   │   ├── flag_helpers.go       # Flag type constraints, cloning, parsing helpers
│   │   ├── flow_context.go       # BranchingFlowContext for command path tracking
│   │   ├── flow_context_access.go # Flow context helpers (typed value access)
│   │   ├── glamour.go            # Markdown help rendering via glamour/v2
│   │   ├── manpage.go            # Man page generation via mango
│   │   ├── middleware.go         # Middleware chain pattern
│   │   ├── output.go             # Rich output (OutputTable, OutputResult, shape-aware errors, dynamic format help)
│   │   ├── prompts.go            # Interactive prompts via huh/v2
│   │   ├── scope.go              # DI scope wrapping samber/do/v2
│   │   ├── spinner.go            # Terminal spinner middleware
│   │   ├── telemetry.go          # OpenTelemetry middleware
│   │   ├── type_handler.go       # Extensible type registry
│   │   ├── type_handler_kinds.go # Primitive kind handlers (string/bool/int/uint/float/slice)
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
├── AGENTS.md                     # This file
├── FEATURES.md                   # Feature status
├── TODO_LIST.md                  # Remaining tasks
├── .golangci.yml                 # Lint configuration
├── flake.nix                     # Nix dev shell, formatter, checks
├── flake.lock                    # Nix lock file
└── README.md                     # User documentation
```

### Package Guidelines

| Package           | Purpose       | Importable? | Coverage |
| ----------------- | ------------- | ----------- | -------- |
| `pkg/cmdguard/v2` | Type-safe API | Yes         | ~85%     |
| `pkg/testutil`    | Test helpers  | Yes         | —        |

---

## Key Dependencies

| Library                                     | Purpose              | Version |
| ------------------------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`                    | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`                   | Dependency injection | v2.0.0  |
| `github.com/spf13/pflag`                    | Flag parsing         | v1.0.10 |
| `charm.land/fang/v2`                        | Cobra styling        | v2.0.1  |
| `charm.land/huh/v2`                         | Interactive prompts  | v2.0.3  |
| `charm.land/glamour/v2`                     | Markdown rendering   | v2.0.0  |
| `go.opentelemetry.io/otel/trace`            | OpenTelemetry spans  | v1.44.0 |
| `github.com/larsartmann/go-output`          | Rich output formats  | v0.9.0  |
| `github.com/larsartmann/samber-do-auditlog` | DI audit logging     | v0.0.1  |

---

## API Reference

See [docs/API.md](docs/API.md) for the full API reference (constructors, options, methods, DI, middleware, error handling, version/doctor commands).

Quick reference: `NewCLI[T]`, `NewCommand[T,F]`, `NewParentCommand[T,F]`, `AddCommand`, `Execute`. All functions return errors — zero panics. See [pkg.go.dev](https://pkg.go.dev/github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2) for godoc.

---

## Coding Standards

### Go Conventions

- **Go 1.26** - Use modern Go features
- **gofumpt** formatting (via `golangci-lint fmt`)
- **Error handling** - Always check errors, wrap with `fmt.Errorf("context: %w", err)`
- **No panics** in v2 library code
- **Functional options** pattern for configuration
- **Constructor pattern** - All Command creation via `NewCommand`/`NewParentCommand`, struct fields unexported

### Testing

- `t.Parallel()` in every test function and subtest (paralleltest linter)
- `//nolint:paralleltest` for tests using `t.Setenv` or capturing `os.Stdout`
- `//nolint:fatcontext` at file level for test files with context in closures
- Table-driven tests: `tests := []struct{...}` pattern
- Two test packages: `v2` (internal, access private helpers) and `v2_test` (external)

### Test Commands

```bash
go test ./... -count=1 -timeout 120s -race     # All tests with race detection
go test ./... -count=1 -timeout 120s -cover     # Coverage report
golangci-lint run ./...                          # Lint (0 issues)
go build ./...                                   # Verify build
```

---

## Architecture Decisions

### v2 Design Principles

1. **Single type parameter** - `CLI[T]` only parameterizes on config; each command has its own flags type
2. **No Panics** - All operations return errors
3. **DI-Powered** - samber/do/v2 for dependency injection
4. **Typed Flags** - Struct tags for flag definitions
5. **Standalone AddCommand** - Function (not method) to support per-command flag types
6. **Env tags** - `env:"VAR_NAME"` struct tag reads from environment
7. **Counting flags** - `count:"true"` tag enables -v/-vv/-vvv pattern
8. **Signal handling** - `WithSignalHandling[T]()` for graceful shutdown
9. **Rich output** - OutputTable/OutputResult with 16 formats via go-output registries
10. **Copy-on-write registries** (v2.7.0+) — `FlagRegistry` shares global type/validator registries via copy-on-write; reads use the shared maps, writes trigger a lazy clone. `RegisterTypeHandler()`/`RegisterValidator()` write to global defaults (visible to instances that haven't cloned); `FlagRegistry.RegisterTypeHandler()`/`FlagRegistry.RegisterFlagValidator()` trigger COW clone and write to instance-local maps
11. **$EDITOR support** - `EditInEditor()` for user input editing
12. **Typo suggestions** - `SuggestFlag`/`SuggestCommand` with Levenshtein
13. **Full sentinel coverage** - All 40+ errors identifiable via `errors.Is()`
14. **Generic helpers** - `textMarshal[T]`/`textUnmarshal[T]`, `renderAndWrite`/`marshalAndWrite`, `branchWithCtx`
15. **Spinner middleware** `SpinnerMiddleware[T](title)` shows a lipgloss-styled spinner on stderr; skips when not a terminal
16. **Glamour help** — `WithGlamourHelp[T]()` renders command `Long` and `Example` fields via `charm.land/glamour/v2` markdown; uses `RenderWithEnvironmentConfig` (checks `GLAMOUR_STYLE` env var, defaults to `"dark"`); applied recursively to all commands at registration time
17. **Telemetry middleware** — `TelemetryMiddleware[T](tracer)` creates an OpenTelemetry span per command; requires `go.opentelemetry.io/otel/trace.Tracer`; `WithTelemetry[T](tracer)` is the convenience CLI option
18. **Audit log integration** — `WithAuditLog[T](plugin)` wires `samber-do-auditlog` into the DI injector; `AuditLogCommand[T](cli)` provides an `audit-log` subcommand with HTML/JSON/NDJSON/Mermaid export; `cli.AuditLog()` for programmatic access

### Key Gotchas

1. `t.Setenv` + `t.Parallel()` = panic — use `//nolint:paralleltest`
2. `PostRunE` is NOT called when `RunE` errors (Cobra behavior)
3. `NoFlags` is `type NoFlags = struct{}` — use `(NoFlags{})` with parens for comparisons
4. fang provides styled output by default; `--no-color` flag or `NO_COLOR` env var disables color (fang respects NO_COLOR automatically via colorprofile; `--no-color` sets `NO_COLOR=1` for fang to pick up)
5. `AddCommand` calls `cmd.Validate()` as defense-in-depth even though constructors already validate
6. **envPrefix propagation** — `WithEnvPrefix` sets prefix on root AND command-level flags (fixed in v2.2)
7. **Counting flags** — must use `int` type with `count:"true"` tag; don't reuse flag names from root config
8. **Prompt tag** — `prompt:"Question?"` on a struct field enables interactive prompting when the flag is missing and `WithPromptOnMissing` is set on the command. Bool fields use `huh.NewConfirm`, enum fields (with `values` tag) use `huh.NewSelect`, all others use `huh.NewInput`
9. **SuggestFlag API** — returns `(string, bool)` since v2.2 (breaking change from string-only)
10. **Instance-scoped registries** — `FlagRegistry` clones `typeRegistry` and `validatorRegistry` from globals at creation time; package-level `RegisterTypeHandler()`/`RegisterValidator()` write to the global defaults template, not to existing instances. Use `FlagRegistry.RegisterTypeHandler()` for per-instance customization.
11. **Direct go-output usage** — Users use `output.FormatTable`, `output.FormatJSON`, etc. directly from `github.com/larsartmann/go-output`; cmdguard only re-exports the `OutputFormat = output.Format` type alias for convenience. No `ParseOutputFormat`, `SupportedFormats`, `IsFormatSupported`, or `Format*` constant re-exports.
12. **Deprecated APIs (removed)** — `IsExecutable()` removed; use `HasHandler()`. `FlowContextAccessor` was removed in v2.3.0 — use `GetBranchingFlowContext(ctx)` directly
13. **Typed branching** — `BranchWithDuration(name, time.Duration)` and `BranchWithDeadlineTime(name, time.Time)` are the only branching methods (string-based `BranchWithTimeout`/`BranchWithDeadline` removed in v2.3.0)
14. **Regex validation cache** — `validateRegex` caches compiled patterns in `sync.Map`; global state, tests must not run in parallel
15. **Exit codes** — `ExecuteAndExit` checks for `ExitCoder` interface; use `NewExitError(code, err)` for custom exit codes
16. **`--no-color` + NO_COLOR** — `--no-color` persistent flag is registered by default; `cli.NoColor()` returns true if passed; `Execute()` sets `NO_COLOR=1` for fang to pick up. `NO_COLOR` env var is also respected automatically by fang's colorprofile.
17. **Strict validation** — `WithStrictValidation[T]()` requires `WithShort` on all commands; enforced at `AddCommand` time
18. **Draconian validation** — `WithDraconianValidation[T]()` is superset of strict + requires `WithExample` on leaf commands; parent commands are exempt
19. **Config validation** — `WithConfigValidation[T](fn)` runs after root flag parsing but before any command handler; blocks execution on error
20. **Args validation** — `WithExactArgs`/`WithMinimumArgs`/etc. use cobra's built-in arg validators; runs during command execution, not at registration
21. **Spinner non-TTY** — `SpinnerMiddleware` auto-skips when `os.Stderr` is not a terminal; use `SpinnerConfig{Writer: ...}` to override
22. **Glamour v2 env-based theme** — `WithGlamourHelp[T]()` now uses `RenderWithEnvironmentConfig` which checks `GLAMOUR_STYLE` env var, defaulting to "dark"; the string `"auto"` is no longer a valid glamour theme name in v2
23. **Telemetry context propagation** — `TelemetryMiddleware` starts a span but cannot propagate the new context to the handler due to the `next func() error` middleware API signature; child spans must use the original context passed to the handler
24. **FullPath populated at execution time** — `CommandInfo.FullPath` is set via `cobra.CommandPath()` inside the handler closure, NOT at command registration; it's empty in unit tests unless you call the handler through a cobra execution
25. **Glamour idempotent** — `applyGlamourIfEnabled` resets `glamourHelp=false` after applying to prevent double-rendering (which would wrap ANSI codes inside ANSI codes); calling Execute twice is safe
26. **outputEnabled removed** — The unused `outputEnabled` field was removed from `CLI[T]`; use `outputFormat != ""` to detect if output formatting is configured
27. **NewExitError returns (\***ExitError**, **error\**) — validates 0-255 range; breaking change from `*ExitError`
28. **NewScopeFromInjector returns (\***Scope**, **error\*\*) — nil injector returns error; breaking change from nil dereference
29. **Sentinel wrapping** — All 40+ errors use `fmt.Errorf("%w: ...", sentinel)` for `errors.Is()` chainability
30. **Config file precedence** — `WithConfigFile[T](paths...)` loads config BEFORE flag registration; config values become the new tag defaults, so flags/env still override them
31. **Config file paths** — supports `$ENV` expansion and `~` expansion; missing files are silently skipped
32. **Config file `--config` override** — if the config struct has a `config` flag, its value overrides the default search paths from `WithConfigFile`
33. **Config file flat only (v1)** — JSON/YAML/TOML loaders detect top-level keys matching `flag` tag names; nested structs in config files are not yet supported
34. **Nix flake limited** — `flake.nix` only provides devShell + formatter + format check (no `buildGoModule` or vet checks); could be extended now that go-output is published
35. **Glamour v2 no `"auto"` theme** — `charm.land/glamour/v2` removed the `"auto"` theme; use empty string (env-based via `GLAMOUR_STYLE`) or explicit theme like `"dark"`; `WithGlamourHelp` now sets theme to `""` for env-based detection
36. **DoctorCommand uses HealthCheckResults** — `DoctorCommand[T]` calls `HealthCheckResultsWithContext(ctx)` which returns `map[string]error`; DI services with `do.HealthcheckerWithContext` are included automatically; custom checks via `WithDoctorCheck` run after DI checks
37. **configload single file** — YAML/TOML/JSON/Auto loaders consolidated into `configload/loader.go`; uses `genericLoader` with pluggable `unmarshalFunc`; TOML import aliased as `toml` to avoid conflict with local `cmdguard` import alias
38. **Zero panics** — All Must* functions removed in v2.5.0. Every function returns errors. No `MustNewCommand`, `MustInvoke`, `MustParse*`, etc.
39. **Package signature** — `Package[T any](scope *Scope, name, short string, defaults T, opts ...CLIOption[T])` takes a pre-existing `*Scope` as the first parameter. Callers must create the scope with `NewScope(name)` before calling `Package`. This separates scope construction from CLI construction, eliminating signature duplication with `NewCLI`.
40. **WithGracefulShutdown** — `WithGracefulShutdown[T]()` enables graceful DI shutdown on SIGINT/SIGTERM. Implies `WithSignalHandling`. Services implementing `do.ShutdownerWithError` are shut down in reverse invocation order. `WithSignalHandling` only cancels context and does NOT trigger DI shutdown
41. **Override + CloneScope** — `Override[T](scope, provider)` and `OverrideValue[T](scope, value)` replace services in a scope for testing. `CloneScope(scope)` creates a copy with same registrations but no invoked state. Standard pattern: clone → override → invoke
42. **GoDuration default validation** — `RegisterGoDurationHandler()` now validates the default value at registration time (returns error for non-empty invalid defaults), consistent with bool/int/uint/float handlers; empty defaults are allowed (zero value)
43. **ErrLogLevel/ErrLogFormat error chain** — `ParseLogLevel`/`ParseLogFormat` now wrap errors with their respective sentinels (`ErrLogLevel`/`ErrLogFormat`), so `errors.Is(err, v2.ErrLogLevel)` works; the chain is `ErrLogLevel → EnumError → ErrInvalidEnum`
44. **Removed sentinels** — `ErrNoFlags`, `ErrTooFewArgs`, `ErrTooManyArgs` were declared but never used in any code path; removed to reduce API surface.
45. **Removed re-exports** — `ParseOutputFormat()`, `SupportedFormats()`, `IsFormatSupported()`, and 16 `Format*` constant re-exports removed; users import `output.Format*` directly from go-output. `IsExecutable()` removed; use `HasHandler()`.
46. **configload.Auto()** — tries YAML → TOML → JSON sequentially (not file-extension based); since JSON is valid YAML, JSON data is handled by the YAML parser first; use `LoaderForPath()` for precise format detection when the file extension is known
47. **ShutdownAll error chain** — `Shutdown()` wraps with `ErrServiceConstruction` once; `ShutdownAll` collects these without additional wrapping (fixed double-wrap in v2.4.0)
48. **Arg validators return errors** — `WithExactArgs`, `WithMinimumArgs`, `WithMaximumArgs`, `WithRangeArgs` now set an error on the command if given invalid args (negative n, min > max) instead of panicking. The error is surfaced by `NewCommand`/`NewParentCommand`.
49. **NewScopeWithOpts** — `NewScopeWithOpts(name, opts)` creates scope with `do.InjectorOpts` for custom logging, lifecycle hooks, health check timeouts. `WithDILogging[T](logf)` is the CLI convenience option.
50. **Audit log integration** — `WithAuditLog[T](plugin)` wires `samber-do-auditlog` hooks into the CLI's injector via `buildInjectorOpts()`, which merges DILogging and audit hooks. `AuditLogCommand[T](cli)` creates an `audit-log` subcommand supporting 4 formats (html, json, ndjson, mermaid). Returns `ErrAuditLogNotEnabled` when plugin is nil — callers should check with `errors.Is`. `cli.AuditLog()` returns the plugin; `cli.AuditLogReport()` returns a snapshot. `AuditLogCommand` gracefully degrades in tests without the plugin.
51. **buildInjectorOpts** — Merges `diLogf` and `auditLog` into a single `*do.InjectorOpts`. Returns nil when neither is configured (uses default injector). This replaced the old inline `NewScopeWithOpts` call in `cli.initialize()`.
52. **No replace directives** — go-output v0.9.0 and samber-do-auditlog v0.0.2 are fully published; no local replace directives needed.
53. **Fang integration (ADR-001)** — `WithCLIVersion` auto-pipes to `fang.WithVersion`; `WithCLICommit` auto-pipes to `fang.WithCommit`. Users should NOT use `WithFangOptions(fang.WithVersion(...))` alongside `WithCLIVersion` — this would create duplicate fang version opts. `fang.WithNotifySignal` is intentionally skipped because cmdguard's `WithSignalHandling`/`WithGracefulShutdown` provides DI-aware signal handling that fang cannot (see `docs/adr/001-fang-integration-strategy.md`).
54. **16 output formats via go-output v0.9.0 registries** — `RenderTableData` (all 16 formats) and `RenderAnyData` (JSON, YAML, TOML) via thread-safe `formatRegistry[T]`. Users import `output.FormatTable` etc. directly; cmdguard only re-exports `OutputFormat` type alias. `OutputResult()` provides shape-aware error messages. `OutputTable()` uses `AddRowChecked()` for fail-fast row validation. `--output` flag help is auto-generated from `RegisteredTableDataFormats()`. `RegisteredFormats()` exposes registered formats for callers.
55. **Deduplicated validators** — `validateEmail` and `validateURL` in `flags_validate.go` delegate to `ParseEmail()` and `ParseURL()` respectively, eliminating duplicate parsing logic.
56. **errors.AsType (Go 1.26)** — `output.go` uses `errors.AsType[*T]` instead of `errors.As(err, &v)` for consistency with `cli.go:298`. All new code should use `errors.AsType`.
57. **Validation error aggregation** — `ValidateConfig` uses `errors.Join(append([]error{ErrConfigValidation}, errs...)...)` so individual validation errors are reachable via `errors.Is`. Previous `%v` formatting lost the chain.
58. **errors_audit.go** — `ErrAuditLogNotEnabled` and `ErrInvalidOutputFormat` moved from `errors.go` to `errors_audit.go`, matching the per-domain split pattern (command/config/di/flags/audit).
59. **Copy-on-write registries (v2.7.0)** — `FlagRegistry` no longer eagerly clones `typeRegistry`/`validatorRegistry`. Instead it shares the global maps via `share()` and clones lazily on first write via `register()`. This reduces NewCLI by 48% (5.8µs vs 11µs) and saves 10 allocs per command. The `owned bool` and `parent *typeRegistry` fields control the COW state. Behavioral change: global registrations via `RegisterTypeHandler()`/`RegisterValidator()` are now visible to instances created before the registration (as long as they haven't triggered a lazy clone).
60. **Cached os.UserHomeDir()** — `config_file.go` caches the home directory via `sync.OnceValue` (`cachedHomeDir`). The home directory is immutable during a process lifetime, so this eliminates redundant syscalls when multiple config paths use `~/` expansion.
61. **Iterator methods (iter.Seq)** — `TagsSeq()`, `FlagNamesSeq()`, `PathSeq()`, `ChildrenSeq()` provide zero-allocation alternatives to `Tags()`, `FlagNames()`, `Path()`, `Children()`. The old methods still return defensive copies for backward compatibility. Use the `iter.Seq` variants when you only need to range over the data.

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
