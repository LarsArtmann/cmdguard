# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [2.5.0] - 2026-06-10

### Added

- **Runtime type guards in `dispatchParse`** — Verifies return type matches handler target
- **Short tag validation** — `short` tag must be exactly 1 character; returns error at registration
- **Nil tracer guard in `TelemetryMiddleware`** — Returns error instead of nil dereference
- **`BranchingFlowContext.SetValue` child safety** — Skips children that have local key set
- **Mutable slice protection** — `Tags()` and `Path()` return cloned slices
- **Arg validator error returns** — `WithExactArgs`, `WithMinimumArgs`, etc. return errors instead of panicking for invalid args

### Changed

- **Removed all Must\* panic-inducing functions** — Zero panics in library code. Every function returns errors:
  - `MustNewCommand`, `MustNewParentCommand` (use `NewCommand`, `NewParentCommand`)
  - `MustNewCLI`, `MustAddCommand` (use `NewCLI`, `AddCommand`)
  - `MustVersionCommand` (use `VersionCommand`)
  - `MustDoctorCommand` (use `DoctorCommand`)
  - `MustInvoke[T]`, `MustInvokeNamed[T]` (use `Invoke[T]`)
  - `MustGet[T]`, `RequireBranchingFlowContext` (use `GetBranchingFlowContext`)
  - `MustParse[T]`, `MustParseDuration`, `MustParseLogLevel`, `MustParseLogFormat`, `MustParseEnum`
  - `MustParseURL`, `MustParseEmail`, `MustParsePort`, `MustParseFilePath`, `MustParseHostPort`
- **Split `errors.go` into domain-specific files** — `errors_command.go`, `errors_flags.go`, `errors_config.go`, `errors_di.go`
- **Unified `fieldValueToString`/`formatFieldValue`** — Single canonical `formatFieldValue()` in `flag_helpers.go`
- **`doc.go` "never panics" claim now truthful** — No Must\* functions remain
- **Consolidated `command_suggest.go` into `flags_suggest.go`** — Single file for all typo suggestions

### Fixed

- **`setStringField` panic on type mismatch** — Added `AssignableTo` guard before `field.Set`
- **`ExitError.Error()` nil panic** — Added nil guard for `e.Err`
- **Double-wrapping of `ErrServiceConstruction` in `ShutdownAll`** — Fixed in v2.4.0

### Removed

- All 16+ Must\* panic variants (see "Changed" section above)

---

## [2.4.0] - 2026-06-03

### Added

- **`--no-color` flag + `NO_COLOR` support** - `cli.NoColor()` returns true when `--no-color` passed or `NO_COLOR=1` env var is set

### Changed

- **Removed `FlowContextAccessor`** - Dead API with zero consumers; use `GetBranchingFlowContext(ctx)` directly
- **Removed string-based `BranchWithTimeout`/`BranchWithDeadline`** - Use typed `BranchWithDuration(name, time.Duration)` and `BranchWithDeadlineTime(name, time.Time)` instead
- **`TimingMiddleware` callback signature** - Now includes `error` parameter to distinguish success vs failure timing
- **Migrated to `charm.land` vanity imports** - `charmbracelet/glamour` → `charm.land/glamour/v2`
- **Replaced `gopkg.in/yaml.v3`** with `github.com/go-faster/yaml` for YAML config loading
- **Modernized to Go 1.26** - `errors.As` → `errors.AsType` in examples

### Fixed

- Resolved `tparallel` lint issue in `TestCLINoColor` by extracting `t.Setenv` subtest

---

## [2.3.0] - 2026-06-01

### Added

- **Interactive prompts** - `WithPromptOnMissing[T,F]()` with `prompt:"Question?"` struct tag via charmbracelet/huh
- **Glamour markdown help** - `WithGlamourHelp[T]()` and `WithGlamourHelpTheme[T](theme)` render command Long/Example as markdown
- **Spinner middleware** - `SpinnerMiddleware[T](title)` and `SpinnerMiddlewareWithConfig[T](config)` with TTY auto-detection
- **Telemetry middleware** - `TelemetryMiddleware[T](tracer)` creates OpenTelemetry span per command; `WithTelemetry[T](tracer)` convenience option
- **Config file loading** - `WithConfigFile[T](paths...)` for JSON, `WithConfigFileLoader[T](loader, paths...)` for YAML/TOML via `configload` sub-package
- **`FullPath` in CommandInfo** - Full command path (`root.subcmd.leaf`) available in middleware
- **`Phase` typed enum** - Replaces `CommandInfo.Phase string` with type-safe enum
- **Taskctl example** - Single comprehensive showcase replacing 13 scattered examples
- **Nix flake** - `flake.nix` with devShell (Go 1.26, gopls, golangci-lint), treefmt formatter, format check
- **Fuzz tests** - 7 fuzz targets for all value type parsers

### Changed

- **Architecture hardening** - Extracted `handlerConfig[T,F]` from 8-param wire function, consolidated 5 error types into `labeledError`
- **File splitting** - Split `type_handler.go` (481 lines), `command.go` (403 lines), `flow_context.go` (396 lines) into focused files
- **Deduplication** - Eliminated all semantic clone groups (0 remaining)
- **Error wrapping** - All 40+ errors use `fmt.Errorf("%w: ...", sentinel)` for `errors.Is()` chainability
- **Contextual errors** - Added context wrapping to CLI, command, and flag initialization errors
- **Value types** - Normalized `IsEmpty()` across all types, extracted shared `requireNonEmpty` helper
- **Removed `Ptr[T]`** - Use Go 1.26 built-in `new(v)` instead

### Fixed

- **CommandInfo concurrency** - Move copy inside handler closure for goroutine safety
- **Config file race** - Eliminate race in `resolveConfigFlag`
- **Enum validation** - Validate Enum values in SetField instead of silently bypassing allowed list
- **Required tag parsing** - Propagate error instead of silently ignoring
- **MergeConfigs** - True deep copy for reference types
- **Flag cloning** - Deep copy in `cloneFlags` to prevent shared mutable state
- **Validator registry** - Clone from global defaults per FlagRegistry instance

---

## [2.2.0] - 2026-04-30

### Added

- **`env:"VAR"` struct tag** - Environment variable binding with `WithEnvPrefix[T]("MYAPP_")` prefix support
- **`count:"true"` struct tag** - Counting flags: `-vvv` → 3
- **12 output formats** - go-output integration: table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, and more
- **`EditInEditor(ctx, content)`** - Open `$EDITOR` for user input (now with `context.Context`)
- **`WithFang[T](bool)`** - Proper name for styled help (deprecates `WithColor`)
- **Subcommand typo suggestions** - `SuggestCommand` with Levenshtein distance
- **`WithSignalHandling[T]()`** - SIGINT/SIGTERM context cancellation
- **Instance-scoped validators** - `FlagRegistry.RegisterFlagValidator()` for per-instance customization
- **Custom value types** - URL, Email, Port, FilePath, HostPort with full validation
- **Shell completion** - `WithCompletion[T,F](fn)` and `WithValidArgs[T,F](args...)`
- **Man page generation** - `cli.ManPage()`, `cli.WriteManPage()`, `GenerateManPageCommand[T](cli)`
- **`WithOutputFormat[T]()`** - Auto `--output` flag with format selection
- **Command groups** - `WithGroup[T](id, title)` for organized help output
- **Flag validation** - `validate:"email,min=5"` tag with built-in + custom validators

### Changed

- **TypeHandler registry** - Unified 3-way split brain into single registry
- **Validator registry** - Instance-scoped, removing global mutable state
- **Constructor pattern** - All Command creation via `NewCommand`/`NewParentCommand`, struct fields unexported
- **0 lint issues** - From 113 to 0, added comprehensive `.golangci.yml`
- **Refactored output.go** - Monolithic switch → format renderer registry

### Fixed

- **55 race conditions** - `sync.RWMutex` on `globalTypeRegistry`
- **`Ptr[T]` function** - Returned zero-valued pointer instead of pointer to v
- **`SuggestFlag` API** - Returns `(string, bool)` for consistency
- **envPrefix propagation** - Now reaches command-level FlagRegistry
- **Local go-output replace** - Removed after tagging v0.1.0

### Removed

- **Dead code** - `parseCustomDefault`, `wrapErr`, `parseField`, `parseAndSetLog*`
- **Scattered examples** - Consolidated into single `examples/taskctl/` showcase

---

## [2.1.0] - 2026-03-28

### Added

- **CLI[T] New API** - Type-safe CLI builder with generic config and flags
- **BranchingFlowContext** - Tracks command execution path with context propagation
- **Option[T] type** - Functional options pattern for CLI configuration
- **SimpleCLI alias** - Backward compatibility alias for CLI[T]
- **Functional options** - CLIOption[T] for configuring CLI instances
- **WithCLIScope option** - Inject existing DI scope
- **WithLong option** - Set long description via option
- **Comprehensive flow_context tests** - 344 lines, 44 tests
- **CLI[T] integration tests**

### Changed

- **Deprecated GuardedCommand** - Use CLI[T] instead (see migration guide)
- **Updated examples/typed to use SimpleCLI pattern**
- **Improved godoc for public APIs**
- **Updated README.md with links to new docs**

### Fixed

- **flow_context cancel bug** - Self-cancel tracking now works correctly
- **wrapcheck errors** - json.Marshal/Unmarshal now properly annotated
- **go.mod compatibility** - Fixed for Go 1.26
- **go:fix inline directive** - Removed duplicate in type_helpers.go

### Removed

- **AddCommandFunc** - Redundant, use AddCommand instead
- **Dead code packages** - Removed pkg/apperrors, pkg/testutil

### Documentation

- **Migration Guide v1→v2** - docs/MIGRATION_v1_v2.md
- **Quickstart Guide** - docs/QUICKSTART.md
- **README updates** - Links to new documentation

---

## [2.0.0] - 2026-03-22

### Added

- **v2 API** - Type-safe API with dependency injection
- **samber/do/v2 integration** - Full DI support with scopes
- **Typed errors** - Sentinel errors for precise error handling
- **Struct-based flags** - FlagRegistry with struct tag support
- **Comprehensive tests** - 90%+ coverage across packages
- **BDD-style tests** - Behavior documentation via Ginkgo/Gomega
- **Benchmark tests** - Performance validation

### Features

- Generic Command[T, F] type
- GuardedCommand[T, F] for type-safe CLI building
- DI-powered service management
- Flag parsing and validation
- Typo suggestions for unknown flags
- Subcommand support with different flag types
- PreRunE/PostRunE hooks
- Scoped providers for plugin architecture

### Custom Types

- Enum with validation
- Duration with parsing
- LogLevel (debug, info, warn, error)
- LogFormat (text, json)

---

## [0.1.0] - 2026-02-20

### Added

- Initial release of cmdguard v2
- Type-safe CLI construction with generics
- Dependency injection via samber/do/v2
- Flag binding with struct tags
- Full Cobra integration

[Unreleased]: https://github.com/larsartmann/cmdguard/compare/v2.4.0...HEAD
[2.4.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.4.0
[2.3.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.3.0
[2.2.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.2.0
[2.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.1.0
[2.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.0.0
[0.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v0.1.0
