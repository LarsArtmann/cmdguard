# TODO List

**Updated:** 2026-06-01
**Status:** v2.3.0-dev — release-ready

## Completed

### Phase 1–9: All Complete (see CHANGELOG.md for v2.2.0 and v2.3.0 details)

- [x] All core features implemented and tested
- [x] All architecture hardening complete
- [x] All documentation updated
- [x] Nix flake with devShell, formatter, and format check
- [x] CI with pinned golangci-lint, codecov, Nix check, benchmarks
- [x] Release automation workflow
- [x] 333 tests (1084 cases), 83.6% coverage, 0 lint issues, 0 race conditions

## Remaining Work

### CI/CD

- [ ] Add `CODECOV_TOKEN` secret to GitHub repo settings (required for codecov upload)

### Future (v3.0+)

- [ ] Plugin system for custom validators and type handlers

### Future Cleanup (API-breaking, defer to v3.0)

- [ ] Make NoFlags a distinct named type (not type alias)
- [ ] Change TimingMiddleware callback to include error
- [ ] Remove string-based BranchWithTimeout/BranchWithDeadline (replaced by typed alternatives)
- [x] Remove FlowContextAccessor (thin wrapper with no added value)
- [ ] Rename Get[T]/MustGet[T] to more specific names
- [ ] Make RegisterInScope generic instead of `...any`
- [ ] Remove or redesign Package() for error-safe DI integration
