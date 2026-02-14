# TODO_LIST.md - cmdguard Project Tasks

**Last Updated:** 2026-02-14
**Purpose:** Track remaining tasks for the cmdguard Guard API implementation.

---

## Completed (Phase 1)

|| Task | Completed |
||------|-----------|
|| Remove `cmd/` folder | ✅ Done |
|| Redesign public API for single-step initialization | ✅ Done |
|| Add compile-time validation (panic on invalid) | ✅ Done |
|| Fix errcheck violations | ✅ Done |
|| Remove orphaned internal packages | ✅ Done |
|| Rewrite README for Guard API | ✅ Done |
|| Rewrite FEATURES.md for Guard API | ✅ Done |

---

## Current Focus (Phase 2 - Testing)

### High Priority

|| Task | Status | Notes |
||------|--------|-------|
|| Add tests for `pkg/cmdguard` | ❌ NOT_DONE | GuardedCommand has 0% coverage |
|| Add tests for `internal/logging` | ❌ NOT_DONE | logging package has 0% coverage |

### Test Cases Needed for GuardedCommand

- [ ] `New()` creates GuardedCommand with defaults
- [ ] `New()` loads config from environment
- [ ] `AddCommand()` panics on invalid command (no handler)
- [ ] `AddCommand()` panics on command without name
- [ ] `AddCommand()` accepts valid command
- [ ] `AddCommand()` panics after Execute() called
- [ ] `AddSubcommand()` panics on invalid child
- [ ] `AddSubcommand()` adds to parent correctly
- [ ] `Execute()` runs command successfully
- [ ] `ExecuteAndExit()` exits with code 0 on success
- [ ] `ExecuteAndExit()` exits with code 1 on error
- [ ] Strict mode requires RunE (not Run)
- [ ] Parent commands don't need handlers
- [ ] `validate` command works
- [ ] `version` command works
- [ ] Log-level validation in PreRunE
- [ ] `Command()` returns underlying cobra.Command
- [ ] `Config()` returns config
- [ ] `IsStrictMode()` returns correct value

---

## Medium Priority

|| Task | Status | Notes |
||------|--------|-------|
|| Add version injection at build time | ❌ NOT_DONE | Version hardcoded to "0.1.0" |
|| Add justfile for common tasks | ❌ NOT_DONE | No justfile exists |
|| Create examples directory | ❌ NOT_DONE | No example applications |
|| Add CONTRIBUTING.md | ❌ NOT_DONE | No contribution guide |
|| Improve test coverage to 80%+ | 🔄 PARTIAL | config at 94.1%, others at 0% |

---

## Lower Priority

|| Task | Status | Notes |
||------|--------|-------|
|| Add CI/CD workflow | ❌ NOT_DONE | No GitHub Actions |
|| Add JSON logging option | ❌ NOT_DONE | Only text handler |
|| Plugin system for custom validators | ❌ NOT_DONE | Future enhancement |
|| Enhanced flag validation | ❌ NOT_DONE | Type validation, required flags |
|| Performance benchmarks | ❌ NOT_DONE | Not yet needed |
|| Release automation | ❌ NOT_DONE | Manual releases |

---

## Documentation Updates Needed

|| Task | Status | Notes |
||------|--------|-------|
|| Update AGENTS.md for current architecture | ❌ NOT_DONE | Still references old structure |
|| Clean up docs/planning/ folder | ❌ NOT_DONE | Transformation plan completed |
|| Clean up docs/status/ folder | ❌ NOT_DONE | Historical status reports |

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
| Phase 2 (Testing) | 2 | 0 | 2 |
| Medium Priority | 5 | 0 | 5 |
| Lower Priority | 6 | 0 | 6 |
| Documentation | 3 | 0 | 3 |
| **TOTAL** | **23** | **7** | **16** |

---

## Changelog

### 2026-02-14
- Rewrote TODO_LIST.md for Guard API
- Marked Phase 1 items as complete
- Updated project structure
- Removed obsolete items (framework-specific tasks)
- Added specific test cases needed for GuardedCommand
- Simplified to focus on actionable items

### Earlier 2026-02-14 (Historical)
- Initial creation tracking framework tasks
- Analyzed multiple docs files
- Tracked transformation plan progress
- All framework tasks now complete or obsolete
