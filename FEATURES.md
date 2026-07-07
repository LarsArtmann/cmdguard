# cmdguard Features

---

## Legend

|| Status | Meaning |
|| ----------------------- | -------------------------------------------------------- |
|| ✅ FULLY_FUNCTIONAL | Feature works as designed, tested, and documented |
|| ⚠️ PARTIALLY_FUNCTIONAL | Feature works but has limitations, gaps, or known issues |
|| 📝 PLANNED | Feature is designed but not yet implemented |
|| 📦 SUB-MODULE | Feature lives in an optional importable sub-module |

---

## v3 API (pkg/cmdguard/v3)

### CLI[T]

| Feature                                     | Status              | Notes                                |
| ------------------------------------------- | ------------------- | ------------------------------------ |
| `NewCLI[T](name, short, defaults, opts...)` | ✅ FULLY_FUNCTIONAL | Creates typed CLI, returns errors    |
| `AddCommand(cli, cmd)`                      | ✅ FULLY_FUNCTIONAL | Adds typed subcommand, returns error |
| `Execute(ctx)`                              | ✅ FULLY_FUNCTIONAL | Runs command with context            |
| `ExecuteWithArgs(ctx, args)`                | ✅ FULLY_FUNCTIONAL | For testing                          |
| `ExecuteAndExit(ctx)`                       | ✅ FULLY_FUNCTIONAL | Run and os.Exit (respects ExitCoder) |
| `Scope()`                                   | ✅ FULLY_FUNCTIONAL | Returns DI scope                     |
| `Config()`                                  | ✅ FULLY_FUNCTIONAL | Returns typed config \*T             |
| `Shutdown(ctx)`                             | ✅ FULLY_FUNCTIONAL | Graceful shutdown                    |
| `HealthCheck()`                             | ✅ FULLY_FUNCTIONAL | Runs health checks                   |
| `RootCommand()`                             | ✅ FULLY_FUNCTIONAL | Returns underlying cobra.Command     |

### CLI Options

| Option                         | Status              | Notes                                         |
| ------------------------------ | ------------------- | --------------------------------------------- |
| `WithCLIVersion(v)`            | ✅ FULLY_FUNCTIONAL | Set version string                            |
| `WithCLILong(desc)`            | ✅ FULLY_FUNCTIONAL | Set long description                          |
| `WithCLIScope(scope)`          | ✅ FULLY_FUNCTIONAL | Custom DI scope                               |
| `WithSilenceErrors()`          | ✅ FULLY_FUNCTIONAL | Suppress cobra error printing                 |
| `WithSilenceUsage()`           | ✅ FULLY_FUNCTIONAL | Suppress usage on error (now the default)     |
| `WithFang(bool)`               | ✅ FULLY_FUNCTIONAL | Enable/disable fang styling                   |
| `WithFangOptions(opts...)`     | ✅ FULLY_FUNCTIONAL | Pass fang options                             |
| `WithMiddleware[T](mw...)`     | ✅ FULLY_FUNCTIONAL | Add command middleware                        |
| `WithGroup(id,title)`          | ✅ FULLY_FUNCTIONAL | Command groups in help                        |
| `WithEnvPrefix(pfx)`           | ✅ FULLY_FUNCTIONAL | Prefix for env var lookups                    |
| `WithSignalHandling()`         | ✅ FULLY_FUNCTIONAL | Auto SIGINT/SIGTERM ctx cancellation          |
| `WithGracefulShutdown()`       | ✅ FULLY_FUNCTIONAL | Graceful DI service shutdown on signals       |
| `WithDILogging(logf)`          | ✅ FULLY_FUNCTIONAL | Internal DI container logging                 |
| `WithCLICommit(c)`             | ✅ FULLY_FUNCTIONAL | Git commit hash (auto-piped to fang)          |
| `WithFangErrorHandler(fn)`     | ✅ FULLY_FUNCTIONAL | Custom fang error display                     |
| `WithFangColorScheme(fn)`      | ✅ FULLY_FUNCTIONAL | Custom fang color theme                       |
| `WithOutputFormat(fmt)`        | ✅ FULLY_FUNCTIONAL | Auto --output flag with format selection      |
| `WithConfigValidation[T](fn)`  | ✅ FULLY_FUNCTIONAL | Validate config after flag parsing            |
| `WithStrictValidation()`       | ✅ FULLY_FUNCTIONAL | Require short desc on all commands            |
| `WithDraconianValidation()`    | ✅ FULLY_FUNCTIONAL | Strict + examples on leaf commands            |
| `glamour.WithHelp()`           | 📦 SUB-MODULE       | Markdown rendering for help text (auto theme) |
| `WithPostFlagParse[T](fns...)` | ✅ FULLY_FUNCTIONAL | Post-parse hook: DI init, session storage     |
| `WithCleanup[T](fns...)`       | ✅ FULLY_FUNCTIONAL | Post-RunE cleanup that fires on error too     |

### Cobra Escape Hatch (Raw Cobra Subcommands)

For consumers that register raw `*cobra.Command` subcommands via
`cli.RootCommand().AddCommand` (gradual migration from plain Cobra). These APIs
bridge the typed cmdguard world to raw cobra handlers.

| Feature                          | Status              | Notes                                                      |
| -------------------------------- | ------------------- | ---------------------------------------------------------- |
| `ConfigFromContext[T](ctx)`      | ✅ FULLY_FUNCTIONAL | Type-safe config retrieval for raw cobra RunE handlers     |
| `ArgsFromContext(ctx)`           | ✅ FULLY_FUNCTIONAL | Positional args for RunE handlers                          |
| `RegisterLocalCommandFlags(cmd)` | ✅ FULLY_FUNCTIONAL | Register root's local-scoped flags on a subcommand         |
| `RegisterScopedFlags(cmd)`       | ✅ FULLY_FUNCTIONAL | Register flags by scope (persistent vs local) on a command |
| `RegisterLocalFlags(cmd)`        | ✅ FULLY_FUNCTIONAL | Register only local-scoped flags on a command              |

### Command[T, F]

| Feature                           | Status              | Notes                                                    |
| --------------------------------- | ------------------- | -------------------------------------------------------- |
| `NewCommand` / `NewParentCommand` | ✅ FULLY_FUNCTIONAL | Constructors with validation, zero panics                |
| `RunE`, `PreRunE`, `PostRunE`     | ✅ FULLY_FUNCTIONAL | Type-safe handlers                                       |
| `Validate()`                      | ✅ FULLY_FUNCTIONAL | Called by constructors                                   |
| Command options (21 total)        | ✅ FULLY_FUNCTIONAL | WithShort, WithFlags, WithPreRunE, args validators, etc. |

### Flag System

| Feature                        | Status              | Notes                                                   |
| ------------------------------ | ------------------- | ------------------------------------------------------- |
| Struct tag flags               | ✅ FULLY_FUNCTIONAL | `flag:"name" short:"n" default:"val" help:"desc"`       |
| `env:"VAR"` struct tag         | ✅ FULLY_FUNCTIONAL | Environment variable binding                            |
| `count:"true"` struct tag      | ✅ FULLY_FUNCTIONAL | Counting flags: -vvv → 3                                |
| Short flags                    | ✅ FULLY_FUNCTIONAL | `short:"n"` for `-n`                                    |
| Required flags                 | ✅ FULLY_FUNCTIONAL | `required:"true"` tag                                   |
| `validate:"email,min=5"` tag   | ✅ FULLY_FUNCTIONAL | Built-in + custom validators                            |
| Flag typo suggestions          | ✅ FULLY_FUNCTIONAL | Levenshtein distance-based                              |
| Subcommand typo suggestions    | ✅ FULLY_FUNCTIONAL | "did you mean?" for unknown subcommands                 |
| Instance-scoped validators     | ✅ FULLY_FUNCTIONAL | FlagRegistry.RegisterFlagValidator() (COW)              |
| TypeHandler registry           | ✅ FULLY_FUNCTIONAL | Extensible type dispatch system (COW)                   |
| `RegisterTypeHandler()`        | ✅ FULLY_FUNCTIONAL | Register custom flag types                              |
| Iterator methods (`iter.Seq`)  | ✅ FULLY_FUNCTIONAL | TagsSeq, FlagNamesSeq, PathSeq, ChildrenSeq             |
| Integer overflow validation    | ✅ FULLY_FUNCTIONAL | int8/16/32, uint8/16 range-checked → ErrIntegerOverflow |
| Scoped flags (`local:"true"`)  | ✅ FULLY_FUNCTIONAL | Root-only flags not inherited by subcommands            |
| Hidden flags (`hidden:"true"`) | ✅ FULLY_FUNCTIONAL | Exclude from --help, stay functional                    |

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

| Feature                                 | Status              | Notes                                         |
| --------------------------------------- | ------------------- | --------------------------------------------- |
| `NewScope(name)`                        | ✅ FULLY_FUNCTIONAL | Creates DI scope                              |
| `NewScopeWithOpts(name, opts)`          | ✅ FULLY_FUNCTIONAL | Scope with custom InjectorOpts (logging, etc) |
| `Provide[T]`, `ProvideValue[T]`         | ✅ FULLY_FUNCTIONAL | Register services                             |
| `Invoke[T]`, `InvokeNamed[T]`           | ✅ FULLY_FUNCTIONAL | Get services                                  |
| `Override[T]`, `OverrideValue[T]`       | ✅ FULLY_FUNCTIONAL | Replace services for testing                  |
| `CloneScope(scope)`                     | ✅ FULLY_FUNCTIONAL | Clone scope for test isolation                |
| `Child(name)`                           | ✅ FULLY_FUNCTIONAL | Hierarchical scopes                           |
| `RootScope()`                           | ✅ FULLY_FUNCTIONAL | Navigate to root from any child               |
| `Shutdown`, `ShutdownAll`               | ✅ FULLY_FUNCTIONAL | Graceful service shutdown                     |
| `HealthCheck`, `HealthCheckWithContext` | ✅ FULLY_FUNCTIONAL | Lifecycle management                          |

### Rich Output (go-output)

| Feature                 | Status              | Notes                                                                                |
| ----------------------- | ------------------- | ------------------------------------------------------------------------------------ |
| `OutputResult()`        | ✅ FULLY_FUNCTIONAL | Shape-aware rendering with go-output v0.30.1 registries                              |
| `OutputTable()`         | ✅ FULLY_FUNCTIONAL | Convenience for table data with AddRowChecked validation                             |
| `RegisteredFormats()`   | ✅ FULLY_FUNCTIONAL | Dynamic format discovery from registered marshalers                                  |
| 16 output formats       | ✅ FULLY_FUNCTIONAL | table/json/csv/tsv/md/xml/d2/yaml/html/tree/mermaid/dot/jsonl/asciidoc/toml/plantuml |
| Dynamic `--output` help | ✅ FULLY_FUNCTIONAL | Auto-generated from RegisteredTableDataFormats()                                     |

### Middleware

| Feature                           | Status              | Notes                               |
| --------------------------------- | ------------------- | ----------------------------------- |
| `TimingMiddleware`                | ✅ FULLY_FUNCTIONAL | Log command execution duration      |
| `RecoveryMiddleware`              | ✅ FULLY_FUNCTIONAL | Recover from panics in handlers     |
| `spinner.Middleware[T]`           | 📦 SUB-MODULE       | Terminal spinner during execution   |
| `spinner.MiddlewareWithConfig[T]` | 📦 SUB-MODULE       | Configurable spinner (frames/speed) |
| `telemetry.Middleware[T]`         | 📦 SUB-MODULE       | OpenTelemetry span per command      |
| `CommandInfo.FullPath`            | ✅ FULLY_FUNCTIONAL | Full command path for middleware    |
| Custom middleware                 | ✅ FULLY_FUNCTIONAL | `func(ctx, cfg, info, next) error`  |

### Shell Completion

|                          | Feature             | Status                                  | Notes |
| ------------------------ | ------------------- | --------------------------------------- | ----- |
| `WithCompletion(fn)`     | ✅ FULLY_FUNCTIONAL | Dynamic shell completion                |
| `WithValidArgs(args...)` | ✅ FULLY_FUNCTIONAL | Static valid arguments                  |
| `CompletionFunc` type    | ✅ FULLY_FUNCTIONAL | Compatible with cobra ValidArgsFunction |

### Man Page Generation (manpage sub-module)

| Feature                       | Status        | Notes                       |
| ----------------------------- | ------------- | --------------------------- |
| `manpage.Generate[T](cli, n)` | 📦 SUB-MODULE | Generate roff man page      |
| `manpage.Write[T](cli, w, n)` | 📦 SUB-MODULE | Write man page to io.Writer |
| `manpage.GenerateCommand[T]`  | 📦 SUB-MODULE | Create `man` subcommand     |

### Markdown Help (glamour sub-module)

| Feature                               | Status        | Notes                                   |
| ------------------------------------- | ------------- | --------------------------------------- |
| `glamour.WithHelp()`                  | 📦 SUB-MODULE | Render command Long/Example as markdown |
| `glamour.WithHelpTheme(theme)`        | 📦 SUB-MODULE | Override the glamour theme              |
| `glamour.RenderMarkdown(md)`          | 📦 SUB-MODULE | Render markdown with auto theme         |
| `glamour.RenderMarkdownWithTheme(md)` | 📦 SUB-MODULE | Render with specific glamour theme      |

### Positional Arguments

| Feature                   | Status              | Notes                             |
| ------------------------- | ------------------- | --------------------------------- |
| `WithExactArgs(n)`        | ✅ FULLY_FUNCTIONAL | Require exactly n positional args |
| `WithMinimumArgs(n)`      | ✅ FULLY_FUNCTIONAL | Require at least n args           |
| `WithMaximumArgs(n)`      | ✅ FULLY_FUNCTIONAL | Allow at most n args              |
| `WithRangeArgs(min, max)` | ✅ FULLY_FUNCTIONAL | Require between min and max args  |
| `WithNoArgs()`            | ✅ FULLY_FUNCTIONAL | Reject any positional args        |
| `WithArgs(fn)`            | ✅ FULLY_FUNCTIONAL | Custom cobra.PositionalArgs       |

### Interactive Prompts (huh)

| Feature                         | Status              | Notes                                    |
| ------------------------------- | ------------------- | ---------------------------------------- |
| `WithPromptOnMissing()`         | ✅ FULLY_FUNCTIONAL | Prompt for missing `prompt`-tagged flags |
| `prompt:"Question?"` struct tag | ✅ FULLY_FUNCTIONAL | Marks field for interactive prompting    |
| `PromptString(title, default)`  | ✅ FULLY_FUNCTIONAL | Text input (via PromptRunner interface)  |
| `PromptSelect(title, options)`  | ✅ FULLY_FUNCTIONAL | Selection (via PromptRunner interface)   |
| `PromptConfirm(title)`          | ✅ FULLY_FUNCTIONAL | Yes/no (via PromptRunner interface)      |
| `prompts.HuhRunner`             | 📦 SUB-MODULE       | huh/v2 PromptRunner implementation       |
| Bool fields → confirm prompt    | ✅ FULLY_FUNCTIONAL | Automatic prompt type selection          |
| Enum fields → select prompt     | ✅ FULLY_FUNCTIONAL | Automatic prompt type selection          |

### Doctor Command

| Feature                          | Status              | Notes                       |
| -------------------------------- | ------------------- | --------------------------- |
| `DoctorCommand[T](cli, opts...)` | ✅ FULLY_FUNCTIONAL | Typed doctor subcommand     |
| `WithDoctorCheck[T](name, run)`  | ✅ FULLY_FUNCTIONAL | Add custom diagnostic check |
| `WithDoctorShort[T](desc)`       | ✅ FULLY_FUNCTIONAL | Custom short description    |
| `WithDoctorLong[T](desc)`        | ✅ FULLY_FUNCTIONAL | Custom long description     |
| `WithDoctorGroupID[T](id)`       | ✅ FULLY_FUNCTIONAL | Command group ID            |

### Health Check Results

| Feature                                    | Status              | Notes                               |
| ------------------------------------------ | ------------------- | ----------------------------------- |
| `Scope.HealthCheckResults()`               | ✅ FULLY_FUNCTIONAL | Per-service health map              |
| `Scope.HealthCheckResultsWithContext(ctx)` | ✅ FULLY_FUNCTIONAL | Per-service health map with context |
| `CLI.HealthCheckResults()`                 | ✅ FULLY_FUNCTIONAL | Delegates to Scope                  |
| `CLI.HealthCheckResultsWithContext(ctx)`   | ✅ FULLY_FUNCTIONAL | Delegates to Scope                  |

### Plugin System

| Feature                         | Status              | Notes                                             |
| ------------------------------- | ------------------- | ------------------------------------------------- |
| `Plugin` interface              | ✅ FULLY_FUNCTIONAL | Bundle custom type handlers + validators          |
| `PluginRegistrar`               | ✅ FULLY_FUNCTIONAL | Scoped `TypeHandler()`/`Validator()` registration |
| `RegisterPlugin(plugin)`        | ✅ FULLY_FUNCTIONAL | Apply to global registries                        |
| `WithPlugin(plugin)`            | ✅ FULLY_FUNCTIONAL | Apply per-instance (CLI option)                   |
| `FlagRegistry.RegisterPlugin()` | ✅ FULLY_FUNCTIONAL | Apply per-FlagRegistry                            |

### Documentation Generation

| Feature               | Status              | Notes                                            |
| --------------------- | ------------------- | ------------------------------------------------ |
| `cli.GenerateDocs(w)` | ✅ FULLY_FUNCTIONAL | Markdown docs for full command tree to io.Writer |

### Audit Log Export

| Feature                               | Status              | Notes                                                                    |
| ------------------------------------- | ------------------- | ------------------------------------------------------------------------ |
| `WithAuditLog(plugin)`                | ✅ FULLY_FUNCTIONAL | Wire samber-do-auditlog into DI injector                                 |
| `ExportAuditLog[T](cli, cfg)`         | ✅ FULLY_FUNCTIONAL | Write audit snapshot to file or io.Writer                                |
| `AuditLogFormat` strong type          | ✅ FULLY_FUNCTIONAL | Validated enum with `ParseAuditLogFormat()` + `Valid()`                  |
| 11 export formats                     | ✅ FULLY_FUNCTIONAL | html, json, ndjson, csv, tsv, mermaid, dot, d2, plantuml, tree, htmltree |
| `AuditLogServiceByName[T](cli)`       | ✅ FULLY_FUNCTIONAL | Query a named service's audit info                                       |
| `AuditLogFailedServices[T](cli)`      | ✅ FULLY_FUNCTIONAL | List services that failed to construct                                   |
| `cli.AuditLog()` / `AuditLogReport()` | ✅ FULLY_FUNCTIONAL | Programmatic access to the plugin + snapshot                             |

### Version Command

| Feature                             | Status              | Notes                                |
| ----------------------------------- | ------------------- | ------------------------------------ |
| `VersionCommand[T](cli)`            | ✅ FULLY_FUNCTIONAL | Typed version subcommand             |
| `GenerateVersionCommand[T](cli, w)` | ✅ FULLY_FUNCTIONAL | Raw cobra command with custom writer |

### Helpers

| Feature          | Status              | Notes                     |
| ---------------- | ------------------- | ------------------------- |
| `ValueOrDefault` | ✅ FULLY_FUNCTIONAL | Nil-safe value access     |
| `MergeConfigs`   | ✅ FULLY_FUNCTIONAL | Deep merge config structs |

### Error Handling

| Feature                   | Status              | Notes                                                    |
| ------------------------- | ------------------- | -------------------------------------------------------- |
| 60 sentinel errors        | ✅ FULLY_FUNCTIONAL | ErrInvalidCommand, ErrMissingHandler, etc.               |
| Typed errors              | ✅ FULLY_FUNCTIONAL | CommandError, FlagError, ServiceError, etc.              |
| `ExitCoder` / `ExitError` | ✅ FULLY_FUNCTIONAL | Custom exit codes for ExecuteAndExit                     |
| `ExitCode(err) int`       | ✅ FULLY_FUNCTIONAL | Public exit-code mapping (nil→0, ExitCoder→code, else→1) |

---

## Dependencies

### Core (direct)

| Dependency                                  | Version | Status              | Purpose                  |
| ------------------------------------------- | ------- | ------------------- | ------------------------ |
| `github.com/spf13/cobra`                    | v1.10.2 | ✅ FULLY_FUNCTIONAL | CLI framework            |
| `github.com/samber/do/v2`                   | v2.0.0  | ✅ FULLY_FUNCTIONAL | Dependency injection     |
| `github.com/spf13/pflag`                    | v1.0.10 | ✅ FULLY_FUNCTIONAL | Flag parsing             |
| `charm.land/fang/v2`                        | v2.0.1  | ✅ FULLY_FUNCTIONAL | Cobra styling            |
| `github.com/larsartmann/go-output`          | v0.30.1 | ✅ FULLY_FUNCTIONAL | Rich output (16 formats) |
| `github.com/larsartmann/samber-do-auditlog` | v0.4.0  | ✅ FULLY_FUNCTIONAL | DI audit logging         |
| `github.com/knadh/koanf/v2`                 | v2.3.5  | ✅ FULLY_FUNCTIONAL | Config file loading      |

### Optional Sub-Modules (isolated — import only what you need)

| Sub-module  | Dependency                       | Version | Purpose                 |
| ----------- | -------------------------------- | ------- | ----------------------- |
| `glamour`   | `charm.land/glamour/v2`          | v2.0.1  | Markdown help rendering |
| `prompts`   | `charm.land/huh/v2`              | v2.0.3  | Interactive prompts     |
| `spinner`   | `charm.land/lipgloss/v2`         | v2.0.5  | Terminal spinner        |
| `telemetry` | `go.opentelemetry.io/otel/trace` | v1.44.0 | OpenTelemetry tracing   |
| `manpage`   | `muesli/mango` + `mango-cobra`   | v0.2.0  | Man page generation     |

---

## Testing

| Package                      | Coverage  | Status  |
| ---------------------------- | --------- | ------- |
| `pkg/cmdguard/v3`            | ~87.3%    | ✅ Good |
| `pkg/cmdguard/v3/configload` | ~87.5%    | ✅ Good |
| Benchmarks                   | 26 total  | ✅ Good |
| Fuzz tests                   | 7 targets | ✅ Good |

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

### Config File Loading

| Feature                                | Status              | Notes                                                        |
| -------------------------------------- | ------------------- | ------------------------------------------------------------ |
| `WithConfigFile[T](paths...)`          | ✅ FULLY_FUNCTIONAL | JSON loader, core package                                    |
| `WithConfigFileLoader[T](loader, ...)` | ✅ FULLY_FUNCTIONAL | Custom loader (YAML/TOML via configload)                     |
| `$ENV` and `~` expansion               | ✅ FULLY_FUNCTIONAL | Path expansion before loading                                |
| `--config` flag override               | ✅ FULLY_FUNCTIONAL | Overrides default search paths                               |
| Missing file = silent skip             | ✅ FULLY_FUNCTIONAL | Not an error when default path missing                       |
| `configload.YAML()`                    | ✅ FULLY_FUNCTIONAL | YAML loader sub-package                                      |
| `configload.TOML()`                    | ✅ FULLY_FUNCTIONAL | TOML loader sub-package                                      |
| `configload.Auto()`                    | ✅ FULLY_FUNCTIONAL | Sequential YAML→TOML→JSON (not extension-based)              |
| Nested struct config                   | ✅ FULLY_FUNCTIONAL | Inner structs flattened; FieldTag.Index tracks reflect path  |
| `configload.KoanfLoader()`             | ✅ FULLY_FUNCTIONAL | Nested config objects via koanf (e.g. `{"db":{"host":"x"}}`) |

### Flag Priority Chain

```
explicit flag → env:"VAR" (with optional prefix) → config file → default value
```

---

**Last updated 2026-06-22.**
