# TODO List

**Updated:** 2026-07-23
**Status:** v3.0.0+ — 0 lint issues, 0 panics, 87.8% coverage, 58 sentinel errors
**Tests:** 470 test functions (1434 runs), 26 benchmarks, 7 fuzz targets
**Lint:** 0 issues
**Sub-modules:** 4 optional modules (glamour, prompts, spinner, telemetry)

> Open work only. Completed items live in `CHANGELOG.md`.

---

## Deferred (requires API-breaking semver bump or external access)

- [ ] **#6** Add flake.nix sub-module builds — needs Nix expertise for multi-module build
- [ ] **#7** Move koanf to optional sub-module — API-breaking for `configload.KoanfLoader()` consumers
- [ ] **#10** Middleware context propagation — changes `Middleware[T]` func signature (v3.1+)
- [ ] **#12** Deduplicate `jsonLoader` in flake.nix — low priority
- [ ] **#15-18** API renames (`Get`→`GetService`, `RegisterInScope` generic, `Package` redesign, `SetConfig` removal) — v3.1+
- [ ] **#23** Second example app — low ROI (2h)
- [ ] **#26** `CODECOV_TOKEN` secret — requires GitHub repo owner access
- [ ] **#28** Fuzz test corpus expansion — low priority (7 targets exist, seeds added but minimal)
- [ ] **#30** `gopls` infertypeargs sweep — cosmetic (~5 remaining info-level diagnostics)

---

## Future Ideas (no timeline)

| #   | Task                                                           | Category         |
| --- | -------------------------------------------------------------- | ---------------- |
| 31  | Extract flag-tags to `github.com/larsartmann/flagtags`         | Refactor / reuse |
| 32  | Service-owned config design (ADR) — services own typed config  | Architecture     |
| 33  | Command-level audit middleware — audit every command execution | Feature          |
| 34  | Built-in audit-log subcommand (`myapp audit-log --format d2`)  | Feature          |
| 35  | Consider making fang optional (plain cobra fallback)           | Architecture     |
| 36  | `FlagRegistry` interface abstraction                           | API design       |
| 37  | Custom per-flag validation hooks (beyond `validate` tag)       | Feature          |
| 38  | Enhanced flag validation enums                                 | Feature          |
| 39  | Metrics/hooks for custom observability (beyond OpenTelemetry)  | Feature          |
| 40  | Branded-ID example app                                         | Example          |
| 41  | Test all examples in CI                                        | CI               |
| 42  | Benchmark regression thresholds in CI                          | CI               |

---

## Notes

- All metrics above are verified against the current codebase; `go test ./...` reports 470 test functions and 1434 pass events across all modules.
- `manpage` was removed from the workspace in commit `34a0c6e`; the remaining optional sub-modules are glamour, prompts, spinner, and telemetry.
- Done items are **not** listed here — they are recorded in `CHANGELOG.md` under the relevant version.
