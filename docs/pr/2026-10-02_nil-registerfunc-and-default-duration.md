## Fix the nil-RegisterFunc silent flag vanish + make stdlib time.Duration work by default

### Problem (found dogfooding in nsfw-classifier)

Registering a `TypeHandlerFunc` with a nil `RegisterFunc` (easy to do when only
`ParseFunc` seems needed) passes `RegisterTypeHandler` validation, then
`dispatchRegister` calls `Register`, which silently no-ops. The flag is never
added to the FlagSet and the CLI later fails with "unknown flag: --retention" —
at invocation time, far from the cause. We shipped exactly this bug.

### Fix

- `dispatchRegister` rejects a `TypeHandlerFunc` with nil `RegisterFunc` at
  CLI-build time: `type handler for time.Duration has nil RegisterFunc; it
  would register nothing and the flag "retention" would silently vanish`.
  The low-level `Register` stays nil-tolerant (existing pinned test intact).
- `Parse` with nil `ParseFunc` returns a descriptive error instead of panicking;
  `Default` with nil `DefaultFunc` yields nil instead of panicking.
- stdlib `time.Duration` is now in the default registry (same handler as the
  opt-in `RegisterGoDurationHandler`): native duration pflag, `5s` parsing,
  duration defaults. Explicit overrides still win. A consumer's 25-line
  hand-rolled registration goes away.

### Tests

- `TestDispatchRegister_NilRegisterFuncRejected` — build-time rejection, flag absent
- `TestTypeHandlerFunc_NilParseFuncErrors` / `_NilDefaultFuncErrors` — no panics
- `TestGoDurationHandlerDefaultOn` — fresh registry handles `time.Duration` end-to-end
- full `go test ./...` green; CHANGELOG updated under Unreleased
