# ADR-002: Lint Strategy and Exclusion Policy

**Status:** Accepted\
**Date:** 2026-07-11\
**Deciders:** Lars Artmann

## Context

cmdguard uses golangci-lint 2.x with 40+ enabled linters. As the codebase grew, lint exclusions accumulated — some legitimate (design decisions), others lazy shortcuts that masked real issues. This ADR documents the strategy for maintaining zero lint issues through real fixes, not silencing linters.

## Decision

### Principle: Fix code, don't silence linters

Every lint exclusion must be either:

1. **A documented design decision** — the lint rule conflicts with a deliberate architectural choice (e.g., package-level registries for the COW pattern)
2. **An ireturn allow-list entry** — a legitimate interface return from an external library (e.g., `do.Injector`, `koanf.Parser`)

### Current exclusions (4 per-file rules + 4 ireturn allow-list entries)

| Rule | File                | Linter           | Justification                                      |
| ---- | ------------------- | ---------------- | -------------------------------------------------- |
| 1    | `type_handler.go`   | gochecknoglobals | `globalTypeRegistry` — COW pattern foundation      |
| 2    | `flags_validate.go` | gochecknoglobals | `globalValidators`, `regexCache` — COW pattern     |
| 3    | `cli_command.go`    | gochecknoglobals | `argsKey`, `configKey` — Go context key convention |
| 4    | `example_test.go`   | forbidigo        | `fmt.Println` — godoc examples must print          |

ireturn allow list: `error`, `empty`, `anon`, `stdlib`, `generic` (built-in), plus `ConfigFileLoader`, `TypeHandler`, `do.Injector`, `koanf.Parser`.

### Regression signal

Track the exclusion count. If it increases, investigate whether the new exclusion is a real fix or a shortcut. The count is documented in AGENTS.md and should only decrease over time.

### What was fixed (not excluded)

- wrapcheck in `output.go` — external errors wrapped via `wrapIfError` helper
- wrapcheck in `type_handler.go` — `dispatchRegister` returns wrapped with `fmt.Errorf`
- funlen in `type_handler_kinds.go` — `registerKinds()` split into 7 focused helpers
- funlen in `type_handler_custom.go` — split into `registerEnumTypes()` + `registerValueTypes()`
- cyclop/funlen in `cli.go` — `initialize()` split into `ensureScope()` + `setupPersistentPreRun()`
- paralleltest — `t.Parallel()` added to all applicable test functions

## Consequences

- Zero lint issues achieved through real code fixes
- 4 remaining exclusions are all documented design decisions
- Exclusion count tracked as a regression metric
- Any new exclusion requires justification in this ADR
