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

<!-- check-rows:ledger -->

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

<!-- check-rows:ledger -->

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

| Item                                                                                                                      | Source                | Priority            |
| ------------------------------------------------------------------------------------------------------------------------- | --------------------- | ------------------- |
| ~~Add `CODECOV_TOKEN` secret to GitHub~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.** | ~~TODO_LIST~~         | ~~Low (CI config)~~ |
| ~~Plugin system for custom validators~~ done — shipped — plugin.go, RegisterPlugin/WithPlugin (FEATURES.md)               | ~~TODO_LIST~~         | ~~v3.0~~            |
| ~~Config file nested struct support~~ done — nested config structs supported (AGENTS.md)                                  | ~~TODO_LIST~~         | ~~v3.0~~            |
| ~~v3.0 API-breaking cleanup (4 items)~~ **Superseded — v4.0.0 shipped the breaking cleanup (2026-07-28).**                | ~~TODO_LIST/ROADMAP~~ | ~~v3.0~~            |
| ~~Fuzz testing (flags_parse, config_parsing)~~ done — 8 fuzz targets with testdata/fuzz seeds (AGENTS.md)                 | ~~ROADMAP~~           | ~~Medium~~          |
| ~~Documentation generator~~ done — GenerateDocs/docgen.go ships (FEATURES.md Documentation Generation)                    | ~~ROADMAP~~           | ~~Low~~             |
| ~~koanf integration for config auto-loading~~ done — KoanfLoader in koanf_loader.go (AGENTS.md)                           | ~~ROADMAP~~           | ~~Low~~             |
| ~~Structured JSON error output~~ done — cli_errors_json.go structured JSON error output                                   | ~~ROADMAP~~           | ~~Low~~             |
| ~~Contribution guide~~ done — CONTRIBUTING.md exists at repo root                                                         | ~~ROADMAP~~           | ~~Low~~             |

---

## d) TOTALLY FUCKED UP

**Nothing.** All changes compile, pass tests, pass lint, and preserve exact same coverage.

---

## e) WHAT WE SHOULD IMPROVE

### Codebase Health

1. ~~**Remaining 9 t30 clones are all idiomatic test patterns** — BDD scenarios with different hooks, validation tests with different constraints, and example code. These are structurally similar but semantically distinct. Further extraction would hurt readability.~~ **NOT-DO — accepted as intentional similarity — art-dupl reports 0 semantic clone groups today.**

2. ~~**FEATURES.md is stale** — Still says "v2.3.0-dev" in header, should be "v2.4.0"~~ done — FEATURES.md rewritten repeatedly since; refreshed again in the 2026-11-07 docs-health pass

3. ~~**configload/ package has 0% coverage** — The optional YAML/TOML loaders have no tests~~ **Superseded — configload deleted; KoanfLoader is the single loader (AGENTS.md).**

4. ~~**No fuzz tests** — Flag parsing and config parsing are security-adjacent (user input). Fuzz tests would be valuable.~~ done — 8 fuzz targets with testdata/fuzz seeds (AGENTS.md)

5. ~~**pre-commit hooks broken** — `git commit --no-verify` required (documented in AGENTS.md)~~ done — replaced by the BuildFlow pre-commit hook (.githooks/pre-commit)

### Process

6. ~~**No CODECOV_TOKEN** — Coverage can't be uploaded to codecov.io from CI~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.**
7. ~~**flake.nix limited** — Only devShell + formatter, no `buildGoModule` or vet checks~~ **NOT-DO — AGENTS.md design decision — flake.nix has no build/vet checks.**

---

## f) Top #25 Things We Should Get Done Next

### High Impact — Code Quality

1. ~~**Add fuzz tests for `flags_parse.go`** — ParseFlagTags and flag value parsing handle arbitrary user input~~ done — FuzzParseFlagTags + testdata/fuzz seeds exist
2. ~~**Add fuzz tests for `config_parsing.go`** — DefaultValue and tag parsing edge cases~~ done — FuzzSetField/FuzzParseFlagTags cover tag + default parsing
3. ~~**Add tests for `configload/` package** — Currently 0% coverage, optional YAML/TOML loaders untested~~ **Superseded — configload deleted; KoanfLoader is the single loader (AGENTS.md).**
4. ~~**Fix pre-commit hooks** — Currently requires `--no-verify`, should work cleanly~~ done — replaced by the BuildFlow pre-commit hook (.githooks/pre-commit)
5. ~~**Add `buildGoModule` to `flake.nix`** — Proper nix build instead of just devShell~~ **NOT-DO — AGENTS.md design decision — flake.nix has no build/vet checks.**

### High Impact — v3.0 Preparation

6. ~~**Write v3.0 API design document** — Plan breaking changes before implementing~~ done — docs/planning/2026-07-06_06-54_v3-superb-cli-redesign.md; v3.0.0 shipped 2026-07-07
7. ~~**Make `NoFlags` a distinct named type** — Not `type NoFlags = struct{}`, break alias~~ done — NoFlags is a distinct named type in v4 (AGENTS.md gotcha)
8. **Rename `Get[T]`/`MustGet[T]`** — Too generic for DI scope methods **→ ROADMAP v5 candidate #6; MustGet no longer exists in v4**
9. **Make `RegisterInScope` generic** — Replace `...any` with proper type parameters **→ TODO_LIST T5 (TODO(v5), scope.go)**
10. **Remove or redesign `Package()`** — Error-safe DI integration needs rethinking **→ ROADMAP v5 candidate #7**

### Medium Impact — Documentation

11. ~~**Update FEATURES.md header** — Change "v2.3.0-dev" → "v2.4.0"~~ done — FEATURES.md maintained through v4; refreshed 2026-11-07
12. ~~**Add contribution guide** — `CONTRIBUTING.md` with PR process, code style, testing requirements~~ done — CONTRIBUTING.md exists at repo root
13. ~~**Create v3.0 migration guide outline** — `docs/MIGRATION_V2_TO_V3.md` skeleton~~ done — docs/MIGRATION_v2_v3.md exists
14. ~~**Add issue/PR templates** — `.github/ISSUE_TEMPLATE/`, `.github/PULL_REQUEST_TEMPLATE.md`~~ done — .github/ISSUE_TEMPLATE/ + .github/PULL_REQUEST_TEMPLATE.md exist
15. ~~**Add godoc examples** — Runnable examples for core API (NewCLI, NewCommand, etc.)~~ done — pkg/cmdguard/v4/example_test.go godoc examples

### Medium Impact — Features

16. ~~**Plugin system for custom validators** — `validate:"custom"` with registry~~ done — Plugin system shipped (plugin.go; FEATURES.md)
17. ~~**Config file nested struct support** — Current loaders only handle flat keys~~ done — nested config structs supported (AGENTS.md)
18. ~~**Structured JSON error output** — `--output=json` should format errors as JSON~~ done — cli_errors_json.go structured JSON error output
19. ~~**Enhanced flag validation enums** — Beyond `validate:"email,min=5"` to richer constraints~~ done — v4 validation supports values/min/max/minlen/maxlen/dmin/dmax + custom validators
20. **Add `Result[T]` type** — Rust-like Result monad for error handling **→ ROADMAP raw ideas (Result[T] monad)**

### Lower Impact — CI/Infra

21. ~~**Add CODECOV_TOKEN secret** — Enable coverage upload from GitHub Actions~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.**
22. ~~**Test all examples in CI** — Verify examples/ compile and test in CI pipeline~~ **Superseded — superseded by the 2026-11-07 harvest — CI examples job re-tracked in TODO_LIST.**
23. **Extract flag-related code to `flagtags`** — Standalone library for struct tag parsing **→ ROADMAP raw ideas (extract flagtags standalone library)**
24. **Add metrics/hooks for custom observability** — Beyond OpenTelemetry spans **→ ROADMAP raw ideas (metrics/hooks observability)**
25. ~~**Deprecate v1 API timeline** — Announce removal schedule~~ **Superseded — the v1 package no longer exists — pkg/ contains only v4 (v4.0.0 module migration).**

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

## Resolution (2026-07-18)

The "Should we start v3.0 development now?" open question (§g) was answered yes: v3.0.0 shipped 2026-07-07. This dedup work targeted `pkg/cmdguard/v2/` files; the v3 redesign restructured much of the surrounding shape. Coverage 82.8% → 87.6%.
