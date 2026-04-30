# cmdguard v2.2.0 — Comprehensive Micro-Task Plan

**Created:** 2026-04-30 03:15
**Status:** Planning — Awaiting Confirmation
**Scope:** ALL remaining work to make cmdguard SUPERB
**Total:** 72 micro-tasks across 5 phases, ~12h total

---

## Methodology

1. **Sources deduplicated:** Status report §E (18 items), §F (25 items), TODO_LIST.md (15 items), roadmap remaining tasks, dead code audit, coverage gaps
2. **Duplicates removed:** Items already completed (T1-T8, T12, T15, T16, T20, T24) excluded
3. **Each task ≤ 12 minutes** — single focused unit of work
4. **Sorted by:** Impact × Customer Value ÷ Effort (highest first)
5. **Phased by dependency:** Tests/fixes first → API polish → docs → new features → infra

---

## Phase Legend

| Phase | Focus | Tasks | Time |
|-------|-------|-------|------|
| **P1** | Critical Tests & Dead Code | #1–14 | ~2.5h |
| **P2** | API Polish & Fixes | #15–27 | ~2h |
| **P3** | Documentation & Examples | #28–40 | ~2h |
| **P4** | New Features | #41–62 | ~5h |
| **P5** | Infrastructure & CI | #63–72 | ~1.5h |

---

## Master Task Table

| # | Task | Impact | Effort | Files | Phase | Category |
|---|------|--------|--------|-------|-------|----------|
| 1 | Tests: type_handler.go — `registerCustomTypes()` coverage | ★★★★★ | 12m | type_handler.go, type_handler_test.go | P1 | test |
| 2 | Tests: type_handler.go — `registerCountHandler()` coverage | ★★★★★ | 12m | type_handler.go, type_handler_test.go | P1 | test |
| 3 | Tests: type_handler.go — `dispatchRegister` with count flag | ★★★★★ | 12m | type_handler.go, type_handler_test.go | P1 | test |
| 4 | Tests: type_handler.go — `lookupHandler` fallback (byKind) | ★★★★★ | 12m | type_handler.go, type_handler_test.go | P1 | test |
| 5 | Tests: type_handler.go — `RegisterTypeHandler` public API | ★★★★★ | 10m | type_handler.go, type_handler_test.go | P1 | test |
| 6 | Tests: output.go — `ParseOutputFormat` all 12 formats | ★★★★★ | 12m | output.go, output_test.go | P1 | test |
| 7 | Tests: output.go — `ParseOutputFormat` invalid/empty input | ★★★★★ | 10m | output.go, output_test.go | P1 | test |
| 8 | Tests: output.go — `OutputResult` with table format | ★★★★★ | 12m | output.go, output_test.go | P1 | test |
| 9 | Tests: output.go — `OutputTable` + `OutputStyledTable` | ★★★★★ | 12m | output.go, output_test.go | P1 | test |
| 10 | Tests: env tag integration — `t.Setenv` priority chain | ★★★★★ | 12m | flags_parse.go, flags_parse_test.go | P1 | test |
| 11 | Tests: env tag with prefix (`WithEnvPrefix`) | ★★★★ | 10m | cli_options.go, cli_test.go | P1 | test |
| 12 | Tests: counting flag — `-vvv` → 3, `-v` → 1, no flag → 0 | ★★★★ | 10m | flags.go, flags_test.go | P1 | test |
| 13 | Tests: editor.go — `EditInEditor` with mock `$EDITOR` | ★★★★ | 12m | editor.go, editor_test.go | P1 | test |
| 14 | Tests: command_suggest.go — `SuggestCommand` | ★★★★ | 10m | command_suggest.go, command_suggest_test.go | P1 | test |
| 15 | Dead code: remove `parseCustomDefault` from config_parsing.go | ★★★ | 5m | config_parsing.go:193 | P2 | cleanup |
| 16 | Dead code: remove `parseAndSetLogLevel` from config_setfield.go | ★★★ | 5m | config_setfield.go:188 | P2 | cleanup |
| 17 | Dead code: remove `parseAndSetLogFormat` from config_setfield.go | ★★★ | 5m | config_setfield.go:193 | P2 | cleanup |
| 18 | Dead code: remove `parseAndSetDuration` from config_setfield.go | ★★★ | 5m | config_setfield.go:198 | P2 | cleanup |
| 19 | Fix: consistent suggest API — `SuggestFlag` returns `string`, `SuggestCommand` returns `(string, bool)` — normalize | ★★★ | 10m | flags_suggest.go, command_suggest.go | P2 | api |
| 20 | Fix: move counting flag handler from `typeRegistry` struct to `registerKinds()` for consistency | ★★★ | 10m | type_handler.go | P2 | refactor |
| 21 | Fix: `wireSubcommandSuggestions` — improve unknown-command interception beyond `SetFlagErrorFunc` | ★★★★ | 12m | cli_command.go | P2 | fix |
| 22 | Deprecation: add runtime warning for `WithColor` (log.Println when used) | ★★ | 10m | cli_options.go | P2 | api |
| 23 | Feature: `WithOutputFormat[T]` CLI option — auto-adds `--output` flag | ★★★★ | 12m | cli_options.go, output.go | P2 | api |
| 24 | Feature: wire `--output` flag to `OutputResult` in command execution | ★★★★ | 12m | cli.go, output.go | P2 | api |
| 25 | Tests: `GenerateHelp()` and `FlagNames()` in flags_suggest.go | ★★★ | 12m | flags_suggest.go, flags_suggest_test.go | P2 | test |
| 26 | Fix: 32-bit int safety — verify `reflect.Convert` handles int64→int on 32-bit | ★★★ | 12m | type_handler.go, flags_parse.go | P2 | fix |
| 27 | Fix T27: implement `do.HealthcheckerWithContext` on `CLI[T]` to fix DI scope hang | ★★★ | 12m | cli.go, scope.go | P2 | fix |
| 28 | Docs: rewrite QUICKSTART.md for v2.2 API (NewCLI, AddCommand, constructors) | ★★★★ | 12m | docs/QUICKSTART.md | P3 | docs |
| 29 | Docs: update README.md with v2.2 features (env, output, counting, signals) | ★★★★★ | 12m | README.md | P3 | docs |
| 30 | Docs: add env tags example (`examples/env-tags/main.go`) | ★★★ | 12m | examples/env-tags/ | P3 | docs |
| 31 | Docs: add go-output example (`examples/output/main.go`) | ★★★ | 12m | examples/output/ | P3 | docs |
| 32 | Docs: add counting flags example (`examples/counting/main.go`) | ★★ | 10m | examples/counting/ | P3 | docs |
| 33 | Docs: add DI pattern example (`examples/di-patterns/main.go`) | ★★★ | 12m | examples/di-patterns/ | P3 | docs |
| 34 | Docs: add error handling example (`examples/error-handling/main.go`) | ★★★ | 12m | examples/error-handling/ | P3 | docs |
| 35 | Docs: add signal handling example (`examples/signals/main.go`) | ★★★ | 10m | examples/signals/ | P3 | docs |
| 36 | Docs: update AGENTS.md to v2.2 final (remove stale refs, add output/editor env) | ★★★ | 10m | AGENTS.md | P3 | docs |
| 37 | Docs: update FEATURES.md coverage numbers post-testing | ★★ | 8m | FEATURES.md | P3 | docs |
| 38 | Docs: add `WithOutputFormat` to API reference in AGENTS.md | ★★ | 8m | AGENTS.md | P3 | docs |
| 39 | Docs: update examples/ README.md | ★★ | 10m | examples/README.md | P3 | docs |
| 40 | Docs: update TODO_LIST.md — mark completed items | ★ | 5m | TODO_LIST.md | P3 | docs |
| 41 | Feature: config file auto-loading — add koanf + YAML parser deps | ★★★★★ | 10m | go.mod | P4 | feature |
| 42 | Feature: config file auto-loading — create `config_file.go` with load logic | ★★★★★ | 12m | config_file.go | P4 | feature |
| 43 | Feature: config file auto-loading — `WithConfigFile[T](path)` CLI option | ★★★★★ | 12m | cli_options.go, config_file.go | P4 | feature |
| 44 | Feature: config file auto-loading — `WithConfigFileFlags[T]()` auto `--config` flag | ★★★★ | 10m | cli_options.go | P4 | feature |
| 45 | Feature: config file auto-loading — merge chain (file → env → flag → arg) | ★★★★★ | 12m | config_file.go, flags_parse.go | P4 | feature |
| 46 | Feature: config file auto-loading — .env file support | ★★★★ | 10m | config_file.go | P4 | feature |
| 47 | Feature: config file auto-loading — TOML parser support | ★★★★ | 10m | config_file.go | P4 | feature |
| 48 | Feature: config file — tests (YAML load, priority, missing file, invalid) | ★★★★★ | 12m | config_file_test.go | P4 | feature |
| 49 | Feature: interactive prompts — add `huh` dependency | ★★★★ | 5m | go.mod | P4 | feature |
| 50 | Feature: interactive prompts — create `prompts.go` infrastructure | ★★★★ | 12m | prompts.go | P4 | feature |
| 51 | Feature: interactive prompts — `PromptString`, `PromptSelect`, `PromptConfirm` | ★★★★ | 12m | prompts.go | P4 | feature |
| 52 | Feature: interactive prompts — `WithPromptOnMissing[T,F]()` command option | ★★★★ | 12m | command.go, prompts.go | P4 | feature |
| 53 | Feature: interactive prompts — tests and example | ★★★★ | 12m | prompts_test.go, examples/prompts/ | P4 | feature |
| 54 | Feature: shell completion — `WithCompletion[T,F](fn)` command option | ★★★ | 12m | command.go, cli_command.go | P4 | feature |
| 55 | Feature: shell completion — wire to Cobra `ValidArgsFunction` | ★★★ | 12m | cli_command.go | P4 | feature |
| 56 | Feature: shell completion — Bash/Zsh/Fish/PowerShell generators | ★★★ | 12m | completion.go | P4 | feature |
| 57 | Feature: shell completion — tests and example | ★★★ | 10m | completion_test.go | P4 | feature |
| 58 | Feature: Markdown help — add `glamour` dep + `WithMarkdownHelp[T]()` | ★★★ | 12m | go.mod, help_markdown.go | P4 | feature |
| 59 | Feature: Markdown help — render Long descriptions as markdown | ★★★ | 12m | help_markdown.go | P4 | feature |
| 60 | Feature: spinner/progress — add `bubbles` dep + `SpinnerMiddleware` | ★★★ | 12m | go.mod, middleware_spinner.go | P4 | feature |
| 61 | Feature: spinner/progress — `WithSpinner[T](msg)` CLI option + Progress helper | ★★★ | 12m | middleware_spinner.go, progress.go | P4 | feature |
| 62 | Feature: telemetry middleware — `TelemetryMiddleware[T]` with OTel spans | ★★ | 12m | middleware_telemetry.go | P4 | feature |
| 63 | Infra: fix go-output dependency — remove local replace, use GOPRIVATE or publish | ★★★★★ | 12m | go.mod, AGENTS.md | P5 | infra |
| 64 | Infra: create v2.2.0 release tag and release notes | ★★★★ | 12m | — | P5 | release |
| 65 | Infra: set up release automation (goreleaser or Go module proxy) | ★★★ | 12m | .goreleaser.yml | P5 | infra |
| 66 | Infra: add codecov integration | ★★★ | 10m | .github/workflows/ | P5 | infra |
| 67 | Infra: fix pre-commit hooks | ★★★ | 12m | .pre-commit-config.yaml | P5 | infra |
| 68 | Infra: add benchmark regression detection script | ★★ | 10m | benchmarks/ | P5 | infra |
| 69 | Infra: merge test helper files (test_helpers_test.go + testhelpers_test.go) | ★★ | 12m | *_test.go | P5 | infra |
| 70 | Infra: add comprehensive performance benchmarks | ★★ | 12m | benchmarks/ | P5 | infra |
| 71 | Infra: man page generation via mango-cobra | ★★ | 12m | manpage.go, cli_options.go | P5 | feature |
| 72 | Infra: plugin system for custom validators and type handlers | ★★ | 12m | plugin.go | P5 | feature |

---

## Summary Statistics

| Metric | Value |
|--------|-------|
| Total micro-tasks | 72 |
| Estimated total time | ~12h |
| Phase P1 (tests/cleanup) | 14 tasks, ~2.5h |
| Phase P2 (API/fixes) | 13 tasks, ~2h |
| Phase P3 (docs/examples) | 13 tasks, ~2h |
| Phase P4 (new features) | 22 tasks, ~4h |
| Phase P5 (infra/CI) | 10 tasks, ~1.5h |
| Max task duration | 12m |
| Min task duration | 5m |
| Avg task duration | ~10m |

---

## Dependency Chain

```
P1 (tests) → P2 (fixes/api) → P3 (docs) → P4 (features) → P5 (infra)
                              ↑                 ↑
                              └── can parallel ──┘
```

- P1 must complete first — tests validate all existing code
- P2 depends on P1 — fixes may change behavior being tested
- P3 can overlap with P2 — docs can be written while API changes land
- P4 features depend on P2 API surface being stable
- P5 infra can run in parallel with P3/P4

---

## Execution Priority Within Each Phase

### P1 — Tests & Cleanup (execution order)

**Round 1: Highest-impact untested code** (~1h)
1. #1 → #2 → #3 → #4 → #5 (type_handler full coverage)
2. #6 → #7 → #8 → #9 (output.go full coverage)

**Round 2: Integration tests** (~45m)
3. #10 → #11 (env tag tests)
4. #12 (counting flag test)
5. #13 → #14 (editor + suggest tests)

### P2 — API Polish (execution order)

**Round 3: Dead code + consistency** (~30m)
6. #15 → #16 → #17 → #18 (remove 4 dead functions)
7. #19 (consistent suggest API)
8. #20 (move count handler)

**Round 4: Features & fixes** (~1h)
9. #21 (wireSubcommandSuggestions fix)
10. #23 → #24 (WithOutputFormat + wiring)
11. #22 (WithColor deprecation warning)
12. #25 (GenerateHelp tests)
13. #26 (32-bit int safety)
14. #27 (T27 CLI-in-DI fix)

### P3 — Documentation (execution order)

**Round 5: Critical docs** (~1h)
15. #28 (QUICKSTART rewrite)
16. #29 (README update)

**Round 6: Examples** (~1h)
17. #30 → #31 → #32 (env, output, counting examples)
18. #33 → #34 → #35 (DI, error, signal examples)

**Round 7: Polish docs** (~30m)
19. #36 → #37 → #38 → #39 → #40 (AGENTS, FEATURES, examples README, TODO)

### P4 — New Features (execution order)

**Round 8: Config file loading** (~1.5h)
20. #41 → #42 → #43 → #44 → #45 → #46 → #47 → #48 (koanf integration)

**Round 9: Interactive prompts** (~1h)
21. #49 → #50 → #51 → #52 → #53 (huh integration)

**Round 10: Completion + Help** (~1h)
22. #54 → #55 → #56 → #57 (shell completion)
23. #58 → #59 (markdown help)

**Round 11: Polish features** (~45m)
24. #60 → #61 (spinner/progress)
25. #62 (telemetry middleware)

### P5 — Infrastructure (execution order)

**Round 12: Critical infra** (~30m)
26. #63 (fix go-output dep — BLOCKS CI)
27. #64 (release tag)
28. #67 (fix pre-commit hooks)

**Round 13: CI + tooling** (~1h)
29. #65 (release automation)
30. #66 (codecov)
31. #68 → #70 (benchmarks)
32. #69 (merge test helpers)
33. #71 → #72 (man pages, plugin system)

---

## Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| go-output local replace blocks CI | 🔴 Critical | #63 must be resolved before any CI |
| T27 CLI-in-DI still broken | 🟡 Medium | #27 — needs Healthchecker impl |
| Coverage 72% (down from 82%) | 🟡 Medium | P1 brings it back to 80%+ |
| koanf adds 6+ transitive deps | 🟡 Medium | #41 — evaluate impact on binary size |
| huh/bubbles add heavy TUI deps | 🟠 Low | P4 — optional, can defer |
| Pre-commit hooks have errors | 🟠 Low | #67 — existing, not regression |

---

_This plan covers ALL remaining TODOs from: status report §E + §F, TODO_LIST.md, roadmap, dead code audit, coverage gaps. No items were excluded._
