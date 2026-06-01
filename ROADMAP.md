# ROADMAP

**Generated:** 2026-04-05
**Updated:** 2026-06-01
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

## Completed (v2.2–v2.3)

- [x] Add GitHub Actions workflow
- [x] Add badge to README
- [x] Add more custom types: URL, Email, Port, FilePath, HostPort
- [x] Add middleware support (`TimingMiddleware`, `RecoveryMiddleware`)
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

---

## Advanced Types

- [ ] Add Result[T] type for error handling
- [ ] Add Validated[T] wrapper with validation functions
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
- [ ] Add enhanced flag validation enums

---

## Configuration

- [ ] Config File Auto-Loading integration with koanf
- [ ] Replace `internal/config` with koanf
- [ ] Replace `internal/logging` with charmbracelet/log

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

## go-output Dependency Architecture

> **Status:** Integrated in v2.3.0 via `pkg/cmdguard/v2/output.go` (re-export wrapper).  
> **Question:** Should output rendering be optional or mandatory?

### Current State

`go-output` is a published dependency (v0.6.1) providing 12 output formats (table, json, csv,
yaml, markdown, xml, html, d2, mermaid, dot, tree, tsv). cmdguard wraps it in `output.go` to
provide `OutputResult`, `OutputTable`, `OutputStyledTable`, and the `--output` global flag
(`WithOutputFormat`).

**Dependency weight:** 11 go-output sub-modules + 99 transitive edges in cmdguard's mod graph.
Every consumer pays this cost whether they use `--output` or not.

### Options for v3.0

| Approach | Pros | Cons | Effort |
|----------|------|------|--------|
| **Keep as-is** | Zero code churn; `go install` works | Heavy dep tree for simple CLIs | None |
| **Extract to sub-package** (`pkg/cmdguard/v2/output`) | Consumers who don't need `--output` avoid import | `CLI[T].OutputFormat()` moves; breaking | M |
| **Split go-output into interface + backends** | Interface in core, heavy backends optional | Requires go-output redesign first | L |
| **Use stdlib `encoding/*` as fallback** | Zero external deps for json/csv/yaml | Loses styled tables, d2, mermaid | M |

### Recommendation

Keep as-is for v2.x. For v3.0, evaluate extracting `output.go` and `cli_output.go` into a
separate sub-package so `CLI[T]` doesn't carry the dependency unless the consumer imports the
output package explicitly.

---

## CI/CD & Release

- [ ] Set up release automation
- [ ] Add codecov integration
- [ ] Set up CI/CD pipeline
- [ ] Add pre-commit hooks
- [ ] Create contribution guide
- [ ] Deprecate v1 API timeline
- [ ] Remove testify/ginkgo completely

---

## Code Review Items

- [x] Review all `any` usages in package
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

- [ ] Add structured JSON error output for `--output=json`
- [x] Add `--no-color` flag + NO_COLOR support
- [ ] Add issue/PR templates
- [ ] Test all examples in CI
- [ ] Metrics/telemetry integration
- [ ] Custom validation hooks
