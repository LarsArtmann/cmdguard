# Deduplication Sprint Status Report

**Date:** 2026-04-16 04:49  
**Branch:** master  
**Commit:** 84cea37  
**Task:** Deduplicate code clones identified by `art-dupl --semantic -t 15 --sort total-tokens`

---

## Executive Summary

Previous session partially completed test deduplication work but **left the codebase in a broken state** — `go vet` fails with `undefined: noOpRunEForTestCLIConfigWithFlags` at `cli_lifecycle_test.go:197`. The function was deleted but one call site was missed. Clone count went from 86 → 86 (no improvement in group count yet, though inline lambda count reduced).

---

## a) FULLY DONE

1. **Added `noOpRunE[T, F]` generic to `test_helpers_test.go`** (package `v2`)  
   - Generic no-op RunE function that works with any config/flags type combination
   - File: `pkg/cmdguard/v2/test_helpers_test.go:95-97`

2. **Refactored `noOpHandler()` and `noOpHandlerForTestAppConfig()`** to delegate to `noOpRunE`  
   - Eliminates duplicate inline lambda definitions
   - File: `pkg/cmdguard/v2/test_helpers_test.go:99-107`

3. **Replaced 7 inline no-op lambdas in `cli_groups_test.go`** with `noOpRunE[testConfig, NoFlags]`  
   - Lines: 31, 40, 49, 102, 202, 209, 216

4. **Replaced 3 inline no-op lambdas in `middleware_test.go`** with `noOpRunE[testConfig, NoFlags]`  
   - Lines: 166-168, 205-207, 309

5. **Replaced 1 reference in `cli_lifecycle_test.go`** at line 165 with `NoOpRunEWithFlags[testCLIConfig, testFlags]()`  
   - File: `pkg/cmdguard/v2/cli_lifecycle_test.go:159`

6. **Unrelated: Added explicit type param `buildChain[T]` in `cli_command.go`**  
   - File: `pkg/cmdguard/v2/cli_command.go:168`

---

## b) PARTIALLY DONE

1. **`cli_lifecycle_test.go` deduplication** — INCOMPLETE AND BROKEN  
   - The local `noOpRunEForTestCLIConfigWithFlags` function was deleted (lines 10-14)
   - Line 159 was updated to use `NoOpRunEWithFlags[testCLIConfig, testFlags]()`
   - **BUT line 197 still references the deleted function** → compile error
   - Status: **BROKEN, MUST FIX IMMEDIATELY**

---

## c) NOT STARTED

### Scope Test Deduplication (Clone Group: 9 instances)
- `scope_lifecycle_test.go` — 4x `if err := ProvideValue(scope, "value"); err != nil { t.Fatalf(...) }`
- `scope_child_test.go:209-211` — same pattern in `assertChildInheritsParent`
- `scope_integration_test.go:78-80` — same pattern
- `scope_provide_basic_test.go:127,129 + 171,173 + 175,177` — same pattern
- **Solution:** Add `mustProvideValue[T any](t *testing.T, scope *Scope, value T)` helper

### Scope Constructor Pattern (Clone Group: 5 instances)
- `scope_provide_basic_test.go:18,20 + 32,34`
- `scope_provide_named_test.go:26,28 + 46,48`
- `scope_scoped_test.go:16,18`
- Pattern: `if err := ProvideValue(scope, X); err != nil { t.Fatalf(...) }`

### Scope Creation Pattern (Clone Group: 7 instances)
- `scope_child_test.go:20,22 + 61,63`
- `scope_scoped_test.go:66,68`
- `scope_new_test.go:19,21 + 54,56`
- `flow_context_branch_test.go:32,34 + 96,98`
- Pattern: `if child == nil { t.Fatal("expected child to not be nil") }`

### Flow Context Value Pattern (Clone Group: 10 instances)
- `flow_context_options_test.go:18,20 + 35,37 + 39,41`
- `flow_context_value_test.go:22,24 + 41,43 + 61,63`
- `option_test.go:151,153 + 188,190 + 207,209 + 214,216`
- Pattern: 3-line assertion blocks for context values

### Config/Option Assertion Pattern (Clone Group: 8 instances)
- `config_default_test.go:70,72`
- `config_tags_test.go:24,26 + 62,64 + 85,87 + 136,138`
- `option_test.go:381,383 + 393,395 + 405,407`

### CLI NewCLI Pattern (Clone Group: 18 instances)
- `cli_core_new_test.go:25,27 + 46,48 + 50,52`
- `cli_lifecycle_test.go:22,24`
- `duration_test.go:78,80 + 126,128`
- `types_filepath_test.go:59,61 + 77,79`
- `types_hostport_test.go:23,25 + 77,79 + 88,90 + 134,136`
- `types_port_test.go:21,23 + 84,86 + 139,141 + 174,176`

### Type MarshalText/UnmarshalText Pattern (Clone Group: 7 instances)
- `duration_test.go:163,165`
- `enum_test.go:159,161`
- `types_email_test.go:89,91`
- `types_filepath_test.go:104,106`
- `types_hostport_test.go:119,121`
- `types_port_test.go:159,161`
- `types_url_test.go:108,110`

### NoOpRunE Remaining (Clone Group: 32 instances)
- `command.go:45,49,53,178,186,195` — production code, not dedupable
- `cli_cobra_command_test.go:20,63`
- `cli_exec_test.go:45`
- `cli_hooks_test.go:86,89`
- `cli_lifecycle_test.go:154,198`
- `command_options_test.go:158`
- `middleware_test.go:50,91,126,232,262`
- `test_helpers_test.go:53,112`
- `testhelpers_test.go:36,43`

### Flag Error Pattern (Clone Group: 7 instances)
- `cli_cobra_command_test.go:167,169`
- `cli_flags_test.go:132,134 + 147,149 + 162,164 + 177,179 + 192,194`
- `command_validate_test.go:69,71`

### Examples Code (Clone Group: 23 instances)
- `examples/typed/main.go` + `examples/validation/main.go` + integration tests
- These are example/demo code — lower priority

### CLI Lifecycle NoOp Remaining (Clone Group: 4 instances)
- `cli_cobra_command_test.go:44`
- `cli_exec_test.go:106`
- `cli_lifecycle_test.go:167,210`

---

## d) TOTALLY FUCKED UP

1. **`cli_lifecycle_test.go:197` — UNDEFINED FUNCTION REFERENCE**  
   ```
   vet: pkg/cmdguard/v2/cli_lifecycle_test.go:197:11: undefined: noOpRunEForTestCLIConfigWithFlags
   ```
   The function `noOpRunEForTestCLIConfigWithFlags` was deleted in the diff but line 197 still calls it. **The codebase does not compile.**

2. **Go build cache corruption** — recurring `cache entry not found` errors require `GOCACHE=$(mktemp -d)` workaround. Not related to our code but slows iteration.

3. **No commits were made during the previous session** — all changes are uncommitted unstaged, meaning any crash/loss would lose all work.

4. **Clone count unchanged** — Despite replacing ~11 inline lambdas, `art-dupl` still reports 86 groups because many of those clones appear in places we haven't touched yet.

---

## e) WHAT WE SHOULD IMPROVE

### Process Improvements
1. **Commit after every self-contained change** — the user explicitly asked for this
2. **Run `go vet` after every edit** — would have caught the broken reference immediately
3. **Don't start new work until current change compiles** — the session kept going despite a broken state
4. **Use `GOCACHE=$(mktemp -d)` for Go commands** — work around the cache corruption
5. **Actually run tests before declaring done** — tests were never run successfully

### Architecture Improvements
6. **`mustProvideValue[T]` helper** — eliminates 9+ instances of `ProvideValue` + error check
7. **`mustNewScope(t, name)` helper** — eliminates repeated `NewScope` + nil check pattern
8. **`mustAddCommand[T, F](t, cli, cmd)` helper** — already exists as `addCommand` but underused
9. **Consider a shared `testMarshalUnmarshal[T]` generic** — could unify type test patterns
10. **`pkg/testutil` already has assertion helpers** — should use them more consistently instead of manual `if err != nil { t.Fatalf(...) }`

### Code Quality
11. **`assertNotPanic` in `scope_integration_test.go` duplicates `testutil.AssertPanics`** — should use the existing utility
12. **Two separate test helper files** (`test_helpers_test.go` in `v2` and `testhelpers_test.go` in `v2_test`) — by design due to Go package rules, but could be better documented
13. **`makeHookRunE` couldn't be genericized** due to Go type inference limitations — this is a language limitation, not fixable

---

## f) Top 25 Things We Should Get Done Next

### Priority 1: Fix broken code (MUST DO FIRST)
1. Fix `cli_lifecycle_test.go:197` — replace `noOpRunEForTestCLIConfigWithFlags` with `NoOpRunEWithFlags`
2. Run `go vet` and verify clean compilation
3. Run full test suite and verify all pass
4. Commit: "fix: replace remaining deleted function reference in cli_lifecycle_test.go"

### Priority 2: Commit existing work
5. Commit all current deduplication changes with detailed message
6. Re-run `art-dupl` to get baseline numbers post-commit

### Priority 3: Scope test dedup (high impact, medium effort)
7. Add `mustProvideValue[T any](t *testing.T, scope *Scope, value T)` to `test_helpers_test.go`
8. Replace `ProvideValue` + error check in `scope_lifecycle_test.go` (4 instances)
9. Replace in `scope_child_test.go:assertChildInheritsParent` (1 instance)
10. Replace in `scope_integration_test.go` (1 instance)
11. Replace in `scope_provide_basic_test.go` (3 instances)
12. Run tests, commit

### Priority 4: Remaining noOpRunE replacements in `v2` package
13. Replace in `cli_cobra_command_test.go:20,63` (2 instances)
14. Replace in `cli_exec_test.go:45` (1 instance)
15. Replace in `cli_hooks_test.go:86,89` (2 instances)
16. Replace in `command_options_test.go:158` (1 instance)
17. Replace in `middleware_test.go:50,91,126,232,262` (5 instances)
18. Run tests, commit

### Priority 5: Flow context and option test dedup
19. Add `assertFlowContextValue` helper for 3-line assertion pattern
20. Replace in `flow_context_options_test.go`, `flow_context_value_test.go`, `option_test.go`
21. Run tests, commit

### Priority 6: Config/option assertion dedup
22. Add helper for `if err != nil { t.Fatalf("expected no error, got: %v", err) }` pattern
23. Replace in `config_default_test.go`, `config_tags_test.go`, `option_test.go`
24. Run tests, commit

### Priority 7: Cross-cutting improvements
25. Replace `assertNotPanic` in `scope_integration_test.go` with `testutil.AssertPanics` (inverted)
26. Run final `art-dupl` and document improvement
27. Update `AGENTS.md` with deduplication patterns
28. `git push`

---

## g) Top #1 Question

**None that I can't figure out.** The path forward is clear: fix the broken reference, commit, then systematically deduplicate working from highest-impact clone groups downward. The only external blocker was the Go cache corruption, which has a known workaround (`GOCACHE=$(mktemp -d)`).

---

## Clone Detection Baseline

```
art-dupl --semantic -t 15 --sort total-tokens
→ 86 clone groups total
→ Top group: 32 instances (noOpRunE pattern)
→ Second group: 23 instances (CLI constructor pattern in examples)
→ Third group: 18 instances (NewCLI + assertion pattern)
```

## Files Modified (Uncommitted)

| File | Change | Status |
|------|--------|--------|
| `test_helpers_test.go` | Added `noOpRunE[T,F]`, refactored `noOpHandler*` | Clean |
| `cli_groups_test.go` | Replaced 7 inline lambdas | Clean |
| `middleware_test.go` | Replaced 3 inline lambdas | Clean |
| `cli_lifecycle_test.go` | Deleted function, replaced 1/2 references | **BROKEN** |
| `cli_command.go` | Added explicit type param `buildChain[T]` | Clean |

## Test Status

**UNKNOWN** — `go vet` fails with `undefined: noOpRunEForTestCLIConfigWithFlags`. Tests cannot compile.
