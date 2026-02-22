# cmdguard Features

**Last Updated:** 2026-02-14  
**Version:** 2.0.0  
**Go Version:** 1.26.0

---

## Legend

| Status                  | Meaning                                                  |
| ----------------------- | -------------------------------------------------------- |
| ✅ FULLY_FUNCTIONAL     | Feature works as designed, tested, and documented        |
| ⚠️ PARTIALLY_FUNCTIONAL | Feature works but has limitations, gaps, or known issues |
| 🔧 BROKEN               | Feature does not work or is non-functional               |
| 📝 PLANNED              | Feature is designed but not yet implemented              |
| 🗑️ DEPRECATED           | Feature exists but is scheduled for removal              |
| ❓ UNKNOWN              | Status cannot be determined                              |

---

## v2 API (pkg/cmdguard/v2)

The v2 API is a complete rewrite offering type-safe CLI construction with dependency injection.

### GuardedCommand[T]

Type-safe CLI with generic configuration type.

| Feature                                       | Status              | Notes                                |
| --------------------------------------------- | ------------------- | ------------------------------------ |
| `New[T](name, short, defaults)`               | ✅ FULLY_FUNCTIONAL | Creates typed CLI, never panics      |
| `NewWithLong[T](name, short, long, defaults)` | ✅ FULLY_FUNCTIONAL | Creates CLI with long description    |
| `AddCommand(cmd)`                             | ✅ FULLY_FUNCTIONAL | Adds typed subcommand, returns error |
| `AddCommandFunc(fn)`                          | ✅ FULLY_FUNCTIONAL | Lazy command initialization          |
| `Execute(ctx)`                                | ✅ FULLY_FUNCTIONAL | Runs command with context            |
| `ExecuteWithArgs(ctx, args)`                  | ✅ FULLY_FUNCTIONAL | For testing                          |
| `ExecuteAndExit(ctx)`                         | ✅ FULLY_FUNCTIONAL | Runs and calls os.Exit               |
| `Scope()`                                     | ✅ FULLY_FUNCTIONAL | Returns DI injector                  |
| `ScopeStruct()`                               | ✅ FULLY_FUNCTIONAL | Returns Scope for advanced ops       |
| `Config()`                                    | ✅ FULLY_FUNCTIONAL | Returns typed config \*T             |
| `SetConfig(cfg)`                              | ✅ FULLY_FUNCTIONAL | Updates configuration                |
| `Shutdown(ctx)`                               | ✅ FULLY_FUNCTIONAL | Graceful shutdown                    |
| `HealthCheck()`                               | ✅ FULLY_FUNCTIONAL | Runs health checks                   |
| `RootCommand()`                               | ✅ FULLY_FUNCTIONAL | Returns underlying cobra.Command     |

### Command[T]

Type-safe command definition with typed flags.

| Feature                           | Status              | Notes                     |
| --------------------------------- | ------------------- | ------------------------- |
| `Use`, `Short`, `Long` fields     | ✅ FULLY_FUNCTIONAL | Standard command metadata |
| `Flags any`                       | ✅ FULLY_FUNCTIONAL | Struct with flag tags     |
| `RunE func(ctx, *T, flags)`       | ✅ FULLY_FUNCTIONAL | Type-safe handler         |
| `PreRunE` / `PostRunE`            | ✅ FULLY_FUNCTIONAL | Lifecycle hooks           |
| `Commands []Command[T]`           | ✅ FULLY_FUNCTIONAL | Nested subcommands        |
| `Hidden`, `Deprecated`            | ✅ FULLY_FUNCTIONAL | Visibility options        |
| `Aliases`, `Version`              | ✅ FULLY_FUNCTIONAL | Additional metadata       |
| `Validate()`                      | ✅ FULLY_FUNCTIONAL | Command validation        |
| `NewCommand(use, short, opts...)` | ✅ FULLY_FUNCTIONAL | Constructor with options  |

### Flag System

| Feature                         | Status              | Notes                                             |
| ------------------------------- | ------------------- | ------------------------------------------------- |
| Struct tag flags                | ✅ FULLY_FUNCTIONAL | `flag:"name" short:"n" default:"val" help:"desc"` |
| Type inference                  | ✅ FULLY_FUNCTIONAL | string, int, bool, float64 supported              |
| Short flags                     | ✅ FULLY_FUNCTIONAL | `short:"n"` for `-n`                              |
| Default values                  | ✅ FULLY_FUNCTIONAL | `default:"value"` tag                             |
| Help text                       | ✅ FULLY_FUNCTIONAL | `help:"description"` tag                          |
| Required flags                  | ✅ FULLY_FUNCTIONAL | `required:"true"` tag, validated at runtime       |
| Flag typo suggestions           | ✅ FULLY_FUNCTIONAL | Levenshtein distance-based suggestions            |
| `SuggestFlag(available, input)` | ✅ FULLY_FUNCTIONAL | Returns closest match for typos                   |
| FlagRegistry                    | ✅ FULLY_FUNCTIONAL | Parse and validate flags                          |

### Dependency Injection (Scope)

| Feature                       | Status              | Notes                |
| ----------------------------- | ------------------- | -------------------- |
| `NewScope(name)`              | ✅ FULLY_FUNCTIONAL | Creates DI scope     |
| `Provide(scope, constructor)` | ✅ FULLY_FUNCTIONAL | Register service     |
| `ProvideValue(scope, value)`  | ✅ FULLY_FUNCTIONAL | Register value       |
| `Invoke[T](scope)`            | ✅ FULLY_FUNCTIONAL | Get service          |
| `Child(name)`                 | ✅ FULLY_FUNCTIONAL | Create child scope   |
| `Shutdown(ctx)`               | ✅ FULLY_FUNCTIONAL | Cleanup services     |
| `HealthCheck()`               | ✅ FULLY_FUNCTIONAL | Check service health |
| `IsRoot()`                    | ✅ FULLY_FUNCTIONAL | Check if root scope  |
| `Path()`                      | ✅ FULLY_FUNCTIONAL | Scope hierarchy path |

### Error Handling

| Feature                    | Status              | Notes                                      |
| -------------------------- | ------------------- | ------------------------------------------ |
| Typed errors               | ✅ FULLY_FUNCTIONAL | ErrInvalidCommand, ErrMissingHandler, etc. |
| Error wrapping             | ✅ FULLY_FUNCTIONAL | Compatible with errors.Is/As               |
| NewCommandError            | ✅ FULLY_FUNCTIONAL | Command-specific errors                    |
| NewServiceError            | ✅ FULLY_FUNCTIONAL | DI service-specific errors                 |
| No panics                  | ✅ FULLY_FUNCTIONAL | All operations return errors               |
| FlagError with suggestion  | ✅ FULLY_FUNCTIONAL | Includes typo suggestion in error message  |
| NewFlagErrorWithSuggestion | ✅ FULLY_FUNCTIONAL | Creates FlagError with suggestion text     |

### Helper Types

|                            | Feature             | Status                              | Notes |
| -------------------------- | ------------------- | ----------------------------------- | ----- |
| `LogLevel` type            | ✅ FULLY_FUNCTIONAL | Enum for debug/info/warn/error      |
| `LogLevel.SlogLevel()`     | ✅ FULLY_FUNCTIONAL | Converts to slog.Level              |
| `LogLevel.UnmarshalText()` | ✅ FULLY_FUNCTIONAL | Validates against allowed values    |
| `Enum[T]` type             | ✅ FULLY_FUNCTIONAL | Generic enum with validation        |
| `NoFlags` type             | ✅ FULLY_FUNCTIONAL | Sentinel for commands without flags |

---

## v1 API (pkg/cmdguard)

The v1 Guard API provides panic-at-construction validation.

| Feature            | Status              | Notes                              |
| ------------------ | ------------------- | ---------------------------------- |
| `New(name, short)` | ✅ FULLY_FUNCTIONAL | Creates guarded root command       |
| `AddCommand(cmd)`  | ✅ FULLY_FUNCTIONAL | Adds subcommand, panics if invalid |
| `Execute(ctx)`     | ✅ FULLY_FUNCTIONAL | Runs command with context          |
| `IsStrictMode()`   | ✅ FULLY_FUNCTIONAL | Returns strict mode status         |

---

## Dependencies

| Dependency                      | Version | Status              | Purpose              |
| ------------------------------- | ------- | ------------------- | -------------------- |
| `github.com/spf13/cobra`        | v1.10.2 | ✅ FULLY_FUNCTIONAL | CLI framework        |
| `github.com/samber/do/v2`       | v2.0.0  | ✅ FULLY_FUNCTIONAL | Dependency injection |
| `github.com/charmbracelet/fang` | v0.4.4  | ✅ FULLY_FUNCTIONAL | Cobra styling        |
| `github.com/stretchr/testify`   | v1.11.1 | ✅ FULLY_FUNCTIONAL | Testing              |

---

## Testing

| Package            | Coverage | Status  |
| ------------------ | -------- | ------- |
| `pkg/cmdguard/v2`  | 90%+     | ✅ Good |
| `pkg/cmdguard`     | 94.3%    | ✅ Good |
| `internal/config`  | 95.7%    | ✅ Good |
| `internal/logging` | 100%     | ✅ Good |

---

## Architecture

### v2 API Design

```
v2.New[AppConfig]("myapp", "My CLI", AppConfig{})
    └── GuardedCommand[AppConfig]
        ├── AddCommand(Command[AppConfig]{...}) - returns error
        ├── Scope() - DI injector for services
        ├── Execute(ctx) - run CLI
        └── Shutdown(ctx) - cleanup
```

**Key Principles:**

1. **Type Safety** - Generic type parameter for config
2. **No Panics** - All operations return errors
3. **DI-Powered** - samber/do/v2 integration
4. **Typed Flags** - Struct tags for flag definitions

---

## Feature Roadmap

### Phase 1: v2 Foundation ✅ COMPLETE

- [x] Implement GuardedCommand[T] with generics
- [x] Implement Command[T] with typed flags
- [x] Implement Scope for DI
- [x] Implement FlagRegistry with struct tags
- [x] Comprehensive error types

### Phase 2: v2 Testing ✅ COMPLETE

- [x] Test errors.go
- [x] Test types.go
- [x] Test config.go
- [x] Test flags.go
- [x] Test scope.go
- [x] Test command.go
- [x] Test guard.go

### Phase 3: v2 Polish ✅ COMPLETE

- [x] Add typed example (examples/typed/main.go)
- [x] Update documentation

### Phase 4: Beyond (Long Term)

- [ ] Plugin system for custom validators
- [ ] Enhanced flag validation (enums, custom validators)
- [ ] Performance benchmarks
- [ ] Release automation

---

## Honest Assessment Summary

### What Works Well ✅

- Type-safe CLI construction with generics
- No panics - all operations return errors
- DI integration with samber/do/v2
- Typed flags with struct tags
- Required flag validation
- Flag typo suggestions with Levenshtein distance
- LogLevel to slog.Level conversion
- Comprehensive test coverage
- Clean public API
- Example tests demonstrating API usage

### What Needs Work ⚠️

- Documentation could be more comprehensive

### What's Missing 🔧

- Plugin system for custom validators (planned)
- Enhanced flag validation (planned)

### Overall Status

**cmdguard v2 is production-ready.**

The v2 API successfully delivers type-safe, DI-powered CLI construction without panics. All core packages have comprehensive test coverage.

**Recommendation:** v2.0.0 release candidate.

---

_This document reflects the v2 API implementation. Last updated 2026-02-15._
