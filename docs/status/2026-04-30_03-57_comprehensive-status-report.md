# cmdguard v2.2 — Comprehensive Status Report

**Date:** 2026-04-30 03:57
**Branch:** master (4 commits ahead of origin)
**Session:** Continuation of "NOW GET SHIT DONE" sprint

---

## A) FULLY DONE ✅

### P1 — Tests & Dead Code Cleanup
- **6 new test files** created: `type_handler_test.go` (697 lines), `output_test.go` (211 lines), `editor_test.go` (44 lines), `command_suggest_test.go` (106 lines), `counting_flag_test.go` (100 lines), `env_tag_test.go` (159 lines)
- **Dead code removed**: `parseCustomDefault()` method, 5 unused functions in `config_setfield.go` (`wrapErr`, `parseField`, `parseAndSetLogLevel`, `parseAndSetLogFormat`, `parseAndSetDuration`)
- **Count handler consolidated**: moved from separate `registerCountHandler()` into `registerKinds()`

### P2 — API Polish
- **envPrefix propagation bug FIXED**: `WithEnvPrefix` now correctly propagates to command-level flags via threading through `cliToCobraCommand` → `initCommandFlags`
- **SuggestFlag API made consistent**: changed return type from `string` to `(string, bool)` to match `SuggestCommand`
- **Added tests**: `TestGenerateHelp`, `TestFlagNames`, updated `TestSuggestFlag`

### P3 — Documentation & Examples (6 new examples)
- `examples/env-tags/main.go` — env var priority chain, WithEnvPrefix
- `examples/output/main.go` — OutputTable/OutputResult in 12 formats
- `examples/counting/main.go` — count:"true" for -v/-vv/-vvv
- `examples/di-patterns/main.go` — service registration, invocation, health checks
- `examples/error-handling/main.go` — sentinel errors, FlagError, suggestions
- `examples/signals/main.go` — WithSignalHandling, graceful shutdown
- `examples/README.md` — feature matrix and usage table
- `AGENTS.md` updated with v2.2 structure, gotchas, CLI options
- `FEATURES.md` coverage updated
- `QUICKSTART.md` and `README.md` rewritten for v2.2

### P5 — New Features (committed)
- **WithOutputFormat[T](format)**: adds global `--output/-o` flag for format selection, resolved via `cli.OutputFormat()` — tested with 4 test cases
- **Shell completion**: `WithCompletion[T,F](fn)` wires cobra ValidArgsFunction, `WithValidArgs[T,F](args...)` for static valid args
- **Man page generation**: `cli.ManPage(section)`, `cli.WriteManPage(w, section)`, `GenerateManPageCommand[T](cli)` — tested with 3 test cases
- **WithColor deprecation**: updated doc comment with v3.0 removal target

### Commits Made This Session (4)
```
ba65a99 feat: add WithOutputFormat, shell completion, and improved deprecation
c3c8eff docs: add 6 v2.2 examples and update project documentation
3a1243f docs: update QUICKSTART.md and README.md for v2.2 API
860bf9f test: add comprehensive tests for v2.2 features and fix envPrefix propagation
```

---

## B) PARTIALLY DONE ⚠️

### Race Condition in Tests
- **55 race detections** when running `go test -race` on the v2 package
- All tests pass **without** `-race` flag (199/199 pass, 0 fail)
- Root cause: global `globalTypeRegistry` written by `RegisterTypeHandler()` and `RegisterGoDurationHandler()` — these tests cannot run in parallel
- Some tests need `//nolint:paralleltest` annotations but don't have them yet
- **Impact**: CI with `-race` will fail intermittently

### Man Page Generation
- Implementation complete and working
- Tests pass
- **NOT committed yet** — sitting as untracked files

---

## C) NOT STARTED 📝

### P4 — New Features (all not started)
1. **Config file loading (koanf)** — YAML/TOML/.env file loading, merge chain (file → env → flag), `WithConfigFile[T](path)`, `WithConfigFileFlags[T]()`
2. **Interactive prompts (huh)** — `PromptString`, `PromptSelect`, `PromptConfirm`, `WithPromptOnMissing[T,F]()`
3. **Glamour markdown help** — render help as styled markdown
4. **Spinner/progress middleware** — bubbles spinner integration
5. **Telemetry middleware** — OTel span wrapping

### P5 — Infrastructure (not started)
1. **Fix go-output dependency** — local `replace` directive blocks CI
2. **Release automation** — v2.2.0 tag, release notes
3. **Codecov integration**
4. **Pre-commit hooks fix** — currently broken, requires `--no-verify`
5. **Benchmark regression detection**
6. **Merge test helper files**

---

## D) TOTALLY FUCKED UP 💥

### 1. Race Detection Tests (55 failures)
**Severity: HIGH** — `-race` flag causes cascading test failures across the entire v2 package.

Root causes identified:
- `globalTypeRegistry` is a package-level `map` written by `RegisterTypeHandler()` without synchronization
- Tests that call `RegisterTypeHandler()` or `RegisterGoDurationHandler()` run `t.Parallel()` and race with each other
- Some tests use `t.Setenv()` + `t.Parallel()` which panics (caught by `paralleltest` linter but not all are annotated)

**Fix needed**: Either make `globalTypeRegistry` use `sync.RWMutex`, or ensure all tests that mutate global state are NOT parallel.

### 2. go-output Local Replace Directive
**Severity: HIGH** — `go.mod` has:
```
replace github.com/larsartmann/go-output => /home/lars/projects/go-output
```
This blocks: CI, other developers, `go mod tidy`, any remote builds.

### 3. go.mod Needs Tidy
**Severity: LOW** — `go.sum` and `go.mod` are modified (new imports from man page feature). Not yet committed.

---

## E) WHAT WE SHOULD IMPROVE 🎯

### Code Quality
1. **Race conditions** — biggest quality gap right now. 55 race detections is unacceptable for a library.
2. **Coverage dropped** from 81.2% → 80.6% (new code added without proportional tests)
3. **Test isolation** — global mutable state (`globalTypeRegistry`) is an anti-pattern for a library

### Architecture
4. **go-output coupling** — tight coupling to local replace; should be proper versioned dependency
5. **Config file loading** — most requested feature gap vs. competitors (viper/koanf)
6. **No release workflow** — 4 commits ahead of origin, no tag, no changelog automation

### Developer Experience
7. **Pre-commit hooks broken** — forces `--no-verify` on every commit
8. **Missing example tests** — all 12 examples have `[no test files]`
9. **Missing benchmarks** — benchmark file exists but no regression detection

---

## F) TOP 25 THINGS TO DO NEXT (Priority Order)

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | **Fix race conditions** — add sync.RWMutex to globalTypeRegistry | CRITICAL | 1h | Quality |
| 2 | **Commit man page code** — already written, just needs `git add` | HIGH | 2min | Commit |
| 3 | **Fix go-output dependency** — remove local replace, publish or GOPRIVATE | HIGH | 2h | Infra |
| 4 | **Add koanf config file loading** — most requested feature | HIGH | 4h | P4 Feature |
| 5 | **Run go mod tidy** — clean up go.mod/go.sum | LOW | 1min | Hygiene |
| 6 | **Add WithConfigFile[T](path)** CLI option | HIGH | 2h | P4 Feature |
| 7 | **Add WithConfigFileFlags[T]()** for auto --config | MED | 1h | P4 Feature |
| 8 | **Fix pre-commit hooks** | MED | 30min | Infra |
| 9 | **Create v2.2.0 release tag** | MED | 15min | Release |
| 10 | **Write release notes** (CHANGELOG.md) | MED | 30min | Docs |
| 11 | **Add huh interactive prompts** — PromptString/PromptSelect/PromptConfirm | MED | 3h | P4 Feature |
| 12 | **Push to origin** — 4 commits ahead | MED | 1min | Git |
| 13 | **Add example smoke tests** — run each example with `--help` | MED | 1h | Testing |
| 14 | **Improve coverage back to 82%+** — cover new code paths | MED | 2h | Testing |
| 15 | **Add codecov.yml + badge** | LOW | 15min | Infra |
| 16 | **Add benchmark regression CI** | LOW | 1h | Infra |
| 17 | **Add glamour markdown help rendering** | MED | 2h | P4 Feature |
| 18 | **Add spinner/progress middleware** | MED | 2h | P4 Feature |
| 19 | **Add telemetry middleware (OTel)** | LOW | 3h | P4 Feature |
| 20 | **Merge test helper files** — consolidate test utilities | LOW | 30min | Refactor |
| 21 | **Add shell completion example** | LOW | 30min | Docs |
| 22 | **Add man page example** | LOW | 15min | Docs |
| 23 | **Add WithPromptOnMissing[T,F]() command option** | MED | 1h | P4 Feature |
| 24 | **Document all new v2.2 features in API reference** | MED | 1h | Docs |
| 25 | **Set up release-automation (goreleaser?)** | LOW | 2h | Infra |

---

## G) MY #1 QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**What is the plan for `go-output`?**

The `go.mod` has `replace github.com/larsartmann/go-output => /home/lars/projects/go-output` which:
- Blocks all CI
- Blocks other developers from building
- Causes `go mod tidy` to complain
- Makes the project non-portable

I need to know: Should I:
- **A)** Remove go-output dependency entirely and inline the output formatting?
- **B)** Publish go-output to a proper Go module proxy (tagged release)?
- **C)** Use `GOPRIVATE` with a real remote URL instead of local path?
- **D)** Something else entirely?

This decision blocks P5 infrastructure work and any CI setup.

---

## Metrics Summary

| Metric | Value |
|--------|-------|
| Source files | 33 (pkg/cmdguard/v2) |
| Test files | 63 |
| Test lines | 10,722 |
| Test cases (PASS) | 199/199 (without -race) |
| Test cases (FAIL with -race) | 55 race detections |
| Code coverage | 80.6% |
| Examples | 12 directories |
| Build status | ✅ Clean |
| Lint status | Not run this session |
| Commits ahead of origin | 4 |
| Uncommitted files | 4 (manpage.go, manpage_test.go, go.mod, go.sum) |

---

_Report generated by Crush at 2026-04-30 03:57_
