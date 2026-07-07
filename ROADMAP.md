# ROADMAP

**Updated:** 2026-07-07
**Purpose:** Aspirational items with no concrete timeline

---

## Completed (v2.2–v2.8)

- [x] **Zero panics** — All Must\* functions removed in v2.5.0; every function returns errors
- [x] Remove all panic-inducing functions (16 deleted: MustNewCommand, MustNewParentCommand, MustNewCLI, MustAddCommand, MustVersionCommand, MustDoctorCommand, MustInvoke, MustInvokeNamed, MustGet, RequireBranchingFlowContext, MustParse, MustParseDuration, MustParseLogLevel, MustParseLogFormat, MustParseEnum, MustParseURL, MustParseEmail, MustParsePort, MustParseFilePath, MustParseHostPort)
- [x] Deep codebase review and improvement sprint (11 phases)
- [x] Error system overhaul: 60+ sentinels, domain-specific error files
- [x] Runtime type guards in dispatchParse
- [x] Short tag length validation
- [x] Nil tracer guard in TelemetryMiddleware
- [x] BranchingFlowContext.SetValue child safety
- [x] Mutable slice protection (Tags/Path return clones)
- [x] Arg validator error returns instead of panics
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
- [x] Rich output with 16 formats via `go-output`
- [x] Review all `any` usages in package
- [x] Remove string-based `BranchWithTimeout`/`BranchWithDeadline`
- [x] Remove deprecated `FlowContextAccessor` API

---

## Completed (v3.0.0) — 2026-07-07

The v3.0 redesign shipped on `github.com/larsartmann/cmdguard/v3`:

- [x] Create `pkg/cmdguard/v3/` directory (migrated via `git mv`)
- [x] Implement non-generic `CLIOption` / `CommandOption` (eliminate type-param explosion)
- [x] `NewCommand` / `NewParentCommand` positional-flags signature (full type inference)
- [x] Write tests for v3 implementation (457 functions, 1430 runs, 87.3% coverage)
- [x] Create v3 example (taskctl)
- [x] Extract 5 optional sub-modules: `glamour`, `manpage`, `prompts`, `spinner`, `telemetry`
- [x] Go workspace (`go.work`, 6 modules)
- [x] Write `docs/MIGRATION_v2_v3.md`
- [x] Make `NoFlags` a distinct named type
- [x] Fix `telemetry.WithTelemetry` to return non-generic `CLIOption`

> **Note:** several v2.x features (Spinner, Glamour, Telemetry, Manpage, Prompts impl)
> were extracted into optional sub-modules in v3.0 rather than removed — they remain
> available via `github.com/larsartmann/cmdguard/<module>`. `EditInEditor` and
> `Result[T]`/`Validated[T]` were removed entirely as non-CLI concerns.

### Deferred v3.0 API-Breaking Cleanup

These remain open for a future v3.x or v4:

- [ ] **Rename `Get[T]`/`MustGet[T]`** — `Get` is too generic; should be
      `GetService[T]` or similar. Breaking: every consumer's import surface changes.
- [ ] **Make `RegisterInScope` generic** — currently takes `...any`; should be
      `RegisterInScope[T](scope, provider)`. Breaking: signature change.
- [ ] **Remove or redesign `Package()`** — takes a pre-existing `*Scope` which is
      an unusual API shape; should be reworked for error-safe DI. Breaking.
- [ ] **Remove `SetConfig`** — mutating a CLI's config after construction is
      unsafe (FlagRegistry isn't re-initialized). Breaking but removes a footgun.

### v3.0 Extraction: `flagtags` Library

- [ ] **Extract flag-tag parsing to `github.com/larsartmann/flagtags`**

  **Rationale:** The struct-tag parsing, validation, and type-handler registry
  are ~2000 lines of self-contained, reusable code with zero cmdguard-specific
  dependencies. Extracting it would:
  - Enable reuse in non-CLI contexts (HTTP handlers, config loaders)
  - Reduce cmdguard's compile surface
  - Allow independent versioning of the tag-parsing layer

  **Design:** `flagtags.Parse(v any) ([]Tag, error)`, `flagtags.RegisterHandler()`,
  `flagtags.RegisterValidator()`. cmdguard would re-export or wrap these.

  **Risk:** Medium — requires careful API stabilization before extraction.
  Defer until the tag format is frozen at v3.0.

---

## Advanced Types

- [x] Add `Result[T]` type for error handling (v2.8)
- [x] Add `Validated[T]` wrapper with validation functions (v2.8)
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

- [x] Implement plugin system for custom validators and type handlers (v2.8)
- [x] Add `Plugin` interface + `PluginRegistrar` + `RegisterPlugin` / `WithPlugin`
- [ ] Add `FlagRegistry` interface abstraction
- [ ] Custom validation hooks

---

## CLI Enhancements

- [x] Add Progress/Spinner Type using charmbracelet/bubbles
- [ ] Add enhanced flag validation enums

---

## Configuration

- [x] Config File Auto-Loading integration with koanf
- [x] Config file nested struct support (v2.8 — ParseFlagTags recurses into nested structs)

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

> **Status:** Integrated in v2.x via `pkg/cmdguard/v3/output.go`. Published at v0.17.0.

`go-output` provides 16 output formats as an independent module tree (root + 9 direct sub-modules + 3 indirect, all pinned in lockstep). Consumers already only pay the dependency cost when they use `--output`.

---

## CI/CD & Release

- [x] Set up release automation
- [ ] Add codecov integration (needs `CODECOV_TOKEN` secret)
- [x] Set up CI/CD pipeline
- [x] Create contribution guide
- [x] Create issue/PR templates
- [ ] Deprecate v1 API timeline

---

## Future Ideas

- [x] Add structured JSON error output for `--output=json`
- [x] Add issue/PR templates
- [ ] Test all examples in CI
- [ ] Extract flag-related code to standalone `flagtags` library
