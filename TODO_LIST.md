# TODO_LIST.md - cmdguard Project Tasks

**Last Updated:** 2026-02-14
**Purpose:** Track remaining tasks for the cmdguard Guard API implementation.

---

## Completed (Phase 1)

| Task | Completed |
|------|-----------|
| Remove `cmd/` folder | DONE |
| Redesign public API for single-step initialization | DONE |
| Add compile-time validation (panic on invalid) | DONE |
| Fix errcheck violations | DONE |
| Remove orphaned internal packages | DONE |
| Rewrite README for Guard API | DONE |
| Rewrite FEATURES.md for Guard API | DONE |

---

## Completed (Phase 2 - Testing)

| Task | Status | Notes |
|------|--------|-------|
| Add tests for `pkg/cmdguard` | DONE | GuardedCommand has 66.7% coverage |
| Add tests for `internal/logging` | DONE | logging package has 100% coverage |
| Update AGENTS.md for current architecture | DONE | Updated 2026-02-14 |

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

## Completed (Phase 3 - Build & Examples)

| Task | Status | Notes |
|------|--------|-------|
| Add version injection at build time | DONE | Use ldflags or `just build-version X.Y.Z` |
| Add justfile for common tasks | DONE | build, test, lint, verify, etc. |
| Create examples directory | DONE | basic, advanced, guarded examples |
| Clean up docs/planning/ folder | DONE | Archived (deleted) |
| Clean up docs/status/ folder | DONE | Archived (deleted) |
| Improve test coverage to 80%+ | DONE | config 94.1%, cmdguard 66.7%, logging 100% |

---

## Medium Priority

| Task | Status | Notes |
|------|--------|-------|
| Add CONTRIBUTING.md | NOT_DONE | No contribution guide |
| Add tests for ExecuteAndExit() | NOT_DONE | Requires os.Exit testing |

---

## Lower Priority

| Task | Status | Notes |
|------|--------|-------|
| Add CI/CD workflow | NOT_DONE | No GitHub Actions |
| Add JSON logging option | NOT_DONE | Only text handler |
| Plugin system for custom validators | NOT_DONE | Future enhancement |
| Enhanced flag validation | NOT_DONE | Type validation, required flags |
| Performance benchmarks | NOT_DONE | Not yet needed |
| Release automation | NOT_DONE | Manual releases |

---

## Current Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   └── guarded_command.go    # Public API (290 lines)
├── internal/
│   ├── config/               # Config loading (provider.go)
│   └── logging/              # slog integration (logger.go)
├── examples/
│   ├── basic/main.go         # Simple CLI example
│   ├── advanced/main.go      # Nested commands example
│   └── guarded/main.go       # Validation demo
├── tests/integration/        # Integration tests
├── docs/
│   └── CLI_DESIGN_PRINCIPLES.md
├── README.md                 # Updated for Guard API
├── FEATURES.md               # Updated for Guard API
├── TODO_LIST.md              # This file
└── justfile                  # Build commands
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
| Phase 3 (Build & Examples) | 6 | 6 | 0 |
| Medium Priority | 2 | 0 | 2 |
| Lower Priority | 6 | 0 | 6 |
| **TOTAL** | **24** | **16** | **8** |

---

## Changelog

### 2026-02-14 (Phase 3 Complete)
- Added version injection at build time (ldflags support)
- Added `just build-version <version>` recipe
- Added `Version()` function to get version programmatically
- Verified justfile exists with all needed commands
- Verified examples directory has 3 working examples
- Archived docs/planning/ and docs/status/ folders
- Pushed all changes to remote

### 2026-02-14 (Phase 2 Complete)
- Added tests for logging package (100% coverage)
- Updated AGENTS.md for current architecture
- Marked logging tests as complete
- Updated coverage stats

### 2026-02-14 (Phase 1 Complete)
- Added tests for GuardedCommand (20 test cases, 66.7% coverage)
- Fixed .gitignore to use /cmdguard (root binary only)
- Marked most test cases as complete
- Fixed table markdown syntax

### 2026-02-14 (Initial)
- Rewrote TODO_LIST.md for Guard API
- Marked Phase 1 items as complete
- Updated project structure
- Removed obsolete items (framework-specific tasks)
