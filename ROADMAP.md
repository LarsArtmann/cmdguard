# ROADMAP

**Updated:** 2026-06-03
**Purpose:** Aspirational items with no concrete timeline

---

## Completed (v2.2–v2.4)

- [x] Add GitHub Actions workflow
- [x] Add badge to README
- [x] Add more custom types: URL, Email, Port, FilePath, HostPort
- [x] Add middleware support (`TimingMiddleware`, `RecoveryMiddleware`, `SpinnerMiddleware`, `TelemetryMiddleware`)
- [x] Create command groups feature (`WithGroup`, `WithGroupID`)
- [x] Environment Variable Binding with `env` struct tags
- [x] Config file support JSON/YAML/TOML
- [x] Shell completion helpers (`WithCompletion`)
- [x] Man page generation (`GenerateManPageCommand`)
- [x] `BranchingFlowContext` for command path tracking
- [x] `EditInEditor` for `$EDITOR` integration
- [x] Counting flags (`count:"true"`)
- [x] Positional args validators (`WithExactArgs`, `WithRangeArgs`, etc.)
- [x] `ExitCoder` / `NewExitError` for custom exit codes
- [x] `WithStrictValidation` / `WithDraconianValidation`
- [x] `WithConfigValidation` for post-parse config validation
- [x] Consumer test harness (`pkg/testutil`)
- [x] Cobra migration guide (`docs/MIGRATION_FROM_COBRA.md`)
- [x] Framework comparison (`docs/COMPARISON.md`)
- [x] Interactive prompts via `huh` (`WithPromptOnMissing`, `prompt` tag)
- [x] Markdown help rendering via `glamour` (`WithGlamourHelp`)
- [x] Terminal spinner (`SpinnerMiddleware`)
- [x] OpenTelemetry integration (`WithTelemetry`, `TelemetryMiddleware`)
- [x] `--no-color` flag + `NO_COLOR` support
- [x] Rich output with 12+ formats via `go-output`
- [x] Review all `any` usages in package
- [x] Remove string-based `BranchWithTimeout`/`BranchWithDeadline`
- [x] Remove deprecated `FlowContextAccessor` API

---

## v3.0 Major Redesign

- [ ] Create v3.0 API design document
- [ ] Create `pkg/cmdguard/v3/` directory
- [ ] Implement core types, CLI, commands, flags, scope, options for v3
- [ ] Write tests for v3 implementation
- [ ] Create v3 examples
- [ ] Write `MIGRATION_V2_TO_V3.md`

### v3.0 API-Breaking Cleanup

- [ ] Make `NoFlags` a distinct named type (not `type NoFlags = struct{}`)
- [ ] Rename `Get[T]`/`MustGet[T]` to more specific names
- [ ] Make `RegisterInScope` generic instead of `...any`
- [ ] Remove or redesign `Package()` for error-safe DI integration
- [ ] Make `NoFlags` a distinct named type (not type alias)

---

## Advanced Types

- [ ] Add `Result[T]` type for error handling
- [ ] Add `Validated[T]` wrapper with validation functions
- [ ] Create example application for branded IDs

---

## Documentation Generation

- [ ] Create `examples/docs-generator/main.go`
- [ ] Define `FlagDoc` struct
- [ ] Add `GenerateDocs()` method to CLI
- [ ] Implement markdown documentation generator
- [ ] Add API examples to godoc

---

## Plugin System

- [ ] Implement plugin system for custom validators
- [ ] Add validation interface abstraction
- [ ] Add `FlagRegistry` interface
- [ ] Custom validation hooks

---

## CLI Enhancements

- [x] Add Progress/Spinner Type using charmbracelet/bubbles
- [ ] Add enhanced flag validation enums

---

## Configuration

- [ ] Config File Auto-Loading integration with koanf
- [ ] Config file nested struct support
- [ ] Replace `internal/logging` with charmbracelet/log

---

## Fuzz Testing

- [x] Add fuzz tests to `flags_parse.go`
- [x] Add fuzz tests to `config_parsing.go`
- [ ] Add fuzz test corpus in `testdata/fuzz/` directories

---

## Metrics & Observability

- [x] Telemetry integration (OpenTelemetry)
- [ ] Metrics/hooks for custom observability

---

## go-output Dependency Architecture

> **Status:** Integrated in v2.x via `pkg/cmdguard/v2/output.go`. Published at v0.6.2.

`go-output` provides 12 output formats. Consider extracting to a sub-package in v3.0 so consumers only pay the dependency cost when they use `--output`.

---

## CI/CD & Release

- [x] Set up release automation
- [ ] Add codecov integration (needs `CODECOV_TOKEN` secret)
- [x] Set up CI/CD pipeline
- [x] Create contribution guide
- [ ] Deprecate v1 API timeline

---

## Future Ideas

- [ ] Add structured JSON error output for `--output=json`
- [x] Add issue/PR templates
- [ ] Test all examples in CI
- [ ] Extract flag-related code to standalone `flagtags` library
