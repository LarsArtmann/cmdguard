# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, TODO(v5) markers,
> status report harvests, and code analysis. Updated as work completes.

---

## v5 Breaking Changes (Deferred)

These items have `TODO(v5)` markers in source code. They are public API
renames/removals from the 2026-07-18 naming review (plus follow-ups) that
cannot ship in v4.x without breaking downstream consumers.

| #  | Task                                               | File                     | Status  |
| -- | -------------------------------------------------- | ------------------------ | ------- |
| T1 | Rename `CommandInfo` → `CommandMetadata`           | `middleware.go`          | 🔴 TODO |
| T2 | Rename `TypeHandler` → `TypeCodec`                 | `type_handler.go`        | 🔴 TODO |
| T3 | Rename `PromptRunner` → `HuhPrompter` (or similar) | `prompts/prompts.go`     | 🔴 TODO |
| T4 | Remove `SetConfig` (now `// Deprecated:`)          | `cli_accessors.go`       | 🔴 TODO |
| T5 | Replace `RegisterInScope` with generic variant (now `// Deprecated:`) | `scope.go` | 🔴 TODO |

Further v5 candidates (Get[T] rename, Package[T] redesign) live in
ROADMAP.md §"Additional v5 candidates" — keep them there, not here.

> These must stay as `TODO` (not `NOTE`) so `grep TODO` finds them. The godox
> linter exclusion for `TODO(v5)` is deliberately narrow.

---

## Blocked on Upstream

| #  | Task                                     | Blocker                                                                                                                                  | Priority | Status     |
| -- | ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------- |
| F1 | Command-level audit middleware           | Upstream API is implemented: `Plugin.RecordCommand` / `Recorder.RecordCommand` / `EventTypeCommand` in `../samber-do-auditlog` (Unreleased, committed locally). **Unblock:** push + tag samber-do-auditlog ≥ v0.11.0, bump cmdguard `go.mod`, then add `AuditMiddleware[T]` (records PhaseBefore/PhaseAfter + duration + error per command execution). | Medium   | 🟡 BLOCKED |
