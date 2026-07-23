# ROADMAP

**Updated:** 2026-07-23
**Purpose:** Aspirational items with no concrete timeline

---

## v1 Deprecation

v1 (`pkg/cmdguard/v1`) is deprecated as of v3.0.0. It receives no new features
and only critical security fixes. v1 will be removed in v4.0.0 (no earlier than
2026-12-31). Consumers should migrate to v3 via `docs/MIGRATION_v2_v3.md`.

---

## Deferred v3.x / v4 API-Breaking Cleanup

These remain open for a future v3.x or v4:

- [ ] **Rename `Get[T]`/`MustGet[T]`** — `Get` is too generic; should be
      `GetService[T]` or similar. Breaking: every consumer's import surface changes.
- [ ] **Make `RegisterInScope` generic** — currently takes `...any`; should be
      `RegisterInScope[T](scope, provider)`. Breaking: signature change.
- [ ] **Remove or redesign `Package()`** — takes a pre-existing `*Scope` which is
      an unusual API shape; should be reworked for error-safe DI. Breaking.
- [ ] **Remove `SetConfig`** — mutating a CLI's config after construction is
      unsafe (FlagRegistry isn't re-initialized). Breaking but removes a footgun.

---

## Future Ideas

- [ ] **Extract flag-tag parsing to `github.com/larsartmann/flagtags`**
- [ ] **Branded-ID example application**
- [ ] **Add `FlagRegistry` interface abstraction**
- [ ] **Custom per-flag validation hooks (beyond `validate` tag)**
- [ ] **Enhanced flag validation enums**
- [ ] **Metrics/hooks for custom observability (beyond OpenTelemetry)**
- [ ] **Test all examples in CI**
- [ ] **Benchmark regression thresholds in CI**

---

## Deferred from 2026-07-18 Audit Closure

The following items were identified during the multi-skill audit session (2026-07-18) and
explicitly deferred with rationale. They are NOT blocking; they are documented here to stop
them from haunting the TODO list.

| #   | Item                                                                                              | Deferral rationale                                                                                         |
| --- | ------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| 1   | Re-run 4 skills with reference files loaded                                                       | More reports = more debt without action layer. Re-audit in a dedicated future session AFTER closure ships. |
| 2   | Run 4 additional skills (brutal-self-review, library-deep-dive, status-report, docs-health BUILD) | New audits, not closure. Belong in a future session.                                                       |
| 3   | TypeHandler to TypeCodec rename + Deprecated alias                                                | Public API break. v4 only. `TODO(v4)` marker added at `type_handler.go:13`.                                |
| 4   | ConfigFile branded type                                                                           | YAGNI until a consumer needs it.                                                                           |
| 5   | Extract koanf into configload sub-module                                                          | YAGNI until consumer asks or core LOC > 12k. Currently ~8k.                                                |
| 6   | Split v3 into v3 + v3/internal/                                                                   | LOC trigger (12k) not met. Premature split adds boundary friction.                                         |
| 7   | Fuzz corpus expansion                                                                             | Existing 7 targets have minimal corpus. Valuable but not closure.                                          |
| 8   | Audit examples/taskctl/main_test.go (876 lines)                                                   | Test-smell audit is a separate concern.                                                                    |
| 9   | CONTRIBUTING.md refresh                                                                           | Not blocking. Verify-then-decide in a future pass.                                                         |
| 10  | Verify git-town.toml + library-policy.yaml                                                        | Config sanity, not closure.                                                                                |
| 11  | Update WHAT_THIS_PROJECT_IS_ABOUT.md + _NOT.md                                                    | Living docs; belongs in docs-health, not this closure.                                                     |
| 12  | Schedule re-run after v3.1 ships                                                                  | Re-audit after the next minor release.                                                                     |

---

## Notes

- All shipped items previously listed here have been moved to `CHANGELOG.md` under their release version.
- `manpage` was removed from the workspace in commit `34a0c6e`; optional sub-modules are now glamour, prompts, spinner, and telemetry.
