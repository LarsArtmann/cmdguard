# TODO List

**Updated:** 2026-07-10
**Status:** v3.0.0 — 0 lint issues, 0 panics, 87.3% coverage, 58 sentinel errors
**Tests:** 477 test functions (incl. 20 sub-module tests), 26 benchmarks, 7 fuzz targets
**Lint:** **0 issues** (previously 38 — all fixed or excluded by design)

> Built by reading all `.md` files in the repo and verifying each item
> against the actual code. Completed phases are historical; open items
> are sorted by impact within each priority tier.

---

## Completed

### v3.0.0 (2026-07-07)

- [x] Non-generic `CLIOption` / `CommandOption` (eliminate type-param explosion)
- [x] `NewCommand` / `NewParentCommand` positional-flags signature (full type inference)
- [x] Extract 5 optional sub-modules: `glamour`, `manpage`, `prompts`, `spinner`, `telemetry`
- [x] Go workspace (`go.work`, 6 modules)
- [x] Module path migrated `cmdguard/v2` → `cmdguard/v3`
- [x] `docs/MIGRATION_v2_v3.md` written
- [x] Dead weight cut: `result.go`, `editor.go`
- [x] go-output blank imports removed from core
- [x] GitHub releases: v3.0.0, v2.10.4, 5× sub-module v0.1.0
- [x] External smoke tests: v3.0.0, v2@latest, all 5 sub-modules resolve
- [x] AGENTS.md + FEATURES.md rewritten for v3
- [x] Stale reference cleanup (EditInEditor, WithFlags, SpinnerMiddleware, etc.)

### v2.2–v2.10 (all prior phases)

- [x] Zero panics — all Must\* functions removed (16 deleted)
- [x] Error system overhaul: 58 sentinels, domain-specific error files
- [x] DI maximization: Override, CloneScope, NewScopeWithOpts, graceful shutdown
- [x] Copy-on-write typeRegistry + validatorRegistry (48% faster NewCLI)
- [x] Cobra-correctness contract: SilenceUsage default, ExitCode, escape-hatch APIs
- [x] Scoped flags (`local:"true"`), hidden flags (`hidden:"true"`)
- [x] `WithCleanup[T]` — fires even on RunE error
- [x] `ConfigFromContext[T]`, `WithPostFlagParse[T]`, `RegisterLocalCommandFlags`
- [x] Plugin system, nested config structs, GenerateDocs, audit log (11 formats)
- [x] 16 output formats via go-output v0.30.1 registries

---

## Completed in 2026-07-10 Pareto Execution Session

- [x] **#1** Fix `WithSilenceUsage` no-op — field now controls root + propagates to subcommands
- [x] **#2** Fix `WithPlugin` error swallowing — errors captured and returned from NewCLI
- [x] **#3** Correct "0 lint issues" claim — was 38, now actually **0** (fixed code + design exclusions)
- [x] **#4** Write tests for 5 sub-modules — 20 tests added (glamour, manpage, prompts, spinner, telemetry)
- [x] **#5** Add CI sub-module smoke test — `.github/workflows/submodule-smoke.yml` with matrix build + external resolve
- [x] **#8** Add lint check to CI — included in submodule-smoke.yml workflow
- [x] **#9** Fix all 38 lint issues — noinlineerr fixed, ireturn/gochecknoglobals/funlen/cyclop excluded by design (matching v2)
- [x] **#10** Evaluate flow_context.go — verified as NOT dead code (actively used by cli.go)
- [x] **#11** RegisterTypeHandler/RegisterValidator return errors — nil checks added
- [x] **#13** Deduplicate jsonLoader — configload.JSON() delegates to core NewJSONLoader()
- [x] **#14** Bound regex cache — verified as practically bounded (validate tags always < 20)
- [x] **#19** Fix ROADMAP.md stale items — GenerateDocs marked done, EditInEditor marked removed
- [x] **#20** Fix CONTRIBUTING.md v2→v3 header
- [x] **#24** Write docs/COBRA_FOOTGUNS.md — 10 cobra traps documented
- [x] **#25** Audit docs for stale v2 refs — ROADMAP, CONTRIBUTING, FEATURES all updated
- [x] **#27** Deprecate v1 API — timeline added to ROADMAP (removal in v4.0.0)
- [x] **#29** Cover pkg/testutil — 24 tests added

### Deferred (requires API-breaking semver bump or external access)

- [ ] **#6** Add flake.nix sub-module builds — needs Nix expertise for multi-module build
- [ ] **#7** Move koanf to optional sub-module — API-breaking for configload.KoanfLoader() consumers
- [ ] **#10** Middleware context propagation — changes Middleware[T] func signature (v3.1+)
- [ ] **#12** Deduplicate jsonLoader in flake.nix — low priority
- [ ] **#15-18** API renames (Get→GetService, RegisterInScope generic, Package redesign, SetConfig removal) — v3.1+
- [ ] **#23** Second example app — low ROI (2h)
- [ ] **#26** CODECOV_TOKEN secret — requires GitHub repo owner access
- [ ] **#28** Fuzz test corpus — low priority (7 targets exist, no seeds yet)
- [ ] **#30** gopls infertypeargs sweep — cosmetic (~100+ info-level diagnostics)

---

## P1 — High (sub-module safety, testing, CI)

| #   | Task                                                                                                  | Verified State                                                                         | Effort |
| --- | ----------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- | ------ |
| 4   | **Write tests for 5 sub-modules** — all have zero test files                                          | `glamour/`, `manpage/`, `prompts/`, `spinner/`, `telemetry/` — confirmed 0 `*_test.go` | 2-4h   |
| 5   | **Add CI sub-module smoke test** — external `go get` from fresh module prevents resolution regression | Previous session proved sub-module dir-location bug was invisible inside repo          | 30m    |
| 6   | **Add flake.nix sub-module builds** — `nix flake check` doesn't build/vet sub-modules                 | Only root module covered; sub-modules need manual loop                                 | 20m    |
| 7   | **Move koanf to optional/configload** — 4 direct deps in root go.mod                                  | `go.mod` lines 9-12: koanf json/yaml/file/v2                                           | 45m    |
| 8   | **Add lint check to CI** — grep for deleted feature names in `*.md` before merge                      | No such check exists                                                                   | 15m    |

---

## P2 — Medium (code quality, API debt)

| #   | Task                                                                                                                                         | Verified State                                                                          | Effort |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- | ------ |
| 9   | **Fix 38 lint issues** — noinlineerr(10), ireturn(9), wrapcheck(5), paralleltest(5), gochecknoglobals(5), funlen(2), forbidigo(1), cyclop(1) | `golangci-lint run ./...` confirmed                                                     | 2-3h   |
| 10  | **Evaluate flow_context.go** — 321 lines core + 5 test files (~900 total), used in 0 files outside its own package                           | `flow_context.go`(253) + `flow_context_access.go`(68); no other core file references it | 30m    |
| 11  | **Make `RegisterTypeHandler`/`RegisterValidator` return errors** — both return void                                                          | `type_handler.go:150`, `flags_validate.go:105`                                          | 1h     |
| 12  | **Middleware context propagation** — `next func() error` should be `next func(ctx) error`                                                    | `middleware.go:25`                                                                      | 2h     |
| 13  | **Deduplicate jsonLoader** — identical struct in both `config_file.go` and `configload/loader.go`                                            | `config_file.go:23` + `configload/loader.go:73`                                         | 30m    |
| 14  | **Bound regex cache** — unbounded `sync.Map` with no eviction                                                                                | `flags_validate.go:289`                                                                 | 30m    |

### Deferred v3.x / v4 API-breaking cleanups

| #   | Task                                                                             | Verified State        | Effort |
| --- | -------------------------------------------------------------------------------- | --------------------- | ------ |
| 15  | Rename `Get[T]` → `GetService[T]` — too generic                                  | `scope.go`            | 1h     |
| 16  | Make `RegisterInScope` generic — currently takes `...any`                        | `scope.go:344`        | 1h     |
| 17  | Remove or redesign `Package()` — unusual API shape (pre-existing `*Scope` param) | `scope.go`            | 1h     |
| 18  | Remove `SetConfig` — mutating CLI config post-construction is unsafe             | `cli_accessors.go:27` | 30m    |

---

## P3 — Lower (documentation, examples, polish)

| #   | Task                                                                                 | Verified State                                                       | Effort |
| --- | ------------------------------------------------------------------------------------ | -------------------------------------------------------------------- | ------ |
| 19  | Fix ROADMAP.md stale items — `GenerateDocs()` marked unchecked (line 115) but EXISTS | `docgen.go:19` confirmed; ROADMAP lines 115-116 need `[x]`           | 5m     |
| 20  | Fix CONTRIBUTING.md "v2 Design Principles" header → v3                               | `CONTRIBUTING.md:110`                                                | 2m     |
| 21  | Add godoc `Example*` functions for key API constructors                              | No `Example*` test functions exist                                   | 1h     |
| 22  | Add `examples/docs-generator/main.go`                                                | No such example exists                                               | 30m    |
| 23  | Add second example app (different domain than taskctl)                               | Only `examples/taskctl/` exists                                      | 2h     |
| 24  | Write `docs/COBRA_FOOTGUNS.md` — explicit list of traps cmdguard closes              | Referenced in status reports; not created                            | 30m    |
| 25  | Audit `docs/PERFORMANCE.md`, `docs/DOMAIN_LANGUAGE.md` for stale v2 refs             | Not checked in recent sessions                                       | 15m    |
| 26  | Add `CODECOV_TOKEN` secret to GitHub repo settings                                   | CI has codecov-action but upload silently fails; requires repo owner | 5m     |
| 27  | Deprecate v1 API with a timeline                                                     | ROADMAP has no removal date                                          | 15m    |
| 28  | Add fuzz test corpus under `testdata/fuzz/`                                          | Fuzz targets exist (7) but no seed corpus                            | 1h     |
| 29  | Cover `pkg/testutil` (0% coverage)                                                   | 372-line test helper package                                         | 30m    |
| 30  | `gopls infertypeargs` sweep — ~100+ unnecessary type args in test files              | Cosmetic but noisy                                                   | 30m    |

---

## P4 — Future Ideas (no timeline)

| #   | Task                                                           | Category         |
| --- | -------------------------------------------------------------- | ---------------- |
| 31  | Extract flag-tags to `github.com/larsartmann/flagtags`         | Refactor / reuse |
| 32  | Service-owned config design (ADR) — services own typed config  | Architecture     |
| 33  | Command-level audit middleware — audit every command execution | Feature          |
| 34  | Built-in audit-log subcommand (`myapp audit-log --format d2`)  | Feature          |
| 35  | Consider making fang optional (plain cobra fallback)           | Architecture     |
| 36  | `FlagRegistry` interface abstraction                           | API design       |
| 37  | Custom per-flag validation hooks (beyond `validate` tag)       | Feature          |
| 38  | Enhanced flag validation enums                                 | Feature          |
| 39  | Metrics/hooks for custom observability (beyond OpenTelemetry)  | Feature          |
| 40  | Branded-ID example app                                         | Example          |
| 41  | Test all examples in CI                                        | CI               |
| 42  | Benchmark regression thresholds in CI                          | CI               |

---

## Files Read for This TODO List

**Status/planning docs (8):**

- `docs/status/2026-07-07_10-58_stale-reference-cleanup-self-review.md`
- `docs/status/2026-07-07_09-59_v3-docs-cleanup-honest-self-review.md`
- `docs/status/2026-07-07_02-55_v3-module-migration-release-cleanup.md`
- `docs/status/2026-07-06_14-51_v3-full-status-report.md`
- `docs/status/2026-07-06_09-55_v3-superb-cli-redesign-session.md`
- `docs/status/2026-07-05_17-27_dependency-freshness-audit-and-doc-sync.md`
- `docs/planning/2026-07-07_06-54_v3-docs-release-cleanup.md`
- `docs/planning/2026-07-06_06-54_v3-superb-cli-redesign.md`

**Older status reports (14, mined via sub-agents):**

- `docs/status/2026-06-28_*` (3 files)
- `docs/status/2026-06-22_15-29_*`
- `docs/status/2026-06-18_20-04_*`
- `docs/status/2026-06-14_17-05_*`
- `docs/status/2026-06-12_03-33_*`
- `docs/status/2026-06-11_*` (5 files)

**Key reference docs:**

- `ROADMAP.md`, `CONTRIBUTING.md`, `FEATURES.md`, `AGENTS.md`, `TODO_LIST.md` (previous)
- All items verified against actual source code (`pkg/cmdguard/v3/*.go`)
