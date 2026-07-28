# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

## [4.0.0] - 2026-07-28

### Breaking Changes

- **Module path changed to `github.com/larsartmann/cmdguard/v4`** — update all import paths from `/v3` to `/v4`. The package directory moved from `pkg/cmdguard/v3/` to `pkg/cmdguard/v4/`. Package name changed from `v3` to `v4`. Update import aliases from `v3` to `v4` (e.g., `v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"`).
- **Config loading consolidated to a single KoanfLoader** — the `configload` sub-package (`configload.YAML()`, `configload.TOML()`, `configload.JSON()`, `configload.Auto()`, `configload.LoaderForPath()`, `configload.NewKoanfLoader()`) and the core `jsonLoader`/`NewJSONLoader()` have been deleted. `WithConfigFile(paths...)` now creates a `KoanfLoader` that auto-detects JSON/YAML/TOML by file extension. KoanfLoader lives in the `v4` package (`koanf_loader.go`), uses koanf only as a format parser, converts to JSON, then reuses the shared `loadConfigFromJSON` processing path for case-insensitive nested struct matching. Dependency changes: `go-faster/yaml` and `pelletier/go-toml/v2` demoted from direct to `// indirect`; `koanf/providers/file` eliminated in favor of an in-repo `bytesProvider`; `koanf/parsers/toml` (`v0.1.0`) added.

### Migration from v3

1. **Module path**: `github.com/larsartmann/cmdguard/v3` → `github.com/larsartmann/cmdguard/v4`
2. **Import paths**: `cmdguard/v3/pkg/cmdguard/v3` → `cmdguard/v4/pkg/cmdguard/v4`
3. **Package alias**: `v3 "..."` → `v4 "..."` and all `v3.Function()` → `v4.Function()`
4. **Config loading**: Replace `configload.YAML()`/`configload.TOML()`/`configload.Auto()` + `WithConfigFileLoader(loader, paths...)` with just `WithConfigFile(paths...)`. Replace `configload.NewKoanfLoader(paths...)` with `v4.NewKoanfLoader(paths...)`.

### Changed

- **go-output upgraded** from v0.31.1 to v0.35.0 (thread-safe format registries, 16 output formats).
- **samber-do-auditlog upgraded** from v0.7.0 to v0.8.1.
- Sub-module dependency updates across all 4 sub-modules (glamour, prompts, spinner, telemetry).
- Deferred `TODO(v4)` markers updated to `TODO(v5)` — the naming-review renames (TypeHandler → TypeCodec, CommandInfo → CommandMetadata, PromptRunner → HuhPrompter) are deferred to v5.

## [3.1.0] - 2026-07-25

### Added

- **Audit logging integration** — `WithAuditLog(plugin)` wires `samber-do-auditlog` hooks into the DI injector via `buildInjectorOpts()`; `cli.AuditLog()`/`cli.AuditLogReport()` for programmatic access; `AuditLogServiceByName`/`AuditLogFailedServices` query helpers; `ExportAuditLog[T]` supports 11 formats (html, json, ndjson, csv, tsv, mermaid, dot, d2, plantuml, tree, htmltree). No built-in subcommand — consumers export via their own flag/env pattern.
- **Prompts sub-module** (`github.com/larsartmann/cmdguard/prompts`) — `huh/v2` interactive prompt runner implementing the core `PromptRunner` interface (bool → confirm, enum → select, else → input).
- **Sub-module releases** tagged at `v0.1.0` for `glamour/v0.1.0`, `prompts/v0.1.0`, `spinner/v0.1.0`, and `telemetry/v0.1.0`.
- **CLI error handling test coverage** — comprehensive tests for error paths, exit codes, and `ExitCoder`/`NewExitError` contracts.

### Changed

- **CLI options restructured** — command configuration and options handling refactored for cleaner internal wiring.
- **Audit log integration refreshed** — updated to `samber-do-auditlog v0.7.0` API (`auditlog.ServiceName` removed; `ServiceByName` now takes a plain string).
- **Migration guide and living docs refreshed** — `docs/MIGRATION_v2_v3.md`, `AGENTS.md`, `FEATURES.md`, `TODO_LIST.md`, and `ROADMAP.md` updated to reflect the current 4-sub-module workspace and verified metrics.
- **Internal traversal consolidated** — repeated internal parsing/traversal logic deduplicated.

### Fixed

- **Sub-module external resolution** — the 4 sub-modules were moved from `pkg/cmdguard/<name>/` to the repo root (`<name>/`) so their module paths (`github.com/larsartmann/cmdguard/<name>`) resolve for external consumers. Previously the `replace` directives in the root `go.mod` only worked locally; consumers got "missing go.mod at revision". Each sub-module's `go.mod` now requires the real published `cmdguard/v3 v3.0.0` (not the placeholder).
- **telemetry sub-module compile error** — `WithTelemetry` returned `v3.CLIOption[T]` (non-existent generic); now returns non-generic `v3.CLIOption`.
- **prompts sub-module lint violation** — `make([]huh.Option[string], len(options))` replaced with a zero-length append slice to satisfy the `makezero` linter.

### Removed

- **Manpage sub-module** — removed from the workspace (`34a0c6e`). The feature was not worth the maintenance surface for v3; consumers who need man pages can use `muesli/mango-cobra` directly.

## [3.0.0] - 2026-07-07

> **⚠️ Version correction:** The v3 redesign was initially mis-tagged as `v2.11.0` on the `/v2` module path — a semver violation since the API is breaking. That tag has been deleted and retracted. The module path is now `github.com/larsartmann/cmdguard/v3`. If you pulled `v2.11.0`, switch to `v3.0.0` (update import paths from `/v2` to `/v3`). The v2 line continues at `v2.10.4` (retracts the bad `v2.11.0`).

### Breaking Changes — v3 API Redesign

#### Command API: Non-generic options + type inference

- `CommandOption` is now non-generic. All metadata options (`WithShort`, `WithLong`, `WithExample`, `WithGroupID`, `WithNoArgs`, etc.) require zero type parameters.
- `NewCommand` signature changed: `NewCommand(use string, flags F, runE func(ctx, *T, F) error, opts ...CommandOption)` — flags are now the second positional argument, enabling full type inference.
- `WithFlags` option deleted entirely — pass flags positionally to `NewCommand`.
- `NewParentCommand` signature changed: `NewParentCommand[T](use, long string, flags F, opts ...CommandOption)` — subcommands via `WithSubcommands(cmds...)` option.
- `WithPreRunE`, `WithPostRunE`, `WithSubcommands` are generic functions that return non-generic `CommandOption` — type safety preserved via sealed interface pattern.

**Before (7 type params per command):**

```go
v2.NewCommand[AppConfig, *ListFlags]("list", handler,
    v2.WithShort[AppConfig, *ListFlags]("List tasks"),
    v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
)
```

**After (zero type params):**

```go
v3.NewCommand("list", &ListFlags{}, handler,
    v3.WithShort("List tasks"),
)
```

#### Mono-repo modularization: optional sub-modules

Heavy dependencies extracted into optional importable modules. Core direct deps reduced from 30 to 13.

| Module    | Import path                                 | Deps isolated              |
| --------- | ------------------------------------------- | -------------------------- |
| Telemetry | `github.com/larsartmann/cmdguard/telemetry` | OpenTelemetry SDK          |
| Manpage   | `github.com/larsartmann/cmdguard/manpage`   | mango/roff                 |
| Glamour   | `github.com/larsartmann/cmdguard/glamour`   | chroma/goldmark/bluemonday |
| Prompts   | `github.com/larsartmann/cmdguard/prompts`   | huh/bubbles/bubbletea      |
| Spinner   | `github.com/larsartmann/cmdguard/spinner`   | lipgloss                   |

Core extension hooks: `WithHelpTransform[T]()`, `PromptRunner` interface + `SetPromptRunner()`.

#### Removed from core

- `result.go` — sum types (Result[T], Validated[T]) not a CLI concern
- `editor.go` — EditInEditor $EDITOR support (marginal feature)
- `telemetry.go` — moved to telemetry sub-module
- `glamour.go` — moved to glamour sub-module
- `spinner.go` — moved to spinner sub-module
- `manpage.go` — moved to manpage sub-module
- 10 go-output blank imports removed from `output.go`

## [2.10.4] - 2026-07-07

### Fixed

- **Retract `v2.11.0`** — added `retract` directive to `go.mod`. The v3 redesign was incorrectly tagged `v2.11.0` on the `/v2` module path (a semver violation: the API is breaking). That tag is deleted and retracted. **Consumers should migrate to `v3.0.0`** (update import paths from `/v2` to `/v3`). See the [v3.0.0] entry and the [migration guide](docs/MIGRATION_v2_v3.md).

### Changed

- Created `release/v2.10` maintenance branch (home for any future v2.x patches).
- This is the final planned release on the `/v2` module path.

## [2.10.3] - 2026-07-06

### Changed

- Upgrade `github.com/larsartmann/go-output` to v0.30.1
- Upgrade `github.com/larsartmann/samber-do-auditlog` to v0.4.0

### Note

- This tag was originally orphaned (on no branch). It now lives on `release/v2.10`.
- The same dependency upgrades are included in `v3.0.0`.

## [2.10.2] - 2026-07-05

### Changed

- Bump `go-output` from v0.23.0 to v0.23.3 (all 10 direct sub-modules + 2 indirect: `escape`, `daghtml`)
- Bump `samber-do-auditlog` from v0.3.0 to v0.3.1
- Bump `go-toml/v2` from v2.4.2 to v2.4.3
- Bump `lipgloss/v2` from v2.0.4 to v2.0.5
- Bump `bubbles/v2` from v2.1.0 to v2.1.1
- Bump `bubbletea/v2` from v2.0.7 to v2.0.8

### Docs

- Sync AGENTS.md and FEATURES.md dependency tables to actual go.mod versions (go-output v0.17.2→v0.23.3, samber-do-auditlog v0.3.0→v0.3.1)
- Fix stale `Version: 2.8.1` in package doc comment
- Update go-output sub-modules note: `enum`/`envdetect` absorbed into core; indirect set is now `escape` + `daghtml`

## [2.10.1] - 2026-07-02

### Changed

- Bump `go-output` from v0.17.2 to v0.23.0 (all 10 direct sub-modules updated)
- Bump `go-toml/v2` from v2.4.0 to v2.4.2

### Removed

- Drop `go-output/enum` and `go-output/envdetect` indirect dependencies (absorbed into go-output core)

### Docs

- Normalize table column padding across all markdown files (FEATURES.md, README.md, TODO_LIST.md, docs/)
- Update .gitignore and flake.lock

## [2.10.0] - 2026-06-28

This release refocuses cmdguard on its founding mission — "make consumers use
Cobra correctly." The trigger was auditing BuildFlow, the primary consumer,
which had built four workarounds around cmdguard's gaps. Each new API below
replaces a specific workaround. The flagship example also no longer exits 0 on
failure or double-prints errors.

### Added

- **`SilenceUsage = true` by default** — Cobra's #1 footgun (dumping full usage
  after every command error) is now off by default. Fang already forced this
  true; now fang-off mode matches. `--help` is unaffected. Closes the core
  "use Cobra correctly" contract
- **`ExitCode(err) int`** — public exit-code mapping (nil → 0, `ExitCoder` → its
  code, else → 1). `ExecuteAndExit` now uses it; consumers can call it directly
- **Scoped flags (`local:"true"` tag)** — root-only flags that are NOT inherited
  by subcommands. Keeps subcommand `--help` focused. Replaces the manual
  "register this flag group again on each subcommand" pattern
- **`hidden:"true"` flag tag** — excludes a flag from `--help` while keeping it
  fully functional. Replaces hardcoded `flag.Hidden = true` lists keyed by name
- **`ConfigFromContext[T](ctx) (*T, bool)`** — type-safe retrieval of the
  resolved config for raw `*cobra.Command` subcommands (the "escape hatch" added
  via `cli.RootCommand().AddCommand`). Eliminates the hand-rolled context-key
  session systems consumers previously needed
- **`WithPostFlagParse[T](fns ...)`** — a hook that runs after flag parsing and
  config validation but before any command handler. Use for DI initialisation,
  session storage, logging setup. Replaces the manual "save + wrap cmdguard's
  `PersistentPreRunE`" workaround
- **`WithCleanup[T](fns ...)`** — a hook that runs after a command's `RunE`
  completes, including when `RunE` errors. Closes the Cobra gap where neither
  `PostRunE` nor `PersistentPostRunE` fire on `RunE` error. The hook receives
  the command, the resolved config, and the `RunE` error (nil on success); the
  original error is never swallowed (cleanup errors are joined, so both stay
  reachable via `errors.Is`). Covers both cmdguard-managed commands and raw
  cobra subcommands (escape hatch) by wrapping `RunE` at execute time

### Changed

- **`examples/taskctl`** — the flagship example now exits non-zero on command
  failure (`v2.ExitCode(execErr)` instead of `os.Exit(0)`) and no longer
  double-prints errors. It now models the correct contract it teaches
- **`WithPostFlagParse` execution order** — parse flags → store config in
  context → `configValidate` → `postFlagParse` hooks → command handler

### Fixed

- **Flagship example exit code** — `ExecuteAndExit` now respects `ExitCoder`;
  the example previously exited 0 even on handler errors

### Security

- **Go directive bumped `1.26.3` → `1.26.4`** — fixes stdlib CVEs
  GO-2026-5037, GO-2026-5038, GO-2026-5039 (one reachable via
  `ExitError.Error` → `crypto/x509`). Unblocked now that nixpkgs ships
  `go_1_26 >= 1.26.4`

### Removed

- **`makezero` linter disabled** in `.golangci.yml` — it directly conflicts with
  staticcheck `S1019` (makezero wants `make([]T, 0, N)` + append; staticcheck
  flags the equivalent `make([]T, N, N)` as redundant). Keeping both is impossible

## [2.9.0] - 2026-06-22

### Changed

- **`samber-do-auditlog` v0.1.0 → v0.3.0** — v0.2.0 added Plugin-level Mermaid/DOT/D2/PlantUML wrappers (eliminating cmdguard's local adapter functions and `writeReportToFile` helper); v0.3.0 added Tree/HTMLTree export formats. All upgrades are additive and non-breaking
- **`go-output` v0.17.1 → v0.17.2** — All 10 direct sub-modules bumped in lockstep; the 3 indirect modules (`enum`, `escape`, `envdetect`) remain at v0.17.1 (latest tag available for each). Non-breaking
- **Transitive dependency refresh** — All indirect dependencies updated to latest via `go get -u all` + `go mod tidy`; `go.sum` resynced

### Added

- **4 new audit log export formats** — `d2`, `plantuml`, `tree`, `htmltree` join the existing 7 (html, json, ndjson, csv, tsv, mermaid, dot), bringing the total to 11. All use Plugin-level methods directly (no adapter functions needed since samber-do-auditlog v0.2.0)

### Removed

- **Audit log wrapper functions** — `exportMermaidReportToFile`, `writeMermaidReport`, `exportDOTReportToFile`, `writeDOTReport`, and `writeReportToFile` removed; superseded by Plugin-level `ExportToMermaid`/`WriteMermaid`/`ExportToDOT`/`WriteDOT` methods added in samber-do-auditlog v0.2.0

### Security

- **govulncheck identifies 3 stdlib vulnerabilities in Go 1.26.3** (GO-2026-5037, GO-2026-5038, GO-2026-5039), all fixed in Go 1.26.4. One (GO-2026-5037) is reachable via `ExitError.Error` → `crypto/x509`. Bumping `go.mod` to `go 1.26.4` is deferred until nixpkgs packages `go_1_26 >= 1.26.4`; currently the nix sandbox cannot auto-download the toolchain during `nix flake check`. Consumers building with Go 1.26.4+ are not affected

## [2.8.1] - 2026-06-21

### Changed

- **`go-output` v0.13.0 → v0.17.0** — All 13 modules (root + 9 direct sub-modules + 3 indirect) updated to v0.17.0 in lockstep. The v0.14–v0.17 renderer sub-modules (d2, delimited, graph, markdown, markup, plantuml, serialization, table, tree, enum, escape, envdetect) were missing from the Go proxy due to a go-output release tagging gap; tagged and pushed them to fix the gap. No cmdguard code changes required — all consumed APIs (`RenderTableData`, `RenderAnyData`, `TableData`, `ParseFormat`, `RegisteredTableDataFormats`, `Format`) are stable across the range

## [2.8.0] - 2026-06-19

### Added

- **`Result[T]` and `Validated[T]` sum types** — Explicit error-handling types inspired by Rust's `Result<T, E>`. `Ok(value)` / `Err(err)` constructors, `IsOk()`/`IsErr()` predicates, `Value()`/`Expect()`/`UnwrapOr()` accessors. Zero panics. `Validated[T]` wraps a value plus a slice of validation errors (partial success pattern)
- **`Plugin` system** — `Plugin` interface for bundling custom type handlers and validators. `PluginRegistrar` exposes scoped `TypeHandler()`/`Validator()` registration. `RegisterPlugin()` applies to global registries; `FlagRegistry.RegisterPlugin()` applies per-instance
- **`GenerateDocs()`** — `cli.GenerateDocs(w)` writes markdown documentation for the full command tree (synopsis, usage, flags, examples) to any `io.Writer`
- **`ExportAuditLog[T]` helper** — Reusable function that writes the audit log snapshot in HTML/JSON/NDJSON/Mermaid/CSV/TSV/DOT format to a file or `io.Writer`. Consumers no longer need to implement the format switch themselves
- **`AuditLogFormat` strong type** — Validated enum (`html`, `json`, `ndjson`, `mermaid`, `csv`, `tsv`, `dot`) with `ParseAuditLogFormat()` constructor and `Valid()` method. Replaces raw string format selection
- **`ErrUnsupportedAuditLogFormat`** — Sentinel error for invalid format values, classified as `"audit"` in the JSON error system
- **Integer overflow validation** — `int8`/`int16`/`int32`/`uint8`/`uint16` flag values are range-checked at parse time; returns `ErrIntegerOverflow` instead of silently wrapping
- **Nested struct config support** — Config structs can contain nested structs; inner fields are discovered and flattened for flag registration and config-file loading. `Index` field on `FieldTag` tracks the reflect path
- **Koanf-based config loader** — `configload.KoanfLoader()` handles nested config objects (e.g. `{"db":{"host":"x"}}` → `--db-host`) via `github.com/knadh/koanf`

### Changed

- **`samber-do-auditlog` v0.0.4 → v0.1.0** — Consumed from the Go module proxy; local `replace` directive removed. v0.1.0 is API-compatible (all surfaces cmdguard uses are in the stable set). Adds CSV/TSV/DOT export wired through the new map-dispatch registry
- **Audit export dispatch** — `exportAuditLogToFile`/`exportAuditLogToWriter` refactored from per-format switch statements to a single map-of-exporters per direction, dropping cyclomatic complexity below the lint threshold and making new formats a one-line addition
- **`go-output` v0.12.0 → v0.13.0** — All 15 sub-modules updated. `markdown/` and `tree/` extracted into standalone sub-modules; `output.go` imports them explicitly to preserve `FormatMarkdown` and `FormatTree` availability

### Removed

- **`AuditLogCommand[T]`** — Built-in `audit-log` subcommand removed. Consumers implement their own export via flags/env + `ExportAuditLog[T]` (see BuildFlow for reference)
- **`AuditLogOption`, `WithAuditLogShort`, `WithAuditLogLong`, `WithAuditLogGroupID`** — Command-specific options removed with the subcommand
- **`ErrAuditLogNotEnabled`, `ErrInvalidOutputFormat`** — Sentinel errors removed with the subcommand
- **`errors_audit.go`** — File removed; audit format errors now in `auditlog.go`

### Fixed

- **Type-registry data race** — `RegisterTypeHandler()`/`RegisterValidator()` now use a locked accessor instead of direct map mutation, eliminating the race detected under `-race` during concurrent test runs

---

## [2.7.0] - 2026-06-17

### Added

- **Copy-on-write registries** — `FlagRegistry` shares global type/validator registries lazily; clones only on first write. Reduces `NewCLI` by 48% (~5.8µs vs 11µs) and saves 10 allocations per command
- **Cached `os.UserHomeDir()`** — `sync.OnceValue` eliminates redundant syscalls during config path `~/` expansion
- **Iterator methods (`iter.Seq`)** — `TagsSeq()`, `FlagNamesSeq()`, `PathSeq()`, `ChildrenSeq()` provide zero-allocation traversal alternatives to slice-returning variants
- **COW isolation tests and benchmarks** — 6 correctness tests and dedicated benchmarks for copy-on-write behavior
- **Performance analysis report** — Comprehensive HTML report at `docs/research/performance-analysis.html`

### Changed

- **Dependency updates** — `samber-do-auditlog` v0.0.2 → v0.0.4, `go-toml/v2` v2.2.0 → v2.4.0, `chroma/v2` v2.14.0 → v2.27.0

---

## [2.6.1] - 2026-06-14

### Added

- **Restored `MustNewCommand` and `MustNewParentCommand`** — Convenience wrappers for quick prototyping (callers who need zero-panics should use `NewCommand`/`NewParentCommand`)

### Changed

- **Dependency updates** — Updated transitive dependency pins in `go.mod` and `go.sum`

---

## [2.6.0] - 2026-06-12

### Added

- **go-output v0.9.0 integration** — Delegated `RenderTableData` and `RenderAnyData` to go-output registries; removed 167 lines of duplicated formatting logic
- **4 new output formats** — JSONL, AsciiDoc, TOML, PlantUML (16 total formats)
- **fang integration** — `WithCLIVersion`/`WithCLICommit` auto-pipe to fang; `WithFangErrorHandler` and `WithFangColorScheme` exposed
- **koanf-based config loader** — `configload` sub-package with JSON, YAML, and TOML support via `LoaderForPath()` and `Auto()`
- **Graceful shutdown** — `WithGracefulShutdown[T]()` triggers DI service shutdown on SIGINT/SIGTERM via `do.ShutdownerWithError`
- **DI test helpers** — `Override[T]`, `OverrideValue[T]`, `CloneScope()` for test isolation
- **DI logging** — `WithDILogging[T](logf)` and `NewScopeWithOpts(name, opts)` for custom injector configuration
- **Doctor command** — `DoctorCommand[T](cli)` with `HealthCheckResults` and custom `WithDoctorCheck` diagnostics
- **MustParse value type helpers** — `MustParseDuration`, `MustParseLogLevel`, `MustParseLogFormat`, `MustParseEnum` for API consistency
- **ADR-001** — Documented fang integration strategy in `docs/adr/001-fang-integration-strategy.md`
- **Comprehensive error reference** — `docs/ERROR_REFERENCE.md` with 62 sentinels

### Changed

- **Zero panics guarantee** — Removed all 16+ `Must*` panic-inducing functions from library code; every function returns errors
- **`NoFlags` distinct type** — Changed from `struct{}` alias to proper named type
- **Removed `WithColor`** — Deprecated option removed; use `WithFang` instead
- **`IsExecutable()` removed** — Use `HasHandler()` instead
- **Error file split** — `errors.go` → `errors_command.go`, `errors_flags.go`, `errors_config.go`, `errors_di.go`, `errors_audit.go`
- **API reference extracted** — Moved from `AGENTS.md` (649 lines) to dedicated `docs/API.md`
- **Go 1.26 modernization** — `errors.As` → `errors.AsType` throughout codebase
- **Test infrastructure deduplicated** — Zero semantic clone groups remaining at threshold 30

### Fixed

- **`ShutdownAll` double-wrapping** — `ErrServiceConstruction` no longer wrapped twice
- **`NO_COLOR` restoration** — Environment variable restored after execution instead of permanently mutated
- **`flow_context.SetValue` child safety** — Skips children with locally-set keys
- **Config `Auto()` format detection** — Tries YAML → TOML → JSON instead of only JSON
- **`ErrLogLevel`/`ErrLogFormat` error chain** — Parse errors now properly chain to their sentinels

---

## [2.5.0] - 2026-06-10

### Added

- **Runtime type guards in `dispatchParse`** — Verifies return type matches handler target
- **Short tag validation** — `short` tag must be exactly 1 character; returns error at registration
- **Nil tracer guard in `TelemetryMiddleware`** — Returns error instead of nil dereference
- **`BranchingFlowContext.SetValue` child safety** — Skips children that have local key set
- **Mutable slice protection** — `Tags()` and `Path()` return cloned slices
- **Arg validator error returns** — `WithExactArgs`, `WithMinimumArgs`, etc. return errors instead of panicking for invalid args

### Changed

- **Removed all Must\* panic-inducing functions** — Zero panics in library code. Every function returns errors:
  - `MustNewCommand`, `MustNewParentCommand` (use `NewCommand`, `NewParentCommand`)
  - `MustNewCLI`, `MustAddCommand` (use `NewCLI`, `AddCommand`)
  - `MustVersionCommand` (use `VersionCommand`)
  - `MustDoctorCommand` (use `DoctorCommand`)
  - `MustInvoke[T]`, `MustInvokeNamed[T]` (use `Invoke[T]`)
  - `MustGet[T]`, `RequireBranchingFlowContext` (use `GetBranchingFlowContext`)
  - `MustParse[T]`, `MustParseDuration`, `MustParseLogLevel`, `MustParseLogFormat`, `MustParseEnum`
  - `MustParseURL`, `MustParseEmail`, `MustParsePort`, `MustParseFilePath`, `MustParseHostPort`
- **Split `errors.go` into domain-specific files** — `errors_command.go`, `errors_flags.go`, `errors_config.go`, `errors_di.go`
- **Unified `fieldValueToString`/`formatFieldValue`** — Single canonical `formatFieldValue()` in `flag_helpers.go`
- **`doc.go` "never panics" claim now truthful** — No Must\* functions remain
- **Consolidated `command_suggest.go` into `flags_suggest.go`** — Single file for all typo suggestions

### Fixed

- **`setStringField` panic on type mismatch** — Added `AssignableTo` guard before `field.Set`
- **`ExitError.Error()` nil panic** — Added nil guard for `e.Err`
- **Double-wrapping of `ErrServiceConstruction` in `ShutdownAll`** — Fixed in v2.4.0

### Removed

- All 16+ Must\* panic variants (see "Changed" section above)

---

## [2.4.0] - 2026-06-03

### Added

- **`--no-color` flag + `NO_COLOR` support** - `cli.NoColor()` returns true when `--no-color` passed or `NO_COLOR=1` env var is set

### Changed

- **Removed `FlowContextAccessor`** - Dead API with zero consumers; use `GetBranchingFlowContext(ctx)` directly
- **Removed string-based `BranchWithTimeout`/`BranchWithDeadline`** - Use typed `BranchWithDuration(name, time.Duration)` and `BranchWithDeadlineTime(name, time.Time)` instead
- **`TimingMiddleware` callback signature** - Now includes `error` parameter to distinguish success vs failure timing
- **Migrated to `charm.land` vanity imports** - `charmbracelet/glamour` → `charm.land/glamour/v2`
- **Replaced `gopkg.in/yaml.v3`** with `github.com/go-faster/yaml` for YAML config loading
- **Modernized to Go 1.26** - `errors.As` → `errors.AsType` in examples

### Fixed

- Resolved `tparallel` lint issue in `TestCLINoColor` by extracting `t.Setenv` subtest

---

## [2.3.0] - 2026-06-01

### Added

- **Interactive prompts** - `WithPromptOnMissing[T,F]()` with `prompt:"Question?"` struct tag via charmbracelet/huh
- **Glamour markdown help** - `WithGlamourHelp[T]()` and `WithGlamourHelpTheme[T](theme)` render command Long/Example as markdown
- **Spinner middleware** - `SpinnerMiddleware[T](title)` and `SpinnerMiddlewareWithConfig[T](config)` with TTY auto-detection
- **Telemetry middleware** - `TelemetryMiddleware[T](tracer)` creates OpenTelemetry span per command; `WithTelemetry[T](tracer)` convenience option
- **Config file loading** - `WithConfigFile[T](paths...)` for JSON, `WithConfigFileLoader[T](loader, paths...)` for YAML/TOML via `configload` sub-package
- **`FullPath` in CommandInfo** - Full command path (`root.subcmd.leaf`) available in middleware
- **`Phase` typed enum** - Replaces `CommandInfo.Phase string` with type-safe enum
- **Taskctl example** - Single comprehensive showcase replacing 13 scattered examples
- **Nix flake** - `flake.nix` with devShell (Go 1.26, gopls, golangci-lint), treefmt formatter, format check
- **Fuzz tests** - 7 fuzz targets for all value type parsers

### Changed

- **Architecture hardening** - Extracted `handlerConfig[T,F]` from 8-param wire function, consolidated 5 error types into `labeledError`
- **File splitting** - Split `type_handler.go` (481 lines), `command.go` (403 lines), `flow_context.go` (396 lines) into focused files
- **Deduplication** - Eliminated all semantic clone groups (0 remaining)
- **Error wrapping** - All 40+ errors use `fmt.Errorf("%w: ...", sentinel)` for `errors.Is()` chainability
- **Contextual errors** - Added context wrapping to CLI, command, and flag initialization errors
- **Value types** - Normalized `IsEmpty()` across all types, extracted shared `requireNonEmpty` helper
- **Removed `Ptr[T]`** - Use Go 1.26 built-in `new(v)` instead

### Fixed

- **CommandInfo concurrency** - Move copy inside handler closure for goroutine safety
- **Config file race** - Eliminate race in `resolveConfigFlag`
- **Enum validation** - Validate Enum values in SetField instead of silently bypassing allowed list
- **Required tag parsing** - Propagate error instead of silently ignoring
- **MergeConfigs** - True deep copy for reference types
- **Flag cloning** - Deep copy in `cloneFlags` to prevent shared mutable state
- **Validator registry** - Clone from global defaults per FlagRegistry instance

---

## [2.2.0] - 2026-04-30

### Added

- **`env:"VAR"` struct tag** - Environment variable binding with `WithEnvPrefix[T]("MYAPP_")` prefix support
- **`count:"true"` struct tag** - Counting flags: `-vvv` → 3
- **12 output formats** - go-output integration: table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, and more
- **`EditInEditor(ctx, content)`** - Open `$EDITOR` for user input (now with `context.Context`)
- **`WithFang[T](bool)`** - Proper name for styled help (deprecates `WithColor`)
- **Subcommand typo suggestions** - `SuggestCommand` with Levenshtein distance
- **`WithSignalHandling[T]()`** - SIGINT/SIGTERM context cancellation
- **Instance-scoped validators** - `FlagRegistry.RegisterFlagValidator()` for per-instance customization
- **Custom value types** - URL, Email, Port, FilePath, HostPort with full validation
- **Shell completion** - `WithCompletion[T,F](fn)` and `WithValidArgs[T,F](args...)`
- **Man page generation** - `cli.ManPage()`, `cli.WriteManPage()`, `GenerateManPageCommand[T](cli)`
- **`WithOutputFormat[T]()`** - Auto `--output` flag with format selection
- **Command groups** - `WithGroup[T](id, title)` for organized help output
- **Flag validation** - `validate:"email,min=5"` tag with built-in + custom validators

### Changed

- **TypeHandler registry** - Unified 3-way split brain into single registry
- **Validator registry** - Instance-scoped, removing global mutable state
- **Constructor pattern** - All Command creation via `NewCommand`/`NewParentCommand`, struct fields unexported
- **0 lint issues** - From 113 to 0, added comprehensive `.golangci.yml`
- **Refactored output.go** - Monolithic switch → format renderer registry

### Fixed

- **55 race conditions** - `sync.RWMutex` on `globalTypeRegistry`
- **`Ptr[T]` function** - Returned zero-valued pointer instead of pointer to v
- **`SuggestFlag` API** - Returns `(string, bool)` for consistency
- **envPrefix propagation** - Now reaches command-level FlagRegistry
- **Local go-output replace** - Removed after tagging v0.1.0

### Removed

- **Dead code** - `parseCustomDefault`, `wrapErr`, `parseField`, `parseAndSetLog*`
- **Scattered examples** - Consolidated into single `examples/taskctl/` showcase

---

## [2.1.0] - 2026-03-28

### Added

- **CLI[T] New API** - Type-safe CLI builder with generic config and flags
- **BranchingFlowContext** - Tracks command execution path with context propagation
- **Option[T] type** - Functional options pattern for CLI configuration
- **SimpleCLI alias** - Backward compatibility alias for CLI[T]
- **Functional options** - CLIOption[T] for configuring CLI instances
- **WithCLIScope option** - Inject existing DI scope
- **WithLong option** - Set long description via option
- **Comprehensive flow_context tests** - 344 lines, 44 tests
- **CLI[T] integration tests**

### Changed

- **Deprecated GuardedCommand** - Use CLI[T] instead (see migration guide)
- **Updated examples/typed to use SimpleCLI pattern**
- **Improved godoc for public APIs**
- **Updated README.md with links to new docs**

### Fixed

- **flow_context cancel bug** - Self-cancel tracking now works correctly
- **wrapcheck errors** - json.Marshal/Unmarshal now properly annotated
- **go.mod compatibility** - Fixed for Go 1.26
- **go:fix inline directive** - Removed duplicate in type_helpers.go

### Removed

- **AddCommandFunc** - Redundant, use AddCommand instead
- **Dead code packages** - Removed pkg/apperrors, pkg/testutil

### Documentation

- **Migration Guide v1→v2** - docs/MIGRATION_v1_v2.md
- **Quickstart Guide** - docs/QUICKSTART.md
- **README updates** - Links to new documentation

---

## [2.0.0] - 2026-03-22

### Added

- **v2 API** - Type-safe API with dependency injection
- **samber/do/v2 integration** - Full DI support with scopes
- **Typed errors** - Sentinel errors for precise error handling
- **Struct-based flags** - FlagRegistry with struct tag support
- **Comprehensive tests** - 90%+ coverage across packages
- **BDD-style tests** - Behavior documentation via Ginkgo/Gomega
- **Benchmark tests** - Performance validation

### Features

- Generic Command[T, F] type
- GuardedCommand[T, F] for type-safe CLI building
- DI-powered service management
- Flag parsing and validation
- Typo suggestions for unknown flags
- Subcommand support with different flag types
- PreRunE/PostRunE hooks
- Scoped providers for plugin architecture

### Custom Types

- Enum with validation
- Duration with parsing
- LogLevel (debug, info, warn, error)
- LogFormat (text, json)

---

## [1.0.0] - 2026-04-30

Stability commitment release.

### Changed

- Dependency updates (charmbracelet/x/exp transitive dependencies)

---

## [0.2.0] - 2026-04-08

### Added

- `WithSilenceErrors`, `WithSilenceUsage`, `WithColor` CLI options
- Comprehensive tests for flag helper functions, cliToCobraCommand edge cases, error paths

### Changed

- README and AGENTS.md rewritten for v2.1 API patterns

---

## [0.1.0] - 2026-02-20

### Added

- Initial release of cmdguard
- Type-safe CLI construction with generics
- Dependency injection via samber/do/v2
- Flag binding with struct tags
- Full Cobra integration

[Unreleased]: https://github.com/larsartmann/cmdguard/compare/v4.0.0...HEAD
[4.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v4.0.0
[3.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v3.1.0
[3.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v3.0.0
[2.10.4]: https://github.com/larsartmann/cmdguard/releases/tag/v2.10.4
[2.10.3]: https://github.com/larsartmann/cmdguard/releases/tag/v2.10.3
[2.10.2]: https://github.com/larsartmann/cmdguard/releases/tag/v2.10.2
[2.10.1]: https://github.com/larsartmann/cmdguard/releases/tag/v2.10.1
[2.10.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.10.0
[2.9.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.9.0
[2.8.1]: https://github.com/larsartmann/cmdguard/releases/tag/v2.8.1
[2.8.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.8.0
[2.7.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.7.0
[2.6.1]: https://github.com/larsartmann/cmdguard/releases/tag/v2.6.1
[2.6.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.6.0
[2.5.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.5.0
[2.4.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.4.0
[2.3.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.3.0
[2.2.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.2.0
[2.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.1.0
[2.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.0.0
[1.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v1.0.0
[0.2.0]: https://github.com/larsartmann/cmdguard/releases/tag/v0.2.0
[0.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v0.1.0
