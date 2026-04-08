# TODO Table View

**Updated:** 2026-04-05
**Status:** v2.1.0 — Code complete. 135 items pruned to remaining work.

## Completed This Sprint

| Task                                                    | Status              |
| ------------------------------------------------------- | ------------------- |
| Fix CLI[T] AddCommand flag parsing (cloneAndParseFlags) | ✅ Done             |
| Refactor nestif in flag_helpers.go                      | ✅ Done             |
| err113 linter check                                     | ✅ Clean (0 issues) |
| t.Parallel() in all v2 tests                            | ✅ Done             |
| Tests for initialize error paths                        | ✅ Done             |
| Tests for cliToCobraCommand edge cases                  | ✅ Done             |
| Tests for flag helper functions                         | ✅ Done             |
| WithSilenceErrors/WithSilenceUsage/WithColor options    | ✅ Done             |
| pkg/errors → pkg/errtypes, BaseError → CodedError       | ✅ Done             |
| Migrate all callers to NewCLI/AddCommand                | ✅ Done             |
| Remove deprecated GuardedCommand[T,F] (1,624 lines)     | ✅ Done             |
| Rename guard*\* to cli*\* and flag_helpers              | ✅ Done             |
| README.md rewrite for v2.1 API                          | ✅ Done             |
| AGENTS.md update for v2.1 API                           | ✅ Done             |
| FEATURES.md cleanup                                     | ✅ Done             |
| flags_parse_test.go complexity check                    | ✅ Clean (0 issues) |
| v1 t.Parallel() check                                   | ✅ Already had it   |
| testify/ginkgo removal check                            | ✅ 0 instances      |
| Benchmarks check                                        | ✅ Already migrated |
| Examples check                                          | ✅ All 4 migrated   |
| err113 linter audit                                     | ✅ Clean            |
| exhaustruct exclusions                                  | ✅ Done             |
| paralleltest fixes                                      | ✅ Done             |

## Remaining Work

### 🟡 Medium Priority

| Task                                | Effort | Impact |
| ----------------------------------- | ------ | ------ |
| Improve flag suggestion algorithm   | 15 min | Medium |
| Migrate remaining testify to stdlib | 15 min | Medium |
| Add fuzz tests to flags_parse.go    | 20 min | Medium |
| Add fuzz tests to config_parsing.go | 20 min | Medium |

### 📚 Documentation

| Task                                    | Effort | Impact |
| --------------------------------------- | ------ | ------ |
| API Reference (godoc examples)          | 20 min | Medium |
| Update docs/QUICKSTART.md for v2.1      | 15 min | Medium |
| Update docs/MIGRATION_v1_v2.md for v2.1 | 15 min | Medium |
| DI Pattern Example                      | 15 min | Medium |
| Error Handling Example                  | 15 min | Medium |

### 🎯 Examples

| Task                        | Effort | Impact |
| --------------------------- | ------ | ------ |
| Database connection example | 30 min | High   |
| Lifecycle hook examples     | 20 min | Medium |
| Advanced DI Example         | 20 min | Medium |

### 📊 Performance

| Task                     | Effort | Impact |
| ------------------------ | ------ | ------ |
| Comprehensive benchmarks | 30 min | Medium |
| Benchmark regression CI  | 10 min | High   |

### ⚙️ Release & CI

| Task                            | Effort | Impact |
| ------------------------------- | ------ | ------ |
| v2.1.0 release tag and notes    | 15 min | High   |
| Release automation              | 20 min | High   |
| Codecov integration             | 15 min | High   |
| Fix pre-commit hooks (5 errors) | 30 min | High   |
| Migrate benchmarks to NewCLI    | 15 min | Medium |

### 🔮 Future (v3.0+)

| Task                         | Effort | Impact |
| ---------------------------- | ------ | ------ |
| Plugin system for validators | 30 min | Low    |
| Enhanced flag validation     | 20 min | Low    |
| Config file auto-loading     | 30 min | Low    |
| Shell completion helpers     | 20 min | Low    |
| Result[T] type               | 25 min | Low    |
| Progress/Spinner (bubbles)   | 30 min | Low    |
| Command groups feature       | 30 min | Low    |

## Statistics

| Metric                  | Count |
| ----------------------- | ----- |
| Completed this sprint   | 23    |
| Remaining work items    | 23    |
| Future/deferred (v3.0+) | 7     |
