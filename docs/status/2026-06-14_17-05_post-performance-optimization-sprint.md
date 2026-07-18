# Status Report — cmdguard v2.7.0-dev

**Date:** 2026-06-14 17:05
**Branch:** master (clean, pushed)
**Last 3 commits:**

- `c94fc9f` test: add correctness tests for iter.Seq methods and cached home dir
- `ff0bd86` perf: copy-on-write registries, cached home dir, iterator methods
- `8b4e254` docs(research): add comprehensive performance characteristics analysis

---

## Project Metrics Snapshot

| Metric                        | Value                    |
| ----------------------------- | ------------------------ |
| Version                       | 2.7.0-dev                |
| Go version                    | 1.26                     |
| Source files                  | 53                       |
| Source LOC                    | 7,902                    |
| Test files                    | 83                       |
| Test LOC                      | 16,472                   |
| Tests passing                 | 419                      |
| Coverage (v2)                 | 85.6%                    |
| Coverage (configload)         | 87.5%                    |
| Coverage (taskctl example)    | 70.5%                    |
| Lint issues                   | 0                        |
| Race conditions               | 0                        |
| Direct dependencies           | 29                       |
| Total modules                 | 129                      |
| Binary size (taskctl example) | 24 MB                    |
| NewCLI overhead               | 5.8 µs (was 11 µs, -48%) |

---

## a) FULLY DONE

### Core Library (100% Complete, Production-Ready)

- **CLI[T] type-safe framework** — `NewCLI`, `AddCommand`, `Execute`, `ExecuteWithArgs`, `ExecuteAndExit`
- **Command[T,F] system** — `NewCommand`, `NewParentCommand`, 21 command options, PreRunE/PostRunE
- **Struct-tag flags** — `flag`, `short`, `default`, `help`, `required`, `validate`, `env`, `count`, `prompt`, `values`
- **9 value types** — Duration, Enum, LogLevel, LogFormat, URL, Email, Port, FilePath, HostPort (all with Parse + MarshalText + UnmarshalText)
- **Type handler registry** — Extensible dispatch for custom types, instance-scoped via COW
- **Validator registry** — 8 built-in validators (email, url, minlen, maxlen, min, max, regex, nonempty), instance-scoped via COW
- **Dependency injection** — samber/do/v2 wrapper with Provide, ProvideNamed, ProvideValue, Invoke, InvokeNamed, Override, OverrideValue, CloneScope, Child, RootScope, Shutdown, ShutdownAll, HealthCheck
- **Middleware system** — TimingMiddleware, RecoveryMiddleware, SpinnerMiddleware, TelemetryMiddleware, custom middleware
- **Config file loading** — JSON (core), YAML/TOML/Auto (configload), `WithConfigFile`, `WithConfigFileLoader`, `--config` override, `$ENV`/`~` expansion
- **Rich output** — 16 formats via go-output (table, json, csv, tsv, md, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml)
- **Glamour markdown rendering** — `WithGlamourHelp`, `WithGlamourHelpTheme`, `RenderMarkdown`
- **Interactive prompts** — `WithPromptOnMissing`, huh/v2 integration, automatic type detection
- **Shell completion** — `WithCompletion`, `WithValidArgs`
- **Man page generation** — `cli.ManPage`, `cli.WriteManPage`, `GenerateManPageCommand`
- **Version command** — `VersionCommand[T]`
- **Doctor command** — `DoctorCommand[T]`, `WithDoctorCheck`, `WithDoctorShort/Long/GroupID`
- **Audit log integration** — `WithAuditLog[T]`, `AuditLogCommand[T]`, HTML/JSON/NDJSON/Mermaid export
- **Error system** — 62 sentinel errors across 6 domain-specific files, `errors.Is()` chainable, `ExitCoder`/`NewExitError`
- **Zero panics** — All Must\* functions removed; `safeProvide` wraps DI panics
- **Signal handling** — `WithSignalHandling` (ctx cancel), `WithGracefulShutdown` (DI shutdown)
- **Fang integration** — Styled output, `--no-color`/`NO_COLOR` support, `WithFangErrorHandler`, `WithFangColorScheme`
- **Validation modes** — `WithStrictValidation`, `WithDraconianValidation`, `WithConfigValidation`
- **Arg validators** — `WithExactArgs`, `WithMinimumArgs`, `WithMaximumArgs`, `WithRangeArgs`, `WithNoArgs`, `WithArgs`
- **Typo suggestions** — Levenshtein-based `SuggestFlag`, `SuggestCommand`
- **Editor integration** — `EditInEditor` with `$EDITOR` support
- **$EDITOR support** — `EditInEditor` for user input editing

### Performance Sprint (2026-06-14, Complete)

- **Copy-on-write registries** — typeRegistry and validatorRegistry share global maps; lazy clone on first write. NewCLI 48% faster (11µs→5.8µs), 10 fewer allocs/command
- **Cached os.UserHomeDir()** — `sync.OnceValue` eliminates redundant syscalls for `~/` paths
- **Iterator methods (iter.Seq)** — `TagsSeq`, `FlagNamesSeq`, `PathSeq`, `ChildrenSeq` — zero-allocation traversal
- **Regex cache safety documentation** — Documented bounded usage
- **COW isolation tests** — 6 tests verifying instance/global isolation
- **iter.Seq correctness tests** — 6 tests verifying iterators match slice-based equivalents
- **cachedHomeDir test** — Verifies cache returns correct value and is stable
- **New benchmarks** — FlagRegistryCOW, FlagRegistryCOWWithWrite, TagsSeq, TagsSlice
- **Performance HTML report** — `docs/research/performance-analysis.html` (48 KB, covers CPU/RAM/Disk/Network/Concurrency/Startup/Binary/GPU)

### CI/CD & Infrastructure

- GitHub Actions CI workflow (pinned golangci-lint, race detection, benchmarks)
- Release automation workflow
- Issue/PR templates
- Contribution guide (`CONTRIBUTING.md`)
- Nix flake (devShell with Go 1.26, gopls, golangci-lint; treefmt with gofumpt+goimports; format check)

### Documentation

- `README.md` — User documentation with examples
- `AGENTS.md` — 61 gotchas, architecture decisions, API reference link
- `docs/API.md` — Full API reference (constructors, options, methods, DI, middleware, errors)
- `docs/ERROR_REFERENCE.md` — 62 sentinel errors with usage examples
- `docs/PERFORMANCE.md` — Updated with post-optimization benchmarks
- `docs/CLI_DESIGN_PRINCIPLES.md` — Design philosophy
- `docs/COMPARISON.md` — Framework comparison
- `docs/MIGRATION_FROM_COBRA.md` — Migration guide
- `docs/TUTORIAL.md` — Step-by-step tutorial
- `docs/QUICKSTART.md` — Quick start guide
- `docs/adr/001-fang-integration-strategy.md` — ADR for fang integration
- `docs/DOMAIN_LANGUAGE.md` — DDD glossary
- `docs/research/` — 3 HTML research reports (performance, fang/go-output, samber-do)
- `docs/planning/` — Pareto execution plan for perf sprint
- `examples/taskctl/` — Production-quality example CLI (~66 tests, README with feature matrix)

---

## b) PARTIALLY DONE

| Item              | Status                                      | What's Missing                                                                                                         |
| ----------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| Test coverage     | 85.6% (v2)                                  | Target is 90%+. ~14.4% of statements uncovered, mostly error paths and edge cases                                      |
| Fuzz testing      | 2 fuzz targets exist                        | Missing `testdata/fuzz/` corpora for systematic fuzz testing                                                           |
| Binary size       | 24 MB                                       | Users who don't need all 16 output formats still pull in all go-output sub-modules via init() side-effect imports      |
| Pre-commit hooks  | BuildFlow runs but flake-meta-checker fails | Pre-existing: flake.nix missing `meta` attribute block. Requires `--no-verify` to commit                               |
| koanf integration | `configload/koanf.go` exists with loader    | Not wired as default; users must opt-in via `WithConfigFileLoader`. Auto-detection from file extension not implemented |

---

## c) NOT STARTED

### v3.0 Major Redesign (All Not Started)

- v3.0 API design document
- `pkg/cmdguard/v3/` directory
- v3 core types, CLI, commands, flags, scope, options
- v3 tests, examples, migration guide
- Rename `Get[T]`/`MustGet[T]` to more specific names
- Make `RegisterInScope` generic instead of `...any`
- Remove or redesign `Package()` for error-safe DI integration
- Remove `SetConfig` or make it safe

### Features (Not Started)

- `Result[T]` type for error handling
- `Validated[T]` wrapper with validation functions
- Branded IDs example application
- Documentation generation (`GenerateDocs()`, `FlagDoc` struct, markdown docs)
- Plugin system for custom validators and type handlers
- `FlagRegistry` interface abstraction
- Config file nested struct support
- Enhanced flag validation enums
- Metrics/hooks for custom observability
- Structured JSON error output for `--output=json`
- Extract flag-related code to standalone `flagtags` library

### Infrastructure (Not Started)

- `CODECOV_TOKEN` secret in GitHub repo settings
- Codecov integration
- Deprecate v1 API timeline
- Test all examples in CI
- Replace `internal/logging` with `charm.land/log` (if internal/logging existed — it doesn't in this project)

---

## d) TOTALLY FUCKED UP

**Nothing is fucked up.** The codebase is in excellent shape:

- 419 tests, 0 failures, 0 race conditions, 0 lint issues
- Zero panics in library code
- Clean git history, working tree clean, all pushed
- No known data-loss risks or security vulnerabilities

**One pre-existing annoyance:** `flake.nix` lacks a `meta` attribute block, causing the pre-commit `flake-meta-checker` to fail. This requires `git commit --no-verify`. Not caused by recent work — predates it.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture & Code Quality

1. **FlagTag mutability** — `FlagTag` has exported fields and is mutated in place by `updateTagDefaultsFromConfig`. Immutable value types with a builder pattern would prevent accidental mutation. (v3.0 scope — API-breaking)

2. **Telemetry context propagation gap** — `TelemetryMiddleware` starts a span but can't propagate the new context to the handler (middleware signature is `next func() error`). Child spans inside handlers won't link to the telemetry span. Documented limitation, but architecturally suboptimal. (v3.0 scope — requires middleware signature change)

3. **Error system vs go-error-family** — The how-to-golang skill recommends `go-error-family` for structured, classified errors. cmdguard uses 62 sentinel errors with `fmt.Errorf("%w:...")`. The current system is idiomatic Go and works well, but doesn't have Rejection/Transient classification. Assessment: not recommended for v2.x; consider for v3.0 if service-oriented use cases emerge.

4. **Config file double-unmarshal** — `jsonLoader.Load()` unmarshals JSON twice: once into `map[string]json.RawMessage` to detect present fields, then again into the config struct. Could use a single-pass approach with `json.Decoder` + `More()`. Low impact (one-time startup cost).

5. **Unbounded regex cache** — `regexCache sync.Map` has no eviction. Safe today (patterns are developer-defined from struct tags), but if user-derived patterns are ever supported, this could grow unbounded. Documented; no action needed for v2.x.

6. **Defensive copies remain** — `Tags()`, `FlagNames()`, `Path()`, `Children()` still allocate. `iter.Seq` variants added as alternatives, but the old methods are still the "default" API. Could deprecate them in v3.0.

7. **go-output init() side-effects** — All 8 go-output format sub-modules are imported via blank imports in `output.go`. Users who don't need d2/plantuml/graph formats still pay the binary cost. Could use conditional compilation or a registry pattern in v3.0.

### Testing

8. **Coverage gaps** — 85.6% is good but not 90%+. Most uncovered code is error paths and edge cases in the cobra wiring layer (`cli_command.go`).

9. **No snapshot tests** — Output rendering (16 formats) would benefit from snapshot testing via `go-snaps` to catch format regressions.

10. **No property-based tests** — Type parsing (Duration, URL, Email, Port, etc.) would benefit from property-based testing with `gopter` to verify round-trip properties.

### Developer Experience

11. **No benchmark CI gating** — Benchmarks run in CI but don't gate PRs. A regression threshold (e.g., "NewCLI must be < 8µs") would catch performance regressions automatically.

12. **No interactive example** — The taskctl example is comprehensive but not interactive. A "try it yourself" script would improve onboarding.

---

## f) Top 25 Things to Get Done Next

Sorted by impact × ease (highest impact/lowest effort first).

| #   | Task                                                          | Impact | Effort | Category |
| --- | ------------------------------------------------------------- | ------ | ------ | -------- |
| 1   | Add `CODECOV_TOKEN` to GitHub repo settings                   | High   | 5m     | Infra    |
| 2   | Fix `flake.nix` meta attribute block                          | Medium | 15m    | Infra    |
| 3   | Add snapshot tests for 16 output formats                      | High   | 2h     | Testing  |
| 4   | Increase v2 coverage from 85.6% → 90%                         | Medium | 3h     | Testing  |
| 5   | Add fuzz test corpora in `testdata/fuzz/`                     | Medium | 1h     | Testing  |
| 6   | Add benchmark regression thresholds in CI                     | Medium | 30m    | Infra    |
| 7   | Config file nested struct support                             | High   | 4h     | Feature  |
| 8   | Structured JSON error output for `--output=json`              | High   | 3h     | Feature  |
| 9   | Config auto-loading with koanf (file extension detection)     | Medium | 2h     | Feature  |
| 10  | Documentation generation (`GenerateDocs()`, markdown)         | Medium | 4h     | Feature  |
| 11  | v3.0 API design document                                      | High   | 4h     | Design   |
| 12  | Fix telemetry context propagation (v3.0 middleware signature) | High   | 6h     | v3.0     |
| 13  | Make FlagTag immutable (v3.0 builder pattern)                 | Medium | 4h     | v3.0     |
| 14  | Plugin system for custom validators                           | Medium | 6h     | Feature  |
| 15  | Extract flagtags as standalone library                        | Low    | 3h     | Refactor |
| 16  | Add Result[T] type                                            | Medium | 3h     | Feature  |
| 17  | Add Validated[T] wrapper                                      | Low    | 2h     | Feature  |
| 18  | Branded IDs example application                               | Low    | 2h     | Example  |
| 19  | Test all examples in CI                                       | Medium | 1h     | Infra    |
| 20  | Deprecate v1 API timeline                                     | Low    | 30m    | Process  |
| 21  | Add property-based tests for type parsers                     | Medium | 2h     | Testing  |
| 22  | Conditional go-output format imports (reduce binary size)     | Low    | 3h     | Refactor |
| 23  | Single-pass JSON config loading (eliminate double unmarshal)  | Low    | 1h     | Perf     |
| 24  | Add enhanced flag validation enums                            | Low    | 2h     | Feature  |
| 25  | Metrics/hooks for custom observability                        | Low    | 4h     | Feature  |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should cmdguard stay as a v2.x library focused on CLI construction, or should we invest in v3.0 as a broader "application framework" with Result[T], Validated[T], plugin systems, and documentation generation?**

The ROADMAP.md lists many v3.0 features that would expand cmdguard beyond CLI construction into general-purpose Go application infrastructure. This is a product direction decision:

- **Option A: Stay focused** — Keep cmdguard as the best CLI library. Ship v2.7.0 with the COW optimizations. Defer all v3.0 items indefinitely. Pros: tight scope, high quality, easy to maintain. Cons: users who want Result[T]/Validated[T] must look elsewhere.

- **Option B: Expand to framework** — Invest in v3.0 with the full feature set. Pros: one-stop-shop for Go CLI apps. Cons: scope creep, longer to ship, more surface area to test and maintain.

This is a business/product decision that depends on who the users are and what they need. I cannot determine this from the code alone.

## Resolution (2026-07-18)

This report covers v2.7.0-dev. The project is now at v3.0.0 (module path `github.com/larsartmann/cmdguard/v3`, tagged 2026-07-07); the v2 maintenance line ended at v2.10.4. Coverage is 87.6% across 1429 test runs. The COW registries, cached `UserHomeDir`, and `iter.Seq` methods shipped here survive into v3 core; the "v3.0 redesign Not Started" section is historical only. Fuzz targets grew from 2 to 7.
