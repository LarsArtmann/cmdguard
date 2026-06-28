# TODO List

**Updated:** 2026-06-28
**Status:** v2.10.0 — zero panics, 86.6% coverage, 0 lint issues, 0 race conditions, 16 output formats, 11 audit log formats, copy-on-write registries, nested config, plugin system, Result/Validated types, GenerateDocs, cobra-correctness contract (SilenceUsage default, ExitCode, escape-hatch APIs)
**Tests:** 473 test functions (1355 runs incl. subtests), 26 benchmarks, 1 fuzz file, 0 build errors

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

- [x] Remove all Must\* panic-inducing functions (16 functions deleted)
- [x] Update FEATURES.md: remove Must\* entries, update metrics
- [x] Update README.md: remove Must\* examples, update tagline
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

### Phase 16: Performance Optimization Sprint (2026-06-14)

- [x] Performance analysis: comprehensive HTML report at `docs/research/performance-analysis.html`
- [x] Copy-on-write typeRegistry — eliminates 10 allocs/command, 48% faster NewCLI
- [x] Copy-on-write validatorRegistry — same COW pattern
- [x] Cache `os.UserHomeDir()` via `sync.OnceValue` — eliminates redundant syscalls
- [x] Iterator methods: `TagsSeq()`, `FlagNamesSeq()`, `PathSeq()`, `ChildrenSeq()` — zero-allocation traversal
- [x] Document regex cache safety bounds
- [x] Add COW isolation tests (6 tests)
- [x] Add benchmarks: FlagRegistryCOW, FlagRegistryCOWWithWrite, TagsSeq, TagsSlice
- [x] Update PERFORMANCE.md with post-optimization numbers
- [x] Update AGENTS.md gotchas (#59-61: COW, cached home dir, iterators)

### Phase 17: Dependency Maximization & Doc Sync (2026-06-17)

- [x] Verify all direct dependencies at latest published versions
- [x] Upgrade indirect deps: `rogpeppe/go-internal` v1.15.0, `charmbracelet/x/conpty`, `charmbracelet/x/exp/golden`, `go-output/testhelpers/graphtest`
- [x] Mark koanf config loading as completed (was already integrated, stale in remaining work)
- [x] Sync stale version references across FEATURES.md, README.md, ROADMAP.md, AGENTS.md

### Phase 18: Audit Log Format Expansion & Doc Sync (2026-06-21)

- [x] Upgrade `samber-do-auditlog` v0.1.0 → v0.3.0 (4 new export formats: d2, plantuml, tree, htmltree → 11 total)
- [x] Remove local audit-log adapter functions (superseded by Plugin-level methods)
- [x] Refresh all transitive dependencies via `go get -u all` + `go mod tidy`
- [x] Fix stale audit log format lists in docs (7 → 11 formats)
- [x] Brutal self-review report + Go 1.26.4 security TODO (govulncheck GO-2026-5037/5038/5039)

### Phase 19: Cobra-Correctness Contract & Escape-Hatch APIs (2026-06-28)

Mission pivot back to "make consumers use Cobra correctly," driven by auditing BuildFlow (the primary consumer) and replacing its four workarounds with first-class APIs.

- [x] Close the cobra-correctness contract: `SilenceUsage=true` by default, public `ExitCode(err)` helper, flagship example no longer exits 0 on failure or double-prints errors
- [x] Add scoped flags (`local:"true"` tag) — root-only flags not inherited by subcommands
- [x] Add `hidden:"true"` flag tag — exclude from `--help`, stay functional
- [x] Add `ConfigFromContext[T](ctx)` — type-safe config retrieval for raw cobra subcommands (the escape hatch)
- [x] Add `WithPostFlagParse[T](fns ...)` — post-parse hook (DI init, session storage)
- [x] Disable `makezero` linter (directly conflicts with staticcheck S1019)
- [x] Bump `go.mod` → `go 1.26.4` (nixpkgs now ships it; fixes GO-2026-5037/5038/5039)
- [x] Untrack stray `taskctl-audit.html` generated artifact + gitignore generated example HTML
- [x] Update CHANGELOG.md, FEATURES.md, docs/API.md for all v2.10.0 changes

## Remaining Work — Priority Sorted

### P0: Open

| #   | Task                                                                                                              | Files              | Effort |
| --- | ----------------------------------------------------------------------------------------------------------------- | ------------------ | ------ |
| 20  | Add `CODECOV_TOKEN` secret to GitHub repo settings (requires repo owner — cannot be set programmatically)         | GitHub settings    | 5m     |

### P1: Future (v3.0+)

| #   | Task                                                                                     | Category |
| --- | ---------------------------------------------------------------------------------------- | -------- |
| 21  | ~~Plugin system for custom validators and type handlers~~ ✅ DONE (v2.8)                 |
| 22  | ~~Config file nested struct support~~ ✅ DONE (v2.8)                                     |
| 23  | ~~Documentation generation (GenerateDocs, markdown)~~ ✅ DONE (v2.8)                     |
| 24  | ~~Advanced types: Result[T], Validated[T]~~ ✅ DONE (v2.8)                               |
| 25  | ~~Structured JSON error output for `--output=json`~~ ✅ DONE (v2.7)                      |
| 26  | Extract flag-related code to standalone `flagtags` library                               | Refactor |
| 27  | ~~Consider extracting `go-output` to sub-package~~ ✅ DONE (already external at v0.12.0) |

### P2: Future Cleanup (API-breaking, defer to v3.0)

| #   | Task                                                           |
| --- | -------------------------------------------------------------- |
| 30  | Rename `Get[T]`/`MustGet[T]` to more specific names            |
| 31  | Make `RegisterInScope` generic instead of `...any`             |
| 32  | Remove or redesign `Package()` for error-safe DI integration   |
| 33  | Remove `SetConfig` or make it safe (reinitialize FlagRegistry) |
