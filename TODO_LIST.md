# TODO List

**Updated:** 2026-04-30
**Status:** v2.2.0 — SUPERB CLI/TUI features implemented.

## Completed ✅

### Phase 1: Foundation (1% → 51%)
- [x] Unify type dispatch into TypeHandler registry (eliminate 3-way split brain)
- [x] Fix custom type registration in `setStringField` (URL/Email/Port/FilePath/HostPort)
- [x] Make validator registry instance-scoped (remove global mutable state)
- [x] Fix `Ptr[T]` function (returned zero-valued pointer instead of pointer to v)
- [x] Clean up 52 leftover `_test.go_test` artifact files

### Phase 2: Core Features (4% → 64%)
- [x] Add `env:"VAR"` struct tag support with `WithEnvPrefix[T]`
- [x] Add subcommand typo suggestions (`SuggestCommand`)
- [x] Add `WithSignalHandling[T]()` for SIGINT/SIGTERM context cancellation
- [x] Add go-output integration (12 output formats)
- [x] Add `count:"true"` struct tag for counting flags (-vvv → 3)
- [x] Add `EditInEditor()` $EDITOR helper
- [x] Add `WithFang[T]()` as proper name (deprecate `WithColor`)
- [x] Add fuzz tests for all value type parsers and parsing functions

### Phase 2: Previous Sprint
- [x] Fix CLI[T] AddCommand flag parsing (cloneAndParseFlags pattern)
- [x] Refactor nestif complexity in flag_helpers.go
- [x] Fix err113 dynamic error issues
- [x] Add t.Parallel() to all v2 tests
- [x] Add tests for `initialize` error paths
- [x] Add tests for `cliToCobraCommand` edge cases
- [x] Add tests for flag helper functions
- [x] Add WithSilenceErrors, WithSilenceUsage, WithColor CLI options
- [x] Remove v1 API, internal packages, v1 integration tests (3,841 lines)
- [x] Remove Option[T]/Result[T] ghost types (1,501 lines)
- [x] Remove 6 ghost koanf dependencies
- [x] Fix nilnil, forcetypeassert, exhaustive, err113 lint issues
- [x] Update all docs (AGENTS.md, README.md, FEATURES.md) to remove v1 refs
- [x] Archive 31 old status reports
- [x] Delete .go_test template artifacts and empty internal/ directory
- [x] Clean .golangci.yml stale references

## Remaining Work

### 📚 Documentation
- [ ] Update docs/QUICKSTART.md for v2.2 API (env tags, output, counting flags, etc.)
- [ ] DI Pattern Example in docs/
- [ ] Error Handling Example in docs/

### 📊 Performance
- [ ] Add comprehensive performance benchmarks
- [ ] Add benchmark regression detection to CI

### ⚙️ Release & CI
- [ ] Create v2.2.0 release tag and notes
- [ ] Set up release automation
- [ ] Add codecov integration
- [ ] Fix pre-commit hooks (currently pre-existing errors)

### 🔮 Future (v3.0+)
- [ ] Config file auto-loading with koanf (YAML/TOML/.env)
- [ ] Interactive prompts (huh integration) with `WithPromptOnMissing`
- [ ] Spinner/progress middleware (bubbles)
- [ ] Shell completion helpers (`WithCompletion[T,F]`)
- [ ] Markdown help rendering (glamour)
- [ ] Man page generation via mango-cobra
- [ ] Telemetry middleware (OpenTelemetry spans)
- [ ] Plugin system for custom validators and type handlers
