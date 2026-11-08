# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Dates are in YYYY-MM-DD format (ISO 8601).

## [Unreleased]

### Added

- **Stable machine error codes (`ErrorCode(err)`, envelope `error_code`).**
  Every sentinel now carries a stable lowercase snake_case code
  (`required_flag`, `config_unknown_key`, …) that survives Go identifier
  renames, so scripts and agents can switch on failures without parsing
  messages. The machine-readable error envelope (`--output json|yaml|toml`)
  gains an additive `error_code` field; the human error path is unchanged.
  Sentinel facts (identifier, envelope class, code, hint) live in ONE
  canonical table, so hint coverage, error classification, and codes cannot
  drift apart — and run-time errors that previously classified as `"unknown"`
  in the envelope (invalid URL/email/duration values, command timeouts,
  editor failures) now carry their real class.
- **Generated, drift-guarded error reference.** `docs/ERROR_REFERENCE.md`'s
  sentinel table is rendered from the same canonical fact table via
  `GenerateErrorReference()`; a test fails when the committed doc stops
  matching the code. Regenerate with
  `go test ./pkg/cmdguard/v4 -run TestErrorReferenceDoc -update-golden`.
- **Strict config keys, opt in (`WithStrictConfigKeys()`).** Unknown keys in a
  config file are now an error instead of a silent no-op, with the offending
  dotted key, the loaded file path, and a "did you mean" suggestion for likely
  typos (`ErrConfigUnknownKey`). Accepted spellings match what config loading
  already accepts (canonical field name, flag name, and case/hyphen/underscore
  variants), and nested/array keys are named in full (`servers[0].port`).
  Off by default in v4 because a shared config file may legitimately carry
  extra keys; **`WithBatteries()` turns it on** as part of the production
  defaults (v5 flips the global default). `doctor` reports config problems
  (missing file, parse failure, unknown keys) alongside DI health via the new
  `CLI.ConfigDiagnostics()`, and `WithDoctorSchemaValidator` attaches an
  external JSON Schema check.
- **Config provenance: the new `config explain` command (`WithConfigExplainCommand()`).**
  Answers "where did this value come from?" per field: flag, environment,
  config file, tag default, or unset — resolved with the same precedence the
  runtime uses. Table output for humans, `--output json` for scripts.
  `CLI.ExplainConfig()` exposes the same data programmatically. Fields tagged
  `secret:"true"` are redacted to `***` in every view while still reporting
  their source, so provenance can be audited without leaking credentials.
- **`secret:"true"` struct tag.** Marks a config value that must never be
  logged or committed. It appears in the machine contract (`describe`), the
  generated docs, the config schema, and drives `config explain` redaction.
- **`schema` sub-module** (`github.com/larsartmann/cmdguard/schema`): emits a
  JSON Schema draft 2020-12 document describing a config struct (`FromConfig`)
  or a whole CLI (`FromCLI`), including nested objects, arrays, enums,
  min/max constraints, and cmdguard hints (`flag`, `env`, format, `secret`)
  under `x-cmdguard`. Properties always reject unknown keys
  (`additionalProperties: false`), matching strict loading — a schema typo is
  the same typo the binary rejects. Output is deterministic and editor-friendly
  (`Emit`/`Marshal` write indented JSON). Experimental surface.
- **`v4.ConfigKeys(cfg)`**: the public vocabulary of a config struct — every
  accepted key with its flag, env var, type, format, constraints, default, and
  nested shape. It is the single source shared by strict-key checking and the
  schema emitter, so the two can never disagree.

### Fixed

- **`min`/`max` validate rules now resolve in the flag type's domain.**
  Duration-typed flags (`v4.Duration` and `time.Duration`) compare bounds as
  durations, so `validate:"min=10s,max=1m"` means ten seconds and one minute —
  previously the numeric validators parsed the value text as a plain number
  and failed with "invalid validator parameter" (or, for `time.Duration`,
  silently compared nanoseconds). Textual parsed types (URL, Email, HostPort,
  FilePath) fall back to length bounds instead of failing to parse as
  numbers. String kinds keep the existing length semantics, numeric kinds and
  Port keep numeric semantics. The same resolution now applies to
  `ValidateConfig`, which previously validated the raw tag text without the
  string length rewrite.
- **Impossible flag defaults now fail at construction, not on the first
  run.** `NewFlagRegistry` (and therefore `NewCommand`/`NewCLI`) validates
  every tag default against that tag's `values` and `validate` rules and
  rejects violations with the new `ErrInvalidFlagDefault` sentinel. Defaults
  whose `validate` tag names a validator registered after construction are
  skipped at construction time and remain enforced at run time, so the
  register-after-construction pattern keeps working.
- **Config values in camelCase, lowercase, or snake_case no longer lose the
  race against flag defaults.** Key matching accepted `{"LogLevel": "debug"}`
  (case-insensitive JSON decoding) while the "which fields did the file set?"
  bookkeeping matched only the exact flag spelling (`log-level`), so those
  fields were treated as unset and flag parsing overwrote the config value
  with the tag default. Every accepted spelling is now normalized to the
  canonical key before decoding, and the same accepted-spelling vocabulary is
  used for both the value assignment and the bookkeeping.
- **Generated markdown docs now carry exactly the facts the machine contract
  carries.** `GenerateDocs` previously rendered from pflag's `FlagUsages`,
  a second representation that silently drifted from `describe` (missing env
  vars, required markers, accepted values, counting and secret markers).
  It now renders from the same source as the JSON contract: an eight-column
  table (flag, type, env, default, values, required, secret, help), with inherited
  flags in their own section and table-escaping so pipes in help text cannot
  break the layout. Root-only CLIs (no subcommands) no longer produce empty
  output.

- **`ExitCode` no longer maps a failing error to exit 0 when an `ExitCoder`
  reports 0 (#11).** Third-party error types (e.g. go-error-family) implement
  `ExitCoder` with 0 as their "no explicit override" sentinel; `extractExitCode`
  trusted it verbatim, so failed CLIs exited 0. A zero code from an `ExitCoder`
  is now treated as "no override" and falls back to the generic failure code 1
  (a nil error still yields 0). The JSON error envelope `code` field follows the
  same rule.
- **Flag values no longer leak across repeated `ExecuteWithArgs` calls on a
  reused CLI (#13).** pflag marks a flag `Changed` forever after its first parse
  and never restores values, so a second execute re-ran with the first parse's
  flags (including write-action flags like `--apply` that the new args never
  mentioned). Every `Execute` now resets per-flag state first: `Changed`
  markers, values (stringSlice values are rebuilt because pflag's append
  semantics make `Set(DefValue)` corrupt them), and the typed config, which is
  restored to its post-constructor snapshot (defaults plus config-file
  overlay). One-execute-per-process CLIs are unaffected. Caveat: raw cobra
  escape-hatch flags of pflag's other appending kinds (intSlice and friends)
  keep their residue value; only their `Changed` marker resets.

### Added

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
- **lint: CG007 missing short description, CG008 missing CLI version.** The
  end-goal baseline: a command without `WithShort` renders as a blank help
  line, and a CLI without `WithCLIVersion` cannot identify its build. Both are
  warnings, skip `_test.go` files and calls with spread option lists, and are
  suppressible with the usual `//cmdguard-lint:ignore`. Ships to consumers
  with `lint/v0.2.0`.

---

## [4.2.1] - 2026-11-04

### Changed

- **Dependency updates across the workspace** (no API changes): koanf v2.4.1
  (parsers json v1.0.1, yaml v1.1.1), go-output v0.38.7,
  samber-do-auditlog v0.12.1, lipgloss v2.0.6, and matching sub-module bumps
  (lint: go-finding v1.15.1, toolsdk v1.16.0, go-linter-sdk v0.4.0).

---

## [4.2.0] - 2026-11-02

**Behavior changes to review before upgrading:**

- **Flag tags are enforced now.** `required`, `values`, and `validate` run at
  execution time against resolved values (flag, env, default). CLIs that
  silently accepted missing required flags or invalid enum/range values will
  start failing with typed errors. Discovery commands (`describe`, `version`,
  `doctor`, `completion`, `help`) are exempt from `required` enforcement, and
  `min`/`max` on string flags now mean length.
- **Validation error text changed.** Failures read
  `Command "add": flag "title": value too short: "ab" must be at least 3 characters.`
  — no more stacked `preparing run context / validating flags for / ...`
  prefixes. Anything matching on the old wording (tests, parsing) must update.
- **Errors teach.** Human-readable errors print one actionable `Hint:` line
  (disable with `WithoutErrorHints()`); machine-format error envelopes carry
  the same sentence as `hint`.

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
- **Machine-readable CLI contract.** `cli.Describe()` and
  `cli.WriteDescription(w)` emit a versioned JSON schema (`schema_version`)
  of the whole command tree, including flag type, default, env var (with prefix
  applied), enum values, required/count/hidden/inherited, and the command tree
  with aliases, groups and deprecations. A hidden `describe` command is
  registered automatically at `Execute` time so agents and tooling can drive a
  CLI without scraping `--help`; a consumer-defined `describe` command wins and
  `WithoutDescribeCommand()` opts out. Additive within a schema major.
- **Atomic batch command registration.** `AddCommands[T](cli, cmds...)`
  validates and converts every command before attaching any, accepts
  heterogeneous flags types through the sealed `AddableCommand[T]` interface,
  rejects nil entries, and keys duplicate detection on the cobra command name
  so `"list [path]"` and `"list"` collide. `AddCommand` shares the same
  prepare/attach path and no longer reserves a name when conversion fails.
- **Actionable error hints.** `ErrorHint(err)` maps every sentinel used by the
  JSON error classifier to one fix sentence, and the machine-readable error
  envelope now carries it as `hint`. Human-readable errors also get a `Hint:`
  line (fang's default view and the plain cobra path); `WithoutErrorHints()`
  turns it off. A coverage test fails when a newly classified sentinel lacks a
  hint.
- **Production batteries.** `WithBatteries()` opts into strict command
  validation, signal handling with graceful DI shutdown, error hints, and
  auto-registered `version` (when `WithCLIVersion` is set) and `doctor`
  commands; consumer-defined commands of the same name win. `NewApp` wraps
  construction plus command registration and discards the half-built CLI when
  the build callback fails.

### Fixed

- **`required`, `values`, and `validate` tags are now enforced at run time.**
  `FlagRegistry.ValidateFlags` had tests but no production caller, so the tags
  were parsed and documented but never checked: a missing required flag fell
  through to the handler, and an invalid enum or `min`/`max` value was accepted.
  Validation now runs after parsing for the root and every command registry.
  `min`/`max` on string fields mean length (previously the text was parsed as a
  number, producing "invalid validator parameter"). Environment- and
  default-sourced values are validated too; `describe`, `version`, `doctor`,
  `completion`, and `help` skip required enforcement so discovery always works.
- **Readable validation errors.** Removed the stacked
  "preparing run context / validating flags for / validating flag / required
  flag not set" prefixes; a failure now reads
  `Command "add": flag "title": value too short: "ab" must be at least 3 characters.`

### Changed

- **Dependency bumps:** `samber-do-auditlog` v0.10.0 → v0.11.0 (command
  execution events), `golang.org/x/net` v0.59.0 → v0.60.0 in core + glamour
  (GO-2026-6617/6611/6603, govulncheck-verified reachable).
- **website:** `postcss-nested` overridden to ^8.0.1 in
  `website/pnpm-workspace.yaml`, pulling patched `postcss-selector-parser`
  ≥7.1.6 (GHSA-rj75-hqrm-r3gf, moderate). `pnpm audit` clean; `astro build`
  green.
- **website:** dependency refresh clearing the 2026-11-02 Dependabot batch —
  astro 7.4.4, vite 8.3.5, typescript 7.0.2, html-validate 11.18.0, plus
  transitive `undici` ≥8.12.2, `zod` ≥4.6.6, `fast-uri` ≥3.1.9. Full
  `pnpm audit` clean (19 findings, 6 high → 0); site builds 22 pages.

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

## [2.10.3] - 2026-07-07

### Changed

- go-output v0.30.1, samber-do-auditlog v0.4.0 (`d006757`). Entry added 2026-11-07 — the tag existed without a CHANGELOG entry (found by the docs-health sweep).

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
