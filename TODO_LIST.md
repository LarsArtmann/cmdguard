# TODO List

**Updated:** 2026-04-18
**Status:** v2.1.0 — Post-v1 removal hardening in progress.

## Completed (Multi-Session Sprint) ✅

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
- [x] Clean .golangci.yml stale references (guard_flags.go, errtypes, koanf, testify, ginkgo)

## Remaining Work

### 🔴 High Priority

- [ ] Fix benchmarks to use `NewCLI` instead of deprecated `New`
- [ ] Unify type dispatch into `TypeHandler` registry (eliminate 3-way split brain)
- [ ] Fix custom type registration in `flags.go` (only handles 4 of 8 types)
- [ ] Add fuzz tests to value type parsers (URL, Email, Port, FilePath, HostPort)
- [ ] Make validator registry instance-scoped (remove global mutable state)

### 🟡 Medium Priority

- [ ] Add fuzz tests to `flags_parse.go` and `config_parsing.go`
- [ ] Improve flag suggestion algorithm
- [ ] Merge split test helper files (`test_helpers_test.go` + `testhelpers_test.go`)

### 📚 Documentation

- [x] API Reference documentation (godoc examples)
- [ ] Update docs/QUICKSTART.md for v2.1 API
- [ ] DI Pattern Example in docs/
- [ ] Error Handling Example in docs/

### 📊 Performance

- [ ] Add comprehensive performance benchmarks
- [ ] Add benchmark regression detection to CI

### ⚙️ Release & CI

- [ ] Create v2.1.0 release tag and notes
- [ ] Set up release automation
- [ ] Add codecov integration
- [ ] Fix pre-commit hooks (currently pre-existing errors)

### 🔮 Future (v3.0+)

- [ ] Config file auto-loading with koanf
- [ ] Shell completion helpers
- [ ] Progress/Spinner Type (bubbles)
