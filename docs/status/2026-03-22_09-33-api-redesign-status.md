# cmdguard v2.1.0 API Redesign - Status Report

**Date:** 2026-03-22 09:33
**Status:** IN PROGRESS - Phase 1 Started, Build passing

---

## Executive Summary

The project is redesigning the cmdguard v2 API to v2.1.0. The main goal is to simplify the type system, make DI optional, and improve the developer experience.

**Current State:** Tests pass ✅ | Build passes ✅ | Implementation NOT started

---

## A) FULLY DONE ✅

1. **Planning Documents Created**
   - `docs/API_DESIGN_REVIEW.md` - Comprehensive API design research (52KB)
   - `docs/planning/2026-03-22_09-14-api-redesign-v2.1.md` - Implementation plan with 96 subtasks

2. **Research Completed**
   - Analyzed 5 projects using cmdguard: timesheets, projects-management-automation, docs-organizer, code-duplicate-analyzer
   - Identified key issues: P0 redundant type param, P0 confusing name, P1 forced DI, P2 `any` usage

3. **Pareto Analysis**
   - 1% → 51%: Type-safe flag definitions
   - 4% → 64%: Remove F param, make DI optional
   - 20% → 80%: Rename, functional options, consolidation

4. **Git Commits Pushed**
   - `3b4d5ea` - docs: add nil-safety design and deprecation strategy
   - `ffa818c` - docs: add v2.1 API redesign planning document
   - `0ef4e32` - docs: add comprehensive API design review
   - Earlier commits: test framework migration from ginkgo to native Go testing

5. **Test Coverage**: 89.0% (target: 90%+)

---

## b) PARTIALLY DONE ⚠️

1. **Phase 1: CLI[T] Type Creation**
   - Status: File created but deleted/lost during session transition
   - File: `pkg/cmdguard/v2/cli.go`
   - Issue: Need to recreate and integrate properly

2. **Documentation Updates**
   - `docs/API_DESIGN_REVIEW.md` - Modified but not committed
   - Added nil-safety patterns and deprecation strategy

---

## c) NOT STARTED ❌

1. **Phase 2: Functional Options** - WithDI(), WithVersion(), WithLong()
2. **Phase 3: AddCommand accepting any flags type**
3. **Phase 4: Remove AddAnyCommand**
4. **Phase 5: Fix NewFlagRegistry generic**
5. **Phase 6: Consolidate Scope() methods**
6. **Phase 7: Remove AddCommandFunc**
7. **Phase 8: CLI type alias for backward compat**
8. **Phase 9: Update tests**
9. **Phase 10: Update examples**
10. **Phase 11: Update documentation**
11. **Phase 12: Final verification**

---

## d) TOTALLY FUCKED UP 💥

**NONE** - The codebase is in good shape!

- Tests pass: ✅
- Build passes: ✅
- No breaking changes committed: ✅
- All planning done: ✅

---

## e) WHAT WE SHOULD IMPROVE 📈

1. **Better session continuity** - Files created in one session were next
2. **Smaller commits** - Commit planning docs separately from code changes
3. **Test-first approach** - Write tests before implementation
4. **Incremental phases** - Each phase should be committed separately

---

## f) TOP #25 THINGS TO GET DONE NEXT

| Priority | Task                                                     | Impact   | Effort |
| -------- | -------------------------------------------------------- | -------- | ------ |
| 1        | Create `cli.go` with `CLI[T]` type                       | Critical | 15min  |
| 2        | Add `AddCommand` to `CLI[T]` accepting any flags         | Critical | 20min  |
| 3        | Create `cli_command.go` for command handling             | High     | 20min  |
| 4        | Add functional options: WithDI, WithVersion, WithLong    | High     | 15min  |
| 5        | Write tests for `CLI[T]`                                 | High     | 30min  |
| 6        | Update `guard_exec.go` for CLI[T]                        | High     | 15min  |
| 7        | Create backward compat type alias `GuardedCommand = CLI` | Medium   | 5min   |
| 8        | Fix `NewFlagRegistry` to be generic                      | Medium   | 20min  |
| 9        | Remove `AddAnyCommand` (replaced by AddCommand)          | Medium   | 10min  |
| 10       | Consolidate Scope() methods                              | Medium   | 15min  |
| 11       | Remove `AddCommandFunc`                                  | Low      | 5min   |
| 12       | Update examples/typed to use CLI[T]                      | High     | 20min  |
| 13       | Update examples/di to use CLI[T]                         | Medium   | 15min  |
| 14       | Update examples/advanced-flags                           | Medium   | 15min  |
| 15       | Create MIGRATION.md guide                                | High     | 30min  |
| 16       | Update README.md for v2.1                                | High     | 20min  |
| 17       | Update AGENTS.md API section                             | Medium   | 15min  |
| 18       | Add `Package()` for samber/do integration                | Medium   | 20min  |
| 19       | Add flag documentation generator                         | Medium   | 30min  |
| 20       | Write integration tests for CLI[T]                       | High     | 30min  |
| 21       | Verify 90%+ test coverage                                | Critical | 15min  |
| 22       | Run linter and fix issues                                | Medium   | 30min  |
| 23       | Update Version constant to 2.1.0                         | Low      | 5min   |
| 24       | Create examples/docs-generator                           | Low      | 30min  |
| 25       | Final manual testing of all examples                     | Medium   | 20min  |

---

## g) TOP #1 QUESTION 🤔

**Should we keep `GuardedCommand[T,F]` as a type alias to `CLI[T]` for backward compatibility, OR should we create a separate v3 package?**

Options:

1. **Type alias** - `type GuardedCommand[T, F] = CLI[T]` - but F would be ignored, potentially confusing
2. **Separate v3** - Create `pkg/cmdguard/v3` with clean API
3. **Keep both** - Maintain `GuardedCommand[T,F]` alongside `CLI[T]` with deprecation notice

**My recommendation:** Keep both with deprecation notice. Type aliases with different type param counts don't work well in Go.

---

## Current Git Status

```
Modified (not staged):
  docs/API_DESIGN_REVIEW.md
  pkg/cmdguard/v2/flags_parse_test.go
  pkg/cmdguard/v2/flags_registry_test.go
```

---

## Next Immediate Actions

1. Commit current documentation changes
2. Create `cli.go` properly (was lost in session transition)
3. Implement `AddCommand` for `CLI[T]`
4. Write tests for new API
5. Update examples
