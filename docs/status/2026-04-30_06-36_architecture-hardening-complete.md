# cmdguard — Full Status Report

**Date:** 2026-04-30 06:36 CEST
**Branch:** master (up to date with origin)
**Version:** v2.2.0
**Go:** 1.26

---

## Executive Summary

cmdguard is a Go library for building validated Cobra CLI applications with type-safe dependency injection. After three intensive sessions, the project is in **excellent shape**: 0 lint issues, 0 race conditions, 199 tests passing, 80.9% coverage, clean build, published dependencies. The codebase underwent significant architecture hardening this session — fixing real bugs, eliminating duplication, and improving type safety.

---

## Current Metrics

| Metric | Value | Status |
|--------|-------|--------|
| Build | `go build ./...` — 0 errors | ✅ Clean |
| Tests | 199 passing with `-race` | ✅ Clean |
| Lint | `golangci-lint run ./...` — 0 issues | ✅ Clean |
| Vet | `go vet ./...` — 0 issues | ✅ Clean |
| Coverage | 80.9% (`pkg/cmdguard/v2`) | ✅ Good |
| Race conditions | 0 | ✅ Clean |
| Source files | 96 in `pkg/cmdguard/v2/` (33 production, 63 test) | |
| Total lines | 16,323 in v2 package | |
| Examples | 12 directories | |
| Dependencies | All published (no local replace) | ✅ Clean |
| CI | GitHub Actions (build+test+lint) | ✅ Active |
| Uncommitted | Auto-formatting drift only (golines, go.mod tidy) | ⚠️ Trivial |

---

## a) FULLY DONE

### Bugs Fixed (This Session — 12 commits pushed)

| # | Bug | File | Severity |
|---|-----|------|----------|
| 1 | **BranchingFlowContext double-cancellation** — `Cancel()` called both `b.cancels[i]()` and `child.Cancel()`, invoking the same cancel func twice | `flow_context.go` | 🔴 Real bug |
| 2 | **Enum.Allowed() returns internal slice** — callers could mutate the enum's allowed values | `types_enum.go` | 🟡 Defensive |
| 3 | **RecoveryMiddleware loses stack trace** — only captured panic value, not the stack | `middleware.go` | 🟡 Debug quality |

### Architecture Improvements (This Session)

| # | Change | Files | Impact |
|---|--------|-------|--------|
| 4 | Use `errors.Join` in `Scope.ShutdownAll` | `scope.go` | Proper error chains for `errors.Is`/`errors.As` |
| 5 | Fix `Scope.Path()` allocation | `scope.go` | Collect-then-reverse instead of prepend-per-iteration |
| 6 | Extract shared `lookupFlagInCommand` | `flags.go`, `flags_parse.go` | Eliminated duplicated local→persistent lookup |
| 7 | `getFieldValue` uses `fmt.Stringer` | `config.go` | No more hardcoded type switch per custom type |
| 8 | Table-driven type handler registration | `type_handler.go` | **-52 lines**, 5 identical handlers → loop |
| 9 | `makeEnumLikeHandler` for LogLevel/LogFormat | `type_handler.go` | Deduplicated enum-like help formatting |
| 10 | `parseAndSetValue` delegates to `SetField` | `flags_parse.go` | **-20 lines**, removed duplicated reflect logic |
| 11 | `map[string]struct{}` for command set | `cli.go` | Idiomatic Go |
| 12 | `SetConfig` WARNING documentation | `cli_accessors.go` | Footgun documented |

### New Features

| # | Feature | Files |
|---|---------|-------|
| 13 | NewParentCommand example (`examples/subcommands/`) | `main.go` |
| 14 | Shareable pre-commit hook (`scripts/pre-commit`) | `pre-commit` |

### Prior Sessions (Also Done)

- Unified type dispatch into TypeHandler registry
- `env:"VAR"` struct tag support with `WithEnvPrefix`
- Subcommand typo suggestions
- `WithSignalHandling[T]()` for SIGINT/SIGTERM
- go-output integration (12 output formats)
- `count:"true"` struct tag for counting flags
- `EditInEditor()` with `context.Context`
- Shell completion wiring
- Man page generation via mango-cobra
- Fix all 55 race conditions (`sync.RWMutex`)
- Remove local go-output replace (tagged v0.1.0)
- Achieve 0 lint issues (was 113)
- Output.go format renderer registry
- Sentinel errors `ErrUnsupportedFormat`, `ErrFormatRequiresTypedData`
- 6 working examples (env-tags, output, counting, di-patterns, error-handling, signals)
- Fuzz tests for value type parsers
- GitHub Actions CI workflow

---

## b) PARTIALLY DONE

| Item | Status | What's Left |
|------|--------|-------------|
| Pre-commit hooks | Script exists at `scripts/pre-commit` but must be manually copied to `.git/hooks/`. No `lefthook`/`husky`/`pre-commit` framework integration. | Auto-install mechanism |
| Benchmarks | 18 benchmark functions exist in `benchmarks/guard_bench_test.go` but TODO_LIST says "add CLI construction/flag parsing/command execution benchmarks" — unclear if existing ones cover these. | Verify coverage, add if missing |
| CI codecov | Not integrated | Add `codecov` step to GitHub Actions |

---

## c) NOT STARTED

| # | Item | Effort | Impact |
|---|------|--------|--------|
| 1 | Codecov integration in CI | Small | Medium |
| 2 | v2.2.0 release tag and notes | Small | High |
| 3 | Release automation (goreleaser?) | Medium | Medium |
| 4 | Benchmark regression detection in CI | Medium | Medium |
| 5 | Config file auto-loading with koanf | Large | High |
| 6 | Interactive prompts (huh integration) | Large | Medium |
| 7 | Spinner/progress middleware (bubbles) | Medium | Low |
| 8 | Glamour markdown help rendering | Medium | Medium |
| 9 | Telemetry middleware (OpenTelemetry) | Medium | Low |
| 10 | Plugin system for custom validators | Large | Medium |

---

## d) TOTALLY FUCKED UP — Nothing!

No regressions, no broken tests, no unfixable issues. The codebase is in its best state ever.

---

## e) WHAT WE SHOULD IMPROVE

### High Priority (Architecture Debt)

1. **`RegisterInScope` accepts `...any`** — loses all type safety. Uses runtime type switching. Should be generic or provide typed variants.

2. **`Package()` panics on error** — undermines the error-safe DI pattern. The `do.Package` contract forces a void function, but panicking in library code is a smell. Consider a different integration pattern.

3. **`BranchWithTimeout`/`BranchWithDeadline` accept `string` params** — parsing strings at runtime when the caller has `time.Duration`/`time.Time` is a code smell. Should accept typed parameters.

4. **`FlowContextAccessor` is a thin wrapper** — has 3 methods that all delegate to `BranchingFlowContext`. No clear benefit over using `BranchingFlowContext` directly. Adds API surface for no reason.

5. **`Get[T]`/`MustGet[T]` naming** — `Get` is extremely generic and will collide with other packages. Should be `GetFlowValue[T]` or similar.

6. **`IsExecutable` is just `HasHandler`** — redundant method, same logic, adds API surface.

7. **`NoFlags` is a type alias** (`type NoFlags = struct{}`) — users who accidentally pass an empty struct literal get `NoFlags` behavior silently. A distinct `type NoFlags struct{}` would be clearer.

### Medium Priority (Code Quality)

8. **`validateMin`/`validateMax` and `validateMinLen`/`validateMaxLen` near-duplicates** in `flags_validate.go` — could collapse into parameterized validators.

9. **Two separate validation execution paths** in `flags_validate.go` — `runValidateTag` and `parseValidateRulesWithRegistry` are different entry points for the same concept. Risk of divergence.

10. **`validateRegex` compiles regex on every call** — should cache or compile once.

11. **`formatFieldValue` in `flags_validate.go`** — many branches produce identical output. Could collapse with a default.

12. **`derefPointerToStruct` is exported** from `config_parsing.go` — internal utility that leaks implementation detail.

### Low Priority (Polish)

13. **22 one-liner accessor methods** on `Command` (~60 lines of boilerplate) — could use code generation.

14. **`wireSubcommandSuggestions`** in `cli_command.go` sets `SetFlagErrorFunc` on root but all handlers just `return err` unchanged — effectively a no-op.

15. **`outputState` mutex in `cli_output.go`** — CLI execution is single-threaded. Unnecessary complexity.

16. **Hardcoded format help string** in `cli_output.go` — duplicates the valid format list from `ParseOutputFormat`.

### v3.0 API-Breaking Candidates

17. **`TimingMiddleware` callback should include error** — so middleware can distinguish success vs failure timing.

18. **`ExecuteAndExit` always exits with code 1** — no way to propagate structured exit codes.

19. **`WithColor` deprecated** — still exported, should be removed in v3.

---

## f) Top #25 Things to Do Next

Sorted by **impact × effort** (highest first):

| # | Task | Effort | Impact | Type |
|---|------|--------|--------|------|
| 1 | **Tag v2.2.0 release** with release notes | 30min | 🔴 High | Release |
| 2 | **Commit the auto-formatting drift** (golines, go.mod tidy) | 5min | 🟡 Medium | Housekeeping |
| 3 | **Add `.gitignore` entry** for compiled binaries (env-tags, subcommands) | 2min | 🟡 Medium | Housekeeping |
| 4 | **Verify benchmark coverage** — do existing benchmarks cover CLI construction, flag parsing, command execution? | 30min | 🟡 Medium | Testing |
| 5 | **Add codecov to CI** | 30min | 🟡 Medium | CI |
| 6 | **Set up goreleaser** for release automation | 2hr | 🟡 Medium | CI |
| 7 | **Make `RegisterInScope` generic** instead of `...any` | 1hr | 🟡 Medium | Architecture |
| 8 | **Change `BranchWithTimeout`/`BranchWithDeadline`** to accept `time.Duration`/`time.Time` | 30min | 🟡 Medium | API |
| 9 | **Rename `Get[T]`/`MustGet[T]`** to `GetFlowValue[T]`/`MustGetFlowValue[T]` | 15min | 🟢 Low | API |
| 10 | **Remove `FlowContextAccessor`** — use `BranchingFlowContext` directly | 30min | 🟢 Low | Cleanup |
| 11 | **Remove `IsExecutable`** method from Command | 5min | 🟢 Low | API |
| 12 | **Collapse `validateMin`/`validateMax`** into parameterized validators | 1hr | 🟢 Low | Dedup |
| 13 | **Cache regex compilation** in `validateRegex` | 15min | 🟢 Low | Perf |
| 14 | **Unexport `derefPointerToStruct`** | 5min | 🟢 Low | Encapsulation |
| 15 | **Add `.gitignore`** for example binaries | 2min | 🟢 Low | Housekeeping |
| 16 | **Remove `wireSubcommandSuggestions`** no-op | 15min | 🟢 Low | Dead code |
| 17 | **Remove `outputState` mutex** — document single-threaded assumption | 15min | 🟢 Low | Simplicity |
| 18 | **Generate help string from valid format list** in cli_output.go | 30min | 🟢 Low | DRY |
| 19 | **Make `NoFlags` a distinct named type** | 30min | 🟢 Low | API (breaking) |
| 20 | **Add `TimingMiddleware` error parameter** | 30min | 🟢 Low | API (breaking) |
| 21 | **Config file auto-loading with koanf** | 1-2 days | 🔴 High | Feature |
| 22 | **Interactive prompts (huh)** | 1-2 days | 🟡 Medium | Feature |
| 23 | **Glamour markdown help rendering** | 4hr | 🟡 Medium | Feature |
| 24 | **Telemetry middleware (OpenTelemetry)** | 4hr | 🟢 Low | Feature |
| 25 | **Plugin system for custom validators** | 1-2 days | 🟡 Medium | Feature |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should `NoFlags` remain a type alias or become a distinct type?**

Currently `NoFlags = struct{}` (alias). Changing to `NoFlags struct{}` (distinct type) would:
- ✅ Prevent accidental `struct{}` from matching `NoFlags` behavior
- ✅ Allow adding methods to `NoFlags` later (e.g., `String()`)
- ❌ Break any code comparing `Command[T, struct{}]` with `Command[T, NoFlags]`
- ❌ Users who wrote `v2.Command[Config, struct{}]` would need to change to `v2.Command[Config, v2.NoFlags]`

This is an API-breaking change. **Is the v2 API considered stable enough that we should defer this to v3.0, or is it still acceptable to break?**

---

## Session History

### Session 1 (prior)
- 113 lint → 0 lint
- 55 race conditions → 0
- Output.go registry refactor
- go-output v0.1.0 published

### Session 2 (prior)
- v2.2.0 features complete
- 6 new examples
- CI workflow
- Documentation updated

### Session 3 (this session — 12 commits)
- Fixed 3 bugs (double-cancel, Enum mutation, lost stack traces)
- 9 architecture improvements (deduplication, delegation, type safety)
- 2 new features (subcommands example, pre-commit hook)
- Documentation updated

### All Commits (24 total since work began)

```
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
c3c8eff docs: add 6 v2.2 examples and update project documentation
```

---

## Uncommitted Changes (Trivial)

- `go.mod`/`go.sum` — dependency version drift from `go mod tidy` (indirect charmbracelet deps)
- `middleware.go` — golines auto-formatting (long line split)
- `type_handler.go` — golines auto-formatting (long line split)
- Stale binaries: `env-tags`, `subcommands` (compiled examples, should be gitignored)

---

_Generated by Crush — 2026-04-30_
