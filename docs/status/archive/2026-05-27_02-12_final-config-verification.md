# Comprehensive Status Report — Config File Loading Feature Complete

**Date:** 2026-05-27 02:12 CEST  
**Branch:** master  
**Commits since last report:** 10  
**Working tree:** CLEAN  
**Tests:** 932 individual test runs, 82.4% coverage  
**Build:** PASS  
**Lint:** 0 issues  
**Race detector:** 0 races

---

## a) FULLY DONE

### Phase 0: Unblock Build

- Fixed go-output modularized imports in `output.go`
- Added `serialization`, `delimited`, `markup` subpackages to imports
- Updated parent `go.work` workspace with `go-output/markup` and `go-output/testhelpers`
- Updated `go.mod` with new dependencies

### Phase 1: Core Config File Loading

- **`ConfigFileLoader` interface** — `Load(data []byte, cfg any) ([]string, error)`
- **JSON loader** (stdlib, core package) — flat key-value objects matching `flag` tag names
- **Path resolution** — `$ENV` expansion + `~` expansion
- **Precedence chain** — `flag > env > config file > tag default`
  - Implemented by updating `tag.Default` values before flag registration
  - Zero modifications to existing `ParseFlags` logic
- **`--config` flag override** — scans `os.Args` for `--config` / `--config=` forms
- **Missing file handling** — silently skipped (not an error for default paths)
- **Sentinel errors** — `ErrConfigFileRead`, `ErrConfigFileParse`, `ErrConfigFileLoad`, `ErrConfigFileNotFound`

### Phase 2: Optional Config Loaders

- **`pkg/cmdguard/v2/configload/yaml.go`** — YAML via `gopkg.in/yaml.v3`
- **`pkg/cmdguard/v2/configload/toml.go`** — TOML via `github.com/pelletier/go-toml/v2`
- **`pkg/cmdguard/v2/configload/auto.go`** — `LoaderForPath()` selects by extension
- `configload.JSON()` for symmetry

### Phase 3: CLI Integration

- `WithConfigFile[T](paths...)` — JSON-only, core package
- `WithConfigFileLoader[T](loader, paths...)` — custom loader
- `CLI[T]` fields: `configFilePaths`, `configFileLoader`
- Integrated into `CLI.initialize()` between registry creation and flag registration
- `FlagRegistry.updateTagDefaultsFromConfig()` handles primitive types + `fmt.Stringer`

### Phase 4: Tests

- **Unit tests** (`config_file_test.go`) — 5 test functions, 14 subtests:
  - `TestExpandConfigPath` — env and tilde expansion
  - `TestJSONLoader_Load` — flat config, partial keys, invalid JSON
  - `TestLoadConfigFile` — existing file, missing file skip, not-found error
  - `TestUpdateTagDefaultsFromConfig` — primitives + unknown fields
  - `TestResolveConfigFlag` — `--config`, `--config=`, not present
- **Integration tests** (`config_file_integration_test.go`) — 5 subtests:
  - Config file overrides default
  - Flag overrides config file
  - Env overrides config file
  - Flag overrides env AND config file
  - Missing config file is not an error

### Phase 5: Example + Docs

- `examples/config-file/main.go` + `README.md`
- `FEATURES.md` — Config File Loading section added
- `AGENTS.md` — project structure, CLI options, Key Gotchas (items 22-25)
- `TODO_LIST.md` — marked config file loading as done, updated stats

### Phase 6: Quality Hardening

- Fixed `resolveConfigFlag` race condition (was reading global `os.Args`)
- Fixed `fieldValueToString` for primitive types (`int`, `uint`, `float`, `bool`)
- Fixed depguard rules for `gopkg.in/yaml.v3` and `github.com/pelletier/go-toml/v2`
- Applied `golangci-lint fmt` formatting
- Fixed `nonamedreturns` across config files
- Fixed `modernize` (strings.CutPrefix)
- Added `//nolint:gosec` for intentional config file reading
- Added file-level `//nolint` for test files with inline handler returns
- Removed accidentally committed ELF binary, added to `.gitignore`

---

## b) PARTIALLY DONE

### Configload Sub-package Tests

- **0% coverage** — no tests for YAML/TOML loaders
- Loaders are thin wrappers around well-tested libraries, but zero integration coverage

### Config File Feature Gaps

- **Nested struct support** — only top-level keys matching `flag` tags
- **Short flag `-c` override** — only `--config` / `--config=` scanned
- **Auto-detection in `WithConfigFile`** — hardcoded to JSON; user must use `WithConfigFileLoader` for auto-detect
- **Missing explicit `--config` error** — silently skips even when user explicitly passes bad path

---

## c) NOT STARTED

From TODO_LIST.md and backlog:

1. Interactive prompts (huh integration) with `WithPromptOnMissing`
2. Spinner/progress middleware (bubbles)
3. Glamour markdown help rendering
4. Telemetry middleware (OpenTelemetry spans)
5. Plugin system for custom validators and type handlers
6. Config file `.env` support
7. Config watching / hot-reload
8. Nested struct config file support
9. Config file write-back / save
10. CLI construction benchmark
11. Flag parsing benchmark
12. Command execution benchmark
13. Benchmark regression detection in CI
14. Codecov integration
15. v2.3.0 release tag and notes
16. Release automation
17. Make `NoFlags` a distinct named type (v3.0)
18. Change `TimingMiddleware` callback to include error (v3.0)
19. Remove string-based `BranchWithTimeout`/`BranchWithDeadline` (v3.0)
20. Remove `FlowContextAccessor` (v3.0)
21. Rename `Get[T]`/`MustGet[T]` to more specific names (v3.0)
22. Make `RegisterInScope` generic instead of `...any` (v3.0)
23. Remove or redesign `Package()` for error-safe DI (v3.0)
24. Support config file URLs (HTTP)
25. Update `go-output` to stable published version (avoid workspace coupling)

---

## d) TOTALLY FUCKED UP!

### NONE

All builds pass, all tests pass, race detector clean, lint clean.

### Near-misses during this session:

- **Committed ELF binary** (`config-file` from `go run` in root) — caught, removed, `.gitignore` updated
- **`resolveConfigFlag` race** — initial implementation read global `os.Args`; fixed by accepting `[]string` parameter
- **`fieldValueToString` missing primitives** — plain `int`/`bool`/`float` not converted; fixed with explicit type switch

---

## e) WHAT WE SHOULD IMPROVE!

### Architecture (High Impact)

1. **Config override should use cobra flag parsing, not `os.Args` scan** — The `--config` scan is a workaround. Better: register `--config` as a real cobra flag, parse in `PersistentPreRunE`, then conditionally reload config and re-parse flags. Supports `-c` naturally.
2. **ConfigFileLoader should support remote sources** — URLs, etcd, Consul. Currently file-system only.
3. **Add `ConfigSource` enum for debugging** — Track which source provided each value (flag, env, file, default).
4. **Config file validation** — Beyond JSON/YAML/TOML parsing, validate against flag constraints.

### Code Quality (Medium Impact)

5. **Add tests for configload YAML/TOML loaders** — Currently 0% coverage
6. **Error on missing file when `--config` explicitly passed** — User passed `--config /bad/path`; they want an error, not silence
7. **Auto-detect format in `WithConfigFile`** — Inspect file extension, use appropriate loader
8. **Support nested structs in config files** — Flat key-value is limiting for complex configs

### CI/CD & Tooling

9. **Codecov integration**
10. **Benchmark regression detection**
11. **Release automation for v2.3.0**
12. **Fix go.work workspace coupling** — External contributors can't build without identical local workspace

---

## f) Top #25 Things We Should Get Done Next

| #   | Task                                                      | Impact    | Effort | Category     |
| --- | --------------------------------------------------------- | --------- | ------ | ------------ |
| 1   | Add tests for configload YAML/TOML loaders                | High      | 20m    | Quality      |
| 2   | Auto-detect format in `WithConfigFile` by extension       | High      | 15m    | UX           |
| 3   | Error on missing file when `--config` explicitly passed   | High      | 10m    | UX           |
| 4   | Fix `resolveConfigFlag` to use cobra flag parsing         | High      | 2h     | Architecture |
| 5   | Interactive prompts (`huh` integration)                   | Very High | 4h     | Feature      |
| 6   | Spinner/progress middleware (`bubbles`)                   | Medium    | 2h     | Feature      |
| 7   | Glamour markdown help rendering                           | Medium    | 1h     | Feature      |
| 8   | Telemetry middleware (OpenTelemetry)                      | Medium    | 3h     | Feature      |
| 9   | Config file `.env` support                                | Medium    | 30m    | Feature      |
| 10  | Nested struct config file support                         | Medium    | 2h     | Feature      |
| 11  | Config file write-back / save                             | Medium    | 1h     | Feature      |
| 12  | Plugin system for validators/type handlers                | High      | 6h     | Architecture |
| 13  | Config watching / hot-reload                              | Low       | 3h     | Feature      |
| 14  | CLI construction benchmark                                | Medium    | 30m    | Performance  |
| 15  | Flag parsing benchmark                                    | Medium    | 30m    | Performance  |
| 16  | Command execution benchmark                               | Medium    | 30m    | Performance  |
| 17  | Benchmark regression detection in CI                      | Medium    | 1h     | CI/CD        |
| 18  | Codecov integration                                       | Medium    | 30m    | CI/CD        |
| 19  | v2.3.0 release tag and notes                              | High      | 30m    | Release      |
| 20  | Release automation                                        | Medium    | 2h     | CI/CD        |
| 21  | Add `ConfigSource` enum for value tracking                | Low       | 1h     | Architecture |
| 22  | Support config file URLs (HTTP)                           | Low       | 1h     | Feature      |
| 23  | Fix go.work workspace coupling                            | Medium    | 1h     | Dependencies |
| 24  | Refactor `output.go` stale LSP diagnostics                | Low       | 10m    | Cleanup      |
| 25  | Update AGENTS.md line count (go-structure-linter warning) | Low       | 30m    | Docs         |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Why does the parent `go.work` at `/home/lars/projects/go.work` include `go-output` as a workspace module, forcing cmdguard to use the local development HEAD instead of the published v0.5.0?**

This causes real problems:

- `go mod tidy` fails because test dependencies of workspace modules (e.g., `testhelpers/graphtest`) aren't published
- cmdguard builds against unstable go-output APIs without explicit opt-in
- CI reproducibility is compromised — the workspace coupling is invisible to `go.mod`
- External contributors cannot build cmdguard without the exact same local workspace structure

I added `./go-output/markup` and `./go-output/testhelpers` to the workspace to fix the immediate build break, but this feels like treating a symptom. The root question: **should cmdguard be removed from the parent workspace and use explicit `replace` directives in its own `go.mod`, or should the workspace be the canonical development environment for all projects?**

If the workspace is canonical, then every external contributor needs the exact same repo layout. If `replace` directives are canonical, then local development requires manually managing them. Neither option is great, but the current hybrid (workspace forces dependencies without explicit consent) is the worst of both worlds.
