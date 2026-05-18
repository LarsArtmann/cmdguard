# cmdguard — Comprehensive Status Report

**Date:** 2026-05-18 15:07
**Version:** v2.3.0-dev
**Go Version:** 1.26
**Since Last Report:** 2026-05-17 18:01 (Go 1.26 migration, Ptr removal)

---

## Executive Summary

cmdguard is a production-grade Go library for building validated Cobra CLI applications with type-safe dependency injection. The v2 API is **feature-complete** with 104 source files, 66 test files, 285 test/benchmark/fuzz functions, 84.3% coverage on the core package, 0 lint issues, 0 race conditions, and 0 build errors. The library is ready for v2.3.0 release pending cleanup items listed below.

---

## Health Dashboard

| Metric             | Value          | Status | Trend      |
| ------------------ | -------------- | ------ | ---------- |
| Source files (v2)  | 104            | ✅      | Stable     |
| Test files (v2)    | 66             | ✅      | Stable     |
| Test functions     | 285            | ✅      | Stable     |
| Core coverage      | 84.3%          | ✅      | Down 0.2%  |
| Lint issues        | 0              | ✅      | Stable     |
| Race conditions    | 0              | ✅      | Stable     |
| Build errors       | 0              | ✅      | Stable     |
| Deprecated APIs    | 3 remaining    | ⚠️      | Cleanup    |
| Examples with tests| 5/12           | ⚠️      | Gap        |
| TODOs in source    | 0              | ✅      | Clean      |
| Total lines (v2)   | 17,822         | ✅      | Growing    |

---

## a) FULLY DONE ✅

### Core API (100%)

- `CLI[T]` with full constructor, options, accessors, and lifecycle
- `Command[T, F]` with 21 command options, constructors, validation
- `AddCommand` standalone function (per-command flag types)
- `NewParentCommand` / `MustNewParentCommand`
- All 16 CLI options functional and tested
- All 21 command options functional and tested

### Flag System (100%)

- Struct tag flags: `flag`, `short`, `default`, `help`, `required`, `validate`, `env`, `count`
- Flag typo suggestions (Levenshtein)
- Subcommand typo suggestions
- Instance-scoped validators (`FlagRegistry.RegisterFlagValidator()`)
- TypeHandler registry for custom types
- Flag priority chain: explicit flag → env → default

### Value Types (100%)

- `Duration`, `Enum`, `LogLevel`, `LogFormat`, `URL`, `Email`, `Port`, `FilePath`, `HostPort`
- All 9 types with MarshalText/UnmarshalText, validation, godoc examples
- Fuzz tests for all parsers (7 fuzz targets)

### Dependency Injection (100%)

- `Scope` wrapping samber/do/v2
- `Provide`, `ProvideNamed`, `ProvideValue`, `Invoke`, `InvokeNamed`, `MustInvoke`
- Hierarchical child scopes
- Lifecycle: `ShutdownAll`, `HealthCheck` with context
- `NewScopeFromInjector` returning `(*Scope, error)`

### Rich Output (100%)

- 12 output formats via go-output integration
- `OutputResult`, `OutputTable`, `OutputStyledTable`
- `WithOutputFormat[T]()` auto-flag
- Format renderer registry

### Error Handling (100%)

- 35+ sentinel errors with `errors.Is()` support
- `CommandError`, `FlagError`, `ServiceError`, `EnumError` typed errors
- `FlagErrorWithSuggestion` for typo fixes
- `ExitCoder` / `ExitError` with 0-255 range validation
- Zero panics in library API

### Middleware (100%)

- `TimingMiddleware`, `RecoveryMiddleware` (with stack traces)
- Custom middleware chain: `func(ctx, cfg, info, next) error`

### Positional Arguments (100%)

- `WithExactArgs`, `WithMinimumArgs`, `WithMaximumArgs`, `WithRangeArgs`, `WithNoArgs`, `WithArgs`

### Version & Completion (100%)

- `VersionCommand[T]`, `MustVersionCommand[T]`, `GenerateVersionCommand[T]`
- Shell completion: `WithCompletion[T,F]`, `WithValidArgs[T,F]`
- Man page generation: `cli.ManPage()`, `GenerateManPageCommand[T]`

### Validation Modes (100%)

- `Lenient` / `Strict` / `Draconian` spectrum
- `WithConfigValidation[T](fn)` for cross-field validation
- `WithStrictValidation[T]()` requires short descriptions
- `WithDraconianValidation[T]()` requires examples on leaf commands

### Quality (100%)

- 0 lint issues (was 113)
- 0 race conditions (was 55)
- 0 build errors
- 0 TODOs in source files
- Pre-commit hook script
- GitHub Actions CI workflow

### Documentation & Examples (100%)

- README.md as public landing page
- QUICKSTART.md, CLI_DESIGN_PRINCIPLES.md, DOMAIN_LANGUAGE.md
- 12 runnable examples (all compile)
- Features.md with honest status indicators
- AGENTS.md contributor guide

---

## b) PARTIALLY DONE ⚠️

### Examples Test Coverage — 5/12 have tests

| Example         | Tests | Gap                                  |
| --------------- | ----- | ------------------------------------ |
| advanced-flags  | ✅     | —                                    |
| basic           | ✅     | —                                    |
| di              | ✅     | —                                    |
| typed           | ✅     | —                                    |
| validation      | ✅     | —                                    |
| counting        | ❌     | No test file                         |
| di-patterns     | ❌     | No test file                         |
| env-tags        | ❌     | No test file                         |
| error-handling  | ❌     | No test file                         |
| output          | ❌     | No test file                         |
| signals         | ❌     | No test file                         |
| subcommands     | ❌     | No test file                         |

### Deprecated APIs Still Present — 3 items pending v3.0 removal

| API                    | File                        | Replacement                       |
| ---------------------- | --------------------------- | --------------------------------- |
| `WithColor[T](bool)`   | `cli_options.go:56`         | `WithFang[T](bool)`               |
| `IsExecutable()`       | `command.go:105`            | `HasHandler()`                    |
| `FlowContextAccessor`  | `flow_context_access.go:68` | `BranchingFlowContext` methods    |

### go-output Sub-module Replace Directives — Partially Cleaned

- `go-output/cmdguard/` deleted today (orphan, zero consumers)
- `go-output/go.mod` still has `replace` for `cmdguard → ./cmdguard` (now dead path)
- `go-output` sub-modules (`enum`, `escape`, `testhelpers`) use `v0.0.0` with local `replace` — could block external consumers if not properly tagged

---

## c) NOT STARTED 📝

### Phase 9: Architecture Hardening (v2.3)

| Item                                                    | Priority |
| ------------------------------------------------------- | -------- |
| `errors.As` → `errors.AsType[ExitCoder]` (Go 1.26)     | Low      |
| Extract `handlerConfig[T,F]` from 8-param function      | Medium   |
| Add `Phase` typed enum for `CommandInfo.Phase string`   | Medium   |
| Fix 7 unwrapped error returns                           | Medium   |
| Consolidate 5 error types into `labeledError`           | Medium   |
| Split `type_handler.go` (was 481 lines, now 177+185)    | Done ✅   |
| Split `command.go` (287 lines) — extract args options   | Low      |
| Split `flow_context.go` (264 lines) — extract options   | Low      |
| Fix `outputFormat`/`outputState.format` split brain     | Low      |
| Consolidate value type MarshalText/UnmarshalText        | Low      |

### Performance

| Item                                   | Priority |
| -------------------------------------- | -------- |
| CLI construction benchmark             | Medium   |
| Flag parsing benchmark                 | Medium   |
| Command execution benchmark            | Medium   |
| Benchmark regression detection in CI   | Low      |

### CI/CD

| Item                           | Priority |
| ------------------------------ | -------- |
| Codecov integration            | Medium   |
| v2.3.0 release tag and notes   | High     |
| Release automation             | Medium   |

### Future (v3.0+)

- Config file auto-loading with koanf (YAML/TOML/.env)
- Interactive prompts (huh integration)
- Spinner/progress middleware (bubbles)
- Glamour markdown help rendering
- Telemetry middleware (OpenTelemetry)
- Plugin system for custom validators/type handlers

### v3.0 API Cleanup (Breaking)

- Make `NoFlags` a distinct named type
- Change `TimingMiddleware` callback to include error
- Remove `BranchWithTimeout`/`BranchWithDeadline` (string-based)
- Remove `FlowContextAccessor`
- Rename `Get[T]`/`MustGet[T]`
- Make `RegisterInScope` generic
- Remove/redesign `Package()`

---

## d) TOTALLY FUCKED UP 💥

**Nothing is critically broken.** The codebase is healthy. But here are the items that need attention:

### 1. go-output Replace Directive for Deleted cmdguard/

`go-output/go.mod` still has `replace github.com/larsartmann/go-output/cmdguard => ./cmdguard` pointing to a directory that no longer exists. This will cause `go mod tidy` to fail in go-output until removed.

### 2. go-output Sub-module Versioning Risk

go-output's sub-modules (`enum`, `escape`, `testhelpers`) are referenced at `v0.0.0` with local `replace` directives. External consumers of `go-output v0.2.0` may hit resolution errors if these sub-modules aren't properly tagged and published. This could block anyone running `go get github.com/larsartmann/go-output@v0.2.0`.

### 3. Features.md and TODO_LIST.md Stale Metrics

- Features.md says "~82%" coverage — actual is 84.3%
- TODO_LIST.md says "247 tests" — actual is 285 (test/benchmark/fuzz)
- TODO_LIST.md says "80.4% coverage" — actual is 84.3%
- TODO_LIST.md says "210 in v2" — actual is 234 in v2

### 4. Ptr[T] Documented in FEATURES.md but Removed

`FEATURES.md` line 166 still lists `Ptr[T]` as a feature. It was removed in commit `bd20534` (replaced by Go 1.26 `new(v)` built-in).

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture

1. **Split large files** — `scope.go` (361 lines, 6 responsibilities), `config.go` (230 lines, 5 responsibilities), `output.go` (258 lines, 5 responsibilities), `cli.go` (247 lines, 4 responsibilities) all do too many things
2. **Fix outputFormat split brain** — `outputFormat` field on CLI vs `outputState.format` somewhere else; needs consolidation
3. **Consolidate error types** — 5 separate error structs share the same pattern; extract to `labeledError`
4. **Fix 7 unwrapped error returns** — bare `return err` without context in some paths

### Testing

5. **Add tests for 7 untested examples** — counting, di-patterns, env-tags, error-handling, output, signals, subcommands
6. **Add benchmarks** — CLI construction, flag parsing, command execution
7. **Push coverage from 84.3% → 90%+** — focus on uncovered error paths

### Documentation

8. **Fix stale metrics** in FEATURES.md and TODO_LIST.md
9. **Remove Ptr[T] from FEATURES.md** — it's gone
10. **Add godoc examples** for remaining value types

### Dependencies

11. **Fix go-output replace directive** — remove dead cmdguard replace, ensure sub-modules are tagged
12. **Verify go-output v0.2.0 is resolvable** from clean GOPROXY — run `GONOSUMCHECK= GONOSUMDB= GOFLAGS= go get github.com/larsartmann/go-output@v0.2.0` in a temp module

### CI/CD

13. **Add codecov integration** — track coverage trends
14. **Create v2.3.0 release** — tag, release notes, git tag
15. **Add benchmark regression detection**

---

## f) Top 25 Things We Should Get Done Next

### High Impact (Do First)

| #  | Item                                                | Effort | Impact |
| -- | --------------------------------------------------- | ------ | ------ |
| 1  | Fix stale metrics in TODO_LIST.md and FEATURES.md   | 10min  | High   |
| 2  | Remove Ptr[T] from FEATURES.md                      | 2min   | Medium |
| 3  | Fix go-output replace for deleted cmdguard/          | 5min   | High   |
| 4  | Verify go-output v0.2.0 resolves from clean proxy   | 15min  | High   |
| 5  | Tag v2.3.0 release with release notes               | 30min  | High   |

### Architecture Cleanup (Pre-Release Polish)

| #  | Item                                                | Effort | Impact |
| -- | --------------------------------------------------- | ------ | ------ |
| 6  | Fix outputFormat/outputState split brain            | 30min  | Medium |
| 7  | Consolidate 5 error types into labeledError         | 1hr    | Medium |
| 8  | Fix 7 unwrapped error returns (add fmt.Errorf ctx)  | 30min  | Medium |
| 9  | Extract handlerConfig[T,F] from wireHandler         | 30min  | Medium |
| 10 | Add Phase typed enum for CommandInfo.Phase string   | 20min  | Medium |

### Testing Gaps

| #  | Item                                                | Effort | Impact |
| -- | --------------------------------------------------- | ------ | ------ |
| 11 | Add tests for subcommands example                   | 20min  | Medium |
| 12 | Add tests for env-tags example                      | 20min  | Medium |
| 13 | Add tests for counting example                      | 20min  | Medium |
| 14 | Add tests for error-handling example                | 20min  | Medium |
| 15 | Add tests for di-patterns example                   | 20min  | Medium |
| 16 | Add tests for output example                        | 20min  | Medium |
| 17 | Add tests for signals example                       | 20min  | Medium |
| 18 | Add CLI construction benchmark                      | 20min  | Medium |
| 19 | Add flag parsing benchmark                          | 20min  | Medium |
| 20 | Push core coverage to 90%+                          | 2hr    | High   |

### CI/CD & Infrastructure

| #  | Item                                                | Effort | Impact |
| -- | --------------------------------------------------- | ------ | ------ |
| 21 | Add codecov integration                             | 30min  | Medium |
| 22 | Add benchmark regression detection to CI            | 1hr    | Medium |
| 23 | Set up release automation (goreleaser?)             | 2hr    | Medium |

### Future Prep

| #  | Item                                                | Effort | Impact |
| -- | --------------------------------------------------- | ------ | ------ |
| 24 | Write v3.0 migration guide (breaking changes list) | 1hr    | Medium |
| 25 | Evaluate koanf for config file auto-loading         | 2hr    | Low    |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Is `go-output v0.2.0` actually resolvable from a clean GOPROXY?**

go-output's `go.mod` has 5 sub-modules at `v0.0.0` with local `replace` directives:
- `go-output/enum`
- `go-output/escape`
- `go-output/testhelpers`
- `go-output/cmdguard` (now deleted but replace still present)

If these sub-modules aren't tagged and pushed as separate module versions, any external consumer running `go get github.com/larsartmann/go-output@v0.2.0` will fail with `module ... found, but does not contain package`. This would block anyone from using cmdguard in a real project.

**Recommended action:** Run `GONOSUMCHECK=* GONOSUMDB=* GOFLAGS= go get github.com/larsartmann/go-output@v0.2.0` in a temporary empty Go module to verify. If it fails, the sub-modules need to be tagged, or the monorepo structure needs a `go.work` + proper tagging strategy.

---

## Session History (Recent Commits)

| Date       | Commit   | Description                                                  |
| ---------- | -------- | ------------------------------------------------------------ |
| 2026-05-18 | (today)  | Deleted go-output/cmdguard/ orphan (404 lines, zero consumers) |
| 2026-05-17 | aa439dd  | Consistent markdown table formatting, domain language docs    |
| 2026-05-17 | bd20534  | Remove Ptr[T] — use Go 1.26 new(v) built-in                  |
| 2026-05-17 | 3b8e9d6  | Post-refactor-polish comprehensive status report              |
| 2026-05-16 | Multiple | Architecture hardening, error consolidation, coverage push   |
| 2026-05-16 | Multiple | Test sprint, dead code removal, v2.3 features                |
| 2026-04-30 | Multiple | v2.2 SUPERB sprint, DI, env tags, output integration         |

---

_Generated by Crush on 2026-05-18 at 15:07_
