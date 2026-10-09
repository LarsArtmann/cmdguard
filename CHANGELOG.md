# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Dates are in YYYY-MM-DD format (ISO 8601).

## [Unreleased]

### Added

- **Command-level audit middleware.** `AuditMiddleware[T]` records a
  PhaseBefore event before each command runs and a PhaseAfter event after it
  completes, carrying the wall-clock duration in milliseconds and the
  command's error (nil on success), via samber-do-auditlog v0.11.0's
  `Plugin.RecordCommand`. Command events surface in `Report.Events` and the
  NDJSON stream (`Event.IsCommand()` distinguishes them from service events).
  `WithAuditMiddleware[T](plugin)` wires it into the middleware chain — pair
  it with `WithAuditLog(plugin)` using the same plugin. A nil plugin makes
  the middleware a transparent passthrough.
- **lint: baseline/ratchet mode.** `--write-baseline` snapshots current
  findings into `.cmdguard-lint-baseline.json`; subsequent runs fail only on
  findings beyond the baseline (per rule+file, line-drift tolerant, stale
  entries reported so the ratchet can tighten). Library surface:
  `LoadBaseline`/`WriteBaseline`/`ApplyBaseline`.
- **lint: cross-file name tracing for CG003/CG004/CG006.** Package-level CLI
  variables, Execute errors, and constructor errors are unioned per directory
  group (test files never merge with prod files; function-local names never
  leak), so a constructor in one file plus its misuse in another is caught.
  FP budget held: cmdguard itself and examples stay at 0 findings.
- **lint: provider options.** The BuildFlow provider declares `enable` and
  `disable` (comma-separated rule IDs, same semantics as the CLI flags via
  the new `lint.FilterByIDs`).
- **lint: generated README rule table.** `rules --markdown` emits the rule
  table from `allRuleDefs()`; `rules_table_test.go` fails the build when the
  README section drifts.

### Fixed

- Nothing yet.

### Changed

- **Dependency bumps:** `samber-do-auditlog` v0.10.0 → v0.11.0 (command
  execution events), `golang.org/x/net` v0.59.0 → v0.60.0 in core + glamour
  (GO-2026-6617/6611/6603, govulncheck-verified reachable).
- **website:** `postcss-nested` overridden to ^8.0.1 in
  `website/pnpm-workspace.yaml`, pulling patched `postcss-selector-parser`
  ≥7.1.6 (GHSA-rj75-hqrm-r3gf, moderate). `pnpm audit` clean; `astro build`
  green.

---

## [4.1.0] - 2026-10-04

### Fixed

- **Nil `TypeHandlerFunc.RegisterFunc` no longer silently drops the flag.** A handler
  that registers nothing made the flag vanish from the CLI, surfacing far from the
  cause as "unknown flag" at invocation time (found dogfooding in nsfw-classifier).
  `dispatchRegister` now rejects such handlers at CLI-build time with an error naming
  the type and flag. `Parse`/`Default` with nil funcs return errors/nil instead of
  panicking. The low-level `Register` stays nil-tolerant for derived handlers.
- **Repeated `[]string` flags no longer parse pflag's bracketed render form.** A
  changed slice flag is read through `flag.Value.String()`, which renders `"[a,b]"`;
  splitting that string on commas handed downstream handlers elements like `"[100"`
  and `"101]"` instead of `"100"` and `"101"` (seen live with a consumer's repeatable
  `--document` ID flag). The slice parse now trims the brackets before splitting,
  with an end-to-end regression test through `ExecuteWithArgs`.

### Changed

- **stdlib `time.Duration` is now handled by default.** A fresh registry (and the
  global template) includes the duration handler previously available only via the
  opt-in `RegisterGoDurationHandler()`: native pflag duration registration, `5s`-style
  parsing, duration defaults. Explicit `RegisterTypeHandler` overrides still win.

### CI

- Workflows install Go 1.27 (matching the repo's `go 1.27` floor) and pin
  golangci-lint to v2.14.0, the version whose `config verify` passes the repo config.

---

## [lint/v0.1.0] - 2026-10-04

First stable release of the `lint` sub-module: a usage linter for cmdguard-based CLIs.

- 6 syntactic rules (CG001–CG006): execute bypass (`fang.Execute` around cmdguard),
  stale major import, constructor panic (`MustNew*`), runtime `SetVersion`,
  duplicate version options, execute-error reprint
- Pure `go/parser` AST walk — no type info, no build required
- BuildFlow provider (`lint/provider`, toolsdk spec `cmdguard-lint`, module fan-out)
- Dogfood CLI `cmd/cmdguard-lint` (built with cmdguard v4 itself) with text, JSON,
  and SARIF 2.1.0 output (round-trips via `finding.FindingsFromSARIF`)
- Suppressions: `//cmdguard-lint:ignore <ID> <reason>` (reason required) on the
  finding line or the line above
- Per-directory memoized analysis with `ClearCache()` for long-lived processes
- `CurrentMajor` constant tracks the latest cmdguard major for CG002
- Validated: timesheets corpus 18/18 findings at hand-audited locations, 0 false
  positives on cmdguard/erraudit/go-structure-linter/branching-flow
- 34 test functions; golangci-lint clean under the repo config

---

## [4.0.2] - 2026-08-06

### Changed

- Dependency alignment release: sub-module and core dependency bumps to align
  with the v4.0.1 sub-module releases (no API changes).

---

## [glamour/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` to `v4.CLIOption`. Update import aliases from `v3` to `v4`.

---

## [prompts/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` to `v4.CLIOption`. Update import aliases from `v3` to `v4`.

---

## [spinner/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` / `v3.Middleware[T]` to `v4.CLIOption` / `v4.Middleware[T]`. Update import aliases from `v3` to `v4`.

### Fixed

- Added `sync.Once` and `sync.Mutex` to `textSpinner` for concurrency-safe stop behavior.

---

## [telemetry/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` / `v3.Middleware[T]` to `v4.CLIOption` / `v4.Middleware[T]`. Update import aliases from `v3` to `v4`.

---

## [4.0.2] - 2026-08-06

### Changed

- Dependency alignment release: sub-module and core dependency bumps to align
  with the v4.0.1 sub-module releases (no API changes).

---

## [4.0.1] - 2026-08-06

### Fixed

- `NO_COLOR` env var restore — `applyNoColorIfSet` now uses `os.LookupEnv` instead of `os.Getenv` to correctly distinguish "unset" from "set to empty string" when restoring the environment after `--no-color` flag execution. Previously, a pre-existing `NO_COLOR=""` would be unset instead of restored.
- `doc.go` godoc example — fixed stale `v3` import alias and "v2 constructors" reference.
- README.md sub-modules code example — fixed missing `"time"` import (used `time.Millisecond` without importing the `time` package).

### Changed

- `go-output` upgraded v0.35.0 → v0.37.0.
- `koanf` patched: v2.3.5 → v2.3.6, parsers/json v1.0.0 → v1.0.1, parsers/yaml v1.1.0 → v1.1.1.
- **Documentation v3→v4 drift fix** — 36 user-facing and contributor-facing documents updated from stale v3 API references to v4 (QUICKSTART, TUTORIAL, WHAT_THIS_PROJECT_IS_ABOUT, WHAT_THIS_PROJECT_IS_NOT, MIGRATION_FROM_COBRA, COMPARISON, PERFORMANCE, doc.go, README, website guides). 23 website source files updated.
- AGENTS.md project structure corrected from stale v3 references to v4; flightrecorder sub-module documented (package table, lint strategy, design principles, gotchas).
- Nix flake inputs updated; go.mod `replace` directives consolidated into a single block.
- `TODO_LIST.md` and `ROADMAP.md` created.

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
