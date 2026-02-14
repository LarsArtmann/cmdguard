# cmdguard Architecture Review

> **Review Date:** 2026-02-14  
> **Reviewer:** Crush (AI Assistant)  
> **Library:** github.com/larsartmann/cmdguard  
> **Version:** 0.1.0 (pre-release)

---

## Executive Summary

**cmdguard** is a Go framework for building validated Cobra CLI applications with integrated dependency injection, lifecycle management, and runtime validation. It combines four powerful libraries:

| Library | Purpose |
|---------|---------|
| `spf13/cobra` | CLI command framework |
| `samber/do/v2` | Dependency injection container |
| `knadh/koanf` | Hierarchical configuration management |
| `charmbracelet/fang` | Styled error output and execution |

### Overall Assessment

| Criterion | Score | Notes |
|-----------|-------|-------|
| **Architecture** | 8/10 | Clean separation, good patterns |
| **Code Quality** | 7/10 | Minor issues, one broken API |
| **Test Coverage** | 5/10 | Core packages untested |
| **Documentation** | 9/10 | Excellent principles doc |
| **Production Readiness** | 6/10 | Needs fixes before v1.0 |

---

## Module Structure

```
cmdguard/
├── cmd/cmdguard/           # CLI application entry point
│   └── main.go             # Setup and execution
├── pkg/cmdguard/           # Public API
│   └── public_api.go       # Application facade
└── internal/
    ├── commands/           # Cobra command registry
    │   └── root.go         # Registry + command setup
    ├── config/             # Koanf configuration
    │   ├── provider.go     # Config loading
    │   └── provider_test.go
    ├── di/                 # Dependency injection
    │   └── module.go       # Service wiring
    ├── logging/            # Structured logging
    │   └── logger.go       # Slog setup
    └── validation/         # Command/flag validation
        ├── registry.go     # Metadata storage
        ├── registry_test.go
        ├── validator.go    # Validation logic
        └── validator_test.go
```

---

## Core Components

### 1. Application Facade (`pkg/cmdguard`)

The `Application` struct is the primary entry point:

```go
type Application struct {
    module      *di.Module
    registry    *commands.Registry
    validator   *validation.Validator
    config      *config.Config
    rootCmd     *cobra.Command
    initialized bool
}
```

**Lifecycle Flow:**

```
New() → Initialize() → Validate() → Execute() → Shutdown()
   ↓         ↓              ↓           ↓           ↓
  nil    ProvideSVcs   ValidateAll   FangExec   Cleanup
```

**Public API Surface:**

| Method | Purpose | Status |
|--------|---------|--------|
| `New()` | Constructor | ✅ |
| `Initialize()` | Setup DI container | ✅ |
| `InitializeWithOptions()` | Setup with options | ⚠️ (1 broken) |
| `Validate()` / `MustValidate()` | Run validation | ✅ |
| `Execute()` / `ExecuteAndExit()` | Run command | ✅ |
| `Shutdown()` | Graceful cleanup | ✅ |
| `HealthCheck()` | Service health | ✅ |
| `AddCommand()` | Runtime commands | ✅ |

---

### 2. DI Module (`internal/di`)

Uses `samber/do/v2` for service management.

**Service Registration:**

```go
// Lazy services - singleton, created on first use
do.Provide(injector, config.NewConfig)
do.Provide(injector, validation.NewRegistry)
do.Provide(injector, validation.NewValidator)
do.Provide(injector, commands.NewRegistry)

// Transient - new instance per injection
do.ProvideTransient(injector, validation.NewFlagValidator)
```

**Service Dependencies:**

```
Config ←──────┬─── Registry
              │
Validator ←───┴─── Registry
       │
       └── Config
```

---

### 3. Command Registry (`internal/commands`)

Manages the Cobra command tree with integrated validation.

**Registry Responsibilities:**
- Root command creation
- Subcommand registration
- Global flag setup (`--config`, `--log-level`, `--strict`)
- PreRun validation (log-level enum check)
- Fang-styled execution

**Flag Design (Good Example):**

```go
root.PersistentFlags().StringP("config", "c", "", "Config file path")
root.PersistentFlags().StringP("log-level", "l", "info", "Log level: debug, info, warn, error")
root.PersistentFlags().BoolP("strict", "s", false, "Enable strict mode validation")
```

✅ Boolean flags are boolean  
✅ Short flags provided  
✅ Clear defaults  
✅ Enum validation in PreRunE

---

### 4. Configuration (`internal/config`)

**Loading Hierarchy** (highest priority last):

1. Default values (in code)
2. Config file (`config.yaml`)
3. Environment variables (`CMDGUARD_*`)
4. Custom config file via `CMDGUARD_CONFIG_FILE`

**Config Structure:**

```go
type Config struct {
    StrictMode bool   `koanf:"strict_mode"`
    ConfigFile string `koanf:"config_file"`
    LogLevel   string `koanf:"log_level"`
}
```

⚠️ **Issue:** `NewConfigWithCommand()` is never used - flag integration incomplete.

---

### 5. Validation System (`internal/validation`)

Two-tier validation architecture:

#### Registry (`validation.Registry`)
- Thread-safe command/flag metadata storage
- `sync.RWMutex` for concurrent access
- Extracts flag info from Cobra commands

#### Validator (`validation.Validator`)
- Validates all registered commands have handlers
- Validates all flags are properly bound
- Recursively validates command trees

**Validation Rules:**

| Check | Description |
|-------|-------------|
| Handler presence | Leaf commands must have Run/RunE |
| Flag binding | All declared flags must be bound |
| Subcommand handling | Parent commands can omit handlers |

---

## Critical Issues

### 🔴 P0: Broken `WithValidationHook` Option

**Location:** `pkg/cmdguard/public_api.go:130-135`

```go
func WithValidationHook(hook func() error) Option {
    return func(a *Application) error {
        // Store hook for later execution  // <-- NEVER STORED!
        return nil
    }
}
```

**Problem:** Hook is accepted but never stored or executed.

**Fix:**
```go
type Application struct {
    // ... existing fields ...
    validationHooks []func() error  // Add this
}

func WithValidationHook(hook func() error) Option {
    return func(a *Application) error {
        a.validationHooks = append(a.validationHooks, hook)
        return nil
    }
}

func (a *Application) Validate() error {
    // ... existing validation ...
    for _, hook := range a.validationHooks {
        if err := hook(); err != nil {
            return fmt.Errorf("validation hook failed: %w", err)
        }
    }
    return nil
}
```

---

### 🟡 P1: Incomplete Context Handling

**Location:** `internal/di/module.go:123-127`

```go
func (m *Module) ShutdownWithContext(ctx context.Context) error {
    // For now, delegate to regular shutdown
    // Context handling would be added in production
    return m.Shutdown()
}
```

**Problem:** Context timeout/deadline is completely ignored.

**Fix:** Implement proper context-aware shutdown with timeout channels.

---

### 🟡 P1: Silent Config Errors

**Location:** `internal/config/provider.go:32-35`

```go
_ = k.Load(file.Provider("config.yaml"), yaml.Parser())
_ = k.Load(env.Provider("CMDGUARD_", ".", nil), nil)
```

**Problem:** Errors silently ignored. Users won't know if config file is missing.

**Fix:** Add structured logging for config load attempts.

---

### 🟢 P2: Error Aggregation

**Location:** `internal/di/module.go:115-117`

```go
if len(errs) > 0 {
    return errs[0] // Only returns first error
}
```

**Problem:** Other shutdown errors are lost.

**Fix:** Use `errors.Join(errs...)` (Go 1.20+)

---

## Test Coverage Analysis

| Package | Lines | Tests | Coverage | Status |
|---------|-------|-------|----------|--------|
| `internal/config` | 130 | 3 funcs, 15 cases | ~80% | ✅ Good |
| `internal/validation` | 274 | 2 files, 15+ cases | ~75% | ✅ Good |
| `internal/commands` | 168 | 0 | 0% | 🔴 Missing |
| `internal/di` | 142 | 0 | 0% | 🔴 Missing |
| `pkg/cmdguard` | 213 | 0 | 0% | 🔴 Missing |

**Recommendation:** Add integration tests for full lifecycle.

---

## Design Patterns Assessment

### ✅ Well Implemented

| Pattern | Implementation | Quality |
|---------|---------------|---------|
| Dependency Injection | `samber/do/v2` | Excellent |
| Registry | `validation.Registry` | Good |
| Builder | `Application` + `Option` | Good |
| Facade | `pkg/cmdguard` public API | Good |

### ⚠️ Could Improve

| Pattern | Current | Suggested |
|---------|---------|-----------|
| Observer | Hooks broken | Fix or remove |
| Strategy | Validation hardcoded | Allow custom validators |
| Factory | Direct construction | Factory methods for testability |

---

## API Design Review

### Public API Surface

```go
// Constructor
func New() *Application

// Initialization
func (a *Application) Initialize() error
func (a *Application) InitializeWithOptions(opts ...Option) error

// Options
type Option func(*Application) error
func WithCommand(cmd *cobra.Command) Option                    // ✅
func WithValidationHook(hook func() error) Option               // ❌ BROKEN

// Validation
func (a *Application) Validate() error
func (a *Application) MustValidate()

// Execution
func (a *Application) Execute(ctx context.Context) error
func (a *Application) ExecuteAndExit(ctx context.Context)

// Lifecycle
func (a *Application) Shutdown() error
func (a *Application) HealthCheck() error

// Accessors
func (a *Application) Root() *cobra.Command
func (a *Application) Registry() *commands.Registry
func (a *Application) Config() *config.Config
func (a *Application) Validator() *validation.Validator
func (a *Application) Injector() do.Injector
func (a *Application) IsStrictMode() bool

// Runtime modification
func (a *Application) AddCommand(cmd *cobra.Command)
```

### Missing API (Suggestions)

```go
// Additional options that would be useful:
func WithConfigFile(path string) Option
func WithLogger(logger *slog.Logger) Option
func WithShutdownTimeout(timeout time.Duration) Option
func WithPreRunHook(hook func() error) Option
func WithPostRunHook(hook func() error) Option
func WithPanicRecovery(enabled bool) Option

// Validation customization:
func (a *Application) AddValidator(v CommandValidator)
```

---

## Performance Considerations

### Current State

| Aspect | Assessment |
|--------|------------|
| Memory | Lightweight, minimal allocations |
| Startup | DI resolution on first use (lazy) |
| Concurrency | Thread-safe via mutexes |
| Validation | O(n) where n = commands + flags |

### Potential Optimizations

1. **Validation caching** - Skip re-validation if tree unchanged
2. **Lazy config loading** - Only load config when accessed
3. **Parallel validation** - For large command trees

---

## Security Assessment

| Check | Status | Notes |
|-------|--------|-------|
| Input validation | ✅ | Log level enum validation |
| Config file path | ⚠️ | No path traversal check |
| Error messages | ✅ | No sensitive data exposed |
| Dependency versions | ✅ | No known vulnerabilities |

---

## Production Readiness Checklist

### Must Have (Blocking)

- [ ] Fix `WithValidationHook` implementation
- [ ] Add tests for `internal/commands`
- [ ] Add tests for `internal/di`
- [ ] Add tests for `pkg/cmdguard`
- [ ] Implement `ShutdownWithContext` properly
- [ ] Add config file path validation

### Should Have (Recommended)

- [ ] Structured logging for config loading
- [ ] Use `errors.Join()` for multiple errors
- [ ] Integration tests for full lifecycle
- [ ] Benchmark tests for validation
- [ ] Example applications

### Nice to Have

- [ ] Plugin system for custom validators
- [ ] Metrics/monitoring hooks
- [ ] Config hot-reload
- [ ] Shell completion generation

---

## Recommendations

### Short Term (This Week)

1. **Fix the broken hook option** - Either implement or remove
2. **Add missing tests** - Focus on `commands` and `di` packages
3. **Complete context handling** - Respect timeout in shutdown

### Medium Term (This Month)

1. **Integration tests** - Full app lifecycle testing
2. **Documentation** - Add usage examples and patterns
3. **Error handling** - Use `errors.Join()` consistently

### Long Term (Next Quarter)

1. **Plugin architecture** - Allow custom validators
2. **Observability** - Metrics and tracing hooks
3. **v1.0 release** - Stable API guarantee

---

## Conclusion

cmdguard is a **promising framework** with solid architectural foundations. It demonstrates:

- ✅ Good separation of concerns
- ✅ Proper use of established libraries
- ✅ Thoughtful CLI design principles
- ✅ Clean public API

However, it needs:

- 🔴 Critical bug fixes (broken hook option)
- 🔴 Comprehensive test coverage
- 🟡 Minor API completeness improvements

**Verdict:** 6/10 production readiness. With recommended fixes, easily 8/10.

**Recommendation:** Fix P0/P1 issues and add tests before any production use.

---

## Appendix: Code Metrics

```
Total Go Files:     14
Total Lines:        ~1,500
Test Files:         3
Test Lines:         ~500
Dependencies:       6 direct, 27 indirect
Build Time:         <1s
Test Time:          <1s
Binary Size:        ~15MB (with dependencies)
```

---

*Review generated by Crush AI Assistant*  
*Assisted-by: Crush via crush@charm.land*
