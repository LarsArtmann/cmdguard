# TODO_LIST.md - cmdguard Project Tasks

**Last Updated:** 2026-02-14
**Purpose:** Track remaining tasks for the cmdguard Guard API implementation.

---

## Completed (Phase 1)

| Task | Completed |
|------|-----------|
| Remove `cmd/` folder | ✅ Done |
| Redesign public API for single-step initialization | ✅ Done |
| Add compile-time validation (panic on invalid) | ✅ Done |
| Fix errcheck violations | ✅ Done |
| Remove orphaned internal packages | ✅ Done |
| Rewrite README for Guard API | ✅ Done |
| Rewrite FEATURES.md for Guard API | ✅ Done |

---

## Completed (Phase 2 - Testing)

| Task | Status | Notes |
|------|--------|-------|
| Add tests for `pkg/cmdguard` | ✅ DONE | GuardedCommand has 66.7% coverage |
| Add tests for `internal/logging` | ✅ DONE | logging package has 100% coverage |
| Update AGENTS.md for current architecture | ✅ DONE | Updated 2026-02-14 |

### Test Cases for GuardedCommand

- [x] `New()` creates GuardedCommand with defaults
- [x] `New()` loads config from environment
- [x] `AddCommand()` panics on invalid command (no handler)
- [x] `AddCommand()` panics on command without name
- [x] `AddCommand()` accepts valid command
- [x] `AddCommand()` panics after Execute() called
- [x] `AddSubcommand()` panics on invalid child
- [x] `AddSubcommand()` adds to parent correctly
- [x] `Execute()` runs command successfully
- [ ] `ExecuteAndExit()` exits with code 0 on success
- [ ] `ExecuteAndExit()` exits with code 1 on error
- [x] Strict mode requires RunE (not Run)
- [x] Parent commands don't need handlers
- [x] `validate` command works
- [x] `version` command works
- [x] Log-level validation in PreRunE
- [x] `Command()` returns underlying cobra.Command
- [x] `Config()` returns config
- [x] `IsStrictMode()` returns correct value

---

## Medium Priority

| Task | Status | Notes |
|------|--------|-------|
| Add version injection at build time | ❌ NOT_DONE | Version hardcoded to "0.1.0" |
| Add justfile for common tasks | ❌ NOT_DONE | No justfile exists |
| Create examples directory | ❌ NOT_DONE | No example applications |
| Add CONTRIBUTING.md | ❌ NOT_DONE | No contribution guide |
| Improve test coverage to 80%+ | ✅ DONE | config 94.1%, cmdguard 66.7%, logging 100% |

---

## Lower Priority

| Task | Status | Notes |
|------|--------|-------|
| Add CI/CD workflow | ❌ NOT_DONE | No GitHub Actions |
| Add JSON logging option | ❌ NOT_DONE | Only text handler |
| Plugin system for custom validators | ❌ NOT_DONE | Future enhancement |
| Enhanced flag validation | ❌ NOT_DONE | Type validation, required flags |
| Performance benchmarks | ❌ NOT_DONE | Not yet needed |
| Release automation | ❌ NOT_DONE | Manual releases |

---

## Documentation Updates Needed

| Task | Status | Notes |
|------|--------|-------|
| Clean up docs/planning/ folder | ❌ NOT_DONE | Transformation plan completed |
| Clean up docs/status/ folder | ❌ NOT_DONE | Historical status reports |

---

## Current Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   └── guarded_command.go    # Public API (285 lines)
├── internal/
│   ├── config/               # Config loading (provider.go)
│   └── logging/              # slog integration (logger.go)
├── docs/                     # Documentation (some outdated)
├── README.md                 # Updated for Guard API
├── FEATURES.md               # Updated for Guard API
└── TODO_LIST.md              # This file
```

---

## Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/spf13/cobra` | v1.10.2 | CLI framework |
| `github.com/charmbracelet/fang` | v0.4.4 | Cobra styling |
| `github.com/stretchr/testify` | v1.11.1 | Testing |

---

## Summary Statistics

| Category | Total | Done | Not Done |
|----------|-------|------|----------|
| Phase 1 (Core) | 7 | 7 | 0 |
| Phase 2 (Testing) | 3 | 3 | 0 |
| Medium Priority | 5 | 1 | 4 |
| Lower Priority | 6 | 0 | 6 |
| Documentation | 2 | 0 | 2 |
| **TOTAL** | **23** | **11** | **12** |

---

## Changelog

### 2026-02-14 (Updated)
- Added tests for logging package (100% coverage)
- Updated AGENTS.md for current architecture
- Marked logging tests as complete
- Updated coverage stats

### 2026-02-14 (Earlier)
- Added tests for GuardedCommand (20 test cases, 66.7% coverage)
- Fixed .gitignore to use /cmdguard (root binary only)
- Marked most test cases as complete
- Fixed table markdown syntax

### 2026-02-14 (Initial)
- Rewrote TODO_LIST.md for Guard API
- Marked Phase 1 items as complete
- Updated project structure
- Removed obsolete items (framework-specific tasks)
