# Go Modularization Assessment — cmdguard

**Date:** 2026-06-10
**Current State:** Single go.mod (monolith)

## 1. Current State

| Property     | Value                |
| ------------ | -------------------- |
| go.mod files | 1 (root)             |
| go.work      | None                 |
| Packages     | 7                    |
| Source files | 50 non-test, 80 test |
| Total lines  | 8,463                |

## 2. Should This Project Modularize?

| Signal                       | Weight | Present? | Notes                                                |
| ---------------------------- | ------ | -------- | ---------------------------------------------------- |
| Small project                | High   | **Yes**  | 7 packages, single domain                            |
| No external consumers        | Medium | **No**   | Published library at github.com/larsartmann/cmdguard |
| Prototype/spike              | High   | No       | v2.4.0, stable API                                   |
| All packages change together | High   | **Yes**  | Single package `v2` contains 90% of code             |

**Score: 2 High signals** → Per the go-modularize skill: "Consider partial modularization — extract only the core/domain module to establish a clean API surface."

## 3. Package Coupling Analysis

```
pkg/cmdguard/v2/          ← Core library (34 source files, flat package)
pkg/cmdguard/v2/configload/ ← Config loaders (imports parent)
pkg/cmdguard/v2/testutil/   ← Test helpers (imports parent)
pkg/testutil/              ← Generic test assertions (no internal deps)
examples/taskctl/          ← Example app (imports v2)
benchmarks/                ← Benchmarks (imports v2)
tests/integration/         ← Integration tests (imports v2)
```

**All packages depend on `pkg/cmdguard/v2` (one-way). No circular deps.**

## 4. Modularization Options

### Option A: No modularization (recommended for now)

The project is a **single-purpose library** with 7 packages in a clean hierarchy. All packages change together because they're all part of one library's API. The internal coupling is managed by Go's package system.

**Rationale:**

- Under 10 source packages
- Single domain (CLI framework)
- All packages change together
- Published as a single library — consumers want one `go get`
- Module split would add 3-5 go.mod files with no real benefit

### Option B: Partial split (future consideration)

If the library grows significantly, consider:

```
cmdguard/
├── go.work
├── core/                    # go.mod — types, interfaces, errors
│   └── types.go, errors.go, config.go
├── pkg/cmdguard/v2/         # go.mod — main library
│   └── depends on core/
└── configload/              # go.mod — optional config loaders
    └── depends on core/
```

**Benefits:** Consumers who only need types/errors can import `core/` without pulling in cobra/fang/huh/otel.
**Costs:** More go.mod management, potential semver drift, CI complexity.

### Option C: Full split (not recommended)

Splitting into 5+ modules (core, flags, types, middleware, output) would be over-modularization. These packages always change together and have no independent consumers.

## 5. Verdict

**Do NOT modularize now.** The project is too small and too tightly coupled to benefit. The single-package design (`v2`) is actually a strength — it keeps the API surface simple for consumers.

Revisit when:

- Package count exceeds 15
- Clear independent consumers emerge (e.g., someone only wants the type validation without cobra)
- The library exceeds 15,000 lines
