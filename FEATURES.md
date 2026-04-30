# cmdguard Features

**Last Updated:** 2026-04-30
**Version:** 2.2.0
**Go Version:** 1.26

---

## Legend

|| Status | Meaning |
|| ----------------------- | -------------------------------------------------------- |
|| ✅ FULLY_FUNCTIONAL | Feature works as designed, tested, and documented |
|| ⚠️ PARTIALLY_FUNCTIONAL | Feature works but has limitations, gaps, or known issues |
|| 📝 PLANNED | Feature is designed but not yet implemented |

---

## v2 API (pkg/cmdguard/v2)

### CLI[T]

| Feature                                     | Status              | Notes                                |
| ------------------------------------------- | ------------------- | ------------------------------------ |
| `NewCLI[T](name, short, defaults, opts...)` | ✅ FULLY_FUNCTIONAL | Creates typed CLI, never panics      |
| `AddCommand(cli, cmd)`                      | ✅ FULLY_FUNCTIONAL | Adds typed subcommand, returns error |
| `Execute(ctx)`                              | ✅ FULLY_FUNCTIONAL | Runs command with context            |
| `ExecuteWithArgs(ctx, args)`                | ✅ FULLY_FUNCTIONAL | For testing                          |
| `ExecuteAndExit(ctx)`                       | ✅ FULLY_FUNCTIONAL | Runs and calls os.Exit               |
| `Scope()`                                   | ✅ FULLY_FUNCTIONAL | Returns DI scope                     |
| `Config()`                                  | ✅ FULLY_FUNCTIONAL | Returns typed config \*T             |
| `Shutdown(ctx)`                             | ✅ FULLY_FUNCTIONAL | Graceful shutdown                    |
| `HealthCheck()`                             | ✅ FULLY_FUNCTIONAL | Runs health checks                   |
| `RootCommand()`                             | ✅ FULLY_FUNCTIONAL | Returns underlying cobra.Command     |

### CLI Options

| Option                    | Status              | Notes                                    |
| ------------------------- | ------------------- | ---------------------------------------- |
| `WithCLIVersion[T](v)`    | ✅ FULLY_FUNCTIONAL | Set version string                       |
| `WithCLILong[T](desc)`    | ✅ FULLY_FUNCTIONAL | Set long description                     |
| `WithCLIScope[T](scope)`  | ✅ FULLY_FUNCTIONAL | Custom DI scope                          |
| `WithSilenceErrors[T]()`  | ✅ FULLY_FUNCTIONAL | Suppress cobra error printing            |
| `WithSilenceUsage[T]()`   | ✅ FULLY_FUNCTIONAL | Suppress usage on error                  |
| `WithFang[T](bool)`       | ✅ FULLY_FUNCTIONAL | Enable/disable fang styling              |
| `WithColor[T](bool)`      | 🗑️ DEPRECATED       | Use WithFang instead                     |
| `WithFangOptions[T]()`    | ✅ FULLY_FUNCTIONAL | Pass fang options                        |
| `WithMiddleware[T]()`     | ✅ FULLY_FUNCTIONAL | Add command middleware                   |
| `WithGroup[T](id,title)`  | ✅ FULLY_FUNCTIONAL | Command groups in help                   |
| `WithEnvPrefix[T](pfx)`   | ✅ FULLY_FUNCTIONAL | Prefix for env var lookups               |
| `WithSignalHandling[T]()` | ✅ FULLY_FUNCTIONAL | Auto SIGINT/SIGTERM ctx cancellation     |
| `WithOutputFormat[T]()`   | ✅ FULLY_FUNCTIONAL | Auto --output flag with format selection |

### Command[T, F]

| Feature                                   | Status              | Notes                                   |
| ----------------------------------------- | ------------------- | --------------------------------------- |
| `NewCommand` / `NewParentCommand`         | ✅ FULLY_FUNCTIONAL | Constructors with validation            |
| `MustNewCommand` / `MustNewParentCommand` | ✅ FULLY_FUNCTIONAL | Panic variants                          |
| `RunE`, `PreRunE`, `PostRunE`             | ✅ FULLY_FUNCTIONAL | Type-safe handlers                      |
| `Validate()`                              | ✅ FULLY_FUNCTIONAL | Called by constructors                  |
| Command options (15 total)                | ✅ FULLY_FUNCTIONAL | WithShort, WithFlags, WithPreRunE, etc. |

### Flag System

| Feature                      | Status              | Notes                                             |
| ---------------------------- | ------------------- | ------------------------------------------------- |
| Struct tag flags             | ✅ FULLY_FUNCTIONAL | `flag:"name" short:"n" default:"val" help:"desc"` |
| `env:"VAR"` struct tag       | ✅ FULLY_FUNCTIONAL | Environment variable binding                      |
| `count:"true"` struct tag    | ✅ FULLY_FUNCTIONAL | Counting flags: -vvv → 3                          |
| Short flags                  | ✅ FULLY_FUNCTIONAL | `short:"n"` for `-n`                              |
| Required flags               | ✅ FULLY_FUNCTIONAL | `required:"true"` tag                             |
| `validate:"email,min=5"` tag | ✅ FULLY_FUNCTIONAL | Built-in + custom validators                      |
| Flag typo suggestions        | ✅ FULLY_FUNCTIONAL | Levenshtein distance-based                        |
| Subcommand typo suggestions  | ✅ FULLY_FUNCTIONAL | "did you mean?" for unknown subcommands           |
| Instance-scoped validators   | ✅ FULLY_FUNCTIONAL | FlagRegistry.RegisterFlagValidator()              |
| TypeHandler registry         | ✅ FULLY_FUNCTIONAL | Extensible type dispatch system                   |
| `RegisterTypeHandler()`      | ✅ FULLY_FUNCTIONAL | Register custom flag types                        |

### Value Types

| Feature     | Status              | Notes                          |
| ----------- | ------------------- | ------------------------------ |
| `Duration`  | ✅ FULLY_FUNCTIONAL | time.Duration wrapper          |
| `Enum`      | ✅ FULLY_FUNCTIONAL | Validated enum values          |
| `LogLevel`  | ✅ FULLY_FUNCTIONAL | debug/info/warn/error          |
| `LogFormat` | ✅ FULLY_FUNCTIONAL | text/json                      |
| `URL`       | ✅ FULLY_FUNCTIONAL | Validated URL                  |
| `Email`     | ✅ FULLY_FUNCTIONAL | Validated email                |
| `Port`      | ✅ FULLY_FUNCTIONAL | 1-65535, named ports           |
| `FilePath`  | ✅ FULLY_FUNCTIONAL | Cleaned paths, existence check |
| `HostPort`  | ✅ FULLY_FUNCTIONAL | host:port validation           |

### Dependency Injection

| Feature                      | Status              | Notes                |
| ---------------------------- | ------------------- | -------------------- |
| `NewScope(name)`             | ✅ FULLY_FUNCTIONAL | Creates DI scope     |
| `Provide[T]`, `ProvideValue` | ✅ FULLY_FUNCTIONAL | Register services    |
| `Invoke[T]`                  | ✅ FULLY_FUNCTIONAL | Get services         |
| `Child(name)`                | ✅ FULLY_FUNCTIONAL | Hierarchical scopes  |
| `Shutdown`, `HealthCheck`    | ✅ FULLY_FUNCTIONAL | Lifecycle management |

### Rich Output (go-output)

| Feature               | Status              | Notes                                                   |
| --------------------- | ------------------- | ------------------------------------------------------- |
| `OutputResult()`      | ✅ FULLY_FUNCTIONAL | Render data in configured format                        |
| `OutputTable()`       | ✅ FULLY_FUNCTIONAL | Convenience for table data                              |
| `OutputStyledTable()` | ✅ FULLY_FUNCTIONAL | Lipgloss-styled terminal tables                         |
| `ParseOutputFormat()` | ✅ FULLY_FUNCTIONAL | String → Format conversion                              |
| 12 output formats     | ✅ FULLY_FUNCTIONAL | table/json/csv/tsv/md/xml/d2/yaml/html/tree/mermaid/dot |

### Middleware

| Feature              | Status              | Notes                              |
| -------------------- | ------------------- | ---------------------------------- |
| `TimingMiddleware`   | ✅ FULLY_FUNCTIONAL | Log command execution duration     |
| `RecoveryMiddleware` | ✅ FULLY_FUNCTIONAL | Recover from panics in handlers    |
| Custom middleware    | ✅ FULLY_FUNCTIONAL | `func(ctx, cfg, info, next) error` |

### Shell Completion

|                               | Feature             | Status                                  | Notes |
| ----------------------------- | ------------------- | --------------------------------------- | ----- |
| `WithCompletion[T,F](fn)`     | ✅ FULLY_FUNCTIONAL | Dynamic shell completion                |
| `WithValidArgs[T,F](args...)` | ✅ FULLY_FUNCTIONAL | Static valid arguments                  |
| `CompletionFunc` type         | ✅ FULLY_FUNCTIONAL | Compatible with cobra ValidArgsFunction |

### Man Page Generation

|                                  | Feature             | Status                      | Notes |
| -------------------------------- | ------------------- | --------------------------- | ----- |
| `cli.ManPage(section)`           | ✅ FULLY_FUNCTIONAL | Generate roff man page      |
| `cli.WriteManPage(w, section)`   | ✅ FULLY_FUNCTIONAL | Write man page to io.Writer |
| `GenerateManPageCommand[T](cli)` | ✅ FULLY_FUNCTIONAL | Create `man` subcommand     |

### Helpers

| Feature          | Status              | Notes                       |
| ---------------- | ------------------- | --------------------------- |
| `EditInEditor`   | ✅ FULLY_FUNCTIONAL | Open content in $EDITOR     |
| `Ptr[T]`         | ✅ FULLY_FUNCTIONAL | Pointer helper              |
| `ValueOrDefault` | ✅ FULLY_FUNCTIONAL | Nil-safe value access       |
| `MustParse[T]`   | ✅ FULLY_FUNCTIONAL | Panic-on-fail for constants |
| `MergeConfigs`   | ✅ FULLY_FUNCTIONAL | Deep merge config structs   |

### Error Handling

| Feature                   | Status              | Notes                                       |
| ------------------------- | ------------------- | ------------------------------------------- |
| 30+ sentinel errors       | ✅ FULLY_FUNCTIONAL | ErrInvalidCommand, ErrMissingHandler, etc.  |
| Typed errors              | ✅ FULLY_FUNCTIONAL | CommandError, FlagError, ServiceError, etc. |
| FlagError with suggestion | ✅ FULLY_FUNCTIONAL | Includes typo suggestion in error message   |
| No panics (library API)   | ✅ FULLY_FUNCTIONAL | All operations return errors                |

---

## Dependencies

| Dependency                         | Version | Status              | Purpose              |
| ---------------------------------- | ------- | ------------------- | -------------------- |
| `github.com/spf13/cobra`           | v1.10.2 | ✅ FULLY_FUNCTIONAL | CLI framework        |
| `github.com/samber/do/v2`          | v2.0.0  | ✅ FULLY_FUNCTIONAL | Dependency injection |
| `github.com/spf13/pflag`           | v1.0.10 | ✅ FULLY_FUNCTIONAL | Flag parsing         |
| `charm.land/fang/v2`               | v2.0.1  | ✅ FULLY_FUNCTIONAL | Cobra styling        |
| `github.com/larsartmann/go-output` | latest  | ✅ FULLY_FUNCTIONAL | Rich output formats  |

---

## Testing

| Package           | Coverage  | Status  |
| ----------------- | --------- | ------- |
| `pkg/cmdguard/v2` | 80.9%     | ✅ Good |
| Fuzz tests        | 7 targets | ✅ Good |

---

## Architecture

### TypeHandler Registry

```
TypeHandler interface {
    Register(flags, tag) error   // Add flag to pflag.FlagSet
    Parse(value, tag) (any, error)  // Parse string → Go value
    Default(tag) any              // Compute default value
}
```

All type dispatch (primitives + 9 custom types) flows through a single registry.
Custom types can be added via `RegisterTypeHandler(reflect.Type, TypeHandler)`.

### Flag Priority Chain

```
explicit flag → env:"VAR" (with optional prefix) → default value
```

---

**Last updated 2026-04-30.**
