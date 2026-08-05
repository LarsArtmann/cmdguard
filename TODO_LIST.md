# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, TODO(v5) markers, and
> code analysis. Updated as work completes.

---

## Status Legend

| Status | Meaning |
| ------ | ------- |
| 🔴 TODO | Not started |
| 🟡 IN PROGRESS | Actively being worked on |
| 🟢 DONE | Completed |

---

## v5 Breaking Changes (Deferred)

These items have `TODO(v5)` markers in source code. They are public API renames
from the 2026-07-18 naming review that cannot ship in v4.x without breaking
downstream consumers.

| # | Task | File | Status |
|---|------|------|--------|
| T1 | Rename `CommandInfo` → `CommandMetadata` | `middleware.go:40` | 🔴 TODO |
| T2 | Rename `TypeHandler` → `TypeCodec` | `type_handler.go:13` | 🔴 TODO |
| T3 | Rename `PromptRunner` → `HuhPrompter` (or similar) | `prompts/prompts.go:27` | 🔴 TODO |

> These must stay as `TODO` (not `NOTE`) so `grep TODO` finds them. The godox
> linter exclusion for `TODO(v5)` is deliberately narrow.

---

## Partially Functional Items

| # | Task | Priority | Status |
|---|------|----------|--------|
| P1 | Remove `SetConfig(cfg)` — unsafe post-construction mutation without re-initializing FlagRegistry | High | 🔴 TODO |
| P2 | Fix `Get[T]` → rename to `GetService[T]` (v5 breaking change) | Medium | 🔴 TODO |
| P3 | Fix `RegisterInScope(parent, name, ...any)` — takes `...any` instead of being generic | Medium | 🔴 TODO |
| P4 | Redesign `Package[T](scope, ...)` — unusual API shape with pre-existing `*Scope` param | Low | 🔴 TODO |
| P5 | Fix middleware context propagation — `next func() error` doesn't propagate context, blocks timeout/cancellation middleware | Medium | 🔴 TODO |

---

## Planned Features

| # | Task | Priority | Status |
|---|------|----------|--------|
| F1 | Command-level audit middleware — only DI lifecycle events captured; command-level events not implemented | Medium | 🔴 TODO |

---

## Technical Debt

| # | Task | Priority | Status |
|---|------|----------|--------|
| D1 | Improve `pkg/testutil` coverage (currently 49.6%) | Low | 🔴 TODO |
| D2 | Improve `examples/taskctl` coverage (currently 68.2%, below core's 87.8%) | Low | 🔴 TODO |

---

## Recently Completed

| # | Task | Date |
|---|------|------|
| C1 | Documentation drift fix — updated all user-facing docs from v3 to v4 API | 2026-08-05 |
| C2 | AGENTS.md project structure corrected from v3 to v4 | 2026-08-05 |
| C3 | CHANGELOG.md rebuilt with real version history (v0.1.0 through v4.0.0) | 2026-08-05 |
| C4 | TODO_LIST.md and ROADMAP.md created (were referenced but missing) | 2026-08-05 |
