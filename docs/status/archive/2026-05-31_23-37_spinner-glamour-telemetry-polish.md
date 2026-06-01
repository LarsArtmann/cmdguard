# Status Report — 2026-05-31 23:37

**Project:** cmdguard v2.3.0-dev
**Branch:** master (clean)
**Tests:** 272 passing, 0 failures, 83.3% coverage
**Lint:** 0 issues (golangci-lint 2.x)
**Race:** 0 detected

---

## a) FULLY DONE ✅

### Core Library (v2 API)

- **CLI[T] + Command[T,F]** — Type-safe DI-powered CLI framework. Single type param on CLI, per-command flags via Command. Production-grade, no panics.
- **Flag System** — Struct tag flags (`flag`, `short`, `default`, `help`, `env`, `required`, `validate`, `values`, `count`, `prompt`). Instance-scoped type handler and validator registries. Typo suggestions (Levenshtein).
- **9 Value Types** — Duration, Email, Enum, FilePath, HostPort, LogLevel, LogFormat, Port, URL. All with MarshalText/UnmarshalText, ParseXxx constructors.
- **Dependency Injection** — samber/do/v2 integration via Scope. Provide/Invoke/ProvideValue. Lifecycle hooks (HealthCheck, Shutdown).
- **Middleware Chain** — buildChain with ordered wrapping. TimingMiddleware, RecoveryMiddleware (with stack traces), SpinnerMiddleware, TelemetryMiddleware.
- **BranchingFlowContext** — Hierarchical command path tracking with value propagation, typed timeout/deadline branches.
- **Config Files** — JSON built-in, optional YAML/TOML loaders. Path expansion ($ENV, ~). Config values become defaults, flags/env still override.
- **Output System** — 12+ formats via go-output (table, json, csv, yaml, markdown, xml, d2, etc.). WithOutputFormat[T] for auto --output flag.
- **Error System** — 40+ sentinel errors with errors.Is() chainability. Typed errors (CommandError, ServiceError, FlagError). ExitCoder + ExitError for custom exit codes.
- **Validation** — WithStrictValidation, WithDraconianValidation, WithConfigValidation. Phase-aware validation at AddCommand time.

### Recent Features (This Sprint)

- **SpinnerMiddleware** — Terminal spinner with Lipgloss styling, goroutine-based, auto-skips non-TTY, configurable frames/interval via SpinnerMiddlewareWithConfig.
- **Glamour Help Rendering** — WithGlamourHelp[T]() (auto theme) and WithGlamourHelpTheme[T]("dark") for markdown rendering of Long/Example fields. Idempotent (guards against double-render on repeated Execute calls).
- **TelemetryMiddleware** — OpenTelemetry span per command phase. Span names include command path: "deploy run", "myapp database migrate pre-run". SpanKindServer, command.fullpath attribute. WithTelemetry[T](tracer) convenience option.
- **FullPath in CommandInfo** — Populated at execution time via cobra.CommandPath(). Available to all middleware for hierarchical tracing/logging.

### Testing & Quality

- **272 tests**, 0 failures, 83.3% coverage, 0 race conditions
- **0 lint issues** (golangci-lint 2.x with 50+ linters)
- **Fuzz tests** for all value type parsers
- **Benchmarks** package exists
- **Integration tests** in tests/integration/

### Examples & Documentation

- **Consolidated taskctl example** — Single comprehensive example demonstrating every major feature (DI, flags, output, middleware, spinner, glamour, lifecycle, error handling, signal handling, config files, completion, hidden/deprecated commands, aliases, args validators, BranchingFlowContext, EditInEditor, version command).
- **AGENTS.md** — Full contributor guide with API reference, architecture decisions, gotchas, coding standards, test commands.
- **FEATURES.md** — Complete feature inventory with status indicators.
- **TODO_LIST.md** — Phase 1-9 all complete.

---

## b) PARTIALLY DONE ⚠️

### SpinnerMiddleware
- Works correctly for TTY vs non-TTY detection
- **Missing:** No test for actual TTY spinner rendering (would require terminal emulation). Only tested with bytes.Buffer (non-TTY path via direct newTextSpinner) and middleware-level tests for skip/error behavior.
- **Missing:** No validation on SpinnerConfig (e.g., negative interval, empty title, nil writer).

### TelemetryMiddleware
- Works with any trace.Tracer, creates phase-named spans
- **Missing:** No test verifying span parent-child relationships in a multi-command execution. Only unit-tested with mock tracer.
- **Missing:** WithTelemetry[T] is a convenience wrapper but doesn't offer any config beyond the tracer (e.g., custom span naming, attribute filtering).

### Glamour Help Rendering
- Works, idempotent, theme-customizable
- **Missing:** No test for the auto theme (which falls back to notty in non-TTY and doesn't transform markdown). Only "dark" theme is tested in E2E.
- **Missing:** applyGlamourHelp silently swallows rendering errors (only checks `err == nil`). Could log a warning.

### Value Types Architecture
- All 9 types work, but they have **repetitive patterns** (ParseXxx, String, MarshalText, UnmarshalText, error constructors). This was noted in Phase 9 as "consolidate patterns" but the consolidation was only partial.
- LogLevel and LogFormat are type aliases over Enum, which is a good pattern but inconsistent — Email, URL, Port, etc. don't follow this pattern.

---

## c) NOT STARTED 📝

### Performance
- CLI construction benchmarks
- Flag parsing benchmarks
- Command execution benchmarks
- Benchmark regression detection in CI

### CI/CD
- codecov integration
- v2.3.0 release tag and notes
- Release automation

### Future (v3.0+)
- Plugin system for custom validators and type handlers

### Future Cleanup (v3.0, API-breaking)
- Make NoFlags a distinct named type (not type alias)
- Change TimingMiddleware callback to include error
- Remove string-based BranchWithTimeout/BranchWithDeadline
- Remove FlowContextAccessor
- Rename Get[T]/MustGet[T] to more specific names
- Make RegisterInScope generic instead of ...any
- Remove or redesign Package()

---

## d) TOTALLY FUCKED UP 💥

### Nothing is catastrophically broken.

However, there are **architectural concerns**:

1. **glamourHelp field is unexported but tested externally** — `TestWithGlamourHelp` in `glamour_test.go` (same package `v2`) accesses `cli.glamourHelp` directly. This works because the test is in the same package, but it means the field is effectively public within the package but not documented as part of the API surface.

2. **CommandInfo mutation in closure** — In `cli_command.go`, `info` is captured by reference in the handler closure and mutated (`info.FullPath = c.CommandPath()`). This is safe because each handler gets its own copy via `info := cfg.info`, but the pattern is subtle and could confuse future contributors.

3. **spinnerFrames() is a function returning a new slice on every call** — Not cached, allocates on every call. Minor but creates unnecessary GC pressure if called frequently.

---

## e) WHAT WE SHOULD IMPROVE 🔄

### Architecture

1. **Value type consolidation** — The 9 value types (Duration, Email, Enum, FilePath, HostPort, LogLevel, LogFormat, Port, URL) all follow the same pattern but with slightly different error handling. A shared `ValidatedString` or `ValidatedValue[T]` generic base could eliminate 60%+ of the boilerplate.

2. **SpinnerConfig validation** — Currently accepts any input. Should validate: non-negative interval, non-nil writer, non-empty frames slice, non-empty title.

3. **Middleware composition API** — Currently middleware is a flat slice. For complex apps, conditional middleware (e.g., spinner only on certain commands) requires manual middleware logic. A `MiddlewareWhen` or `MiddlewareFor` helper could help.

4. **Config file nested struct support** — Currently flat only. Nested YAML/TOML configs are common in real applications.

5. **Error context propagation in glamour** — applyGlamourHelp silently ignores rendering errors. Should at minimum log to stderr.

### Code Quality

6. **spinnerFrames() should be a package-level var** — Not a function returning new slice each time. Already flagged by gochecknoglobals linter (as `defaultSpinnerFrames` doesn't exist yet).

7. **Test coverage gaps** — configload package is 0% covered. prompts.go (huh integration) has minimal test coverage.

8. **E2E/integration tests are thin** — The tests/integration package exists but has minimal coverage. No full CLI lifecycle tests (create CLI, add commands, execute, check output, shutdown).

### Documentation

9. **AGENTS.md is stale on test count** — Says "224 tests" but we now have 272. The "84.5% coverage" claim is also stale (now 83.3%).

10. **FEATURES.md is missing entries** — WithGlamourHelpTheme, SpinnerMiddlewareWithConfig, SpinnerConfig, FullPath in CommandInfo, RenderMarkdown, RenderMarkdownWithTheme are not listed.

---

## f) Top 25 Things We Should Get Done Next

**Sorted by Impact × Ease (Pareto)**

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Update AGENTS.md test count (224→272) and coverage (84.5%→83.3%) | High | 5min | Docs |
| 2 | Update FEATURES.md with spinner/glamour/telemetry entries | High | 15min | Docs |
| 3 | Add SpinnerConfig validation (negative interval, nil writer) | Medium | 15min | Quality |
| 4 | Fix spinnerFrames() to be a package-level var (not function) | Low | 5min | Lint/Perf |
| 5 | Add CLI construction benchmark | Medium | 30min | Perf |
| 6 | Add flag parsing benchmark | Medium | 30min | Perf |
| 7 | Add command execution benchmark | Medium | 30min | Perf |
| 8 | Add error logging in applyGlamourHelp on render failure | Medium | 10min | Quality |
| 9 | Add test for SpinnerConfig with nil writer | Medium | 10min | Testing |
| 10 | Add test for SpinnerConfig with zero interval | Medium | 10min | Testing |
| 11 | Add test for glamour "auto" theme behavior in non-TTY | Low | 15min | Testing |
| 12 | Add E2E test: full CLI lifecycle (create→execute→shutdown) | High | 1hr | Testing |
| 13 | Add telemetry test: multi-command span parent-child | Medium | 30min | Testing |
| 14 | Consolidate value type ParseXxx patterns into shared helper | High | 2hr | Architecture |
| 15 | Add configload tests (currently 0% coverage) | Medium | 1hr | Testing |
| 16 | Add prompts.go tests (huh integration) | Medium | 1hr | Testing |
| 17 | Add benchmark regression detection to CI | Medium | 30min | CI/CD |
| 18 | Add codecov integration | Medium | 30min | CI/CD |
| 19 | Create v2.3.0 release tag and notes | High | 1hr | Release |
| 20 | Add MiddlewareFor helper (conditional middleware per command) | Medium | 2hr | Feature |
| 21 | Add nested config struct support for YAML/TOML | High | 4hr | Feature |
| 22 | Add Plugin system for validators/type handlers | Medium | 4hr | Architecture |
| 23 | Set up release automation (goreleaser) | Medium | 2hr | CI/CD |
| 24 | v3.0 API cleanup (NoFlags distinct type, remove deprecated APIs) | High | 4hr | Architecture |
| 25 | Full integration test suite for taskctl example | Medium | 2hr | Testing |

---

## g) Top #1 Question I Cannot Figure Out Myself 🤔

**What is the target audience for the telemetry feature?**

The TelemetryMiddleware requires a `trace.Tracer` from the OpenTelemetry SDK, but cmdguard does NOT re-export or vendor the OTel SDK — consumers must import `go.opentelemetry.io/otel` themselves. This means:

- The example cannot demonstrate a working tracer setup without adding heavy OTel SDK dependencies (exporters, resources) to go.mod
- Consumers must understand OTel initialization to use this feature
- The feature is essentially "middleware that creates spans" — the hard part (SDK setup, exporters, sampling) is left to the user

**Should we:**
1. Keep it as-is (thin middleware, user brings their own tracer) — simplest, but the example can't demonstrate it
2. Add a helper like `NewStdoutTracerProvider()` that wraps the OTel SDK boilerplate — adds deps but makes the example work
3. Document it as "requires OTel SDK setup, see example" without adding the SDK to the example

This affects whether we add OTel SDK deps to go.mod (currently only `otel` and `otel/trace` are there, not the SDK/exporters).

---

## Git State

**Last commit:** `0ad8774` feat(examples/taskctl): add spinner middleware, glamour help rendering, and markdown descriptions
**Working tree:** clean
**Untracked:** none
**Modified:** none
