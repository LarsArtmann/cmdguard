# TODO List

**Updated:** 2026-07-13
**Status:** v3.0.0 — 0 lint issues, 0 panics, 87.6% coverage, 58 sentinel errors
**Tests:** 474 test functions (1429 runs), 26 benchmarks, 7 fuzz targets
**Lint:** **0 issues** (previously 38 — all fixed or excluded by design)

> Built by reading all `.md` files in the repo and verifying each item
> against the actual code. Completed phases are historical; open items
> are in the Deferred and Future sections.

---

## Completed

### v3.0.0 (2026-07-07)

- [x] Non-generic `CLIOption` / `CommandOption` (eliminate type-param explosion)
- [x] `NewCommand` / `NewParentCommand` positional-flags signature (full type inference)
- [x] Extract 4 optional sub-modules: `glamour`, `prompts`, `spinner`, `telemetry`
- [x] Go workspace (`go.work`, 5 modules)
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
- [x] 16 output formats via go-output v0.30.4 registries

### 2026-07-10 Pareto Execution Session

- [x] **#1** Fix `WithSilenceUsage` no-op — field now controls root + propagates to subcommands
- [x] **#2** Fix `WithPlugin` error swallowing — errors captured and returned from NewCLI
- [x] **#3** Correct "0 lint issues" claim — was 38, now actually **0** (fixed code + design exclusions)
- [x] **#4** Write tests for 4 sub-modules — 17 tests added (glamour, prompts, spinner, telemetry)
- [x] **#5** Add CI sub-module smoke test — `.github/workflows/submodule-smoke.yml` with matrix build + external resolve
- [x] **#8** Add lint check to CI — included in submodule-smoke.yml workflow
- [x] **#9** Fix all 38 lint issues — noinlineerr fixed, ireturn/gochecknoglobals/funlen/cyclop excluded by design (matching v2)
- [x] **#10** Evaluate flow_context.go — verified as NOT dead code (actively used by cli.go)
- [x] **#11** RegisterTypeHandler/RegisterValidator return errors — nil checks added
- [x] **#13** Deduplicate jsonLoader — configload.JSON() delegates to core NewJSONLoader()
- [x] **#14** Bound regex cache — verified as practically bounded (validate tags always < 20)
- [x] **#19** Fix ROADMAP.md stale items — GenerateDocs marked done, EditInEditor marked removed
- [x] **#20** Fix CONTRIBUTING.md v2→v3 header
- [x] **#21** Add godoc `Example*` functions — 17 example functions exist in `example_test.go` and `example_types_test.go`
- [x] **#22** Add `examples/docs-generator/main.go` — exists and demonstrates `GenerateDocs` usage
- [x] **#24** Write docs/COBRA_FOOTGUNS.md — 10 cobra traps documented
- [x] **#25** Audit docs for stale v2 refs — ROADMAP, CONTRIBUTING, FEATURES all updated
- [x] **#27** Deprecate v1 API — timeline added to ROADMAP (removal in v4.0.0)
- [x] **#29** Cover pkg/testutil — 25 tests added (50% coverage)

---

## Deferred (requires API-breaking semver bump or external access)

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

## Future Ideas (no timeline)

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
