# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, TODO(v5) markers,
> status report harvests, and code analysis. Updated as work completes.

---

## v5 Breaking Changes (Deferred)

These items have `TODO(v5)` markers in source code. They are public API renames
from the 2026-07-18 naming review that cannot ship in v4.x without breaking
downstream consumers.

| #   | Task                                               | File                    | Status    |
| --- | -------------------------------------------------- | ----------------------- | --------- |
| T1  | Rename `CommandInfo` → `CommandMetadata`           | `middleware.go:40`      | 🔴 TODO   |
| T2  | Rename `TypeHandler` → `TypeCodec`                 | `type_handler.go:13`    | 🔴 TODO   |
| T3  | Rename `PromptRunner` → `HuhPrompter` (or similar) | `prompts/prompts.go:27` | 🔴 TODO   |

> These must stay as `TODO` (not `NOTE`) so `grep TODO` finds them. The godox
> linter exclusion for `TODO(v5)` is deliberately narrow.

---

## Documentation Drift (v3→v4 Residuals)

Items from the 2026-08-05 documentation drift fix that were not completed.
Each verified against current source.

| #   | Task                                                                  | Evidence                  | Priority | Status    |
| --- | -------------------------------------------------------------------- | ------------------------- | -------- | --------- |
| D3  | Fix `CONTRIBUTING.md` — `v3` → `v4` (line 101 test package, line 115 heading) | `CONTRIBUTING.md:101,115` | High     | 🔴 TODO   |
| D4  | Fix `docs/ERROR_REFERENCE.md` title — "v2" → "v4"                    | `docs/ERROR_REFERENCE.md:1` | High     | 🔴 TODO   |
| D5  | Create `docs/MIGRATION_v3_v4.md` — dedicated v3→v4 migration guide   | _(file does not exist)_   | Medium   | 🔴 TODO   |
| D6  | Add `manpage` removal note to `docs/MIGRATION_v2_v3.md` §3           | `docs/MIGRATION_v2_v3.md` | Low      | 🔴 TODO   |

---

## Partially Functional Items

| #   | Task                                                                                                                       | Priority | Status    |
| --- | -------------------------------------------------------------------------------------------------------------------------- | -------- | --------- |
| P1  | Remove `SetConfig(cfg)` — unsafe post-construction mutation without re-initializing FlagRegistry                           | High     | 🔴 TODO   |
| P2  | Fix `Get[T]` → rename to `GetService[T]` (v5 breaking change)                                                              | Medium   | 🔴 TODO   |
| P3  | Fix `RegisterInScope(parent, name, ...any)` — takes `...any` instead of being generic                                      | Medium   | 🔴 TODO   |
| P4  | Redesign `Package[T](scope, ...)` — unusual API shape with pre-existing `*Scope` param                                     | Low      | 🔴 TODO   |
| P5  | Fix middleware context propagation — `next func() error` doesn't propagate context, blocks timeout/cancellation middleware | Medium   | 🔴 TODO   |

---

## Planned Features

| #   | Task                                                                                                     | Priority | Status    |
| --- | -------------------------------------------------------------------------------------------------------- | -------- | --------- |
| F1  | Command-level audit middleware — only DI lifecycle events captured; command-level events not implemented | Medium   | 🔴 TODO   |

---

## Technical Debt

| #   | Task                                                                                  | Evidence                 | Priority | Status    |
| --- | ------------------------------------------------------------------------------------- | ------------------------ | -------- | --------- |
| D1  | Improve `pkg/testutil` coverage (currently 49.6%)                                     | `pkg/testutil/`          | Low      | 🔴 TODO   |
| D2  | Improve `examples/taskctl` coverage (currently 68.2%, below core's 87.8%)             | `examples/taskctl/`      | Low      | 🔴 TODO   |
| D7  | Clean up residual git corruption — `git fsck` reports broken links + invalid reflog entry `3e483b3b` | `git fsck --connectivity-only` | Medium   | 🔴 TODO   |
| D8  | Tag `flightrecorder` v0.1.0 — other 4 sub-modules are tagged, this one is not         | `git tag -l`             | Medium   | 🔴 TODO   |
| D9  | Restore lost flightrecorder godoc examples (`ExampleWithFlightRecorderRecorder`, `ExampleRecorder_CaptureToWriter`) — lost in git corruption | `flightrecorder/example_test.go` | Low      | 🔴 TODO   |
| D10 | Test `go tool trace` parseability of flightrecorder snapshots                         | `flightrecorder/`        | Low      | 🔴 TODO   |
