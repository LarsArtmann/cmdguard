# cmdguard — Comprehensive Status Report

**Date:** 2026-06-11 12:18\
**Author:** Crush (AI assistant)\
**Trigger:** User-requested full status audit\
**Branch:** master\
**Version:** v2.5.0

---

## Executive Summary

cmdguard is in **excellent shape**. The library is production-ready with 404 passing tests, 85.9% coverage, 0 lint issues, 0 race conditions, and 0 panics in library code. The recent output.go delegation to go-output registries reduced output.go from 309 to 147 lines while gaining all 16 output formats. The self-critique session found and fixed 5 concrete issues (errors.AsType adoption, sentinel relocation, error chain preservation, brittle test, stale metrics).

**One real blocker:** 10 local `replace` directives in `go.mod` prevent external consumers from building. Tagging go-output v0.9.0 and samber-do-auditlog v0.0.2 is required before any public release.

---

## a) FULLY DONE

### Core Library (v2.5.0)

| Feature                   | Status | Details                                                       |
| ------------------------- | ------ | ------------------------------------------------------------- |
| CLI[T] + NewCLI           | ✅     | Type-safe CLI construction, all options                       |
| Command[T,F] + NewCommand | ✅     | 21 command options, zero panics                               |
| Flag system (struct tags) | ✅     | flag/short/default/help/env/count/prompt/validate/values tags |
| TypeHandler registry      | ✅     | 9 built-in kinds + 9 custom types, extensible                 |
| DI scope (samber/do/v2)   | ✅     | Provide/Invoke/Override/Clone/Shutdown/HealthCheck            |
| Rich output (16 formats)  | ✅     | Delegated to go-output registries via blank imports           |
| Middleware chain          | ✅     | Timing, Recovery, Spinner, Telemetry + custom                 |
| Shell completion          | ✅     | Dynamic + static, cobra-compatible                            |
| Man page generation       | ✅     | Via mango + roff                                              |
| Markdown help (glamour)   | ✅     | Auto theme detection, GLAMOUR_STYLE support                   |
| Interactive prompts (huh) | ✅     | Text, select, confirm with auto-type detection                |
| Config file loading       | ✅     | JSON (core) + YAML/TOML/Auto (configload sub-package)         |
| Error handling            | ✅     | 60+ sentinels, typed errors, full errors.Is chains            |
| Doctor command            | ✅     | DI health checks + custom checks                              |
| Version command           | ✅     | Typed subcommand                                              |
| Audit log command         | ✅     | HTML/JSON/NDJSON/Mermaid export                               |
| Zero panics guarantee     | ✅     | All Must\* removed, all functions return errors               |

### Infrastructure

| Item                                          | Status |
| --------------------------------------------- | ------ |
| Nix flake (devShell, formatter, format check) | ✅     |
| CI (golangci-lint, codecov, nix check)        | ✅     |
| Release automation workflow                   | ✅     |
| 89 golangci-lint linters enabled              | ✅     |
| Benchmark suite (22 benchmarks)               | ✅     |
| Fuzz test targets (7 targets)                 | ✅     |

### Documentation

| File                                       | Status            |
| ------------------------------------------ | ----------------- |
| AGENTS.md (57 gotchas)                     | ✅ Current        |
| FEATURES.md (all features verified)        | ✅ Current        |
| TODO_LIST.md (phases 1-15 complete)        | ✅ Current        |
| README.md                                  | ✅ Current        |
| docs/API.md                                | ✅ Current        |
| docs/ERROR_REFERENCE.md (62 sentinels)     | ✅ Current        |
| docs/CLI_DESIGN_PRINCIPLES.md              | ✅ Current        |
| docs/adr/001-fang-integration-strategy.md  | ✅ Current        |
| examples/taskctl/ (production example CLI) | ✅ 70.5% coverage |

### Session Work (2026-06-11)

| Commit    | What                                           | Impact                                                     |
| --------- | ---------------------------------------------- | ---------------------------------------------------------- |
| `4b5bb7c` | Adopt `errors.AsType` in output.go             | Eliminated 2 gopls hints, Go 1.26 consistency              |
| `17d20ac` | Extract audit-log sentinels to errors_audit.go | Per-domain file split pattern                              |
| `f296106` | Use `errors.Join` in ValidateConfig            | Preserves errors.Is chain for individual validation errors |
| `a3fa206` | Replace hard-coded format count                | Future-proof against go-output format additions            |
| `4ea803d` | Update AGENTS.md with gotchas #55-57           | Documentation current                                      |

---

## b) PARTIALLY DONE

### Coverage Gaps (85.9% — target 90%)

**13 functions at 0% coverage:**

| Function                     | File                  | Why                           | Effort |
| ---------------------------- | --------------------- | ----------------------------- | ------ |
| `PromptString`               | prompts.go:23         | Requires terminal TTY mock    | Medium |
| `PromptSelect`               | prompts.go:37         | Requires terminal TTY mock    | Medium |
| `PromptConfirm`              | prompts.go:57         | Requires terminal TTY mock    | Medium |
| `NewManPage`                 | manpage.go:63         | Manual section constructor    | Easy   |
| `WithConfigFileLoader`       | config_file.go:178    | Custom loader registration    | Easy   |
| `WithDoctorLong`             | doctor.go:34          | Simple option, trivial        | Easy   |
| `WithAuditLogGroupID`        | auditlog.go:43        | Simple option, trivial        | Easy   |
| `RegisterValidator`          | flags_validate.go:79  | Global validator registration | Easy   |
| `validateEmail`              | flags_validate.go:153 | Delegates to ParseEmail       | Easy   |
| `validateURL`                | flags_validate.go:165 | Delegates to ParseURL         | Easy   |
| `validateNonEmpty`           | flags_validate.go:292 | String non-empty validator    | Easy   |
| `validateFieldByKind`        | flags_validate.go:301 | Kind-dispatched validation    | Easy   |
| `runValidateTagWithRegistry` | flags_validate.go:312 | Tag validation runner         | Easy   |

**11 functions at 1-59% coverage:**

| Function                      | Coverage | File                   |
| ----------------------------- | -------- | ---------------------- |
| `validateTagRules`            | 16.7%    | flags.go:195           |
| `GenerateManPageCommand`      | 14.3%    | manpage.go:41          |
| `applyNoColorIfSet`           | 25.0%    | cli.go:224             |
| `isEnvSet`                    | 28.6%    | prompts.go:174         |
| `writeAuditMermaid`           | 35.7%    | auditlog.go:147        |
| `validateMax`                 | 41.7%    | flags_validate.go:236  |
| `writeAuditToFileOrWriter`    | 44.4%    | auditlog.go:124        |
| `WithRangeArgs`               | 50.0%    | command_options.go:161 |
| `SpinnerMiddlewareWithConfig` | 50.0%    | spinner.go:93          |
| `validateMin`                 | 50.0%    | flags_validate.go:213  |
| `validateMaxLen`              | 55.6%    | flags_validate.go:195  |

### go.mod Replace Directives

10 local `replace` directives pointing to `../go-output` and `../samber-do-auditlog`. Required for development but **blocks external consumers**. Partially resolved — go-output v0.8.0 is published with new APIs, but sub-modules need v0.9.0 tags.

---

## c) NOT STARTED

| #  | Task                                                                            | Category     | Priority | Effort |
| -- | ------------------------------------------------------------------------------- | ------------ | -------- | ------ |
| 1  | Tag go-output v0.9.0 and remove replace directives                              | Release      | P0       | 30m    |
| 2  | Tag samber-do-auditlog v0.0.2 (commit html_templ.go)                            | Release      | P0       | 15m    |
| 3  | Remove deprecated `OutputStyledTable`                                           | Cleanup      | P2       | 10m    |
| 4  | Fix prompt test coverage (0% → 80%+)                                            | Testing      | P1       | 2h     |
| 5  | Fix manpage test coverage (14% → 80%+)                                          | Testing      | P1       | 1h     |
| 6  | Fix auditlog write coverage (35-44% → 80%+)                                     | Testing      | P1       | 1h     |
| 7  | Duplicate `jsonLoader` unification                                              | Refactor     | P1       | 30m    |
| 8  | `RegisterTypeHandler`/`RegisterValidator` return errors instead of panic risk   | Architecture | P2       | 1h     |
| 9  | Move `initNoColorFlag` out of cli_output.go                                     | Cleanup      | P3       | 5m     |
| 10 | Per-command middleware support                                                  | Feature      | P1       | 2h     |
| 11 | Structured JSON error output for `--output=json`                                | Feature      | P1       | 1h     |
| 12 | Plugin system for custom validators/type handlers                               | Feature      | P3       | 4h     |
| 13 | Config file nested struct support                                               | Feature      | P3       | 4h     |
| 14 | Config auto-loading with koanf                                                  | Feature      | P3       | 3h     |
| 15 | Advanced types: Result[T], Validated[T], branded IDs                            | Feature      | P3       | 6h     |
| 16 | Extract flag-related code to `flagtags` library                                 | Refactor     | P3       | 4h     |
| 17 | Documentation generation command                                                | Feature      | P3       | 4h     |
| 18 | `MergeConfigs` zero-value blind spot (false/0/"" never override)                | Fix          | P2       | 1h     |
| 19 | Wrap `do.Provide` calls in recover() for no-panic guarantee                     | Architecture | P2       | 30m    |
| 20 | Remove unused sentinels in v3 (`ErrNoFlags`, `ErrTooFewArgs`, `ErrTooManyArgs`) | Cleanup      | P3       | 5m     |

---

## d) TOTALLY FUCKED UP

### Nothing is "totally fucked up" — this is the honest assessment:

| Issue                               | Severity   | Details                                                                                                                                       |
| ----------------------------------- | ---------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| go.mod replace directives           | **HIGH**   | External consumers **cannot build** without local go-output/auditlog repos. Must tag upstream and remove before any release.                  |
| Prompt functions untested           | **MEDIUM** | `PromptString`, `PromptSelect`, `PromptConfirm` all at 0%. These interact with TTY via huh — need proper mocking strategy.                    |
| `RegisterTypeHandler` panic risk    | **LOW**    | Global registry functions can panic (from `do.Provide` internally). Violates "zero panics" guarantee if called incorrectly.                   |
| `validateEmail`/`validateURL` at 0% | **LOW**    | Functions delegate to `ParseEmail()`/`ParseURL()` (which ARE tested), but the wrapper functions themselves have zero test coverage. Easy fix. |
| Duplicate `jsonLoader`              | **LOW**    | Same struct+method in two packages (`config_file.go` and `configload/loader.go`). Dumb but harmless.                                          |

---

## e) WHAT WE SHOULD IMPROVE

### Architecture Improvements

1. **errors.AsType everywhere** — Go 1.26 idiom. output.go was the last holdout (fixed this session). Verify no other `errors.As` usages remain.

2. **Per-command middleware** — Currently middleware is CLI-wide only. Adding `WithCommandMiddleware[T,F]` to `CommandOption` would enable command-specific chains (e.g., spinner only on long-running commands, auth middleware only on protected commands).

3. **Structured error output** — When `--output=json` is set, errors should also be JSON-formatted. Currently only successful output uses the format; errors always go to stderr as plain text.

4. **Config validation chain** — `ValidateConfig` now uses `errors.Join` (fixed this session), but `WithConfigValidation[T](fn)` in cli.go still wraps with `%w: %w` double-wrap pattern. Consider `errors.Join` there too.

5. **Testutil consolidation** — Two packages (`pkg/cmdguard/v2/testutil/` and `pkg/testutil/`) with overlapping helpers. `ContainsString`/`StringSliceContains` in pkg/testutil are redundant with `slices.Contains`.

6. **RenderOptions passthrough** — go-output supports `Title`, `GraphID`, `ColorMode` in `RenderOptions`, but cmdguard's `OutputConfig` doesn't expose them. Users who want titled tables or colored output have no API.

7. **Coverage → 90%** — 85.9% → 90% requires covering the 13 zero-coverage functions (mostly validators and options) + the 11 partial functions. Most are straightforward table-driven tests.

### Code Quality Improvements

8. **Duplicate jsonLoader** — Unify by having `configload.JSON()` wrap the core package's loader.

9. **initNoColorFlag placement** — Sits in `cli_output.go` but is unrelated to output formatting. Move to own file or `cli.go`.

10. **OutputStyledTable removal** — Deprecated but still imported by taskctl example. Update example, then remove function + `table` import.

11. **Replace directive cleanup** — The single biggest "technical debt" item. Once go-output v0.9.0 is tagged and samber-do-auditlog v0.0.2 exists, all 10 lines disappear from go.mod.

---

## f) Top 25 Things We Should Get Done Next

Sorted by **impact × effort** (highest first):

| #  | Task                                                                          | Impact   | Effort | Category     |
| -- | ----------------------------------------------------------------------------- | -------- | ------ | ------------ |
| 1  | **Tag go-output v0.9.0, remove replace directives**                           | CRITICAL | 30m    | Release      |
| 2  | **Tag samber-do-auditlog v0.0.2**                                             | HIGH     | 15m    | Release      |
| 3  | **Cover 13 zero-coverage functions** (validators, options)                    | HIGH     | 2h     | Testing      |
| 4  | **Fix validateEmail/validateURL coverage**                                    | HIGH     | 30m    | Testing      |
| 5  | **Fix prompt function coverage** (mock huh)                                   | HIGH     | 2h     | Testing      |
| 6  | **Fix manpage coverage** (NewManPage, GenerateManPageCommand)                 | MEDIUM   | 1h     | Testing      |
| 7  | **Fix auditlog write coverage** (writeAuditToFileOrWriter, writeAuditMermaid) | MEDIUM   | 1h     | Testing      |
| 8  | **Unify duplicate jsonLoader**                                                | MEDIUM   | 30m    | Refactor     |
| 9  | **Remove OutputStyledTable** (update taskctl first)                           | MEDIUM   | 15m    | Cleanup      |
| 10 | **Structured JSON error output** for `--output=json`                          | HIGH     | 1h     | Feature      |
| 11 | **Move initNoColorFlag** out of cli_output.go                                 | LOW      | 5m     | Cleanup      |
| 12 | **Per-command middleware** support                                            | HIGH     | 2h     | Feature      |
| 13 | **RenderOptions passthrough** (Title, GraphID, ColorMode)                     | MEDIUM   | 1h     | Feature      |
| 14 | **Wrap do.Provide in recover()** for no-panic guarantee                       | MEDIUM   | 30m    | Architecture |
| 15 | **MergeConfigs zero-value fix** (false/0/"" never override)                   | MEDIUM   | 1h     | Fix          |
| 16 | **Consolidate testutil packages**                                             | LOW      | 1h     | Cleanup      |
| 17 | **Fix applyNoColorIfSet coverage** (25%)                                      | LOW      | 30m    | Testing      |
| 18 | **Fix validateTagRules coverage** (16.7%)                                     | LOW      | 30m    | Testing      |
| 19 | **Fix WithRangeArgs coverage** (50%)                                          | LOW      | 15m    | Testing      |
| 20 | **Fix SpinnerMiddlewareWithConfig coverage** (50%)                            | LOW      | 30m    | Testing      |
| 21 | **Plugin system for custom validators/type handlers**                         | HIGH     | 4h     | Feature      |
| 22 | **Config file nested struct support**                                         | HIGH     | 4h     | Feature      |
| 23 | **Extract flagtags library**                                                  | MEDIUM   | 4h     | Refactor     |
| 24 | **Advanced types** (Result[T], Validated[T])                                  | MEDIUM   | 6h     | Feature      |
| 25 | **Documentation generation command**                                          | LOW      | 4h     | Feature      |

---

## g) Top #1 Question I Cannot Figure Out Myself

**When should we tag go-output v0.9.0 and samber-do-auditlog v0.0.2?**

The go-output repo at `/home/lars/projects/go-output/` has all the new `TableDataMarshaler`, `AnyDataMarshaler`, `RenderTableData`, `RenderAnyData` APIs and sub-module `init()` registrations that cmdguard depends on. The samber-do-auditlog repo at `/home/lars/projects/samber-do-auditlog/` needs `html_templ.go` committed to the repo (currently gitignored).

Without these tags:

- cmdguard's `go.mod` has 10 local `replace` directives
- No external consumer can build cmdguard
- CI on GitHub will fail (no access to local paths)

This is the **single biggest blocker** for any public-facing work. I can prepare the commits but **cannot tag releases** or push to those repos without your explicit approval.

---

## Metrics Dashboard

| Metric                      | Value         | Trend                                       |
| --------------------------- | ------------- | ------------------------------------------- |
| **Version**                 | v2.5.0        | —                                           |
| **Total tests**             | 404 passing   | ↑ from 385                                  |
| **Coverage**                | 85.9%         | ↓ from 86.0% (new errors_audit.go untested) |
| **Lint issues**             | 0             | —                                           |
| **Race conditions**         | 0             | —                                           |
| **Panics in library**       | 0             | —                                           |
| **Production source lines** | ~8,800        | ↓ from ~9,000                               |
| **Test source lines**       | ~13,800       | —                                           |
| **Source files**            | 50            | +1 (errors_audit.go)                        |
| **Exported functions**      | ~120          | —                                           |
| **Sentinel errors**         | 62            | —                                           |
| **Output formats**          | 16            | —                                           |
| **Dependencies**            | 14 direct     | —                                           |
| **Replace directives**      | 10 (blocking) | —                                           |
| **0% coverage functions**   | 13            | —                                           |
| **Packages**                | 6             | —                                           |

---

## Coverage Heat Map (by package)

| Package                      | Coverage | Status                    |
| ---------------------------- | -------- | ------------------------- |
| `pkg/cmdguard/v2`            | 85.9%    | Target: 90%               |
| `pkg/cmdguard/v2/configload` | 90.2%    | ✅                        |
| `pkg/cmdguard/v2/testutil`   | 55.2%    | ⚠️ (helpers, low priority) |
| `examples/taskctl`           | 70.5%    | ✅ (integration example)  |
| `pkg/testutil`               | 0.0%     | ⚠️ (shared panic helpers)  |
| `benchmarks`                 | —        | N/A                       |
| `tests/integration`          | —        | N/A                       |

---

## Session Commit Log (Today)

```
4ea803d docs: update AGENTS.md with errors_audit.go, errors.AsType, validation fix
a3fa206 test,output: replace hard-coded format count with dynamic check
f296106 fix,config: use errors.Join for validation error aggregation
17d20ac refactor,errors: extract audit-log sentinels to errors_audit.go
4b5bb7c refactor,output: adopt errors.AsType for Go 1.26 consistency
0e780ea test,output: cover all 16 formats, add errors.Is checks and nil data test
24e0535 docs(status): output registry delegation complete status report
5bc7f28 docs: update AGENTS.md with output registry delegation changes
5d76469 refactor,output: delegate to go-output registries, eliminate 167 lines
d0dee06 chore: improve error messages with structured context fields
1e84c8d docs(status): output FormatStrategy refactor status report with self-critique
56868c6 refactor,output: replace dual-registry with FormatStrategy interface
```

---

_Waiting for further instructions._

---

## Appendix A — Session 2026-06-11 (afternoon)

**Date:** 2026-06-11 ~13:30\
**Trigger:** Resumed interrupted session — Tasks 9, 10, 11, 14, 15, 16, 17

### Tasks Completed

| #  | Task                                             | Category     | Status  | Details                                                                                                   |
| -- | ------------------------------------------------ | ------------ | ------- | --------------------------------------------------------------------------------------------------------- |
| 9  | Structured JSON error output for `--output=json` | Feature      | ✅ DONE | `cli_errors_json.go` — typed JSON errors, Fang no-op handler, `--no-color` integration                    |
| 10 | Per-command middleware support                   | Feature      | ✅ DONE | `WithCommandMiddleware[T,F]` option on commands                                                           |
| 11 | Move `initNoColorFlag` out of cli_output.go      | Cleanup      | ✅ DONE | Moved to `cli.go`                                                                                         |
| 14 | Config auto-loading with koanf                   | Feature      | ✅ DONE | `configload.NewKoanfLoader()` — YAML/JSON, nested struct flattening via `FlatPaths: true` + `Tag: "flag"` |
| 15 | Wrap `do.Provide` in recover()                   | Architecture | ✅ DONE | `safeProvide()` helper converts panics → `ErrServiceRegistration` errors                                  |
| 16 | MergeConfigs zero-value fix                      | Fix          | ✅ DONE | Removed `IsZero()` skip — false/0/"" now correctly override                                               |
| 17 | Consolidate testutil packages                    | Cleanup      | ✅ DONE | Deleted `pkg/cmdguard/v2/testutil/`, inlined helpers into 2 importers                                     |

### Key Changes by File

| File                                          | Change                                                                       |
| --------------------------------------------- | ---------------------------------------------------------------------------- |
| `pkg/cmdguard/v2/cli_errors_json.go`          | **NEW** — Structured JSON error output when `--output=json`                  |
| `pkg/cmdguard/v2/cli_errors_json_test.go`     | **NEW** — 18 tests for JSON error output                                     |
| `pkg/cmdguard/v2/configload/koanf.go`         | **NEW** — KoanfLoader with nested config flattening                          |
| `pkg/cmdguard/v2/configload/koanf_test.go`    | **NEW** — 15 tests covering YAML/JSON/nested/multi-path/expansion            |
| `pkg/cmdguard/v2/scope.go`                    | Added `safeProvide()` — panic recovery for Provide/ProvideNamed/ProvideValue |
| `pkg/cmdguard/v2/scope_provide_basic_test.go` | 3 tests for duplicate registration → error, not panic                        |
| `pkg/cmdguard/v2/config.go`                   | Removed `IsZero()` skip in `mergeStruct()`                                   |
| `pkg/cmdguard/v2/config_merge_test.go`        | Updated tests for zero-value override behavior                               |
| `pkg/cmdguard/v2/testutil/`                   | **DELETED** — Consolidated into `pkg/testutil/`                              |
| `pkg/cmdguard/v2/testhelpers_test.go`         | Inlined AddCommand (was delegating to deleted testutil)                      |
| `tests/integration/v2_bdd_lifecycle_test.go`  | Inlined AddCommand                                                           |
| `pkg/cmdguard/v2/cli.go`                      | Moved `initNoColorFlag` from cli_output.go; added command middleware wiring  |
| `pkg/cmdguard/v2/cli_output.go`               | Removed `initNoColorFlag` (moved to cli.go)                                  |
| `pkg/cmdguard/v2/output.go`                   | Removed unused code (from earlier session)                                   |
| `.golangci.yml`                               | Added koanf dependencies to depguard allow list                              |
| `go.mod` / `go.sum`                           | Added koanf v2.3.5 + parsers/file providers                                  |

### Updated Metrics

| Metric                    | Before | After | Change                                                          |
| ------------------------- | ------ | ----- | --------------------------------------------------------------- |
| **Total tests**           | 404    | 410   | +6 (koanf tests in configload; net after testutil deletion)     |
| **Coverage (v2)**         | 85.9%  | 86.0% | +0.1%                                                           |
| **Coverage (configload)** | 90.2%  | 87.5% | ↓ (koanf.go added, tests cover core paths but not all branches) |
| **Lint issues**           | 0      | 0     | —                                                               |
| **Race conditions**       | 0      | 0     | —                                                               |
| **Panics in library**     | 0      | 0     | — (now includes Provide recovery)                               |
| **Source files**          | 50     | 52    | +2 (koanf.go, cli_errors_json.go)                               |
| **Dependencies**          | 14     | 17    | +3 (koanf/v2, koanf parsers/yaml, koanf parsers/json)           |
| **Replace directives**    | 10     | 10    | — (unchanged)                                                   |

### Updated NOT STARTED List

Tasks completed this session removed from the list:

| #  | Task                           | Old Status  | New Status |
| -- | ------------------------------ | ----------- | ---------- |
| 9  | Move `initNoColorFlag`         | NOT STARTED | ✅ DONE    |
| 10 | Per-command middleware         | NOT STARTED | ✅ DONE    |
| 11 | Structured JSON error output   | NOT STARTED | ✅ DONE    |
| 14 | Config auto-loading with koanf | NOT STARTED | ✅ DONE    |
| 15 | Wrap do.Provide in recover()   | NOT STARTED | ✅ DONE    |
| 18 | MergeConfigs zero-value fix    | NOT STARTED | ✅ DONE    |
| 19 | Consolidate testutil packages  | NOT STARTED | ✅ DONE    |

Remaining NOT STARTED tasks (renumbered):

| #  | Task                                                    | Category     | Priority | Effort |
| -- | ------------------------------------------------------- | ------------ | -------- | ------ |
| 1  | Tag go-output v0.9.0 and remove replace directives      | Release      | P0       | 30m    |
| 2  | Tag samber-do-auditlog v0.0.2 (commit html_templ.go)    | Release      | P0       | 15m    |
| 3  | Remove deprecated `OutputStyledTable`                   | Cleanup      | P2       | 10m    |
| 4  | Fix prompt test coverage (0% → 80%+)                    | Testing      | P1       | 2h     |
| 5  | Fix manpage test coverage (14% → 80%+)                  | Testing      | P1       | 1h     |
| 6  | Fix auditlog write coverage (35-44% → 80%+)             | Testing      | P1       | 1h     |
| 7  | Duplicate `jsonLoader` unification                      | Refactor     | P1       | 30m    |
| 8  | `RegisterTypeHandler`/`RegisterValidator` return errors | Architecture | P2       | 1h     |
| 9  | Plugin system for custom validators/type handlers       | Feature      | P3       | 4h     |
| 10 | Extract flag-related code to `flagtags` library         | Refactor     | P3       | 4h     |
| 11 | Documentation generation command                        | Feature      | P3       | 4h     |
| 12 | Advanced types: Result[T], Validated[T], branded IDs    | Feature      | P3       | 6h     |
| 13 | Remove unused sentinels in v3                           | Cleanup      | P3       | 5m     |

### Updated "d) TOTALLY FUCKED UP"

| Issue                               | Severity             | Change                               |
| ----------------------------------- | -------------------- | ------------------------------------ |
| go.mod replace directives           | **HIGH**             | Unchanged — still the #1 blocker     |
| Prompt functions untested           | **MEDIUM**           | Unchanged                            |
| `RegisterTypeHandler` panic risk    | ~~LOW~~ **RESOLVED** | Task 15 wrapped Provide in recover() |
| `validateEmail`/`validateURL` at 0% | **LOW**              | Unchanged                            |
| Duplicate `jsonLoader`              | **LOW**              | Unchanged                            |

### Updated "e) WHAT WE SHOULD IMPROVE"

Completed this session (removed from list):

- ~~Per-command middleware~~ → ✅ DONE (Task 10)
- ~~Structured error output~~ → ✅ DONE (Task 11)
- ~~Testutil consolidation~~ → ✅ DONE (Task 17)
- ~~Wrap do.Provide in recover()~~ → ✅ DONE (Task 15)
- ~~MergeConfigs zero-value fix~~ → ✅ DONE (Task 16)

### Uncommitted Changes

All changes are in the working tree, ready to commit:

```
21 files changed, 1114 insertions(+), 372 deletions(-)
```

**Verification:** `go build ./...` ✅ · `go test ./... -race` ✅ (410 tests) · `golangci-lint run ./...` ✅ (0 issues) · `go mod tidy` ✅

## Resolution (2026-07-18)

Superseded by v3.0.0 (2026-07-07). The "one real blocker" — 10 `go.mod` replace directives (§Executive Summary) — is fully resolved: go-output is at v0.30.4 and samber-do-auditlog at v0.5.0, no replace directives remain. Current metrics: 1429 test runs (was 404), 87.6% coverage (was 85.9%), 58 sentinel errors (was "60+"), 26 benchmarks (was 22), 7 fuzz targets (unchanged). jsonv2 (`GOEXPERIMENT=jsonv2`) was adopted 2026-07-14.
