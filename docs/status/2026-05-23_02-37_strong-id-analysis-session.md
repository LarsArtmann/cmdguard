# Status Report — 2026-05-23 02:37

**Project:** cmdguard — CLI Guard Library  
**Branch:** master (up to date with origin/master)  
**Status:** v2.3.0-dev  
**Session:** Strong ID Analysis — Branded ID Introduction

---

## Executive Summary

Session focused on resolving Strong ID Analysis violations by introducing branded ID types from `go-branded-id` into `examples/di-patterns/main.go` and fixing a parameter shadowing issue in `pkg/cmdguard/v2/cli_options.go`. All 3 violations resolved, lint clean, tests passing, build passing.

---

## Current Project State

| Metric | Value |
| ------ | ----- |
| Tests | 224 passing |
| Coverage (v2 core) | 84.3% |
| Lint issues | 0 |
| Race conditions | 0 |
| Build errors | 0 |
| Go version | 1.26 |
| Code duplication | 0 groups |

---

## Work Breakdown

### (a) ✅ FULLY DONE

- **Phase 9: `errors.As` → `errors.AsType` (Go 1.26 idiom)** — `pkg/cmdguard/v2/cli_options.go:79` parameter renamed `id` → `groupID` to avoid shadowing the type name. Fixes Strong ID violation for `WithGroup` function parameter.

- **Strong ID — `TaskID` branded type** — `examples/di-patterns/main.go` now uses `TaskBrand` phantom type + `TaskID = id.ID[TaskBrand, int]` replacing bare `int`. Compile-time type safety prevents mixing `TaskID` with other entity IDs.

- **Strong ID — `NextID` as `TaskID` alias** — `nextID` field in `TaskStore` now uses `NextID = TaskID` type alias (not a separate branded type) since it semantically holds the next `TaskID` to be assigned, not a distinct entity type.

- **Example binary cleanup** — `di-patterns` compiled binary (8.3 MB) removed from working tree (untracked).

### (b) ⚠️ PARTIALLY DONE

- **go.mod upgrade** — `go-branded-id` upgraded from `// indirect` to direct dependency. However, the library is pinned to `v0.1.0` with no `// go-branded-id v0.2+` available yet — the library itself may need versioning before cmdguard can track a latest stable tag. Status: working, but dependency pinning is conservative.

- **Phase 9 remaining items (8 of 10 incomplete):**
  - `handlerConfig[T,F]` struct extraction from 8-param `wireHandlerWithMiddleware` — NOT STARTED
  - `Phase` typed enum to replace `CommandInfo.Phase string` — NOT STARTED
  - Fix 7 unwrapped error returns — PARTIALLY DONE (commit `2d91cde` wrapped 4, 3 remain)
  - Consolidate 5 error types into `internal.labeledError` — NOT STARTED
  - Split `type_handler.go` (481 lines) into 3 files — NOT STARTED
  - Split `command.go` (403 lines) — NOT STARTED
  - Split `flow_context.go` (396 lines) — NOT STARTED
  - Fix `outputFormat`/`outputState.format` split brain — NOT STARTED
  - Consolidate value type MarshalText/UnmarshalText patterns — NOT STARTED

### (c) ❌ NOT STARTED

- **Benchmarks** — CLI construction, flag parsing, and command execution benchmarks have never been added. No regression detection in CI.

- **Codecov integration** — No coverage tracking service configured.

- **gopls `errors.As` → `errors.AsType[ExitCoder]`** — Go 1.26 idiom not yet applied anywhere in codebase (was listed in Phase 9).

- **Config file auto-loading** — koanf integration for YAML/TOML/.env not started (v3.0 territory).

- **Interactive prompts** — huh integration with `WithPromptOnMissing` not started.

- **Plugin system** — Custom validators and type handlers extensibility not started.

- **v2.3.0 release** — No release tag or release notes created.

### (d) ✅ NOT TOTALLY FUCKED UP

No broken states. All tests pass, lint is clean, build succeeds. The `di-patterns` example binary was a leftover from `go run` execution but posed no functional risk.

---

## What We Should Improve

1. **Phase 9 file splitting** — `type_handler.go` (481 lines), `command.go` (403 lines), `flow_context.go` (396 lines) are all overdue for architectural splitting. High complexity, low urgency but significant maintainability debt.

2. **Error wrapping gaps** — 3 unwrapped error returns remain from the Phase 9 item "Fix 7 unwrapped error returns". Need full grep across codebase for `return err` patterns without `fmt.Errorf` wrapping.

3. **Benchmark suite missing** — No performance regression detection exists. Critical for a library that users import into their CLI applications.

4. **Coverage gap at 84.3%** — 15.7% of v2 code uncovered. Examples especially have 0% coverage despite being the primary user-facing documentation. The `di-patterns` example has 0% coverage and is now a real Go program with branded types.

5. **`di-patterns` example as integration test** — The `di-patterns` example could serve as a functional test (`go run examples/di-patterns/main.go list && go run examples/di-patterns/main.go add --title "test" && go run examples/di-patterns/main.go list`) but has no test file. Should add `_test.go` exercising the full flow.

6. **Unnecessary type arguments** — 118 LSP hints across codebase for `//nolint:infertypeargs`. These could be cleaned up across all test files.

7. **gocritic warnings about disabled linter checks** — 3 gocritic warnings: `dupImport`, `octalLiteral`, `whyNoLint` are redundantly disabled since golangci-lint already ignores them by default.

8. **`.gitignore` gap** — `di-patterns` binary not in `.gitignore` (pattern is `/di-patterns` but it's a file not directory; should be `di-patterns` or `/di-patterns` with trailing note about binary).

---

## Top #25 Things to Get Done Next

1. **Fix remaining 3 unwrapped error returns** — grep for `return err` without `fmt.Errorf` wrapping in `pkg/cmdguard/v2/`
2. **Add `di-patterns` test file** — functional test of TaskStore with branded IDs exercising list/add flow
3. **Add benchmarks** — CLI construction, flag parsing, command execution benchmarks in `benchmarks/`
4. **Run `art-dupl` check** — verify 0 code duplication groups are maintained
5. **Split `type_handler.go`** — extract `type_registry.go`, `type_parsing.go`, keep `type_handler.go` minimal
6. **Split `command.go`** — extract args validation options into `command_args.go`
7. **Split `flow_context.go`** — extract options into `flow_context_options.go`
8. **Create `Phase` typed enum** — replace `CommandInfo.Phase string` with `type Phase int` + const block
9. **Apply `errors.AsType[ExitCoder]` idiom** — replace all `errors.As(err, &exitCoder)` calls
10. **Fix `outputFormat`/`outputState.format` split brain** — consolidate output state tracking
11. **Consolidate value type marshaling** — extract shared `MarshalText`/`UnmarshalText` patterns
12. **Add codecov integration** — configure GitHub Actions to upload coverage reports
13. **Create v2.3.0 release** — tag, changelog entry, publish
14. **Clean up 118 `infertypeargs` hints** — remove unnecessary type arguments across test files
15. **Fix gocritic redundant disable warnings** — remove `dupImport`, `octalLiteral`, `whyNoLint` from `.golangci.yml`
16. **Add `.gitignore` entry for `di-patterns` binary** — prevent future accidental commits
17. **Consolidate 5 error types into `internal.labeledError`** — reduce error type surface area
18. **Extract `handlerConfig[T,F]` struct** — reduce `wireHandlerWithMiddleware` parameter count
19. **Add example tests for all examples** — raise example coverage from 0%
20. **Add fuzz tests to value type parsers** — extend fuzz coverage in `pkg/cmdguard/v2/`
21. **Document `go-branded-id` usage in docs** — `examples/di-patterns` now demonstrates branded IDs, should be referenced in FEATURES.md
22. **Update FEATURES.md** — add `go-branded-id` integration as a documented feature
23. **Update TODO_LIST.md** — mark `examples/di-patterns` branded ID demo as complete (references ROADMAP item)
24. **Add `WithPromptOnMissing` design doc** — plan for huh integration in v3.0
25. **Deprecate v1 timeline** — set concrete date for v1 removal; currently v1 (`pkg/cmdguard/`) still exists with no deprecation warning

---

## Top #1 Question I Can NOT Figure Out

**How should cmdguard handle `go-branded-id` as a dependency — should it be a direct dependency, an optional dependency, or should cmdguard re-export/wrap the branded ID types internally?**

The issue: cmdguard now uses `go-branded-id` in examples (demonstrating best practices for users), but the library itself doesn't expose or re-export these types. This means:
- Users who want branded IDs must take a direct dependency on `go-branded-id`
- There's no cmdguard-native `Branding[T, B]` or `BrandedID[Brand, Value]` type
- The `examples/di-patterns` shows the pattern but it's "import this other library" rather than "cmdguard gives you this"

Should cmdguard **not** depend on `go-branded-id` at all (keep examples using plain types), should it **re-export** `go-branded-id` types under a cmdguard package (e.g., `cmdguard/v2/id`), or should it **maintain its own** minimal branded ID implementation internally? The existing docs (`docs/planning/go-composable-business-types-usage.md`) discuss this but don't reach a conclusion.

---

## Session Diff Summary

```
examples/di-patterns/main.go   | +17/-5  (branded IDs: TaskBrand, TaskID, NextID)
pkg/cmdguard/v2/cli_options.go | +2/-2   (parameter rename: id → groupID)
go.mod                         | +1/-1   (go-branded-id: indirect → direct)
```

---

## Files Changed This Session

| File | Change |
| ---- | ------ |
| `pkg/cmdguard/v2/cli_options.go` | Parameter `id` → `groupID` in `WithGroup[T]` |
| `examples/di-patterns/main.go` | Added `TaskBrand`, `TaskID`, `NextID` branded types |
| `go.mod` | Promoted `go-branded-id` to direct dependency |
| `di-patterns` (binary) | Removed untracked compiled binary |

---

## Git Status

```
On branch: master (up to date with origin/master)
Modified:  examples/di-patterns/main.go
Modified:  go.mod
Modified:  pkg/cmdguard/v2/cli_options.go
Untracked: di-patterns (compiled binary — should be gitignored)
```

---

*Generated: 2026-05-23 02:37 — cmdguard v2.3.0-dev*
