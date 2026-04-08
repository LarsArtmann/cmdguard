# TODO List

**Updated:** 2026-04-05
**Status:** v2.1.0 — Code complete. Remaining: examples, docs, CI polish.

## Completed This Sprint ✅

- [x] Fix CLI[T] AddCommand flag parsing (cloneAndParseFlags pattern)
- [x] Refactor nestif complexity in flag_helpers.go
- [x] Fix err113 dynamic error issues (already clean)
- [x] Add t.Parallel() to all v2 tests
- [x] Add tests for `initialize` error paths
- [x] Add tests for `cliToCobraCommand` edge cases
- [x] Add tests for flag helper functions
- [x] Add WithSilenceErrors, WithSilenceUsage, WithColor CLI options
- [x] Rename pkg/errors to pkg/errtypes, BaseError to CodedError
- [x] Migrate all callers to NewCLI/AddCommand API
- [x] Remove deprecated GuardedCommand[T,F] code (1,624 lines)
- [x] Rename guard*\* files to cli*\* and flag_helpers
- [x] Update README.md with v2.1 API
- [x] Update AGENTS.md with v2.1 API patterns
- [x] Update FEATURES.md (remove deprecated section)

## Remaining Work

### 🟡 Medium Priority

- [ ] Improve flag suggestion algorithm
- [ ] Migrate remaining testify usage to stdlib (if any remains)
- [ ] Add fuzz tests to flags_parse.go and config_parsing.go

### 📚 Documentation

- [ ] API Reference documentation (godoc examples)
- [ ] Update docs/QUICKSTART.md for v2.1 API
- [ ] Update docs/MIGRATION_v1_v2.md for v2.1 API
- [ ] DI Pattern Example in docs/
- [ ] Error Handling Example in docs/

### 🎯 Examples

- [ ] Add example with real database connection
- [ ] Add lifecycle hook examples
- [ ] Advanced DI Example

### 📊 Performance

- [ ] Add comprehensive performance benchmarks
- [ ] Add benchmark regression detection to CI

### ⚙️ Release & CI

- [ ] Create v2.1.0 release tag and notes
- [ ] Set up release automation
- [ ] Add codecov integration
- [ ] Fix pre-commit hooks (currently 5 pre-existing errors)
- [ ] Migrate benchmarks from deprecated v2.New to v2.NewCLI

### 🔮 Future (v3.0+)

- [ ] Plugin system for custom validators
- [ ] Enhanced flag validation (enums, custom validators)
- [ ] Config file auto-loading with koanf
- [ ] Shell completion helpers
- [ ] Result[T] type for error handling
- [ ] Progress/Spinner Type (bubbles)
- [ ] Command groups feature
