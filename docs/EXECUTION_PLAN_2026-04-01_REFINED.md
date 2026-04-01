# Refined Execution Plan - 2026-04-01

## Self-Reflection: What I Could Have Done Better

### 1. What I Forgot

- **Didn't add t.Parallel()** to the new tests in `types_custom_test.go`
- **Didn't check if the custom types integrate** with FlagRegistry default value handling
- **Didn't add example usage** in example_test.go or examples/ directory
- **Didn't update FEATURES.md** to document the new types
- **Didn't check for code duplication** - URL, Email, Port parsing patterns are similar
- **Didn't split large files** - types_custom.go (591 lines) and types_custom_test.go (621 lines) exceed 350 line limit

### 2. What Could Be Improved

- **Type architecture**: Custom types share similar patterns (Parse/MustParse, validation, marshaling) - could be abstracted
- **Error messages**: Could be more consistent and helpful
- **Documentation**: Missing godoc examples for new types
- **Integration tests**: No tests showing custom types used as CLI flags
- **Benchmarks**: No performance measurements for parsing operations

### 3. Architecture Improvements

- **Generic Validated[T] wrapper**: Could encapsulate validation logic
- **Type-safe ID types**: Like UserID, OrderID using branded types pattern
- **Result[T] type**: For explicit error handling without exceptions

### 4. Established Libraries to Leverage

- **github.com/google/uuid**: For UUID type (instead of custom string)
- **github.com/oklog/ulid**: For ULID type
- **github.com/shopspring/decimal**: For Money/Decimal type
- **github.com/go-playground/validator/v10**: For struct validation tags

---

## Comprehensive Multi-Step Execution Plan

### Phase 1: Immediate Fixes (Impact: High, Effort: Low) - TODAY

| Step | Task                                               | Impact | Effort | Time   |
| ---- | -------------------------------------------------- | ------ | ------ | ------ |
| 1.1  | Add t.Parallel() to types_custom_test.go           | High   | Low    | 10 min |
| 1.2  | Add godoc examples for custom types                | High   | Low    | 15 min |
| 1.3  | Update FEATURES.md with new types                  | High   | Low    | 10 min |
| 1.4  | Fix missing default value handling in custom types | High   | Low    | 15 min |
| 1.5  | Split types_custom.go into separate files          | Medium | Low    | 20 min |

### Phase 2: Code Quality (Impact: Medium, Effort: Medium) - THIS WEEK

| Step | Task                                          | Impact | Effort | Time   |
| ---- | --------------------------------------------- | ------ | ------ | ------ |
| 2.1  | Create generic Validated[T] wrapper type      | High   | Medium | 30 min |
| 2.2  | Refactor custom types to use Validated[T]     | High   | Medium | 45 min |
| 2.3  | Add integration tests (CLI with custom types) | High   | Medium | 30 min |
| 2.4  | Create Result[T] type for error handling      | Medium | Medium | 25 min |
| 2.5  | Add example programs using custom types       | Medium | Medium | 30 min |

### Phase 3: Infrastructure (Impact: High, Effort: Medium) - THIS WEEK

| Step | Task                                     | Impact | Effort | Time   |
| ---- | ---------------------------------------- | ------ | ------ | ------ |
| 3.1  | Create GitHub Actions CI workflow        | High   | Medium | 30 min |
| 3.2  | Add t.Parallel() to all test files       | High   | Low    | 45 min |
| 3.3  | Add benchmarks for parsing operations    | Medium | Medium | 30 min |
| 3.4  | Migrate testify tests to stdlib          | Medium | High   | 60 min |
| 3.5  | Fix file size limits (split large files) | Medium | Medium | 40 min |

### Phase 4: Advanced Features (Impact: Medium, Effort: High) - NEXT WEEK

| Step | Task                                    | Impact | Effort | Time   |
| ---- | --------------------------------------- | ------ | ------ | ------ |
| 4.1  | Add UUID type (using google/uuid)       | Medium | Low    | 20 min |
| 4.2  | Add Money/Decimal type                  | Medium | Medium | 30 min |
| 4.3  | Create branded ID types (UserID, etc.)  | Medium | Medium | 30 min |
| 4.4  | Add Config file auto-loading with koanf | High   | High   | 60 min |
| 4.5  | Add env struct tag support              | Medium | Medium | 40 min |

---

## Detailed Implementation Plan

### Step 1.1: Add t.Parallel() to types_custom_test.go

**Why**: All table-driven tests should run in parallel for faster execution.

**How**:

```go
func TestURL(t *testing.T) {
    t.Parallel()  // ADD THIS
    t.Run("ParseURL valid", func(t *testing.T) {
        t.Parallel()  // ADD THIS TO ALL SUBTESTS
        // ...
    })
}
```

### Step 1.2: Add godoc examples

**Why**: Examples show up in godoc and help users understand the API.

**How**: Add to `example_test.go`:

```go
func ExampleParseURL() {
    u, err := v2.ParseURL("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(u.Hostname())
    // Output: example.com
}
```

### Step 1.4: Fix default value handling

**Why**: Custom types need to handle default values from struct tags.

**Check**: Does `flag:"endpoint" default:"https://api.example.com"` work?

**Issue**: The current `parseAndSetCustom` doesn't handle default values - it only parses user input.

**Fix needed in flags.go**: When registering flags, default values for custom types need special handling.

### Step 1.5: Split types_custom.go

**Files to create**:

- `types_url.go` - URL type (120 lines)
- `types_email.go` - Email type (80 lines)
- `types_port.go` - Port type (100 lines)
- `types_filepath.go` - FilePath type (120 lines)
- `types_hostport.go` - HostPort type (80 lines)

Delete `types_custom.go` after splitting.

### Step 2.1: Create Validated[T] wrapper

**Design**:

```go
type Validated[T any] struct {
    value T
    valid bool
    validate func(T) error
}

func NewValidated[T any](value T, validate func(T) error) (Validated[T], error)
func (v Validated[T]) Value() (T, bool)
func (v Validated[T]) Must() T
```

**Benefits**:

- Consistent validation pattern
- Composable validators
- Type-safe

### Step 2.4: Create Result[T] type

**Design** (similar to Rust's Result):

```go
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](value T) Result[T]
func Err[T any](err error) Result[T]
func (r Result[T]) IsOk() bool
func (r Result[T]) IsErr() bool
func (r Result[T]) Unwrap() T
func (r Result[T]) UnwrapOr(default T) T
func (r Result[T]) Map(f func(T) T) Result[T]
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T]
```

### Step 3.1: GitHub Actions CI

**Workflow features**:

- Run on PRs and pushes to main
- Go 1.26
- Run tests with race detector
- Run tests with coverage
- Run linter (golangci-lint)
- Check go mod tidy
- Upload coverage to codecov

### Step 3.2: Add t.Parallel() pattern

**Pattern for each test file**:

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    t.Run("case", func(t *testing.T) {
        t.Parallel()
        // test code
    })
}
```

---

## Code Review: What Exists That We Can Reuse

### Existing Patterns to Leverage

1. **Duration type** (`types.go:53-96`)
   - Pattern: `ParseX(s string) (X, error)`
   - Pattern: `FromX(x) X` constructor
   - Pattern: Methods for common operations
   - Pattern: MarshalText/UnmarshalText

2. **Enum type** (`types.go:12-48`)
   - Pattern: `allowed` slice for validation
   - Pattern: `IsEmpty()` check

3. **Option[T]** (`types.go:267-459`)
   - Pattern: Generic wrapper type
   - Pattern: Map, Filter, And, Or operations
   - Pattern: JSON/Text marshaling

4. **FlagRegistry** (`flags.go`)
   - Pattern: `parseAndSetCustom` for type dispatch
   - Pattern: Type-based switch using `reflect.TypeFor[X]()`

5. **Error types** (`errors.go`)
   - Pattern: Sentinel errors
   - Pattern: Wrapped errors with context

### Existing Libraries in go.mod

- `github.com/samber/do/v2` - DI (already used)
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/fang` - Styling
- `github.com/knadh/koanf/v2` - Config (could use more)

### Libraries to Add

- `github.com/google/uuid` - UUID validation
- `github.com/shopspring/decimal` - Money/Decimal

---

## Priority Queue (Sorted by Impact/Effort Ratio)

### P0 - Critical (Do First)

1. ✅ Add t.Parallel() to new tests
2. ✅ Fix default value handling for custom types
3. ✅ Split large files (line limit)

### P1 - High Value

4. Add godoc examples
5. Update FEATURES.md
6. Create GitHub Actions CI
7. Add t.Parallel() to all tests

### P2 - Medium Value

8. Create Validated[T] wrapper
9. Refactor types to use Validated[T]
10. Add benchmarks
11. Create Result[T] type

### P3 - Future Work

12. Add UUID type
13. Add Money type
14. Add branded ID types
15. Config file auto-loading

---

## Questions for Future Consideration

1. **Should we use generics more?**
   - Validated[T] could replace individual validation logic
   - Parser[T] interface for all parseable types

2. **Should we add more charmbracelet integration?**
   - bubbles for spinners/progress bars
   - lipgloss for styling

3. **Should we support more config formats?**
   - TOML (using github.com/BurntSushi/toml)
   - YAML (already supported via koanf)

4. **Should we add validation tags?**
   - `validate:"email"` instead of Email type
   - `validate:"url"` instead of URL type
   - More flexible but less type-safe

---

## Progress Tracking

| Phase   | Status         | Items Complete | Items Total |
| ------- | -------------- | -------------- | ----------- |
| Phase 1 | 🔄 In Progress | 0              | 5           |
| Phase 2 | ⏳ Not Started | 0              | 5           |
| Phase 3 | ⏳ Not Started | 0              | 5           |
| Phase 4 | ⏳ Not Started | 0              | 5           |

**Last Updated**: 2026-04-01
