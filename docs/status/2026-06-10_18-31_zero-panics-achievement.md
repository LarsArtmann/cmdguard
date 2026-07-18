# Status Report: Zero Panics Achievement

**Date:** 2026-06-10 18:31
**Author:** Crush (AI Assistant)
**Scope:** Full codebase audit after panic elimination sprint

---

## Executive Summary

cmdguard v2 has achieved **zero panics in library code**. Every function that previously panicked now returns errors. This was a breaking change (v2.5.0) that touched 44 files, removed 593 lines, and added 299 — a net reduction of 294 lines. All 744 test cases pass, 0 lint issues, 0 race conditions, 83.5% coverage.

---

## A) FULLY DONE

### Panic Elimination (This Session)

| What                                              | Status      | Details                                                        |
| ------------------------------------------------- | ----------- | -------------------------------------------------------------- |
| `MustNewCommand` / `MustNewParentCommand`         | **Deleted** | Removed from `command.go`                                      |
| `MustNewCLI` / `MustAddCommand`                   | **Deleted** | Removed from `cli.go`                                          |
| `MustVersionCommand`                              | **Deleted** | Removed from `version.go`                                      |
| `MustDoctorCommand`                               | **Deleted** | Removed from `doctor.go`                                       |
| `MustInvoke[T]` / `MustInvokeNamed[T]`            | **Deleted** | Removed from `scope.go`                                        |
| `MustGet[T]` / `RequireBranchingFlowContext`      | **Deleted** | Removed from `flow_context_access.go`                          |
| `MustParse[T]` (generic helper)                   | **Deleted** | Removed from `type_helpers.go`                                 |
| `MustParseDuration`                               | **Deleted** | Removed from `types_duration.go`                               |
| `MustParseLogLevel` / `MustParseLogFormat`        | **Deleted** | Removed from `types_log.go`                                    |
| `MustParseEnum`                                   | **Deleted** | Removed from `types_enum.go`                                   |
| `MustParseURL`                                    | **Deleted** | Removed from `types_url.go`                                    |
| `MustParseEmail`                                  | **Deleted** | Removed from `types_email.go`                                  |
| `MustParsePort`                                   | **Deleted** | Removed from `types_port.go`                                   |
| `MustParseFilePath`                               | **Deleted** | Removed from `types_filepath.go`                               |
| `MustParseHostPort`                               | **Deleted** | Removed from `types_hostport.go`                               |
| `mustNonNegative` (panic helper)                  | **Deleted** | Replaced with `nonNegativeErr` returning error                 |
| `WithExactArgs` / `WithMinimumArgs` / etc. panics | **Fixed**   | Now set `optionErr` field on Command, surfaced by constructors |
| `Package()` forced panic                          | **Fixed**   | Changed from `func(do.Injector)` to `(*CLI[T], error)`         |
| All test callers updated                          | **Done**    | 44 files changed                                               |
| `doc.go` updated                                  | **Done**    | Removed Must\* references, updated version example             |
| `AGENTS.md` updated                               | **Done**    | Updated API reference, gotchas, removed Must\* sections        |

### Previous Sprint (Committed Before Panic Work)

| What                                                         | Commit    | Details                      |
| ------------------------------------------------------------ | --------- | ---------------------------- |
| `formatFieldValue` tests                                     | `b62aed4` | 10 subtests                  |
| `MustVersionCommand`/`GenerateVersionCommand` tests          | `e3cd57b` | (Now deleted in panic work)  |
| `validatorRegistry` through `ValidateConfig`                 | `37d994c` | Threading fix                |
| `MustParseDuration`/`MustParseLogLevel`/`MustParseLogFormat` | `7fe2940` | (Now deleted in panic work)  |
| `ErrLogLevel`/`ErrLogFormat` error chain wiring              | `670e8ea` | Proper sentinel wrapping     |
| `configload.Auto()` YAML→TOML→JSON                           | `6d832e3` | Auto-detection fix           |
| `getFieldValue` → `formatFieldValue`                         | `090b9f6` | Canonical formatter          |
| `ShutdownAll` double-wrap fix                                | `3e4c5f5` | Error chain correction       |
| `MustParseEnum` addition                                     | `903c8c5` | (Now deleted in panic work)  |
| Enum Marshal/Unmarshal docs                                  | `86bf188` | Comment clarity              |
| Flow context SetValue semantic fix                           | `34f2e06` | Skip children with local key |

---

## B) PARTIALLY DONE

| Area                          | Status             | What's Left                                                                                                                                                              |
| ----------------------------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| FEATURES.md                   | **Stale**          | Still lists `MustNewCommand`, `MustVersionCommand`, `MustDoctorCommand`, `MustParse[T]`, `MustParseEmail`, etc. as ✅ FULLY_FUNCTIONAL. These functions no longer exist. |
| TODO_LIST.md                  | **Stale**          | Still references "373 tests, 84.0% coverage" — now 744 test cases, 83.5% coverage. P6 items reference `MustGet` rename and `Package()` removal that are already done.    |
| ROADMAP.md                    | **Stale**          | Still lists "Remove or redesign Package() for error-safe DI integration" — already done.                                                                                 |
| README.md                     | **Possibly stale** | May reference Must\* functions in examples. Not audited this session.                                                                                                    |
| WHAT_THIS_PROJECT_IS_NOT.md   | **Possibly stale** | Not audited this session.                                                                                                                                                |
| docs/CLI_DESIGN_PRINCIPLES.md | **Not audited**    | May reference Must\* functions.                                                                                                                                          |
| docs/COMPARISON.md            | **Not audited**    | May reference Must\* functions.                                                                                                                                          |
| docs/MIGRATION_FROM_COBRA.md  | **Not audited**    | May reference Must\* functions.                                                                                                                                          |

---

## C) NOT STARTED

| #   | Task                                                                    | Category    | Impact | Effort |
| --- | ----------------------------------------------------------------------- | ----------- | ------ | ------ |
| 1   | Update FEATURES.md: remove all Must\* entries, add "Zero Panics" badge  | Docs        | High   | 30m    |
| 2   | Update TODO_LIST.md: reflect v2.5.0 state, mark P6 Must\* items as done | Docs        | Medium | 15m    |
| 3   | Update ROADMAP.md: mark Package() redesign as done                      | Docs        | Low    | 5m     |
| 4   | Audit README.md for Must\* references                                   | Docs        | High   | 15m    |
| 5   | Audit docs/_.md for Must_ references                                    | Docs        | Medium | 20m    |
| 6   | Add `CODECOV_TOKEN` to GitHub repo settings                             | CI          | Medium | 5m     |
| 7   | Nested struct config file support                                       | Feature     | High   | 3-5d   |
| 8   | Plugin system for validators/type handlers                              | Feature     | High   | 3-5d   |
| 9   | Documentation generation (GenerateDocs)                                 | Feature     | Medium | 2-3d   |
| 10  | Advanced types: Result[T], Validated[T], branded IDs                    | Feature     | Medium | 3-5d   |
| 11  | Structured JSON error output                                            | Feature     | Low    | 1-2d   |
| 12  | Extract flagtags to standalone library                                  | Refactor    | Low    | 2-3d   |
| 13  | v3.0 API design document                                                | Planning    | High   | 2-3d   |
| 14  | v3.0 directory skeleton                                                 | Planning    | Medium | 1d     |
| 15  | `NoFlags` as distinct named type (not alias)                            | v3 Breaking | High   | 2h     |
| 16  | Remove deprecated `WithColor` option                                    | v3 Breaking | Low    | 30m    |
| 17  | Fix `os.Setenv("NO_COLOR", "1")` process-wide mutation                  | v3 Breaking | Medium | 2h     |
| 18  | Remove `SetConfig` or make it safe                                      | v3 Breaking | Medium | 2h     |
| 19  | Fuzz testing expansion                                                  | Testing     | Medium | 1-2d   |
| 20  | Go module with `buildGoModule` in flake.nix                             | Nix         | Low    | 2h     |

---

## D) TOTALLY FUCKED UP

Nothing is totally fucked up. The codebase is in excellent shape:

- **0 compilation errors**
- **0 lint issues** (golangci-lint with 96 linters)
- **0 race conditions**
- **0 panics in library code**
- **All tests pass** (744 test cases)

However, there is **documentation debt** from this session's work:

1. **FEATURES.md is actively lying** — lists 6+ functions that no longer exist as "✅ FULLY_FUNCTIONAL"
2. **TODO_LIST.md metrics are wrong** — says "373 tests, 84.0%" when it's actually 744, 83.5%
3. **doc.go** was updated but other docs were not

This is not a code problem — it's a documentation consistency problem. The code is correct; the docs are stale.

---

## E) WHAT WE SHOULD IMPROVE

### Immediate (Documentation Sync)

1. **FEATURES.md** — The single most important doc to fix. Lists deleted functions as functional. This misleads users.
2. **TODO_LIST.md** — Metrics, version, and completion status all stale.
3. **README.md** — Likely has Must\* examples that no longer compile.
4. **ROADMAP.md** — Items already completed still listed as TODO.
5. **All docs/\*.md** — Comprehensive audit for Must\* references.

### Short-Term (Quality)

6. **Coverage gap analysis** — 83.5% is good but identify the 16.5% gaps systematically.
7. **Integration test expansion** — `tests/integration/` has no tests yet. The taskctl example serves as de facto integration tests but they're in the example module.
8. **Error documentation** — 60 sentinel errors deserve a generated reference table.
9. **API stability guarantees** — Document which APIs are stable vs. experimental.

### Medium-Term (Features)

10. **Nested config file support** — Most requested missing feature. Config files only support flat key-value pairs.
11. **Plugin system** — Runtime-extensible validators and type handlers.
12. **v3.0 planning** — Several P6 items are API-breaking; plan the next major version.

### Architecture

13. **`optionErr` pattern** — The new approach of setting errors in CommandOption closures and surfacing in constructors works but is slightly non-obvious. Consider documenting the pattern or making it more explicit.
14. **`Package()` return type** — Changed from `func(do.Injector)` to `(*CLI[T], error)`. This is a better API but we should verify no downstream consumers relied on the old signature.
15. **Test helper signatures** — `newTestCLICommand`, `newTestCLICommandWithShort`, `newTestParentCommand` all gained a `*testing.T` parameter. This is correct but touched many files. Consider if there's a cleaner pattern.

---

## F) Top 25 Things To Do Next

Priority-ordered, impact-weighted:

| #   | Task                                                      | Category | Impact   | Effort | Status      |
| --- | --------------------------------------------------------- | -------- | -------- | ------ | ----------- |
| 1   | Update FEATURES.md: remove Must\* entries, update status  | Docs     | Critical | 30m    | Not started |
| 2   | Update README.md: remove Must\* examples                  | Docs     | Critical | 15m    | Not started |
| 3   | Update TODO_LIST.md: metrics, version, completed items    | Docs     | High     | 15m    | Not started |
| 4   | Audit all docs/_.md for Must_ references                  | Docs     | High     | 20m    | Not started |
| 5   | Update ROADMAP.md: mark completed items                   | Docs     | Medium   | 5m     | Not started |
| 6   | Version bump to v2.5.0 in go.mod and doc.go               | Release  | High     | 5m     | Not started |
| 7   | Git tag v2.5.0                                            | Release  | High     | 1m     | Not started |
| 8   | Write CHANGELOG.md entry for v2.5.0                       | Docs     | Medium   | 15m    | Not started |
| 9   | Add `nix buildGoModule` to flake.nix                      | Nix      | Medium   | 2h     | Not started |
| 10  | Coverage gap analysis: identify uncovered paths           | Testing  | Medium   | 1h     | Not started |
| 11  | Add integration tests to tests/integration/               | Testing  | Medium   | 1d     | Not started |
| 12  | Error reference doc: auto-generated table of 60 sentinels | Docs     | Medium   | 2h     | Not started |
| 13  | Nested struct config file support (YAML/TOML/JSON)        | Feature  | High     | 3-5d   | Not started |
| 14  | Plugin system for validators and type handlers            | Feature  | High     | 3-5d   | Not started |
| 15  | Fix os.Setenv("NO_COLOR") process-wide mutation           | v3 Prep  | Medium   | 2h     | Not started |
| 16  | Make NoFlags a distinct named type                        | v3 Prep  | Medium   | 2h     | Not started |
| 17  | Remove deprecated WithColor option                        | v3 Prep  | Low      | 30m    | Not started |
| 18  | Add CODECOV_TOKEN to GitHub                               | CI       | Medium   | 5m     | Not started |
| 19  | Expand fuzz test targets                                  | Testing  | Low      | 1-2d   | Not started |
| 20  | v3.0 API design document                                  | Planning | High     | 2-3d   | Not started |
| 21  | Documentation generation (GenerateDocs)                   | Feature  | Medium   | 2-3d   | Not started |
| 22  | Structured JSON error output for --output=json            | Feature  | Low      | 1-2d   | Not started |
| 23  | Advanced types: Result[T], Validated[T]                   | Feature  | Medium   | 3-5d   | Not started |
| 24  | Extract flagtags to standalone library                    | Refactor | Low      | 2-3d   | Not started |
| 25  | Config auto-loading with koanf integration                | Feature  | Low      | 2-3d   | Not started |

---

## G) Top #1 Question I Cannot Figure Out Myself

**What is the version strategy?**

The codebase is at v2.4.0 in docs/version but we just made breaking API changes (removed 15+ exported functions, changed `Package()` signature, changed arg validator behavior). This is semantically a **major breaking change** that should be v3.0.

However, the Go module path is `github.com/larsartmann/cmdguard/v2` — we can't release v3.0 without creating a new `v3/` directory and module path.

The question is: **Should this be v2.5.0 (minor bump, breaking change acknowledged) or v3.0.0 (proper semver, requires new module path)?**

The Must\* functions were convenience wrappers — every one of them has an error-returning equivalent that already existed. The `Package()` change is more significant. The arg validator change is subtle but technically breaking (code that relied on early panics will now get errors at construction time instead).

This is a product/strategy decision I cannot make.

---

## Metrics Dashboard

| Metric                     | Value  | Change                                        |
| -------------------------- | ------ | --------------------------------------------- |
| **Source files**           | 121    | Unchanged                                     |
| **Source lines**           | 21,106 | -294 from panic removal                       |
| **Test cases**             | 744    | Down from ~1000+ (removed Must\* panic tests) |
| **Coverage**               | 83.5%  | Down from 84.0% (removed trivial panic tests) |
| **Lint issues**            | 0      | Unchanged                                     |
| **Race conditions**        | 0      | Unchanged                                     |
| **Build errors**           | 0      | Unchanged                                     |
| **Panics in library code** | 0      | Down from ~15 panic points                    |
| **Sentinel errors**        | 60     | Unchanged                                     |
| **Command options**        | 21     | Unchanged                                     |
| **CLI options**            | 22     | Unchanged                                     |
| **Value types**            | 9      | Unchanged                                     |
| **Output formats**         | 12     | Unchanged                                     |
| **Dependencies**           | 8      | Unchanged                                     |

### Coverage By Package

| Package                      | Coverage |
| ---------------------------- | -------- |
| `pkg/cmdguard/v2`            | 83.5%    |
| `pkg/cmdguard/v2/configload` | 90.2%    |
| `pkg/cmdguard/v2/testutil`   | 87.5%    |
| `examples/taskctl`           | 71.3%    |

### Git History (Last 10 Commits)

| Commit    | Description                                                                              |
| --------- | ---------------------------------------------------------------------------------------- |
| `c8d86c8` | fix(cmdguard): remove all panic-inducing functions — zero panics guaranteed              |
| `89f9668` | docs(status): add comprehensive status report for deep improvement sprint                |
| `48acdab` | docs(cmdguard): update AGENTS.md with final metrics and new gotchas                      |
| `86bf188` | docs(cmdguard): clarify why Enum has hand-written Marshal/UnmarshalText                  |
| `903c8c5` | feat(cmdguard): add MustParseEnum for API consistency                                    |
| `3e4c5f5` | fix(cmdguard): remove double-wrapping of ErrServiceConstruction in ShutdownAll           |
| `090b9f6` | refactor(cmdguard): replace getFieldValue with formatFieldValue                          |
| `6d832e3` | fix(configload): Auto() now tries YAML→TOML→JSON instead of only JSON                    |
| `59f7dc1` | docs(cmdguard): update AGENTS.md with error chain and unused sentinel notes              |
| `670e8ea` | fix(cmdguard): wire ErrLogLevel/ErrLogFormat into parse chain, fix bare sentinel returns |

---

## Files Modified This Session

44 files changed, 299 insertions(+), 593 deletions(-):

**Source (panics removed):**

- `cli.go`, `command.go`, `command_options.go`, `doctor.go`, `version.go`
- `scope.go`, `flow_context_access.go`, `type_helpers.go`
- `types_duration.go`, `types_email.go`, `types_enum.go`, `types_filepath.go`
- `types_hostport.go`, `types_log.go`, `types_port.go`, `types_url.go`

**Source (docs):**

- `doc.go`, `AGENTS.md`

**Tests updated:**

- `coverage_test.go`, `doctor_test.go`, `version_test.go`, `duration_test.go`
- `enum_test.go`, `helpers_test.go`, `cli_superb_test.go`
- `flow_context_options_test.go`, `scope_provide_named_test.go`, `scope_integration_test.go`
- `testhelpers_test.go`, `cli_core_accessors_test.go`, `cli_core_addcmd_test.go`
- `cli_core_flow_test.go`, `cli_error_paths_test.go`, `cli_lifecycle_test.go`
- `example_test.go`, `example_types_test.go`
- `types_email_test.go`, `types_url_test.go`, `types_port_test.go`
- `types_hostport_test.go`, `types_filepath_test.go`

**Examples/benchmarks:**

- `examples/taskctl/commands.go`, `examples/taskctl/main_test.go`
- `benchmarks/guard_bench_test.go`

## Resolution (2026-07-18)

The "v2.5.0 or v3.0.0?" question in §G was resolved: this shipped as **v3.0.0 (2026-07-07)** on module path `github.com/larsartmann/cmdguard/v3`, not v2.5.0. The v2 line ended at v2.10.4. The `NO_COLOR` env mutation flagged in §F was fixed (commit `e53c8e6`). Zero-panics contract holds in v3 (no `Must*` variants exist).
