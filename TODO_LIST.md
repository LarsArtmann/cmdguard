# TODO List

**Generated:** 2026-04-05
**Purpose:** Actionable items for next 2-4 weeks

## 🔴 High Priority (This Sprint)

- [ ] Fix CLI[T] AddCommand flag parsing bug using cloneAndParseFlags pattern (source: pkg/cmdguard/v2/cli.go:190)
- [ ] Refactor flags_parse_test.go complexity
- [ ] Refactor nestif complexity in flags parsing
- [ ] Fix err113 dynamic error wrapping issues

## 🟡 Medium Priority

- [ ] Improve flag suggestion algorithm
- [ ] Improve error types (more specific error categories)
- [ ] Update README.md with ID usage examples
- [ ] Update AGENTS.md integration patterns
- [ ] Decide on Ginkgo vs stdlib testing and update AGENTS.md
- [ ] Add t.Parallel() to guarded_command_test.go tests
- [ ] Add tests for `initialize` error paths
- [ ] Add tests for `cliToCobraCommand` edge cases
- [ ] Add tests for `cloneAndParseFlags` error paths

## 🧪 Testing & Refactoring

- [ ] Migrate errors_test.go (remove testify)
- [ ] Migrate config_test.go (remove testify)
- [ ] Migrate guard_test.go from testify to stdlib
- [ ] Split guarded_command_test.go (669 lines)
- [ ] Split v2_mixed_flags_test.go (662 lines)
- [ ] Split flags.go (358 lines)
- [ ] Split config.go (352 lines)
- [ ] Split flags_test.go (678 lines)
- [ ] Split guard_test.go (1103 lines)
- [ ] Split config_test.go (452 lines)
- [ ] Split types_test.go (438 lines)

## 📚 Documentation

- [ ] API Reference documentation
- [ ] DI Pattern Example
- [ ] Mixed Flags Example

## ⚙️ Configuration & Options

- [ ] Remove AddAnyCommand (superseded by AddCommand)
- [ ] Make scope creation lazy
- [ ] Create compatibility shims if needed
- [ ] Add Middleware Support
- [ ] Add WithColor option for fang integration
- [ ] Add more CLI[T] options

## 📊 Performance

- [ ] Benchmark: Command Creation
- [ ] Benchmark: Flag Parsing
- [ ] Benchmark: DI Resolution
- [ ] Benchmark Report
- [ ] Add comprehensive performance benchmarks
- [ ] Add benchmark regression detection to CI

## 🎯 Examples

- [ ] Add example/basic unit tests
- [ ] Add example/typed unit tests
- [ ] Add example with real database connection
- [ ] Add example with HTTP server
- [ ] Add lifecycle hook examples
- [ ] Advanced DI Example
- [ ] Middleware Example
- [ ] Testing Example
- [ ] Error Handling Example

## 🛠️ Linting & Code Quality

- [ ] Reduce cyclomatic complexity (cyclop)
- [ ] Extract constants (goconst)
- [ ] Split funlen functions
- [ ] Rename BaseError to avoid inheritance hint
- [ ] Audit error message consistency
- [ ] Add context to exec.Command instances

## 🔍 Testing

- [ ] Run full test suite
- [ ] Run linter
- [ ] Verify build passes
- [ ] Manual testing of examples
- [ ] Add integration tests for command handlers with IDs
- [ ] Add tests/acceptance/ directory with user-journey tests

## ⚡ Quick Wins

- [ ] Add short flags for common options
- [ ] Validate enum values for --log-level
- [ ] Show defaults in help text
