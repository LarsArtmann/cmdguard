# Roadmap

> Long-term direction and raw ideas not yet refined into actionable tasks.
> Items here are NOT commitments — they're directions worth exploring.
> See [TODO_LIST.md](TODO_LIST.md) for actionable, bounded work.

---

## v5 Major Release

The v5 thesis and full workstream plan (P1 inevitable quality, P2 zero-bug by
construction, P3 fantastic UX, P4 minimum code): [docs/planning/2026-11-01_v5-end-goal-paved-road.md](docs/planning/2026-11-01_v5-end-goal-paved-road.md).
Renames are the smallest part of v5, not the point of it.

**Owner decisions (2026-11-02):**

- **Release vehicle:** enforced flag-tag validation ships as **v4.2.0 with a
  migration note** (documented behavior change); v5 stays reserved for the
  defaults + renames major.
- **Module path:** v5.0.0 lives at **`github.com/larsartmann/cmdguard/v5`**
  (standard Go semver; consumers rewrite imports via codemod). Scaffold = plan
  task M11.
- **v5 defaults:** strict validation, signals + graceful shutdown, and
  batteries all flip **on** by default; every individual opt-out
  (`WithoutStrictValidation` etc.) remains available.

The following breaking changes are tracked with `TODO(v5)` markers in source code:

| # | Change                                                                | Rationale                                                              | Source Location      |
| - | --------------------------------------------------------------------- | ---------------------------------------------------------------------- | -------------------- |
| 1 | Rename `CommandInfo` → `CommandMetadata`                              | "Info" suffix is vague; `Metadata` is precise                          | `middleware.go`      |
| 2 | Rename `TypeHandler` → `TypeCodec`                                    | "Handler" is generic; `Codec` captures parse/default dual role         | `type_handler.go`    |
| 3 | Rename `PromptRunner` → `HuhPrompter` (or similar)                    | "Runner" suffix is generic                                             | `prompts/prompts.go` |
| 4 | Remove `SetConfig(cfg)` — now `// Deprecated:`                        | Unsafe post-construction mutation without re-initializing FlagRegistry | `cli_accessors.go`   |
| 5 | Replace `RegisterInScope` with generic variant — now `// Deprecated:` | Erases provider types; use `Child` + `do.Provide` meanwhile            | `scope.go`           |

Additional v5 candidates from partially functional items:

| # | Change                                           | Rationale                                                                                                                    |
| - | ------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| 6 | Rename `Get[T]` → `GetService[T]`                | `Get` is too generic for a DI scope                                                                                          |
| 7 | Redesign `Package[T](scope, ...)`                | Unusual API shape with pre-existing `*Scope` param                                                                           |
| 8 | Unify `Middleware[T]` and `ContextMiddleware[T]` | Two signatures coexist since v4 delivered context propagation additively; v5 can collapse to the context-threading signature |

### Flag-representation consolidation (v5 scope, from ADR-002)

A flag is described four times (`FlagTag` → `flagMetadata` → `FlagDescription`
→ `pflag.Flag`); the inventory and the accepted gaps live in
[docs/design/2026-11-02_flag-representation-gap-matrix.md](docs/design/2026-11-02_flag-representation-gap-matrix.md).
v4 tolerates the split because ADR-002 pins one producer per fact and the
drift gate enforces the docgen hop. v5 collapses the plumbing:

**Target model.**

- One registration-time `flagMeta` value per flag, produced once from struct
  tags, stored in the cobra annotation as the single durable copy.
- `FlagDescription` becomes a thin view: contract facts from `flagMeta`,
  runtime facts (`Deprecated`, `Hidden`) from pflag.
- `FlagTag` stays the parser's output type but stops being re-read after
  registration (today `config explain` and `ConfigKeys` re-parse tags from
  the struct; they would consume the same `flagMeta`).
- The config-file vocabulary (`v4.ConfigKey`) remains a separate type — file
  shapes (nested objects, indexed arrays) are not flag shapes — but shares
  the canonical key.
- The contract's `Type` field is revisited: pflag's value-type string today,
  with the JSON type exposed alongside so schema consumers stop guessing.

**Breaking changes for v5.**

- `FlagTag.Type reflect.Type` leaves the public struct (leaky implementation
  detail; consumers use the string contract).
- Strict config keys flip to default-on (decision recorded 2026-11-02;
  `WithStrictConfigKeys()` stays for discoverability; a `Without*` opt-out
  ships with it).
- `flagMetadata` JSON may be reshaped freely — it is internal annotation
  payload, not a consumer surface (ADR-002 §1 applies to `FlagDescription`,
  not to annotations).

**Migration and effort.** Consumers of the JSON contract are unaffected
(additive rule). Go-API consumers touching `FlagTag.Type` are rare (the lint
corpus showed none); a codemod replaces it with the contract string.
Estimate: 2–3 days including the strict-default flip, migration note, and
consumer re-verification on erraudit + timesheets.

---

## Architectural Directions

### Middleware Context Propagation

**Delivered additively in v4** (2026-10-04): `ContextMiddleware[T]` threads a
`next func(context.Context) error` through the chain, wired via
`WithContextMiddleware[T]`; context middleware run outside plain middleware so
derived contexts (timeouts, cancellation, tracing scopes) reach every inner
middleware and the handler. `TimeoutMiddleware[T]` ships as the reference
implementation (`ErrCommandTimeout`). Remaining v5 work: unify `Middleware[T]`
and `ContextMiddleware[T]` into one signature (breaking).

### Command-Level Audit Middleware

The upstream recording API is implemented in ../samber-do-auditlog
(`Plugin.RecordCommand` / `EventTypeCommand`, Unreleased). Blocked on that
library being pushed + tagged (≥ v0.11.0); then bump cmdguard's go.mod and add
`AuditMiddleware[T]` capturing command execution events (start, end, duration,
error). See TODO_LIST.md "Blocked on Upstream".

### Internal Package Split

If `pkg/cmdguard/v4` grows beyond 12k LOC, consider splitting into `v4` +
`v4/internal/` to separate public API surface from implementation. Current
size is below the trigger threshold — this is YAGNI until growth demands it.

---

## Deferred from 2026-07-18 Audit Closure

These items were explicitly deferred during the post-audit closure to prevent
the long tail from haunting the TODO list. Each has a rationale for deferral.

| #  | Item                                                                                              | Rationale                                                                                             |
| -- | ------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| 1  | Re-run 4 skills with reference files loaded                                                       | More reports without action = more debt. Re-audit in a dedicated session after closure ships.         |
| 2  | Run 4 additional skills (brutal-self-review, library-deep-dive, status-report, docs-health BUILD) | New audits, not closure of existing findings.                                                         |
| 3  | `TypeHandler` / `TypeHandlerFunc` rename + `// Deprecated` alias                                  | Public API break — deferred to v5. See TODO(v5) markers above.                                        |
| 4  | `ConfigFile` branded type                                                                         | YAGNI until a consumer needs it.                                                                      |
| 5  | Extract koanf into a sub-module                                                                   | YAGNI — no consumer has asked; 4 koanf deps are already isolated. LOC trigger not met.                |
| 6  | Split `v4` into `v4` + `v4/internal/`                                                             | LOC trigger (12k) not met. Premature split adds boundary friction.                                    |
| 7  | Fuzz corpus expansion                                                                             | Existing 8 targets have minimal corpus; valuable but not blocking.                                    |
| 8  | Audit `examples/taskctl/main_test.go` (876 lines)                                                 | Test-smell audit is a separate concern.                                                               |
| 9  | `CONTRIBUTING.md` refresh                                                                         | Not blocking; verify-then-decide in a future pass.                                                    |
| 10 | Verify `git-town.toml` + `library-policy.yaml`                                                    | Config sanity, not closure.                                                                           |
| 11 | ~~Update `WHAT_THIS_PROJECT_IS_ABOUT.md` + `_NOT.md`~~                                            | ✅ Done 2026-08-05 — updated to v4 API.                                                               |
| 12 | Schedule re-run after v3.1 ships                                                                  | v3.1 has shipped; v4.0.0 has shipped. Next re-audit window: after v5 or when feature growth warrants. |

---

## Raw Ideas (Not Refined)

These are unformed ideas that may or may not become features:

### Core Library

- **Config file watching** — hot-reload config on file change (inotify/fsnotify)
- **CUE for the machine contract** — revisit only when a consumer needs to VALIDATE the contract itself (not config files) against a schema, AND cue's Go codegen round-trips cmdguard's struct tags losslessly, AND the build-time cost is measured. Analysis: `docs/status/2026-11-02_04-52_cue-adoption-analysis-and-self-review.md`; decision record: ADR-002 §6. The struct tags are the single declaration surface today — a parallel CUE schema would be a second source of truth until all three criteria hold.
- **Plugin marketplace** — community-contributed type handlers and validators
- **gRPC middleware sub-module** — command-level gRPC tracing
- **Web-based CLI preview** — render command tree as interactive HTML
- **Shell completion v2** — richer dynamic completion with type-aware suggestions
- **Benchmark dashboard** — track performance across releases
- **Benchmark CI gating** — detect performance regressions automatically by comparing benchmark results against a baseline in CI
- **v3→v4 migration guide** — ~~dedicated `docs/MIGRATION_v3_v4.md`~~ done (2026-08-06)

### Harvested from 2026-0X status reports (2026-11-07 docs-health annotate pass)

Recurring ideas that surfaced across the June–September reports and were routed
here instead of re-minted as TODO rows. Each bullet is the routing target its
`→ ROADMAP` annotation points at.

- **Structured logging adapter (slog)** — core-side slog integration beyond the flightrecorder handler
- **charmbracelet/log integration** — structured logging for middleware/DI internals
- **Observability hooks (metrics, custom tracing)** — middleware-level metric hooks beyond OpenTelemetry spans
- **Metrics/hooks observability** — same class as above; kept as one idea
- **`Mark(err, sentinel)` identity stamping** — cockroachdb-style error marking as an alternative to `RegisterClassification`
- **Family-specific error format constructors** — `NewRejectionf`-style helpers (go-error-family adjacent)
- **`Compose` returning the worst error Family** — richer than `errors.Join`
- **`Result[T]` monad** — Rust-style result type (existed in v2.8, removed in the v3/v4 slim-down; revive only on demand)
- **`Validated[T]` wrapper** — validation-carrying type, same revival bar as `Result[T]`
- **Typed slice flag handlers** — `[]int`, `[]Port` beyond the `[]string`/raw-cobra escape hatch
- **Typed flow-context keys** — `BranchingFlowContext` `SetValue` takes `any`; generic keys would prevent runtime type mismatches
- **`Validator[T]` composition API** — `.And()`/`.Or()` combinators instead of string-parsed `validate` params
- **Lazy `Branch()` copy-on-write** — replace `maps.Clone` per branch with parent-lookup (pathological-tree perf only)
- **Remove the `setStringField` Enum registry bypass** — config decode special-cases `Enum` outside the TypeHandler registry (`config_setfield.go`)
- **`PortFromInt` accepting `uint16`** — still takes `int` for backward compat (`types_port.go`)
- **`iter.Seq` output streaming** — streaming rows into `OutputTable`
- **Output format snapshot tests** — golden files for the 16 formats
- **Property-based type-parser tests** — round-trip properties for Duration/URL/Email/Port parsers
- **Version embedding via `runtime/debug.ReadBuildInfo`/`-ldflags`** — accurate `--version` without a released-tag string
- **Telemetry shutdown span integration** — OpenTelemetry span around graceful shutdown
- **Telemetry → `ContextMiddleware[T]` migration** — documented future work in AGENTS.md (span-context propagation)
- **Branded-ID example app** — showcase app built on go-branded-id
- **Plugin-based custom type handler example** — runnable example for `RegisterTypeHandler`
- **Extract a `flagtags` standalone library** — struct-tag parsing as its own module (~2000 LOC); revisit only with a real second consumer
- **Test all examples in CI** — dedicated CI job building/testing `examples/` (harvested 2026-11-07 into TODO_LIST; this bullet is the long-term home)

Harvested from the June-10/11 output-sprint reports (2026-11-07 sweep, second batch):

- **`RenderOptions` passthrough** — `OutputConfig` is still `{Format, Writer}` (`pkg/cmdguard/v4/output.go`); Title/GraphID/ColorMode passthrough was proposed three times and never tracked
- **Command-scoped middleware** — v4 ships only CLI-wide `Middleware[T]`/`ContextMiddleware[T]`; per-command chains (e.g. spinner on one command) remain impossible (`WithCommandMiddleware` existed in the v2 line and was not carried forward)
- **`docs/TEST_HELPERS.md`** — inventory of test helpers spread across 7+ v4 test files
- **Codecov badge in README** — coverage uploads tokenless; badge is a one-line README change
- **govulncheck/gosec in CI** — both are manual-only today (see TODO_LIST S2 for the toolchain blocker)
- **art-dupl/dedup enforcement in CI** — clone-count regression protection existed only as a manual ritual

Harvested from the July planning/review docs (2026-11-07 sweep, third batch):

- **go.work portability** — the workspace pins 13 absolute `/home/lars/projects/go-output/*` paths, so the repo builds only on this machine; fine for local dev (AGENTS.md documents the workspace), but a fresh clone needs `go.work` trimming before it builds
- **`cmdguard/lite` stripped sub-package** — <10 transitive modules, stdlib-`flag` compatible (from the 2026-07-16 upd evaluation)
- **Single-command fast path `Run[T,F](...)`** — no `NewCLI` + `AddCommand` ceremony for one-command CLIs (upd evaluation)
- **`PrintHelp(flagSet, ...)` incremental adoption** — render help for a plain `flag.FlagSet` as an entry point (upd evaluation)
- **Website signature element** — the hero's identifying visual (before/after split, terminal-style error, construction animation), explicitly deferred to "a future design sprint" (2026-07-18 frontend-design review Finding 3)

Harvested from the early-July website/USP reports (2026-11-07 sweep, routed as "port X to website" / "marketing&infra"):

Website guide ports:

- **Port docs/TUTORIAL.md to the website**
- **Port docs/COMPARISON.md to the website**
- **Port docs/PERFORMANCE.md to the website**
- **Port docs/CLI_DESIGN_PRINCIPLES.md to the website**
- **Port docs/MIGRATION_v2_v3.md to the website**
- **Port docs/COBRA_FOOTGUNS.md to the website**
- **Doctor-Version guide** on the website (`WithDoctorCheck` + `VersionCommand` patterns)
- **Signal-handling guide** on the website (`WithSignalHandling`/`WithGracefulShutdown` patterns)

Launch marketing & infra:

- **Logo mark** for the project (shield-style, matches the hero)
- **OG image generation** (`astro-og-canvas`) for social previews
- **Lighthouse CI** (`lighthouse.yml` + `lighthouserc.json`)
- **Firebase preview channels** on PRs
- **DNS** for the appspot (`cmdguard.web.app` alias) and apex→www redirect
- **`FIREBASE_SERVICE_ACCOUNT`** in GitHub secrets (replaces `firebase login` in CI)
- **Uptime monitoring** for `cmdguard.web.app` and `cmdguard.lars.software`
- **`SECURITY.md`**
- **`.github/FUNDING.yml`** (GitHub Sponsors)
- **`CITATION.cff`** for academic citations
- **Sitemap `lastmod`** entries
- **`getting-started/full-tutorial.mdx`** end-to-end walkthrough
- **`dedup.sh`** (jscpd) duplication check for the website
- **Launch blog post** + social announcements (r/golang, HN, X/LinkedIn)
- **Awesome Go / awesome-go-cli submissions** + Go project wiki listing
- **Sub-module badges** in README (version badges per sub-module)
- **Website flake app** (`nix run .#website`) wrapping the pnpm build

Documentation & CI guards (routed as "guards/note" and "stale-doc guard"):

- **check-dep-versions guard** — CI script comparing docs' dependency tables against `go.mod`
- **dep-bump doc-sync discipline** — every dependency PR updates FEATURES/AGENTS in the same change
- **Derive dependency tables from go.mod** — single source instead of hand-maintained tables
- **Simplify the go-output note** — the indirect-deps footnote in FEATURES is stale-prone; derive it
- **Audit-log subcommand** — built-in opt-in `audit-log export` command (consumers hand-roll it today)
- **Service-owned config ADR** — document why DI services (not the CLI config struct) own connection config
- **Per-command flag scoping (`WithProvider`)** — scoped DI overrides per command tree
- **Strict flag audit mode** — CI mode that fails on flags lacking env/default provenance
- **Post-rewrite semantic doc audit** — a docs-health pass after every formatter rewrite (the 2026-11-05 table-strike lesson)
- **Stale-doc-ref CI guard** — grep for references to deleted/moved docs in living docs
- **Semver CI check** — enforce tag/CHANGELOG/module-path coherence before release
- **godoclint** — godoc convention linting across exported symbols
- **`docs/STATUS.md` index** — one-page pointer to current status reports (reports themselves stay historical)
- **docs-health CI gate** — run `docs-gates` (or equivalent) on docs-touching PRs
- **markdownlint + lychee** — markdown style + link checking in CI

Other:

- **musttag compliance evaluation** for sub-module config structs
- **No-op `PromptConfirm` test** — rewrite or delete the identified no-op (TODO_LIST H2)
- **CONTRIBUTING expansion** — doc-verification section, testing-patterns section, release-cutting section, GOCOVERDIR coverage notes
- **goat ASCII lifecycle diagram** — ASCII sequence render of the CLI lifecycle (recurred in three 2026-08-06 reports)
- **nix-flake-check badge** in README
- **`WithCleanup` benchmark** — measure the cleanup-hook path
- **Profile `examples/taskctl`** — pprof pass over the flagship
- **QUICKSTART flightrecorder mention** — the quickstart never mentions the flightrecorder option
- **cmdguard-lint ServiceName-pollution rule** — a rule (or doc) forbidding `samber-do-auditlog` ServiceName pollution
- **Suppression expiry (`ExpiresAt`)** — cmdguard-lint suppressions with an expiry date
- **Promote the HTML-report structural validator** into `scripts/validate-html-report.py`
- **Consumer-skill drift guard test** — a test failing when `.agents/skills/cmdguard/SKILL.md` source-line references rot (AGENTS.md requires same-change updates)
- **cmdguard-lint suppression FP hardening** — directive-in-string-literal false positives for AST-comment suppressions
- **cmdguard-lint `command_name` field** — dedicated field over ServiceRef overload (upstream samber-do-auditlog schema 0.4.0)
- **lint/go.mod `go-humanize` require audit** — direct vs indirect classification hygiene
- **Document taskctl's manual option configuration** — why the flagship wires options by hand instead of helpers
- **Headless browser in the dev shell** — availability check for HTML-report rendering tests

### Flight Recorder Enhancements

- **`MaxSnapshots` config** — rate limiting / disk protection against runaway trace file creation
- **`CaptureReasonPanic`** — capture on panic recovery (currently only slow/error)
- **`CaptureReasonTimeout`** — capture on context-deadline
- **`Sync()` method** — flush pending captures without stopping the recorder
- **`Recorder.Status()`** — snapshot stats (started, captures, last capture time)
- **Configurable timestamp format** — let users choose precision or timezone
- **`WithFlightRecorderIf[T](cond)`** — custom capture predicates
- **Structured logging** — `slog` handler instead of printf-style
- **Metric hooks** — capture count, bytes written, capture duration
- **`CaptureOnSignal`** — capture on SIGINT/SIGTERM before shutdown
- **gzip compression** — compress snapshots on write for disk savings
- **Trace upload hook** — post-capture callback for remote storage
- **`flightrecorder/d2`** — trace visualization export
- **`flightrecorder/pprof`** — convert trace to pprof profile
- **env-based config** — `WithFlightRecorderEnvVar` for environment-driven setup
- **`go tool trace` parse validation** — test that produced snapshots actually parse
- **CI fuzzing** — `-fuzztime` in CI for continuous fuzzing
- **Separate file-creation vs write errors** — `Capture` error taxonomy
- **`os.CreateTemp` pattern** — collision-free snapshot filenames
- **`Config.Validate()`** — validate config before `Start`
- **DoctorCommand integration** — recorder health as a doctor check
- **`ExportAuditLog`-style trace export config** — trace export options
- **`WithOnError` capture trigger** — capture on the error-callback path
- **`CaptureAsync` variant** — `CaptureAsync(ctx, name, reason) <-chan error`
- **`gotraceui` integration** — alternative trace viewer
- **`MinFreeDisk` config** — skip capture below a disk-space threshold
- **Exported `SanitizeFilename`** — useful for other modules
- **`Start()` returning a cleanup `func()`** — defer ergonomics
- **`WithFlightRecorderAutoStop`** — stop the recorder after the first capture
