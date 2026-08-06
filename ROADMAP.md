# Roadmap

> Long-term direction and raw ideas not yet refined into actionable tasks.
> Items here are NOT commitments — they're directions worth exploring.
> See [TODO_LIST.md](TODO_LIST.md) for actionable, bounded work.

---

## v5 Major Release

The next major version will focus on API clarity and correctness. The following
breaking changes are tracked with `TODO(v5)` markers in source code:

| #   | Change                                             | Rationale                                                      | Source Location         |
| --- | -------------------------------------------------- | -------------------------------------------------------------- | ----------------------- |
| 1   | Rename `CommandInfo` → `CommandMetadata`           | "Info" suffix is vague; `Metadata` is precise                  | `middleware.go:40`      |
| 2   | Rename `TypeHandler` → `TypeCodec`                 | "Handler" is generic; `Codec` captures parse/default dual role | `type_handler.go:13`    |
| 3   | Rename `PromptRunner` → `HuhPrompter` (or similar) | "Runner" suffix is generic                                     | `prompts/prompts.go:27` |

Additional v5 candidates from partially functional items:

| #   | Change                                                    | Rationale                                                              |
| --- | --------------------------------------------------------- | ---------------------------------------------------------------------- |
| 4   | Remove `SetConfig(cfg)`                                   | Unsafe post-construction mutation without re-initializing FlagRegistry |
| 5   | Rename `Get[T]` → `GetService[T]`                         | `Get` is too generic for a DI scope                                    |
| 6   | Fix `RegisterInScope(parent, name, ...any)` to be generic | Takes `...any` instead of type-safe generics                           |
| 7   | Redesign `Package[T](scope, ...)`                         | Unusual API shape with pre-existing `*Scope` param                     |

---

## Architectural Directions

### Middleware Context Propagation

Current middleware uses `next func() error` — context is NOT propagated to the
next middleware in the chain. This blocks timeout/cancellation middleware. A
redesign to `next func(ctx context.Context) error` would enable context-aware
middleware but is a breaking change to the `Middleware[T]` interface.

### Command-Level Audit Middleware

Only DI lifecycle events are currently captured by the audit log. A
command-level audit middleware would capture command execution events (start,
end, duration, error). See FEATURES.md "Command-level audit middleware" (PLANNED).

### Internal Package Split

If `pkg/cmdguard/v4` grows beyond 12k LOC, consider splitting into `v4` +
`v4/internal/` to separate public API surface from implementation. Current
size is below the trigger threshold — this is YAGNI until growth demands it.

---

## Deferred from 2026-07-18 Audit Closure

These items were explicitly deferred during the post-audit closure to prevent
the long tail from haunting the TODO list. Each has a rationale for deferral.

| #   | Item                                                                                              | Rationale                                                                                             |
| --- | ------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| 1   | Re-run 4 skills with reference files loaded                                                       | More reports without action = more debt. Re-audit in a dedicated session after closure ships.         |
| 2   | Run 4 additional skills (brutal-self-review, library-deep-dive, status-report, docs-health BUILD) | New audits, not closure of existing findings.                                                         |
| 3   | `TypeHandler` / `TypeHandlerFunc` rename + `// Deprecated` alias                                  | Public API break — deferred to v5. See TODO(v5) markers above.                                        |
| 4   | `ConfigFile` branded type                                                                         | YAGNI until a consumer needs it.                                                                      |
| 5   | Extract koanf into a sub-module                                                                   | YAGNI — no consumer has asked; 4 koanf deps are already isolated. LOC trigger not met.                |
| 6   | Split `v4` into `v4` + `v4/internal/`                                                             | LOC trigger (12k) not met. Premature split adds boundary friction.                                    |
| 7   | Fuzz corpus expansion                                                                             | Existing 8 targets have minimal corpus; valuable but not blocking.                                    |
| 8   | Audit `examples/taskctl/main_test.go` (876 lines)                                                 | Test-smell audit is a separate concern.                                                               |
| 9   | `CONTRIBUTING.md` refresh                                                                         | Not blocking; verify-then-decide in a future pass.                                                    |
| 10  | Verify `git-town.toml` + `library-policy.yaml`                                                    | Config sanity, not closure.                                                                           |
| 11  | ~~Update `WHAT_THIS_PROJECT_IS_ABOUT.md` + `_NOT.md`~~                                            | ✅ Done 2026-08-05 — updated to v4 API.                                                               |
| 12  | Schedule re-run after v3.1 ships                                                                  | v3.1 has shipped; v4.0.0 has shipped. Next re-audit window: after v5 or when feature growth warrants. |

---

## Raw Ideas (Not Refined)

These are unformed ideas that may or may not become features:

### Core Library

- **Config file watching** — hot-reload config on file change (inotify/fsnotify)
- **Plugin marketplace** — community-contributed type handlers and validators
- **gRPC middleware sub-module** — command-level gRPC tracing
- **Web-based CLI preview** — render command tree as interactive HTML
- **Shell completion v2** — richer dynamic completion with type-aware suggestions
- **Benchmark dashboard** — track performance across releases
- **v3→v4 migration guide** — dedicated `docs/MIGRATION_v3_v4.md` (v2→v3 got one; v3→v4 only has a CHANGELOG entry)

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
