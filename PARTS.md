# cmdguard Component Analysis

**Last Updated:** 2026-02-28
**Version:** 2.0.0
**Status:** Analysis Complete

---

## Executive Summary

cmdguard contains several components with extraction potential. This document analyzes each component, compares against alternatives, and identifies unique value propositions.

**Key Findings:**

| Component             | Extraction Potential | Unique Value                               | Recommendation                  |
| --------------------- | -------------------- | ------------------------------------------ | ------------------------------- |
| Type-Safe Flag System | **High**             | Struct tags + generics + Cobra integration | Extract as standalone lib       |
| DI Scope Wrapper      | **Medium**           | Simplified samber/do/v2 API                | Keep internal, document pattern |
| Error Types           | **Low**              | CLI-specific sentinel errors               | Keep internal                   |
| Config Provider       | **Low**              | Simple env var loading                     | Use koanf instead               |
| Logging Setup         | **Low**              | Basic slog wrapper                         | Use charmbracelet/log directly  |

---

## Component Analysis

### 1. Type-Safe Flag System (`pkg/cmdguard/v2/flags*.go`)

**Lines of Code:** ~400
**Dependencies:** cobra, pflag, reflect
**Complexity:** Medium

#### What It Does

```go
type GreetFlags struct {
    Name   string `flag:"name" short:"n" default:"World" help:"Name to greet"`
    Count  int    `flag:"count" short:"c" default:"1" help:"Number of times"`
    Shout  bool   `flag:"shout" short:"s" default:"false" help:"Uppercase output"`
}

// Automatic registration, parsing, and validation
registry := NewFlagRegistry(&GreetFlags{})
registry.RegisterFlags(cmd)
registry.ValidateFlags(cmd)
```

#### Unique Features

| Feature             | cmdguard             | Kong             | go-flags   | pflag      |
| ------------------- | -------------------- | ---------------- | ---------- | ---------- |
| Struct tags         | ✅                   | ✅               | ✅         | ❌         |
| Cobra integration   | ✅ Native            | Manual           | Manual     | ✅ Native  |
| Type inference      | ✅ Automatic         | Mapper interface | Reflection | Manual     |
| Required validation | ✅ `required:"true"` | ✅               | ✅         | ❌         |
| Typo suggestions    | ✅ Levenshtein       | ❌               | ❌         | ❌         |
| Generic type safety | ✅ Compile-time      | ❌ Runtime       | ❌ Runtime | ❌ Runtime |

#### Alternatives

**Kong** (`github.com/alecthomas/kong`) - ~2.4k stars

- Declarative struct tags, type-safe mappers
- No Cobra integration (separate CLI framework)
- No typo suggestions

**go-flags** (`github.com/jessevdk/go-flags`) - ~2.6k stars

- Extensive struct tag support
- No Cobra integration
- Legacy API design

**pflag** (stdlib for Cobra) - Industry standard

- No struct tag support
- Manual flag registration

#### Extraction Recommendation: **HIGH**

**Proposed Library:** `github.com/larsartmann/flagtags`

```go
// Standalone API
import "github.com/larsartmann/flagtags"

type Config struct {
    Port int `flag:"port" short:"p" default:"8080" help:"Server port"`
}

flagtags.Register(cmd, &Config{})
flagtags.Parse(cmd, &Config{})
flagtags.Validate(cmd, &Config{})
```

**Value Proposition:**

1. **Only library with native Cobra/pflag integration + struct tags**
2. **Typo suggestions** - unique among flag libraries
3. **Generic compile-time safety** - unlike reflection-based alternatives
4. **Zero-boilerplate** - one function call for registration

---

### 2. DI Scope Wrapper (`pkg/cmdguard/v2/scope.go`)

**Lines of Code:** ~270
**Dependencies:** samber/do/v2
**Complexity:** Low-Medium

#### What It Does

Wraps samber/do/v2 with CLI-friendly patterns:

```go
// Scoped DI with hierarchy
scope := NewScope("myapp")
child := scope.Child("worker")

// Generic providers
Provide(scope, func(i do.Injector) (*Database, error) { ... })
ProvideValue(scope, &Logger{})

// Generic invocation
db, err := Invoke[*Database](scope)

// Lifecycle
scope.HealthCheck()
scope.Shutdown(ctx)
```

#### Unique Features

| Feature           | cmdguard Scope | samber/do/v2 | fx  | wire   |
| ----------------- | -------------- | ------------ | --- | ------ |
| Scope hierarchy   | ✅             | ✅ Native    | ❌  | ❌     |
| Generic API       | ✅             | ✅           | ❌  | ✅     |
| CLI-focused       | ✅             | ❌           | ❌  | ❌     |
| Health checks     | ✅             | ✅           | ✅  | ❌     |
| Graceful shutdown | ✅             | ✅           | ✅  | Manual |

#### Alternatives

**samber/do/v2** - Direct use

- Same functionality, more verbose
- No scope path tracking
- No CLI convenience methods

**uber-go/fx** - ~5.8k stars

- Lifecycle hooks (OnStart/OnStop) - **gap in cmdguard**
- No hierarchical scopes
- Requires module pattern

**google/wire** - ~13k stars

- Compile-time safety
- No runtime scopes
- Requires code generation

#### Extraction Recommendation: **MEDIUM**

**Keep internal** but document as reusable pattern.

**Missing features to add:**

- `OnStart` / `OnStop` lifecycle hooks (fx pattern)
- Per-command scope factory
- Ordered shutdown priority

---

### 3. Error Types (`pkg/cmdguard/v2/errors.go`)

**Lines of Code:** ~210
**Dependencies:** None (stdlib)
**Complexity:** Low

#### What It Does

CLI-specific sentinel errors with context:

```go
// Sentinel errors
ErrInvalidCommand, ErrMissingHandler, ErrFlagParseFailed

// Contextual wrappers
NewCommandError(name, err)      // "command 'x': ..."
NewFlagError(name, err)         // "flag 'x': ..."
NewFlagErrorWithSuggestion(...) // "flag 'x': ... (did you mean --y?)"
NewConfigError(field, err)      // "config field 'x': ..."
NewServiceError(type, err)      // "service 'x': ..."
```

#### Unique Features

- **Typo suggestions in FlagError** - unique
- **CLI-specific error types** - command, flag, config, service
- **errors.Is/As compatible** - proper Go error wrapping

#### Alternatives

**cockroachdb/errors** - General purpose

- No CLI-specific types
- More features than needed

**pkg/errors** - Deprecated

#### Extraction Recommendation: **LOW**

Too small to justify separate library. Keep internal.

---

### 4. Config Provider (`internal/config/provider.go`)

**Lines of Code:** ~80
**Dependencies:** None (stdlib)
**Complexity:** Low

#### What It Does

Simple environment variable loading:

```go
cfg := config.Load()
// Reads CMDGUARD_LOG_LEVEL, CMDGUARD_STRICT_MODE, etc.
```

#### Alternatives

**koanf** (`github.com/knadh/koanf/v2`) - Recommended by HOW_TO_GOLANG.md

- Multiple format support (YAML, JSON, TOML, ENV)
- Hot reload capable
- No global state
- Rich provider ecosystem

#### Extraction Recommendation: **LOW**

**Replace with koanf** per library policy. Current implementation is too simple.

---

### 5. Logging Setup (`internal/logging/logger.go`)

**Lines of Code:** ~130
**Dependencies:** slog (stdlib)
**Complexity:** Low

#### What It Does

Basic slog wrapper with format/level parsing:

```go
logger := logging.NewLogger("text", "debug")
// or "json", "info", etc.
```

#### Alternatives

**charmbracelet/log** - Recommended by HOW_TO_GOLANG.md

- Full slog handler implementation
- Styled output
- Context integration
- Already a dependency via fang

#### Extraction Recommendation: **LOW**

**Replace with charmbracelet/log directly** per library policy.

---

## Proposed Extractions

### High Priority: `flagtags`

**Repository:** `github.com/larsartmann/flagtags`

**Scope:**

- `pkg/cmdguard/v2/flags.go`
- `pkg/cmdguard/v2/flags_parse.go`
- `pkg/cmdguard/v2/flags_suggest.go`
- `pkg/cmdguard/v2/flags_registry_test.go`

**API Design:**

```go
package flagtags

import "github.com/spf13/cobra"

// Register adds flags to command from struct tags
func Register(cmd *cobra.Command, cfg any) error

// Parse populates struct from command flags
func Parse(cmd *cobra.Command, cfg any) error

// Validate checks required flags and enum values
func Validate(cmd *cobra.Command, cfg any) error

// Suggest returns closest flag name for typos
func Suggest(available []string, input string) string

// FlagTag represents parsed struct tag
type FlagTag struct {
    Name     string
    Short    string
    Default  string
    Help     string
    Required bool
    Values   []string  // For enums
    Type     reflect.Type
}
```

**Unique Value:**

1. Only library combining struct tags + Cobra integration + typo suggestions
2. Generic compile-time safety
3. Zero boilerplate

---

## Keep Internal

### DI Scope (`scope.go`)

Keep internal but document pattern. Add missing features:

```go
// Add lifecycle hooks (fx pattern)
type LifecycleHook struct {
    OnStart func(ctx context.Context) error
    OnStop  func(ctx context.Context) error
}

func (s *Scope) RegisterLifecycle(hook LifecycleHook) error

// Add per-command scope factory
func (g *GuardedCommand[T, F]) CommandScope(cmdName string) *Scope
```

### Error Types (`errors.go`)

Keep internal. Too small for extraction.

### Config Provider (`config/provider.go`)

**Replace with koanf** - align with library policy.

### Logging (`logging/logger.go`)

**Replace with charmbracelet/log** - align with library policy.

---

## Alternative Libraries Summary

### CLI Frameworks with Type-Safe Flags

| Library  | Stars     | Struct Tags | DI Integration | Fail-Fast     |
| -------- | --------- | ----------- | -------------- | ------------- |
| Kong     | ~2.4k     | ✅          | Manual         | Parse-time    |
| go-flags | ~2.6k     | ✅          | None           | Parse-time    |
| Cobra+Fx | ~38k+5.8k | ❌          | ✅ Manual      | ValidateApp() |
| Wire+CLI | ~13k      | Varies      | ✅ Generated   | Compile-time  |

### DI Libraries with Scope Support

| Library      | Stars | Scope Hierarchy | Generic API | Lifecycle         |
| ------------ | ----- | --------------- | ----------- | ----------------- |
| samber/do/v2 | ~2k   | ✅              | ✅          | ✅                |
| uber-go/fx   | ~5.8k | ❌              | ❌          | ✅ OnStart/OnStop |
| google/wire  | ~13k  | ❌              | ✅          | Manual            |

---

## Action Items

### Phase 1: Document Patterns (Low Effort)

- [ ] Document DI scope pattern in docs/
- [ ] Add lifecycle hook examples

### Phase 2: Align with Library Policy (Medium Effort)

- [ ] Replace `internal/config` with koanf
- [ ] Replace `internal/logging` with charmbracelet/log

### Phase 3: Extract flagtags (High Effort, High Value)

- [ ] Create `github.com/larsartmann/flagtags` repository
- [ ] Extract flag-related code
- [ ] Add standalone tests
- [ ] Document API
- [ ] Update cmdguard to use flagtags

---

## References

- [HOW_TO_GOLANG.md](/Users/larsartmann/projects/library-policy/HOW_TO_GOLANG.md) - Library policy
- [FEATURES.md](FEATURES.md) - cmdguard feature status
- [README.md](README.md) - cmdguard documentation
