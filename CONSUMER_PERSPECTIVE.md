# Consumer Perspective: What cmdguard is Missing

> An honest audit from the perspective of a Go developer evaluating cmdguard for their next project.
> **Date:** 2026-05-27 | **Based on:** v2.3.0-dev, 264 tests, 81.2% coverage

---

## Critical Missing Pieces (Adoption Blockers)

### 1. No Migration Guide from Plain Cobra

cmdguard's primary audience is **existing Cobra users** who want type safety. Yet there's no "migrating from Cobra" guide. The only migration doc covers v1→v2 (internal API evolution). A Cobra user asking "how hard is it to adopt cmdguard incrementally?" finds no answer.

**What's needed:** `docs/MIGRATION_FROM_COBRA.md` showing step-by-step how to wrap an existing Cobra app, starting with just flags and progressively adopting DI, lifecycle hooks, and validation.

### 2. No Comparison with Alternatives

A developer evaluating CLI frameworks will compare cmdguard against Kong, urfave/cli, go-flags, and sflags. The README only shows "raw Cobra vs. cmdguard" — a comparison against the weakest alternative. The competitive analysis exists internally in `PARTS.md` but is hidden from consumers.

**What's needed:** A "Why cmdguard?" section in the README with a comparison table against the 2-3 most popular alternatives, or a dedicated `COMPARISON.md` linked from the README.

### 3. No API Stability Guarantee

The CHANGELOG says "adheres to Semantic Versioning" but there's no explicit stability guarantee. A team considering cmdguard for production needs to know:

- Is the v2 API stable? Will it break before v3?
- What's the deprecation policy for marked items?
- What's the expected timeline between major versions?

**What's needed:** A versioning/stability statement. Even a single sentence: "The v2 API is stable and will only receive additive changes until v3."

### 4. No Test Harness for Consumers

Every example test manually constructs `CLI[T]`, calls `ExecuteWithArgs`, and captures output. The `pkg/testutil` package exists but only has `panic_test_helpers.go` — nothing that helps a consumer test their CLI.

Competitors provide testing utilities: urfave/cli has `cli.NewContext()`, Kong has `kong.Exit(func(s int))`. cmdguard consumers repeat the same boilerplate in every test file.

**What's needed:** A public `testutil` package with `NewTestCLI[T]`, `CaptureOutput`, and `TestResult` (stdout, stderr, exit code, error).

---

## Documentation Gaps (Discovery Blockers)

### 5. No Tutorial or Step-by-Step Guide

QUICKSTART.md covers API surface in 310 lines, but there's no narrative tutorial. A new user doesn't want a reference — they want to **build something real** end-to-end. "Building a Task CLI with cmdguard" would teach DI, flags, subcommands, validation, and output formatting in a single cohesive walkthrough.

### 6. README is 25+ Public APIs Behind

The README documents the initial feature set but hasn't kept up with v2.2/v2.3 additions. Missing from the README:

- `WithConfigValidation`, `WithStrictValidation`, `WithDraconianValidation`
- `WithOutputFormat`, `WithConfigFile`, `WithConfigFileLoader`
- `WithEnvPrefix`
- `EditInEditor`
- Shell completion (`WithCompletion`, `WithValidArgs`)
- Man page generation
- `BranchingFlowContext` API
- Middleware (`TimingMiddleware`, `RecoveryMiddleware`)
- Positional args validators
- `ExitCoder` / `NewExitError`
- `WithGroup` / `WithGroupID`

A consumer reading the README sees maybe 60% of the library's capabilities.

### 7. No `doc.go` for pkg.go.dev

The package-level godoc is 2 lines on `errors.go`. When a developer visits pkg.go.dev, they see a bare function list — no overview, no getting-started snippet, no linked examples. A `doc.go` with a comprehensive package overview would dramatically improve the API documentation experience.

### 8. No godoc Examples for Key APIs

Go's `Example*()` test functions render in pkg.go.dev. cmdguard has 13 examples but none for the most important APIs: `NewCLI`, `Provide`/`Invoke` (DI), `OutputTable`, middleware, or error handling. These are the APIs consumers struggle with most.

---

## Example Gaps (Learning Blockers)

### 9. 12+ Features Have No Example

| Feature                                               | Has Example? |
| ----------------------------------------------------- | ------------ |
| `WithStrictValidation`                                | No           |
| `WithConfigValidation`                                | No           |
| `VersionCommand` / `MustVersionCommand`               | No           |
| Positional args (`WithExactArgs`, etc.)               | No           |
| `WithOutputFormat`                                    | No           |
| Middleware (`TimingMiddleware`, `RecoveryMiddleware`) | No           |
| Man page generation                                   | No           |
| Shell completion (`WithCompletion`)                   | No           |
| `EditInEditor`                                        | No           |
| `ExitCoder` / `ExitError`                             | No           |
| `WithGroup` / `WithGroupID`                           | No           |
| `MustNewCommand` / `MustNewParentCommand`             | No           |

**What's needed:** Either a "kitchen sink" example (`examples/superb/`) or individual examples per feature.

### 10. No "Real World" Example

All 12 examples are toy demonstrations — "hello world", "greet someone", "list users". There's no realistic CLI showing how cmdguard handles a production scenario: multiple services, error recovery, config files, signal handling, and rich output all working together. This is what convinces a senior engineer that the library is production-ready.

---

## Feature Gaps (Competitive Disadvantages)

### 11. No Structured Error Output

When `--output=json` is set, data output is JSON but errors are plain text. The typed error hierarchy (`CommandError`, `FlagError`, `ServiceError`) would serialize beautifully but isn't wired for it. This matters for CI/automation consumers who parse CLI output programmatically.

### 12. No `NO_COLOR` Documentation or `--no-color` Flag

Most production CLIs respect `NO_COLOR` (no-color.org) and provide `--no-color`. Fang/lipgloss handles it implicitly, but this isn't documented. There's no explicit flag and no mention in README.

### 13. No Performance Story

14 benchmarks exist in `benchmarks/` but no results are documented. A consumer comparing CLI frameworks has no way to assess cmdguard's overhead. Even a simple statement like "cmdguard adds <1ms overhead over raw Cobra" with benchmark numbers would help.

---

## Infrastructure Gaps (Trust Blockers)

### 14. Stale ROADMAP.md

The ROADMAP lists items as unchecked that are already complete: "Add GitHub Actions workflow" (exists), "Add badge to README" (exists), "Add more custom types: URL, Email, Port, FilePath" (all exist), "Add middleware support" (exists), "Create command groups feature" (exists), "Environment Variable Binding with env struct tags" (exists), "Config file support YAML/TOML" (exists). This signals to consumers that the project may not be actively maintained or that docs are an afterthought.

### 15. No SECURITY.md

No security policy for reporting vulnerabilities. Standard for any production library.

### 16. No Issue/PR Templates

No `.github/ISSUE_TEMPLATE/` or `.github/PULL_REQUEST_TEMPLATE.md`. Increases friction for community contributions and bug reports.

### 17. Examples Not Tested in CI

6 of 12 examples have no test files. The CI runs `go test ./...` which compiles examples with tests, but nothing verifies that examples actually work as documented.

### 18. go-output Uses Local Replace Directive

`go.mod` has a local replace directive for `github.com/larsartmann/go-output`. This blocks anyone else from building the project without manually setting up the local dependency — a critical issue for any consumer trying to `go get` or contribute.

---

## Prioritized Action Items

| Priority | Item                                     | Effort | Impact                      |
| -------- | ---------------------------------------- | ------ | --------------------------- |
| P0       | Remove go-output local replace directive | Low    | Unblocks adoption           |
| P0       | Add Cobra migration guide                | Medium | Targets core audience       |
| P0       | Update README with missing 25+ APIs      | Medium | Shows full capability       |
| P1       | Add `doc.go` with comprehensive godoc    | Low    | Improves pkg.go.dev         |
| P1       | Add consumer test harness (`NewTestCLI`) | Medium | Unblocks testing            |
| P1       | Add comparison table/COMPARISON.md       | Low    | Helps evaluation            |
| P1       | Fix stale ROADMAP.md                     | Low    | Signals active maintenance  |
| P1       | Add API stability statement              | Low    | Builds trust                |
| P2       | Add kitchen-sink or real-world example   | Medium | Proves production readiness |
| P2       | Add godoc examples for key APIs          | Medium | Improves discoverability    |
| P2       | Add step-by-step tutorial                | High   | Best onboarding experience  |
| P2       | Document performance / benchmark results | Low    | Helps evaluation            |
| P2       | Add SECURITY.md                          | Low    | Standard practice           |
| P3       | Add `--no-color` flag + documentation    | Low    | CLI best practice           |
| P3       | Add structured JSON error output         | Medium | CI/automation support       |
| P3       | Add issue/PR templates                   | Low    | Community infrastructure    |
| P3       | Test all examples in CI                  | Low    | Reliability                 |
