# Comprehensive Status Report

**Date:** 2026-05-27 01:43 CEST  
**Branch:** master  
**Commits since last report:** 8  
**Tests:** 869 individual test runs, 82.4% coverage  
**Build:** PASS  
**Race detector:** 0 races

---

## a) FULLY DONE

### 1. Fixed Broken Build (go-output modularization)

- Updated `pkg/cmdguard/v2/output.go` to use new go-output subpackages:
  - `github.com/larsartmann/go-output/serialization` (MarshalJSON, MarshalYAML)
  - `github.com/larsartmann/go-output/delimited` (MarshalTSV, NewCSVWriter)
  - `github.com/larsartmann/go-output/markup` (MarshalXMLFromTableData, NewHTMLRenderer)
- Added missing submodules to parent `go.work` workspace
- Updated `go.mod` with new dependencies

### 2. Config File Loading — Core Feature

- **`ConfigFileLoader` interface** — `Load(data []byte, cfg any) ([]string, error)`
- **JSON loader** (stdlib) — flat key-value objects matching `flag` tag names
- **Path resolution** — `$ENV` expansion + `~` expansion via `os.ExpandEnv` / `os.UserHomeDir`
- **Precedence chain** — `flag > env > config file > tag default`
  - Implemented by updating `tag.Default` values before flag registration
  - Reuses existing `ParseFlags` without modification
- **`--config` flag override** — scans `os.Args` for `--config` / `--config=` forms
- **Missing file handling** — silently skipped (not an error)
- **Sentinel errors** — `ErrConfigFileRead`, `ErrConfigFileParse`, `ErrConfigFileLoad`, `ErrConfigFileNotFound`

### 3. Config File Loaders — Optional Sub-package

- **`pkg/cmdguard/v2/configload/yaml.go`** — YAML loader via `gopkg.in/yaml.v3`
- **`pkg/cmdguard/v2/configload/toml.go`** — TOML loader via `github.com/pelletier/go-toml/v2`
- **`pkg/cmdguard/v2/configload/auto.go`** — `LoaderForPath()` selects by extension
- `configload.JSON()` provided for symmetry

### 4. CLI Integration

- `WithConfigFile[T](paths...)` — JSON-only, core package
- `WithConfigFileLoader[T](loader, paths...)` — custom loader
- Integrated into `CLI.initialize()` between registry creation and flag registration
- `FlagRegistry.updateTagDefaultsFromConfig()` updates tag defaults from loaded values
- `CLI[T]` fields: `configFilePaths`, `configFileLoader`

### 5. Tests

- **Unit tests** (`config_file_test.go`):
  - `TestExpandConfigPath` — env and tilde expansion
  - `TestJSONLoader_Load` — flat config, partial keys, invalid JSON
  - `TestLoadConfigFile` — existing file, missing file skip, not-found error
  - `TestUpdateTagDefaultsFromConfig` — primitive types + unknown fields
  - `TestResolveConfigFlag` — `--config`, `--config=`, not present
- **Integration tests** (`config_file_integration_test.go`):
  - `TestConfigFilePrecedence` — 5 subtests covering full precedence chain
  - Config file overrides default
  - Flag overrides config file
  - Env overrides config file
  - Flag overrides env AND config file
  - Missing config file is not an error

### 6. Example

- `examples/config-file/` with `main.go` + `README.md`
- Demonstrates `WithConfigFile`, `--config` override, precedence

### 7. Documentation Updates

- **FEATURES.md** — added "Config File Loading" section with status table
- **AGENTS.md** — updated project structure, CLI options table, Key Gotchas (items 22-25)
- **TODO_LIST.md** — marked config file loading as done, updated test count

---

## b) PARTIALLY DONE

### Lint / Formatting

- `golangci-lint fmt` applied to most files
- Some `nlreturn` warnings remain in `config_file_integration_test.go` (non-blocking)
- `gosec` suppressed with `//nolint:gosec` for intentional config file reading

### Config File Feature Gaps

- **Nested struct support** — JSON/YAML/TOML loaders only detect top-level keys matching `flag` tags
- **Short flag `-c` override** — only `--config` / `--config=` are scanned; `-c` short form not supported for override
- **Auto-detection in CLI** — `configload.Auto()` exists but `WithConfigFile` doesn't auto-detect; user must use `WithConfigFileLoader(configload.Auto(), paths...)`

---

## c) NOT STARTED

From TODO_LIST.md and feature backlog:

1. Interactive prompts (huh integration) with `WithPromptOnMissing`
2. Spinner/progress middleware (bubbles)
3. Glamour markdown help rendering
4. Telemetry middleware (OpenTelemetry spans)
5. Plugin system for custom validators and type handlers
6. Config file auto-loading with `.env` support
7. Config watching / hot-reload
8. CLI construction benchmark
9. Flag parsing benchmark
10. Command execution benchmark
11. Benchmark regression detection in CI
12. Codecov integration
13. v2.3.0 release tag and notes
14. Release automation
15. Make `NoFlags` a distinct named type (v3.0)
16. Change `TimingMiddleware` callback to include error (v3.0)
17. Remove string-based `BranchWithTimeout`/`BranchWithDeadline` (v3.0)
18. Remove `FlowContextAccessor` (v3.0)
19. Rename `Get[T]`/`MustGet[T]` to more specific names (v3.0)
20. Make `RegisterInScope` generic instead of `...any` (v3.0)
21. Remove or redesign `Package()` for error-safe DI (v3.0)

---

## d) TOTALLY FUCKED UP!

### NONE

All builds pass, all tests pass, race detector is clean. No regressions introduced.

### Near-misses during this session:

- **Accidentally committed ELF binary** (`config-file` from `go run` in root) — caught and removed, added to `.gitignore`
- **`resolveConfigFlag` race condition** — initial implementation read global `os.Args`, causing data races in parallel tests. Fixed by accepting `[]string` parameter.
- **`fieldValueToString` missing primitives** — initial implementation only used `getFieldValue` which relies on `fmt.Stringer`. Plain `int`/`bool`/`float` fields weren't being converted. Fixed by adding primitive type handling.

---

## e) WHAT WE SHOULD IMPROVE!

### Architecture

1. **Config file loading should not read `os.Args` at all** — The `--config` override scan is a hack. Better design: register `--config` as a true cobra flag, parse it in `PersistentPreRunE`, then conditionally re-load config and re-parse flags. This would also support `-c` naturally.
2. **ConfigFileLoader interface is synchronous only** — No support for async/remote config sources (etcd, Consul, HTTP URL).
3. **No config file validation** — We validate flags but don't validate config file contents beyond JSON/YAML/TOML parsing.
4. **No config file write-back** — Can't save modified config back to disk.

### Code Quality

5. **nlreturn warnings in tests** — 6 occurrences in `config_file_integration_test.go`
6. **gci formatting drift** — Some imports not perfectly grouped
7. **Test coverage for configload sub-package** — 0% coverage; no tests for YAML/TOML loaders

### API Design

8. **`WithConfigFile` should auto-detect format** — Currently hardcoded to JSON; should inspect extension
9. **Config file paths should support URLs** — `https://example.com/config.yaml` would be useful
10. **Missing config file with explicit `--config` should be an error** — Currently silently skips; if user explicitly passes `--config /bad/path`, they probably want an error.

---

## f) Top #25 Things We Should Get Done Next

| #   | Task                                                            | Impact    | Effort | Category     |
| --- | --------------------------------------------------------------- | --------- | ------ | ------------ |
| 1   | Fix `nlreturn` warnings in config tests                         | Low       | 5m     | Quality      |
| 2   | Add tests for configload YAML/TOML loaders                      | High      | 20m    | Quality      |
| 3   | Auto-detect format in `WithConfigFile` by extension             | High      | 15m    | UX           |
| 4   | Error on missing file when `--config` explicitly passed         | High      | 10m    | UX           |
| 5   | Interactive prompts (`huh` integration)                         | Very High | 4h     | Feature      |
| 6   | Spinner/progress middleware (`bubbles`)                         | Medium    | 2h     | Feature      |
| 7   | Glamour markdown help rendering                                 | Medium    | 1h     | Feature      |
| 8   | Telemetry middleware (OpenTelemetry)                            | Medium    | 3h     | Feature      |
| 9   | Config file `.env` support                                      | Medium    | 30m    | Feature      |
| 10  | Plugin system for validators/type handlers                      | High      | 6h     | Architecture |
| 11  | Nested struct config file support                               | Medium    | 2h     | Feature      |
| 12  | Config file write-back / save                                   | Medium    | 1h     | Feature      |
| 13  | Config watching / hot-reload                                    | Low       | 3h     | Feature      |
| 14  | CLI construction benchmark                                      | Medium    | 30m    | Performance  |
| 15  | Flag parsing benchmark                                          | Medium    | 30m    | Performance  |
| 16  | Command execution benchmark                                     | Medium    | 30m    | Performance  |
| 17  | Benchmark regression detection in CI                            | Medium    | 1h     | CI/CD        |
| 18  | Codecov integration                                             | Medium    | 30m    | CI/CD        |
| 19  | v2.3.0 release tag and notes                                    | High      | 30m    | Release      |
| 20  | Release automation                                              | Medium    | 2h     | CI/CD        |
| 21  | Fix `output.go` stale LSP diagnostics                           | Low       | 10m    | Cleanup      |
| 22  | Add `ConfigSource` enum for debugging                           | Low       | 1h     | Architecture |
| 23  | Support config file URLs (HTTP)                                 | Low       | 1h     | Feature      |
| 24  | Refactor `resolveConfigFlag` to use cobra flag parsing          | Medium    | 2h     | Architecture |
| 25  | Update `go-output` to stable v0.5.0 to avoid workspace coupling | Medium    | 1h     | Dependencies |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Why does the parent `go.work` at `/home/lars/projects/go.work` include `go-output` as a workspace module, forcing cmdguard to use the local development HEAD instead of the published v0.5.0?**

This causes several problems:

- `go mod tidy` fails because test dependencies of workspace modules (e.g., `testhelpers/graphtest`) aren't published
- cmdguard builds against unstable go-output APIs without opt-in
- CI reproducibility is compromised because the workspace coupling is invisible to `go.mod`

I added `./go-output/markup` and `./go-output/testhelpers` to the workspace to fix the immediate build break, but this feels like treating a symptom. The root question: **should cmdguard be removed from the workspace and use explicit `replace` directives in its own `go.mod`, or should the workspace be the canonical development environment for all projects?** The current setup makes it impossible for external contributors to build cmdguard without also having the exact same local workspace structure.
