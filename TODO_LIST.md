# TODO_LIST.md - cmdguard Project Tasks

**Last Updated:** 2026-02-15
**Purpose:** Track tasks for the cmdguard v2 API implementation.

---

## Completed (v2 Implementation)

| Task | Status | Notes |
|------|--------|-------|
| Design v2 API with generics | ✅ DONE | GuardedCommand[T], Command[T] |
| Implement errors.go | ✅ DONE | Typed errors, no panics |
| Implement types.go | ✅ DONE | Common types and interfaces |
| Implement config.go | ✅ DONE | Typed configuration |
| Implement flags.go | ✅ DONE | FlagRegistry with struct tags |
| Implement scope.go | ✅ DONE | DI scope with samber/do/v2 |
| Implement command.go | ✅ DONE | Command[T] definition |
| Implement guard.go | ✅ DONE | GuardedCommand[T] implementation |

---

## Completed (v2 Testing)

| Task | Status | Notes |
|------|--------|-------|
| Test errors.go | ✅ DONE | 142 lines of tests |
| Test types.go | ✅ DONE | 346 lines of tests |
| Test config.go | ✅ DONE | 360 lines of tests |
| Test flags.go | ✅ DONE | 488 lines of tests |
| Test scope.go | ✅ DONE | 458 lines of tests |
| Test command.go | ✅ DONE | 399 lines of tests |
| Test guard.go | ✅ DONE | 565 lines of tests |

---

## Completed (v2 Polish)

| Task | Status | Notes |
|------|--------|-------|
| Add typed example | ✅ DONE | examples/typed/main.go |
| Update FEATURES.md | ✅ DONE | v2 API documentation |
| Update TODO_LIST.md | ✅ DONE | This file |

---

## Remaining Tasks

### Low Priority

| Task | Status | Notes |
|------|--------|-------|
| Update README.md for v2 | ⏳ PENDING | Add v2 examples |
| Update AGENTS.md for v2 | ⏳ PENDING | Document v2 patterns |
| Add more examples | ⏳ PENDING | DI patterns, advanced flags |
| Plugin system for custom validators | ⏳ PENDING | Future enhancement |
| Enhanced flag validation | ⏳ PENDING | Enums, custom validators |
| Performance benchmarks | ⏳ PENDING | Not yet needed |
| Release automation | ⏳ PENDING | Manual releases sufficient |

---

## Current Project Structure

```
cmdguard/
├── .github/workflows/ci.yml    # CI/CD pipeline
├── pkg/cmdguard/
│   ├── v2/                     # v2 API (recommended)
│   │   ├── errors.go           # Typed errors
│   │   ├── types.go            # Common types
│   │   ├── config.go           # Configuration
│   │   ├── flags.go            # Flag registry
│   │   ├── scope.go            # DI scope
│   │   ├── command.go          # Command[T]
│   │   └── guard.go            # GuardedCommand[T]
│   ├── guarded_command.go      # v1 API
│   └── guarded_command_test.go
├── internal/
│   ├── config/                 # Configuration (95.7% coverage)
│   └── logging/                # Logging (100% coverage)
├── examples/
│   ├── basic/main.go           # Simple CLI example
│   ├── advanced/main.go        # Nested commands example
│   ├── guarded/main.go         # v1 panic demo
│   └── typed/main.go           # v2 API demo
├── AGENTS.md                   # Developer guide
├── CONTRIBUTING.md             # Contribution guide
├── FEATURES.md                 # Feature documentation
├── README.md                   # User documentation
├── TODO_LIST.md                # This file
├── go.mod                      # Dependencies
└── justfile                    # Build automation
```

---

## Summary

**Status: v2.0.0 READY ✅**

The v2 API is complete with comprehensive test coverage:
- 7 implementation files (~1,700 lines)
- 7 test files (~2,700 lines)
- Complete example demonstrating all features

**Key v2 Features:**
- Type-safe with generics (`GuardedCommand[T]`)
- No panics - all operations return errors
- DI integration with samber/do/v2
- Typed flags with struct tags

**Remaining work is documentation polish only.**
