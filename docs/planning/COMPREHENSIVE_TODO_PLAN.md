# Comprehensive TODO Plan

**Generated:** 2026-06-01
**Status:** v2.3.0-dev — 272 tests, 83.3% coverage, 0 lint issues

---

## Priority Legend

| Priority | Meaning |
|----------|---------|
| **P0** | Blockers / broken right now |
| **P1** | High value, ships the release |
| **P2** | Important polish, post-release |
| **P3** | Nice-to-have, backlog |
| **v3** | API-breaking, deferred to v3.0 |

## Effort Legend

| Effort | Time |
|--------|------|
| **XS** | ≤ 5 min |
| **S** | ≤ 12 min |
| **M** | ≤ 30 min |
| **L** | ≤ 60 min |

---

## Plan

| # | Priority | Task | Effort | Category | Why |
|---|----------|------|--------|----------|-----|
| 1 | **P0** | Stale `config_file.go` has uncommitted broken change referencing nonexistent `FilterSetFields` — revert or implement | S | Bug | Build is broken on HEAD+working-tree |
| 2 | **P1** | Update CHANGELOG.md with v2.2.0 + v2.3.0 entries (currently stops at v2.1.0) | L | Docs | Release blocker — no changelog since March |
| 3 | **P1** | Create v2.3.0 release tag and GitHub release notes | M | Release | Ships the product |
| 4 | **P1** | Pin `golangci-lint` version in CI (currently `version: latest` — non-reproducible) | XS | CI | Build can break any day upstream |
| 5 | **P1** | Fix CONTRIBUTING.md outdated testing section (still references `testify`, `just`, Go 1.24) | S | Docs | Already done — verify committed |
| 6 | **P1** | Fix README.md stale Examples section (lists 13 deleted example dirs) | S | Docs | Already done — verify committed |
| 7 | **P1** | Add `.direnv/` to `.gitignore` | XS | Cleanup | Already done — verify committed |
| 8 | **P1** | Delete stale `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md` (says "Not Yet Implemented" but flake exists) | XS | Cleanup | Already done — verify committed |
| 9 | **P1** | Update AGENTS.md test count (272) and coverage (83.3%) — verify numbers match `go test -cover` | S | Docs | Stale numbers mislead |
| 10 | **P1** | Update FEATURES.md test count (says 22 benchmarks, 7 fuzz — verify) | S | Docs | Accuracy |
| 11 | **P1** | Update FEATURES.md "Last updated" date (says 2026-05-17) | XS | Docs | Stale |
| 12 | **P2** | Add codecov badge + integration to CI | M | CI | Visibility into coverage trends |
| 13 | **P2** | Add Nix CI job alongside existing GitHub Actions (Option A from proposal) | M | CI | Reproducible CI |
| 14 | **P2** | Add CLI construction benchmark (`BenchmarkNewCLI`) | S | Perf | Baseline for regression |
| 15 | **P2** | Add flag parsing benchmark (`BenchmarkParseFlags`) | S | Perf | Baseline for regression |
| 16 | **P2** | Add command execution benchmark (`BenchmarkExecute`) | S | Perf | Baseline for regression |
| 17 | **P2** | Add benchmark regression detection to CI (compare against baseline) | M | CI/Perf | Prevent silent perf degradation |
| 18 | **P2** | Set up release automation (goreleaser or manual tag-based) | M | Release | Reproducible releases |
| 19 | **P2** | Delete stale `justfile` (AGENTS.md says deprecated; flake.nix replaces it) | S | Cleanup | Dead code confuses contributors |
| 20 | **P2** | Clean up `docs/status/` — 30 historical status reports, archive or delete | S | Cleanup | Cluttered docs tree |
| 21 | **P2** | Delete `docs/dedup-report.html` (untracked artifact) | XS | Cleanup | Generated artifact in repo |
| 22 | **P2** | Fix README.md "220+ tests" → "272+ tests" and "~82% coverage" → "~83.3%" | S | Docs | Already done — verify committed |
| **23** | **P2** | Verify README.md references match actual docs (TUTORIAL.md, QUICKSTART.md, etc. all exist) | XS | Docs | Broken links hurt credibility |
| 24 | **P2** | Remove stale `.gitignore` entries for deleted example dirs (`/advanced-flags`, `/di`, etc.) | S | Cleanup | Dead entries in gitignore |
| 25 | **P2** | Add `config_file.go` FilterSetFields function OR revert working tree to committed version | S | Bug | Pre-existing uncommitted breakage |
| 26 | **P3** | Plugin system for custom validators and type handlers | L | Feature | v3.0+ roadmap |
| 27 | **v3** | Make NoFlags a distinct named type (not type alias) | S | Cleanup | API-breaking |
| 28 | **v3** | Change TimingMiddleware callback to include error | S | Cleanup | API-breaking |
| 29 | **v3** | Remove string-based BranchWithTimeout/BranchWithDeadline | S | Cleanup | API-breaking, typed alternatives exist |
| 30 | **v3** | Remove FlowContextAccessor (thin wrapper, no added value) | XS | Cleanup | API-breaking |
| 31 | **v3** | Rename Get[T]/MustGet[T] to more specific names | S | Cleanup | API-breaking |
| 32 | **v3** | Make RegisterInScope generic instead of `...any` | M | Cleanup | API-breaking |
| 33 | **v3** | Remove or redesign Package() for error-safe DI integration | M | Cleanup | API-breaking |

---

## Summary

| Priority | Count | Est. Total Time |
|----------|-------|-----------------|
| **P0** (Blocker) | 1 | 12 min |
| **P1** (Release) | 10 | ~2.5 hrs |
| **P2** (Polish) | 13 | ~3.5 hrs |
| **P3** (Backlog) | 1 | 60 min |
| **v3** (Deferred) | 7 | ~2 hrs |
| **Total** | 33 | ~8.5 hrs |

### Execution Order

1. **Fix P0** — broken build
2. **Verify P1 items 5–8** — already done in this session, confirm committed
3. **P1 items 2–4, 9–11** — docs + CI pin + changelog
4. **P1 item 3** — release tag
5. **P2 batch** — benchmarks, CI, cleanup
6. **v3 items** — defer
