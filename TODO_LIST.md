# TODO List

**Updated:** 2026-04-30
**Status:** v2.2.0 — 199 tests, 80.6% coverage, 0 lint issues, 0 race conditions

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
- [x] Remove dead code: parseCustomDefault, wrapErr, parseField, parseAndSetLog*
- [x] Fix SuggestFlag API to return (string, bool) for consistency
- [x] Move count handler into registerKinds() for consistency

### Phase 3: Documentation & Examples
- [x] Update docs/QUICKSTART.md for v2.2 API
- [x] Update README.md with v2.2 features
- [x] Add 6 working examples (env-tags, output, counting, di-patterns, error-handling, signals)
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

## Remaining Work

### ⚡ Performance
- [ ] Add CLI construction benchmark
- [ ] Add flag parsing benchmark
- [ ] Add command execution benchmark
- [ ] Add benchmark regression detection to CI

### ⚙️ CI/CD
- [ ] Add GitHub Actions CI workflow (build, test, lint)
- [ ] Add codecov integration
- [ ] Fix pre-commit hooks
- [ ] Create v2.2.0 release tag and notes
- [ ] Set up release automation

### 📚 Examples
- [ ] Add NewParentCommand example

### 🔮 Future (v3.0+)
- [ ] Config file auto-loading with koanf (YAML/TOML/.env)
- [ ] Interactive prompts (huh integration) with `WithPromptOnMissing`
- [ ] Spinner/progress middleware (bubbles)
- [ ] Glamour markdown help rendering
- [ ] Telemetry middleware (OpenTelemetry spans)
- [ ] Plugin system for custom validators and type handlers
