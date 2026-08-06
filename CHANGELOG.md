# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Dates are in YYYY-MM-DD format (ISO 8601).

## [Unreleased]

### Added

- **flightrecorder sub-module** — independently importable module (`github.com/larsartmann/cmdguard/flightrecorder`) built on Go 1.25+ `runtime/trace.FlightRecorder`. Provides continuous in-memory trace buffering, context-aware snapshot capture, and automatic slow-command / error capture middleware. Zero external dependencies (stdlib only). Sixth workspace module. Full API surface: `Recorder` (Start/Stop/Enabled/WriteTo/Capture/CaptureToWriter), `Config` (8 fields with defaults), `Middleware[T]`, `WithFlightRecorder[T]`, `WithFlightRecorderRecorder[T]`, sentinel errors (`ErrAlreadyStarted`, `ErrNotEnabled`), 48 test functions (41 tests + 3 godoc examples), 3 benchmarks, 1 fuzz target, 96.1% coverage.
- CLI lifecycle improvements — enhanced command shutdown behavior and lifecycle test coverage.
- `TODO_LIST.md` and `ROADMAP.md` — created (were referenced in AGENTS.md but missing).

### Changed

- **Documentation v3→v4 drift fix** — 36 user-facing and contributor-facing documents updated from stale v3 API references to v4 (QUICKSTART, TUTORIAL, WHAT_THIS_PROJECT_IS_ABOUT, WHAT_THIS_PROJECT_IS_NOT, MIGRATION_FROM_COBRA, COMPARISON, PERFORMANCE, doc.go, README, website guides). 23 website source files updated.
- AGENTS.md project structure corrected from stale v3 references to v4; flightrecorder sub-module documented (package table, lint strategy, design principles, gotchas).
- `CHANGELOG.md` rebuilt with real version history (v0.1.0 through v4.0.0) derived from git tags and tag messages.
- `go-output` upgraded v0.35.0 → v0.37.0.
- Nix flake inputs updated; go.mod `replace` directives consolidated into a single block.

### Fixed

- `doc.go` godoc example — fixed stale `v3` import alias and "v2 constructors" reference.
- README.md sub-modules code example — fixed missing `"time"` import (used `time.Millisecond` without importing the `time` package).
- flightrecorder `evaluateCapture` — error reason now takes precedence over slow reason when both conditions are met (previously slow overwrote error).
- flightrecorder timestamp format — changed to nanosecond precision to prevent same-second concurrent captures from clobbering each other's files.

---

## [flightrecorder/v0.1.0] - 2026-08-06

First stable release of the `flightrecorder` sub-module. Tagged alongside v4.0.0 core.

- Zero external dependencies (Go 1.25+ stdlib `runtime/trace`)
- Continuous in-memory trace buffering with configurable capacity
- Context-aware snapshot capture: `Recorder.Capture`, `Recorder.CaptureToWriter`
- Automatic slow-command (`CaptureOnSlow`+`SlowThreshold`) and error (`CaptureOnError`) capture middleware
- `WithFlightRecorder[T]` and `WithFlightRecorderRecorder[T]` CLI options
- Process-wide singleton with `sync.WaitGroup` coordination
- 48 test functions, 3 benchmarks, 1 fuzz target, 96.1% coverage

---

## [4.0.0] - 2026-07-28

### Breaking Changes

- **Module path:** `github.com/larsartmann/cmdguard/v3` → `github.com/larsartmann/cmdguard/v4`
- **Package directory:** `pkg/cmdguard/v3/` → `pkg/cmdguard/v4/`
- **Package name:** `v3` → `v4` (update all import aliases)
- **Config loading:** `configload` sub-package deleted; `WithConfigFile(paths...)` now auto-detects JSON/YAML/TOML via `KoanfLoader`
- `NewJSONLoader` deleted; `configload` sub-package deleted

### Changed

- Consolidated config loading into unified `KoanfLoader` path resolution
- go-output upgraded v0.31.1 → v0.35.0
- samber-do-auditlog upgraded v0.7.0 → v0.8.1

---

## [3.1.0] - 2026-07-25

### Fixed

- Auditlog diagram export build break

---

## [3.0.0] - 2026-07-07

### Breaking Changes — Breaking API Redesign

Corrects the v2.11.0 mis-release that put breaking changes on a `/v2` path.

- **Non-generic `CLIOption` / `CommandOption`** — type inference via positional flags eliminates the v2 "7 type params per command" explosion. Metadata options (`WithShort`, `WithLong`, `WithExample`, etc.) take zero type parameters. Only lifecycle hooks (`WithPreRunE`, `WithPostRunE`, `WithSubcommands`) remain generic.
- **Module path:** `github.com/larsartmann/cmdguard/v2` → `github.com/larsartmann/cmdguard/v3`
- `NewCommand` / `NewParentCommand` API refactored — flags passed positionally, no `WithFlags` option needed
- Sealed lifecycle hook interfaces — compile-time type safety restored

### Added

- 5 extracted optional sub-modules (`glamour`, `manpage`, `prompts`, `spinner`, `telemetry`) — core has zero dependencies on these
- `go.work` multi-module workspace for unified local builds
- Audit logging with `samber-do-auditlog` integration

---

## [2.10.4] - 2026-07-07

### Fixed

- Retracts mis-released v2.11.0 (breaking changes were incorrectly published on a `/v2` path; corrected in v3.0.0)

---

## [2.10.0] - 2026-06-28

### Added — Cobra-Correctness Contract + Escape-Hatch APIs

Refocus on the founding mission: "make consumers use Cobra correctly."

- `SilenceUsage = true` by default (cobra's #1 footgun, off by default)
- `ExitCode(err) int` — public exit-code mapping function
- Scoped flags (`local:"true"`) — root-only flags not inherited by subcommands
- `hidden:"true"` flag tag — exclude from `--help`, stay functional
- `ConfigFromContext[T]` — type-safe config retrieval for raw cobra subcommands
- `WithPostFlagParse[T]` — post-parse hook (DI init, session storage)
- `WithCleanup[T]` — post-RunE cleanup that fires even when RunE errors

### Changed

- Go directive bumped 1.26.3 → 1.26.4 (CVEs GO-2026-5037/5038/5039)
- 457 tests, 1430 runs, 26 benchmarks, 7 fuzz targets, 86.7% coverage

---

## [2.9.0] - 2026-06-22

### Added

- 4 new audit-log export formats (d2, plantuml, tree, htmltree) — total now 11
- samber-do-auditlog v0.1.0 → v0.3.0
- go-output v0.17.1 → v0.17.2

### Changed

- 430 test functions, 26 benchmarks, 7 fuzz targets, 86.6% coverage, 0 lint issues

---

## [2.8.0] - 2026-06-19

### Changed

- Release v2.8.0 — incremental improvements and dependency updates

---

## [2.7.0] - 2026-06-17

### Changed

- Release v2.7.0 — incremental improvements

---

## [2.6.0] - 2026-06-12

### Added

- `RegisteredFormats()` — dynamic format discovery from registered marshalers
- Shape-aware error messages in `OutputResult()`
- `OutputTable()` uses `AddRowChecked()` for fail-fast row validation
- Dynamic `--output` flag help from `RegisteredTableDataFormats()`
- go-output v0.9.0 with generic registries

### Removed

- `IsExecutable()` — use `HasHandler()`
- 16 `Format*` constant re-exports — use `output.Format*` directly
- `ParseOutputFormat()` — use `output.ParseFormat()`
- `SupportedFormats()` — use `output.AllFormats`
- `IsFormatSupported()` — use `format.IsValid()`
- `ErrNoFlags`, `ErrTooFewArgs`, `ErrTooManyArgs` sentinels

### Changed

- 407+ tests passing, 85.9% coverage, 0 lint issues, 0 race conditions

---

## [2.5.0] - 2026-06-05

### Fixed

- Module path fixed for Go major version compatibility

---

## [1.0.0] - 2026-04-30

### Changed

- Stable release — dependency updates and go.sum alignment

---

## [0.2.0] - 2026-04-08

### Changed

- Minor update

---

## [0.1.0] - 2026-02-14

### Added — Initial Release

First stable release of cmdguard with the Guard API.

- Single-step initialization with `cmdguard.New()`
- Compile-time validation (panic on invalid commands)
- Built-in version and validate commands
- Environment-based configuration
- Strict mode for RunE enforcement
- 94% coverage on config package, 66% coverage on cmdguard package
- 3 direct dependencies (minimal external footprint)
