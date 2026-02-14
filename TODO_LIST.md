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
| docs/status/2026-02-14_04-19_CLI_IMPROVEMENTS.md | ⏳ Pending | - |
| docs/status/2026-02-14_05-44_POST_DECISION_STATE.md | ⏳ Pending | - |
| docs/status/2026-02-14_04-25_COMPREHENSIVE_STATUS.md | ⏳ Pending | - |
| docs/status/2026-02-14_04-11_BUILD_VERIFICATION_COMPLETE.md | ⏳ Pending | - |
| docs/planning/2026-02-14_04-21_CMDGUARD_TRANSFORMATION_PLAN.md | ⏳ Pending | - |
| docs/status/2026-02-14_03-48_INITIAL_IMPLEMENTATION_COMPLETE.md | ⏳ Pending | - |

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

### 2. Code Duplication [NOT_DONE]
**Source:** AGENTS.md (verified in code)
**Details:** Version command logic duplicated in `main.go:86-92` and `root.go:148-158`
**Action:** Extract to single source of truth

### 3. Manual DI Wiring [NOT_DONE]
**Source:** AGENTS.md → Transformation Plan
**File:** `pkg/cmdguard/public_api.go:90`
**Details:** `SetValidator()` manual wiring still exists
**Action:** Use proper constructor injection

### 4. False Documentation Claims [NOT_DONE]
**Source:** README.md (lines 164-165)
**Claims:**
- "No duplicate command names" validation - NOT IMPLEMENTED
- "No conflicting aliases" validation - NOT IMPLEMENTED
**Action:** Either implement these validations or remove from documentation

### 5. Broken Code Examples in Documentation [NOT_DONE]
**Source:** README.md (lines 45-77, 81-95)
**Details:** Code examples missing `"fmt"` import, will not compile
**Action:** Fix examples to be copy-pasteable

### 6. WithValidationHook is a No-Op [NOT_DONE]
**Source:** README.md (lines 154-155) → `pkg/cmdguard/public_api.go:130-135`
**Details:** Hook is documented but implementation discards it (returns `nil`)
**Action:** Either implement or remove from API

### 7. Incomplete Context Handling [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P1)
**File:** `internal/di/module.go:123-127`
**Details:** Context passed to `Shutdown()` is completely ignored
**Action:** Implement proper context-aware shutdown

### 8. Silent Config Errors [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P1)
**File:** `internal/config/provider.go:32-35`
**Details:** Errors discarded with `_` assignment
**Action:** Properly handle/return config loading errors

### 9. No Error Aggregation [NOT_DONE]
**Source:** ARCHITECTURE_REVIEW.md (P2)
**File:** `internal/di/module.go:115-117`
**Details:** Returns only `errs[0]` instead of using `errors.Join()`
**Action:** Aggregate multiple errors with `errors.Join()`

---

## Disproved Issues (Remove from Documentation)

| Claim | Reality | Action |
|-------|---------|--------|
| Unused `log/slog` import in root.go:7 | `slog` IS used at lines 23, 141, 154 | Remove from AGENTS.md |
| Missing `charmbracelet/log` dependency | `logger.go` only imports stdlib (`log/slog`, `os`) | Remove from AGENTS.md |

---

## Missing Documentation (Exists but not documented)

| Feature | Location | Action |
|---------|----------|--------|
| `app.IsStrictMode()` | `pkg/cmdguard/public_api.go:191` | Add to README API Reference |
| `app.AddCommand()` | `pkg/cmdguard/public_api.go:212` | Add to README API Reference |
| Built-in `validate` command | Auto-added by `SetupCommands()` | Document in README |
| Built-in `version` command | Auto-added by `SetupCommands()` | Document in README |

---

## Planned Changes (From Transformation Plan)

### High Priority - Core Architecture

| Task | Status | Notes |
|------|--------|-------|
| Remove `cmd/` folder | ❌ NOT_DONE | Folder still exists at `cmd/cmdguard/` |
| Redesign public API | ❌ NOT_DONE | Still uses multi-step Initialize/Validate pattern |
| Add compile-time validation | ❌ NOT_DONE | No interception of Cobra calls |
| Fix DI usage | ❌ NOT_DONE | Manual wiring via `SetValidator()` exists |
| Add justfile | ❌ NOT_DONE | No justfile found |

### Medium Priority - Code Quality

| Task | Status | Notes |
|------|--------|-------|
| Improve test coverage (80%+) | 🔄 PARTIAL | config at ~47.6% per plan |
| Fix 4 errcheck violations | ❌ NOT_DONE | See Critical Issues #1 |
| Define guard architecture in ARCHITECTURE.md | ❌ NOT_DONE | Task 1.2 |
| Fix code duplication (7 clone groups) | ❌ NOT_DONE | Task 3.3 |
| Add integration tests | ❌ NOT_DONE | Task 3.4 |
| Create examples directory | ❌ NOT_DONE | Task 3.5 |

### Lower Priority - Future Work

| Task | Status | Notes |
|------|--------|-------|
| Add CI/CD workflow | ❌ NOT_DONE | Task 4.1 |
| Benchmarks | ❌ NOT_DONE | Beyond 80% |
| Fuzz testing | ❌ NOT_DONE | Beyond 80% |
| Release automation | ❌ NOT_DONE | Beyond 80% |
| Container image | ❌ NOT_DONE | Beyond 80% |
| Homebrew formula | ❌ NOT_DONE | Beyond 80% |
| Plugin system | ❌ NOT_DONE | Beyond 80% |
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
| Critical Issues | 9 | 0 | 0 | 9 |
| High Priority | 5 | 0 | 0 | 5 |
| Medium Priority | 6 | 0 | 1 | 5 |
| Lower Priority | 19 | 0 | 0 | 19 |
| API Options | 7 | 0 | 0 | 7 |
| Missing Docs | 4 | 0 | 0 | 4 |
| CLI UX Violations | 3 | 0 | 0 | 3 |
| Decisions | 3 | 3 | 0 | 0 |
| **TOTAL** | **56** | **3** | **1** | **52** |

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
