# TODO List

**Updated:** 2026-06-10
**Status:** v2.5.0 — zero panics, 85.4% coverage, 0 lint issues, 0 race conditions, 16 output formats
**Tests:** 395+ passing, 0 build errors

## Completed

### Phase 1–9: All Complete

- [x] All core features implemented and tested
- [x] All architecture hardening complete
- [x] All documentation updated
- [x] Nix flake with devShell, formatter, and format check
- [x] CI with pinned golangci-lint, codecov, Nix check, benchmarks
- [x] Release automation workflow

### Phase 10: Post-Release Maintenance (2026-06-08)

- [x] Fix `flake.nix` infinite recursion (`goPkg = goPkg` → `pkgs.go_1_26`)
- [x] Update FEATURES.md version from v2.3.0-dev to v2.4.0
- [x] Fix test count metrics (357, not 356)
- [x] Add gofumpt and goimports to flake.nix treefmt
- [x] Add `Scope.HealthCheckResults()` / `HealthCheckResultsWithContext()`
- [x] Add `CLI.HealthCheckResults()` / `HealthCheckResultsWithContext()`
- [x] Add `DoctorCommand[T]` convenience helper
- [x] DRY configload: extract `genericLoader` (3 files → 1)
- [x] Add configload tests: YAML, TOML, JSON, Auto, LoaderForPath (22 tests)
- [x] Consolidate `command_suggest.go` into `flags_suggest.go`
- [x] Update taskctl example: manual health → DoctorCommand
- [x] Update docs: FEATURES.md, TODO_LIST.md, AGENTS.md

### Phase 11: Codebase Review (2026-06-10)

- [x] Code quality scan: build, lint, duplication analysis — 0 issues
- [x] Full code review: all 50 source files reviewed
- [x] Architecture review: modularity (8.5/10), scalability (9/10), composability (7/10)
- [x] Architecture visualization: D2 diagrams (current + improved)
- [x] Docs freshness check: fixed stale items in AGENTS.md, FEATURES.md
- [x] Naming review: 9/10 quality — 3 minor issues
- [x] Architecture deepening: 6 candidates identified
- [x] Go modularization assessment: NOT recommended (project too small)
- [x] Features audit: all features verified against code
- [x] TODO list rebuilt from all .md sources

### Phase 12: Zero Panics (2026-06-10)

- [x] Remove all Must* panic-inducing functions (16 functions deleted)
- [x] Update FEATURES.md: remove Must* entries, update metrics
- [x] Update README.md: remove Must* examples, update tagline
- [x] Update all docs for zero-panics guarantee

### Phase 13: samber/do v2 Utilization Sprint (2026-06-10)

- [x] Research: full samber/do v2 API surface audit (54 public symbols)
- [x] Add `WithGracefulShutdown[T]()` — graceful DI shutdown on SIGINT/SIGTERM
- [x] Add `Override[T]` / `OverrideValue[T]` — replace services for testing
- [x] Add `CloneScope(scope)` — clone DI scope for test isolation
- [x] Add `NewScopeWithOpts(name, opts)` — create scope with custom InjectorOpts
- [x] Add `WithDILogging[T](logf)` — DI container internal logging
- [x] Update `WithSignalHandling` doc to clarify context-only behavior
- [x] Add research report: `docs/research/samber-do-v2-utilization.html`
- [x] Update taskctl example: `WithSignalHandling` → `WithGracefulShutdown`
- [x] Add Clone+Override test example in taskctl
- [x] samber/do utilization: 24% → 43% (13 → 23 of 54 API symbols)

### Phase 14: Post-Sprint Cleanup (2026-06-10)

- [x] Make `NoFlags` a distinct named type (not struct{} alias) — P6 #29
- [x] Remove deprecated `WithColor` option — use `WithFang` instead — P6 #34
- [x] Fix `NO_COLOR` env var restored after execution instead of permanently mutated — P6 #35
- [x] Extract API reference from AGENTS.md to docs/API.md
- [x] Add comprehensive error reference (62 sentinels) to docs/ERROR_REFERENCE.md
- [x] Consolidate type handler registration with helper function
- [x] Add tests for 6 previously 0%-covered functions
- [x] Add `Scope.RootScope()` accessor
- [x] Add DI benchmarks: NewScopeWithOpts, CloneScope, ProvideInvokeCycle
- [x] Clean up unused types and imports in test files

### Phase 15: Library Integration Sprint (2026-06-10)

- [x] Add `WithCLICommit[T](commit)` — auto-pipes to fang
- [x] Add `WithFangErrorHandler[T](handler)` and `WithFangColorScheme[T](cs)`
- [x] Add 4 new output formats: JSONL, AsciiDoc, TOML, PlantUML (16 total)
- [x] Document new APIs in doc.go and docs/API.md
- [x] Deduplicate `validateEmail`/`validateURL` — delegate to `ParseEmail`/`ParseURL`
- [x] Fix `HostPort.IsEmpty()` coupling — use `hp.port.IsEmpty()` instead of `hp.port.port`
- [x] Add IsEmpty() tests for Duration, LogLevel, LogFormat, Port
- [x] Add ADR-001 for fang integration strategy

## Remaining Work — Priority Sorted

### P0: Open

| #  | Task                                                   | Files           | Effort |
| --- | ------------------------------------------------------ | --------------- | ------ |
| 20 | Add `CODECOV_TOKEN` secret to GitHub repo settings     | GitHub settings | 5m     |

### P1: Future (v3.0+)

| #  | Task                                                        | Category |
| --- | ----------------------------------------------------------- | -------- |
| 21  | Plugin system for custom validators and type handlers       | Feature  |
| 22  | Config file nested struct support                           | Feature  |
| 23  | Documentation generation (GenerateDocs, markdown, API docs) | Feature  |
| 24  | Advanced types: Result[T], Validated[T], branded IDs        | Feature  |
| 25  | Config auto-loading with koanf integration                  | Feature  |
| 26  | Structured JSON error output for `--output=json`            | Feature  |
| 27  | Extract flag-related code to standalone `flagtags` library  | Refactor |
| 28  | Consider extracting `go-output` to sub-package             | Refactor |

### P2: Future Cleanup (API-breaking, defer to v3.0)

| #  | Task                                                          |
| --- | ------------------------------------------------------------- |
| 30  | Rename `Get[T]`/`MustGet[T]` to more specific names          |
| 31  | Make `RegisterInScope` generic instead of `...any`            |
| 32  | Remove or redesign `Package()` for error-safe DI integration |
| 33  | Remove `SetConfig` or make it safe (reinitialize FlagRegistry) |
