# TODO List

**Updated:** 2026-06-03
**Status:** v2.4.0 — release-ready

## Completed

### Phase 1–9: All Complete

- [x] All core features implemented and tested
- [x] All architecture hardening complete
- [x] All documentation updated
- [x] Nix flake with devShell, formatter, and format check
- [x] CI with pinned golangci-lint, codecov, Nix check, benchmarks
- [x] Release automation workflow
- [x] 357 tests, 82.8% coverage, 0 lint issues, 0 race conditions

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
