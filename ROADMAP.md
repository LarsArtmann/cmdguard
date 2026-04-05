# ROADMAP

**Generated:** 2026-04-05
**Purpose:** Aspirational items with no concrete timeline

---

## v3.0 Major Redesign

> See `v3.0-major-redesign-plan.md` for details

- [ ] Create v3.0 API design document
- [ ] Create `pkg/cmdguard/v3/` directory
- [ ] Implement `errors.go` for v3
- [ ] Implement `types.go` for v3
- [ ] Implement `cli.go` for v3
- [ ] Implement `command.go` for v3
- [ ] Implement `options.go` for v3
- [ ] Implement `flags.go` for v3
- [ ] Implement `flags_parse.go` for v3
- [ ] Implement `flags_validate.go` for v3
- [ ] Implement `scope.go` for v3
- [ ] Implement `scope_provide.go` for v3
- [ ] Implement `cli_exec.go` for v3
- [ ] Write tests for v3 implementation
- [ ] Create v3 examples
- [ ] Write MIGRATION_V2_TO_V3.md

---

## Advanced Types

- [ ] Add Result[T] type for error handling
- [ ] Add Validated[T] wrapper with validation functions
- [ ] Add more custom types: URL, Email, Port, FilePath
- [ ] Create example application for branded IDs

---

## Documentation Generation

- [ ] Create examples/docs-generator/main.go
- [ ] Define FlagDoc struct
- [ ] Add GenerateDocs() method to CLI
- [ ] Implement markdown documentation generator
- [ ] Add GenerateDocsToFile() helper
- [ ] Add API examples to godoc
- [ ] Add flag validation examples

---

## Plugin System

- [ ] Implement plugin system for custom validators
- [ ] Add validation interface abstraction
- [ ] Add FlagRegistry interface
- [ ] Custom validation hooks

---

## CLI Enhancements

- [ ] Add Progress/Spinner Type using charmbracelet/bubbles
- [ ] Add Shell Completion Helpers
- [ ] Create command groups feature
- [ ] Add enhanced flag validation enums

---

## Configuration

- [ ] Config File Auto-Loading integration with koanf
- [ ] Environment Variable Binding with env struct tags
- [ ] Replace `internal/config` with koanf
- [ ] Replace `internal/logging` with charmbracelet/log
- [ ] Config file support YAML/TOML

---

## Fuzz Testing

- [ ] Add fuzz tests to flags_parse.go
- [ ] Add fuzz tests to config_parsing.go
- [ ] Add fuzz test corpus in testdata/fuzz/ directories

---

## Metrics & Observability

- [ ] Metrics/telemetry integration

---

## Standalone Library

- [ ] Create `github.com/larsartmann/flagtags` repository
- [ ] Extract flag-related code to standalone library

---

## CI/CD & Release

- [ ] Set up release automation
- [ ] Add GitHub Actions workflow
- [ ] Add codecov integration
- [ ] Add badge to README
- [ ] Set up CI/CD pipeline
- [ ] Add pre-commit hooks
- [ ] Create contribution guide
- [ ] Deprecate v1 API timeline
- [ ] Remove testify/ginkgo completely

---

## Code Review Items

- [ ] Review all `any` usages in package
- [ ] Document DI patterns
- [ ] Document DI scope pattern in docs/
- [ ] Document error handling strategy
- [ ] Review gochecknoglobals
- [ ] Review recvcheck
- [ ] Review unparam
- [ ] Review other examples for duplicate code
- [ ] Configure exhaustruct for external structs
- [ ] Rename test packages to use `_test` suffix

---

## Future Ideas

- [ ] Add middleware support
- [ ] Add more custom types (URL, Email, Port, FilePath)
- [ ] Create command groups feature
- [ ] Implement plugin system for custom validators
- [ ] Add enhanced flag validation enums
- [ ] Custom validation hooks
- [ ] Metrics/telemetry integration
