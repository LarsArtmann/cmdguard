# Comprehensive TODO Table View

## Pending Tasks Extracted from TODO_LIST.md

| ID  | Task Description                                                        | Priority | Source                                                                                          | Estimated Effort | Impact |
| --- | ----------------------------------------------------------------------- | -------- | ----------------------------------------------------------------------------------------------- | ---------------- | ------ |
| 1   | Fix CLI[T] AddCommand flag parsing bug using cloneAndParseFlags pattern | MEDIUM   | pkg/cmdguard/v2/cli.go:190                                                                      | 20 min           | High   |
| 2   | Update `New()` to accept functional options                             | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 3   | Update `Scope()` to return `*Scope`                                     | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 4   | Update `FlagRegistry` to `FlagRegistry[F]` struct                       | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 5   | Update `NewFlagRegistry` to be generic                                  | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 6   | Update `ParseFlags` to be generic                                       | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 7   | Fix usetesting - replace os.Setenv with t.Setenv in tests               | MEDIUM   | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 10 min           | Medium |
| 8   | Update README.md                                                        | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 9   | Update AGENTS.md                                                        | MEDIUM   | API_DESIGN_REVIEW.md, PARETO_EXECUTION_MASTERPLAN.md                                            | 15 min           | High   |
| 10  | Update GoDoc comments                                                   | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | Medium |
| 11  | Update example code                                                     | MEDIUM   | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 12  | Update Examples to Use CLI[T] API                                       | MEDIUM   | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md                                                  | 15 min           | High   |
| 13  | Update examples/typed to use SimpleCLI pattern                          | MEDIUM   | 2026-03-22_14-10_STATUSReport-v2.1.md                                                           | 15 min           | High   |
| 14  | Improve godoc for public APIs                                           | MEDIUM   | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 15  | Fix intrange - use range for integers pattern                           | MEDIUM   | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 10 min           | Low    |
| 16  | Fix gocritic issue                                                      | MEDIUM   | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 17  | Fix revive style issue                                                  | MEDIUM   | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 18  | Remove duplicate //go:fix inline directive                              | MEDIUM   | type_helpers.go:9                                                                               | 5 min            | Low    |
| 19  | Update examples/basic/main.go                                           | MEDIUM   | api-redesign-v2.1.md                                                                            | 15 min           | High   |
| 20  | Update examples/typed/main.go                                           | MEDIUM   | api-redesign-v2.1.md                                                                            | 15 min           | High   |
| 21  | Refactor flags_parse_test.go complexity                                 | MEDIUM   | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | Medium |
| 22  | Refactor nestif                                                         | MEDIUM   | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 23  | Fix err113 dynamic error wrapping issues                                | MEDIUM   | docs/status/2026-03-22_14-12_comprehensive-status.md                                            | 15 min           | Medium |
| 24  | Update FEATURES.md                                                      | MEDIUM   | api-redesign-v2.1.md                                                                            | 10 min           | Medium |
| 25  | Update API_DESIGN_REVIEW.md                                             | MEDIUM   | api-redesign-v2.1.md                                                                            | 15 min           | Low    |
| 26  | Improve flag suggestion algorithm                                       | MEDIUM   | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 27  | Improve error types                                                     | MEDIUM   | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 28  | Update README.md with ID usage examples                                 | LOW      | docs/planning/go-composable-business-types-usage.md                                             | 15 min           | Medium |
| 29  | Update AGENTS.md integration patterns                                   | LOW      | docs/planning/go-composable-business-types-usage.md                                             | 15 min           | Medium |
| 30  | Decide on Ginkgo vs stdlib testing and update AGENTS.md                 | LOW      | BDD_TESTS_REVIEW.md                                                                             | 10 min           | Low    |
| 31  | Review all `any` usages in package                                      | LOW      | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 32  | Document DI patterns                                                    | LOW      | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 33  | Document DI scope pattern in docs/                                      | LOW      | PARTS.md                                                                                        | 10 min           | Medium |
| 34  | Document error handling strategy                                        | LOW      | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Medium |
| 35  | Review gochecknoglobals                                                 | LOW      | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 10 min           | Low    |
| 36  | Review recvcheck                                                        | LOW      | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 10 min           | Low    |
| 37  | Review unparam                                                          | LOW      | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 10 min           | Low    |
| 38  | Create v3.0 API design document                                         | LOW      | 2026-03-22_14-10_STATUSReport-v2.1.md                                                           | 20 min           | Low    |
| 39  | Review other examples for duplicate code                                | LOW      | 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Low    |
| 40  | Remove `F` type parameter from `GuardedCommand[T, F]` → `CLI[T]`        | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 30 min           | High   |
| 41  | Rename `GuardedCommand` → `CLI`                                         | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 42  | Make `AddCommand` accept `Command[T, any]`                              | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 43  | Remove `AddAnyCommand`                                                  | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 10 min           | High   |
| 44  | Remove `AddCommandFunc`                                                 | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 5 min            | Low    |
| 45  | Add `WithDI()` option for opt-in DI                                     | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 46  | Make scope creation lazy                                                | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 20 min           | High   |
| 47  | Remove `ScopeStruct()` method                                           | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 10 min           | High   |
| 48  | Create `pkg/cmdguard/v3/` directory                                     | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 5 min            | Low    |
| 49  | Implement `errors.go` for v3                                            | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 15 min           | Low    |
| 50  | Implement `types.go` for v3                                             | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 15 min           | Low    |
| 51  | Implement `cli.go` for v3                                               | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 52  | Implement `command.go` for v3                                           | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 53  | Implement `options.go` for v3                                           | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 15 min           | Low    |
| 54  | Implement `flags.go` for v3                                             | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 15 min           | Low    |
| 55  | Implement `scope.go` for v3                                             | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 15 min           | Low    |
| 56  | Add tests for optional DI                                               | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | Medium |
| 57  | Add tests for functional options                                        | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | Medium |
| 58  | Verify 90%+ coverage maintained                                         | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 10 min           | High   |
| 59  | Increase v2 coverage from 81.2% to 90%+                                 | UNKNOWN  | docs/status/2026-03-28_00-05_COMPREHENSIVE_STATUS_REPORT.md                                     | 30 min           | High   |
| 60  | Add t.Parallel() to guarded_command_test.go tests                       | UNKNOWN  | pkg/cmdguard/guarded_command_test.go                                                            | 10 min           | Medium |
| 61  | Add tests for `initialize` error paths                                  | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 62  | Add tests for `cliToCobraCommand` edge cases                            | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 63  | Add tests for `cloneAndParseFlags` error paths                          | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 64  | Write tests for v2.1 additions                                          | UNKNOWN  | 2026-03-22_14-10_STATUSReport-v2.1.md                                                           | 20 min           | High   |
| 65  | Add fuzz test corpus in testdata/fuzz/ directories                      | UNKNOWN  | docs/status/2026-03-28_00-05_COMPREHENSIVE_STATUS_REPORT.md                                     | 20 min           | Medium |
| 66  | Add CLI[T] integration tests                                            | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | High   |
| 67  | Migrate errors_test.go (remove testify)                                 | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 68  | Migrate types_test.go (remove testify)                                  | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 69  | Migrate command_test.go (remove testify)                                | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 70  | Migrate config_test.go (remove testify)                                 | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 71  | Migrate guard_test.go from testify to stdlib                            | UNKNOWN  | docs/planning/2026-02-20_COMPREHENSIVE_EXECUTION_PLAN.md                                        | 20 min           | Medium |
| 72  | Split guarded_command_test.go (669 lines)                               | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 25 min           | Low    |
| 73  | Split v2_mixed_flags_test.go (662 lines)                                | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 25 min           | Low    |
| 74  | Split flags.go (358 lines)                                              | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 75  | Split config.go (352 lines)                                             | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 76  | Split flags_test.go (678 lines)                                         | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 25 min           | Low    |
| 77  | Split guard_test.go (1103 lines)                                        | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 35 min           | Low    |
| 78  | Split config_test.go (452 lines)                                        | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 79  | Split types_test.go (438 lines)                                         | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 80  | README v2 Rewrite                                                       | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 45 min           | High   |
| 81  | Quickstart Example                                                      | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 20 min           | High   |
| 82  | Migration Guide v1 → v2                                                 | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md, 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                 | 20 min           | High   |
| 83  | API Reference                                                           | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 20 min           | Medium |
| 84  | DI Pattern Example                                                      | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 85  | Mixed Flags Example                                                     | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Medium |
| 86  | MIGRATION.md guide                                                      | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 87  | Add `Package()` function for samber/do integration                      | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 88  | Add `WithScope()` option to inject existing scope                       | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 89  | Add deprecation type aliases for backward compatibility                 | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | High   |
| 90  | Add `Deprecated:` comments to removed functions                         | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 10 min           | Medium |
| 91  | Create compatibility shims if needed                                    | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 15 min           | Low    |
| 92  | Add Deprecation Path for GuardedCommand                                 | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md                                                  | 10 min           | Medium |
| 93  | Add Middleware Support                                                  | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md                                                  | 30 min           | Medium |
| 94  | Add Progress/Spinner Type using charmbracelet/bubbles                   | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 30 min           | Low    |
| 95  | Add Shell Completion Helpers                                            | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 20 min           | Low    |
| 96  | Add Result[T] type for error handling                                   | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 25 min           | Medium |
| 97  | Add Validated[T] wrapper with validation functions                      | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 25 min           | Medium |
| 98  | Config File Auto-Loading integration with koanf                         | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 30 min           | High   |
| 99  | Environment Variable Binding with env struct tags                       | UNKNOWN  | docs/status/2026-03-28_00-15_EXECUTION_PLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md | 25 min           | High   |
| 100 | Replace `internal/config` with koanf                                    | UNKNOWN  | PARTS.md                                                                                        | 30 min           | High   |
| 101 | Replace `internal/logging` with charmbracelet/log                       | UNKNOWN  | PARTS.md                                                                                        | 25 min           | High   |
| 102 | Add short flags for common options                                      | UNKNOWN  | CLI_DESIGN_PRINCIPLES.md                                                                        | 15 min           | Medium |
| 103 | Validate enum values for --log-level                                    | UNKNOWN  | CLI_DESIGN_PRINCIPLES.md                                                                        | 10 min           | Medium |
| 104 | Show defaults in help text                                              | UNKNOWN  | CLI_DESIGN_PRINCIPLES.md                                                                        | 10 min           | Medium |
| 105 | Add flag suggestions on unknown flag errors                             | UNKNOWN  | CLI_DESIGN_PRINCIPLES.md, 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                       | 15 min           | High   |
| 106 | Create Option[T] type implementation                                    | UNKNOWN  | 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 107 | Benchmark: Command Creation                                             | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 108 | Benchmark: Flag Parsing                                                 | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 109 | Benchmark: DI Resolution                                                | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 110 | Benchmark Report                                                        | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 111 | Add comprehensive performance benchmarks                                | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 30 min           | Medium |
| 112 | Add benchmark regression detection to CI                                | UNKNOWN  | docs/status/2026-03-28_00-05_COMPREHENSIVE_STATUS_REPORT.md                                     | 10 min           | High   |
| 113 | Implement `flags_parse.go` for v3                                       | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 25 min           | Low    |
| 114 | Implement `flags_validate.go` for v3                                    | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 115 | Implement `scope_provide.go` for v3                                     | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 116 | Implement `cli_exec.go` for v3                                          | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 117 | Write tests for v3 implementation                                       | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 30 min           | High   |
| 118 | Create v3 examples                                                      | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 25 min           | Medium |
| 119 | Write MIGRATION_V2_TO_V3.md                                             | UNKNOWN  | v3.0-major-redesign-plan.md                                                                     | 20 min           | Low    |
| 120 | Advanced DI Example                                                     | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 20 min           | Medium |
| 121 | Middleware Example                                                      | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md, 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                 | 20 min           | Medium |
| 122 | Testing Example                                                         | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 15 min           | Low    |
| 123 | Error Handling Example                                                  | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 20 min           | Medium |
| 124 | Add example/basic unit tests                                            | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 125 | Add example/typed unit tests                                            | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 126 | Add example with real database connection                               | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 30 min           | High   |
| 127 | Add example with HTTP server                                            | UNKNOWN  | 2026-03-23_15-34_COMPREHENSIVE_STATUS_REPORT.md                                                 | 30 min           | High   |
| 128 | Create example application for branded IDs                              | UNKNOWN  | docs/planning/go-composable-business-types-usage.md                                             | 30 min           | High   |
| 129 | Add lifecycle hook examples                                             | UNKNOWN  | PARTS.md                                                                                        | 20 min           | Medium |
| 130 | Create examples/docs-generator/main.go                                  | UNKNOWN  | api-redesign-v2.1.md                                                                            | 20 min           | High   |
| 131 | Add flag documentation generator                                        | UNKNOWN  | api-redesign-v2.1.md                                                                            | 20 min           | High   |
| 132 | Define FlagDoc struct                                                   | UNKNOWN  | api-redesign-v2.1.md                                                                            | 10 min           | High   |
| 133 | Add GenerateDocs() method to CLI                                        | UNKNOWN  | api-redesign-v2.1.md                                                                            | 15 min           | High   |
| 134 | Implement markdown documentation generator                              | UNKNOWN  | api-redesign-v2.1.md                                                                            | 20 min           | High   |
| 135 | Add GenerateDocsToFile() helper                                         | UNKNOWN  | api-redesign-v2.1.md                                                                            | 10 min           | High   |
| 136 | Add API examples to godoc                                               | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | High   |
| 137 | Add flag validation examples                                            | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | High   |
| 138 | Remove `NewWithLong()` (superseded by `WithLong()` option)              | UNKNOWN  | API_DESIGN_REVIEW.md                                                                            | 5 min            | High   |
| 139 | Create cliConfig internal struct                                        | UNKNOWN  | api-redesign-v2.1.md                                                                            | 15 min           | High   |
| 140 | Create type alias CLI = GuardedCommand                                  | UNKNOWN  | api-redesign-v2.1.md                                                                            | 5 min            | High   |
| 141 | Add deprecation notice to GuardedCommand                                | UNKNOWN  | api-redesign-v2.1.md                                                                            | 5 min            | High   |
| 142 | Consolidate Scope() methods                                             | UNKNOWN  | api-redesign-v2.1.md                                                                            | 15 min           | High   |
| 143 | Reduce cyclomatic complexity (cyclop)                                   | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 20 min           | Medium |
| 144 | Extract constants (goconst)                                             | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 145 | Configure exhaustruct for external structs                              | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 146 | Split funlen functions                                                  | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | Low    |
| 147 | Rename BaseError to avoid inheritance hint                              | UNKNOWN  | docs/status/2026-03-23_20-51_COMPREHENSIVE_STATUS_REPORT.md                                     | 10 min           | Low    |
| 148 | Rename test packages to use \_test suffix                               | UNKNOWN  | docs/status/2026-03-22_14-12_comprehensive-status.md                                            | 10 min           | Low    |
| 149 | Create `github.com/larsartmann/flagtags` repository                     | UNKNOWN  | PARTS.md                                                                                        | 30 min           | Low    |
| 150 | Extract flag-related code to standalone library                         | UNKNOWN  | PARTS.md                                                                                        | 45 min           | Low    |
| 151 | Add more custom types (URL, Email, Port, FilePath)                      | UNKNOWN  | 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | Medium |
| 152 | Create command groups feature                                           | UNKNOWN  | 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                                                 | 30 min           | Medium |
| 153 | Implement plugin system for custom validators                           | UNKNOWN  | FEATURES.md                                                                                     | 30 min           | High   |
| 154 | Add enhanced flag validation enums                                      | UNKNOWN  | FEATURES.md                                                                                     | 20 n             | Medium |
| 155 | Custom validation hooks                                                 | UNKNOWN  | ARCHITECTURE_REVIEW.md                                                                          | 20 min           | High   |
| 156 | Metrics/telemetry integration                                           | UNKNOWN  | ARCHITECTURE_REVIEW.md                                                                          | 30 min           | High   |
| 157 | Config file support YAML/TOML                                           | UNKNOWN  | ARCHITECTURE_REVIEW.md                                                                          | 25 min           | High   |
| 158 | Changelog v2.0                                                          | UNKNOWN  | PARETO_EXECUTION_MASTERPLAN.md                                                                  | 10 min           | High   |
| 159 | Create v2.1.0 release tag and release notes                             | UNKNOWN  | 2026-03-22_14-10_STATUSReport-v2.1.md                                                           | 15 min           | High   |
| 160 | Add changelog                                                           | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 10 min           | High   |
| 161 | Set up release automation                                               | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md, FEATURES.md                                    | 20 min           | High   |
| 162 | Add GitHub Actions workflow                                             | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | High   |
| 163 | Add codecov integration                                                 | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | High   |
| 164 | Add badge to README                                                     | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 5 min            | High   |
| 165 | Set up CI/CD pipeline                                                   | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 20 min           | High   |
| 166 | Add pre-commit hooks                                                    | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md                                                                  | 15 min           | High   |
| 167 | Create contribution guide                                               | UNKNOWN  | COMPREHENSIVE_STATUS_REPORT.md, PARETO_EXECUTION_MASTERPLAN.md                                  | 20 min           | High   |
| 168 | Deprecate v1 API timeline                                               | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | High   |
| 169 | Remove testify/ginkgo completely                                        | UNKNOWN  | 2026-03-22_14-10_STATUSReport-v2.1.md                                                           | 20 min           | High   |
| 170 | Run full test suite                                                     | UNKNOWN  | api-redesign-v2.1.md                                                                            | 30 min           | High   |
| 171 | Run linter                                                              | UNKNOWN  | api-redesign-v2.1.md                                                                            | 10 min           | High   |
| 172 | Verify build passes                                                     | UNKNOWN  | api-redesign-v2.1.md                                                                            | 10 min           | High   |
| 173 | Manual testing of examples                                              | UNKNOWN  | api-redesign-v2.1.md                                                                            | 20 min           | High   |
| 174 | Add WithColor option for fang integration                               | UNKNOWN  | api-redesign-v2.1.md                                                                            | 10 min           | Medium |
| 175 | Audit error message consistency                                         | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 176 | Add more CLI[T] options                                                 | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Medium |
| 177 | Add fuzz tests to flags_parse.go                                        | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | Medium |
| 178 | Add fuzz tests to config_parsing.go                                     | UNKNOWN  | 2026-03-22_16-58_COMPREHENSIVE_STATUS_REPORT.md                                                 | 20 min           | Medium |
| 179 | Add context to exec.Command instances                                   | UNKNOWN  | docs/status/2026-03-23_20-51_COMPREHENSIVE_STATUS_REPORT.md                                     | 15 min           | Medium |
| 180 | Fuzz test corpus entries                                                | UNKNOWN  | 2026-03-28_02-44_COMPREHENSIVE_STATUS_REPORT.md                                                 | 15 min           | Low    |
| 181 | Add tests for ID serialization in config                                | UNKNOWN  | docs/planning/go-composable-business-types-usage.md                                             | 15 min           | Medium |
| 182 | Add integration tests for command handlers with IDs                     | UNKNOWN  | docs/planning/go-composable-business-types-usage.md                                             | 20 min           | High   |
| 183 | Add tests/acceptance/ directory with user-journey tests                 | UNKNOWN  | BDD_TESTS_REVIEW.md                                                                             | 30 min           | High   |
| 184 | Add validation interface abstraction                                    | UNKNOWN  | docs/planning/2026-02-20_COMPREHENSIVE_EXECUTION_PLAN.md                                        | 20 min           | High   |
| 185 | Add FlagRegistry interface                                              | UNKNOWN  | docs/planning/2026-02-20_COMPREHENSIVE_EXECUTION_PLAN.md                                        | 20 min           | High   |

## Analysis of API_DESIGN_REVIEW.md Implementation Status

Based on my review of API_DESIGN_REVIEW.md, here's what has been implemented vs what remains:

### ✅ IMPLEMENTED (Based on TODO_LIST.md completed items):

- Migration Guide v1→v2 (docs/MIGRATION_v1_v2.md) - DONE
- Quickstart Guide (docs/QUICKSTART.md) - DONE
- Add deprecation notices to GuardedCommand API - DONE
- Bump version to 2.1.0 - DONE
- Fix go.mod for go 1.26 compatibility - DONE
- Verify all linter issues resolved - DONE
- Verify all fuzz tests pass - DONE
- Remove dead code packages (pkg/apperrors, pkg/testutil) - DONE
- Remove duplicate //go:fix inline directive in type_helpers.go - DONE
- Increase v2 coverage to 90.2% (added MarshalText tests) - DONE
- Update README.md with links to new docs - DONE

### ⏳ IN PROGRESS / PENDING (From API_DESIGN_REVIEW.md recommendations):

1. **Type Parameter Simplification** - Remove `F` from `GuardedCommand[T, F]` → `CLI[T]` (Items 40, 41, 42, 43, 44)
2. **DI Optional** - Add `WithDI()` option, make scope creation lazy (Items 45, 46, 47)
3. **Type Safety** - Fix `any` in FlagRegistry (Items 5, 6)
4. **API Consistency** - Update `New()` to accept functional options (Item 2)
5. **Scope Access** - Update `Scope()` to return `*Scope`, remove `ScopeStruct()` (Items 3, 47)
6. **Redundant Methods** - Remove `AddCommandFunc` (Item 44)
7. **samber/do Integration** - Add `Package()` function (Item 87)
8. **Documentation Updates** - Update README.md, AGENTS.md, GoDoc comments, example code (Items 8, 9, 10, 11)

### 📊 PRIORITY SORTING (by Impact/Effort Ratio):

**HIGH IMPACT, LOW EFFORT (Quick Wins - ≤15 min):**

- Remove duplicate //go:fix inline directive (5 min) - Item 18
- Add t.Parallel() to guarded_command_test.go tests (10 min) - Item 60
- Add benchmark regression detection to CI (10 min) - Item 112
- Add badge to README (5 min) - Item 164
- Add changelog (10 min) - Item 160
- Remove `NewWithLong()` (5 min) - Item 138
- Create type alias CLI = GuardedCommand (5 min) - Item 140
- Add deprecation notice to GuardedCommand (5 min) - Item 141
- Remove `AddCommandFunc` (5 min) - Item 44

**HIGH IMPACT, MEDIUM EFFORT (15-30 min):**

- Update `New()` to accept functional options (20 min) - Item 2
- Update `Scope()` to return `*Scope` (15 min) - Item 3
- Update `FlagRegistry` to `FlagRegistry[F]` struct (15 min) - Item 4
- Update `NewFlagRegistry` to be generic (15 min) - Item 5
- Update `ParseFlags` to be generic (15 min) - Item 6
- Make `AddCommand` accept `Command[T, any]` (20 min) - Item 42
- Remove `F` type parameter from `GuardedCommand[T, F]` → `CLI[T]` (30 min) - Item 40
- Rename `GuardedCommand` → `CLI` (20 min) - Item 41
- Remove `AddAnyCommand` (10 min) - Item 43
- Add `WithDI()` option for opt-in DI (20 min) - Item 45
- Make scope creation lazy (20 min) - Item 46
- Remove `ScopeStruct()` method (10 min) - Item 47
- Add `Package()` function for samber/do integration (15 min) - Item 87
- Add `WithScope()` option to inject existing scope (15 min) - Item 88
- Add deprecation type aliases for backward compatibility (15 min) - Item 89
- Add `Deprecated:` comments to removed functions (10 min) - Item 90
- Add Deprecation Path for GuardedCommand (10 min) - Item 92
- Consolidate Scope() methods (15 min) - Item 142
- Create cliConfig internal struct (15 min) - Item 139
- Update FEATURES.md (10 min) - Item 24
- Update README.md (15 min) - Item 8
- Update AGENTS.md (15 min) - Item 9
- Update GoDoc comments (15 min) - Item 10
- Update example code (15 min) - Item 11
- Update Examples to Use CLI[T] API (15 min) - Item 12
- Update examples/typed to use SimpleCLI pattern (15 min) - Item 13
- Update examples/basic/main.go (15 min) - Item 19
- Update examples/typed/main.go (15 min) - Item 20
- Fix err113 dynamic error wrapping issues (15 min) - Item 23
- Add flag suggestions on unknown flag errors (15 min) - Item 105
- Add API examples to godoc (15 min) - Item 136
- Add flag validation examples (15 min) - Item 137
- Add GenerateDocs() method to CLI (15 min) - Item 133
- Add GenerateDocsToFile() helper (10 min) - Item 135
- Add tests for v2.1 additions (20 min) - Item 64
- Add CLI[T] integration tests (20 min) - Item 66
- Verify 90%+ coverage maintained (10 min) - Item 58
- Increase v2 coverage from 81.2% to 90%+ (30 min) - Item 59
- Add tests for optional DI (15 min) - Item 56
- Add tests for functional options (15 min) - Item 57
- Migrate guard_test.go from testify to stdlib (20 min) - Item 71
- Write tests for v3 implementation (30 min) - Item 117
- Create v3 examples (25 min) - Item 118
- Write MIGRATION_V2_TO_V3.md (20 min) - Item 119
- Add example with real database connection (30 min) - Item 126
- Add example with HTTP server (30 min) - Item 127
- Create example application for branded IDs (30 min) - Item 128
- Add integration tests for command handlers with IDs (20 min) - Item 182
- Add tests/acceptance/ directory with user-journey tests (30 min) - Item 183
- Add validation interface abstraction (20 min) - Item 184
- Add FlagRegistry interface (20 min) - Item 185
- Config File Auto-Loading integration with koanf (30 min) - Item 98
- Environment Variable Binding with env struct tags (25 min) - Item 99
- Replace `internal/config` with koanf (30 min) - Item 100
- Replace `internal/logging` with charmbracelet/log (25 min) - Item 101
- Implement plugin system for custom validators (30 min) - Item 153
- Custom validation hooks (20 min) - Item 155
- Metrics/telemetry integration (30 min) - Item 156
- Config file support YAML/TOML (25 min) - Item 157
- Set up release automation (20 min) - Item 161
- Add GitHub Actions workflow (20 min) - Item 162
- Add codecov integration (15 min) - Item 163
- Set up CI/CD pipeline (20 min) - Item 165
- Add pre-commit hooks (15 min) - Item 166
- Create contribution guide (20 min) - Item 167
- Deprecate v1 API timeline (15 min) - Item 168
- Remove testify/ginkgo completely (20 min) - Item 169
- Run full test suite (30 min) - Item 170
- Run linter (10 min) - Item 171
- Verify build passes (10 min) - Item 172
- Manual testing of examples (20 min) - Item 173

**MEDIUM IMPACT, LOW EFFORT (≤15 min):**

- Fix usetesting - replace os.Setenv with t.Setenv in tests (10 min) - Item 7
- Improve godoc for public APIs (15 min) - Item 14
- Fix intrange - use range for integers pattern (10 min) - Item 15
- Fix gocritic issue (15 min) - Item 16
- Fix revive style issue (15 min) - Item 17
- Refactor nestif (15 min) - Item 22
- Improve flag suggestion algorithm (15 min) - Item 26
- Improve error types (15 min) - Item 27
- Update API_DESIGN_REVIEW.md (15 min) - Item 25
- Decide on Ginkgo vs stdlib testing and update AGENTS.md (10 min) - Item 30
- Review all `any` usages in package (20 min) - Item 31
- Document DI patterns (15 min) - Item 32
- Document DI scope pattern in docs/ (10 min) - Item 33
- Document error handling strategy (15 min) - Item 34
- Review gochecknoglobals (10 min) - Item 35
- Review recvcheck (10 min) - Item 36
- Review unparam (10 min) - Item 37
- Create v3.0 API design document (20 min) - Item 38
- Review other examples for duplicate code (15 min) - Item 39
- Add tests for `initialize` error paths (15 min) - Item 61
- Add tests for `cliToCobraCommand` edge cases (15 min) - Item 62
- Add tests for `cloneAndParseFlags` error paths (15 min) - Item 63
- Add fuzz test corpus in testdata/fuzz/ directories (20 min) - Item 65
- Migrate errors_test.go (remove testify) (15 min) - Item 67
- Migrate types_test.go (remove testify) (15 min) - Item 68
- Migrate command_test.go (remove testify) (15 min) - Item 69
- Migrate config_test.go (remove testify) (15 min) - Item 70
- Split guarded_command_test.go (669 lines) (25 min) - Item 72
- Split v2_mixed_flags_test.go (662 lines) (25 min) - Item 73
- Split flags.go (358 lines) (15 min) - Item 74
- Split config.go (352 lines) (15 min) - Item 75
- Split flags_test.go (678 lines) (25 min) - Item 76
- Split guard_test.go (1103 lines) (35 min) - Item 77
- Split config_test.go (452 lines) (15 min) - Item 78
- Split types_test.go (438 lines) (15 min) - Item 79
- API Reference (20 min) - Item 83
- DI Pattern Example (15 min) - Item 84
- Mixed Flags Example (15 min) - Item 85
- MIGRATION.md guide (15 min) - Item 86
- Quickstart Example (20 min) - Item 81
- Migration Guide v1 → v2 (20 min) - Item 82
- Advanced DI Example (20 min) - Item 120
- Middleware Example (20 min) - Item 121
- Testing Example (15 min) - Item 122
- Error Handling Example (20 min) - Item 123
- Add example/basic unit tests (15 min) - Item 124
- Add example/typed unit tests (15 min) - Item 125
- Add lifecycle hook examples (20 min) - Item 129
- Create examples/docs-generator/main.go (20 min) - Item 130
- Add flag documentation generator (20 min) - Item 131
- Define FlagDoc struct (10 min) - Item 132
- Implement markdown documentation generator (20 min) - Item 134
- Benchmark: Command Creation (15 min) - Item 107
- Benchmark: Flag Parsing (15 min) - Item 108
- Benchmark: DI Resolution (15 min) - Item 109
- Benchmark Report (15 min) - Item 110
- Add comprehensive performance benchmarks (30 min) - Item 111
- Implement `flags_parse.go` for v3 (25 min) - Item 113
- Implement `flags_validate.go` for v3 (20 min) - Item 114
- Implement `scope_provide.go` for v3 (20 min) - Item 115
- Implement `cli_exec.go` for v3 (20 min) - Item 116
- Create Option[T] type implementation (15 min) - Item 106
- Benchmark: Command Creation (15 min) - Item 107
- Benchmark: Flag Parsing (15 min) - Item 108
- Benchmark: DI Resolution (15 min) - Item 109
- Benchmark Report (15 min) - Item 110
- Add more custom types (URL, Email, Port, FilePath) (20 min) - Item 151
- Add enhanced flag validation enums (20 min) - Item 154
- Add WithColor option for fang integration (10 min) - Item 174
- Audit error message consistency (15 min) - Item 175
- Add more CLI[T] options (15 min) - Item 176
- Add fuzz tests to flags_parse.go (20 min) - Item 177
- Add fuzz tests to config_parsing.go (20 min) - Item 178
- Add context to exec.Command instances (15 min) - Item 179
- Fuzz test corpus entries (15 min) - Item 180
- Add tests for ID serialization in config (15 min) - Item 181
- Create `github.com/larsartmann/flagtags` repository (30 min) - Item 149
- Extract flag-related code to standalone library (45 min) - Item 150
- Create command groups feature (30 min) - Item 152
- Reduce cyclomatic complexity (cyclop) (20 min) - Item 143
- Extract constants (goconst) (15 min) - Item 144
- Configure exhaustruct for external structs (15 min) - Item 145
- Split funlen functions (15 min) - Item 146
- Rename BaseError to avoid inheritance hint (10 min) - Item 147
- Rename test packages to use \_test suffix (10 min) - Item 148
