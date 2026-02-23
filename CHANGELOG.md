# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Split large test files to comply with 250-line policy:
  - `guard_test.go` (1103 lines) → 9 focused test files
  - `flags_test.go` (678 lines) → 4 focused test files
  - `config_test.go` (452 lines) → 5 focused test files
  - `scope_test.go` (446 lines) → 6 focused test files

## [0.1.0] - 2026-02-20

### Added

- Initial release of cmdguard v2
- Type-safe CLI construction with generics
- Dependency injection via samber/do/v2
- Flag binding with struct tags
- Subcommand support with different flag types
- PreRunE/PostRunE hooks
- Scoped providers for plugin architecture
- Comprehensive test suite with 100% coverage
- Benchmark tests for performance validation
- BDD-style tests for behavior documentation

### Features

- **GuardedCommand**: Type-safe CLI builder with compile-time validation
- **FlagRegistry**: Automatic flag registration from struct tags
- **Scope**: Hierarchical DI containers with proper lifecycle management
- **Command Options**: Functional options pattern for command configuration
- **Custom Types**: Enum, Duration, LogLevel, LogFormat with JSON marshaling
- **Error Handling**: Structured errors with suggestions for invalid flags
- **Integration**: Full cobra integration with standard patterns

[Unreleased]: https://github.com/larsartmann/cmdguard/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v0.1.0
