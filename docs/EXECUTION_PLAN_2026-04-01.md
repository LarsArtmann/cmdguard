# Comprehensive Multi-Step Execution Plan

**Date:** 2026-04-01  
**Project:** cmdguard v2.1.0+  
**Status:** Package() function fixed, ready for next phase

---

## Reflection: What Was Forgotten / Could Be Improved

### 1. Testing Before Committing

**Issue:** The `Package()` function was committed with a duplicate registration bug that would have been caught by running tests.

**Lesson:** Always run tests before committing new features, especially when integrating with DI systems where registration order matters.

### 2. Understanding Initialization Flow

**Issue:** Didn't trace through `NewCLI` → `initialize()` → `ProvideValue()` to realize config was already registered.

**Lesson:** When adding integration helpers like `Package()`, trace the full initialization flow to avoid double-registration.

### 3. Option Precedence Bug

**Issue:** `initialize()` was unconditionally overwriting `cli.scope`, breaking `WithCLIScope` option.

**Lesson:** Options should always be checked before setting defaults. Pattern: `if field == nil { field = default }`

---

## Existing Patterns to Leverage

### Type-Safe Wrapper Pattern (from types.go)

```go
// Option[T] - zero-cost wrapper with semantic meaning
type Option[T any] struct {
    value T
    ok    bool
}

// Enum - validated string with allowed values
type Enum struct {
    value   string
    allowed []string
}

// Duration - wraps time.Duration with parsing
type Duration struct {
    duration time.Duration
}
```

**Key Insights:**

- Use struct wrappers for type safety
- Provide `Parse*` constructors that return `(T, error)`
- Implement standard interfaces (`String()`, `MarshalText()`, etc.)
- Keep zero values meaningful

### Error Pattern (from errors.go)

```go
// Sentinel errors for errors.Is() checking
var ErrInvalidCommand = errors.New("invalid command")

// Rich error types with context
type CommandError struct {
    Op  string
    Key string
    Err error
}

func (e *CommandError) Error() string { ... }
func (e *CommandError) Unwrap() error { return e.Err }
```

### Functional Options Pattern (from cli.go)

```go
type CLIOption[T any] func(*CLI[T])

func WithCLIVersion[T any](version string) CLIOption[T] {
    return func(cli *CLI[T]) {
        cli.version = version
    }
}
```

---

## Multi-Step Execution Plan

Sorted by **Impact vs Effort** matrix:

### 🔴 HIGH IMPACT, LOW EFFORT (Do First)

#### 1. Add t.Parallel() to Test Files

**Work:** ~30 minutes  
**Impact:** Faster test execution  
**Files:** `guarded_command_test.go`, other large test files  
**Pattern:** Add `t.Parallel()` at start of each test, ensure no shared state

#### 2. Fix CLI[T] AddCommand Flag Parsing Bug

**Work:** ~1 hour  
**Impact:** Critical bug fix  
**Source:** TODO_LIST.md mentions bug at `cli.go:190`  
**Pattern:** Use `cloneAndParseFlags` pattern from existing code

#### 3. Add Custom Types (URL, Email, Port, FilePath)

**Work:** ~2 hours  
**Impact:** Enhanced type safety for common CLI inputs  
**Pattern:** Follow `Duration` type pattern from `types.go`

**Example:**

```go
// URL validates and wraps a URL string
type URL struct {
    url *url.URL
    raw string
}

func ParseURL(s string) (URL, error) {
    u, err := url.Parse(s)
    if err != nil {
        return URL{}, NewURLError(s, err)
    }
    return URL{url: u, raw: s}, nil
}
```

### 🟡 HIGH IMPACT, MEDIUM EFFORT

#### 4. Implement Result[T] Type

**Work:** ~3 hours  
**Impact:** Better error handling ergonomics  
**Pattern:** Similar to Rust's Result<T, E> or Go's optional with error

**Design:**

```go
// Result[T] represents either a success value or an error
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](v T) Result[T] { return Result[T]{value: v} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

func (r Result[T]) IsOk() bool { return r.err == nil }
func (r Result[T]) IsErr() bool { return r.err != nil }
func (r Result[T]) Unwrap() (T, error) { return r.value, r.err }
```

#### 5. Add Progress/Spinner Type

**Work:** ~4 hours  
**Impact:** Better UX for long-running operations  
**Library:** Use `charmbracelet/bubbles`  
**Pattern:** Wrap bubbles spinner with our own type

#### 6. Implement Validated[T] Wrapper

**Work:** ~3 hours  
**Impact:** Runtime validation for config fields  
**Pattern:** Decorator pattern around any type

```go
// Validated[T] wraps a value with validation rules
type Validated[T any] struct {
    value T
    valid bool
    rules []Validator[T]
}

type Validator[T any] func(T) error
```

### 🟢 MEDIUM IMPACT, LOW EFFORT

#### 7. Add Short Flags for Common Options

**Work:** ~1 hour  
**Impact:** Better CLI UX  
**Pattern:** Update struct tag parsing in `flags.go`

**Example:**

```go
type Config struct {
    Verbose bool `flag:"verbose" short:"v" default:"false"`
    Output  string `flag:"output" short:"o" default:"stdout"`
}
```

#### 8. Show Defaults in Help Text

**Work:** ~2 hours  
**Impact:** Better user experience  
**Pattern:** Hook into cobra's help generation

#### 9. Add Shell Completion Helpers

**Work:** ~3 hours  
**Impact:** Better CLI UX  
**Pattern:** Use cobra's completion system

### 🟠 MEDIUM IMPACT, MEDIUM EFFORT

#### 10. Replace internal/config with koanf

**Work:** ~6 hours  
**Impact:** More robust config handling  
**Library:** `github.com/knadh/koanf/v2`  
**Pattern:** Already used in AGENTS.md examples

#### 11. Replace internal/logging with charmbracelet/log

**Work:** ~4 hours  
**Impact:** Better log aesthetics  
**Library:** `github.com/charmbracelet/log`  
**Note:** Check if this is worth the dependency

#### 12. Add Middleware Support

**Work:** ~6 hours  
**Impact:** Extensible command processing  
**Pattern:** Chain of responsibility

```go
type Middleware[T any] func(Command[T]) Command[T]

func (cli *CLI[T]) Use(mw ...Middleware[T])
```

### ⚪ LOW IMPACT, HIGH EFFORT (Defer)

#### 13. Create v3 API (pkg/cmdguard/v3)

**Work:** ~40 hours  
**Impact:** Clean slate design  
**Note:** Not needed until v2 limitations become painful

#### 14. Full koanf Integration with Auto-Loading

**Work:** ~16 hours  
**Impact:** Automatic config file watching  
**Note:** Can be done incrementally

#### 15. Plugin System for Validators

**Work:** ~20 hours  
**Impact:** User-extensible validation  
**Note:** Overkill for current needs

---

## Immediate Next Steps (Priority Order)

### Step 1: Fix Known Bugs (This Session)

- [ ] Investigate CLI[T] AddCommand flag parsing bug (cli.go:190)
- [ ] Add t.Parallel() to test files
- [ ] Add tests for initialize error paths

### Step 2: Type Safety Improvements (Next Session)

- [ ] Add URL type with validation
- [ ] Add Email type with validation
- [ ] Add Port type with range validation
- [ ] Add FilePath type with existence checks

### Step 3: UX Improvements (Future Session)

- [ ] Show defaults in help text
- [ ] Add short flags support
- [ ] Add shell completion helpers

### Step 4: Advanced Features (Future)

- [ ] Implement Result[T] type
- [ ] Add Progress/Spinner type
- [ ] Middleware support

---

## Using Established Libraries

### Already Integrated

- ✅ `samber/do/v2` - Dependency injection
- ✅ `charmbracelet/fang` - Cobra styling
- ✅ `spf13/cobra` - CLI framework

### Recommended Additions

- 📦 `charmbracelet/bubbles` - Progress/spinner (lightweight)
- 📦 `knadh/koanf/v2` - Config management (replaces internal/config)
- 📦 `charmbracelet/log` - Structured logging (optional)

### Avoid For Now

- ❌ Heavy validation libraries (keep it simple)
- ❌ Plugin systems (overkill)
- ❌ Complex middleware chains (YAGNI)

---

## Type Architecture Improvements

### Current Strengths

- `Option[T]` - Optional values
- `Enum` - Validated strings
- `Duration` - Time parsing
- `LogLevel` - Typed log levels

### Missing Patterns

- `Result[T]` - Error handling
- `Validated[T]` - Runtime validation
- `URL`, `Email`, `Port` - Domain types
- `NonEmpty[T]` - Collection constraints

### Design Principles

1. **Zero-cost abstractions** - Wrappers should compile away
2. **Parse constructors** - `ParseT(string) (T, error)` pattern
3. **Standard interfaces** - Implement `String()`, `MarshalText()`, etc.
4. **Composable** - Types should chain together

---

## Questions for User

1. **Priority Check:** Should we focus on bug fixes first, or new types?
2. **v3 Question:** Is there interest in starting v3 design, or is v2 sufficient?
3. **Library Choice:** Should we add `charmbracelet/bubbles` for progress spinners?
4. **Config Migration:** Should we prioritize replacing internal/config with koanf?
5. **Testing Strategy:** Want to add benchmarks next, or more integration tests?

---

## Summary

**Immediate Wins:**

- ✅ Package() function fixed and tested
- ✅ WithCLIScope now works correctly
- ✅ Test coverage maintained at 90.2%

**Next Actions:**

1. Fix CLI[T] AddCommand flag parsing bug
2. Add t.Parallel() to tests
3. Add custom types (URL, Email, Port)
4. Consider Result[T] type

**Technical Debt:**

- Some test files are large (need splitting)
- A few files exceed 350 lines
- internal/config could be replaced with koanf

---

_Last Updated: 2026-04-01_
