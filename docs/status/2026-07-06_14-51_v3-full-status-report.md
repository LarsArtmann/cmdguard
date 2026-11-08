# v3 Superb CLI Redesign — Full Status Report

**Date:** 2026-07-06 14:51
**Branch:** master (pushed to remote)
**Commits this session:** 11 (4f9b0ea..9297257)
**All tests:** GREEN (5 packages, 1830 test runs, race detection)
**Coverage:** 87.3%
**Lint:** 0 issues
**Source lines:** 8,110 (non-test)
**Core direct deps:** 13 (was 30)

---

> **Update 2026-07-23:** This report captured the v3 redesign just before release. The remaining "NOT STARTED" doc/release items (README/AGENTS refresh, `docs/MIGRATION_v2_v3.md`, sub-module tests, GitHub Releases, lint cleanup) were completed in the sessions following 2026-07-06 (`83f6602` through `a9c8e82`). The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## a) FULLY DONE

### 1. Type Parameter Explosion Eliminated (Both Command AND CLI Options)

**Commits:** `4f9b0ea`, `e2b125c`, `9297257`

- `CommandOption` is non-generic — all metadata options (`WithShort`, `WithLong`, `WithNoArgs`, etc.) need zero type params
- `CLIOption` is non-generic — all 22+ CLI options (`WithCLIVersion`, `WithEnvPrefix`, `WithStrictValidation`, etc.) need zero type params
- `NewCommand` takes `flags F` as positional arg → Go infers T and F automatically
- `WithFlags` option deleted entirely
- Only 4 generic-returning CLI options remain (`WithConfigValidation`, `WithMiddleware`, `WithPostFlagParse`, `WithCleanup`) — they use sealed interfaces internally

**Before:** `v2.NewCommand[AppConfig, *ListFlags]("list", handler, v2.WithShort[AppConfig, *ListFlags]("..."), v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}))`
**After:** `v2.NewCommand("list", &ListFlags{}, handler, v2.WithShort("..."))`

### 2. Lifecycle Hook Type Safety Restored

**Commit:** `428af36`

Sealed interface pattern: `lifecycleHook` + `typedHook[T,F]`. Unexported `isLifecycleHook()` method prevents external implementations. Safe type assertions — returns nil on mismatch, no panic.

### 3. Mono-Repo Modularization (5 Optional Modules Extracted)

<!-- check-rows:ledger -->
| Module    | Path                      | Deps Isolated              | Status |
| --------- | ------------------------- | -------------------------- | ------ |
| telemetry | `pkg/cmdguard/telemetry/` | OpenTelemetry SDK          | ✅     |
| manpage   | `pkg/cmdguard/manpage/`   | mango/roff                 | ✅     |
| glamour   | `pkg/cmdguard/glamour/`   | chroma/goldmark/bluemonday | ✅     |
| prompts   | `pkg/cmdguard/prompts/`   | huh/bubbles/bubbletea      | ✅     |
| spinner   | `pkg/cmdguard/spinner/`   | lipgloss (code separation) | ✅     |

All 6 modules (core + 5 optional) build independently with `GOWORK=off`.

Extension hooks added to core: `WithHelpTransform`, `PromptRunner` interface + `SetPromptRunner()`.

### 4. Dead Weight Cut

<!-- check-rows:ledger -->
| File                       | Lines | Reason                          | Status |
| -------------------------- | ----- | ------------------------------- | ------ |
| `result.go`                | 147   | Sum types — not a CLI concern   | ✅ Cut |
| `editor.go`                | 60    | $EDITOR support — marginal      | ✅ Cut |
| `telemetry.go`             | 71    | Moved to sub-module             | ✅     |
| `glamour.go`               | 102   | Moved to sub-module             | ✅     |
| `spinner.go`               | 182   | Moved to sub-module             | ✅     |
| `manpage.go`               | 65    | Moved to sub-module             | ✅     |
| 10 go-output blank imports | —     | Users import go-output directly | ✅     |

### 5. Documentation

- **doc.go** — fully updated to v3 API examples
- **CHANGELOG.md** — `[Unreleased]` section documents all breaking changes
- **Pareto plan** — `docs/planning/2026-07-06_06-54_v3-superb-cli-redesign.md`

### 6. go.work Workspace

Multi-module workspace with `use` directives for all 6 modules. Replace directives in root go.mod for all 5 sub-modules.

---

## b) PARTIALLY DONE

### 1. koanf Deps Still in Root go.mod (4 direct deps)

koanf is only used by `configload/koanf.go`, but the deps are in the root go.mod because configload is a sub-package of v2, not a separate module. Moving koanf to its own module or making it optional is the next step.

### 2. AGENTS.md — Partially Updated

Only 2 lines mention the new module structure. The full architecture section still describes the old flat structure.

### 3. README.md — Not Updated

13 references to old API patterns (`WithShort[`, `NewCommand[`, `WithFlags`). Needs comprehensive rewrite.

### 4. FEATURES.md — Not Updated

Zero mentions of the 5 new optional modules. Still describes telemetry, glamour, prompts, spinner, manpage as core features.

### 5. example/taskctl — Stale Comments

README.md and main.go header comments reference deleted features (`WithGlamourHelpTheme`, `EditInEditor`).

---

## c) NOT STARTED

1. **Move koanf to configload module or make optional** — 4 direct deps in root **→ ROADMAP deferred #5 (koanf extraction YAGNI)**
2. ~~**Cut or internalize flow_context** — 321 lines + 5 test files, used in 2 places~~ **NOT-DO — decided to keep — flow_context*.go ships as the branching-flow feature (AGENTS.md).**
3. ~~**Update README.md** with v3 quickstart and clean API examples~~ done — README rewritten through v3/v4 releases
4. ~~**Update AGENTS.md** module structure section for v3~~ done — AGENTS.md documents the 8-module workspace (2026-11-02)
5. ~~**Update FEATURES.md** with module status table~~ done — FEATURES.md SUB-MODULE rows (maintained through v4)
6. ~~**Update example/taskctl/README.md** — stale feature references~~ done — cleaned — examples maintained through v4
7. ~~**Command-level audit middleware** — audit every command execution~~ done — AuditMiddleware[T] shipped (auditlog.go; FEATURES.md)
8. **Built-in audit-log subcommand** — `myapp audit-log --format d2`
9. ~~**Deepen audit log integration** — the user's explicit "I want deeper integration"~~ done — AuditMiddleware[T] + ExportAuditLog API shipped (FEATURES.md)
10. **Service-owned config design** (ADR) — the big architectural vision from the conversation
11. ~~**Release v3.0.0-alpha tag** — current code is a breaking change, needs versioning~~ **Superseded — superseded — v3.0.0 released directly (2026-07-07).**
12. ~~**Update docs/API.md** — references old API~~ **Superseded — docs/API.md was removed in the 2026-08-06 docs refresh; the contract lives in FEATURES + the describe JSON.**
13. ~~**Update docs/QUICKSTART.md** — references old API~~ done — doc rewritten for v4 (updated through v3/v4 releases)
14. ~~**Update docs/MIGRATION_FROM_COBRA.md** — references old API~~ done — doc rewritten for v4 (updated through v3/v4 releases)

---

## d) TOTALLY FUCKED UP (Honest Self-Assessment)

### 1. Example README Still References WithFlags

`examples/taskctl/README.md` line still says `WithFlags` as if it exists. Anyone reading the example docs sees a non-existent API.

### 2. Flow Context Is Orphaned Code

321 lines of `flow_context.go` + 5 test files (total ~900 lines) for a feature that's used in exactly 2 places in core (`cli.go`, `cli_accessors.go`). It adds complexity without clear value. I should have evaluated it for cutting when I cut result.go and editor.go.

### 3. No Migration Guide

We made massive breaking changes (NewCommand signature, WithFlags deleted, CLIOption de-genericized, 5 modules extracted) and there's zero migration guide. Existing consumers are completely stranded.

### 4. lipgloss Still Direct Dep (Through fang)

Lipgloss is in root go.mod because fang pulls it. This isn't really our fault — fang depends on lipgloss. But it means our "13 direct deps" number is slightly misleading since lipgloss + its transitive tree (charmbracelet/x/ansi, etc.) still come in through fang.

### 5. WithSilenceUsage Is Now a No-Op

The `WithSilenceUsage()` option exists but does nothing — `cli.rootCmd.SilenceUsage = true` is hardcoded in NewCLI. The option should either work or be removed.

**Fixed:** v4 silences usage by default; `WithoutSilenceUsage()` opts out (FEATURES.md).

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Service-owned config** — the biggest vision from this session. Services declare their typed config through the DI graph. Flags become optional per-command input, not the primary config mechanism. This would make cmdguard genuinely unique.

2. **koanf as optional** — move to configload sub-module. Core only needs JSON (built into stdlib). **→ ROADMAP deferred #5 (koanf extraction YAGNI)**

3. ~~**Flow context evaluation** — is BranchingFlowContext earning its 900 lines? Consider cutting or making it internal.~~ **NOT-DO — decided to keep — flow_context*.go ships as the branching-flow feature (AGENTS.md).**

4. ~~**WithSilenceUsage fix** — either make it work (override the default) or remove it and document that usage is always silenced.~~ done — usage silenced by default; WithoutSilenceUsage() opts out (FEATURES.md)

5. ~~**fang dependency** — fang pulls lipgloss + mango-cobra indirectly. Consider whether fang should be optional too (users who want styled output import fang module, core uses plain cobra).~~ done — kept — fang stays a direct core dep (charm.land/fang/v2 in go.mod)

### Code Quality

6. ~~**Test coverage for new cliSpec** — the sealed interface pattern needs direct tests for type mismatch behavior (wrong T, wrong F).~~ done — typedHook exercised in cli_cobra_command_test.go

7. ~~**Integration tests for sub-modules** — verify that importing glamour/spinner/etc. actually works from external code, not just go.work.~~ done — .github/workflows/submodule-smoke.yml

8. ~~**gopls infertypeargs warnings** — 186 gopls info diagnostics about unnecessary type arguments in test files. These are cosmetic but noisy.~~ **Superseded — the v3/v4 rewrites replaced those sources.**

---

## f) Next 25 Tasks (Sorted by Impact/Effort)

| #  | Task                                                     | Impact | Effort |
| -- | -------------------------------------------------------- | ------ | ------ |
| ~~1~~  | ~~Fix example/taskctl/README.md stale WithFlags ref~~ done — cleaned — taskctl README has no WithFlags refs | ~~5~~ | ~~5min~~ |
| ~~2~~  | ~~Fix or remove WithSilenceUsage no-op~~ done — usage silenced by default; WithoutSilenceUsage() opts out (FEATURES.md) | ~~6~~ | ~~10min~~ |
| 3  | Move koanf to configload sub-module (4 deps from root) **→ ROADMAP deferred #5 (koanf extraction YAGNI)** | 7      | 30min  |
| ~~4~~  | ~~Update README.md with v3 quickstart~~ done — README rewritten through v3/v4 releases | ~~8~~ | ~~30min~~ |
| ~~5~~  | ~~Update AGENTS.md module structure for v3~~ done — AGENTS.md documents the 8-module workspace (2026-11-02) | ~~7~~ | ~~20min~~ |
| ~~6~~  | ~~Update FEATURES.md with module status~~ done — FEATURES.md SUB-MODULE rows (maintained through v4) | ~~5~~ | ~~15min~~ |
| ~~7~~  | ~~Write v3 migration guide (MIGRATION_TO_V3.md)~~ done — docs/MIGRATION_v2_v3.md exists | ~~8~~ | ~~30min~~ |
| ~~8~~  | ~~Evaluate and cut/internalize flow_context.go~~ **NOT-DO — decided to keep — flow_context*.go ships as the branching-flow feature (AGENTS.md).** | ~~5~~ | ~~30min~~ |
| ~~9~~  | ~~Add command-level audit middleware~~ done — AuditMiddleware[T] shipped (auditlog.go; FEATURES.md) | ~~8~~ | ~~60min~~ |
| 10 | Add built-in audit-log subcommand **→ ROADMAP raw ideas (built-in audit-log subcommand)** | 6      | 45min  |
| 11 | Design service-owned config ADR **→ ROADMAP raw ideas (service-owned config ADR)** | 10     | 60min  |
| ~~12~~ | ~~Update docs/API.md for v3~~ **Superseded — docs/API.md was removed in the 2026-08-06 docs refresh; the contract lives in FEATURES + the describe JSON.** | ~~4~~ | ~~20min~~ |
| ~~13~~ | ~~Update docs/QUICKSTART.md for v3~~ done — doc rewritten for v4 (updated through v3/v4 releases) | ~~4~~ | ~~15min~~ |
| ~~14~~ | ~~Clean up gopls infertypeargs warnings in tests~~ **Superseded — the v3/v4 rewrites replaced those sources.** | ~~2~~ | ~~30min~~ |
| ~~15~~ | ~~Add integration test: import glamour from external~~ done — .github/workflows/submodule-smoke.yml | ~~5~~ | ~~30min~~ |
| ~~16~~ | ~~Add integration test: import spinner from external~~ done — .github/workflows/submodule-smoke.yml | ~~4~~ | ~~15min~~ |
| ~~17~~ | ~~Consider making fang optional (plain cobra fallback)~~ done — kept — fang stays a direct core dep (charm.land/fang/v2 in go.mod) | ~~7~~ | ~~90min~~ |
| ~~18~~ | ~~Update docs/MIGRATION_FROM_COBRA.md for v3~~ done — doc rewritten for v4 (updated through v3/v4 releases) | ~~3~~ | ~~20min~~ |
| ~~19~~ | ~~Tag v3.0.0-alpha~~ **Superseded — superseded — v3.0.0 released directly (2026-07-07).** | ~~6~~ | ~~5min~~ |
| ~~20~~ | ~~Update TODO_LIST.md with v3 status~~ done — TODO_LIST.md maintained (v5 sections, 2026-11-02) | ~~3~~ | ~~10min~~ |
| ~~21~~ | ~~Update ROADMAP.md to reflect v3 direction~~ done — ROADMAP.md reflects v4/v5 direction (2026-11-02) | ~~3~~ | ~~15min~~ |
| ~~22~~ | ~~Add CLI option to wire telemetry module (WithTelemetry)~~ done — WithTelemetry[T] ships in telemetry/ (telemetry.go:73) | ~~4~~ | ~~15min~~ |
| ~~23~~ | ~~Consider cutting docgen.go if rarely used~~ done — kept — GenerateDocs/docgen.go + describe JSON contract ship (FEATURES.md) | ~~2~~ | ~~10min~~ |
| ~~24~~ | ~~Audit all sentinel errors for completeness~~ done — sentinelFacts covers all sentinels (errors_facts.go) | ~~4~~ | ~~30min~~ |
| ~~25~~ | ~~Performance: benchmark NewCLI with cliSpec vs old struct~~ **Superseded — superseded — the pre-cliSpec struct is gone in v4; perf tracked as ROADMAP raw ideas (Benchmark dashboard/CI gating).** | ~~3~~ | ~~30min~~ |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should fang be an optional module too?**

fang pulls lipgloss, mango-cobra, and ~15 transitive deps. If we made fang optional (core uses plain cobra, users who want styled output import the fang module), core direct deps would drop to ~8 (cobra, pflag, do, go-toml, go-yaml, go-output, samber-do-auditlog, x/term).

But fang provides the core "superb CLI" experience — styled help, styled errors, automatic version. Making it optional undermines the "batteries included" promise. The question is: **is the dependency cost worth the UX benefit for users who don't want styled output?**

**Answered:** kept — fang stays a direct core dependency in v4 (`charm.land/fang/v2 v2.0.1` in go.mod); the batteries-included direction won.

## Resolution (2026-07-23)

- §a "5 Optional Modules" extraction was completed before v3.0.0; the workspace later moved to 4 sub-modules after `manpage` was removed (`34a0c6e`).
- §c "NOT STARTED" docs/release tasks (migration guide, GitHub Releases, sub-module tests, lint cleanup) were closed in the 2026-07-07 through 2026-07-11 sessions.
- Remaining architecture items (koanf extraction, middleware context propagation, public API renames) are deferred to v3.1+/v4 and tracked in `ROADMAP.md` "Deferred from 2026-07-18 Audit Closure".
- Current metrics: 470 test functions, 1434 runs, 87.8% core coverage, 0 lint issues.
