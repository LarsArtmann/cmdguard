# TODO List

**Updated:** 2026-06-08
**Status:** v2.4.0 — post-release maintenance

## Completed

### Phase 1–9: All Complete

- [x] All core features implemented and tested
- [x] All architecture hardening complete
- [x] All documentation updated
- [x] Nix flake with devShell, formatter, and format check
- [x] CI with pinned golangci-lint, codecov, Nix check, benchmarks
- [x] Release automation workflow

### Phase 10: Post-Release Maintenance (2026-06-08)

- [x] Fix `flake.nix` infinite recursion (`goPkg = goPkg` → `pkgs.go_1_26`)
- [x] Update FEATURES.md version from v2.3.0-dev to v2.4.0
- [x] Fix test count metrics (357, not 356)
- [x] Add gofumpt and goimports to flake.nix treefmt
- [x] Add `Scope.HealthCheckResults()` / `HealthCheckResultsWithContext()`
- [x] Add `CLI.HealthCheckResults()` / `HealthCheckResultsWithContext()`
- [x] Add `DoctorCommand[T]` / `MustDoctorCommand[T]` convenience helper
- [x] DRY configload: extract `genericLoader` (3 files → 1)
- [x] Add configload tests: YAML, TOML, JSON, Auto, LoaderForPath (22 tests)
- [x] Consolidate `command_suggest.go` into `flags_suggest.go`
- [x] Update taskctl example: manual health → DoctorCommand
- [x] Update docs: FEATURES.md, TODO_LIST.md, AGENTS.md

## Remaining Work

### CI/CD

- [ ] Add `CODECOV_TOKEN` secret to GitHub repo settings (required for codecov upload)

### Future (v3.0+)

- [ ] Plugin system for custom validators and type handlers
- [ ] Config file nested struct support

### Future Cleanup (API-breaking, defer to v3.0)

- [ ] Make NoFlags a distinct named type (not type alias)
- [ ] Rename Get[T]/MustGet[T] to more specific names
- [ ] Make RegisterInScope generic instead of `...any`
- [ ] Remove or redesign Package() for error-safe DI integration
