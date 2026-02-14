# TODO_LIST.md - cmdguard Project Tasks

**Last Updated:** 2026-02-14
**Purpose:** Comprehensive tracking of all TODOs, issues, and planned changes for the cmdguard project.

---

## Files Analyzed

| File | Status | Date Analyzed |
|------|--------|---------------|
| AGENTS.md | ✅ Done | 2026-02-14 |
| README.md | ✅ Done | 2026-02-14 |
| docs/ARCHITECTURE_REVIEW.md | ✅ Done | 2026-02-14 |
| docs/CLI_DESIGN_PRINCIPLES.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_05-28_DECISION_POINT.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_04-19_CLI_IMPROVEMENTS.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_05-44_POST_DECISION_STATE.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_04-25_COMPREHENSIVE_STATUS.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_04-11_BUILD_VERIFICATION_COMPLETE.md | ✅ Done | 2026-02-14 |
| docs/planning/2026-02-14_04-21_CMDGUARD_TRANSFORMATION_PLAN.md | ✅ Done | 2026-02-14 |
| docs/status/2026-02-14_03-48_INITIAL_IMPLEMENTATION_COMPLETE.md | ✅ Done | 2026-02-14 |

---

## Key Decisions (From DECISION_POINT.md)

| Decision | Status | Notes |
|----------|--------|-------|
| Keep framework approach (Option A) | ✅ ADOPTED | Not rebuilding as guard library |
| README describes framework | ✅ DONE | API reference present |
| Documentation aligned with implementation | ✅ DONE | Claims verified |

---

### 1. Errcheck Violations [NOT_DONE]
**Source:** AGENTS.md (verified in code)
**Files:**
- `cmd/cmdguard/main.go:64,80,90,103`
- `internal/commands/root.go:142,155`
**Action:** Add error checking for `fmt.Fprintln`/`fmt.Fprintf` calls

### 2. False Documentation Claims [NOT_DONE]
**Source:** README.md (lines 164-165)
**Claims:**
- "No duplicate command names" validation - NOT IMPLEMENTED
- "No conflicting aliases" validation - NOT IMPLEMENTED
**Action:** Either implement these validations or remove from documentation

### 3. Broken Code Examples in Documentation [NOT_DONE]
**Source:** README.md (lines 45-77, 81-95)
**Details:** Code examples missing `"fmt"` import, will not compile
**Action:** Fix examples to be copy-pasteable

### 4. WithValidationHook is a No-Op [NOT_DONE]
**Source:** README.md (lines 154-155) → `pkg/cmdguard/public_api.go:130-135`
**Details:** Hook is documented but implementation discards it (returns `nil`)
**Action:** Either implement or remove from API
**Related:** Lower Priority "Design validation hook interface", "Allow custom validators"

### 5. Architecture Diagram Outdated [NOT_DONE]
**Source:** POST_DECISION_STATE.md (Technical Debt)
**Details:** Architecture diagram doesn't include logging layer
**Action:** Update architecture diagram with logging layer

### 6. Incomplete Context Handling [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P1)
**File:** `internal/di/module.go:123-127`
**Details:** Context passed to `Shutdown()` is completely ignored
**Action:** Implement proper context-aware shutdown

### 7. Silent Config Errors [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P1)
**File:** `internal/config/provider.go:32-35`
**Details:** Errors discarded with `_` assignment
**Action:** Properly handle/return config loading errors

### 8. No Error Aggregation [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P2)
**File:** `internal/di/module.go:115-117`
**Details:** Returns only `errs[0]` instead of using `errors.Join()`
**Action:** Aggregate multiple errors with `errors.Join()`

**Note:** Items "Code Duplication" and "Manual DI Wiring" were moved to Medium/High Priority sections to avoid duplication.

---

## Disproved Issues (Remove from Documentation)

| Claim | Reality | Action |
|-------|---------|--------|
| Unused `log/slog` import in root.go:7 | `slog` IS used at lines 23, 141, 154 | Remove from AGENTS.md |
| Missing `charmbracelet/log` dependency | `logger.go` only imports stdlib (`log/slog`, `os`) | Remove from AGENTS.md |
| "No Structured Logging" (INITIAL_IMPLEMENTATION_COMPLETE.md) | `internal/logging/logger.go` exists with full `slog` implementation | Mark as historical error |

---

## Missing Documentation (Exists but not documented)

| Feature | Location | Action |
|---------|----------|--------|
| `app.IsStrictMode()` | `pkg/cmdguard/public_api.go:191` | Add to README API Reference |
| `app.AddCommand()` | `pkg/cmdguard/public_api.go:212` | Add to README API Reference |
| Built-in `validate` command | Auto-added by `SetupCommands()` | Document in README |
| Built-in `version` command | Auto-added by `SetupCommands()` | Document in README |
| API usage examples | README | Add examples section per POST_DECISION_STATE.md |
| CONTRIBUTING.md | Project root | Create per POST_DECISION_STATE.md |
| Version tagging strategy | Documentation | Define per POST_DECISION_STATE.md |

---

## Superseded by Option A Decision

| Task | Original Location | Reason Obsolete |
|------|-------------------|-----------------|
| Remove `cmd/` folder | High Priority | Framework needs CLI entry point |
| Redesign public API | High Priority | Framework pattern is being kept |
| Add compile-time validation | High Priority | Guard-library specific feature |
| Define guard architecture in ARCHITECTURE.md | Medium Priority | Framework approach doesn't need separate guard architecture |

---

## Planned Changes (From Transformation Plan)

### High Priority - Core Architecture

| Task | Status | Notes |
|------|--------|-------|
| Fix DI usage | ❌ NOT_DONE | Manual wiring via `SetValidator()` exists |
| Add justfile | ❌ NOT_DONE | No justfile found |

### Medium Priority - Code Quality

| Task | Status | Notes |
|------|--------|-------|
| Improve test coverage (80%+) | 🔄 PARTIAL | config at ~47.6% per plan |
| Fix 4 errcheck violations | ❌ NOT_DONE | See Critical Issues #1 |
| Fix code duplication (7 clone groups) | ❌ NOT_DONE | Task 3.3 |
| Add integration tests | ❌ NOT_DONE | Task 3.4 |
| Create examples directory | ❌ NOT_DONE | Task 3.5 |
| Add command registry tests | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Add logging output tests | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Add flag type validation | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Add command dependency validation | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Implement strict mode for missing flags | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Add error context to validation errors | ❌ NOT_DONE | BUILD_VERIFICATION_COMPLETE.md |
| Add context variables to error messages | ❌ NOT_DONE | BUILD_VERIFICATION_COMPLETE.md |

### Lower Priority - Future Work

| Task | Status | Notes |
|------|--------|-------|
| Add CI/CD workflow | ❌ NOT_DONE | Task 4.1 |
| CI/CD Go version matrix testing | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Benchmarks | ❌ NOT_DONE | Beyond 80% |
| Fuzz testing | ❌ NOT_DONE | Beyond 80% |
| Release automation | ❌ NOT_DONE | Beyond 80% |
| Go module versioning | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Release notes automation | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Container image | ❌ NOT_DONE | Beyond 80% |
| Homebrew formula | ❌ NOT_DONE | Beyond 80% |
| Plugin system | ❌ NOT_DONE | Beyond 80% |
| Middleware support for commands | ❌ NOT_DONE | POST_DECISION_STATE.md Plugin System |
| Design validation hook interface | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Allow custom validators | ❌ NOT_DONE | POST_DECISION_STATE.md |
| Telemetry | ❌ NOT_DONE | Beyond 80% |
| Metrics | ❌ NOT_DONE | Beyond 80% |
| Tracing | ❌ NOT_DONE | Beyond 80% |
| Security audit | ❌ NOT_DONE | Beyond 80% |
| Complete test suite | ❌ NOT_DONE | Beyond 80% |
| Validation caching | ❌ NOT_DONE | Skip re-validation if tree unchanged |
| Parallel validation | ❌ NOT_DONE | For large command trees |
| Lazy config loading | ❌ NOT_DONE | Only load when accessed |
| Path traversal validation | ❌ NOT_DONE | For config file paths |
| Refactor: slices.Contains | ❌ NOT_DONE | root.go:52 loop simplification |
| Phantom type recommendations | ❌ NOT_DONE | Optional: 12 violations in BUILD_VERIFICATION |

### Suggested API Options (Not Implemented)

| Option | Status | Notes |
|--------|--------|-------|
| `WithConfigFile(path string)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `WithLogger(logger *slog.Logger)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `WithShutdownTimeout(timeout)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `WithPreRunHook(hook func() error)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `WithPostRunHook(hook func() error)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `WithPanicRecovery(enabled bool)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |
| `AddValidator(v CommandValidator)` | ❌ NOT_DONE | ARCHITECTURE_REVIEW suggestion |

---

## Test Coverage Gaps (0% Coverage)

| Package | Has Tests | Action |
|---------|-----------|--------|
| `internal/commands` | ❌ NO | Add `*_test.go` files |
| `internal/di` | ❌ NO | Add `*_test.go` files |
| `pkg/cmdguard` | ❌ NO | Add `*_test.go` files |
| `internal/config` | ✅ YES | ~47.6% coverage |
| `internal/validation` | ✅ YES | Tests exist |

---

## CLI UX Compliance (Per CLI_DESIGN_PRINCIPLES.md)

### Compliant ✅
| Principle | Evidence |
|-----------|----------|
| Boolean flags use `BoolP()` | `root.go:46` - `BoolP("strict", "s", false, ...)` |
| Short flags for common options | All 3 flags have short versions: `-c`, `-l`, `-s` |
| Enum validation in PreRunE | `root.go:49-58` validates log-level |
| Consistent kebab-case naming | All flags: `config`, `log-level`, `strict` |

### Violations ❌
| Principle | File | Issue |
|-----------|------|-------|
| Copy-pasteable examples | `root.go:33-41` | No `Example:` field on commands |
| Default values in help | `root.go:45` | `--log-level` doesn't show "(default: info)" |
| Unknown flag suggestions | N/A | No "Did you mean --xyz?" feature |

---

## Summary Statistics

| Category | Total | Done | Partial | Not Done |
|----------|-------|------|---------|----------|
| Critical Issues | 8 | 0 | 0 | 8 |
| High Priority | 2 | 0 | 0 | 2 |
| Medium Priority | 12 | 0 | 1 | 11 |
| Lower Priority | 24 | 0 | 0 | 24 |
| API Options | 7 | 0 | 0 | 7 |
| Missing Docs | 7 | 0 | 0 | 7 |
| CLI UX Violations | 3 | 0 | 0 | 3 |
| Decisions | 3 | 3 | 0 | 0 |
| Superseded (Option A) | 4 | - | - | - |
| **TOTAL ACTIONABLE** | **66** | **3** | **1** | **62** |

---

## Changelog

### 2026-02-14
- Initial creation of TODO_LIST.md
- Analyzed AGENTS.md
- Verified issues against actual code
- Disproved 2 claims (unused import, missing dependency)
- Analyzed README.md
- Found 3 new critical issues (false docs, broken examples, no-op hook)
- Found 4 undocumented API methods
- Analyzed ARCHITECTURE_REVIEW.md
- Found 3 new critical issues (context handling, silent errors, error aggregation)
- Found test coverage gaps (0% in commands/di/pkg)
- Added 7 suggested API options
- Added 5 performance/code quality improvements
- Analyzed CLI_DESIGN_PRINCIPLES.md
- Found 3 CLI UX violations (missing examples, missing defaults, no flag suggestions)
- 4 principles compliant (BoolP, short flags, enum validation, kebab-case)
- Analyzed DECISION_POINT.md
- Documented 3 key decisions (keep framework, README update, docs alignment)
- All action items already tracked - no new TODOs
- Analyzed POST_DECISION_STATE.md
- Added 1 critical issue (architecture diagram outdated)
- Added 5 medium priority items (tests, validation enhancements)
- Added 6 lower priority items (CI/CD matrix, plugin system)
- Added 3 documentation items (API examples, CONTRIBUTING.md, version tagging)
- Analyzed COMPREHENSIVE_STATUS.md
- No new actionable TODOs (all items already tracked)
- Disproved 1 claim (AGENTS.md exists - was claimed missing)
- Analyzed BUILD_VERIFICATION_COMPLETE.md
- Added 2 medium priority items (error context, context variables)
- Added 1 lower priority item (phantom types - optional)
- Disproved AGENTS.md missing claim (file now exists)
- Analyzed CMDGUARD_TRANSFORMATION_PLAN.md
- No new actionable TODOs (all items already tracked in Transformation section)
- 4 tasks marked obsolete due to Option A (framework) decision
- Analyzed INITIAL_IMPLEMENTATION_COMPLETE.md
- No new actionable TODOs (all P0-P3 items already tracked)
- Disproved 1 claim ("No Structured Logging" - logger.go exists with full slog)
- Status doc has historical inaccuracies (line counts, versions) - not actionable
- **DE-DUPLICATION PASS COMPLETE**
- Removed 2 duplicate Critical Issues (#2 Code Duplication, #3 Manual DI Wiring - merged into other sections)
- Renumbered Critical Issues from 10 to 8
- Added "Superseded by Option A Decision" section (4 items)
- Moved 3 High Priority items to Superseded (Remove cmd/, Redesign API, Compile-time validation)
- Moved 1 Medium Priority item to Superseded (Define guard architecture)
- Updated statistics: 74 → 66 actionable items
