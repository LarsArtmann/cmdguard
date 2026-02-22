# Status Report: Must* Function Removal

**Date:** 2026-02-22 09:42
**Status:** In Progress
**Author:** AI Assistant

---

## Executive Summary

Removing all `Must*` functions from v2 package to achieve true "no panics" guarantee. The v2 package is documented as panic-free, but currently contains 5 `Must*` functions that can panic.

---

## Current State

### Must* Functions to Remove

| Function | File | Line | Alternative |
|----------|------|------|-------------|
| `MustNewCommand` | command.go | 222 | `NewCommand` |
| `MustAddCommand` | guard_command.go | 39 | `AddCommand` |
| `MustAddAnyCommand` | guard_command.go | 81 | `AddAnyCommand` |
| `MustEnum` | types.go | 27 | `ParseEnum` |
| `MustDuration` | types.go | 72 | `ParseDuration` |

### Files Requiring Updates

**Source Files (removal):**
- `pkg/cmdguard/v2/command.go`
- `pkg/cmdguard/v2/guard_command.go`
- `pkg/cmdguard/v2/types.go`

**Test Files (usage updates):**
- `pkg/cmdguard/v2/command_test.go` - 4 usages
- `pkg/cmdguard/v2/guard_test.go` - 7 usages
- `pkg/cmdguard/v2/types_test.go` - 8 usages
- `pkg/cmdguard/v2/example_test.go` - 2 usages

**Documentation:**
- `FEATURES.md`
- `docs/FEATURES.md`

---

## Test Status

All tests pass before starting:
```
ok  github.com/larsartmann/cmdguard/pkg/cmdguard/v2  0.722s
ok  github.com/larsartmann/cmdguard/pkg/cmdguard     0.721s
ok  github.com/larsartmann/cmdguard/internal/config  0.228s
ok  github.com/larsartmann/cmdguard/internal/logging 0.417s
```

---

## Rationale

The v2 package philosophy is "fail with errors, not panics." Keeping `Must*` functions:
1. Violates the documented behavior
2. Creates inconsistency in the API
3. Encourages panic-prone patterns

By removing them, we force users to handle errors explicitly, which is the Go way.

---

## Execution Checklist

- [ ] Remove `MustNewCommand` from command.go
- [ ] Update command_test.go to use `NewCommand`
- [ ] Remove `MustAddCommand` and `MustAddAnyCommand` from guard_command.go
- [ ] Update guard_test.go to use `AddCommand` and `AddAnyCommand`
- [ ] Remove `MustEnum` and `MustDuration` from types.go
- [ ] Update types_test.go to use `ParseEnum` and `ParseDuration`
- [ ] Update example_test.go to use non-panicking alternatives
- [ ] Update FEATURES.md documentation
- [ ] Update docs/FEATURES.md documentation
- [ ] Run full test suite
- [ ] Commit changes

---

## Risk Assessment

**Low Risk:** All `Must*` functions have direct non-panicking equivalents. The refactoring is mechanical.

**Test Coverage:** High (88.9% for v2) - changes will be validated.

---

## Notes

- Previous session removed `MustInvoke` from DI integration
- This completes the "no panics" initiative for v2
- Uncommitted file `internal/config/config_bdd_test.go` is unrelated
