# Comprehensive TODO Table View

## Pending Tasks — Aligned with TODO_LIST.md (2026-04-05)

_Last audited against codebase_

### 🔴 High Priority (This Sprint)

| ID  | Task Description                                                        | Priority | Source                                        | Effort  | Impact |
| --- | ----------------------------------------------------------------------- | -------- | --------------------------------------------- | ------- | ------ |
| 1   | Fix CLI[T] AddCommand flag parsing bug using cloneAndParseFlags pattern | HIGH     | pkg/cmdguard/v2/cli.go:190                    | 20 min  | High   |
| 2   | Refactor flags_parse_test.go complexity                                 | HIGH     | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT  | 20 min  | Medium |
| 3   | Refactor nestif complexity in flags parsing                             | HIGH     | COMPREHENSIVE_STATUS_REPORT                   | 15 min  | Low    |
| 4   | Fix err113 dynamic error wrapping issues                                | HIGH     | 2026-03-22_14-12_comprehensive-status          | 15 min  | Medium |

### 🟡 Medium Priority

| ID  | Task Description                                     | Priority | Source                               | Effort  | Impact |
| --- | ---------------------------------------------------- | -------- | ------------------------------------ | ------- | ------ |
| 5   | Improve flag suggestion algorithm                    | MEDIUM   | 2026-03-22_16-58_COMPREHENSIVE_STATUS | 15 min  | Medium |
| 6   | Improve error types (more specific error categories) | MEDIUM   | 2026-03-22_16-58_COMPREHENSIVE_STATUS | 15 min  | Medium |
| 7   | Update README.md with ID usage examples              | MEDIUM   | go-composable-business-types-usage    | 15 min  | Medium |
| 8   | Update AGENTS.md integration patterns                | MEDIUM   | go-composable-business-types-usage    | 15 min  | Medium |
| 9   | Decide on Ginkgo vs stdlib testing and update docs   | MEDIUM   | BDD_TESTS_REVIEW                     | 10 min  | Low    |
| 10  | Add tests for `initialize` error paths               | MEDIUM   | 2026-03-23_15-34_COMPREHENSIVE_STATUS | 15 min  | Medium |
| 11  | Add tests for `cliToCobraCommand` edge cases         | MEDIUM   | 2026-03-23_15-34_COMPREHENSIVE_STATUS | 15 min  | Medium |
| 12  | Add tests for `cloneAndParseFlags` error paths       | MEDIUM   | 2026-03-23_15-34_COMPREHENSIVE_STATUS | 15 min  | Medium |

### 🧪 Testing & Refactoring

| ID  | Task Description                                      | Priority | Source                          | Effort  | Impact |
| --- | ----------------------------------------------------- | -------- | ------------------------------- | ------- | ------ |
| 13  | Increase v2 coverage from 86.7% to 90%+              | TESTING  | Current coverage report         | 30 min  | High   |
| 14  | Migrate errors_test.go (remove testify)               | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Medium |
| 15  | Migrate config_test.go (remove testify)               | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Medium |
| 16  | Migrate guard_test.go from testify to stdlib          | TESTING  | COMPREHENSIVE_EXECUTION_PLAN    | 20 min  | Medium |
| 17  | Split v2_mixed_flags_test.go (662 lines)              | TESTING  | 2026-03-22_16-58_COMPREHENSIVE  | 25 min  | Low    |
| 18  | Split flags.go (237 lines)                            | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Low    |
| 19  | Split config.go (164 lines)                           | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Low    |
| 20  | Split flags_test.go                                   | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 25 min  | Low    |
| 21  | Split guard_test.go                                   | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 35 min  | Low    |
| 22  | Split config_test.go                                  | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Low    |
| 23  | Split types_test.go                                   | TESTING  | PARETO_EXECUTION_MASTERPLAN     | 15 min  | Low    |

### 📚 Documentation

| ID  | Task Description          | Priority | Source                      | Effort  | Impact |
| --- | ------------------------- | -------- | --------------------------- | ------- | ------ |
| 24  | API Reference             | DOC      | PARETO_EXECUTION_MASTERPLAN | 20 min  | Medium |
| 25  | DI Pattern Example        | DOC      | PARETO_EXECUTION_MASTERPLAN | 15 min  | Medium |
| 26  | Mixed Flags Example       | DOC      | PARETO_EXECUTION_MASTERPLAN | 15 min  | Medium |
| 27  | Update GoDoc comments     | DOC      | API_DESIGN_REVIEW           | 15 min  | Medium |
| 28  | Improve godoc for public  | DOC      | 2026-03-23_15-34_COMPREHENS | 15 min  | Medium |
| 29  | Document DI patterns      | DOC      | 2026-03-22_16-58_COMPREHENS | 15 min  | Medium |
| 30  | Document DI scope pattern | DOC      | PARTS.md                    | 10 min  | Medium |
| 31  | Document error handling   | DOC      | COMPREHENSIVE_STATUS_REPORT | 15 min  | Medium |

### ⚙️ Configuration & Options

| ID  | Task Description                          | Priority | Source              | Effort  | Impact |
| --- | ----------------------------------------- | -------- | ------------------- | ------- | ------ |
| 32  | Remove AddAnyCommand (superseded)         | CONFIG   | API_DESIGN_REVIEW   | 10 min  | High   |
| 33  | Make scope creation lazy                  | CONFIG   | API_DESIGN_REVIEW   | 20 min  | High   |
| 34  | Create compatibility shims if needed      | CONFIG   | API_DESIGN_REVIEW   | 15 min  | Low    |
| 35  | Add Middleware Support                     | CONFIG   | EXECUTION_PLAN      | 30 min  | Medium |
| 36  | Add WithColor option for fang integration  | CONFIG   | api-redesign-v2.1   | 10 min  | Medium |
| 37  | Add more CLI[T] options                   | CONFIG   | 2026-03-22_COMPREH  | 15 min  | Medium |
| 38  | Update `New()` to accept functional opts  | CONFIG   | API_DESIGN_REVIEW   | 20 min  | High   |
| 39  | Update FlagRegistry to FlagRegistry[F]    | CONFIG   | API_DESIGN_REVIEW   | 15 min  | High   |
| 40  | Update NewFlagRegistry to be generic      | CONFIG   | API_DESIGN_REVIEW   | 15 min  | High   |
| 41  | Update ParseFlags to be generic           | CONFIG   | API_DESIGN_REVIEW   | 15 min  | High   |
| 42  | Add `WithDI()` option for opt-in DI       | CONFIG   | API_DESIGN_REVIEW   | 20 min  | High   |
| 43  | Consolidate Scope() methods               | CONFIG   | api-redesign-v2.1   | 15 min  | High   |
| 44  | Create cliConfig internal struct          | CONFIG   | api-redesign-v2.1   | 15 min  | High   |
| 45  | Remove NewWithLong (superseded by option) | CONFIG   | API_DESIGN_REVIEW   | 5 min   | High   |

### 📊 Performance

| ID  | Task Description                          | Priority | Source                       | Effort  | Impact |
| --- | ----------------------------------------- | -------- | ---------------------------- | ------- | ------ |
| 46  | Add benchmark regression detection to CI  | PERF     | 2026-03-28_COMPREHENSIVE     | 10 min  | High   |
| 47  | Add comprehensive performance benchmarks  | PERF     | 2026-03-23_15-34_COMPREHENS  | 30 min  | Medium |

### 🎯 Examples

| ID  | Task Description                      | Priority | Source                         | Effort  | Impact |
| --- | ------------------------------------- | -------- | ------------------------------ | ------- | ------ |
| 48  | Add example with real database conn   | EXAMPLE  | 2026-03-23_15-34_COMPREHENSIVE | 30 min  | High   |
| 49  | Add example with HTTP server          | EXAMPLE  | 2026-03-23_15-34_COMPREHENSIVE | 30 min  | High   |
| 50  | Add lifecycle hook examples           | EXAMPLE  | PARTS.md                       | 20 min  | Medium |
| 51  | Advanced DI Example                   | EXAMPLE  | PARETO_EXECUTION_MASTERPLAN    | 20 min  | Medium |
| 52  | Middleware Example                    | EXAMPLE  | PARETO_EXECUTION_MASTERPLAN    | 20 min  | Medium |
| 53  | Testing Example                       | EXAMPLE  | PARETO_EXECUTION_MASTERPLAN    | 15 min  | Low    |
| 54  | Error Handling Example                | EXAMPLE  | PARETO_EXECUTION_MASTERPLAN    | 20 min  | Medium |
| 55  | Create example for branded IDs        | EXAMPLE  | go-composable-business-types   | 30 min  | High   |

### 🛠️ Linting & Code Quality

| ID  | Task Description                     | Priority | Source                         | Effort  | Impact |
| --- | ------------------------------------ | -------- | ------------------------------ | ------- | ------ |
| 56  | Reduce cyclomatic complexity (cyclop) | LINT    | COMPREHENSIVE_STATUS_REPORT    | 20 min  | Medium |
| 57  | Extract constants (goconst)          | LINT    | COMPREHENSIVE_STATUS_REPORT    | 15 min  | Low    |
| 58  | Split funlen functions               | LINT    | COMPREHENSIVE_STATUS_REPORT    | 15 min  | Low    |
| 59  | Rename BaseError (avoid inherit hint)| LINT    | 2026-03-23_20-51_COMPREHENSIVE | 10 min  | Low    |
| 60  | Audit error message consistency      | LINT    | 2026-03-22_16-58_COMPREHENSIVE | 15 min  | Medium |
| 61  | Fix intrange lint issue              | LINT    | COMPREHENSIVE_STATUS_REPORT    | 10 min  | Low    |
| 62  | Fix gocritic issue                   | LINT    | COMPREHENSIVE_STATUS_REPORT    | 15 min  | Low    |
| 63  | Fix revive style issue               | LINT    | COMPREHENSIVE_STATUS_REPORT    | 15 min  | Low    |
| 64  | Configure exhaustruct for ext structs| LINT    | COMPREHENSIVE_STATUS_REPORT    | 15 min  | Low    |
| 65  | Review gochecknoglobals              | LINT    | COMPREHENSIVE_STATUS_REPORT    | 10 min  | Low    |
| 66  | Review recvcheck                     | LINT    | COMPREHENSIVE_STATUS_REPORT    | 10 min  | Low    |
| 67  | Review unparam                       | LINT    | COMPREHENSIVE_STATUS_REPORT    | 10 min  | Low    |

### 🔍 Testing & Verification

| ID  | Task Description                                 | Priority | Source                         | Effort  | Impact |
| --- | ------------------------------------------------ | -------- | ------------------------------ | ------- | ------ |
| 68  | Add integration tests for command handlers + IDs | TEST     | go-composable-business-types   | 20 min  | High   |
| 69  | Add tests/acceptance/ with user-journey tests    | TEST     | BDD_TESTS_REVIEW               | 30 min  | High   |
| 70  | Add fuzz tests to flags_parse.go                 | TEST     | 2026-03-22_16-58_COMPREHENSIVE | 20 min  | Medium |
| 71  | Add fuzz tests to config_parsing.go              | TEST     | 2026-03-22_16-58_COMPREHENSIVE | 20 min  | Medium |
| 72  | Add fuzz test corpus entries                     | TEST     | 2026-03-28_02-44_COMPREHENSIVE | 15 min  | Low    |
| 73  | Add tests for ID serialization in config         | TEST     | go-composable-business-types   | 15 min  | Medium |
| 74  | Add tests for optional DI                        | TEST     | API_DESIGN_REVIEW              | 15 min  | Medium |
| 75  | Add tests for functional options                 | TEST     | API_DESIGN_REVIEW              | 15 min  | Medium |
| 76  | Write tests for v2.1 additions                   | TEST     | 2026-03-22_14-10_STATUSReport  | 20 min  | High   |

### ⚡ Quick Wins

| ID  | Task Description                  | Priority | Source               | Effort | Impact |
| --- | --------------------------------- | -------- | -------------------- | ------ | ------ |
| 77  | Add short flags for common opts   | QUICK    | CLI_DESIGN_PRINCIPLES| 15 min | Medium |
| 78  | Validate enum values for --log    | QUICK    | CLI_DESIGN_PRINCIPLES| 10 min | Medium |
| 79  | Show defaults in help text        | QUICK    | CLI_DESIGN_PRINCIPLES| 10 min | Medium |

### 🔮 Future / Roadmap (v3.0+)

| ID  | Task Description                          | Priority | Source                        | Effort  | Impact |
| --- | ----------------------------------------- | -------- | ----------------------------- | ------- | ------ |
| 80  | Create v3.0 API design document           | ROADMAP  | 2026-03-22_14-10_STATUSReport | 20 min  | Low    |
| 81  | Create `pkg/cmdguard/v3/` directory        | ROADMAP  | v3.0-major-redesign-plan      | 5 min   | Low    |
| 82  | Implement `errors.go` for v3              | ROADMAP  | v3.0-major-redesign-plan      | 15 min  | Low    |
| 83  | Implement `types.go` for v3               | ROADMAP  | v3.0-major-redesign-plan      | 15 min  | Low    |
| 84  | Implement `cli.go` for v3                 | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 85  | Implement `command.go` for v3             | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 86  | Implement `options.go` for v3             | ROADMAP  | v3.0-major-redesign-plan      | 15 min  | Low    |
| 87  | Implement `flags.go` for v3               | ROADMAP  | v3.0-major-redesign-plan      | 15 min  | Low    |
| 88  | Implement `flags_parse.go` for v3         | ROADMAP  | v3.0-major-redesign-plan      | 25 min  | Low    |
| 89  | Implement `flags_validate.go` for v3      | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 90  | Implement `scope.go` for v3               | ROADMAP  | v3.0-major-redesign-plan      | 15 min  | Low    |
| 91  | Implement `scope_provide.go` for v3       | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 92  | Implement `cli_exec.go` for v3            | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 93  | Write tests for v3 implementation         | ROADMAP  | v3.0-major-redesign-plan      | 30 min  | Low    |
| 94  | Create v3 examples                        | ROADMAP  | v3.0-major-redesign-plan      | 25 min  | Low    |
| 95  | Write MIGRATION_V2_TO_V3.md              | ROADMAP  | v3.0-major-redesign-plan      | 20 min  | Low    |
| 96  | Implement plugin system for validators    | ROADMAP  | FEATURES.md                   | 30 min  | Low    |
| 97  | Custom validation hooks                   | ROADMAP  | ARCHITECTURE_REVIEW           | 20 min  | Low    |
| 98  | Metrics/telemetry integration             | ROADMAP  | ARCHITECTURE_REVIEW           | 30 min  | Low    |
| 99  | Config file support YAML/TOML             | ROADMAP  | ARCHITECTURE_REVIEW           | 25 min  | Low    |
| 100 | Create command groups feature             | ROADMAP  | 2026-03-28_02-44_COMPREHENSIV | 30 min  | Low    |
| 101 | Create `github.com/larsartmann/flagtags`  | ROADMAP  | PARTS.md                      | 30 min  | Low    |
| 102 | Extract flag-related code to standalone   | ROADMAP  | PARTS.md                      | 45 min  | Low    |
| 103 | Add Progress/Spinner Type (bubbles)       | ROADMAP  | EXECUTION_PLAN                | 30 min  | Low    |
| 104 | Add Shell Completion Helpers              | ROADMAP  | EXECUTION_PLAN                | 20 min  | Low    |
| 105 | Add Result[T] type for error handling     | ROADMAP  | EXECUTION_PLAN                | 25 min  | Low    |
| 106 | Add Validated[T] wrapper                  | ROADMAP  | EXECUTION_PLAN                | 25 min  | Low    |
| 107 | Config File Auto-Loading with koanf       | ROADMAP  | EXECUTION_PLAN                | 30 min  | Low    |
| 108 | Environment Variable Binding with env tags| ROADMAP  | EXECUTION_PLAN                | 25 min  | Low    |
| 109 | Replace `internal/config` with koanf      | ROADMAP  | PARTS.md                      | 30 min  | Low    |
| 110 | Replace `internal/logging` with charm/log | ROADMAP  | PARTS.md                      | 25 min  | Low    |

### 📋 Release & CI

| ID  | Task Description                         | Priority | Source                       | Effort  | Impact |
| --- | ---------------------------------------- | -------- | ---------------------------- | ------- | ------ |
| 111 | Create v2.1.0 release tag and notes      | RELEASE  | 2026-03-22_14-10_STATUSRepo  | 15 min  | High   |
| 112 | Set up release automation                | RELEASE  | 2026-03-22_16-58_COMPREHENS  | 20 min  | High   |
| 113 | Add codecov integration                  | RELEASE  | 2026-03-22_16-58_COMPREHENS  | 15 min  | High   |
| 114 | Add pre-commit hooks                     | RELEASE  | COMPREHENSIVE_STATUS_REPORT  | 15 min  | High   |
| 115 | Deprecate v1 API timeline                | RELEASE  | 2026-03-22_16-58_COMPREHENS  | 15 min  | High   |
| 116 | Remove testify/ginkgo completely         | RELEASE  | 2026-03-22_14-10_STATUSRepo  | 20 min  | High   |
| 117 | Update FEATURES.md                       | RELEASE  | api-redesign-v2.1            | 10 min  | Medium |
| 118 | Update API_DESIGN_REVIEW.md              | RELEASE  | api-redesign-v2.1            | 15 min  | Low    |
| 119 | Rename test packages to use `_test`      | RELEASE  | 2026-03-22_14-12_comprehens  | 10 min  | Low    |

### 🔧 API Polish

| ID  | Task Description                        | Priority | Source                       | Effort  | Impact |
| --- | --------------------------------------- | -------- | ---------------------------- | ------- | ------ |
| 120 | Review all `any` usages in package      | API      | API_DESIGN_REVIEW            | 20 min  | High   |
| 121 | Add validation interface abstraction    | API      | COMPREHENSIVE_EXECUTION_PLAN | 20 min  | High   |
| 122 | Add FlagRegistry interface              | API      | COMPREHENSIVE_EXECUTION_PLAN | 20 min  | High   |
| 123 | Add GenerateDocs() method to CLI        | API      | api-redesign-v2.1            | 15 min  | High   |
| 124 | Add GenerateDocsToFile() helper         | API      | api-redesign-v2.1            | 10 min  | High   |
| 125 | Define FlagDoc struct                   | API      | api-redesign-v2.1            | 10 min  | High   |
| 126 | Implement markdown documentation gen    | API      | api-redesign-v2.1            | 20 min  | High   |
| 127 | Create examples/docs-generator/main.go  | API      | api-redesign-v2.1            | 20 min  | High   |
| 128 | Add flag documentation generator        | API      | api-redesign-v2.1            | 20 min  | High   |
| 129 | Add API examples to godoc               | API      | 2026-03-22_16-58_COMPREHENS  | 15 min  | High   |
| 130 | Add flag validation examples            | API      | 2026-03-22_16-58_COMPREHENS  | 15 min  | High   |
| 131 | Add enhanced flag validation enums      | API      | FEATURES.md                  | 20 min  | Medium |
| 132 | Update example code to CLI[T] pattern   | API      | API_DESIGN_REVIEW            | 15 min  | High   |
| 133 | Update examples/typed to CLI[T] pattern | API      | 2026-03-22_14-10_STATUSRepo  | 15 min  | High   |
| 134 | Update examples/basic/main.go           | API      | api-redesign-v2.1            | 15 min  | High   |
| 135 | Review other examples for duplicate     | API      | 2026-03-28_02-44_COMPREHENS  | 15 min  | Low    |

---

## Completed Items (Since Last Audit)

| Task                                                | Completed In         |
| --------------------------------------------------- | -------------------- |
| Remove `F` type param → `CLI[T]`                    | v2.1.0               |
| Rename `GuardedCommand` → `CLI`                     | v2.1.0               |
| Make `AddCommand` accept `Command[T, any]`          | v2.1.0               |
| `Scope()` returns `*Scope`                          | v2.1.0               |
| Remove `ScopeStruct()` from CLI[T]                  | v2.1.0               |
| `Package()` function for samber/do                  | v2.1.0               |
| `WithCLIScope()` option                             | v2.1.0               |
| Deprecation notices on GuardedCommand               | v2.1.0               |
| `SimpleCLI[T]` type alias for compat                | v2.1.0               |
| `Option[T]` type                                    | v2.1.0               |
| Custom types (URL, Email, Port, FilePath, HostPort) | v2.1.0               |
| Flag suggestions (flags_suggest.go)                 | v2.1.0               |
| Quickstart guide (docs/QUICKSTART.md)               | v2.1.0               |
| Migration guide v1→v2                               | v2.1.0               |
| Benchmarks (benchmarks/)                            | v2.1.0               |
| CI/CD pipeline (.github/workflows/ci.yml)           | v2.1.0               |
| CI badge in README                                  | v2.1.0               |
| CONTRIBUTING.md                                     | v2.1.0               |
| CHANGELOG.md                                        | v2.1.0               |
| `t.Parallel()` in v2 tests (167 calls)              | v2.1.0               |
| Remove `//go:fix` directive                         | v2.1.0               |
| Fix usetesting (no os.Setenv in tests)              | v2.1.0               |
| Example tests (basic, typed, di, advanced-flags)    | v2.1.0               |
| `exec.CommandContext` usage                         | v2.1.0               |
| README v2 rewrite                                   | v2.1.0               |

---

## Priority Analysis

### HIGH IMPACT / LOW EFFORT (Quick Wins — do first)

| ID  | Task                              | Effort | Why                                    |
| --- | --------------------------------- | ------ | -------------------------------------- |
| 45  | Remove NewWithLong (deprecated)   | 5 min  | Clean API surface                      |
| 32  | Remove AddAnyCommand              | 10 min | Superseded by CLI[T].AddCommand        |
| 78  | Validate enum values for --log    | 10 min | User-facing correctness                |
| 79  | Show defaults in help text        | 10 min | User experience                        |
| 1   | Fix AddCommand flag parsing bug   | 20 min | Correctness bug                        |
| 111 | Create v2.1.0 release tag         | 15 min | Unblock users                          |
| 38  | Update `New()` to accept options  | 20 min | API consistency with NewCLI            |
| 13  | Increase v2 coverage to 90%+      | 30 min | Currently 86.7%                        |

### HIGH IMPACT / MEDIUM EFFORT (Next sprint)

| ID  | Task                             | Effort  | Why                              |
| --- | -------------------------------- | ------- | -------------------------------- |
| 42  | Add `WithDI()` option            | 20 min  | Opt-in DI, cleaner API           |
| 33  | Make scope creation lazy         | 20 min  | Performance & API simplicity     |
| 132 | Update examples to CLI[T]        | 15 min  | Show recommended API             |
| 112 | Release automation               | 20 min  | Sustainable releases             |
| 116 | Remove testify/ginkgo completely | 20 min  | Dependency cleanup               |
| 48  | Database connection example       | 30 min  | High-value documentation         |

### DEFERRED (v3.0 / Low Priority)

All items in the **Roadmap** section (IDs 80–110) are deferred to v3.0 planning.
No implementation work should start until v2.1.0 is released and stable.

---

## Statistics

| Metric              | Count |
| ------------------- | ----- |
| Total pending tasks | 135   |
| High priority       | 4     |
| Medium priority     | 8     |
| Testing/refactoring | 11    |
| Documentation       | 8     |
| Configuration/opts  | 14    |
| Performance         | 2     |
| Examples            | 8     |
| Linting/quality     | 12    |
| Testing/verification| 9     |
| Quick wins          | 3     |
| Roadmap (v3.0+)     | 31    |
| Release/CI          | 9     |
| API polish          | 16    |
| **Completed**       | **25**|

### Estimated Total Effort

| Category             | Items | Est. Time |
| -------------------- | ----- | --------- |
| Quick wins (<15 min) | ~25   | ~4 hrs    |
| Medium (15-30 min)   | ~70   | ~17 hrs   |
| Large (30+ min)      | ~20   | ~10 hrs   |
| Roadmap              | ~31   | ~8 hrs    |
| **Total**            | **135**| **~39 hrs**|
