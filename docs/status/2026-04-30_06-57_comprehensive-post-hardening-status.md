# cmdguard — Comprehensive Status Report

**Date:** 2026-04-30 06:57 CEST
**Branch:** master (up to date with origin)
**Version:** v2.2.0
**Go:** 1.26

---

## Executive Summary

cmdguard is a Go library for building validated Cobra CLI applications with type-safe dependency injection. After four intensive sessions spanning 30+ commits, the project is in **excellent shape**: 881 tests passing, 80.4% coverage, 0 lint issues, 0 race conditions, clean build. This session (6 commits) focused on dead code removal, performance optimization, and API deprecation hygiene.

---

## Current Metrics

| Metric          | Value                                             | Status     |
| --------------- | ------------------------------------------------- | ---------- |
| Build           | `go build ./...` — 0 errors                       | ✅ Clean   |
| Tests (v2)      | 815 passing with `-race`                          | ✅ Clean   |
| Tests (all)     | 881 passing across all packages                   | ✅ Clean   |
| Lint            | `golangci-lint run ./...` — 0 issues              | ✅ Clean   |
| Coverage        | 80.4% (`pkg/cmdguard/v2`)                         | ✅ Good    |
| Race conditions | 0                                                 | ✅ Clean   |
| Source files    | 96 in `pkg/cmdguard/v2/` (33 production, 63 test) |            |
| Total lines     | 27,117 in v2 package (16,348 prod, 10,769 test)   |            |
| Examples        | 12 directories                                    |            |
| Dependencies    | All published (no local replace)                  | ✅ Clean   |
| CI              | GitHub Actions (build+test+lint)                  | ✅ Active  |
| Uncommitted     | `go.sum` drift only                               | ⚠️ Trivial |

---

## a) FULLY DONE

### Session 4 (This Session — 6 commits)

| #   | Change                                                                                                                                                                                                                                                      | File(s)             | Type           |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- | -------------- |
| 1   | **Removed dead code `wireSubcommandSuggestions`** — function set `FlagErrorFunc` on root that just `return err` unchanged (cobra's default). Called once per parent command, redundantly overwriting root each time.                                        | `cli_command.go`    | Dead code      |
| 2   | **Collapsed `formatFieldValue` duplicate cases** — 4 separate cases (Complex64/128, Array/Slice, Map/Struct, Chan/Func, Uintptr/UnsafePointer) all produced identical `fmt.Sprintf("%v", field.Interface())` output. Replaced with single `default` branch. | `flags_validate.go` | Dedup          |
| 3   | **Cached regex compilation in `validateRegex`** — `sync.Map` cache: first call compiles, subsequent hit cache. Optimized for read-heavy pattern (write-once, read-many).                                                                                    | `flags_validate.go` | Perf           |
| 4   | **Removed unnecessary mutex from `outputState`** — CLI execution is sequential (cobra is not concurrent). `sync.Mutex` added overhead for no benefit.                                                                                                       | `cli_output.go`     | Simplification |
| 5   | **Deprecated `IsExecutable()`** — Added `Deprecated:` doc comment pointing to `HasHandler()`. Will remove in v3.                                                                                                                                            | `command.go`        | API hygiene    |
| 6   | **Deprecated `FlowContextAccessor`/`NewFlowContextAccessor`** — Thin wrappers around `BranchingFlowContext`. Added `Deprecated:` comments. Will remove in v3.                                                                                               | `flow_context.go`   | API hygiene    |
| 7   | **Added `BranchWithDuration(name, time.Duration)`** — Typed alternative to string-based `BranchWithTimeout`.                                                                                                                                                | `flow_context.go`   | API            |
| 8   | **Added `BranchWithDeadlineTime(name, time.Time)`** — Typed alternative to string-based `BranchWithDeadline`.                                                                                                                                               | `flow_context.go`   | API            |
| 9   | **Updated AGENTS.md** — Deprecation notices, coverage metrics, gotchas for new APIs.                                                                                                                                                                        | `AGENTS.md`         | Docs           |

### Session 3 (Prior — 12 commits)

- Fixed BranchingFlowContext double-cancellation bug
- Fixed Enum.Allowed() returning mutable internal slice
- Added stack trace capture to RecoveryMiddleware
- Used `errors.Join` in Scope.ShutdownAll
- Fixed Scope.Path() allocation
- Extracted shared `lookupFlagInCommand`
- Replaced hardcoded type switch with `fmt.Stringer`
- Deduplicated type handler registrations (-52 lines)
- Simplified `parseAndSetValue` to delegate to `SetField` (-20 lines)
- Used `map[string]struct{}` for command registration set
- Added NewParentCommand example
- Added shareable pre-commit hook script

### Session 2 (Prior — v2.2.0 features)

- `env:"VAR"` struct tag support with `WithEnvPrefix`
- Subcommand typo suggestions
- `WithSignalHandling[T]()` for SIGINT/SIGTERM
- go-output integration (12 output formats)
- `count:"true"` struct tag for counting flags
- `EditInEditor()` with `context.Context`
- Shell completion wiring
- Man page generation via mango-cobra
- `WithOutputFormat[T]` CLI option
- 6 new examples
- GitHub Actions CI workflow

### Session 1 (Prior — foundation)

- Unified type dispatch into TypeHandler registry
- Fixed custom type registration in `setStringField`
- Made validator registry instance-scoped
- Fixed `Ptr[T]` function
- Cleaned up 52 artifact files
- Achieved 0 lint issues (was 113)
- Fixed all 55 race conditions
- Removed local go-output replace (tagged v0.1.0)
- Output.go format renderer registry

---

## b) PARTIALLY DONE

| Item                 | Status                                                                                                                                                                               | What's Left                        |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------- |
| **Pre-commit hooks** | Script at `scripts/pre-commit`, but must be manually copied to `.git/hooks/`. No `lefthook`/`husky` framework.                                                                       | Auto-install mechanism             |
| **Benchmarks**       | 18 benchmark functions in `benchmarks/guard_bench_test.go`. TODO_LIST says "add CLI construction/flag parsing/command execution benchmarks" — existing ones may already cover these. | Verify coverage, add if missing    |
| **API deprecation**  | `IsExecutable`, `FlowContextAccessor` marked `Deprecated:`. `WithColor` already deprecated.                                                                                          | Actual removal deferred to v3      |
| **Typed branching**  | `BranchWithDuration`/`BranchWithDeadlineTime` added as typed alternatives. String-based originals still present.                                                                     | Remove string-based variants in v3 |

---

## c) NOT STARTED

| #   | Item                                        | Effort | Impact |
| --- | ------------------------------------------- | ------ | ------ |
| 1   | Codecov integration in CI                   | Small  | Medium |
| 2   | v2.2.0 release tag and notes                | Small  | High   |
| 3   | Release automation (goreleaser?)            | Medium | Medium |
| 4   | Benchmark regression detection in CI        | Medium | Medium |
| 5   | Config file auto-loading with koanf         | Large  | High   |
| 6   | Interactive prompts (huh integration)       | Large  | Medium |
| 7   | Spinner/progress middleware (bubbles)       | Medium | Low    |
| 8   | Glamour markdown help rendering             | Medium | Medium |
| 9   | Telemetry middleware (OpenTelemetry)        | Medium | Low    |
| 10  | Plugin system for custom validators         | Large  | Medium |
| 11  | Unexport `derefPointerToStruct`             | Tiny   | Low    |
| 12  | Generate help string from valid format list | Small  | Low    |

---

## d) TOTALLY FUCKED UP — Nothing!

No regressions, no broken tests, no unfixable issues. The codebase is in its best state ever. All 4 sessions have been net-positive with zero rollbacks.

---

## e) WHAT WE SHOULD IMPROVE

### High Priority (Architecture Debt — deferred to v3)

1. **`RegisterInScope` accepts `...any`** — loses all type safety. Uses runtime type switching. Should be generic or provide typed variants. Breaking change.

2. **`Package()` panics on error** — undermines the error-safe DI pattern. The `do.Package` contract forces a void function. Inherent to samber/do integration.

3. **`NoFlags` is a type alias** (`type NoFlags = struct{}`) — users who accidentally pass an empty struct literal get `NoFlags` behavior silently. A distinct `type NoFlags struct{}` would be clearer. Breaking change.

4. **`Get[T]`/`MustGet[T]` naming** — extremely generic, collides with other packages. Should be `GetFlowValue[T]` or similar. Breaking change.

5. **`TimingMiddleware` callback should include error** — so middleware can distinguish success vs failure timing. Breaking change.

6. **`ExecuteAndExit` always exits with code 1** — no way to propagate structured exit codes. Needs new interface pattern.

### Medium Priority (Code Quality — can do now)

7. **`validateMin`/`validateMax` near-duplicates** in `flags_validate.go` — could collapse into parameterized validators. ~15 lines each with different error messages.

8. **Two separate validation execution paths** — `runValidateTag` and `parseValidateRulesWithRegistry` are different entry points for the same concept. Risk of divergence.

9. **`derefPointerToStruct` is exported** from `config_parsing.go` — internal utility that leaks implementation detail. Simple unexport.

10. **Hardcoded format help string** in `cli_output.go` — duplicates the valid format list from `ParseOutputFormat`. Should be generated.

### Low Priority (Polish)

11. **22 one-liner accessor methods** on `Command` (~60 lines of boilerplate) — could use code generation.

12. **Examples have low/no test coverage** — most examples show 0% coverage. Not critical (they're demos) but some integration test coverage would catch import/compilation regressions.

13. **`go.sum` drift** — currently uncommitted. Minor housekeeping.

---

## f) Top #25 Things to Do Next

Sorted by **impact × effort** (highest first):

| #   | Task                                                                                                            | Effort   | Impact    | Type               |
| --- | --------------------------------------------------------------------------------------------------------------- | -------- | --------- | ------------------ |
| 1   | **Tag v2.2.0 release** with CHANGELOG/release notes                                                             | 30min    | 🔴 High   | Release            |
| 2   | **Commit `go.sum` drift**                                                                                       | 2min     | 🟡 Medium | Housekeeping       |
| 3   | **Unexport `derefPointerToStruct`**                                                                             | 5min     | 🟢 Low    | Encapsulation      |
| 4   | **Verify benchmark coverage** — do existing benchmarks cover CLI construction, flag parsing, command execution? | 30min    | 🟡 Medium | Testing            |
| 5   | **Add codecov to CI**                                                                                           | 30min    | 🟡 Medium | CI                 |
| 6   | **Generate format help string from valid format list**                                                          | 30min    | 🟢 Low    | DRY                |
| 7   | **Collapse `validateMin`/`validateMax`** into parameterized validators                                          | 1hr      | 🟢 Low    | Dedup              |
| 8   | **Consolidate validation execution paths** (`runValidateTag` vs `parseValidateRulesWithRegistry`)               | 1hr      | 🟡 Medium | Architecture       |
| 9   | **Set up goreleaser** for release automation                                                                    | 2hr      | 🟡 Medium | CI                 |
| 10  | **Add benchmark regression detection** to CI                                                                    | 1hr      | 🟡 Medium | CI                 |
| 11  | **Add smoke tests for examples** (verify they compile + run)                                                    | 1hr      | 🟡 Medium | Testing            |
| 12  | **Make `RegisterInScope` generic** instead of `...any`                                                          | 1hr      | 🟡 Medium | Architecture       |
| 13  | **Make `NoFlags` a distinct named type**                                                                        | 30min    | 🟡 Medium | API (breaking)     |
| 14  | **Rename `Get[T]`/`MustGet[T]`** to `GetFlowValue[T]`/`MustGetFlowValue[T]`                                     | 15min    | 🟢 Low    | API (breaking)     |
| 15  | **Remove deprecated APIs** (`IsExecutable`, `FlowContextAccessor`, `WithColor`)                                 | 30min    | 🟢 Low    | Cleanup (breaking) |
| 16  | **Add `TimingMiddleware` error parameter**                                                                      | 30min    | 🟢 Low    | API (breaking)     |
| 17  | **Add exit code support** to `ExecuteAndExit`                                                                   | 1hr      | 🟢 Low    | API (breaking)     |
| 18  | **Config file auto-loading with koanf** (YAML/TOML/.env)                                                        | 1-2 days | 🔴 High   | Feature            |
| 19  | **Interactive prompts (huh)** with `WithPromptOnMissing`                                                        | 1-2 days | 🟡 Medium | Feature            |
| 20  | **Glamour markdown help rendering**                                                                             | 4hr      | 🟡 Medium | Feature            |
| 21  | **Telemetry middleware (OpenTelemetry spans)**                                                                  | 4hr      | 🟢 Low    | Feature            |
| 22  | **Spinner/progress middleware (bubbles)**                                                                       | 4hr      | 🟢 Low    | Feature            |
| 23  | **Plugin system for custom validators and type handlers**                                                       | 1-2 days | 🟡 Medium | Feature            |
| 24  | **Generate accessor methods** via `go generate`                                                                 | 2hr      | 🟢 Low    | Code gen           |
| 25  | **Auto-install pre-commit hook** via `go generate` or Makefile                                                  | 1hr      | 🟢 Low    | Tooling            |

---

## g) Top #1 Question I Cannot Figure Out Myself

**What is the release strategy for v2.2.0?**

The codebase is in excellent shape and has been for multiple sessions. The TODO_LIST has "Create v2.2.0 release tag and notes" as a remaining item. Specifically:

1. Should we tag `v2.2.0` now, or wait for more features (e.g., koanf config loading)?
2. Should we use goreleaser, manual tags, or GitHub Releases UI?
3. Do you want a CHANGELOG.md generated from commit history, or manually curated release notes?
4. Should the v3-breaking items (NoFlags, Get[T] rename, RegisterInScope generics) be planned as a separate milestone?

This is a product/project decision that requires your input — I can execute any of these options once decided.

---

## Session History

### Session 1 — Foundation Cleanup

- 113 lint → 0 lint
- 55 race conditions → 0
- Output.go registry refactor
- go-output v0.1.0 published

### Session 2 — v2.2.0 Features

- All v2.2.0 features complete
- 6 new examples
- CI workflow
- Documentation updated

### Session 3 — Architecture Hardening (12 commits)

- 3 bug fixes (double-cancel, Enum mutation, lost stack traces)
- 9 architecture improvements (deduplication, delegation, type safety)
- 2 new features (subcommands example, pre-commit hook)

### Session 4 — Dead Code & Deprecation Sweep (6 commits)

- Removed dead code (`wireSubcommandSuggestions`, `outputState` mutex)
- Performance (regex cache via `sync.Map`)
- Deduplication (`formatFieldValue` collapse)
- API hygiene (deprecated 3 APIs, added 2 typed alternatives)
- Documentation updated

### All Commits (30 total since work began)

```
e36a605 docs: update AGENTS.md with deprecation notices and session improvements
d922875 refactor: deprecate dead APIs, add typed alternatives, fix lint
d3d3178 refactor: remove unnecessary mutex from outputState
96c0216 perf: cache compiled regex patterns in validateRegex
95e4508 refactor: collapse formatFieldValue duplicate cases into default
10580cd refactor: remove dead code wireSubcommandSuggestions
99f5956 chore: housekeeping — format drift, gitignore, status report
770fc50 docs: update TODO_LIST.md and FEATURES.md for architecture hardening
213df86 chore: add shareable pre-commit hook script
334c0df feat: add NewParentCommand example (subcommands)
d52c4c9 refactor: simplify parseAndSetValue to delegate to SetField, document SetConfig
a1aa0f5 style: use map[string]struct{} for command registration set
48678f8 fix: include stack trace in RecoveryMiddleware panic recovery
53e8af2 refactor: deduplicate custom type handler registrations in type_handler.go
e37a69b refactor: replace hardcoded type switch in getFieldValue with fmt.Stringer
a5d098b refactor: extract shared lookupFlagInCommand from duplicated flag lookups
8eb3698 refactor: use errors.Join in Scope.ShutdownAll and fix Path() allocation
edca60d fix: return defensive copy from Enum.Allowed() to prevent mutation
6543945 fix: eliminate double-cancellation in BranchingFlowContext.Cancel()
69b0bf2 ci: streamline GitHub Actions workflow configuration
9f9a586 docs: update TODO_LIST.md, FEATURES.md, AGENTS.md for current state
62947a9 docs: update FEATURES.md and TODO_LIST.md for v2.2.0 completion
facb25f chore: remove local go-output replace, use published v0.1.0
11a67e0 fix: achieve 0 lint issues, fix all race conditions, add context to EditInEditor
17b76da refactor: simplify enumHelp closure and format log format string
12956cf chore: add linting exceptions for specific files in golangci-lint config
472888c style: auto-fix wsl_v5, nlreturn, nolintlint, modernize formatting
ca1e7bd fix: resolve race conditions, refactor output.go registry, fix depguard
1720fa2 feat: add man page generation and comprehensive status report
ba65a99 feat: add WithOutputFormat, shell completion, and improved deprecation
```

---

## Files Modified (Session 4 Only)

| File                                | Change                                                                                |
| ----------------------------------- | ------------------------------------------------------------------------------------- |
| `pkg/cmdguard/v2/cli_command.go`    | Removed `wireSubcommandSuggestions` function and call site                            |
| `pkg/cmdguard/v2/flags_validate.go` | Collapsed `formatFieldValue`, cached regex compilation                                |
| `pkg/cmdguard/v2/cli_output.go`     | Removed `sync.Mutex` from `outputState`                                               |
| `pkg/cmdguard/v2/command.go`        | Deprecated `IsExecutable`                                                             |
| `pkg/cmdguard/v2/flow_context.go`   | Deprecated `FlowContextAccessor`, added `BranchWithDuration`/`BranchWithDeadlineTime` |
| `AGENTS.md`                         | Updated metrics, deprecation notices, gotchas                                         |

---

_Generated by Crush — 2026-04-30_
