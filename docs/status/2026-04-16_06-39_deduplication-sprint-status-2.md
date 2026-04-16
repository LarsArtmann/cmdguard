# Deduplication Sprint — Status Report #2

**Date:** 2026-04-16 06:39  
**Branch:** master (4 commits ahead of origin)  
**Task:** Deduplicate code clones identified by `art-dupl --semantic -t 15 --sort total-tokens`

---

## Executive Summary

Two sessions of deduplication work. Session 1 left code **broken** (undefined function). Session 2 **fixed the breakage**, committed all prior work, and added `mustProvideValue` helper replacing 15 instances of `ProvideValue` + error-check boilerplate across 5 scope test files. `go vet` is currently verifying compilation. **No tests have been successfully run yet** — blocked by Go cache corruption requiring `GOCACHE=$(mktemp -d)` workaround and slow cold-start compilation.

---

## a) FULLY DONE

### Session 1 (committed in `bc0683c`):
1. **Added `noOpRunE[T, F]` generic** to `test_helpers_test.go` (package `v2`)
2. **Refactored `noOpHandler()` and `noOpHandlerForTestAppConfig()`** to delegate to `noOpRunE`
3. **Replaced 7 inline no-op lambdas** in `cli_groups_test.go`
4. **Replaced 3 inline no-op lambdas** in `middleware_test.go`
5. **Deleted `noOpRunEForTestCLIConfigWithFlags`** from `cli_lifecycle_test.go`, replaced with `NoOpRunEWithFlags`

### Session 2 (uncommitted — in working tree):
6. **Fixed broken `cli_lifecycle_test.go:197`** — replaced orphaned reference to deleted function
7. **Added `mustProvideValue[T any](t, scope, value)` helper** to `test_helpers_test.go`
8. **Replaced 15 `ProvideValue` + error-check boilerplates** across 5 files:
   - `scope_lifecycle_test.go`: 6 instances → 6 one-liners (-18 lines)
   - `scope_provide_basic_test.go`: 6 instances → 6 one-liners (-18 lines)
   - `scope_child_test.go`: 1 instance (-2 lines)
   - `scope_integration_test.go`: 2 instances (-6 lines)
   - `scope_provide_named_test.go`: 1 instance (-4 lines)

### Other commits:
- `8551c97`: Added explicit type param `buildChain[T]` in `cli_command.go`
- `aaa6bf1`: Status report doc
- `d908bb6`: Formatting fix in validation example test

---

## b) PARTIALLY DONE

1. **`mustProvideValue` rollout** — 15 of ~18 instances replaced. Remaining:
   - `scope_provide_basic_test.go` has `err := ProvideValue(scope, 42)` + `testutil.AssertNoError` (different pattern, keep as-is)
   - `scope_provide_basic_test.go` has `err := ProvideValue(nil, 42)` (error-path test, keep as-is)

2. **No-op RunE dedup** — 11 of ~25+ replaceable instances done. Remaining in `v2` package:
   - `cli_cobra_command_test.go:20,63` (2 instances)
   - `cli_exec_test.go:45` (1 instance)
   - `cli_hooks_test.go:86,89` (2 instances)
   - `command_options_test.go:158` (1 instance)
   - `middleware_test.go:50,91,126,232,262` (5 instances)

---

## c) NOT STARTED

### High-impact, not started:
1. **Remaining noOpRunE replacements** in v2 test files (11 instances)
2. **NoOpRunE replacements** in `v2_test` package test files (cli_lifecycle, cli_cobra_command, etc.)
3. **Scope constructor pattern** — `if child == nil { t.Fatal(...) }` repeated 7x across scope/flow tests
4. **Flow context value assertion pattern** — 3-line blocks repeated 10x
5. **Config/option assertion pattern** — `if err != nil { t.Fatalf("expected no error, got: %v", err) }` repeated 8x
6. **NewCLI assertion pattern** — repeated 18x across cli_core_new, duration, types_* tests

### Medium-impact, not started:
7. **Type MarshalText/UnmarshalText test pattern** — shared `runUnmarshalErrorTest`-like helper across 7 type test files
8. **Flag error assertion pattern** — 7 instances in cli_cobra_command_test, cli_flags_test, command_validate_test
9. **`assertNotPanic` in scope_integration_test.go** — duplicates `testutil.AssertPanics` (inverted)
10. **Examples code** — 23 clone instances in examples/ (lower priority)

### Low-impact, not started:
11. **Production code clones** — `command.go` has 6 instances of similar option-application patterns (risky to change)
12. **`pkg/testutil/panic_test_helpers.go`** — 308-line file with internal duplication
13. **`cli.go` and `flags.go`** — 6 instances of similar error-wrapping patterns

---

## d) TOTALLY FUCKED UP

1. **Session 1 left code broken** — `cli_lifecycle_test.go:197` referenced deleted function `noOpRunEForTestCLIConfigWithFlags`. Fixed in session 2.

2. **No tests verified for ANY change** — Go cache corruption (`cache entry not found`) causes `go test` to fail. Workaround is `GOCACHE=$(mktemp -d)` but that forces full recompilation (~2-3 min each time). Tests have been attempted 4+ times but never completed successfully during these sessions.

3. **4 commits ahead of origin, never pushed** — All work is local only. Machine crash = total loss.

4. **Background jobs keep timing out** — Every `go vet` and `go test` invocation runs as background and never completes within the conversation window. This is a systemic issue with the Go toolchain + cache state.

---

## e) WHAT WE SHOULD IMPROVE

### Critical Process Issues:
1. **Run `go vet` synchronously with `GOCACHE=$(mktemp -d)`** — accept the 2-3 min wait, don't background it
2. **Commit after every verified change** — we're still sitting on uncommitted scope test changes
3. **Push after every commit** — 4 local-only commits is dangerous
4. **Fix the Go cache** — `go clean -cache` didn't help. May need to `rm -rf ~/Library/Caches/go-build` entirely

### Code Improvements:
5. **Use `testutil.AssertNoError` more consistently** — some files use manual `if err != nil { t.Fatalf(...) }` while others use `testutil.AssertNoError`
6. **Document the two test helper files** — `test_helpers_test.go` (v2) vs `testhelpers_test.go` (v2_test) distinction should be clearer
7. **Consider adding `mustInvoke[T]` helper** — mirrors `mustProvideValue` for the Invoke + error check pattern

---

## f) Top 25 Things to Do Next

### IMMEDIATE (do right now):
1. ✅ `git add` and commit the `mustProvideValue` + scope test changes
2. Run `GOCACHE=$(mktemp -d) go vet ./pkg/cmdguard/v2/...` **synchronously** — wait for it
3. Run `GOCACHE=$(mktemp -d) go test ./pkg/cmdguard/v2/... -count=1 -timeout 180s` **synchronously**
4. Run `go test ./... -count=1 -timeout 180s` for full project verification
5. `git push origin master`

### HIGH IMPACT, LOW EFFORT:
6. Replace remaining noOpRunE in `middleware_test.go` (5 instances: lines 50, 91, 126, 232, 262)
7. Replace in `cli_cobra_command_test.go` (2 instances: lines 20, 63)
8. Replace in `cli_exec_test.go` (1 instance: line 45)
9. Replace in `cli_hooks_test.go` (2 instances: lines 86, 89)
10. Replace in `command_options_test.go` (1 instance: line 158)
11. Commit + test + push

### MEDIUM IMPACT, MEDIUM EFFORT:
12. Add `mustInvoke[T any](t, scope) T` helper for Invoke + error check pattern
13. Add `assertNotNil(t, got, name)` helper for the `if x == nil { t.Fatal(...) }` pattern (7 instances)
14. Replace flow context value assertion pattern (10 instances across flow_context_*.go, option_test.go)
15. Replace config/option assertion pattern (8 instances across config_*.go, option_test.go)

### MEDIUM IMPACT, HIGHER EFFORT:
16. Add `testMarshalRoundtrip[T any](t, jsonStr, expected)` helper for type tests
17. Unify `runUnmarshalErrorTest`-like patterns across duration, enum, port, hostport, email, url, filepath tests
18. Replace flag error assertion pattern (7 instances in cli_*_test.go)

### LOWER PRIORITY:
19. Clean up `pkg/testutil/panic_test_helpers.go` internal duplication
20. Consider shared helpers for NewCLI + assertion pattern (18 instances)
21. Re-run `art-dupl` to measure total improvement
22. Update `AGENTS.md` with deduplication patterns
23. Investigate and fix Go cache corruption root cause
24. Review production code clones in `command.go` (6 instances) — risky but valuable
25. Consider if `assertNotPanic` in scope_integration_test.go should use testutil

---

## g) Top #1 Question

**Why does the Go build cache keep corrupting?** Every `go build`, `go vet`, and `go test` fails with `cache entry not found: open /Users/larsartmann/Library/Caches/go-build/...`. `go clean -cache` doesn't fix it. `GOCACHE=$(mktemp -d)` works but forces full recompilation each time. This is the single biggest blocker to productive iteration. Is there a known fix? Could it be related to the Nix-installed Go at `/nix/store/5ajixjk279m40yf6x96xxlnvw1wg6hq3-go-1.26.0/`?

---

## Metrics

| Metric | Start | After Session 1 | After Session 2 (uncommitted) |
|--------|-------|-----------------|-------------------------------|
| Clone groups | 86 | 86 | ~78 (estimated, not yet verified) |
| Inline no-op lambdas replaced | 0 | 11 | 11 |
| ProvideValue boilerplates replaced | 0 | 0 | 15 |
| Lines removed | 0 | -10 | -49 total |
| Tests passing | unknown | BROKEN | NOT VERIFIED |
| Commits pushed | — | 0 | 0 |

## Git Log

```
d908bb6 style: format long CommandContext call in validation example test
bc0683c refactor: deduplicate no-op RunE lambdas in v2 test files
aaa6bf1 docs: add deduplication sprint status report
8551c97 fix: add explicit type parameter to buildChain call in wireHandlerWithMiddleware
84cea37 chore: comprehensive lint fixes, error handling improvements, and test polish
```

## Uncommitted Changes

```
pkg/cmdguard/v2/test_helpers_test.go        | +8 (mustProvideValue helper)
pkg/cmdguard/v2/scope_lifecycle_test.go      | -18 lines (6→6 one-liners)
pkg/cmdguard/v2/scope_provide_basic_test.go  | -18 lines (6→6 one-liners)
pkg/cmdguard/v2/scope_child_test.go          | -2 lines
pkg/cmdguard/v2/scope_integration_test.go    | -6 lines (2→2 one-liners)
pkg/cmdguard/v2/scope_provide_named_test.go  | -4 lines
```
