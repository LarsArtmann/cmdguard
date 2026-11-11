# v3 Superb CLI Redesign — Session Status Report

**Date:** 2026-07-06 09:55
**Branch:** master
**Commits this session:** 9 (4f9b0ea..1c4fd65)
**All tests:** GREEN (5 packages, 1908 test runs)
**Coverage:** 87.1%

---

> **Update 2026-07-23:** The "NOT STARTED" spinner/output/auditlog extractions and remaining docs were completed in the follow-up sessions from `9297257` to `a9c8e82`. The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## a) FULLY DONE

### 1. Type Parameter Explosion Eliminated

**Commit:** `4f9b0ea`, `e2b125c`

`CommandOption` is now non-generic. All metadata options (`WithShort`, `WithLong`, `WithExample`, `WithGroupID`, `WithNoArgs`, etc.) require zero type parameters. `NewCommand` takes `flags F` as a positional arg so Go infers T and F automatically.

**Before (7 type params per command):**

```go
v2.NewCommand[AppConfig, *ListFlags]("list", handler,
    v2.WithShort[AppConfig, *ListFlags]("List tasks"),
    v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
)
```

**After (zero type params):**

```go
v2.NewCommand("list", &ListFlags{}, handler,
    v2.WithShort("List tasks"),
)
```

`WithFlags` option deleted entirely — flags are positional. All 150+ source files and 40+ test files migrated.

### 2. Lifecycle Hook Type Safety Restored

**Commit:** `428af36`

After the initial de-genericization, PreRunE/PostRunE were stored as `any` — a compile-time safety regression. Fixed with a sealed interface pattern:

- `lifecycleHook` interface with unexported `isLifecycleHook()` method
- `typedHook[T,F]` struct carries the typed function
- `subcommandList` interface + `typedSubcommands[T,F]` for subcommands
- `wireAllHandlers` uses safe `ok` type assertion — returns nil on mismatch, no panic

### 3. Mono-Repo Modularization (4 of ~7 modules extracted)

**Commits:** `29b44a6`, `3541134`, `dea535a`, `11dc1a5`, `1c4fd65`

<!-- check-rows:ledger -->

| Module    | Path                      | Deps Removed                                             | Status      |
| --------- | ------------------------- | -------------------------------------------------------- | ----------- |
| telemetry | `pkg/cmdguard/telemetry/` | otel SDK + exporters                                     | ✅ Done     |
| manpage   | `pkg/cmdguard/manpage/`   | mango/roff (from core code; fang still pulls indirectly) | ✅ Done     |
| glamour   | `pkg/cmdguard/glamour/`   | chroma, goldmark, bluemonday, gorilla/css                | ✅ Done     |
| prompts   | `pkg/cmdguard/prompts/`   | huh, bubbles, bubbletea (full TUI framework)             | ✅ Done     |
| spinner   | —                         | lipgloss (already from fang)                             | NOT STARTED |
| output    | —                         | go-output 10 sub-modules                                 | NOT STARTED |
| auditlog  | —                         | samber-do-auditlog                                       | NOT STARTED |

Core direct deps: **30 → 23** (7 removed). Indirect deps significantly reduced.

**Extension hooks added to core:**

- `HelpTransformFunc` + `WithHelpTransform[T]()` — glamour module implements this
- `PromptRunner` interface + `SetPromptRunner()` — prompts module implements this

### 4. go.work Workspace

**Commit:** `29b44a6`

Multi-module workspace with `use` directives for all 5 modules (root + 4 optional). All modules build with `go build ./...`.

### 5. Pareto Plan Written

**Commit:** `3310586`

`docs/planning/2026-07-06_06-54_v3-superb-cli-redesign.md` — 18 tasks, 89 subtasks, execution graph.

---

## b) PARTIALLY DONE

### 1. go-output Cleanup — 0% code done, design clear

10 blank-imported go-output sub-modules still in `output.go`. Plan: remove blank imports, keep `OutputFormat` type alias, users import go-output directly for `FormatTable`/`FormatJSON`.

### 2. Spinner Extraction — not started

`spinner.go` still in core. Lipgloss is already pulled by fang, so dep tree win is minimal, but code decoupling is valuable.

### 3. Audit Log Deepening — not started

Current audit only captures DI lifecycle events. Planned: command-level audit middleware, built-in `audit-log` subcommand.

### 4. Documentation — not started

AGENTS.md, README.md, FEATURES.md, CHANGELOG.md all reflect v2 API, not the new v3 patterns. doc.go has stale `WithFlags` references.

---

## c) NOT STARTED

1. ~~**Remove go-output blank imports** (output.go — 10 sub-modules)~~ done — shipped — direct go-output usage; only OutputFormat alias re-exported (AGENTS.md)
2. ~~**Extract spinner** to sub-module~~ done — spinner/ sub-module ships (spinner/v0.2.0)
3. ~~**Cut result.go** (sum types — not a CLI concern)~~ done — cut in the v3 slim-down (ROADMAP keeps a Result[T] revival idea)
4. ~~**Cut editor.go** ($EDITOR support — marginal)~~ done — cut — editor.go no longer exists in v4
5. ~~**Command-level audit middleware** (auto-audit every command)~~ done — AuditMiddleware[T] shipped (auditlog.go; FEATURES.md)
6. **Built-in audit-log subcommand** (`myapp audit-log --format d2`)
7. ~~**Update example** to showcase full v3 patterns (service-owned config, clean API)~~ done — examples/ maintained through v4 (taskctl flagship)
8. ~~**Update README.md** with v3 quickstart~~ done — README rewritten through v3/v4 releases
9. ~~**Update AGENTS.md** module structure section~~ done — AGENTS.md documents the 8-module workspace (2026-11-02)
10. ~~**Update FEATURES.md** with v3 changes~~ done — FEATURES.md maintained through v4
11. ~~**CHANGELOG.md** entry for v3 breaking changes~~ done — CHANGELOG [3.0.0] section (2026-07-07)
12. ~~**Lint fixes** (4 formatting issues in new code)~~ done — 0 lint issues (FEATURES.md; CI lint job)

---

## d) TOTALLY FUCKED UP (Honest Self-Assessment)

### 1. CLIOption Still Generic

I deferred making `CLIOption` non-generic because "only 4 of 22 options need explicit `[T]`." That's a weak excuse. The same sealed-interface pattern that fixed lifecycle hooks would work here. The inconsistency — `CommandOption` non-generic, `CLIOption` still generic — is a design smell.

### 2. Stale Documentation in doc.go

`pkg/cmdguard/v2/doc.go` still references `WithFlags[T,F](&flags{})` — an option that **no longer exists**. Anyone reading the package docs sees broken API examples.

### 3. go.work Not Properly Set Up for go.work.sum

`go.work.sum` was not committed. Downstream consumers building with `GOWORK=off` may hit checksum issues. The replace directives in root go.mod for sub-modules are also incomplete — only glamour has one, not manpage/telemetry/prompts (though go.work resolves them in dev).

### 4. No CHANGELOG Entry for Breaking Changes

We made **massive** breaking API changes (NewCommand signature changed, WithFlags deleted, 4 modules extracted) and there's zero CHANGELOG entry. Anyone tracking changes has no idea what happened.

### 5. Lint Issues Introduced

4 new lint issues from the session's edits (gci formatting, golines line length) in:

- `cli_cobra_command_test.go` (2 issues)
- `cli_hooks_test.go` (1 issue)
- `prompts.go` (1 issue)

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Service-owned config** — the biggest design discussion from this session. Services should own their typed config through the DI graph, not through flags. Flags become optional per-command input, not the primary config mechanism. This would make cmdguard genuinely unique.

2. ~~**CLIOption de-genericization** — apply the same sealed-interface pattern from lifecycle hooks. Eliminates `[T]` from all 22 CLI options.~~ done — non-generic options in v4 (sealed-interface pattern)

3. **Flag scope enforcement** — prevent global flags from polluting subcommands. Use the DI graph to determine which flags belong where.

### Code Quality

4. **koanf still in core** (4 direct deps) — should move to configload sub-module or be made optional. **→ ROADMAP deferred #5 (koanf extraction YAGNI)**

5. ~~**flow_context.go** (253+68 lines) — rarely useful, adds complexity. Consider cutting or internalizing.~~ __NOT-DO — decided to keep — flow_context_.go ships as the branching-flow feature (AGENTS.md)._*

6. ~~**Dead weight** — result.go (147 lines, sum types not a CLI concern), editor.go (60 lines, marginal).~~ done — cut — both removed in the v3 slim-down (editor.go gone in v4)

### Module Hygiene

7. ~~**Replace directives incomplete** — only glamour has a replace in root go.mod. manpage/telemetry/prompts rely on go.work only.~~ done — all sub-modules pinned in root go.mod (AGENTS.md publish-shape rule)

8. ~~**go.work.sum not committed** — checksum file missing.~~ done — go.work committed at repo root; workspace builds

9. ~~**GOWORK=off verification** — never ran `GOWORK=off go build` for each sub-module to verify they resolve independently.~~ done — codified — AGENTS.md GOWORK=off publish-shape rule (2026-11-02)

---

## f) Next 25 Tasks (Sorted by Impact/Effort)

| #      | Task                                                                                                                                    | Impact     | Effort    |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------- | ---------- | --------- |
| ~~1~~  | ~~Fix doc.go stale WithFlags references~~ done — doc.go rewritten in v3 (WithFlags deleted)                                             | ~~High~~   | ~~5min~~  |
| ~~2~~  | ~~Fix 4 lint issues (gci/golines formatting)~~ done — 0 lint issues (FEATURES.md; CI lint job)                                          | ~~Medium~~ | ~~10min~~ |
| ~~3~~  | ~~Add CHANGELOG.md entry for v3 breaking changes~~ done — CHANGELOG [3.0.0] section (2026-07-07)                                        | ~~High~~   | ~~15min~~ |
| ~~4~~  | ~~Commit go.work.sum + add missing replace directives~~ done — resolved pre-release (Resolution 2026-07-18); replaces complete          | ~~High~~   | ~~10min~~ |
| ~~5~~  | ~~Remove go-output blank imports from output.go~~ done — direct go-output usage; alias re-export only (AGENTS.md)                       | ~~High~~   | ~~30min~~ |
| ~~6~~  | ~~Cut result.go (sum types)~~ done — cut in the v3 slim-down (ROADMAP Result[T] revival idea)                                           | ~~Medium~~ | ~~10min~~ |
| ~~7~~  | ~~Cut editor.go ($EDITOR)~~ done — cut — editor.go no longer exists in v4                                                               | ~~Low~~    | ~~10min~~ |
| ~~8~~  | ~~Extract spinner to sub-module~~ done — spinner/ sub-module ships (spinner/v0.2.0)                                                     | ~~Low~~    | ~~30min~~ |
| ~~9~~  | ~~Verify GOWORK=off builds for all 5 modules~~ done — codified — AGENTS.md GOWORK=off publish-shape rule                                | ~~High~~   | ~~30min~~ |
| ~~10~~ | ~~Make CLIOption non-generic (sealed interface)~~ done — non-generic options in v4 (sealed-interface pattern)                           | ~~High~~   | ~~45min~~ |
| ~~11~~ | ~~Add command-level audit middleware~~ done — AuditMiddleware[T] shipped (auditlog.go; FEATURES.md)                                     | ~~High~~   | ~~60min~~ |
| 12     | Add built-in audit-log subcommand **→ ROADMAP raw ideas (built-in audit-log subcommand)**                                               | Medium     | 45min     |
| 13     | Move koanf to configload or make optional **→ ROADMAP deferred #5 (koanf extraction YAGNI)**                                            | Medium     | 30min     |
| ~~14~~ | ~~Cut or internalize flow_context.go~~ __NOT-DO — decided to keep — flow_context_.go ships as the branching-flow feature (AGENTS.md)._* | ~~Low~~    | ~~20min~~ |
| ~~15~~ | ~~Update README.md with v3 quickstart~~ done — README rewritten through v3/v4 releases                                                  | ~~High~~   | ~~30min~~ |
| ~~16~~ | ~~Update AGENTS.md module structure~~ done — AGENTS.md documents the 8-module workspace (2026-11-02)                                    | ~~High~~   | ~~20min~~ |
| ~~17~~ | ~~Update FEATURES.md for v3~~ done — FEATURES.md maintained through v4                                                                  | ~~Medium~~ | ~~20min~~ |
| ~~18~~ | ~~Update example to showcase v3 patterns~~ done — examples/ maintained through v4 (taskctl flagship)                                    | ~~Medium~~ | ~~45min~~ |
| 19     | Design service-owned config pattern (ADR) **→ ROADMAP raw ideas (service-owned config ADR)**                                            | Critical   | 60min     |
| 20     | Implement WithProvider for flag scoping **→ ROADMAP raw ideas (flag scoping / WithProvider)**                                           | High       | 90min     |
| 21     | Add strict flag audit mode **→ ROADMAP raw ideas (strict flag audit mode)**                                                             | Medium     | 45min     |
| 22     | Consider cutting go-toml/yaml from core **→ ROADMAP deferred #5 (koanf extraction YAGNI; toml/yaml are koanf parser deps)**             | Medium     | 30min     |
| ~~23~~ | ~~Update MIGRATION_FROM_COBRA.md for v3~~ done — docs/MIGRATION_FROM_COBRA.md exists                                                    | ~~Low~~    | ~~20min~~ |
| ~~24~~ | ~~Add integration test for sub-module imports~~ done — .github/workflows/submodule-smoke.yml                                            | ~~Medium~~ | ~~30min~~ |
| ~~25~~ | ~~Release v3.0.0-alpha tag~~ **Superseded — superseded — v3.0.0 released directly (2026-07-07).**                                       | ~~High~~   | ~~10min~~ |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should the v3 API use `context.Context` as the first parameter in RunE handlers, replacing the separate `ctx` argument?**

Current:

```go
func(ctx context.Context, cfg *T, flags F) error
```

The discussion about service-owned config suggests handlers should resolve services from context:

```go
func(ctx context.Context) error  // cfg + services resolved via DI from ctx
```

But this is a fundamental API shape change that affects every consumer. I cannot determine whether you want to:

- (A) Keep the current 3-arg handler and add DI resolution as an additional layer
- (B) Flatten to `func(ctx) error` and resolve everything (config, flags, services) from context
- (C) Something else entirely

This decision shapes the entire v3 API surface. **Which direction do you want?**

**Answered:** (A) kept — v4 handlers still take `(ctx, cfg, flags)` (`cli_command.go:154`); context threading ships additively via `ContextMiddleware[T]` (2026-10-04, ROADMAP "Middleware Context Propagation"); signature unification is ROADMAP v5 candidate #8.

## Resolution (2026-07-18)

This v3-redesign session report is largely accurate but predates the v3.0.0 release (tagged 2026-07-07). Resolutions since: the §b "Not Started" extractions are now done — the project ships 6 Go modules (core + telemetry, manpage, glamour, prompts, spinner), all at repo root for external resolution. The §d #1 "CLIOption still generic" self-criticism was acted on (sealed-interface pattern applied). §d #3 (go.work.sum/replace directives) resolved pre-release. Current metrics: 87.6% coverage, 1429 test runs, 0 lint issues, jsonv2 migration embraced 2026-07-14.

## Resolution (2026-07-23)

- §3 "NOT STARTED" items (spinner/output sub-modules, docs, `go.work`, replace directives) were completed before v3.0.0.
- Sub-module tests and lint cleanup followed in `f8f3ad4` and `da3b454`.
- Architecture items (koanf extraction, middleware context propagation, API renames) remain deferred and are tracked in `ROADMAP.md`.
- Current metrics: 470 test functions, 1434 runs, 87.8% core coverage, 0 lint issues.
