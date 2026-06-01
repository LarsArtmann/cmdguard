# DELETE TO IMPROVE PLAN — cmdguard

> Deep analysis of everything in this repo that should be deleted, archived, or cleaned up to make the codebase leaner, more honest, and more maintainable.
>
> **Generated:** 2026-06-01
> **Status:** v2.3.0-dev — 333+ tests, ~84% coverage, 0 lint, 0 races

---

## Executive Summary

This repo has **significant doc bloat** from months of AI-assisted development. The pattern is clear: every sprint generated status reports, execution plans, and retrospective documents that were never cleaned up. The **source code is remarkably clean** — no dead code, no TODOs, no stale tests. The problem is almost entirely in documentation and configuration artifacts.

| Category | Files to Delete | Files to Keep | Estimated Cleanup |
|----------|----------------|---------------|-------------------|
| Archived status reports | 59 | 1 | ~59 files |
| Historical execution plans | 17 | 1 | ~17 files |
| Orphaned config/artifacts | 8 | 0 | ~8 files |
| Stale root-level docs | 3 | 4 | ~3 files |
| Stale docs/ files | 7 | 10 | ~7 files |
| Redundant lint config | 1 | 0 | ~1 file |
| **Total** | **~95 files** | **~16 files** | **~95 deletions** |

**Source code deletions:** Minimal — 1 deprecated method (`IsExecutable()`), 0 dead test files, 0 broken tests.

---

## Category 1: Status Report Graveyards (DELETE ALL)

These are AI session handoff documents. They were useful for the session that produced them and have zero ongoing reference value. Every item described in them is either completed or tracked in `TODO_LIST.md`.

### `docs/archive/status/` — 29 files

All from March-April 2026 (v2.1 era). Already archived once. Now just taking up space.

```
docs/archive/status/2026-03-22_14-10_STATUSReport-v2.1.md
docs/archive/status/2026-03-22_14-12_comprehensive-status.md
docs/archive/status/2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-23_20-51_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-23_21-47_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-28_00-05_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-28_00-15_EXECUTION_PLAN.md
docs/archive/status/2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-03-28_08-56_COMPREHENSIVE_EXECUTIVE_STATUS.md
docs/archive/status/2026-03-28_09-54_COMPREHENSIVE_STATUS.md
docs/archive/status/2026-03-28_14-47_FINAL_STATUS.md
docs/archive/status/2026-04-01_03-04_COMPREHENSIVE_EXECUTIVE_STATUS.md
docs/archive/status/2026-04-02_09-26_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-04-05_03-21_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-04-08_21-15_COMPREHENSIVE_STATUS_REPORT.md
docs/archive/status/2026-04-09_06-20_comprehensive-audit.md
docs/archive/status/2026-04-09_11-01_audit-and-fixes.md
docs/archive/status/2026-04-09_session2-continuation.md
docs/archive/status/2026-04-10_09-37_cross-library-integration-review.md
docs/archive/status/2026-04-10_09-58_test-fix-and-project-health.md
docs/archive/status/2026-04-10_10-05_comprehensive-execution-plan.md
docs/archive/status/2026-04-10_10-53_build-fix-and-health-check.md
docs/archive/status/2026-04-10_12-12_comprehensive-analysis-and-execution-plan.md
docs/archive/status/2026-04-12_01-05_session-continuation-validation-example.md
docs/archive/status/2026-04-14_20-38_long-description-enforcement.md
docs/archive/status/2026-04-15_23-31_comprehensive-audit-status.md
docs/archive/status/2026-04-15_23-36_comprehensive-status.md
docs/archive/status/2026-04-16_04-49_deduplication-sprint-status.md
docs/archive/status/2026-04-16_06-39_deduplication-sprint-status-2.md
```

**Action:** Delete entire `docs/archive/` directory (includes the parent `MIGRATION_V1_TO_V2.md` which is superseded by `docs/MIGRATION_FROM_COBRA.md`).

### `docs/status/archive/` — 30 files

All from April-June 2026. Every sprint completed, every item tracked in TODO_LIST.md.

```
docs/status/archive/2026-04-17_19-47_comprehensive-audit-sprint-final.md
docs/status/archive/2026-04-19_13-57_post-v1-removal-comprehensive-retrospective.md
docs/status/archive/2026-04-30_02-07_v2.2.0-SUPERB-sprint-status.md
docs/status/archive/2026-04-30_03-57_comprehensive-status-report.md
docs/status/archive/2026-04-30_06-36_architecture-hardening-complete.md
docs/status/archive/2026-04-30_06-57_comprehensive-post-hardening-status.md
docs/status/archive/2026-05-16_21-38_comprehensive-status-report.md
docs/status/archive/2026-05-16_22-29_mid-session-honest-status.md
docs/status/archive/2026-05-16_23-14_comprehensive-session-resume-status.md
docs/status/archive/2026-05-16_23-23_SUPERB-architecture-hardening-status.md
docs/status/archive/2026-05-17_05-02_comprehensive-status-report.md
docs/status/archive/2026-05-17_06-09_public-release-comprehensive-status.md
docs/status/archive/2026-05-17_06-33_comprehensive-status-report.md
docs/status/archive/2026-05-17_06-35_post-refactor-polish-status.md
docs/status/archive/2026-05-17_18-01_Go-1.26-migration-remove-Ptr-helper.md
docs/status/archive/2026-05-17_smart-consolidations-status.md
docs/status/archive/2026-05-18_15-07_post-orphan-cleanup-comprehensive-status.md
docs/status/archive/2026-05-21_16-49_zero-clones-error-wrapping-dedup.md
docs/status/archive/2026-05-23_02-37_strong-id-analysis-session.md
docs/status/archive/2026-05-27_01-43_config-file-loading-complete.md
docs/status/archive/2026-05-27_02-12_final-config-verification.md
docs/status/archive/2026-05-27_02-24_go-mod-workspace-fix-comprehensive.md
docs/status/archive/2026-05-27_02-41_consumer-perspective-audit-resolution-status.md
docs/status/archive/2026-05-27_03-04_consumer-audit-complete-comprehensive-status.md
docs/status/archive/2026-05-31_21-02_interactive-prompts-huh-complete.md
docs/status/archive/2026-05-31_21-23_Phase-9-architecture-hardening-complete.md
docs/status/archive/2026-05-31_22-34_full-status-report.md
docs/status/archive/2026-05-31_23-37_spinner-glamour-telemetry-polish.md
docs/status/archive/2026-06-01_08-52_polish-and-harden.md
docs/status/archive/2026-06-01_09-27_full-status-report.md
```

**Action:** Delete entire `docs/status/archive/` directory.

### `docs/status/` (active) — keep 1, delete 1

- **Keep:** `2026-06-01_12-07_comprehensive-status.md` (current project state)
- **Delete:** `2026-06-01_09-44_comprehensive-cleanup-and-release-prep.md` (superseded same day by the 12:07 report)

---

## Category 2: Historical Execution Plans (DELETE ALL)

These are sprint plans from Feb-May 2026. Every one is fully executed. Keeping them creates noise and makes it hard to find actual planning documents.

### `docs/planning/` — delete 17, keep 1

| File | Date | Status | Action |
|------|------|--------|--------|
| `2026-02-20_07-02-PARETO_EXECUTION_MASTERPLAN.md` | Feb 20 | Fully executed | **DELETE** |
| `2026-02-20_COMPREHENSIVE_EXECUTION_PLAN.md` | Feb 20 | Fully executed | **DELETE** |
| `2026-03-22_api-redesign-v2.1.md` | Mar 22 | Partially executed, evolved past | **DELETE** |
| `2026-03-22_v2.1-minimal-plan.md` | Mar 22 | Fully executed | **DELETE** |
| `2026-03-22_v3.0-major-redesign-plan.md` | Mar 22 | Not started, misleading | **DELETE** |
| `2026-04-17_post-refactor-audit-polish.md` | Apr 17 | Fully executed | **DELETE** |
| `2026-04-17_brutally-honest-retrospective-and-execution-plan.md` | Apr 17 | Fully executed | **DELETE** |
| `2026-04-18_post-v1-removal-hardening.md` | Apr 18 | Fully executed | **DELETE** |
| `2026-04-30_SUPERB_CLI_TUI_ROADMAP.md` | Apr 30 | Fully executed | **DELETE** |
| `2026-04-30_micro-task-plan.md` | Apr 30 | Mostly executed | **DELETE** |
| `2026-05-16_v2.3-architecture-hardening.md` | May 16 | Fully executed | **DELETE** |
| `2026-05-16_SUPERB-architecture-hardening.md` | May 16 | Fully executed | **DELETE** |
| `2026-05-16_honest-execution-plan.md` | May 16 | Fully executed | **DELETE** |
| `2026-05-16_smart-consolidations.md` | May 16 | Fully executed | **DELETE** |
| `2026-05-17_global-state-elimination-and-coverage.md` | May 17 | Fully executed | **DELETE** |
| `go-composable-business-types-usage.md` | Mar 17 | Never integrated | **DELETE** |
| `examples-consolidation.html` | ? | Generated HTML artifact | **DELETE** |
| `taskctl-improvements.html` | ? | Generated HTML artifact | **DELETE** |
| `COMPREHENSIVE_TODO_PLAN.md` | Jun 1 | Active (33 items, ~4 done) | **KEEP** |

**Action:** Delete 17 files + 1 directory if emptied. Keep only `COMPREHENSIVE_TODO_PLAN.md`.

---

## Category 3: Orphaned Config & Tool Artifacts (DELETE ALL)

These are one-off tool outputs, broken scripts, and local configs that were committed accidentally or are no longer used.

| File | Why Delete |
|------|-----------|
| `.config/metadata.yaml` | Auto-generated metadata, zero references |
| `.auto-deduplicate/false-positives.json` | One-off dedup tool config, references old macOS paths, zero references |
| `.auto-deduplicate/` (directory) | Empty after deleting above |
| `library-policy.yaml` | Local tool config, zero references in CI/scripts/docs |
| `scripts/pre-commit` | Broken for months. AGENTS.md says `--no-verify` required. CI covers all checks. |
| `scripts/` (directory) | Empty after deleting above |
| `report/jscpd-report.json` | Generated jscpd output, zero references, regenerable |
| `report/` (directory) | Empty after deleting above |
| `docs/dedup-report.html` | Generated jscpd HTML report, zero references |
| `docs/reviews/2026-02-20_samber_do_v2_review.md` | All issues fixed, review complete, no actionable content |
| `docs/reviews/` (directory) | Empty after deleting above |
| `pkg/cmdguard/.golangci.yml` | Redundant with root `.golangci.yml` (which is more complete), never loaded by golangci-lint |
| `docs/planning/examples-consolidation.html` | Standalone HTML artifact, zero references |
| `docs/planning/taskctl-improvements.html` | Standalone HTML artifact, zero references |

**Action:** Delete all 12 files and their now-empty parent directories.

---

## Category 4: Stale Root-Level Documentation (DELETE 3)

| File | Why Delete |
|------|-----------|
| `PROGRESS_2026-04-01.md` | Historical snapshot of one day's work. All items tracked elsewhere. |
| `PROJECT_SPLIT_EXECUTIVE_REPORT.md` | Decision made ("don't split"). No future action. Superseded by `docs/COMPARISON.md`. |
| `PARTS.md` | Component extraction analysis. Decision was "don't extract". Superseded by `docs/COMPARISON.md` and `docs/modularization/`. |
| `BDD_TESTS_REVIEW.md` | Stale review asking "use Ginkgo or update AGENTS.md?". Question still unresolved but document is not actionable in its current form. |

| File | Why Keep |
|------|----------|
| `CONSUMER_PERSPECTIVE.md` | **Essential.** Most actionable document for improving adoption. Lists real gaps. |
| `WHAT_THIS_PROJECT_IS_ABOUT.md` | Useful orientation for contributors. Consider merging into README later. |
| `WHAT_THIS_PROJECT_IS_NOT.md` | Prevents scope creep. Short and valuable. |

**Action:** Delete 4 root-level files.

---

## Category 5: Stale Docs Files (DELETE 7)

| File | Why Delete |
|------|-----------|
| `docs/ARCHITECTURE_REVIEW.md` | Reviews v1 architecture only (pre-v2). Superseded by current code. |
| `docs/API_DESIGN_REVIEW.md` | 70% complete, partially stale. Pattern recommendations already implemented. |
| `docs/EXECUTION_PLAN_2026-04-01.md` | Historical execution plan, items done or tracked elsewhere. |
| `docs/EXECUTION_PLAN_2026-04-01_REFINED.md` | Same date, refined version. Both are historical. |
| `docs/DOMAIN_LANGUAGE.md` | Completely unfilled template. No domain terms defined. Pure noise. |
| `docs/architecture_detailed.d2` | Describes v1 internals. Superseded by `architecture.d2`. |
| `docs/architecture_detailed.svg` | Render of the v1 D2 diagram. Unreferenced. |

| File | Why Keep |
|------|----------|
| `docs/COMPARISON.md` | Excellent consumer-facing comparison with alternatives. |
| `docs/QUICKSTART.md` | Core user-facing documentation. |
| `docs/TUTORIAL.md` | Step-by-step onboarding. |
| `docs/FEATURES.md` | Authoritative feature reference. |
| `docs/MIGRATION_FROM_COBRA.md` | Addresses top adoption gap. |
| `docs/CLI_DESIGN_PRINCIPLES.md` | Codified design philosophy. |
| `docs/PERFORMANCE.md` | Benchmark results, consumer-facing. |
| `docs/SUPERB_CLI_GAP_ANALYSIS.md` | Active feature backlog. |
| `docs/architecture.d2` | v2 architecture diagram (needs minor `GuardedCommand`→`CLI` rename). |
| `docs/architecture.svg` | Render of the D2 diagram (unreferenced but keep with d2). |
| `docs/DOMAIN_LANGUAGE.md` | Not yet — see deletion above. |

### Keep but Update

| File | What Needs Fixing |
|------|-------------------|
| `docs/architecture.d2` | References `GuardedCommand[T,F]` — rename to `CLI[T]` |
| `FEATURES.md` (root) | Line 192 lists `Ptr[T]` as `✅ FULLY_FUNCTIONAL` but it was removed in v2.3.0. Mark as removed or delete row. |
| `ROADMAP.md` | "Fuzz Testing" and "Spinner" sections show unchecked items that are already done. Mark as `[x]`. |

---

## Category 6: Source Code Cleanups (Minimal)

The source code is remarkably clean. Almost nothing to delete.

### Dead/Deprecated Code

| Item | Status | Action |
|------|--------|--------|
| `IsExecutable()` method | Deprecated, delegates to `HasHandler()`. Zero callers outside tests/docs. | **Delete in v3** (breaking change). Mark deprecated now if not already. |
| `WithColor` option | Deprecated in favor of `WithFang`. Still used in 18 test call sites. | **Delete in v3.** Migrate tests to `WithFang` first. |

### NOT Dead Code (confirmed alive)

These were checked and confirmed to be real, tested, wired-up code:

- `completion.go` — `WithCompletion`/`WithValidArgs` wired into cobra, used in examples. (But: **zero test coverage** — needs tests, not deletion.)
- `fuzz_test.go` — 7 meaningful fuzz targets with good seed corpora.
- `coverage_test.go` — Tests real exported APIs. Poorly named but legitimate.
- `cli_output.go` — Fully wired into CLI lifecycle, tested in `cli_output_test.go`.
- `doc.go` — Accurate package documentation.

### Test File Consolidation Opportunity

Three test helper files could be merged (not urgent, not deletion):

- `test_helpers_test.go`
- `helpers_test.go`
- `testhelpers_test.go`

---

## Category 7: Keep As-Is (No Changes Needed)

These files/directories are current, accurate, and valuable:

| File/Dir | Why Keep |
|----------|----------|
| `SECURITY.md` | Current, version-specific security policy |
| `CODE_OF_CONDUCT.md` | Standard community file |
| `CONTRIBUTING.md` | Detailed, accurate contributor guide |
| `README.md` | User-facing documentation |
| `CHANGELOG.md` | Up to date through v2.3.0 |
| `TODO_LIST.md` | Current, honest about remaining work |
| `FEATURES.md` | Mostly accurate (fix `Ptr[T]` row) |
| `ROADMAP.md` | Mostly current (mark done items) |
| `AGENTS.md` | Comprehensive AI context |
| `AUTHORS` | Standard open-source file |
| `LICENSE` | Required |
| `flake.nix` / `flake.lock` | Build system |
| `go.mod` / `go.sum` | Dependencies |
| `.golangci.yml` (root) | Authoritative lint config |
| `.github/workflows/` | Functional CI/CD |
| `.github/ISSUE_TEMPLATE/` | Functional issue templates |
| `.github/PULL_REQUEST_TEMPLATE.md` | Functional PR template |
| `.gitattributes` / `.gitignore` | Git config |
| `git-town.toml` | Personal tool config (keep if using git-town) |
| `benchmarks/` | Active v2 benchmarks, used by CI |
| `examples/` | Excellent, complete example app |
| `tests/integration/` | Meaningful v2 integration tests |
| `docs/modularization/` | Unexecuted but valuable proposal |
| `docs/design/` | Foundational v2 architecture vision |
| `docs/QUICKSTART.md` | Core user documentation |
| `docs/TUTORIAL.md` | Step-by-step guide |
| `docs/COMPARISON.md` | Consumer-facing comparison |
| `pkg/cmdguard/v2/testutil/` | Consumer test harness |

---

## Deletion Priority Matrix

### Priority 1: Immediate (zero risk, pure cleanup)

Delete these right now. No breaking changes, no lost information.

1. **`docs/archive/`** — entire directory (29 status reports + 1 stale migration guide)
2. **`docs/status/archive/`** — entire directory (30 status reports)
3. **`docs/status/2026-06-01_09-44_comprehensive-cleanup-and-release-prep.md`** — same-day duplicate
4. **`.auto-deduplicate/`** — one-off tool config
5. **`report/`** — generated jscpd report
6. **`scripts/`** — broken pre-commit hook
7. **`docs/dedup-report.html`** — generated artifact
8. **`docs/reviews/`** — completed review
9. **`pkg/cmdguard/.golangci.yml`** — redundant lint config
10. **`docs/planning/` HTML files** — orphaned artifacts

**Total: ~63 files, ~10 empty directories**

### Priority 2: Near-term (low risk, historical docs)

Delete after confirming no references from external sources.

1. **`docs/planning/` execution plans** — 15 completed plans (keep `COMPREHENSIVE_TODO_PLAN.md`)
2. **`PROGRESS_2026-04-01.md`** — historical snapshot
3. **`PROJECT_SPLIT_EXECUTIVE_REPORT.md`** — decision made
4. **`PARTS.md`** — superseded
5. **`BDD_TESTS_REVIEW.md`** — stale
6. **`docs/ARCHITECTURE_REVIEW.md`** — v1 only
7. **`docs/API_DESIGN_REVIEW.md`** — partially stale
8. **`docs/EXECUTION_PLAN_2026-04-01.md`** + `_REFINED.md` — historical
9. **`docs/DOMAIN_LANGUAGE.md`** — unfilled template
10. **`docs/architecture_detailed.d2`** + `.svg` — v1 diagram

**Total: ~24 files**

### Priority 3: Future (breaking changes, deferred to v3)

1. **`IsExecutable()`** — deprecated method, remove in v3
2. **`WithColor`** — deprecated option, migrate tests then remove in v3
3. **`WithColor` test references** — 18 test call sites to migrate

### Priority 4: Fix, Don't Delete

| File | Fix |
|------|-----|
| `FEATURES.md` line 192 | `Ptr[T]` row: change from `✅ FULLY_FUNCTIONAL` to `🗑️ REMOVED (v2.3.0)` |
| `ROADMAP.md` | Mark fuzz tests, spinner, telemetry as `[x]` done |
| `docs/architecture.d2` | Rename `GuardedCommand[T,F]` → `CLI[T]` |
| `docs/architecture.svg` | Re-render after D2 fix |

---

## Impact Analysis

### What We Lose (Nothing Valuable)

- 59 status reports → All items tracked in `TODO_LIST.md` or already done
- 17 execution plans → All fully executed, conclusions reflected in code
- 12 orphaned artifacts → Zero references, regenerable or broken
- 7 stale docs → Superseded by current docs or describe v1

### What We Gain

- **Clarity:** `docs/` goes from ~100 files to ~15 meaningful ones
- **Discoverability:** No more wading through "COMPREHENSIVE_STATUS_REPORT" noise
- **Honesty:** No stale templates (`DOMAIN_LANGUAGE.md`) or broken scripts (`scripts/pre-commit`)
- **Reduced maintenance:** Fewer files to keep in sync with code changes
- **Faster onboarding:** New contributors see actual documentation, not session transcripts

### Risk Assessment

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Lose historical context for a decision | Low | Key decisions documented in AGENTS.md, CHANGELOG.md, and kept docs |
| Break a reference from external source | Very Low | All deleted files are internal, not user-facing |
| Need to redo a status report | Zero | Status reports are session-scoped by definition |
| Regret deleting a planning doc | Low | Git history preserves everything |

---

## Summary: The Final Tally

| Action | Count |
|--------|-------|
| **Files to delete** | ~95 |
| **Directories to remove** | ~10 (after emptying) |
| **Files to fix (not delete)** | 3 |
| **Source code to delete** | 0 (v2.3.0; 2 items deferred to v3) |
| **Files to keep unchanged** | ~40 |
| **Estimated disk savings** | ~500KB of markdown/HTML noise |

The repo will go from **~140 doc/config files** to **~45 meaningful ones**. The signal-to-noise ratio improves dramatically.

---

*"Perfection is achieved not when there is nothing more to add, but when there is nothing left to take away."* — Antoine de Saint-Exupéry
