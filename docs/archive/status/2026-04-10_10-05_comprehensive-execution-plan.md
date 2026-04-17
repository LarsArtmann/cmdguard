# Comprehensive Execution Plan: cmdguard v2.2.0 Roadmap

**Date:** 2026-04-10 10:05
**Status:** Draft for Review
**Based on:** Cross-library integration review (2026-04-10_09-37)

---

## 1. WHAT I FORGOT / COULD HAVE DONE BETTER

### Forgottens

1. **No compilation verification before the bulk formatting commit** — The prior session committed test files with `:=` shadowing bugs. A `go build ./...` gate would have caught this.
2. **Didn't verify all consumers of slicesEqual** — The `cli_hooks_test.go` had its own `slicesEqual` removed but other files still used the shared helper in `test_helpers_test.go`. Should have grepped for all usages first.
3. **Didn't check example file modifications** — The unstaged changes included `examples/typed/main_advanced_test.go` which had DRY improvements (greetRunE() extraction). These were good changes but weren't explicitly tracked.

### Improvements for Future Sessions

1. **Pre-commit check script** — Always run: `go build ./... && go test ./... -count=1 -timeout 60s`
2. **Change impact analysis** — Before modifying shared helpers, grep for all usages across the codebase
3. **Explicit tracking of all modified files** — Don't just fix the immediate error; audit all related changes
4. **Use `slices` package consistently** — Go 1.21+ has `slices.Equal`, `slices.Contains`, `slices.Sort`. Audit for opportunities.

---

## 2. COMPREHENSIVE MULTI-STEP EXECUTION PLAN

### Phase 1: Critical Fixes (Immediate)

| Step | Task                                  | Files                                  | Verification                       |
| ---- | ------------------------------------- | -------------------------------------- | ---------------------------------- |
| 1.1  | Add pre-commit compilation gate       | `justfile` or `Makefile`               | `just precommit` runs build + test |
| 1.2  | Document the #1 architecture question | `docs/adr/001-dependency-direction.md` | ADR format                         |

### Phase 2: P0 Integration — go-output Companion

| Step | Task                                      | Work Required                         | Impact                   |
| ---- | ----------------------------------------- | ------------------------------------- | ------------------------ |
| 2.1  | **DECISION REQUIRED:** Answer #1 question | Discussion with Lars                  | Blocks everything below  |
| 2.2  | Add "Output Formatting" section to README | `README.md` ~30 lines                 | High user value          |
| 2.3  | Create `examples/output-formats/`         | New dir, ~150 lines                   | Demonstrates integration |
| 2.4  | Verify go-output bridge compiles          | CI check or manual                    | Prevents bitrot          |
| 2.5  | Add output format auto-detection helper   | `pkg/cmdguard/v2/output.go` ~50 lines | UX improvement           |

### Phase 3: P1 Integration — go-business-rules Companion

| Step | Task                                   | Work Required                | Impact               |
| ---- | -------------------------------------- | ---------------------------- | -------------------- |
| 3.1  | Add "Validation" section to README     | `README.md` ~40 lines        | User guidance        |
| 3.2  | Create `examples/validation/`          | New dir, ~200 lines          | Demonstrates pattern |
| 3.3  | Design `WithValidation()` API          | `pkg/cmdguard/v2/command.go` | Optional integration |
| 3.4  | Implement `WithValidation()` option    | ~30 lines                    | Clean API            |
| 3.5  | Add validation to one existing example | `examples/typed/`            | Real usage           |

### Phase 4: Code Quality Improvements

| Step | Task                                | Work Required                 | Impact             |
| ---- | ----------------------------------- | ----------------------------- | ------------------ |
| 4.1  | Triage 160 lint issues              | Review each category          | Cleaner codebase   |
| 4.2  | Fix forbidigo issues in examples    | 20 findings                   | Removes lint noise |
| 4.3  | Improve example coverage            | `examples/typed/`: 3.6% → 50% | Better demos       |
| 4.4  | Use `slices` package everywhere     | Audit + replace               | Modern Go          |
| 4.5  | Use `maps` package where applicable | Go 1.21+ feature              | Modern Go          |

### Phase 5: Type System Architecture Improvements

| Step | Task                                    | Work Required                | Impact                 |
| ---- | --------------------------------------- | ---------------------------- | ---------------------- |
| 5.1  | Add `Option[T]` marshaling support      | `types_option.go`            | JSON interop           |
| 5.2  | Consider `Result[T, E]` type            | New file ~100 lines          | Error handling pattern |
| 5.3  | Add branded ID support                  | `types_branded.go` ~50 lines | Type safety            |
| 5.4  | Review Enum[T] for generic improvements | `enum.go`                    | Cleaner API            |
| 5.5  | Add Duration marshaling helpers         | `duration.go`                | Config interop         |

### Phase 6: Ecosystem & Documentation

| Step | Task                       | Work Required              | Impact                  |
| ---- | -------------------------- | -------------------------- | ----------------------- |
| 6.1  | Create `docs/ECOSYSTEM.md` | New file ~100 lines        | Companion library guide |
| 6.2  | Update FEATURES.md         | Add companion section      | User discovery          |
| 6.3  | Update AGENTS.md           | Integration patterns       | Dev reference           |
| 6.4  | Add CI workflow            | `.github/workflows/ci.yml` | Quality gates           |
| 6.5  | Cross-link companion repos | README badges              | Ecosystem coherence     |

---

## 3. WORK REQUIRED VS IMPACT SORTING

### High Impact, Low Work (Quick Wins)

1. **Add pre-commit gate** — 5 min, prevents bugs
2. **Document go-output bridge** — 30 min, high user value
3. **Create output-formats example** — 1 hr, demonstrates integration
4. **Triage lint issues** — 1 hr, codebase health
5. **Fix forbidigo in examples** — 30 min, reduces noise

### High Impact, Medium Work

6. **Answer #1 architecture question** — Discussion, unblocks P0
7. **Add validation example** — 2 hrs, demonstrates pattern
8. **Implement `WithValidation()`** — 2 hrs, clean API
9. **Improve example coverage** — 3 hrs, better demos
10. **Add CI workflow** — 2 hrs, quality gates

### Medium Impact, Medium Work

11. **Add Option[T] marshaling** — 2 hrs, JSON interop
12. **Add branded ID support** — 2 hrs, type safety
13. **Create ECOSYSTEM.md** — 1 hr, companion guide
14. **Use slices/maps everywhere** — 2 hrs, modern Go
15. **Add Duration helpers** — 1 hr, config interop

### High Impact, High Work

16. **Design Result[T, E] type** — 4 hrs, error handling
17. **Refactor Enum[T] generics** — 4 hrs, API cleanup
18. **Tag go-output v1.0.0** — External repo coordination
19. **Re-license companion libraries** — Legal/licensing
20. **Create cmdguard-showcase repo** — New repo, full demo

---

## 4. EXISTING CODE THAT FITS REQUIREMENTS

### For Validation Integration

- **Existing:** `PreRunE` hooks in `pkg/cmdguard/v2/command.go`
- **Pattern:** Commands already support `PreRunE: func(ctx, cfg, flags) error`
- **Fit:** go-business-rules validation can be called in `PreRunE`
- **Gap:** No first-class `WithValidation()` option ( Phase 3.3 )

### For Output Formatting

- **Existing:** `go-output/cmdguard/` bridge
- **Pattern:** `EnumFlag[T]`, `OutputFormatFlag`, `ColorModeFlag`, `SortByFlag`
- **Fit:** Already designed for cmdguard compatibility
- **Gap:** Not documented in cmdguard's README ( Phase 2.2 )

### For Type Improvements

- **Existing:** `Option[T]` in `types_option.go`
- **Pattern:** Rust-like Some/None with `IsSome()`, `IsNone()`, `Unwrap()`
- **Gap:** No JSON marshaling ( Phase 5.1 )

- **Existing:** `Enum[T]` in `enum.go`
- **Pattern:** String-based enums with validation
- **Gap:** Could benefit from branded type support ( Phase 5.3 )

- **Existing:** `Duration` type
- **Pattern:** Wrapper around `time.Duration` with flag parsing
- **Gap:** Config marshaling helpers ( Phase 5.5 )

### For Branded IDs

- **Existing:** go-output uses `BrandedID[Brand]`
- **Pattern:** `type BrandedID[Brand any] string`
- **Fit:** Could be extracted to shared types or added to cmdguard
- **Gap:** cmdguard has no branded ID support ( Phase 5.3 )

---

## 5. TYPE MODEL ARCHITECTURE IMPROVEMENTS

### Current State

```go
// Option[T] - Rust-like optional type
type Option[T any] struct { value T; valid bool }

// Enum[T] - String-based enum
type Enum struct { value string; allowed []string }

// Duration - Time duration wrapper
type Duration struct { value time.Duration }
```

### Proposed Additions

#### 5.1 Option[T] Marshaling

```go
// Add JSON marshaling for config interop
func (o Option[T]) MarshalJSON() ([]byte, error)
func (o *Option[T]) UnmarshalJSON(data []byte) error
```

#### 5.2 BrandedID[Brand]

```go
// Compile-time ID type safety
package types

type BrandedID[Brand any] string

func NewBrandedID[Brand any](s string) BrandedID[Brand]
func (id BrandedID[Brand]) String() string
func (id BrandedID[Brand]) IsZero() bool

// Usage:
type UserBrand struct{}
type UserID = BrandedID[UserBrand]

var userID UserID = types.NewBrandedID[UserBrand]("usr_123")
// Cannot accidentally use UserID where OrderID is expected
```

#### 5.3 Result[T, E]

```go
// Explicit error handling (similar to Rust Result)
package types

type Result[T any, E error] struct {
    ok  T
    err E
}

func Ok[T any, E error](v T) Result[T, E]
func Err[T any, E error](e E) Result[T, E]
func (r Result[T, E]) IsOk() bool
func (r Result[T, E]) IsErr() bool
func (r Result[T, E]) Unwrap() T          // panics if Err
func (r Result[T, E]) UnwrapOr(default T) T
func (r Result[T, E]) UnwrapOrElse(f func(E) T) T
func (r Result[T, E]) Expect(msg string) T // panics with msg if Err
```

#### 5.4 Enum[T] Generics (if Go allows)

Go doesn't support generic type parameters on methods, so `Enum[T]` is limited. However, we could add:

```go
// TypedEnum is a compile-time safe enum
type TypedEnum[T ~string] struct {
    value T
}

func NewTypedEnum[T ~string](v T) TypedEnum[T]
func (e TypedEnum[T]) Value() T
```

#### 5.5 FlagSet Type

```go
// For complex flag interdependencies
package types

type FlagSet struct {
    flags map[string]Flag
}

func (fs *FlagSet) Register(f Flag)
func (fs *FlagSet) Validate() error // cross-flag validation
func (fs *FlagSet) Parse(args []string) error
```

---

## 6. ESTABLISHED LIBRARIES TO LEVERAGE

### Already Used (Keep Using)

| Library                     | Purpose              | Verdict                    |
| --------------------------- | -------------------- | -------------------------- |
| `github.com/spf13/cobra`    | CLI framework        | Core dependency, essential |
| `github.com/samber/do/v2`   | Dependency injection | Clean DI, keep             |
| `charm.land/fang/v2`        | Cobra styling        | Good for UX, keep          |
| `github.com/knadh/koanf/v2` | Configuration        | Flexible config, keep      |

### Recommended Additions

| Library                                     | Purpose          | Work                    | Impact |
| ------------------------------------------- | ---------------- | ----------------------- | ------ |
| `golang.org/x/exp/slices` → stdlib `slices` | Slice operations | Replace custom helpers  | Medium |
| `golang.org/x/exp/maps` → stdlib `maps`     | Map operations   | Modern Go patterns      | Low    |
| `github.com/samber/mo`                      | Functional types | Option, Result, Either  | Medium |
| `github.com/google/uuid`                    | UUID generation  | For branded ID examples | Low    |
| `github.com/stretchr/testify`               | Test assertions  | More expressive tests   | Medium |

### For Companion Integration

| Library                                 | Companion      | Integration Pattern      |
| --------------------------------------- | -------------- | ------------------------ |
| `github.com/larsartmann/go-output`      | Recommended P0 | Documentation + examples |
| `github.com/artmann/businessrules`      | Recommended P1 | PreRunE + WithValidation |
| `github.com/larsartmann/go-filewatcher` | Consider P2    | Watch command example    |
| `github.com/larsartmann/gogenfilter`    | Consider P3    | Codegen command example  |

### Not Recommended (Per Review)

| Library              | Why Not                      |
| -------------------- | ---------------------------- |
| `universal-workflow` | Heavy deps, different domain |
| `go-cqrs-lite`       | Backend pattern, not CLI     |
| `go-localfirst`      | Reference app, not library   |
| `go-localsync`       | Private deps, sync domain    |
| `go-plugin-mvp`      | MVP, WASM complexity         |

---

## 7. DETAILED QUESTIONS I CANNOT ANSWER

### Critical Blocker

**Q1: What is the intended dependency direction between cmdguard and go-output?**

- Option A: cmdguard imports go-output (tight coupling)
- Option B: go-output imports cmdguard (reverse coupling)
- Option C: Neither imports the other (current state, documentation-only)
- Option D: Separate adapter module (coordination overhead)

**Why I can't decide:** This is a product architecture decision, not a technical one. It depends on:

- Whether cmdguard wants to be "batteries included" or "minimal core"
- Whether go-output is positioned as "the" output library or "a" output library
- Long-term maintenance commitment to the integration

**What I need:** Explicit decision from Lars on which model to pursue.

### Secondary Questions

**Q2: Should branded ID support be in cmdguard or a separate types library?**

- Pro cmdguard: Common CLI need (IDs are everywhere)
- Pro separate: Not CLI-specific, could be reused

**Q3: Should we add testify for tests or keep stdlib-only?**

- Pro testify: `assert.Equal()`, `require.NoError()` are expressive
- Con testify: Another dependency; current stdlib tests work fine

**Q4: What's the timeline for go-output v1.0.0?**

- Blocks P0 documentation (can't recommend unversioned library)
- Needs coordination with go-output repo

**Q5: Are go-filewatcher and gogenfilter intended to be open-sourced?**

- Current proprietary license blocks integration
- Need license change or exclusion from companion ecosystem

---

## 8. NEXT IMMEDIATE ACTIONS (Priority Order)

1. **Answer Q1** (dependency direction) — unblocks all P0 work
2. **Add pre-commit gate** — prevents future compilation errors
3. **Document go-output bridge** — high user value, low work
4. **Create output-formats example** — demonstrates integration
5. **Triage lint issues** — codebase health

---

_This plan is a living document. Update as decisions are made and work progresses._

_Generated at 2026-04-10 10:05_
