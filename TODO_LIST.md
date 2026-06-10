# TODO List

**Updated:** 2026-06-10
**Status:** v2.5.0 — zero panics achieved
**Tests:** 368 passing, 83.5% coverage, 0 lint issues, 0 race conditions

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
- [x] Add `DoctorCommand[T]` / `MustDoctorCommand[T]` convenience helper
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

## Remaining Work — Priority Sorted

### P0: Bugs & Data Correctness

| # | Task | Files | Effort |
|---|------|-------|--------|
| 1 | ~~Fix `MergeConfigs` zero-value sentinel~~ → Investigated: IsZero() skip is intentional for flag merging; documented behavior | config.go | ✅ Done |
| 2 | ~~Fix `Enum.UnmarshalText` validation bypass~~ → Investigated: zero-value Enum bootstrap is intentional design | types_enum.go | ✅ Done |
| 3 | Fix `setStringField` panic on type mismatch: add AssignableTo guard before `field.Set` | config_setfield.go | ✅ Done |
| 4 | Fix `ExitError.Error()` nil panic: add nil guard for `e.Err` | errors.go | ✅ Done |

### P1: Type Safety Improvements

| # | Task | Files | Effort |
|---|------|-------|--------|
| 5 | Add runtime type guards in `dispatchParse`: verify return type matches handler target | type_handler.go | ✅ Done |
| 6 | ~~Fix DefaultFunc error suppression~~ → Investigated: RegisterFunc validates at registration time; safe by design | type_handler_kinds.go | ✅ Done |
| 7 | Validate `short` tag length (must be 1 char) at registration time | config_parsing.go | ✅ Done |
| 8 | Add nil check to `tracer` in `TelemetryMiddleware` | telemetry.go | ✅ Done |

### P2: Architecture Cleanup

| # | Task | Files | Effort |
|---|------|-------|--------|
| 9 | ~~Eliminate global singletons~~ → Investigated: global registries are template pattern, already cloned per-FlagRegistry; documented as intentional | type_handler.go, flags_validate.go | ✅ Done |
| 10 | Unify `fieldValueToString` and `formatFieldValue` into canonical `formatFieldValue()` in flag_helpers.go | flag_helpers.go, config_file.go, flags_validate.go | ✅ Done |
| 11 | Split `errors.go` into domain-specific files: errors_command.go, errors_flags.go, errors_config.go, errors_di.go | errors.go | ✅ Done |
| 12 | Fix `BranchingFlowContext.SetValue` overwriting child local values: skip children with local key set | flow_context.go | ✅ Done |
| 13 | Fix `Tags()` and `Path()` returning internal mutable slices: return `slices.Clone()` | flags.go, flow_context.go | ✅ Done |

### P3: Documentation & Consistency

| # | Task | Files | Effort |
|---|------|-------|--------|
| 14 | Fix `doc.go` "never panics" — Must* functions do panic | doc.go | ✅ Done |
| 15 | Align flag precedence chain across all docs: explicit flag → env → config file → default | README.md, docs/FEATURES.md | ✅ Done |
| 16 | Update `ROADMAP.md`: check off completed fuzz tests, CONTRIBUTING.md, issue/PR templates | ROADMAP.md | ✅ Done |
| 17 | Fix `WHAT_THIS_PROJECT_IS_NOT.md` line 75: config file loading IS provided | WHAT_THIS_PROJECT_IS_NOT.md | ✅ Done |
| 18 | Fix `docs/FEATURES.md`: updated API (YAML() not YAMLLoader{}), added LogFormat type, updated error file split, fixed go-output version | docs/FEATURES.md | ✅ Done |
| 19 | Update `AGENTS.md` v2.3 Design Principles title → v2 Design Principles | AGENTS.md | ✅ Done |

### P4: CI/CD

| # | Task | Files | Effort |
|---|------|-------|--------|
| 20 | Add `CODECOV_TOKEN` secret to GitHub repo settings | GitHub settings | 5m |

### P5: Future (v3.0+)

| # | Task | Category |
|---|------|----------|
| 21 | Plugin system for custom validators and type handlers | Feature |
| 22 | Config file nested struct support | Feature |
| 23 | Documentation generation (GenerateDocs, markdown, API docs) | Feature |
| 24 | Advanced types: Result[T], Validated[T], branded IDs | Feature |
| 25 | Config auto-loading with koanf integration | Feature |
| 26 | Structured JSON error output for `--output=json` | Feature |
| 27 | Extract flag-related code to standalone `flagtags` library | Refactor |
| 28 | Consider extracting `go-output` to sub-package (pay-per-use) | Refactor |

### P6: Future Cleanup (API-breaking, defer to v3.0)

| # | Task |
|---|------|
| 29 | Make `NoFlags` a distinct named type (not type alias) |
| 30 | Rename `Get[T]`/`MustGet[T]` to more specific names |
| 31 | Make `RegisterInScope` generic instead of `...any` |
| 32 | Remove or redesign `Package()` for error-safe DI integration |
| 33 | Remove `SetConfig` or make it safe (reinitialize FlagRegistry) |
| 34 | Remove deprecated `WithColor` option |
| 35 | Fix `os.Setenv("NO_COLOR", "1")` process-wide mutation in `Execute` |
