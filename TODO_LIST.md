# TODO List

**Updated:** 2026-05-16
**Status:** v2.3.0-dev — 247 tests (210 in v2), 80.4% coverage, 0 lint issues, 0 race conditions

## Completed ✅

### Phase 1: Foundation

- [x] Unify type dispatch into TypeHandler registry (eliminate 3-way split brain)
- [x] Fix custom type registration in `setStringField` (URL/Email/Port/FilePath/HostPort)
- [x] Make validator registry instance-scoped (remove global mutable state)
- [x] Fix `Ptr[T]` function (returned zero-valued pointer instead of pointer to v)
- [x] Clean up 52 leftover `_test.go_test` artifact files

### Phase 2: Core Features

- [x] Add `env:"VAR"` struct tag support with `WithEnvPrefix[T]`
- [x] Add subcommand typo suggestions (`SuggestCommand`)
- [x] Add `WithSignalHandling[T]()` for SIGINT/SIGTERM context cancellation
- [x] Add go-output integration (12 output formats)
- [x] Add `count:"true"` struct tag for counting flags (-vvv → 3)
- [x] Add `EditInEditor()` $EDITOR helper (now with context.Context)
- [x] Add `WithFang[T]()` as proper name (deprecate `WithColor`)
- [x] Add fuzz tests for all value type parsers and parsing functions
- [x] Propagate envPrefix to command-level FlagRegistry

### Phase 2: Test & Cleanup Sprint

- [x] Add comprehensive tests for all v2.2 features
- [x] Remove dead code: parseCustomDefault, wrapErr, parseField, parseAndSetLog\*
- [x] Fix SuggestFlag API to return (string, bool) for consistency
- [x] Move count handler into registerKinds() for consistency

### Phase 3: Documentation & Examples

- [x] Update docs/QUICKSTART.md for v2.2 API
- [x] Update README.md with v2.2 features
- [x] Add 6 working examples (env-tags, output, counting, di-patterns, error-handling, signals)
- [x] Add subcommands example demonstrating NewParentCommand
- [x] Update examples/README.md with feature matrix

### Phase 4: New Features

- [x] Add `WithOutputFormat[T]` CLI option (auto --output flag)
- [x] Add shell completion wiring (`WithCompletion[T,F]`, `WithValidArgs[T,F]`)
- [x] Add man page generation via mango-cobra (`cli.ManPage()`, `GenerateManPageCommand`)

### Phase 5: Quality

- [x] Fix all 55 race conditions (sync.RWMutex on globalTypeRegistry)
- [x] Remove local go-output replace directive (tagged v0.1.0)
- [x] Achieve 0 lint issues (was 113)
- [x] Refactor output.go from monolithic switch to format renderer registry
- [x] Add sentinel errors ErrUnsupportedFormat, ErrFormatRequiresTypedData
- [x] Add context.Context to EditInEditor

### Phase 6: Architecture Hardening

- [x] Fix BranchingFlowContext double-cancellation bug
- [x] Fix Enum.Allowed() returning internal slice (defensive copy)
- [x] Use errors.Join in Scope.ShutdownAll for proper error chains
- [x] Fix Scope.Path() allocation (collect-then-reverse)
- [x] Extract shared lookupFlagInCommand (deduplicate flags.go/flags_parse.go)
- [x] Replace hardcoded type switch in getFieldValue with fmt.Stringer
- [x] Deduplicate custom type handler registrations (makeEnumLikeHandler + table-driven)
- [x] Add stack trace capture to RecoveryMiddleware
- [x] Simplify parseAndSetValue to delegate to SetField (remove duplication)
- [x] Use map[string]struct{} for command registration set
- [x] Document SetConfig FlagRegistry desync warning

### Phase 7: Tooling

- [x] Add shareable pre-commit hook script (scripts/pre-commit)
- [x] Add GitHub Actions CI workflow (build, test, lint)
- [x] Add NewParentCommand example

### Phase 8: v2.3 Features

- [x] Add `ExitCoder` interface + `ExitError` struct for custom exit codes in `ExecuteAndExit`
- [x] Add positional argument validators (`WithExactArgs`, `WithMinimumArgs`, `WithMaximumArgs`, `WithRangeArgs`, `WithNoArgs`, `WithArgs`)
- [x] Add `WithConfigValidation[T](fn)` CLI option for cross-field validation
- [x] Add `WithStrictValidation[T]()` CLI option requiring short descriptions on all commands
- [x] Add `VersionCommand[T](cli)`, `MustVersionCommand[T](cli)`, `GenerateVersionCommand[T](cli, w)`

### Phase 9: Architecture Hardening (v2.3)

- [ ] Fix gopls hint: `errors.As` → `errors.AsType[ExitCoder]` (Go 1.26 idiom)
- [ ] Extract `handlerConfig[T,F]` struct from 8-param `wireHandlerWithMiddleware`
- [ ] Add `Phase` typed enum to replace `CommandInfo.Phase string`
- [ ] Fix 7 unwrapped error returns (add `fmt.Errorf` context)
- [ ] Consolidate 5 error types into internal `labeledError`
- [ ] Split `type_handler.go` (481 lines) into 3 files
- [ ] Split `command.go` (403 lines) — extract args options
- [ ] Split `flow_context.go` (396 lines) — extract options
- [ ] Fix `outputFormat`/`outputState.format` split brain
- [ ] Consolidate value type MarshalText/UnmarshalText patterns

## Remaining Work

### ⚡ Performance

- [ ] Add CLI construction benchmark
- [ ] Add flag parsing benchmark
- [ ] Add command execution benchmark
- [ ] Add benchmark regression detection to CI

### ⚙️ CI/CD

- [ ] Add codecov integration
- [ ] Create v2.3.0 release tag and notes
- [ ] Set up release automation

### 🔮 Future (v3.0+)

- [ ] Config file auto-loading with koanf (YAML/TOML/.env)
- [ ] Interactive prompts (huh integration) with `WithPromptOnMissing`
- [ ] Spinner/progress middleware (bubbles)
- [ ] Glamour markdown help rendering
- [ ] Telemetry middleware (OpenTelemetry spans)
- [ ] Plugin system for custom validators and type handlers

### 🧹 Future Cleanup (API-breaking, defer to v3.0)

- [ ] Make NoFlags a distinct named type (not type alias)
- [ ] Change TimingMiddleware callback to include error
- [ ] Remove string-based BranchWithTimeout/BranchWithDeadline (replaced by typed alternatives)
- [ ] Remove FlowContextAccessor (thin wrapper with no added value)
- [ ] Rename Get[T]/MustGet[T] to more specific names
- [ ] Make RegisterInScope generic instead of `...any`
- [ ] Remove or redesign Package() for error-safe DI integration
