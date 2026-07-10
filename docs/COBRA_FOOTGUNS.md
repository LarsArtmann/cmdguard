# Cobra Footguns — What cmdguard Closes

Cobra is a powerful CLI framework, but it has several well-known traps ("footguns")
that catch developers. cmdguard closes these by default or through its API design.

---

## 1. Usage-on-Error (The #1 Footgun)

**Cobra default:** Prints full usage text on every error, burying the actual error message.

**cmdguard fix:** `SilenceUsage = true` by default for root and all subcommands.
The error message is printed exactly once. Use `WithSilenceUsage()` to make the intent
explicit (it's already the default).

## 2. PostRunE Not Called on RunE Error

**Cobra behavior:** `PostRunE` is NOT called when `RunE` returns an error.
Resources opened in `PreRunE` leak.

**cmdguard fix:** `WithCleanup[T](fn)` hooks fire after EVERY command's `RunE`,
including when `RunE` errors. The original error is never swallowed.

## 3. No Context Propagation

**Cobra default:** No built-in signal handling or context cancellation.

**cmdguard fix:** `WithSignalHandling()` cancels context on SIGINT/SIGTERM.
`WithGracefulShutdown()` additionally triggers DI service shutdown in reverse order.

## 4. Flag Parsing Order Ambiguity

**Cobra/pflag:** Flags can be parsed at different times depending on PersistentPreRun
vs Run, making precedence unclear.

**cmdguard fix:** Explicit precedence chain: explicit flag > env var > config file > default.
`WithPostFlagParse[T]()` runs after flag parsing and config validation but before handlers.

## 5. No Type-Safe Config

**Cobra:** Config is typically `map[string]any` or untyped viper bindings.

**cmdguard fix:** Config is a typed struct `T`. Flags, env vars, and config files all
populate the same typed struct. Compile-time safety, no runtime type assertions.

## 6. Error Display Ownership

**Cobra:** Multiple places can print errors (RunE, PersistentPreRunE, the framework itself),
leading to double-printed errors.

**cmdguard fix:** cmdguard owns error display. The error is printed exactly once (fang when
enabled, cobra when disabled). The error returned by `Execute` is for exit-code mapping only.
Consumers must NOT re-print it.

## 7. No Exit Code Control

**Cobra:** `os.Exit(1)` is hardcoded in many examples. Custom exit codes require manual work.

**cmdguard fix:** `ExitCode(err) int` maps errors to exit codes. `NewExitError(code, err)`
returns custom exit codes. `ExecuteAndExit` is the blessed entry point that handles everything.

## 8. Untyped Dependency Injection

**Cobra:** No DI support. Services are typically global variables or manually passed.

**cmdguard fix:** `samber/do/v2` powered DI with typed scopes, `Provide[T]`, `Invoke[T]`,
`Override[T]` for testing, `CloneScope` for isolation.

## 9. No Validation Layers

**Cobra:** No built-in validation for commands, flags, or config.

**cmdguard fix:** `WithStrictValidation()` (requires short descriptions),
`WithDraconianValidation()` (strict + examples on leaf commands),
`WithConfigValidation[T](fn)` (post-parse config validation).

## 10. Untyped Flag Values

**Cobra/pflag:** Flags use `pflag.FlagSet` methods that return `error` but don't validate
semantics (email format, URL format, port range, etc.).

**cmdguard fix:** Struct tags define flags with built-in validators (`validate:"email"`).
Custom value types (`Email`, `URL`, `Port`, `Duration`, etc.) parse and validate at the
type level. Typo suggestions via Levenshtein distance.
