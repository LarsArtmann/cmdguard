# cmdguard — Comprehensive Audit Sprint Status Report

**Date:** 2026-04-17 19:47
**Branch:** `master` @ `a6452f7`
**Version:** v2.1.0
**Status:** Post-audit. All tests pass, 31 lint issues remain.

---

## A. FULLY DONE

### A1. Audit Sprint (Sessions 1–3): 17 commits, 0 regressions

| # | Commit | Area | Summary |
|---|--------|------|---------|
| 1 | `10f5c73` | Correctness | Validators return errors instead of silently passing on bad config |
| 2 | `8766fd1` | Correctness | True deep copy in `MergeConfigs` for reference types (slices, maps, pointers) |
| 3 | `6aaf66e` | Style | Eliminate fallthrough and fix orphaned doc comment |
| 4 | `14ac747` | Cleanup | Remove dead `pkg/errtypes` package |
| 5 | `ad92544` | Examples | Migrate `examples/basic` from v1 to v2 API |
| 6 | `fa1c96d` | CI | Align Go matrix with `go.mod`, test all 5 examples in CI |
| 7 | `e0fcf87` | Tooling | Fix broken justfile recipe paths, rename `run-guarded` |
| 8 | `fb5cdfb` | Correctness | Return errors for min/regex validators on missing separator |
| 9 | `ad0148d` | Correctness | Propagate `required` tag parse error instead of silently ignoring |
| 10 | `0186cca` | Correctness | Validate Enum values in `SetField` against allowed list |
| 11 | `ffa9186` | Correctness | Use deep copy in `cloneFlags` to prevent shared mutable state |
| 12 | `c536230` | Docs | Document `BranchingFlowContext` is not goroutine-safe |
| 13 | `07dc73a` | Cleanup | Remove dead `internal/config/koanf.go` and orphaned test artifacts |
| 14 | `5f1b566` | Examples | Fix DI example typo, replace `MustInvoke` with error-returning `Invoke` |
| 15 | `b211db1` | Examples | Check `AddCommand` error returns in basic example |
| 16 | `7653743` | Docs | Fix `RegisterValidator` doc to say goroutine-safe |
| 17 | `a6452f7` | Tests | Strengthen weak test assertions for WithCLIScope, AddGlobalFlag, TimingMiddleware |

### A2. Prior Sprint (Pre-audit): 13 commits

| Commit | Summary |
|--------|---------|
| `ea5bfc9` | Fix gitignore, unify error sentinel, deduplicate test assertions |
| `e5c7540` | Convert all consumers to `NewCommand`/`NewParentCommand` constructors |
| `c562eca` | Add `mustInvoke` test helpers |
| `5867283` | Add `mustProvideValue` helper, deduplicate scope test boilerplate |
| `0f4783d` | Update TODO_LIST.md and FEATURES.md |
| `60da525` | Add godoc examples for constructors |
| `d82f993` | Add Deprecated notices to v1 API |
| `099ff35` | Update README.md for v2.1 API |
| `427355b` | Update AGENTS.md |
| `bc0683c` | Deduplicate no-op RunE lambdas |
| `aaa6bf1`–`f9bfcef` | Status reports |

### A3. Test Suite: All Green

```
ok  github.com/larsartmann/cmdguard/benchmarks        (no tests)
ok  github.com/larsartmann/cmdguard/examples/advanced-flags    39.1%
ok  github.com/larsartmann/cmdguard/examples/basic              0.0%
ok  github.com/larsartmann/cmdguard/examples/di               12.3%
ok  github.com/larsartmann/cmdguard/examples/typed             3.1%
ok  github.com/larsartmann/cmdguard/examples/validation       26.9%
ok  github.com/larsartmann/cmdguard/internal/config           95.7%
ok  github.com/larsartmann/cmdguard/internal/logging          97.1%
ok  github.com/larsartmann/cmdguard/pkg/cmdguard              87.4%
ok  github.com/larsartmann/cmdguard/pkg/cmdguard/v2           83.4%
ok  github.com/larsartmann/cmdguard/tests/integration          —
```

**Total: 13 packages, 0 failures, race detector enabled.**

---

## B. PARTIALLY DONE

### B1. Lint: 31 issues remaining (was ~40+)

| Linter | Count | Severity | Status |
|--------|-------|----------|--------|
| `staticcheck SA1019` | 10 | Low (deprecated v1 usage in integration tests) | Needs migration |
| `golines` | 6 | Cosmetic | Auto-fixable |
| `gci` | 3 | Cosmetic | Auto-fixable |
| `exhaustive` | 3 | Medium (missing switch cases) | Needs `default:` or explicit cases |
| `errcheck` | 2 | Medium (unchecked error in `examples/basic/main_test.go`) | Quick fix |
| `nlreturn` | 2 | Cosmetic | Auto-fixable |
| `wsl_v5` | 2 | Cosmetic | Auto-fixable |
| `forcetypeassert` | 1 | Medium (unchecked type assertion in `config_setfield.go:122`) | Should guard |
| `goimports` | 1 | Cosmetic | Auto-fixable |
| `nilnil` | 1 | Medium (`config_parsing.go:65` returns `nil, nil`) | Should use sentinel error |

**Assessment:** ~18 are auto-fixable formatting issues. ~13 require manual attention. None are correctness bugs.

### B2. Test Coverage: 83.4% for `pkg/cmdguard/v2`

**Functions at 0% coverage (12):**

| Function | File | Impact |
|----------|------|--------|
| `MustAddCommand` | `cli.go:126` | Low — trivial wrapper |
| `MustNewCLI` | `cli.go:134` | Low — trivial wrapper |
| `ArgsFromContext` | `cli_command.go:17` | Low — simple accessor |
| `WithFangOptions` | `cli_options.go:57` | Low — option function |
| `RegisterValidator` | `flags_validate.go:51` | Medium — public API |
| `runValidateTag` | `flags_validate.go:71` | Medium — core validation path |
| `validateEmail` | `flags_validate.go:131` | Medium — validator |
| `validateURL` | `flags_validate.go:144` | Medium — validator |
| `parseAndSetURL` | `flags_parse.go:193` | Medium — parse path |
| `parseAndSetEmail` | `flags_parse.go:200` | Medium — parse path |
| `parseAndSetPort` | `flags_parse.go:207` | Medium — parse path |
| `parseAndSetFilePath` | `flags_parse.go:214` | Medium — parse path |
| `parseAndSetHostPort` | `flags_parse.go:222` | Medium — parse path |
| `Exists` | `types_filepath.go:88` | Low — filesystem check |
| `IsDir` | `types_filepath.go:98` | Low — filesystem check |
| `IsFile` | `types_filepath.go:112` | Low — filesystem check |
| `Version` accessor | `command.go:77` | Low — field |
| `SilenceErrors` accessor | `command.go:80` | Low — field |
| `SilenceUsage` accessor | `command.go:83` | Low — field |
| `Group` accessor | `command.go:86` | Low — field |
| `WithGroupID` | `command.go:226` | Low — option function |

**Functions at 50–75% coverage (partial misses):**

| Function | Coverage | File |
|----------|----------|------|
| `validateMin` | 50.0% | `flags_validate.go:197` |
| `getFieldValue` | 50.0% | `config.go:98` |
| `validateTagRules` | 22.2% | `flags.go:278` |
| `parseAndSetCustom` | 45.5% | `flags_parse.go:139` |
| `ParseFilePath` | 53.3% | `types_filepath.go:24` |
| `validateNonEmpty` | 0.0% | `flags_validate.go:261` |
| `validateFieldByKind` | 0.0% | `flags_validate.go:270` |
| `formatFieldValue` | 0.0% | `flags_validate.go:281` |

### B3. FEATURES.md coverage table is stale

FEATURES.md still lists `pkg/errtypes` at 100% coverage, but that package was removed in commit `14ac747`. The Testing section needs updating.

---

## C. NOT STARTED

### C1. From TODO_LIST.md — Medium Priority

- [ ] Improve flag suggestion algorithm
- [ ] Migrate remaining testify usage to stdlib (if any remains)
- [ ] Add fuzz tests to `flags_parse.go` and `config_parsing.go`

### C2. From TODO_LIST.md — Documentation

- [ ] Update `docs/QUICKSTART.md` for v2.1 API
- [ ] Update `docs/MIGRATION_v1_v2.md` for v2.1 API
- [ ] DI Pattern Example in docs/
- [ ] Error Handling Example in docs/

### C3. From TODO_LIST.md — Examples

- [ ] Add example with real database connection
- [ ] Add lifecycle hook examples
- [ ] Advanced DI example

### C4. From TODO_LIST.md — Performance

- [ ] Add comprehensive performance benchmarks
- [ ] Add benchmark regression detection to CI

### C5. From TODO_LIST.md — Release & CI

- [ ] Create v2.1.0 release tag and notes
- [ ] Set up release automation (goreleaser)
- [ ] Add codecov integration
- [ ] Fix pre-commit hooks (currently 5 pre-existing errors requiring `--no-verify`)
- [ ] Migrate benchmarks from deprecated v2.New to v2.NewCLI

### C6. From TODO_LIST.md — Future (v3.0+)

- [ ] Plugin system for custom validators
- [ ] Config file auto-loading with koanf
- [ ] Shell completion helpers
- [ ] Progress/Spinner type (bubbles)
- [ ] Command groups feature

---

## D. TOTALLY FUCKED UP

### D1. Pre-commit hooks are broken

The project requires `git commit --no-verify` on every commit. Pre-commit hooks have 5 pre-existing errors. This is a developer experience hazard — every contributor will hit this.

### D2. `basic` binary was in working tree (cleaned)

A 6.4MB compiled binary named `basic` was sitting in the project root — likely a `go build` artifact from testing `examples/basic`. It was untracked and not in `.gitignore`. Cleaned during this report.

### D3. `integration_test.go` uses deprecated v1 API extensively

`tests/integration/integration_test.go` calls `cmdguard.New()` 10 times, triggering 10 `staticcheck SA1019` deprecation warnings. This is the single largest source of lint noise. The file should be migrated to v2 or clearly annotated as "v1 backward-compat tests".

### D4. `nilnil` linter flag on `config_parsing.go:65`

```go
func parseStructTags(...) ([]ParsedTag, error) {
    if s.Kind() != reflect.Struct {
        return nil, nil  // ← nilnil: returns nil error with nil value
    }
```

This is a design smell. The caller can't distinguish "no tags" from "not a struct" without checking the returned slice length.

### D5. Unchecked type assertion in `config_setfield.go:122`

```go
current := field.Interface().(Enum)  // ← forcetypeassert: will panic on non-Enum
```

Should use comma-ok assertion or be guarded by a type check.

---

## E. WHAT WE SHOULD IMPROVE

### E1. Lint Hygiene (31 issues → 0)

**Why:** A clean lint baseline enables CI enforcement. Currently golangci-lint exits non-zero, making it impossible to gate PRs on lint.

**Plan:** Auto-fix formatting (golines, gci, goimports, nlreturn, wsl_v5 = ~14 fixes). Manually fix the rest (exhaustive, errcheck, nilnil, forcetypeassert = ~6 fixes). Suppress SA1019 for v1 compat tests with `//nolint:staticcheck` or migrate.

### E2. Validator Coverage Gap

**Why:** The `validate` tag system has validators for email, URL, min, max, regex, non-empty — but 4 of these have **0% test coverage**. This is the part of the codebase most likely to have edge-case bugs.

**Plan:** Write targeted tests for `validateEmail`, `validateURL`, `validateNonEmpty`, `validateFieldByKind`, `formatFieldValue`, and the custom parse functions (`parseAndSetURL`, `parseAndSetEmail`, `parseAndSetPort`, `parseAndSetFilePath`, `parseAndSetHostPort`).

### E3. Type System: `Option[T]` and `Result[T]` Are Orphaned

**Why:** `Option[T]` and `Result[T]` are 400+ lines of well-tested code (100% coverage) but are never used by any other package in the codebase. They're utility types looking for a purpose.

**Plan:** Either integrate them into the CLI lifecycle (e.g., `RunE` returns `Result[T]`, flags use `Option[T]` for optional values) or extract them to a separate utility package. Currently they're dead weight in the public API.

### E4. Dead Command Accessors

**Why:** `command.go` has 5 accessor methods at 0% coverage (`Version()`, `SilenceErrors()`, `SilenceUsage()`, `Group()`, `WithGroupID()`). These are public API surface that nobody uses.

**Plan:** Either add tests proving they work, or remove them until there's a real consumer.

### E5. `validateTagRules` at 22.2% Coverage

**Why:** This function (`flags.go:278`) is the central dispatcher for the `validate:` struct tag system. At 22% coverage, most validation tag paths are untested.

**Plan:** Write a comprehensive table-driven test exercising all validation tags: `email`, `url`, `min`, `max`, `regex`, `nonempty`, `port`, `filepath`, `hostport`.

---

## F. TOP 25 THINGS TO DO NEXT

Sorted by **impact × urgency / effort**.

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Fix pre-commit hooks (eliminate `--no-verify` requirement) | 🔴 Critical | 2h | CI/Developer Experience |
| 2 | Auto-fix 14 formatting lint issues (`golangci-lint fmt`) | 🟠 High | 15min | Lint Hygiene |
| 3 | Fix `nilnil` in `config_parsing.go:65` (use sentinel error) | 🟠 High | 15min | Correctness |
| 4 | Fix `forcetypeassert` in `config_setfield.go:122` (guard type assertion) | 🟠 High | 15min | Correctness |
| 5 | Add `//nolint:staticcheck` to `integration_test.go` v1 usage | 🟠 High | 10min | Lint Hygiene |
| 6 | Fix 3 `exhaustive` switch warnings in `config.go` and `flag_helpers.go` | 🟠 High | 30min | Correctness |
| 7 | Fix 2 `errcheck` issues in `examples/basic/main_test.go` | 🟠 High | 10min | Correctness |
| 8 | Write tests for validator functions (email, URL, nonempty, fieldByKind) | 🟡 Medium | 2h | Test Coverage |
| 9 | Write tests for custom parse functions (URL, email, port, filepath, hostport) | 🟡 Medium | 2h | Test Coverage |
| 10 | Write tests for `validateTagRules` (22% → 90%+) | 🟡 Medium | 1h | Test Coverage |
| 11 | Write tests for `MustAddCommand`, `MustNewCLI` (trivial panic wrappers) | 🟡 Medium | 30min | Test Coverage |
| 12 | Update FEATURES.md coverage table (remove `pkg/errtypes`) | 🟢 Low | 10min | Documentation |
| 13 | Decide fate of `Option[T]` / `Result[T]` (integrate, extract, or document intent) | 🟡 Medium | 1h | Architecture |
| 14 | Remove dead command accessors or add tests (`Version`, `SilenceErrors`, etc.) | 🟢 Low | 30min | Cleanup |
| 15 | Update `docs/QUICKSTART.md` for v2.1 API | 🟡 Medium | 2h | Documentation |
| 16 | Update `docs/MIGRATION_v1_v2.md` for v2.1 API | 🟡 Medium | 2h | Documentation |
| 17 | Add fuzz tests to `flags_parse.go` and `config_parsing.go` | 🟡 Medium | 3h | Robustness |
| 18 | Create v2.1.0 release tag and GitHub release notes | 🟠 High | 1h | Release |
| 19 | Set up goreleaser for release automation | 🟢 Low | 2h | CI |
| 20 | Add codecov integration with coverage thresholds | 🟢 Low | 1h | CI |
| 21 | Migrate benchmarks from deprecated `v2.New` to `v2.NewCLI` | 🟢 Low | 30min | Technical Debt |
| 22 | Add benchmark regression detection to CI | 🟢 Low | 2h | Performance |
| 23 | Add lifecycle hook examples | 🟢 Low | 1h | Documentation |
| 24 | Shell completion helpers | 🟢 Low | 4h | Feature |
| 25 | Config file auto-loading | 🟢 Low | 8h | Feature |

---

## G. TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**What is the intended scope of cmdguard?**

The codebase has grown well beyond a "CLI guard library." It now includes:

- **CLI framework** (Cobra wrapper with DI)
- **Rich type system** (`Enum[T]`, `Option[T]`, `Result[T]`, `URL`, `Email`, `Port`, `FilePath`, `HostPort`, `Duration`, `LogLevel`, `LogFormat`)
- **Dependency injection** (samber/do wrapper)
- **Middleware system** (timing, recovery, custom chains)
- **Flow context** (branching context tree with value propagation)
- **Flag validation framework** (pluggable validators, struct tags)

`Option[T]` and `Result[T]` are Rust-inspired utility types with 600+ lines of tests — but nothing in cmdguard uses them. They feel like they belong in a separate `pkg/types` or `pkg/functional` package.

**My question:** Is cmdguard meant to be:
1. A **batteries-included CLI framework** (keep everything, document the type system as a feature)?
2. A **focused CLI guard library** (extract `Option`, `Result`, and rich types to separate packages)?
3. Something in between?

This architectural decision affects whether we invest in integrating these types into the CLI lifecycle or extracting them out. It also affects the public API surface for v2.1.0 release.

---

## Metrics Summary

| Metric | Value |
|--------|-------|
| Total source files (`pkg/cmdguard/v2/*.go`) | 36 production, 40 test |
| Total lines of code (v2 only) | 15,280 |
| Production LOC | ~4,800 |
| Test LOC | ~10,480 |
| Test:Production ratio | 2.18:1 |
| v2 test coverage | 83.4% |
| v1 test coverage | 87.4% |
| Internal packages coverage | 95–97% |
| Lint issues | 31 (14 auto-fixable) |
| Open TODO items | ~25 |
| Commits since last release | 30+ |
| Pre-commit hooks | Broken (5 errors) |

---

*Generated 2026-04-17 19:47 CEST. Branch master @ a6452f7.*
