# Modularization Proposal — cmdguard

**Date:** 2026-05-14
**Status:** Draft (post self-review)
**Project:** github.com/larsartmann/cmdguard
**Current Version:** v2.2.0

---

## 1. Executive Summary

cmdguard is a Go library for building validated Cobra CLI applications with type-safe dependency injection. It currently ships as a single `go.mod` monolith with one production package (`pkg/cmdguard/v2`), 12 example packages, 1 test utility package, and 1 integration test package.

This proposal splits cmdguard into **3 sub-modules** (plus root + testutil) with a `go.work`-coordinated workspace, extracting the value types and output rendering into independently versionable modules. The DI scope stays in core after self-review revealed circular dependency risks.

**Expected benefits:**

- **Hard dependency boundaries** — Value types can't accidentally import cobra; output can't depend on DI
- **Independent versioning** — Value types (`types`) and output formatting evolve on their own schedule
- **Cleaner go.mod** — Each module lists only the deps it actually uses
- **Reusability** — `types` and `output` modules can be consumed independently of cmdguard's CLI framework
- **Reduced binary size** for consumers who only need types — they don't pull in cobra/fang

---

## 2. Current State Analysis

### 2.1 Module Landscape

| Module | Path | Internal Deps | External Deps | State |
|--------|------|---------------|---------------|-------|
| `github.com/larsartmann/cmdguard` | `/` (root) | — | cobra, pflag, samber/do/v2, fang/v2, go-output, mango, mango-cobra, roff | Monolith |

### 2.2 Package Dependency Graph (Current)

```
                    ┌──────────────┐
                    │  go.mod root │
                    └──────┬───────┘
                           │
          ┌────────────────┼─────────────────┐
          │                │                  │
    ┌─────▼─────┐   ┌──────▼──────┐   ┌──────▼──────┐
    │ pkg/      │   │ examples/*  │   │ tests/      │
    │ cmdguard/ │   │ (12 pkgs)   │   │ integration │
    │ v2        │   │             │   │             │
    └───────────┘   └──────┬──────┘   └──────┬──────┘
                         │                  │
                    ┌────▼────┐        ┌─────▼────┐
                    │examples │        │ examples │
                    │internal │        │ internal │
                    └─────────┘        └──────────┘

    pkg/testutil ──► (used only by v2 tests)
```

### 2.3 Coupling Analysis

**v2 internal coupling (file-level, production code only):**

```
cli.go ──► scope.go, flags.go, flow_context.go, cli_options.go,
           cli_command.go, cli_output.go, errors.go, command.go,
           middleware.go, output.go

cli_command.go ──► flags.go, flag_helpers.go, command.go,
                   middleware.go, completion.go, errors.go

flags.go ──► config_parsing.go, type_handler.go, flags_validate.go,
             flags_suggest.go, errors.go

type_handler.go ──► config.go, types_*.go, config_parsing.go

scope.go ──► errors.go, cli_options.go, cli.go

output.go ──► errors.go
flow_context.go ──► errors.go
middleware.go ──► errors.go
editor.go ──► (standalone)
type_helpers.go ──► (standalone)
manpage.go ──► cli.go
```

**Key insight:** The coupling is layered, not tangled. No circular dependencies between concern clusters.

### 2.4 God-Package Assessment

`pkg/cmdguard/v2` has:

- **29 non-test files** (~5,000 LOC production code)
- **67 test files** (~10,769 LOC)
- **142 exported symbols**
- **9 logical concern clusters**

Moderate god-package tendency, but kept manageable by clear internal clustering with minimal cross-cluster coupling. Primary justification for splitting is dependency isolation and independent versioning, not code organization.

### 2.5 Files Over 300 Lines

| File | Lines | Concern |
|------|-------|---------|
| type_handler.go | 480 | Flag type dispatch |
| flow_context.go | 395 | Execution context |
| flags_validate.go | 345 | Flag validation |
| scope.go | 342 | DI scope |
| command.go | 340 | Command model |
| output.go | 325 | Output formatting |

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

#### Module 1: `cmdguard-types`

| Field | Content |
|-------|---------|
| **Name & path** | `/types` → `github.com/larsartmann/cmdguard/types` |
| **Purpose** | Validated semantic value types for CLI flags |
| **Dependencies (prod)** | None (zero external deps — stdlib only) |
| **Dependencies (test)** | `github.com/larsartmann/cmdguard/testutil` |
| **Public API** | `Duration`, `Email`, `Enum`, `FilePath`, `HostPort`, `LogLevel`, `LogFormat`, `Port`, `URL` + `Parse`/`MustParse` constructors + type-specific error sentinels + `EnumError`, `DurationError` + `MustParse[T]` helper |
| **Internal packages** | None |
| **External deps** | stdlib only |

**Includes from current v2:**

- `types_duration.go` → `duration.go`
- `types_email.go` → `email.go`
- `types_enum.go` → `enum.go`
- `types_filepath.go` → `filepath.go`
- `types_hostport.go` → `hostport.go`
- `types_log.go` → `log.go`
- `types_port.go` → `port.go`
- `types_url.go` → `url.go`
- `type_helpers.go` → `helpers.go` (only `MustParse[T]`; `Ptr`, `ValueOrDefault`, `EnsureValid` are dead code — delete)
- Error sentinels: `ErrInvalidURL`, `ErrInvalidEmail`, `ErrInvalidPort`, `ErrInvalidFilePath`, `ErrInvalidHostPort`, `ErrInvalidEnum`, `ErrInvalidDuration` + error types `EnumError`, `DurationError` + constructors
- Dead code to delete: `ErrLogLevel`, `ErrLogFormat` (never used)

**Note:** `type_handler.go` stays in core. Core imports types and registers handlers via `RegisterTypeHandler()`. The types module doesn't know about the handler registry.

#### Module 2: `cmdguard-output`

| Field | Content |
|-------|---------|
| **Name & path** | `/output` → `github.com/larsartmann/cmdguard/output` |
| **Purpose** | Multi-format output rendering |
| **Dependencies (prod)** | None internal |
| **Dependencies (test)** | `github.com/larsartmann/cmdguard/testutil` |
| **Public API** | `OutputFormat`, `OutputConfig`, `ParseOutputFormat`, `DefaultOutputConfig`, `OutputResult`, `OutputTable`, `OutputStyledTable`, format constants |
| **Internal packages** | None |
| **External deps** | `github.com/larsartmann/go-output` |

**Includes from current v2:**

- `output.go` (full file moves)
- Error sentinels: `ErrUnsupportedFormat`, `ErrFormatRequiresTypedData`

**Does NOT include:** `cli_output.go` — it defines methods on `CLI[T]` and uses `AddGlobalFlag()`, so it stays in core.

#### Module 3: `cmdguard` (core — root module)

| Field | Content |
|-------|---------|
| **Name & path** | `/` (root) → `github.com/larsartmann/cmdguard` |
| **Purpose** | CLI framework core — CLI[T], Command[T,F], flag system, DI scope, middleware, flow context, cobra integration |
| **Dependencies (prod)** | `cmdguard-types`, `cmdguard-output` |
| **Dependencies (test)** | `github.com/larsartmann/cmdguard/testutil` |
| **Public API** | All CLI/Command types, flag system, scope/DI, middleware, flow context, type handler registry, editor, manpage generation |
| **Internal packages** | None |
| **External deps** | cobra, pflag, samber/do/v2, fang/v2, mango, mango-cobra, roff |

**Stays in core (not extracted):**

- `scope.go` — `Package[T]` creates circular dependency if extracted (calls `NewCLI` and `WithCLIScope`). DI scope is deeply integrated with CLI construction.
- `type_handler.go` — registers handlers for all types; must import types module
- `cli_output.go` — methods on `CLI[T]`, must stay in core
- `config_setfield.go` — references `Enum`, `Duration` directly (imports from types module)
- All errors shared between core and types stay in core; types module defines its own copies

#### Module 4: `cmdguard-testutil`

| Field | Content |
|-------|---------|
| **Name & path** | `/testutil` → `github.com/larsartmann/cmdguard/testutil` |
| **Purpose** | Shared test assertion helpers |
| **Dependencies (prod)** | None |
| **Dependencies (test)** | None |
| **Public API** | `AssertEqual`, `AssertNil`, `AssertPanics`, `AssertFlagRegistered`, `NoOpCobraRun`, etc. |
| **External deps** | `github.com/spf13/cobra` (for cobra-specific test helpers) |

### 3.2 Dependency DAG (Proposed)

```
    ┌──────────────┐
    │ cmdguard     │ ← Core CLI framework
    │ (root)       │    (includes DI scope, type handler,
    └──┬───────┬───┘    middleware, flow context)
       │       │
       ▼       ▼
  ┌────────┐  ┌──────────────┐
  │cmdguard│  │ cmdguard-    │
  │ types  │  │ output       │
  │        │  │              │
  └────────┘  └──────────────┘
       │             │
       ▼             ▼
    (stdlib       (go-output)
     only)

  ┌──────────────────┐
  │ cmdguard-testutil│ ← Standalone (leaf)
  └──────────────────┘
       │
       ▼
    (cobra for test helpers)
```

**DAG verification:**

- `cmdguard` → `types`, `output` (downward, no cycles)
- `types` → stdlib only (leaf)
- `output` → go-output only (leaf)
- `testutil` → cobra only (leaf, test-only dependency)
- **No cycles. No upward dependencies. No bidirectional deps.**

### 3.3 Why DI Scope Stays in Core

**Self-review finding:** `Package[T]` in `scope.go:322` calls both `NewCLI[T]()` and `WithCLIScope[T]()`, creating a circular import:

```
scope.go → NewCLI (cli.go) → *Scope field (scope.go)  ← CIRCULAR
```

Options considered:

1. ~~Extract scope, move `Package[T]` to core~~ — Leaves scope module as just CRUD wrappers around samber/do, not worth a separate module
2. ~~Extract scope, use interface/registry pattern~~ — Over-engineering for a convenience function
3. **Keep scope in core** — Simplest, preserves all functionality, no circular deps

**Decision: Keep scope in core.** The DI scope is deeply integrated with CLI construction (it's a field on `CLI[T]`). Extracting it would require complex interface gymnastics for minimal benefit. The scope's only external dep is `samber/do/v2`, which is already in core's go.mod.

### 3.4 Shared Error Sentinel Strategy

**Self-review finding:** Several error sentinels are used by both types files and core:

| Sentinel | Types usage | Core usage |
|----------|-------------|------------|
| `ErrInvalidURL` | `types_url.go` | `flags_validate.go` |
| `ErrInvalidEmail` | `types_email.go` | `flags_validate.go` |
| `ErrInvalidEnum` | `types_enum.go` | `flags.go` |
| `ErrInvalidDuration` | `types_duration.go` | (tests only) |

**Decision:** The types module defines its own error sentinels. Core defines its own copies for validation. This is acceptable because:

- The sentinels have the same string message
- `errors.Is()` works by identity (pointer), not value — so they're different errors
- Core's validation checks (`flags_validate.go`) already create their own error messages; they don't need to match the types module's sentinel exactly
- Alternatively, core can import the types module's errors (core → types is already in the DAG)

**Preferred approach:** Core imports types module's error sentinels. This avoids duplication and is clean since core already depends on types. The types module owns its error sentinels; core uses them via import.

### 3.5 Replace / Workspace Strategy

**Chosen: `go.work` at repository root.**

```
// go.work
go 1.26.2

use (
    .
    ./types
    ./output
    ./testutil
)
```

**Rules:**

- Never mix `replace` directives AND `go.work` for the same module pair
- `go.work` stays in version control (all modules published together)
- Verify `go mod tidy` works both with and without workspace

### 3.6 Versioning Strategy

**Chosen: Shared version (single git tag `v2.x.x`).**

Rationale:

- Single-team library with tight coupling between core and types/output
- All modules publish together from the same repo
- Consumers import core primarily; types/output are transitive
- Independent semver adds complexity without benefit at this scale

**Tag format:** `v2.3.0` (next release after modularization)

### 3.7 Import Path Impact

| Before | After | Breaking? |
|--------|-------|-----------|
| `github.com/larsartmann/cmdguard/pkg/cmdguard/v2` | Same — unchanged | **No** |
| `github.com/larsartmann/cmdguard/pkg/testutil` | Same — but now separate module | **No** (import path unchanged) |

**New import paths (additive, not replacing):**

- `github.com/larsartmann/cmdguard/types` — value types
- `github.com/larsartmann/cmdguard/output` — output formatting

**Zero breaking changes for existing consumers.** The v2 package re-exports types from sub-modules for backward compatibility.

---

## 4. Self-Review Findings

### 4.1 What We Forgot Initially

1. **`type_handler.go` coupling** — It references every single value type. Must stay in core but imports types module. Registration pattern works because core depends on types (correct DAG direction).

2. **Shared error sentinels** — 4 sentinels are used by both types and core. Resolved by having core import types module's errors.

3. **`Package[T]` circular dependency** — Prevents scope extraction. Kept in core.

4. **Dead code** — `ErrLogLevel`, `ErrLogFormat`, `Ptr`, `ValueOrDefault`, `EnsureValid` are unused. Delete during modularization.

5. **`config_setfield.go`** — Directly imports `Enum`, `Duration`, `ParseEnum`, `FromDuration`. After extraction, these come from the types module import — clean DAG.

6. **`cli_output.go` cannot move with output** — It defines methods on `CLI[T]`. Only `output.go` (rendering logic) moves.

### 4.2 What Could Be Improved

1. **File sizes** — `type_handler.go` (480 lines) and `flow_context.go` (395 lines) exceed the 300-line rule. These should be split regardless of modularization.

2. **Dead code cleanup** — Should happen alongside modularization to avoid carrying unused code into new modules.

3. **`FlagTag` struct location** — Defined in `config.go`, used by `type_handler.go`. Both stay in core, so no issue. If we ever wanted to extract type_handler, we'd need to move FlagTag too.

### 4.3 Cross-Reference with how-to-golang

| Check | Result |
|-------|--------|
| No banned dependencies | ✅ All deps clean (samber/do/v2, cobra, pflag, fang/v2, go-output) |
| Files over 300 lines | ⚠️ 6 files exceed limit — address in execution plan |
| No `any` types | ✅ Uses generics throughout |
| No magic strings | ✅ Constants used |
| Clean root | ⚠️ Root has many stray files (PARTS.md, PROGRESS_2026-04-01.md, etc.) — out of scope |
| Latest Go version | ✅ Go 1.26.2 |

---

## 5. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Shared error sentinels create split brains | Low | Medium | Core imports types module's errors (correct DAG direction) |
| `RegisterTypeHandler` global state | Low | Low | Stays in core; types module has no global state |
| go.work conflicts with consumers | Low | Low | go.work is ignored by consumers |
| Test breakage from package splits | Medium | Low | Fix immediately after each step |
| Re-export maintenance burden | Low | Low | Re-exports are thin type aliases/wrappers |
| Scope extraction attempted (circular deps) | N/A | N/A | Decided to keep scope in core |

---

## 6. Build System Impact

### After Modularization

```bash
# Workspace-aware (default in repo):
go work sync
go test ./... -count=1 -timeout 120s -race

# Per-module (CI isolation):
cd types && go test ./... -race
cd output && go test ./... -race
cd testutil && go test ./... -race
go test ./... -race  # core + examples + integration

# Lint:
golangci-lint run ./...
```

### flake.nix

Update to build each module independently. Not blocking — can be done after.

---

## 7. Key Decisions

1. **go.work over replace directives** — Cleaner, native tooling support
2. **Shared versioning** — Single-team library, tight coupling
3. **Preserve `pkg/cmdguard/v2` import path** — Zero breaking changes
4. **DI scope stays in core** — `Package[T]` creates circular dependency if extracted
5. **Core imports types module's errors** — Avoids sentinel duplication
6. **type_handler stays in core** — It dispatches based on reflect types from all value types
7. **types module has zero external deps** — Maximum reusability
8. **cli_output.go stays in core** — It's CLI infrastructure, not output rendering
9. **Delete dead code during migration** — `ErrLogLevel`, `ErrLogFormat`, `Ptr`, `ValueOrDefault`, `EnsureValid`

---

## 8. What We're NOT Doing

- **Not extracting DI scope** — Circular dependency risk, deeply integrated with CLI
- **Not extracting type_handler** — Tightly coupled to both core config and all value types
- **Not splitting examples into modules** — Demonstration code, not importable
- **Not splitting benchmarks** — They test the core library
- **Not splitting integration tests** — They test the complete system
- **Not changing the public API** — All existing types, functions, and options remain available from `pkg/cmdguard/v2`
- **Not creating `internal/` packages** — Everything remains publicly importable
