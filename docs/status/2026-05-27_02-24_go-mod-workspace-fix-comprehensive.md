# Comprehensive Status Report — go.mod Workspace Fix & v2.3.0-dev Health Check

**Date:** 2026-05-27 02:24 CEST  
**Branch:** master  
**Commits since last report:** 1 (workspace fix)  
**Working tree:** Modified (go.mod, go.sum)  
**Tests:** 932 individual test runs, 84.0% coverage  
**Build:** PASS  
**Lint:** 0 issues  
**Race detector:** 0 races

---

## a) FULLY DONE

### This Session (Immediate Fix)

- **Fixed `go mod tidy` failure** — `go mod tidy` was failing with:
  ```
  github.com/larsartmann/go-output/testhelpers/graphtest: reading .../go.mod at revision testhelpers/graphtest/v0.0.0: unknown revision
  ```
  Root cause: `go-output/serialization` has a test dependency on `go-output/testhelpers/graphtest`, which only exists as a local workspace module (not published to GitHub with a tag). `go mod tidy` tried to resolve it from the remote and failed.

- **Added 11 `replace` directives to `go.mod`** — All `github.com/larsartmann/go-output/*` modules redirected to local relative paths (`../go-output/...`):
  - `go-output` → `../go-output`
  - `go-output/d2` → `../go-output/d2`
  - `go-output/delimited` → `../go-output/delimited`
  - `go-output/enum` → `../go-output/enum`
  - `go-output/escape` → `../go-output/escape`
  - `go-output/graph` → `../go-output/graph`
  - `go-output/markup` → `../go-output/markup`
  - `go-output/serialization` → `../go-output/serialization`
  - `go-output/table` → `../go-output/table`
  - `go-output/testhelpers` → `../go-output/testhelpers`
  - `go-output/testhelpers/graphtest` → `../go-output/testhelpers/graphtest`

- **Verified build integrity** — `go build ./...` compiles cleanly
- **Verified all tests pass** — `go test ./... -count=1 -timeout 120s -race`: all packages green
- **Verified lint** — `golangci-lint run ./...`: 0 issues
- **Ran `go mod tidy`** — Dependencies cleaned: removed unused `charmbracelet/x/exp/golden`, promoted `gopkg.in/yaml.v3` and `pelletier/go-toml/v2` to direct (already used by configload), added `kr/pretty` and `rogpeppe/go-internal` indirects

### Previous Sessions (Still True)

- Config file loading feature (JSON/YAML/TOML, `WithConfigFile`, `WithConfigFileLoader`)
- 14 working examples
- 84.0% test coverage in `pkg/cmdguard/v2`
- 0 lint issues
- 0 race conditions
- Man page generation, shell completion, positional args validators, exit codes, strict/draconian validation
- Full sentinel error coverage (40+ errors)

---

## b) PARTIALLY DONE

### Config File Feature

- **Core:** Fully implemented, tested, and documented
- **Configload sub-package:** 0% coverage — YAML/TOML loaders have no tests (thin wrappers around well-tested libraries)
- **Nested struct support:** Not supported — config files are flat key-value only
- **Short flag `-c` override:** Only `--config` / `--config=` scanned from `os.Args`
- **Auto-detect format in `WithConfigFile`:** Hardcoded to JSON; users must use `WithConfigFileLoader` for auto-detect

### go-output Dependency Management

- **Workspace coupling is FIXED locally** but creates a new problem: external contributors without the exact same local workspace layout (`../go-output`) cannot build this project
- **Replace directives block CI/CD** — GitHub Actions will fail because `../go-output` doesn't exist in CI
- **Need to either:** publish `go-output` tags properly, or vendor, or find another approach

### Documentation

- `FEATURES.md` last updated 2026-05-17 — stale but mostly accurate
- `TODO_LIST.md` last updated 2026-05-16 — partially outdated (Phase 9 items not started)
- `AGENTS.md` comprehensive and current

---

## c) NOT STARTED

From TODO_LIST.md and broader backlog:

1. Interactive prompts (`huh` integration) with `WithPromptOnMissing`
2. Spinner/progress middleware (`bubbles`)
3. Glamour markdown help rendering
4. Telemetry middleware (OpenTelemetry spans)
5. Plugin system for custom validators and type handlers
6. Config file `.env` support
7. Config watching / hot-reload
8. Nested struct config file support
9. Config file write-back / save
10. Config file URL support (HTTP)
11. CLI construction benchmark
12. Flag parsing benchmark
13. Command execution benchmark
14. Benchmark regression detection in CI
15. Codecov integration
16. v2.3.0 release tag and notes
17. Release automation
18. Fix `errors.As` → `errors.AsType[ExitCoder]` (Go 1.26 idiom)
19. Extract `handlerConfig[T,F]` from `wireHandlerWithMiddleware` (8 params)
20. Add `Phase` typed enum to replace `CommandInfo.Phase string`
21. Fix 7 unwrapped error returns (add `fmt.Errorf` context)
22. Consolidate 5 error types into internal `labeledError`
23. Split `type_handler.go` (481 lines) into 3 files
24. Split `command.go` (403 lines) — extract args options
25. Split `flow_context.go` (396 lines) — extract options

---

## d) TOTALLY FUCKED UP!

### The Replace Directive Problem

**Severity:** HIGH for external contributors and CI/CD

The 11 `replace` directives in `go.mod` tie `cmdguard` to a very specific local filesystem layout. If you clone `cmdguard` without also having `go-output` at `../go-output` with the exact same commit state, the build fails.

**Impact:**
- GitHub Actions CI will fail on `go mod download`
- External contributors cannot `go get` or `go install` this module
- `go.work` at parent level partially mitigates this, but only for developers with the exact same monorepo layout

**How we got here:**
`go-output` was modularized into sub-modules (`d2`, `delimited`, `graph`, `markup`, `serialization`, `table`, `testhelpers`, `testhelpers/graphtest`). `cmdguard` imports `go-output` (main) and `go-output/serialization` (for output formats). The `serialization` module has test-only dependencies on `testhelpers/graphtest`, which has no published tag. `go mod tidy` resolves test dependencies of imported packages, so it tries to fetch `testhelpers/graphtest/v0.0.0` and fails.

**Previous attempt failed:** In Phase 5 Quality, we "Removed local go-output replace directive (tagged v0.1.0)" — but that was before `go-output` was further modularized. The modularization introduced new sub-modules that aren't separately tagged.

---

## e) WHAT WE SHOULD IMPROVE!

### Immediate (Critical)

1. **Fix the replace directive / external build problem** — This blocks CI and external usage. Options:
   - **Option A:** Publish proper semver tags for ALL `go-output` sub-modules
   - **Option B:** Remove test-only dependency chain from `serialization` → `graphtest` (upstream fix)
   - **Option C:** Use `go mod edit -exclude` or `-replace` only in CI, not in committed `go.mod`
   - **Option D:** Vendor `go-output` into `cmdguard` (defeats purpose of separate module)
   - **Option E:** Revert `go-output` modularization — make it a single module again

2. **Add tests for configload YAML/TOML loaders** — Currently 0% coverage
3. **Auto-detect format in `WithConfigFile`** by file extension
4. **Error on missing file when `--config` explicitly passed** — silently skipping is wrong UX

### Architecture (High Impact)

5. **Config `--config` should use real cobra flag parsing** — Not `os.Args` string scanning
6. **Add `ConfigSource` enum** — Track which source provided each value (flag, env, file, default)
7. **Plugin system** — Make validators and type handlers truly extensible without modifying library
8. **Consolidate error types** — 5+ similar error structs could become one parameterized type

### Quality (Medium Impact)

9. **Fix 7 unwrapped error returns** — Functions returning bare errors without context
10. **Split oversized files** — `type_handler.go` (481 lines), `command.go` (403 lines), `flow_context.go` (396 lines)
11. **Extract `handlerConfig[T,F]`** — 8-parameter `wireHandlerWithMiddleware` is a smell
12. **Add `Phase` typed enum** — Replace `CommandInfo.Phase string` with compile-time safety

### CI/CD & Release (Medium Impact)

13. **Codecov integration**
14. **Benchmark regression detection**
15. **v2.3.0 release tag and release notes**
16. **Release automation**

---

## f) Top #25 Things We Should Get Done Next

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | **Fix replace directive / external build** | CRITICAL | 1-2h | Dependencies |
| 2 | Add tests for configload YAML/TOML loaders | High | 20m | Quality |
| 3 | Auto-detect format in `WithConfigFile` by extension | High | 15m | UX |
| 4 | Error on missing file when `--config` explicitly passed | High | 10m | UX |
| 5 | Interactive prompts (`huh` integration) | Very High | 4h | Feature |
| 6 | Fix `resolveConfigFlag` to use cobra flag parsing | High | 2h | Architecture |
| 7 | Spinner/progress middleware (`bubbles`) | Medium | 2h | Feature |
| 8 | Glamour markdown help rendering | Medium | 1h | Feature |
| 9 | Telemetry middleware (OpenTelemetry) | Medium | 3h | Feature |
| 10 | Plugin system for validators/type handlers | High | 6h | Architecture |
| 11 | Config file `.env` support | Medium | 30m | Feature |
| 12 | Nested struct config file support | Medium | 2h | Feature |
| 13 | Config file write-back / save | Medium | 1h | Feature |
| 14 | Config watching / hot-reload | Low | 3h | Feature |
| 15 | CLI construction benchmark | Medium | 30m | Performance |
| 16 | Flag parsing benchmark | Medium | 30m | Performance |
| 17 | Command execution benchmark | Medium | 30m | Performance |
| 18 | Benchmark regression detection in CI | Medium | 1h | CI/CD |
| 19 | Codecov integration | Medium | 30m | CI/CD |
| 20 | v2.3.0 release tag and notes | High | 30m | Release |
| 21 | Release automation | Medium | 2h | CI/CD |
| 22 | Consolidate 5 error types into `labeledError` | Medium | 1h | Architecture |
| 23 | Split `type_handler.go` into 3 files | Low | 30m | Cleanup |
| 24 | Split `command.go` — extract args options | Low | 30m | Cleanup |
| 25 | Add `ConfigSource` enum for value tracking | Low | 1h | Architecture |

---

## g) Top #1 Question I Cannot Figure Out Myself

**How should we permanently solve the `go-output` workspace coupling without breaking the modularization?**

The `replace` directives fix the immediate build but create a worse long-term problem (external contributors can't build). The root cause is that `go-output/serialization` (which `cmdguard` imports) has a test-only dependency on `go-output/testhelpers/graphtest`, which has no published tag.

**Options I've considered:**

- **A) Publish tags for all sub-modules** — Requires upstream action on `go-output`. Is there a script or GitHub Action that can auto-tag sub-modules on release?
- **B) Remove the test dependency upstream** — `serialization` tests import `graphtest` for graph serialization tests. Could those tests move to the `graph` package instead?
- **C) Use `GOWORK=off` in CI** — Would force resolution from proxy, but the tags don't exist so it would still fail
- **D) Commit a `go.work` file in `cmdguard`** — Not appropriate for a library module; `go.work` is for local development
- **E) Make `go-output` a single module again** — Reverses the modularization effort. Is the modularization actually providing value?

**What is the intended long-term dependency relationship between `cmdguard` and `go-output`?** Should `cmdguard` pin to stable published versions, or is it expected to co-develop in the same monorepo workspace? If the latter, should we commit a `go.work` file to `cmdguard`? If the former, who publishes the sub-module tags?

---

## Appendix: Raw Metrics

| Metric | Value |
|--------|-------|
| Go version | 1.26.3 |
| Total test runs | 932 |
| v2 package coverage | 84.0% |
| Lint issues | 0 |
| Race conditions | 0 |
| Build status | PASS |
| `go.mod` direct deps | 13 |
| `go.mod` indirect deps | 32 |
| Replace directives | 11 |
| Source files in `pkg/` | 111 |
| Test files in `pkg/` | 68 |
| Example directories | 14 |
| Lines of code (`pkg/`) | ~19,125 |
| Sentinel errors | 40+ |
| Output formats | 12 |
| Value types | 9 |
| Middleware types | 2 built-in + custom |
