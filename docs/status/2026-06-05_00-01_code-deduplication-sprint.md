# Status Report — Code Deduplication Sprint

**Date:** 2026-06-05 00:01
**Branch:** master
**Base:** v2.4.0 (5b12bf5)
**Session Focus:** Semantic code deduplication via `art-dupl`

---

## Executive Summary

Ran `art-dupl --semantic` at threshold 30 and 50, analyzed all 16 clone groups, and eliminated 7 through meaningful refactoring. Net result: **-170 lines**, cleaner test helpers, one production-code deduplication, and zero regressions.

---

## a) FULLY DONE

### Deduplication Refactoring (7 files, 7 clone groups eliminated)

| File                                              | Change                                                         | Clones Eliminated                |
| ------------------------------------------------- | -------------------------------------------------------------- | -------------------------------- |
| `pkg/cmdguard/v2/type_handler.go`                 | Extracted `registerStringFlagFromTag()` helper                 | 3 production-code clones         |
| `pkg/cmdguard/v2/type_handler_custom.go`          | Used new helper for Duration + stringParseTypes loop           | -                                |
| `pkg/cmdguard/v2/type_handler_kinds.go`           | Used new helper for String kind handler                        | -                                |
| `pkg/cmdguard/v2/prompts_test.go`                 | Table-driven `TestWithPromptOnMissing_Integration`             | 1 t50 clone group                |
| `pkg/cmdguard/v2/config_file_integration_test.go` | Extracted `writeTestConfigFile()` helper                       | Reduced clone footprint          |
| `examples/taskctl/main_test.go`                   | Extracted `mustExec()` helper, simplified 21 tests             | 2 large clone groups (13 clones) |
| `pkg/cmdguard/v2/cli_superb_test.go`              | Extracted `newTestCLIWithNoOpCmd()`, simplified 10+ args tests | 4 clone groups                   |

### Quality Gates — ALL PASSING

| Metric                   | Value                     | Status           |
| ------------------------ | ------------------------- | ---------------- |
| Build                    | `go build ./...`          | ✅ 0 errors      |
| Tests                    | `go test ./... -race`     | ✅ ALL PASS      |
| Test cases (v2)          | 273                       | ✅ Unchanged     |
| Test cases (taskctl)     | 62                        | ✅ Unchanged     |
| Test cases (integration) | 17                        | ✅ Unchanged     |
| Coverage (v2)            | 82.8%                     | ✅ Unchanged     |
| Coverage (taskctl)       | 71.1%                     | ✅ Unchanged     |
| Lint                     | `golangci-lint run ./...` | ✅ 0 issues      |
| Race conditions          | `-race` flag              | ✅ 0 detected    |
| Clone groups (t50)       | 1 (was 2)                 | ✅ 50% reduction |
| Clone groups (t30)       | 9 (was 16)                | ✅ 44% reduction |

---

## b) PARTIALLY DONE

### config_file_integration_test.go (t50 clone group)

- Extracted `writeTestConfigFile()` helper to reduce boilerplate
- **Cannot fully eliminate** the 5-clone group because each subtest defines a **different local Config type** with different struct tags (`env:"TEST_NAME"`, `env:"TEST_NAME2"`, or no env tag)
- Go struct tags are compile-time constants — no generic helper can abstract over them
- **Verdict:** Idiomatic test duplication. The `writeTestConfigFile` helper is the maximum reasonable extraction.

---

## c) NOT STARTED

Items from TODO_LIST.md and ROADMAP.md that were NOT in scope for this session:

| Item                                       | Source            | Priority        |
| ------------------------------------------ | ----------------- | --------------- |
| Add `CODECOV_TOKEN` secret to GitHub       | TODO_LIST         | Low (CI config) |
| Plugin system for custom validators        | TODO_LIST         | v3.0            |
| Config file nested struct support          | TODO_LIST         | v3.0            |
| v3.0 API-breaking cleanup (4 items)        | TODO_LIST/ROADMAP | v3.0            |
| Fuzz testing (flags_parse, config_parsing) | ROADMAP           | Medium          |
| Documentation generator                    | ROADMAP           | Low             |
| koanf integration for config auto-loading  | ROADMAP           | Low             |
| Structured JSON error output               | ROADMAP           | Low             |
| Contribution guide                         | ROADMAP           | Low             |

---

## d) TOTALLY FUCKED UP

**Nothing.** All changes compile, pass tests, pass lint, and preserve exact same coverage.

---

## e) WHAT WE SHOULD IMPROVE

### Codebase Health

1. **Remaining 9 t30 clones are all idiomatic test patterns** — BDD scenarios with different hooks, validation tests with different constraints, and example code. These are structurally similar but semantically distinct. Further extraction would hurt readability.

2. **FEATURES.md is stale** — Still says "v2.3.0-dev" in header, should be "v2.4.0"

3. **configload/ package has 0% coverage** — The optional YAML/TOML loaders have no tests

4. **No fuzz tests** — Flag parsing and config parsing are security-adjacent (user input). Fuzz tests would be valuable.

5. **pre-commit hooks broken** — `git commit --no-verify` required (documented in AGENTS.md)

### Process

6. **No CODECOV_TOKEN** — Coverage can't be uploaded to codecov.io from CI
7. **flake.nix limited** — Only devShell + formatter, no `buildGoModule` or vet checks

---

## f) Top #25 Things We Should Get Done Next

### High Impact — Code Quality

1. **Add fuzz tests for `flags_parse.go`** — ParseFlagTags and flag value parsing handle arbitrary user input
2. **Add fuzz tests for `config_parsing.go`** — DefaultValue and tag parsing edge cases
3. **Add tests for `configload/` package** — Currently 0% coverage, optional YAML/TOML loaders untested
4. **Fix pre-commit hooks** — Currently requires `--no-verify`, should work cleanly
5. **Add `buildGoModule` to `flake.nix`** — Proper nix build instead of just devShell

### High Impact — v3.0 Preparation

6. **Write v3.0 API design document** — Plan breaking changes before implementing
7. **Make `NoFlags` a distinct named type** — Not `type NoFlags = struct{}`, break alias
8. **Rename `Get[T]`/`MustGet[T]`** — Too generic for DI scope methods
9. **Make `RegisterInScope` generic** — Replace `...any` with proper type parameters
10. **Remove or redesign `Package()`** — Error-safe DI integration needs rethinking

### Medium Impact — Documentation

11. **Update FEATURES.md header** — Change "v2.3.0-dev" → "v2.4.0"
12. **Add contribution guide** — `CONTRIBUTING.md` with PR process, code style, testing requirements
13. **Create v3.0 migration guide outline** — `docs/MIGRATION_V2_TO_V3.md` skeleton
14. **Add issue/PR templates** — `.github/ISSUE_TEMPLATE/`, `.github/PULL_REQUEST_TEMPLATE.md`
15. **Add godoc examples** — Runnable examples for core API (NewCLI, NewCommand, etc.)

### Medium Impact — Features

16. **Plugin system for custom validators** — `validate:"custom"` with registry
17. **Config file nested struct support** — Current loaders only handle flat keys
18. **Structured JSON error output** — `--output=json` should format errors as JSON
19. **Enhanced flag validation enums** — Beyond `validate:"email,min=5"` to richer constraints
20. **Add `Result[T]` type** — Rust-like Result monad for error handling

### Lower Impact — CI/Infra

21. **Add CODECOV_TOKEN secret** — Enable coverage upload from GitHub Actions
22. **Test all examples in CI** — Verify examples/ compile and test in CI pipeline
23. **Extract flag-related code to `flagtags`** — Standalone library for struct tag parsing
24. **Add metrics/hooks for custom observability** — Beyond OpenTelemetry spans
25. **Deprecate v1 API timeline** — Announce removal schedule

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we start v3.0 development now, or continue polishing v2.x?**

The TODO_LIST and ROADMAP both show v2.4.0 as "release-ready" with all planned v2.x work complete. The remaining v2.x items (codecov token, fuzz tests) are improvements, not features. But v3.0 involves significant API-breaking changes that require careful design.

**What I need to know:** Is there a v2.5 milestone planned (e.g., fuzz tests + configload coverage + plugin system), or should we move directly to v3.0 API design? This determines whether the next sprint focuses on polish or on the new API surface.

---

## Metrics Snapshot

```
Clone groups at t50: 1 (was 2) — 50% reduction
Clone groups at t30: 9 (was 16) — 44% reduction
Lines changed:      +159 / -329 (net -170)
Files modified:     7
Tests passing:      352 (273 v2 + 62 taskctl + 17 integration)
Coverage:           82.8% v2 / 71.1% taskctl
Lint issues:        0
Race conditions:    0
```

---

_Generated by Crush on 2026-06-05_
