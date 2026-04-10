# Cross-Library Integration Review: Pro/Contra for cmdguard

**Date:** 2026-04-10
**Status:** Complete Analysis
**Reviewer:** AI Partner
**Context:** cmdguard v2.1.0 — type-safe Cobra CLI library with DI

---

## Executive Summary

9 libraries were evaluated for integration potential with cmdguard. The evaluation considers: **domain fit** (does it solve a problem cmdguard users have?), **dependency cost**, **maturity**, and **maintainer alignment** (all are Lars's own projects).

**TL;DR — Recommendations:**

| Library | Verdict | Priority |
|---|---|---|
| **go-output** | **INTEGRATE — High Priority** | P0 |
| **go-business-rules** | **INTEGRATE — Medium Priority** | P1 |
| **go-filewatcher** | **INTEGRATE — Low Priority** | P2 |
| **gogenfilter** | **CONSIDER — Niche Use** | P3 |
| **universal-workflow** | **DO NOT INTEGRATE** | — |
| **go-cqrs-lite** | **DO NOT INTEGRATE** | — |
| **go-localfirst** | **DO NOT INTEGRATE** | — |
| **go-localsync** | **DO NOT INTEGRATE** | — |
| **go-plugin-mvp** | **DO NOT INTEGRATE** | — |

---

## Detailed Pro/Contra Per Library

---

### 1. go-output

> *Multi-format output rendering (12 formats: table, JSON, CSV, Markdown, Mermaid, DOT, etc.)*

**Module:** `github.com/larsartmann/go-output`
**Dependencies:** charm.land/lipgloss/v2, go-faster/yaml (minimal)
**Maturity:** Well-tested (unit, integration, fuzz, benchmarks), MIT licensed
**Existing cmdguard bridge:** YES — `go-output/cmdguard/` subpackage already exists

#### PRO

1. **Direct cmdguard bridge already exists** — `go-output/cmdguard/` provides `EnumFlag[T]`, `OutputFormatFlag`, `ColorModeFlag`, `SortByFlag` — these are purpose-built for cmdguard's flag system
2. **Every CLI needs output formatting** — this is a universal concern for CLI applications built with cmdguard
3. **Zero coupling** — the `cmdguard/` subpackage does NOT import cmdguard itself; it provides compatible types via shared interfaces
4. **12 output formats** — table, JSON, CSV, TSV, Markdown, XML, YAML, D2, Mermaid, DOT, HTML, tree — covers virtually any CLI output need
5. **Enum system with CLI-first design** — every enum type has `String()`, `AllowedValues()`, `IsValid()`, `Parse*()` — designed for flag parsing and shell completion
6. **ColorMode with TTY/NO_COLOR/CI detection** — solves a common CLI concern out of the box
7. **Lightweight dependencies** — only lipgloss (terminal styling) and yaml
8. **Phantom-typed IDs** — `BrandedID[Brand]` prevents mixing different ID types at compile time

#### CONTRA

1. **Adds 2 direct dependencies** (lipgloss, go-faster/yaml) to the cmdguard ecosystem — users who don't need output formatting pay the cost
2. **Integration mode unclear** — should cmdguard depend on go-output? Or should go-output's `cmdguard/` package be the only bridge? A circular dependency is architecturally messy
3. **Scope creep risk** — cmdguard is a CLI framework, not an output framework. Adding output formatting could blur the project's identity
4. **No semver tag** — go-output has no versioned release yet, risking breaking changes

#### Verdict: **INTEGRATE (P0)**

The existing `cmdguard/` bridge in go-output proves this was designed for integration. The recommended approach: keep the bridge in go-output (as-is), but **document it prominently in cmdguard's README** and add an example showing `go-output` + `cmdguard` together. Do NOT add a dependency from cmdguard → go-output.

---

### 2. go-business-rules

> *Severity-aware validation framework (Info/Warning/Error/Critical) with fluent rule builders*

**Module:** `github.com/artmann/businessrules` (v1.1.0)
**Dependencies:** ZERO runtime deps (only ginkgo/gomega for tests)
**Maturity:** v1.1.0, extensive tests, clean architecture

#### PRO

1. **Zero runtime dependencies** — perfect for a library like cmdguard that values minimal deps
2. **Solves a real CLI problem** — validating user input, config files, and command arguments with graduated severity (block on errors, warn on warnings, inform on info)
3. **Severity-aware** — unlike standard validators that return pass/fail, this returns violations tagged with severity levels. A CLI can decide to block on errors but proceed with warnings
4. **Pre-built rule constructors** — Numeric, String, Collection, Format, Generic, Composite rules — covers common validation needs
5. **JSON-serializable results** — `ViolationError` and `ValidationResultError` both implement `json.Marshaler`, enabling structured output
6. **Composable** — `All`, `Any`, `When`, `Custom` allow building complex validation chains
7. **Fluent builder** — `NewValidator() → AddRule() → Build()` pattern fits cmdguard's functional options style
8. **Small, focused library** — does one thing well

#### CONTRA

1. **No existing cmdguard integration** — would need new bridge code
2. **Validation is already possible in cmdguard** — users can validate in `PreRunE` hooks with stdlib or any validator. The question is whether cmdguard should provide a recommended validation approach
3. **Naming mismatch** — "business rules" sounds domain-heavy, but the library is actually a general-purpose validation framework. The name might confuse potential users
4. **Added complexity** — cmdguard users who don't need validation would see it as bloat if it becomes a core dependency

#### Verdict: **INTEGRATE (P1)**

Not as a dependency, but as an **optional companion**. Add documentation showing how to use go-business-rules with cmdguard's `PreRunE` hooks. Consider adding a `WithValidation()` option that accepts `Rule` interface. The zero-dependency profile makes this a clean optional integration.

---

### 3. go-filewatcher

> *High-level, composable file system watcher built on fsnotify*

**Module:** `github.com/larsartmann/go-filewatcher`
**Dependencies:** fsnotify (1 direct dependency)
**Maturity:** Well-tested, proprietary license

#### PRO

1. **Extremely minimal dependencies** — only fsnotify
2. **Solves a real CLI pattern** — `watch`/`dev` commands are common in CLI tools (e.g., `go run --watch`, `npm watch`, `air`)
3. **Composable filters** — AND/OR/NOT logic, extensions, globs, regex, min-size
4. **Middleware chain** — logging, recovery, rate-limit, metrics
5. **Debounce** — global and per-path debounce, critical for file watching
6. **Functional options** — matches cmdguard's style exactly
7. **Context-based shutdown** — clean lifecycle management
8. **Phantom types** — `Op` uses branded types for type safety

#### CONTRA

1. **Proprietary license** — cannot be used in open-source projects without permission. This is a blocker for integration with MIT-licensed cmdguard
2. **Niche use case** — only relevant for CLI tools that need file watching (watch mode, hot reload, file processing). Not a universal CLI need
3. **fsnotify is the real dependency** — many projects just use fsnotify directly; go-filewatcher's value-add is the composable filter/middleware system
4. **Example usage in cmdguard would be** a `watch` command that re-runs a handler on file changes — useful but not core to cmdguard's mission

#### Verdict: **INTEGRATE (P2)**

Good fit for a specific use case, but the proprietary license is a concern. If relicensed to MIT, this would be an excellent companion library for building `watch`/`dev` commands with cmdguard. Document as an optional integration pattern.

---

### 4. gogenfilter

> *Detect and filter auto-generated Go code files (sqlc, templ, protobuf, mockgen, etc.)*

**Module:** `github.com/larsartmann/gogenfilter`
**Dependencies:** go-faster/yaml (1 direct dependency)
**Maturity:** Pre-release, extensive tests, proprietary license

#### PRO

1. **Near-zero dependencies** — only YAML parsing for sqlc config discovery
2. **Two-phase detection** — filename-based first (zero I/O), then content-based (reads file). Performance-conscious design
3. **Comprehensive coverage** — sqlc, templ, go-enum, protobuf, mockgen, stringer, and generic `// Code generated by` detection
4. **fs.FS abstraction** — `Filter.WithFS()` makes it testable and embeddable
5. **Structured errors with help text** — `ErrorCode`, `Help()`, `errors.Is` support
6. **Thread-safe metrics** — `FilterStats`, `MetricsMixin` for tracking filter decisions
7. **Pattern matching** — include/exclude glob patterns for custom filtering
8. **Useful for CLI tools** — if cmdguard ever processes Go source files (e.g., a codegen command), this would be essential

#### CONTRA

1. **Proprietary license** — same blocker as go-filewatcher
2. **Pre-release** — no versioned tag, `[Unreleased]` in CHANGELOG
3. **Niche use case for cmdguard** — cmdguard is a CLI framework, not a code analysis tool. Generated file detection is only relevant if cmdguard adds code generation or analysis features
4. **No CI workflow** — noted as TODO in the project
5. **No direct cmdguard connection** — would need a specific feature in cmdguard that processes Go source files

#### Verdict: **CONSIDER (P3)**

Well-engineered but very niche for cmdguard. Only relevant if cmdguard adds code generation or Go-source-analysis features. Keep on the radar, but don't integrate now.

---

### 5. universal-workflow

> *DAG-based workflow orchestration with dependency resolution, parallel execution, and visualization*

**Module:** `github.com/LarsArtmann/universal-workflow`
**Dependencies:** Heavy — bubbletea, lipgloss, cobra, koanf, watermill, samber/do, samber/mo, cockroachdb/errors, graph, otel, posthog, gorilla/websocket, ginkgo/gomega
**Maturity:** v1.0.0, 685+ tests, extensive documentation

#### PRO

1. **Powerful DAG execution** — Kahn's algorithm for topological sort, parallel and sequential modes
2. **Rich feature set** — retry with backoff, conditional branching, security levels, flight recorder, OTEL tracing, Mermaid/DOT export
3. **Well-tested** — 87 test files, 685+ tests
4. **Type-safe** — branded IDs, generic Result[T], phantom types
5. **Could orchestrate multi-command CLI workflows** — e.g., a CLI that runs multiple steps with dependencies

#### CONTRA

1. **Massive dependency footprint** — 15+ direct dependencies including bubbletea, watermill, posthog, otel, websocket. This would explode cmdguard's dependency tree
2. **Completely different domain** — cmdguard is a CLI framework for building command-line tools; universal-workflow is a workflow orchestration engine. They solve different problems
3. **CLI was explicitly removed** — the project has docs titled "CLI_REMOVAL_COMPLETE" and "LIBRARY-FIRST_ARCHITECTURE_ELIMINATING_CLI". It pivoted away from CLI
4. **Overkill for CLI commands** — most CLI commands are simple: parse args → validate → execute → output. DAG-based orchestration is overkill for 99% of CLI use cases
5. **~45% test coverage** — despite 685+ tests, coverage is low due to the massive codebase
6. **Bubble Tea TUI dependency** — terminal visualization is cool but not needed in a CLI framework

#### Verdict: **DO NOT INTEGRATE**

Wrong domain, too heavy. If a cmdguard user needs workflow orchestration, they can import universal-workflow themselves. There's no natural bridge between a CLI flag framework and a DAG workflow engine.

---

### 6. go-cqrs-lite

> *Lightweight CQRS and Event Sourcing library with branded IDs and catalog generation*

**Module:** `github.com/larsartmann/go-cqrs-lite`
**Dependencies:** 3 direct (google/uuid, cockroachdb/errors, go-json-experiment/json)
**Maturity:** Production-ready, comprehensive tests (unit, integration, benchmark, fuzz)

#### PRO

1. **Minimal dependencies** — only 3 direct deps
2. **Well-designed** — clean separation of command/query/event/aggregate, middleware pipeline, branded IDs
3. **AsyncAPI + EventCatalog generation** — auto-generates documentation from Go struct types via reflection
4. **Type-safe IDs** — `id.Of[T]` branded generics prevent mixing AggregateID with EventID at compile time
5. **Event sourcing with replay** — `LoadFromHistory`, `EventSourcedRepository`
6. **Excellent code quality** — 250-line file limit, strict linting, comprehensive test types

#### CONTRA

1. **Completely different domain** — CQRS/Event Sourcing is a backend architecture pattern. cmdguard is a CLI framework. No overlap
2. **No CLI integration point** — the library has no cobra/CLI integration. Its entry points are HTTP handlers
3. **Added complexity for zero benefit** — importing CQRS into a CLI framework makes no sense unless cmdguard becomes a full application framework
4. **Catalog generation is the closest feature** — auto-generating API docs from Go types is interesting but not relevant to CLI flag parsing

#### Verdict: **DO NOT INTEGRATE**

Excellent library, wrong home. go-cqrs-lite belongs in backend services, not CLI frameworks.

---

### 7. go-localfirst

> *Local-first application framework with CRDT sync, event sourcing, and embedded Pebble storage*

**Module:** `github.com/larsartmann/go-localfirst`
**Dependencies:** Heavy — gin, pebble, casbin, websocket, prometheus, samber/do, jwt, htmx, zap
**Maturity:** Reference implementation, not a reusable library

#### PRO

1. **pkg/sync is extractable** — vector clocks, conflict resolution, LWW merge as a zero-dependency package
2. **Good architecture** — clean handler → service → storage layers
3. **CRDT primitives are well-designed** — `VectorClock`, `Operation[T]`, `ConflictResolver[T]`

#### CONTRA

1. **Not a library — it's a reference application** — it's a full Todo app with HTTP handlers, HTMX UI, Docker, auth middleware. Not designed to be imported
2. **Massive dependency tree** — gin, pebble, casbin, websocket, prometheus, jwt. Absolutely inappropriate for a CLI framework
3. **HTTP/REST application** — fundamentally a web server, not a library
4. **go.mod uses `replace` directive** — `replace github.com/larsartmann/go-cqrs-lite => ../go-cqrs-lite` shows it's a local development project, not a published module
5. **Zero relevance to CLI** — local-first sync and CRDTs solve distributed data problems. CLI tools run locally

#### Verdict: **DO NOT INTEGRATE**

go-localfirst is a reference application, not an importable library. The `pkg/sync` package could theoretically be extracted, but it solves distributed data problems that are irrelevant to CLI tools.

---

### 8. go-localsync

> *Pluggable SDK for syncing data from external providers (GitHub, etc.) to local SQLite*

**Module:** `github.com/larsartmann/go-localsync`
**Dependencies:** sqlite, go-github, go-localfirst, go-composable-business-types, ginkgo/gomega
**Maturity:** SDK, depends on private repos

#### PRO

1. **Clean Provider interface** — pluggable data sources
2. **SQLite storage with full JSON fidelity** — good for offline-first apps
3. **Conflict-aware sync** — CRDT-backed via go-localfirst
4. **Branded IDs** — type-safe identifiers

#### CONTRA

1. **Depends on private repositories** — `go-composable-business-types` and `go-localfirst` are private. Would block cmdguard users
2. **Completely different domain** — data synchronization from external APIs. Zero overlap with CLI frameworks
3. **SQLite dependency** — `modernc.org/sqlite` adds significant binary size. Inappropriate for a CLI framework
4. **Specific use case** — syncing GitHub events to local storage. Not general-purpose enough for cmdguard integration

#### Verdict: **DO NOT INTEGRATE**

Solves a completely different problem. No integration point with cmdguard.

---

### 9. go-plugin-mvp

> *WebAssembly plugin system with event-driven architecture*

**Module:** `github.com/larsartmann/go-plugin-mvp` (v0.1.0)
**Dependencies:** extism/go-sdk, watermill, wazero, yaml
**Maturity:** MVP/reference implementation (v0.1.0)

#### PRO

1. **Plugin system concept is relevant to CLIs** — many CLI tools support plugins (kubectl, gh, etc.)
2. **WASM sandboxing** — memory isolation, configurable resource limits
3. **Bidirectional host↔plugin communication** — host functions (KV, logging, events) callable from plugins
4. **Event-driven architecture** — pub/sub with Watermill
5. **Well-structured for an MVP** — typed errors, validation, CI pipeline

#### CONTRA

1. **MVP/proof-of-concept** — v0.1.0, not production-ready. Premature to integrate
2. **Heavy dependencies** — extism, wazero (WASM runtime), watermill. Would add ~50+ transitive dependencies to cmdguard
3. **WASM compilation required** — plugins must be compiled to WASM with TinyGo. Significant developer friction
4. **Over-engineered for cmdguard's scope** — cmdguard is a CLI framework, not a plugin system. If a user needs plugins, they can use go-plugin-mvp independently
5. **No natural integration point** — cmdguard doesn't have a plugin architecture, and adding one would be a fundamental scope expansion
6. **Maintained in `internal/`** — most code is in internal packages, not designed for external consumption

#### Verdict: **DO NOT INTEGRATE**

Interesting concept but premature and out of scope. If cmdguard ever adds a plugin system, go-plugin-mvp could inspire the design, but it should not be a dependency.

---

## Analysis Work Status

### a) FULLY DONE

| Task | Status |
|---|---|
| go-output exploration and analysis | Complete — README, go.mod, source files, cmdguard/ bridge |
| universal-workflow exploration and analysis | Complete — README, go.mod, architecture, test coverage |
| go-business-rules exploration and analysis | Complete — README, go.mod, all source files |
| go-cqrs-lite exploration and analysis | Complete — README, go.mod, all packages |
| go-localfirst exploration and analysis | Complete — README, go.mod, pkg/sync, architecture |
| go-localsync exploration and analysis | Complete — README, go.mod, provider/storage interfaces |
| go-plugin-mvp exploration and analysis | Complete — README, go.mod, plugin architecture, WASM |
| go-filewatcher exploration and analysis | Complete — README, go.mod, options, filters, middleware |
| gogenfilter exploration and analysis | Complete — README, go.mod, filter/detection logic |
| Pro/Contra synthesis for all 9 libraries | Complete |

### b) PARTIALLY DONE

Nothing partially done — all analyses are complete.

### c) NOT STARTED

| Task | Notes |
|---|---|
| Actual integration code for go-output | Report recommends documenting, not code dependency |
| Actual integration code for go-business-rules | Report recommends documenting pattern, not code dependency |
| README updates in cmdguard for companion libraries | Follow-up task |
| Example code showing cmdguard + go-output usage | Follow-up task |

### d) TOTALLY FUCKED UP

- **Agent tool failures** — 5 of 9 initial agent calls failed with "error generating response". Resolved by using direct View/LS tools and single-agent calls instead.
- **No data corruption or analysis errors** — all final analyses are based on direct file reads.

### e) WHAT WE SHOULD IMPROVE

1. **go-output already has a cmdguard bridge** (`go-output/cmdguard/`) — this should be documented in cmdguard's README with a usage example
2. **go-business-rules naming** — consider renaming to something like `go-severity-validation` or `go-graduated-validation` to better communicate its purpose
3. **go-filewatcher and gogenfilter licensing** — both are proprietary; if they're meant to be companion libraries for MIT-licensed cmdguard, they should be relicensed
4. **go-output needs a semver release** — no versioned tag makes it risky to recommend as a companion library
5. **gogenfilter has no CI** — needs CI before being recommended as a companion

### f) Top #25 Things We Should Get Done Next

1. Add "Output Formatting" section to cmdguard README documenting go-output integration
2. Add `examples/output-formats/` showing cmdguard + go-output together
3. Verify go-output's `cmdguard/` bridge works with cmdguard v2.1 API
4. Add "Validation" section to cmdguard README documenting go-business-rules integration
5. Add `examples/validation/` showing cmdguard + go-business-rules in PreRunE hooks
6. Consider adding `WithValidation(rules ...businessrules.Rule)` option to cmdguard
7. Tag go-output with v1.0.0 semver release
8. Investigate go-output circular dependency risk (cmdguard/ subpackage)
9. Add go-output and go-business-rules to cmdguard's FEATURES.md or similar
10. Add "File Watching" section to cmdguard README for go-filewatcher integration
11. Re-evaluate go-filewatcher if relicensed to MIT
12. Re-evaluate gogenfilter if cmdguard adds code generation features
13. Extract go-localfirst/pkg/sync into standalone module if it's meant for reuse
14. Consider a `cmdguard-ecosystem` meta-README documenting all companion libraries
15. Add integration tests verifying go-output cmdguard bridge compiles against cmdguard v2
16. Write ADR documenting the decision NOT to integrate the other 6 libraries
17. Review go-business-rules naming — "businessrules" is misleading for a general validation lib
18. Update cmdguard AGENTS.md with companion library information
19. Consider adding output format auto-detection (TTY → table, pipe → JSON) helper
20. Add ColorMode integration from go-output into cmdguard's global flags
21. Evaluate if go-output's table rendering should be a default for cmdguard commands
22. Consider shared error types between cmdguard and companion libraries
23. Add CI badge cross-links between companion library repos
24. Document minimum Go version compatibility across all companion libraries
25. Create a `cmdguard-showcase` repo demonstrating all companion integrations together

### g) Top #1 Question I Cannot Figure Out Myself

**What is the intended dependency direction between cmdguard and go-output?**

The `go-output/cmdguard/` bridge already exists, but:
- If cmdguard imports go-output → users who don't need output formatting pay the dependency cost
- If go-output imports cmdguard → circular concern (currently avoided since go-output/cmdguard/ uses compatible types, not imports)
- If neither imports the other → the bridge is documentation-only (current state)

Which model is the intended long-term architecture? This determines whether we add go-output as a dependency, keep it documentation-only, or create a separate `cmdguard-output` adapter module.
