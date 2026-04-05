# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.1.0] - 2026-03-28

### Added

- **CLI[T] New API** - Type-safe CLI builder with generic config and flags
- **BranchingFlowContext** - Tracks command execution path with context propagation
- **Option[T] type** - Functional options pattern for CLI configuration
- **SimpleCLI alias** - Backward compatibility alias for CLI[T]
- **Functional options** - CLIOption[T] for configuring CLI instances
- **WithCLIScope option** - Inject existing DI scope
- **WithLong option** - Set long description via option
- **Comprehensive flow_context tests** - 344 lines, 44 tests
- **CLI[T] integration tests**

### Changed

- **Deprecated GuardedCommand** - Use CLI[T] instead (see migration guide)
- **Updated examples/typed to use SimpleCLI pattern**
- **Improved godoc for public APIs**
- **Updated README.md with links to new docs**

### Fixed

- **flow_context cancel bug** - Self-cancel tracking now works correctly
- **wrapcheck errors** - json.Marshal/Unmarshal now properly annotated
- **go.mod compatibility** - Fixed for Go 1.26
- **go:fix inline directive** - Removed duplicate in type_helpers.go

### Removed

- **AddCommandFunc** - Redundant, use AddCommand instead
- **Dead code packages** - Removed pkg/apperrors, pkg/testutil

### Documentation

- **Migration Guide v1→v2** - docs/MIGRATION_v1_v2.md
- **Quickstart Guide** - docs/QUICKSTART.md
- **README updates** - Links to new documentation

---

## [2.0.0] - 2026-03-22

### Added

- **v2 API** - Type-safe API with dependency injection
- **samber/do/v2 integration** - Full DI support with scopes
- **Typed errors** - Sentinel errors for precise error handling
- **Struct-based flags** - FlagRegistry with struct tag support
- **Comprehensive tests** - 90%+ coverage across packages
- **BDD-style tests** - Behavior documentation via Ginkgo/Gomega
- **Benchmark tests** - Performance validation

### Features

- Generic Command[T, F] type
- GuardedCommand[T, F] for type-safe CLI building
- DI-powered service management
- Flag parsing and validation
- Typo suggestions for unknown flags
- Subcommand support with different flag types
- PreRunE/PostRunE hooks
- Scoped providers for plugin architecture

### Custom Types

- Enum with validation
- Duration with parsing
- LogLevel (debug, info, warn, error)
- LogFormat (text, json)

---

## [0.1.0] - 2026-02-20

### Added

- Initial release of cmdguard v2
- Type-safe CLI construction with generics
- Dependency injection via samber/do/v2
- Flag binding with struct tags
- Full Cobra integration

[Unreleased]: https://github.com/larsartmann/cmdguard/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.1.0
[2.0.0]: https://github.com/larsartmann/cmdguard/releases/tag/v2.0.0
[0.1.0]: https://github.com/larsartmann/cmdguard/releases/tag/v0.1.0
