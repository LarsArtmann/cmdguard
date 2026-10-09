# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, TODO(v5) markers,
> status report harvests, and code analysis. Updated as work completes.

---

## v5 Breaking Changes (Deferred)

These items have `TODO(v5)` markers in source code. They are public API
renames/removals from the 2026-07-18 naming review (plus follow-ups) that
cannot ship in v4.x without breaking downstream consumers.

| #  | Task                                                                  | File                 | Status  |
| -- | --------------------------------------------------------------------- | -------------------- | ------- |
| T1 | Rename `CommandInfo` → `CommandMetadata`                              | `middleware.go`      | 🔴 TODO |
| T2 | Rename `TypeHandler` → `TypeCodec`                                    | `type_handler.go`    | 🔴 TODO |
| T3 | Rename `PromptRunner` → `HuhPrompter` (or similar)                    | `prompts/prompts.go` | 🔴 TODO |
| T4 | Remove `SetConfig` (now `// Deprecated:`)                             | `cli_accessors.go`   | 🔴 TODO |
| T5 | Replace `RegisterInScope` with generic variant (now `// Deprecated:`) | `scope.go`           | 🔴 TODO |

Further v5 candidates (Get[T] rename, Package[T] redesign) live in
ROADMAP.md §"Additional v5 candidates" — keep them there, not here.

> These must stay as `TODO` (not `NOTE`) so `grep TODO` finds them. The godox
> linter exclusion for `TODO(v5)` is deliberately narrow.

---

## Blocked on Upstream

| #  | Task                           | Blocker                                                                                                                                                                                                                | Priority | Status     |
| -- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------- |
| F1 | Command-level audit middleware | DONE 2026-10-09: upstream v0.11.0 hit the module proxy (tag pushed); cmdguard `go.mod` bumped to v0.11.0; `AuditMiddleware[T]` + `WithAuditMiddleware[T]` implemented in `auditlog.go` (PhaseBefore/PhaseAfter + duration ms + error, nil-plugin passthrough, FullPath-preferred naming) with 7 tests in `audit_middleware_test.go`; taskctl wires it in `buildApp` and a live `taskctl list` run exported command events (`event_type: command`, before/after, `duration_ms: 1.264`). FEATURES row → FULLY_FUNCTIONAL. | Medium   | ✅ DONE     |

---

## Lint Sub-Module (`lint/`)

Harvested from `docs/status/2026-10-04_07-11_cmdguard-lint-sub-module-build.md` §f
(full 50-item backlog lives there; these are the near-term actionable ones).

| #  | Task                                                         | Notes                                                                                                                                                                                            | Priority | Status      |
| -- | ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | ----------- |
| L1 | Tag `lint/v0.1.0`                                            | DONE 2026-10-04: user approved; CHANGELOG entry, annotated tag pushed, proxy-verified (scratch consumer go get + build + run), GitHub release published                                          | High     | ✅ DONE     |
| L2 | Wire `lint/provider` into the BuildFlow repo                 | DONE 2026-10-04: user approved wire-now; blank import in sdk_imports.go, require v0.1.0 + pinned replace, vendor + vendorHash, wiring test + inventory guards, docs --check green (provider 118) | High     | ✅ DONE     |
| L3 | Cross-file dataflow for CG003/CG004/CG006                    | Name tracing is file-scoped today; constructor in one file + use in another is missed                                                                                                            | Medium   | 🔴 TODO     |
| L4 | Baseline / ratchet mode                                      | Save snapshot of current findings so adoption is incremental (old pass, new fail)                                                                                                                | Medium   | 🔴 TODO     |
| L5 | Provider options (enable/disable via BuildFlow config)       | `provider.Register` spec has no options surface yet                                                                                                                                              | Medium   | 🔴 TODO     |
| L6 | golangci-lint module-plugin distribution layer               | `.custom-gcl.yml` wiring per go-humanize-linter pattern                                                                                                                                          | Low      | 🔴 TODO     |
| L7 | Auto-generate README rule table from `allRuleDefs()`         | Prevents doc drift when rules change (golden test or `go:generate`)                                                                                                                              | Low      | 🔴 TODO     |
| L8 | Bump `CurrentMajor` in `lint/walk.go` when cmdguard v5 ships | Standing duty; CG002 tracks the constant. Rule IDs CG001–CG006 are load-bearing (suppressions key on them) — never renumber                                                                      | —        | 🔵 STANDING |

---

## Security

| #  | Task                                     | Notes                                                                                                                                                                                                                                                                                    | Priority | Status  |
| -- | ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| S1 | Triage the 16 Dependabot vulnerabilities | DONE 2026-10-09: all 25 lifetime alerts (incl. the 16 reported 2026-10-04) show state=fixed on GitHub — 22 auto-resolved 2026-10-04 21:07–21:10 UTC by the dependency bumps pushed that day, 3 earlier, 0 dismissed, 0 open (`gh api .../dependabot/alerts` + `pnpm audit` verified). Residual `postcss-selector-parser <7.1.6` (GHSA-rj75-hqrm-r3gf, moderate, via starlight→expressive-code→postcss-nested) fixed by `postcss-nested: ^8.0.1` override in `website/pnpm-workspace.yaml`; `pnpm audit` clean, `astro build` green (22 pages). | High     | ✅ DONE |
