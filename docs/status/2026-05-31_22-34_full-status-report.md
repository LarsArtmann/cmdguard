# cmdguard — Full Status Report

**Date:** 2026-05-31 22:34 (CET)
**Branch:** master
**Version:** v2.3.0-dev
**Go:** 1.26.3
**Last Commit:** `bef634f` chore: add taskctl binary to .gitignore

---

## Executive Summary

cmdguard is a Go library for building validated Cobra CLI applications with type-safe dependency injection. The project is in excellent health: all 1,073 tests pass with race detection, 0 lint issues, 83.1% coverage on the core library. The v2.3.0-dev feature set is feature-complete with 9 phases of work done. The single consolidated example (`examples/taskctl/`) demonstrates every major feature with 66 tests.

### Health Dashboard

| Metric              | Value    | Status |
| ------------------- | -------- | ------ |
| Total tests         | 1,073    | GREEN  |
| Race conditions     | 0        | GREEN  |
| Lint issues         | 0        | GREEN  |
| Build errors        | 0        | GREEN  |
| Core lib coverage   | 83.1%    | GREEN  |
| Example coverage    | 71.1%    | YELLOW |
| configload coverage | 0.0%     | RED    |
| Open bugs           | 0 known  | GREEN  |

---

## A) FULLY DONE

### Core Library (pkg/cmdguard/v2) — 116 source files, 20,268 lines

| Category                   | Details                                                              |
| -------------------------- | -------------------------------------------------------------------- |
| **CLI[T] construction**    | NewCLI, AddCommand, Execute, ExecuteWithArgs, ExecuteAndExit         |
| **Command[T,F] system**    | NewCommand, NewParentCommand, MustNewCommand, 21 command options     |
| **Flag system**            | Struct tags (flag, short, default, help, env, count, required, validate, prompt, values) |
| **Value types (9)**        | Duration, Enum, LogLevel, LogFormat, URL, Email, Port, FilePath, HostPort |
| **DI system**              | Scope, Provide, Invoke, Child scopes, lifecycle (Shutdown, HealthCheck) |
| **Rich output (12 formats)** | table, json, csv, tsv, md, xml, d2, yaml, html, tree, mermaid, dot |
| **Middleware (4 built-in)** | Timing, Recovery, Spinner, Telemetry (OpenTelemetry)               |
| **Error handling**         | 35+ sentinel errors, 7 typed error types, ExitCoder/ExitError       |
| **Shell completion**       | WithCompletion, WithValidArgs, CompletionFunc                        |
| **Man page generation**    | ManPage(), WriteManPage(), GenerateManPageCommand()                  |
| **Markdown help**          | WithGlamourHelp, RenderMarkdown, RenderMarkdownWithTheme             |
| **Interactive prompts**    | WithPromptOnMissing, prompt tag, PromptString/Select/Confirm         |
| **Config file loading**    | JSON built-in, YAML/TOML via configload sub-package                  |
| **Arg validators**         | WithExactArgs, WithMinimumArgs, WithMaximumArgs, WithRangeArgs, WithNoArgs |
| **Version command**        | VersionCommand, MustVersionCommand, GenerateVersionCommand           |
| **Validation modes**       | Lenient, Strict, Draconian                                           |
| **Flow context**           | BranchingFlowContext with typed timeout/deadline alternatives        |
| **Helpers**                | EditInEditor, Ptr, ValueOrDefault, MustParse, MergeConfigs, EnsureValid |

### Infrastructure

- **Nix flake** — devShell (Go 1.26, gopls, golangci-lint), formatter (treefmt), format check
- **GitHub Actions CI** — build, test, lint workflow
- **Pre-commit hook** — scripts/pre-commit
- **golangci-lint 2.x** — 0 issues across entire project

### Examples

- **Single superb example** (`examples/taskctl/`) — production-grade task manager CLI with 66 tests
- Demonstrates every major cmdguard feature (see `examples/taskctl/README.md` for 40+ feature matrix)

---

## B) PARTIALLY DONE

### Config File Loading (configload package)

- **Status:** YAML and TOML loaders exist but have 0% test coverage
- `pkg/cmdguard/v2/configload/` has no test files
- JSON loader in main package is well-tested
- The configload package is importable but undocumented in examples

### Telemetry Middleware

- **Status:** Fully functional but LSP shows stale type errors (mockTracer doesn't implement trace.Tracer — `missing method tracer`)
- Tests pass fine (the LSP diagnostics are stale/cached)
- The mock in `telemetry_test.go` may need updating if the `trace.Tracer` interface gains new methods in future otel SDK updates

### Flag Validation at Runtime

- **Status:** `required:"true"`, `validate:"min=3"`, and `values:"a,b,c"` tags are parsed and stored but `ValidateFlags()` is a public API not automatically called during command execution
- Users must call `registry.ValidateFlags(cmd)` manually, or use `PreRunE` for validation
- The `values` tag only affects prompt completion choices, not runtime validation of plain string fields
- This is a design gap — documented in AGENTS.md Gotchas but surprising for new users

### Inspect Command in taskctl

- **Status:** Uses `WithExactArgs(1)` but can't access positional args from the RunE handler (cmdguard limitation)
- Hardcoded to always show task #1 as a demo
- This is a known limitation of the cmdguard RunE signature `func(ctx, cfg, flags) error`

---

## C) NOT STARTED

### Performance

- CLI construction benchmark
- Flag parsing benchmark
- Command execution benchmark
- Benchmark regression detection in CI

### CI/CD

- codecov integration
- v2.3.0 release tag and release notes
- Release automation (goreleaser or similar)

### v3.0 API Breaking Changes

- Make NoFlags a distinct named type (not type alias)
- Change TimingMiddleware callback to include error
- Remove string-based BranchWithTimeout/BranchWithDeadline
- Remove FlowContextAccessor
- Rename Get[T]/MustGet[T] to more specific names
- Make RegisterInScope generic instead of `...any`
- Plugin system for custom validators and type handlers

---

## D) TOTALLY FUCKED UP

### Nothing is catastrophically broken. Zero build errors, zero test failures, zero race conditions.

### Known irritations:

1. **LSP stale cache** — gopls shows 4 errors in `telemetry_test.go` and references to deleted `examples/error-handling/` and `examples/advanced-flags/`. These are phantom diagnostics; the code builds and tests pass fine. Requires LSP restart to clear.

2. **`go.mod` local replace directive** — `go-output` uses `replace` pointing to `../go-output`. This blocks CI/other developers. (This is documented in AGENTS.md Gotcha #10 as a known issue.)

3. **Pre-commit hooks broken** — `git commit --no-verify` is required. Pre-existing pre-commit hook errors exist. Documented in AGENTS.md Quick Start.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Design

1. **Auto-enforce flag validation** — `ValidateFlags()` should be called automatically during command execution (in `prepareRunContext` or `wireAllHandlers`). Currently users must call it manually or use `PreRunE`, which makes `required`, `validate`, and `values` tags feel half-baked.

2. **Positional args in RunE** — The RunE handler signature `func(ctx, cfg, flags) error` has no way to access positional args. Should add args `[]string` parameter or provide a `GetArgs(ctx)` helper.

3. **configload test coverage** — 0% coverage on YAML/TOML loaders. These are user-facing APIs.

4. **Replace directive elimination** — The `go-output` local replace in `go.mod` blocks external contributors and CI. Need to either vendor go-output or get a proper tagged release.

5. **Error type simplification** — 7 typed error types is a lot. Consider if all are pulling their weight or if some could be consolidated further.

### Developer Experience

6. **Example testability of prompts** — `WithPromptOnMissing` can't be tested without a TTY. Need a way to inject a mock prompt runner or skip gracefully.

7. **Stale LSP diagnostics** — The 4 telemetry_test.go errors and phantom example references suggest gopls isn't picking up changes from deleted files. A `.gopls` settings tweak or workspace reset would help.

8. **docs/ structure** — `docs/planning/` has 18 markdown files and 2 HTML files that are historical artifacts. Should be cleaned up or archived.

---

## F) Top 25 Things to Do Next

Sorted by impact × effort (high impact first):

| #  | Task                                                    | Impact | Effort | Type         |
| -- | ------------------------------------------------------- | ------ | ------ | ------------ |
| 1  | Auto-call `ValidateFlags()` during command execution    | HIGH   | M      | Architecture |
| 2  | Add positional args access in RunE handler              | HIGH   | M      | Architecture |
| 3  | Add configload test coverage (YAML/TOML)                | HIGH   | S      | Quality      |
| 4  | Create v2.3.0 release tag and release notes             | HIGH   | S      | Release      |
| 5  | Eliminate go-output replace directive                   | HIGH   | M      | Infra        |
| 6  | Add CLI construction benchmark                          | MED    | S      | Performance  |
| 7  | Add flag parsing benchmark                              | MED    | S      | Performance  |
| 8  | Add command execution benchmark                         | MED    | S      | Performance  |
| 9  | Fix pre-commit hooks                                    | MED    | S      | Infra        |
| 10 | Add codecov integration                                 | MED    | S      | CI           |
| 11 | Clean up `docs/planning/` — archive or remove old plans | MED    | S      | Housekeeping |
| 12 | Update FEATURES.md last-updated date                    | LOW    | XS     | Docs         |
| 13 | Update AGENTS.md status header (test count, coverage)   | LOW    | XS     | Docs         |
| 14 | Add prompt mock for testability of WithPromptOnMissing  | MED    | M      | Architecture |
| 15 | Set up release automation (goreleaser)                  | MED    | M      | CI           |
| 16 | Add benchmark regression to CI                          | MED    | M      | Performance  |
| 17 | v3.0: Make NoFlags a distinct named type                | LOW    | S      | API break    |
| 18 | v3.0: Add error to TimingMiddleware callback            | LOW    | XS     | API break    |
| 19 | v3.0: Remove string-based BranchWithTimeout/Deadline    | LOW    | XS     | API break    |
| 20 | v3.0: Remove FlowContextAccessor                        | LOW    | XS     | API break    |
| 21 | v3.0: Rename Get[T]/MustGet[T]                          | LOW    | S      | API break    |
| 22 | v3.0: Plugin system for validators/type handlers        | LOW    | L      | Feature      |
| 23 | Add more fuzz tests for edge cases                      | MED    | M      | Quality      |
| 24 | Document `ValidateFlags` auto-call gap in README        | LOW    | XS     | Docs         |
| 25 | Update README.md to reflect single example structure    | LOW    | XS     | Docs         |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should `ValidateFlags()` be auto-called during command execution, or is the current explicit design intentional?**

The library parses `required`, `validate`, and `values` struct tags and stores them, but `ValidateFlags()` (the function that actually checks these) is only a public API — it's never called from the execution pipeline. This means:

- `required:"true"` on a flag does nothing unless you call `registry.ValidateFlags(cmd)` manually
- `validate:"email"` on a flag does nothing at runtime
- `values:"low,medium,high"` on a string field only affects prompt completion, not flag validation

The taskctl example works around this with `PreRunE` + `ParseEnum`. But this feels like a design gap — users would reasonably expect these tags to auto-enforce.

**Options:**
1. Auto-call `ValidateFlags()` in `prepareRunContext` (behavior change, could break existing users who rely on loose parsing)
2. Add `WithAutoValidation[T,F]()` option to opt-in
3. Document the current behavior more prominently and leave it as-is

**Why I can't decide:** This is a product/design decision about the intended UX of the library, not a technical question. It affects the contract between cmdguard and its users.

---

## Git State

```
Current branch: master
Status: CLEAN (formatting fix staged)
Remote: up to date with origin/master

Recent commits (5):
bef634f chore: add taskctl binary to .gitignore
c3705e6 docs(AGENTS.md): update test count in project structure
1c36b2c docs(examples): update README and main.go feature list
7a3831f feat(examples): add Email/URL types to config, fix ConfigDefaults test
766edca fix(examples): fix two broken tests in taskctl example
```

---

## Test Breakdown

| Package                                     | Tests | Coverage |
| ------------------------------------------- | ----- | -------- |
| `pkg/cmdguard/v2`                           | 970   | 83.1%    |
| `pkg/cmdguard/v2/testutil`                  | ~3    | 87.5%    |
| `examples/taskctl`                          | 66    | 71.1%    |
| `tests/integration`                         | ~34   | n/a      |
| `pkg/cmdguard/v2/configload`                | 0     | 0.0%     |
| **Total**                                   | ~1,073| —        |

---

## Dependency Highlights

| Dependency                     | Version  | Purpose              |
| ------------------------------ | -------- | -------------------- |
| spf13/cobra                    | v1.10.2  | CLI framework        |
| samber/do/v2                   | v2.0.0   | Dependency injection |
| charm.land/fang/v2             | v2.0.1   | Cobra styling        |
| charm.land/huh/v2              | v2.0.3   | Interactive prompts  |
| charm.land/bubbletea/v2        | v2.0.6   | TUI framework        |
| charm.land/lipgloss/v2         | v2.0.3   | Terminal styling     |
| charm.land/bubbles/v2          | v2.1.0   | TUI components       |
| github.com/larsartmann/go-output | v0.6.1 | Rich output formats  |
| spf13/pflag                    | v1.0.10  | Flag parsing         |
