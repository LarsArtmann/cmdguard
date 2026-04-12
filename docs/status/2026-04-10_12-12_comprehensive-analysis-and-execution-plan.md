# Comprehensive Analysis & Execution Plan: cmdguard v2.2.0

**Date:** 2026-04-10 12:12
**Status:** Analysis Complete, Awaiting Execution
**Commits Ahead:** 4 (clean working tree)

---

## 1. WHAT I FORGOT / COULD HAVE DONE BETTER

### Critical Mistakes

1. **Didn't audit staged changes before committing** — I committed `5e027ac` without realizing there were already staged changes that could have had the same `:=` bug I was supposed to be fixing. I should have run `git diff --cached` and `go build` before every commit.

2. **Didn't verify build before assuming broken** — I saw "syntax error" in the build output and assumed the file was broken, when it was actually a transient cache issue or my misreading. Always verify with `go build ./...` directly.

3. **Didn't check for disk space** — The "no space left on device" error wasted time. Should have checked `df -h` early.

4. **Committed linter-generated changes piecemeal** — Instead of batching all the linter formatting changes, I committed them in small pieces. This creates noise in the commit history. Better to run `golangci-lint fmt ./...` once, verify, then commit the batch.

5. **Didn't verify tests with `-race` before declaring done** — The non-race tests passed, but the race detector caught issues I missed. Always run with `-race`.

### What Went Well

1. **Found and fixed the `slicesEqual` issue properly** — After the initial mistake, I grep'd for all usages and fixed them consistently.

2. **Created comprehensive execution plan** — The plan at `2026-04-10_10-05` is solid and covers the right priorities.

3. **Added proper test helpers** — `assertNoError`, `assertPanics`, `assertPanicsWithMessage` are genuinely useful and reduce duplication.

4. **DRY'd up examples with `fatal()` helper** — Consistent error handling pattern across all examples.

---

## 2. COMPREHENSIVE MULTI-STEP EXECUTION PLAN

### Pre-Conditions (Must Do First)

| Step | Task                              | Command                                      | Verification               |
| ---- | --------------------------------- | -------------------------------------------- | -------------------------- |
| P0   | Free disk space                   | `df -h` / clean unnecessary files            | `go clean -cache` succeeds |
| P1   | Run tests with race               | `go test ./... -count=1 -timeout 120s -race` | All 11 pass                |
| P2   | Run lint                          | `golangci-lint run ./...`                    | Count issues               |
| P3   | Commit any pending linter changes | `git add -A && git commit`                   | Clean tree                 |

### Phase 1: Blockers (Before Any Feature Work)

**BLOCKED BY:** Answer to Q1 (dependency direction)

| Step | Task                     | Files                                  | Impact               | Effort |
| ---- | ------------------------ | -------------------------------------- | -------------------- | ------ |
| 1.1  | Answer Q1                | Discussion                             | UNBLOCKS ALL         | 30 min |
| 1.2  | Document decision in ADR | `docs/adr/001-go-output-dependency.md` | Architecture clarity | 30 min |

### Phase 2: P0 — go-output Integration

| Step | Task                              | Files         | Impact                | Effort | Depends |
| ---- | --------------------------------- | ------------- | --------------------- | ------ | ------- |
| 2.1  | Add "Output Formatting" to README | `README.md`   | HIGH - user guidance  | 1 hr   | 1.1     |
| 2.2  | Create `examples/output/`         | New dir       | HIGH - demonstration  | 2 hr   | 1.1     |
| 2.3  | Verify bridge compiles            | Manual test   | MED - validation      | 30 min | 2.2     |
| 2.4  | Add to FEATURES.md                | `FEATURES.md` | MED - discoverability | 30 min | 2.2     |

### Phase 3: P1 — go-business-rules Integration

| Step | Task                          | Files                        | Impact               | Effort | Depends |
| ---- | ----------------------------- | ---------------------------- | -------------------- | ------ | ------- |
| 3.1  | Add "Validation" to README    | `README.md`                  | HIGH - user guidance | 1 hr   | -       |
| 3.2  | Create `examples/validation/` | New dir                      | HIGH - demonstration | 2 hr   | -       |
| 3.3  | Design `WithValidation()` API | `pkg/cmdguard/v2/command.go` | MED - clean API      | 2 hr   | 3.1     |
| 3.4  | Implement `WithValidation()`  | `pkg/cmdguard/v2/command.go` | MED - feature        | 1 hr   | 3.3     |

### Phase 4: P2 — Code Quality

| Step | Task                       | Files                | Impact                | Effort | Depends |
| ---- | -------------------------- | -------------------- | --------------------- | ------ | ------- |
| 4.1  | Triage lint issues         | ALL                  | MED - codebase health | 2 hr   | P3      |
| 4.2  | Fix forbidigo (20 issues)  | Examples             | LOW - reduce noise    | 1 hr   | 4.1     |
| 4.3  | Fix varnamelen (50 issues) | Tests                | LOW - style           | 1 hr   | 4.1     |
| 4.4  | Add CI workflow            | `.github/workflows/` | HIGH - quality gate   | 2 hr   | -       |
| 4.5  | Improve example coverage   | `examples/*`         | MED - demos           | 3 hr   | P3      |

### Phase 5: P3 — Type System Improvements

| Step | Task                       | Files                              | Impact               | Effort | Depends |
| ---- | -------------------------- | ---------------------------------- | -------------------- | ------ | ------- |
| 5.1  | Add `Option[T]` marshaling | `pkg/cmdguard/v2/types_option.go`  | MED - config interop | 1 hr   | P3      |
| 5.2  | Add `BrandedID[Brand]`     | `pkg/cmdguard/v2/types_branded.go` | HIGH - type safety   | 2 hr   | P3      |
| 5.3  | Add `Result[T, E]`         | `pkg/cmdguard/v2/types_result.go`  | MED - error handling | 2 hr   | P3      |
| 5.4  | Audit slices/maps usage    | ALL                                | LOW - modern Go      | 1 hr   | P3      |

### Phase 6: P4 — Ecosystem

| Step | Task                       | Files       | Impact                | Effort | Depends  |
| ---- | -------------------------- | ----------- | --------------------- | ------ | -------- |
| 6.1  | Create `docs/ECOSYSTEM.md` | New         | MED - documentation   | 1 hr   | 2.4, 3.2 |
| 6.2  | Update AGENTS.md           | `AGENTS.md` | MED - dev reference   | 30 min | 2.4, 3.2 |
| 6.3  | Cross-link repos           | All READMEs | LOW - discoverability | 30 min | 6.1      |

---

## 3. WORK REQUIRED VS IMPACT SORTING

### Quick Wins (High Impact, <1hr each)

| Priority | Task                         | Impact       | Effort | Risk |
| -------- | ---------------------------- | ------------ | ------ | ---- |
| 1        | Answer Q1                    | UNBLOCKS ALL | 30 min | LOW  |
| 2        | Fix forbidigo in examples    | LOW          | 1 hr   | LOW  |
| 3        | Document go-output in README | HIGH         | 1 hr   | LOW  |
| 4        | Create output example        | HIGH         | 2 hr   | LOW  |
| 5        | Add CI workflow              | HIGH         | 2 hr   | MED  |

### Medium Effort (1-2hr each)

| Priority | Task                          | Impact | Effort | Risk |
| -------- | ----------------------------- | ------ | ------ | ---- |
| 6        | Create validation example     | HIGH   | 2 hr   | LOW  |
| 7        | Document validation in README | HIGH   | 1 hr   | LOW  |
| 8        | Triage lint issues            | MED    | 2 hr   | LOW  |
| 9        | Add Option[T] marshaling      | MED    | 1 hr   | LOW  |
| 10       | Add BrandedID[Brand]          | HIGH   | 2 hr   | MED  |

### Heavy Lifting (2+hr each)

| Priority | Task                       | Impact | Effort | Risk |
| -------- | -------------------------- | ------ | ------ | ---- |
| 11       | Implement WithValidation() | MED    | 3 hr   | MED  |
| 12       | Add Result[T, E]           | MED    | 2 hr   | MED  |
| 13       | Improve example coverage   | MED    | 3 hr   | LOW  |
| 14       | Create ECOSYSTEM.md        | MED    | 1 hr   | LOW  |
| 15       | Audit slices/maps          | LOW    | 1 hr   | LOW  |

---

## 4. EXISTING CODE THAT FITS REQUIREMENTS

### For Validation (go-business-rules)

**EXISTING:**

```go
// pkg/cmdguard/v2/command.go
type Command[T any, F any] struct {
    PreRunE func(ctx context.Context, cfg *T, flags F) error
    RunE    func(ctx context.Context, cfg *T, flags F) error
    PostRunE func(ctx context.Context, cfg *T, flags F) error
}
```

**FIT:** PreRunE is the perfect place to call businessrules validation

**GAP:** No first-class `WithValidation(rules ...Rule)` option

### For Output Formatting (go-output)

**EXISTING:** `go-output/cmdguard/` bridge already exists with:

- `EnumFlag[T]`
- `OutputFormatFlag`
- `ColorModeFlag`
- `SortByFlag`

**FIT:** These are designed for cmdguard compatibility

**GAP:** Not documented in cmdguard's README

### For Type Improvements

**Option[T]:** `pkg/cmdguard/v2/types_option.go` — exists, needs marshaling

**Enum[T]:** `pkg/cmdguard/v2/enum.go` — exists, could benefit from generics

**Duration:** `pkg/cmdguard/v2/duration.go` — exists, could add marshaling

**BrandedID:** Not in cmdguard, but in go-output. Could extract to shared or add to cmdguard.

---

## 5. TYPE MODEL ARCHITECTURE IMPROVEMENTS

### Current Type System

```
types_option.go   → Option[T] (Some/None pattern)
types_enum.go     → Enum (string-based enums)
types.go          → LogLevel, Duration, NoFlags, Email, URL, etc.
types_branded.go  → MISSING (could add BrandedID)
types_result.go   → MISSING (could add Result[T, E])
```

### Proposed Additions

#### 5.1 BrandedID[Brand]

```go
// Compile-time type safety for IDs
type BrandedID[Brand any] string

func NewBrandedID[Brand any](s string) BrandedID[Brand]
func (id BrandedID[Brand]) String() string
func (id BrandedID[Brand]) IsZero() bool
func (id BrandedID[Brand]) MarshalText() ([]byte, error)
func (id *BrandedID[Brand]) UnmarshalText([]byte) error
```

**Why:** Prevents mixing UserID with OrderID at compile time

#### 5.2 Result[T, E]

```go
type Result[T any, E error] struct { ... }

func Ok[T any, E error](v T) Result[T, E]
func Err[T any, E error](e E) Result[T, E]

func (r Result[T, E]) IsOk() bool
func (r Result[T, E]) IsErr() bool
func (r Result[T, E]) Unwrap() T          // panics if Err
func (r Result[T, E]) UnwrapOr(default T) T
func (r Result[T, E]) UnwrapOrElse(func(E) T) T
func (r Result[T, E]) Map(func(T) T) Result[T, E]
func (r Result[T, E]) MapErr(func(E) E) Result[T, E]
```

**Why:** Explicit error handling, composable, avoids panic-on-nil-pattern

#### 5.3 Option[T] Marshaling

```go
func (o Option[T]) MarshalJSON() ([]byte, error)
func (o *Option[T]) UnmarshalJSON([]byte) error
func (o Option[T]) MarshalText() ([]byte, error)
func (o *Option[T]) UnmarshalText([]byte) error
```

**Why:** Config files (koanf) can serialize/deserialize optional values

#### 5.4 Enum Improvements

Current: `Enum` is a single type with `Allowed()` slice

Possible: `TypedEnum[T ~string]` with compile-time allowed values

**Challenge:** Go doesn't support const generics yet

---

## 6. ESTABLISHED LIBRARIES TO LEVERAGE

### Already Used (Keep)

| Library                     | Purpose              | Status    |
| --------------------------- | -------------------- | --------- |
| `github.com/spf13/cobra`    | CLI framework        | Essential |
| `github.com/samber/do/v2`   | Dependency injection | Essential |
| `charm.land/fang/v2`        | Styling              | Good      |
| `github.com/knadh/koanf/v2` | Configuration        | Flexible  |

### Consider Adding

| Library                       | Purpose                      | Pro                       | Con                       |
| ----------------------------- | ---------------------------- | ------------------------- | ------------------------- |
| `github.com/samber/mo`        | Option, Result, Either types | Saves implementation time | Another dependency        |
| `github.com/stretchr/testify` | Test assertions              | Expressive                | Another dep, stdlib works |
| `golang.org/x/exp/slices`     | Already stdlib `slices`      | Done ✓                    | N/A                       |
| `golang.org/x/exp/maps`       | Already stdlib `maps`        | Done ✓                    | N/A                       |

### Companion Libraries (From Analysis)

| Library             | Recommendation | Integration                    |
| ------------------- | -------------- | ------------------------------ |
| `go-output`         | P0             | Documentation + examples       |
| `go-business-rules` | P1             | PreRunE pattern                |
| `go-filewatcher`    | P2             | Watch command (if re-licensed) |
| `gogenfilter`       | P3             | Codegen (if re-licensed)       |

---

## 7. QUESTIONS I CANNOT ANSWER

### #1 CRITICAL BLOCKER

**Q1: What is the intended dependency direction between cmdguard and go-output?**

**Options:**

- A: cmdguard imports go-output
- B: go-output imports cmdguard
- C: Neither imports the other (current, doc-only)
- D: Separate adapter module

**Why I can't decide:** Product architecture decision. Depends on:

- "Batteries included" vs "minimal core" philosophy
- Long-term maintenance commitment
- Whether go-output is "the" output library or "a" output library

**IMPACT:** ALL P0 work is blocked until this is answered.

### #2 Secondary

**Q2: Should branded ID support be in cmdguard or a separate types library?**

- cmdguard: Common CLI need (every CLI has IDs)
- Separate: Not CLI-specific, could be reused across projects

### #3 Secondary

**Q3: Should we add `github.com/samber/mo` for Option/Result/Either, or implement ourselves?**

- mo: Faster to implement, well-tested
- Ours: Zero deps, matches our exact patterns

### #4 Secondary

**Q4: What's the timeline for go-output v1.0.0?**

- Blocks: Cannot recommend unversioned library
- Impact: P0 documentation quality

### #5 Secondary

**Q5: Should go-filewatcher and gogenfilter be re-licensed to MIT?**

- Current: Proprietary
- Impact: Cannot recommend as companion libraries

---

## 8. IMMEDIATE NEXT ACTIONS

1. **Answer Q1** (UNBLOCKS EVERYTHING)
2. Free disk space, verify tests pass
3. Run `golangci-lint fmt ./...`, commit batch
4. Implement Phase 2 (go-output) or Phase 3 (validation) based on Q1 answer

---

## 9. SESSION METRICS

| Metric               | Value                          |
| -------------------- | ------------------------------ |
| Commits this session | 4                              |
| Files changed        | ~15 files                      |
| Lines changed        | ~300 net                       |
| Working tree         | Clean                          |
| Tests                | 11/11 pass (before disk issue) |
| Blockers             | 1 (Q1)                         |

---

## 10. COMMITS THIS SESSION

```
94b4f9f style(tests): reformat long function signatures
eb03e1c chore: minor formatting and indentation fixes
67506bf docs: add session status report for 2026-04-10_10-53
5e027ac refactor(examples,tests): extract fatal() and test assertion helpers
```

---

**READY FOR:**

- [ ] Answer Q1 to unblock P0
- [ ] Execute phases in order
- [ ] Push when done

_Generated at 2026-04-10 12:12_
