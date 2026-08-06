# Status Report: Flight Recorder Sub-Module Implementation

> **ANNOTATION (2026-08-06):** This report's entire backlog was executed across the
> next two sessions (2026-08-01 20:45 and 21:22). All P0 items (CHANGELOG, FEATURES,
> README, API docs, Package Guidelines table, lint strategy, flightrecorder/README.md)
> shipped at `ba818e3`. All P1 code items (CaptureToWriter, WithFlightRecorderRecorder,
> Capture refactor, integration tests, Start failure test, edge case tests) shipped at
> `ba818e3`. P2/P3 enhancement ideas were harvested into `ROADMAP.md` "Flight Recorder
> Enhancements." Two godoc examples were lost in a subsequent git corruption event
> (see `2026-08-01_22-03` report) and are tracked in `TODO_LIST.md` (D9).

**Date:** 2026-08-01 19:42
**Session Goal:** Leverage Go 1.25 `runtime/trace.FlightRecorder` (from [the Go blog post](https://go.dev/blog/flight-recorder)) in cmdguard
**Status:** Implementation works, tests pass, lint clean — but **documentation, ecosystem integration, and polish are incomplete**

---

## Executive Summary

Built a new `flightrecorder/` sub-module that wraps Go 1.25's `runtime/trace.FlightRecorder`
as cmdguard middleware. It continuously buffers execution traces in memory and auto-captures
`.trace` snapshots when commands are slow or error. The code is solid (91% coverage, 0 races,
0 lint issues), but the surrounding ecosystem work — CHANGELOG, FEATURES.md, README,
benchmarks, examples, go.sum auditing — was **not done**.

---

## a) FULLY DONE (Working & Verified)

| Item                                | Details                                                                                                                                                                                                                                                                                                                                                        |
| ----------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`flightrecorder/recorder.go`**    | `Config` (8 fields with defaults), `Recorder` (Start/Stop/Enabled/WriteTo/Capture), sentinel errors (`ErrAlreadyStarted`, `ErrNotEnabled`), `sanitizeFilename`, `logf`. Thread-safe with `sync.Mutex` + `sync.WaitGroup`.                                                                                                                                      |
| **`flightrecorder/middleware.go`**  | `Middleware[T]` (lazy start, run-phase-only capture, async snapshot), `WithFlightRecorder[T]` (convenience CLIOption), `evaluateCapture` (slow/error trigger logic).                                                                                                                                                                                           |
| **Race condition fix**              | Identified and fixed data race between async `Capture` goroutine and `Stop()` — added `inflight sync.WaitGroup` so `Stop()` waits for in-flight `WriteTo`/`Capture` operations before calling `fr.Stop()`.                                                                                                                                                     |
| **31 tests passing**                | `recorder_test.go` (15 tests): lifecycle, defaults, capture, filenames, cancelled context, custom log func, sentinel errors, sanitize table-driven. `middleware_test.go` (10 tests): non-nil, next passthrough, non-run-phase skip, slow capture, error capture, no-capture-when-clean, full path, lazy start, convenience options. Plus `waitForFile` helper. |
| **91% coverage**                    | Verified with `-cover` flag.                                                                                                                                                                                                                                                                                                                                   |
| **0 lint issues**                   | Passed `golangci-lint run ./...` with the project's strict `.golangci.yml` (100+ linters enabled).                                                                                                                                                                                                                                                             |
| **0 race conditions**               | Passed `-race` flag.                                                                                                                                                                                                                                                                                                                                           |
| **Workspace wired**                 | `go.work` updated, root `go.mod` has `replace` directive.                                                                                                                                                                                                                                                                                                      |
| **AGENTS.md updated**               | Project structure tree, sub-module table, design principles #15 + #20, sub-modules gotchas section, middleware section — all updated with `flightrecorder` references.                                                                                                                                                                                         |
| **`.golangci.yml` exclusion added** | `flightrecorder/.*_test\.go$` excluded from `paralleltest` (process-wide singleton prevents parallel test execution).                                                                                                                                                                                                                                          |
| **`go mod tidy` run**               | Module has valid `go.mod` and `go.sum`.                                                                                                                                                                                                                                                                                                                        |
| **Zero external dependencies**      | Uses only Go stdlib `runtime/trace` — the leanest sub-module in the project.                                                                                                                                                                                                                                                                                   |

---

## b) PARTIALLY DONE (Incomplete)

| Item                        | What Exists                                          | What's Missing                                                                                                                                                                                                                                                                                                        |
| --------------------------- | ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **AGENTS.md documentation** | Sub-module table, gotchas, design principles updated | Package Guidelines table (line ~139) NOT updated — doesn't list `flightrecorder` as a package. Lint strategy section not updated with the new `paralleltest` exclusion. Exclusion count comment not bumped.                                                                                                           |
| **go.sum hygiene**          | `go mod tidy` ran successfully                       | go.sum has 155 lines — suspiciously large for a module whose only dependency is `cmdguard/v4` (which itself is a local replace). The transitive dependency closure from `replace => ../` pulls in cobra/fang/lipgloss/etc as indirect deps. **Not investigated** whether this is correct or a `go mod tidy` artifact. |

---

## c) NOT STARTED (Gaps — Expected But Missing)

### High Priority

1. ~~**CHANGELOG.md** — No `[Unreleased]` entry for the new sub-module. This is a user-facing addition.~~ done at `ba818e3`
2. ~~**FEATURES.md** — `flightrecorder` not listed anywhere in the feature inventory. Status should be DONE.~~ done at `ba818e3`
3. ~~**README.md** — Root README doesn't mention flight recorder capability. Users discovering cmdguard won't know it exists.~~ done at `ba818e3`
4. ~~**`docs/API.md`** — No API documentation for `Recorder`, `Middleware`, `WithFlightRecorder`, `Config`, `CaptureReason`.~~ done at `ba818e3`
5. ~~**No `flightrecorder/README.md`** — Every other sub-module's purpose is documented somewhere; this one has no standalone docs.~~ done at `ba818e3` (recreated after git corruption at `ba818e3`)

### Medium Priority

6. **No example in `examples/`** — No demonstration of how to wire flight recorder into a CLI. The other middleware modules (telemetry, spinner) are used in `examples/taskctl/`.
7. ~~**No benchmarks** — Project has 26 benchmarks in core; flightrecorder has 0. Middleware overhead (start, time measurement, capture goroutine) should be benchmarked.~~ done at `ba818e3` (3 benchmarks added)
8. ~~**No fuzz tests** — Project has 7 fuzz targets; flightrecorder has 0. `sanitizeFilename` and `Capture` (filename construction) are fuzz-worthy.~~ done at `ba818e3` (`FuzzSanitizeFilename` added)
9. ~~**No `go doc` example functions** — The package doc has usage examples in comments, but no runnable `Example*` functions (project convention uses `example_test.go`).~~ partially done at `ba818e3` (3 examples; 2 more lost in git corruption — see TODO_LIST D9)
10. ~~**Package Guidelines table** — `flightrecorder` not added as a row in the `### Package Guidelines` table in AGENTS.md.~~ done at `ba818e3`

### Low Priority

11. ~~**TODO_LIST.md** — No task entries for follow-up work (benchmarks, examples, etc.).~~ done (tracked in CHANGELOG [Unreleased])
12. ~~**ROADMAP.md** — No mention of flight recorder as a shipped feature or future direction.~~ done (enhancement ideas in ROADMAP "Flight Recorder Enhancements")
13. ~~**`flake.nix`** — Not checked whether `nix flake check` needs updating for the new module (probably fine since it's just `treefmt`, but not verified).~~ done at `ba818e3` (`nix flake check` passes)
14. ~~**`nix fmt` / `nix flake check`** — Not run. Only `golangci-lint fmt` was used.~~ done at `ba818e3`

---

## d) TOTALLY FUCKED UP (Nothing)

Nothing is broken, reverted, or non-functional. All code compiles, all tests pass, all lint is clean. The implementation is correct.

**However**, one thing to be honest about:

### Design Smell: `evaluateCapture` captures on BOTH slow AND error

If `CaptureOnError=true` AND `CaptureOnSlow=true` AND a command is both slow AND errors, the current logic captures with `CaptureReasonSlow` (slow check runs second and overwrites `reason`). This means the error snapshot is mislabeled. The blog post's example captures once on the first trigger; this implementation does the same (only one goroutine fires), but the reason labeling is ambiguous. **Not a bug** (the trace data is identical regardless of reason), but the labeling could be more precise (e.g., capture on both conditions, or label as "slow+error").

---

## e) WHAT WE SHOULD IMPROVE

### Code Quality

1. **`Capture` method is 40+ lines** — Could be split into `buildSnapshotPath` + `writeSnapshot` helpers for readability.
2. **No `CaptureToWriter` convenience method** — `Capture` always writes to a file. Users who want to write to `os.Stdout` or a custom `io.Writer` must call `WriteTo` directly and handle filename/log themselves.
3. **No `CaptureReasonPanic`** — The middleware doesn't capture on panics (only slow/error). A panic-capturing variant would require wrapping `RecoveryMiddleware` or adding panic detection in the middleware itself.
4. **No `Sync()` method** — `Stop()` waits for in-flight captures, but there's no standalone "wait for all pending captures to finish" method for users who want to flush without stopping.
5. **Timestamp format hardcoded** — `time.Now().Format("20060102-150405")` is not configurable. Users with high-frequency captures might want millisecond precision or UTC.
6. **No max-captures limit** — A runaway CLI that errors on every invocation could fill disk with trace files. No rate limiting or max-file-count.
7. **Context not propagated to `fr.WriteTo`** — The stdlib's `WriteTo` doesn't accept a context, so a cancelled context in `Capture` only checks `ctx.Err()` at the top, not during the actual write.

### API Design

8. **`Recorder.Config()` returns value, not pointer** — Users can't modify config after creation. This is probably fine (immutable is good), but should be a conscious decision.
9. **`WithFlightRecorder` creates an unexported Recorder** — Users who want to access the recorder (e.g., to call `Capture` manually for custom triggers) can't. The convenience option is too convenient.
10. **No `WithFlightRecorderRecorder[T](rec)` variant** — No CLIOption that takes a pre-created `*Recorder` for users who want shared lifecycle control.

### Testing

11. **No test for "double Start returns ErrAlreadyStarted"** — The sentinel error is tested indirectly via `TestRecorder_StartStopLifecycle`, but `errors.Is(err, ErrAlreadyStarted)` is not asserted.
12. **No test for concurrent `Capture` calls** — Multiple goroutines calling `Capture` simultaneously is not tested (though the mutex/WG should handle it).
13. **No test for "WriteTo after Stop returns ErrNotEnabled"** — Only tested on never-started recorder.
14. **No test for filename collision** — Two captures in the same second get the same timestamp; filename collision behavior is not tested.
15. **No test for empty command name** — `Capture(ctx, "", reason)` falls back to "command", but this isn't tested.
16. **`TestRecorder_ErrorsAreSentinels` doesn't test `ErrAlreadyStarted`** — Only tests `ErrNotEnabled`.

### Integration

17. **No integration test** — The middleware is tested in isolation, never wired through a real `CLI[T]` + `Execute`. The existing `examples/taskctl/main_test.go` pattern (66 tests) would be the model.
18. **No test that `go tool trace` can actually parse the output** — We verify bytes are written, but never validate the trace format. (This would require shelling out to `go tool trace`, which may not be available in all environments.)

---

## f) Up to 50 Things We Should Get Done Next

> **ANNOTATION (2026-08-06):** All P0 items shipped at `ba818e3`. Most P1 items
> shipped at `ba818e3`. P2/P3 enhancement ideas were harvested into `ROADMAP.md`.
> Items left unmarked are still open.

### P0 — Must Do (Ecosystem Completion)

1. ~~Add `CHANGELOG.md` `[Unreleased]` entry for `flightrecorder` sub-module~~ done at `ba818e3`
2. ~~Add `flightrecorder` to `FEATURES.md` as DONE~~ done at `ba818e3`
3. ~~Add flight recorder section to root `README.md`~~ done at `ba818e3`
4. ~~Add `flightrecorder` API reference to `docs/API.md`~~ done at `ba818e3`
5. ~~Add `flightrecorder` row to Package Guidelines table in `AGENTS.md`~~ done at `ba818e3`
6. ~~Update AGENTS.md Lint Strategy section: bump exclusion count, document the new `paralleltest` path exclusion~~ done at `ba818e3`
7. ~~Create `flightrecorder/README.md` with quick-start guide~~ done at `ba818e3`
8. ~~Verify `nix fmt` and `nix flake check` pass with the new module~~ done at `ba818e3`
9. ~~Add `flightrecorder` to `TODO_LIST.md` as completed work~~ done (in CHANGELOG [Unreleased])

### P1 — Should Do (Quality & Polish)

10. ~~Add `ExampleMiddleware` and `ExampleWithFlightRecorder` functions (`example_test.go`)~~ partially done at `ba818e3` (3 examples exist; 2 more lost in git corruption — TODO_LIST D9)
11. ~~Add benchmarks: `BenchmarkMiddleware_Overhead`, `BenchmarkCapture`, `BenchmarkNew`~~ done at `ba818e3`
12. ~~Add fuzz test for `sanitizeFilename`~~ done at `ba818e3`
13. Add fuzz test for `Capture` filename construction
14. ~~Add `CaptureToWriter(writer io.Writer, ...) (int64, error)` method~~ done at `ba818e3`
15. ~~Add `WithFlightRecorderRecorder[T](rec *Recorder)` CLIOption variant~~ done at `ba818e3`
16. ~~Add test asserting `errors.Is(err, ErrAlreadyStarted)` on double-Start~~ done at `ba818e3`
17. ~~Add test for concurrent `Capture` calls from multiple goroutines~~ done at `ba818e3`
18. ~~Add test for `WriteTo` after `Stop` (should return `ErrNotEnabled`)~~ done at `ba818e3`
19. ~~Add test for empty command name fallback to "command"~~ done at `ba818e3`
20. ~~Add test for filename collision (same-second captures)~~ done at `ba818e3`
21. ~~Add integration test: wire through real `CLI[T]` + `Execute`~~ done at `ba818e3` (3 integration tests)
22. Add `CaptureReasonPanic` support (capture on panic recovery) — _in ROADMAP_
23. ~~Split `Capture` method into `buildSnapshotPath` + `writeSnapshot`~~ done at `ba818e3`
24. Add `Sync()` method for flushing pending captures without stopping — _in ROADMAP_

### P2 — Nice to Have (Enhancement)

> **ANNOTATION (2026-08-06):** All P2 items are enhancement ideas harvested into
> `ROADMAP.md` "Flight Recorder Enhancements." None shipped. Left unmarked = open.

25. Add `MaxSnapshots` config field (rate limit / disk protection)
26. Add configurable timestamp format in filenames
27. Add `CaptureReasonTimeout` for context-deadline captures
28. Add `WithFlightRecorderIf[T](cond func(info, elapsed, err) bool)` for custom capture predicates
29. Add structured logging option (slog handler instead of printf-style)
30. Add metric hooks (capture count, bytes written, capture duration)
31. Add `Recorder.Status()` method returning snapshot stats (started, captures, last capture time)
32. Add `flightrecorder` usage to `examples/taskctl/` main.go
33. Add `go tool trace` parsing validation test (if `go` binary available)
34. Add godoc cross-references from `pkg/cmdguard/v4/middleware.go` to the flight recorder sub-module
35. Add `CaptureOnSignal` option (capture on SIGINT/SIGTERM before shutdown)

### P3 — Future Consideration

> **ANNOTATION (2026-08-06):** All P3 items are future ideas harvested into
> `ROADMAP.md` "Flight Recorder Enhancements." None shipped. Left unmarked = open.

36. Consider `flightrecorder/d2` sub-module for trace visualization export
37. Consider integration with `cli.DoctorCommand` (health check for recorder status)
38. Consider integration with `ExportAuditLog` pattern (trace export config)
39. Consider `flightrecorder/pprof` bridge (convert trace to pprof profile)
40. Consider `WithFlightRecorderEnvVar("FLIGHT_RECORDER_")` for env-based config
41. Consider adding trace capture trigger to `WithOnError` callback path
42. Consider `Recorder.CaptureAsync(ctx, name, reason) <-chan error` variant
43. Consider configurable trace compression (gzip snapshots on write)
44. Consider trace upload hook (post-capture callback for remote storage)
45. Consider `flightrecorder/gotraceui` integration for alternative trace viewer
46. Consider `MinFreeDisk` config (skip capture if disk space below threshold)
47. Consider `SanitizeFilename` as exported utility (useful for other modules)
48. Consider `Recorder.Start()` returning a cleanup `func()` for `defer` ergonomics
49. Consider `WithFlightRecorderAutoStop` option (stop recorder after first capture)
50. Consider documenting the 10 MB/s trace data rate in Config.Log output

---

## g) Questions (Cannot Determine Myself)

### 1. Should `flightrecorder` be added to the root `go.mod` `require` block?

Currently, the root module's `go.mod` has a `replace` directive but no `require` entry for
`github.com/larsartmann/cmdguard/flightrecorder`. The other sub-modules (`glamour`, `spinner`)
ARE in the `require` block because core imports them. But core does NOT import `flightrecorder`
(it's a one-way dependency: flightrecorder imports core, not vice versa). The `replace`
directive alone is sufficient for the workspace build. **Question: should the root module
formally depend on its sub-modules even when it doesn't import them, for consistency?**

### 2. Should the flight recorder be integrated into the `examples/taskctl/` example?

The `taskctl` example is described as the "flagship example" with a "production task manager CLI."
Adding flight recorder would demonstrate the feature, but it would also add noise to an example
that's already comprehensive (66 tests). The other middleware modules (telemetry, spinner) are
NOT currently used in `taskctl` either. **Question: is there a convention for which sub-modules
get demonstrated in examples vs. which get their own example directory?**

### 3. Is the `go.sum` bloat expected?

`flightrecorder/go.sum` has 155 lines, the same as `spinner/go.sum` and nearly the same as
`telemetry/go.sum` (159 lines). For a module with zero external dependencies (only
`runtime/trace` from stdlib + local `replace => ../` for cmdguard/v4), this seems excessive.
The transitive closure comes from `cmdguard/v4`'s own dependencies (cobra, fang, lipgloss, etc.)
being pulled in as indirect deps via the `replace` directive. **Question: is this the expected
behavior for all sub-modules in this workspace, or is there a way to produce a leaner go.sum
for modules that only need the core types?**

---

## Session Metrics

| Metric                      | Value                                                                        |
| --------------------------- | ---------------------------------------------------------------------------- |
| Files created               | 4 (`recorder.go`, `middleware.go`, `recorder_test.go`, `middleware_test.go`) |
| Files modified              | 3 (`go.work`, `go.mod`, `.golangci.yml`, `AGENTS.md`)                        |
| Lines of code (source)      | ~300                                                                         |
| Lines of code (tests)       | ~500                                                                         |
| Test functions              | 31 (38 including subtests)                                                   |
| Coverage                    | 91.0%                                                                        |
| Lint issues                 | 0                                                                            |
| Race conditions             | 0 (after fix)                                                                |
| Build errors                | 0                                                                            |
| External dependencies added | 0 (stdlib only)                                                              |
| Time to implement           | ~1 hour                                                                      |
| Iterations on lint          | 3 (initial write → 72 issues → format → 3 issues → fix → 0)                  |
| Race condition fixes        | 1 (async capture vs Stop — added WaitGroup)                                  |

---

## Conclusion

The flight recorder sub-module is **functionally complete and correct**, but **ecosystem
integration is incomplete**. The code quality bar is high (91% coverage, sentinel errors,
thread-safe, zero deps), but the project's documentation conventions (CHANGELOG, FEATURES,
README, API docs, examples, benchmarks, fuzz tests) were not followed for this new addition.
The 50-item backlog ranges from must-do ecosystem completion to future enhancement ideas.
