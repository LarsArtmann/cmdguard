# Smart Consolidations Status Report

**Date:** 2026-05-17
**Session:** Smart DRY consolidations + error quality + validation modes

---

## Summary

Systematic DRY consolidation pass across cmdguard v2, eliminating code duplication, fixing split brains, wiring sentinels, and adding missing validation modes. Zero feature loss, zero regressions.

---

## What Was Done

### a) Commits (this session + resumed)

| Commit | Description |
|--------|-------------|
| `eedc4e0` | Wire sentinels, validate exit codes/args/nil-injector |
| `d5949b8` | DRY consolidations + ValidationMode enum |
| `d50248c` | Consolidate TextMarshaler/Unmarshaler across 8 value types |
| `7ee4451` | Eliminate outputState split brain + version/long dual-write |
| `8f73ae1` | Wrap 19 bare fmt.Errorf calls with sentinel errors |
| `f9a1ae5` | Add 5 draconian validation mode tests |
| `67e2db4` | Fix BenchmarkScopeProvide panic on duplicate registration |
| `55b267c` | Update AGENTS.md and FEATURES.md for v2.3 |

### b) Key Changes

**Split brain eliminations:**
- Removed `outputState` wrapper type — `CLI[T].outputFormat` is the single source of truth
- Extracted `setVersion()`/`setLong()` internal methods — prevents `cli.version`/`cli.rootCmd.Version` drift

**DRY consolidations:**
- `output.go` 326→257 lines: `renderAndWrite`/`marshalAndWrite` generic helpers, registry closures
- `flow_context.go`: `branchWithCtx` helper — 5 Branch methods delegate to single helper
- `type_helpers.go`: `textMarshal[T]`/`textUnmarshal[T]` — eliminated ~48 lines across 8 value types
- `ValidationMode` enum: `Lenient/Strict/Draconian` constants with `>=` comparison

**Error quality:**
- 19 bare `fmt.Errorf` calls now wrap sentinel errors for `errors.Is()` chainability
- All files: `cli.go`, `cli_output.go`, `scope.go`, `config.go`, `flags.go`
- Every error path is now identifiable via sentinel

**Validation safety:**
- `NewExitError` validates 0-255 range, returns `(*ExitError, error)`
- `WithExactArgs/WithMinimumArgs/WithMaximumArgs`: panic on negative `n`
- `WithRangeArgs`: panic on negative min or min > max
- `NewScopeFromInjector`: returns `(*Scope, error)` on nil injector

### c) Test Coverage

- **220+ tests passing** (up from 199)
- **~82% coverage** (up from 80.9%)
- 5 new draconian validation tests covering:
  - Leaf without example fails in draconian mode
  - Leaf with example passes
  - Parent commands exempt from example requirement
  - Draconian enforces strict rules too
  - Strict mode does not require examples
- 22 benchmarks all passing (1 bug fix: ScopeProvide duplicate registration)

### d) Breaking Changes

| Before | After | Impact |
|--------|-------|--------|
| `NewExitError(code, err) *ExitError` | `NewExitError(code, err) (*ExitError, error)` | Callers must handle error return |
| `NewScopeFromInjector(inj, name) *Scope` | `NewScopeFromInjector(inj, name) (*Scope, error)` | Callers must handle error return |

### e) Current Metrics

| Metric | Value |
|--------|-------|
| Tests | 220+ passing |
| Coverage | ~82% |
| Lint issues | 0 |
| Race conditions | 0 |
| Benchmarks | 22 passing |
| Sentinel errors | 40+ |
| Commits ahead of origin | 11 |

### f) Remaining Work

| Priority | Task | Impact |
|----------|------|--------|
| Low | Strict mode: reject flags without `help` tags | Consistency |
| Low | Examples: `examples/superb/main.go` demonstrating all enforcement | Documentation |
| Low | TODO_LIST.md update | Housekeeping |

### g) Decisions Made

- **Skipped `labeledError` consolidation** — 4 lines saved not worth the abstraction (public types, distinct Error() formatting)
- **Kept `renderTableStyled` and `renderTableCSV` as-is** — genuinely different patterns (lipgloss table, streaming writer)
- **Panic for programmer errors** (negative args), error returns for runtime conditions (exit code range)
- **`ErrUnsupportedFormat` replaces `ParseOutputFormat` error** — sentinel-first wrapping in `cli_output.go`
