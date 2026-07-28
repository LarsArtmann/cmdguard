# TODO List

**Updated:** 2026-07-27
**Status:** v3.2.0+ — 0 lint issues, 0 panics, 87.8% coverage
**Tests:** 467 test functions, 26 benchmarks, 7 fuzz targets
**Lint:** 0 issues
**Sub-modules:** 4 optional modules (glamour, prompts, spinner, telemetry)

> Open work only. Completed items live in `CHANGELOG.md` under `[Unreleased]`.

---

## P0 — Blocks shipping

- [ ] **#43** Fix `go.work` go-output pollution — the workspace contains 13
      `use` directives pointing at `/home/lars/projects/go-output/*` (local
      absolute paths). The repo **cannot build on any other machine or in CI**
      without a sibling clone of go-output. Caused by go-output's unpublished
      sub-module pseudo-versions. Options: (a) add `replace` directives in
      `go.mod` for local dev, (b) accept a local-clone requirement and document
      it, (c) fix go-output publishing. Verify with `GOWORK=off go build ./...`
      after the fix. _Source: 2026-07-27 implementation report §d.1 — verified
      still open 2026-07-27._

---

## Deferred (requires API-breaking semver bump or external access)

- [ ] **#6** Add flake.nix sub-module builds — needs Nix expertise for multi-module build
- [ ] **#10** Middleware context propagation — changes `Middleware[T]` func signature (v3.1+)
- [ ] **#15-18** API renames (`Get`→`GetService`, `RegisterInScope` generic, `Package` redesign, `SetConfig` removal) — v3.1+
- [ ] **#23** Second example app — low ROI (2h)
- [ ] **#26** `CODECOV_TOKEN` secret — requires GitHub repo owner access
- [ ] **#28** Fuzz test corpus expansion — low priority (7 targets exist, seeds added but minimal)
- [ ] **#30** `gopls` infertypeargs sweep — cosmetic (~5 remaining info-level diagnostics)

---

## Config loading consolidation follow-ups

These were surfaced by the 2026-07-27 config-loading consolidation
(`docs/status/2026-07-27_01-37_*`). The consolidation itself shipped
(commit `e3e710c`); these are the residual open items.

- [ ] **#44** KoanfLoader edge-case test coverage — koanf→JSON round-trip is
      untested for TOML datetimes, YAML anchors/aliases, deeply nested structs,
      and `int` vs `float64` number preservation. Add targeted tests to
      `koanf_loader_test.go`.
- [ ] **#45** Benchmark regression run — KoanfLoader adds a koanf→JSON marshal
      step vs the old direct `jsonLoader`. Run `go test ./benchmarks/... -bench=.`
      and compare; no baseline was captured before the consolidation.
- [ ] **#46** Review `WithConfigFileLoader` escape hatch (YAGNI) —
      `WithConfigFile` now auto-detects JSON/YAML/TOML, so the custom-loader
      escape hatch may have no remaining use case. Decide: keep (zero-cost
      escape hatch) or remove (smaller API). Removing is a breaking change.
- [ ] **#47** Remove dead `ConfigFileLoader` ireturn allow-list entry —
      `.golangci.yml:255` allow-lists `ConfigFileLoader`, but nothing returns
      the interface anymore (`NewKoanfLoader` returns concrete `*KoanfLoader`;
      `WithConfigFileLoader` only _accepts_ the interface as a parameter).
      Safe to remove; update AGENTS.md exclusion count (4→3 ireturn entries)
      and ADR-002 if proceeding.
- [ ] **#48** Consider making `loadConfigFile` private-only — it is only used
      by the custom-loader fallback path in `loadConfigFileOrSkip`. If custom
      loaders are dropped (#46), this becomes pure internal plumbing.

---

## Notes

- All metrics above are verified against the current codebase: `go test ./...`
  reports 467 `func Test` declarations across all modules (v3 core + sub-modules
  - examples + benchmarks); coverage is 87.8% in `pkg/cmdguard/v3`; lint reports
    0 issues; all 4 sub-modules build cleanly.
- `manpage` was removed from the workspace in commit `34a0c6e`; the remaining
  optional sub-modules are glamour, prompts, spinner, and telemetry.
- The `configload` sub-package was deleted entirely (commit `e3e710c`); config
  loading is now a single `KoanfLoader` in the `v3` package. See `CHANGELOG.md`
  `[Unreleased]`.
- Long-term / unbounded ideas live in `ROADMAP.md` — they are NOT duplicated
  here. When a ROADMAP idea becomes bounded and short-term, it graduates into
  this file.
- Done items are **not** listed here — they are recorded in `CHANGELOG.md`
  under the relevant version.
