# Dependency Graph — cmdguard

**Date:** 2026-05-14
**Status:** Current + Proposed

---

## Current Dependency Graph (Internal Packages)

```
                    ┌──────────────┐
                    │  go.mod root │
                    │ (monolith)   │
                    └──────┬───────┘
                           │
          ┌────────────────┼─────────────────┐
          │                │                  │
    ┌─────▼─────┐   ┌──────▼──────┐   ┌──────▼──────┐
    │ pkg/      │   │ examples/*  │   │ tests/      │
    │ cmdguard/ │   │ (12 pkgs)   │   │ integration │
    │ v2        │   │             │   │             │
    └─────┬─────┘   └──────┬──────┘   └──────┬──────┘
          │                │                  │
          │           ┌────▼────┐        ┌────▼────┐
          │           │examples │        │ (direct │
          │           │internal │        │  v2     │
          │           └─────────┘        │  import)│
          │                              └─────────┘
    ┌─────▼──────┐
    │ pkg/       │
    │ testutil   │
    └────────────┘

    Dependency: examples/* ──► v2
    Dependency: examples/internal ──► v2
    Dependency: tests/integration ──► v2
    Dependency: v2 (tests) ──► testutil
```

### Current External Dependencies

| Dependency                         | Version | Used By                |
| ---------------------------------- | ------- | ---------------------- |
| `charm.land/fang/v2`               | v2.0.1  | v2 (CLI styling)       |
| `github.com/spf13/cobra`           | v1.10.2 | v2 (CLI framework)     |
| `github.com/spf13/pflag`           | v1.0.10 | v2 (flag parsing)      |
| `github.com/samber/do/v2`          | v2.0.0  | v2 (DI)                |
| `github.com/larsartmann/go-output` | v0.2.0  | v2 (output formatting) |
| `github.com/muesli/mango`          | v0.2.0  | v2 (man pages)         |
| `github.com/muesli/mango-cobra`    | v1.3.0  | v2 (man pages)         |
| `github.com/muesli/roff`           | v0.1.0  | v2 (man pages)         |

---

## Proposed Dependency Graph (Post-Modularization)

```
    ┌──────────────────────────────────────────────────────────┐
    │                     go.work                              │
    │  use (., ./types, ./output, ./testutil)                  │
    └──────────────────────────────────────────────────────────┘

    ┌──────────────────┐
    │ cmdguard (root)  │
    │ pkg/cmdguard/v2  │
    │                  │
    │ External deps:   │
    │  - cobra v1.10.2 │
    │  - pflag v1.0.10 │
    │  - samber/do/v2  │
    │  - fang/v2       │
    │  - mango         │
    │  - mango-cobra   │
    │  - roff          │
    └──┬────────────┬──┘
       │            │
       ▼            ▼
    ┌──────────┐  ┌───────────────┐
    │ types    │  │ output        │
    │          │  │               │
    │ Ext:     │  │ Ext:          │
    │  (none)  │  │  - go-output  │
    └──────────┘  └───────────────┘

    ┌──────────────────┐
    │ testutil         │
    │                  │
    │ Ext:             │
    │  - cobra         │
    │  (test-only)     │
    └──────────────────┘

    examples/* ──► cmdguard (root)
    tests/integration ──► cmdguard (root)
    types (tests) ──► testutil
    output (tests) ──► testutil
    cmdguard (root, tests) ──► testutil
```

### Proposed External Dependencies Per Module

| Module            | External Dependencies                                                                                                                 | Lines of go.mod |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| `cmdguard` (root) | cobra, pflag, samber/do/v2, fang/v2, mango, mango-cobra, roff, go-output (transitive via output), types (internal), output (internal) | ~20             |
| `types`           | **None** (stdlib only)                                                                                                                | ~3              |
| `output`          | go-output                                                                                                                             | ~5              |
| `testutil`        | cobra                                                                                                                                 | ~4              |

### Dependency Flow for Key Scenarios

**Consumer wants only value types:**

```
go get github.com/larsartmann/cmdguard/types
→ pulls in: nothing (stdlib only)
```

**Consumer wants output formatting only:**

```
go get github.com/larsartmann/cmdguard/output
→ pulls in: go-output + transitive deps
```

**Consumer wants full CLI framework:**

```
go get github.com/larsartmann/cmdguard
→ pulls in: cobra, pflag, samber/do/v2, fang/v2, mango, etc.
→ transitively: types, output
```

---

## v2 Package Internal Coupling (Detailed)

### Production Code Coupling Matrix

| File → depends on  | errors | config | config_parse | config_set | flags | flags_parse | flags_validate | flags_suggest | flag_helpers | type_handler | type_helpers | types\_\* | scope | flow_ctx | middleware | cli | cli_options | cli_accessors | cli_command | cli_output | output | command | completion | editor | manpage | suggest |
| ------------------ | ------ | ------ | ------------ | ---------- | ----- | ----------- | -------------- | ------------- | ------------ | ------------ | ------------ | --------- | ----- | -------- | ---------- | --- | ----------- | ------------- | ----------- | ---------- | ------ | ------- | ---------- | ------ | ------- | ------- | --- |
| cli.go             | ●      |        |              |            | ●     |             |                |               |              |              |              |           | ●     | ●        | ●          |     | ●           |               | ●           | ●          | ●      | ●       |            |        |         |         |     |
| cli_command.go     | ●      |        |              |            | ●     |             |                |               | ●            | ●            |              |           |       |          | ●          | ●   |             |               |             |            |        | ●       | ●          |        |         |         |
| cli_options.go     |        |        |              |            |       |             |                |               |              |              |              |           | ●     |          | ●          | ●   |             |               |             |            |        |         |            |        |         |         |
| cli_accessors.go   |        |        |              |            |       |             |                |               |              |              |              |           | ●     | ●        |            | ●   |             |               |             |            |        |         |            |        |         |         |
| cli_output.go      |        |        |              |            |       |             |                |               |              |              |              |           |       |          |            | ●   |             | ●             |             |            | ●      |         |            |        |         |         |
| command.go         | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |         |
| config.go          | ●      |        | ●            |            |       |             |                |               |              | ●            |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| config_parsing.go  | ●      |        |              |            |       |             |                |               |              | ●            |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| config_setfield.go | ●      |        |              |            |       |             |                |               |              | ●            |              | ●         |       |          |            |     |             |               |             |            |        |         |            |        |         |
| flags.go           | ●      |        | ●            |            |       |             | ●              | ●             |              | ●            |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |         |
| flags_parse.go     | ●      |        |              | ●          | ●     |             |                |               |              | ●            |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| flags_validate.go  | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| flags_suggest.go   |        |        |              |            | ●     |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| flag_helpers.go    | ●      | ●      |              |            | ●     |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| type_handler.go    |        | ●      | ●            |            |       |             |                |               |              |              |              | ●         |       |          |            |     |             |               |             |            |        |         |            |        |         |
| type_helpers.go    |        |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| scope.go           | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            | ●   | ●           |               |             |            |        |         |            |        |         |
| flow_context.go    | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| middleware.go      | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| output.go          | ●      |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| editor.go          |        |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| manpage.go         |        |        |              |            |       |             |                |               |              |              |              |           |       |          |            | ●   |             |               |             |            |        |         |            |        |         |
| command_suggest.go |        |        |              |            |       |             |                | ●             |              |              |              |           |       |          |            |     |             |               |             |            |        |         |            |        |         |
| completion.go      |        |        |              |            |       |             |                |               |              |              |              |           |       |          |            |     |             |               |             | ●          |        |         |            |        |
| types\_\*.go       | ●      |        |              |            |       |             |                |               |              |              | ●            |           |       |          |            |     |             |               |             |            |        |         |            |        |
