# ADR-001: fang Integration Strategy

**Status:** ACCEPTED
**Date:** 2026-06-10
**Scope:** `pkg/cmdguard/v3/cli.go`, `cli_options.go`, `cli_accessors.go`

## Context

cmdguard uses `charm.land/fang/v2` for styled Cobra output (help, errors, version). fang exposes 6 functional options plus the `Execute` entry point. cmdguard must decide which fang features to integrate, which to skip, and which to replace with its own implementations.

## fang API Surface

| fang Option                       | Purpose                               |
| --------------------------------- | ------------------------------------- |
| `fang.Execute(ctx, cmd, opts...)` | Styled help, errors, auto-completions |
| `fang.WithVersion(v)`             | Version string for `--version`        |
| `fang.WithCommit(c)`              | Git commit hash appended to version   |
| `fang.WithNotifySignal(sig...)`   | Context cancellation on OS signals    |
| `fang.WithoutManpage()`           | Disable auto `man` subcommand         |
| `fang.WithErrorHandler(fn)`       | Custom error display                  |
| `fang.WithColorSchemeFunc(fn)`    | Custom color theme                    |

## Decision

### INTEGRATED (auto-piped from cmdguard options)

**`fang.Execute`** — Always used when `WithFang(true)` (default). Provides styled help output, styled error display, auto-completion, and VT processing on Windows.

**`fang.WithVersion`** — Auto-piped from `WithCLIVersion[T](v)`. A single cmdguard option controls both cmdguard's `VersionCommand` subcommand AND fang's `--version` flag. Users never call `fang.WithVersion` directly.

**`fang.WithCommit`** — Auto-piped from `WithCLICommit[T](c)`. Same principle: one option, both systems.

### INTENTIONALLY SKIPPED (cmdguard provides superior alternatives)

**`fang.WithNotifySignal`** — fang's implementation is a thin wrapper around `signal.NotifyContext(ctx, sig...)` (see `fang.go:169`). cmdguard's `WithSignalHandling[T]()` does the exact same thing — `signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)` in `cli.go:254`. Additionally, cmdguard's `WithGracefulShutdown[T]()` goes further: it cancels the context AND calls `scope.Shutdown(context.WithoutCancel(ctx))` to gracefully shut down all DI services implementing `do.ShutdownerWithError` in reverse invocation order. Using both would create **double signal registration** and **competing context cancellations**.

Reason for own implementation:

1. **DI lifecycle integration** — fang has no concept of DI; cmdguard's signal handling is integrated with samber/do service shutdown
2. **`context.WithoutCancel`** — Graceful shutdown needs a non-cancelled context to complete cleanup; fang doesn't provide this
3. **No double registration** — Both would call `signal.NotifyContext`; having one is correct

**`fang.WithoutManpage`** — fang auto-registers a hidden `man` subcommand using `mango`. cmdguard has its own `ManPage()` / `WriteManPage()` / `GenerateManPageCommand()` API that gives users explicit control over man page generation. The two systems would conflict (two `man` subcommands). If cmdguard wanted to use fang's man page, it would need to disable its own — but cmdguard's approach is richer (custom section numbers, writer control, typed subcommand).

### NOT YET EXPOSED (gaps)

**`fang.WithErrorHandler`** — Currently errors always use fang's default styled output. Power users may want custom error formatting (e.g., structured JSON errors, custom branding). Should be exposed as `WithFangErrorHandler[T](fn)` in a future release.

**`fang.WithColorSchemeFunc`** — Allows custom lipgloss color themes for help/error output. Niche but straightforward. Could be exposed as `WithFangColorScheme[T](fn)`.

## Consequences

**Positive:**

- Single-source-of-truth for version and commit — one option controls both cmdguard and fang
- DI-aware signal handling that fang cannot provide
- No double signal registration or competing context cancellations
- Users don't need to know about fang at all for common operations

**Negative:**

- `WithFangOptions[T](opts...)` still exists for advanced users who want to pass raw fang options — could conflict if they pass `fang.WithVersion` when `WithCLIVersion` was already set (double version in fang opts)
- Future fang features may require similar integration decisions

## Implementation Details

```
WithCLIVersion[T]("1.0.0")
  → cli.version = "1.0.0"           (cmdguard VersionCommand uses this)
  → cli.rootCmd.Version = "1.0.0"   (cobra --version uses this)
  → fangOpts += fang.WithVersion("1.0.0")  (fang --version uses this)

WithCLICommit[T]("abc123")
  → fangOpts += fang.WithCommit("abc123")

WithSignalHandling[T]()
  → signal.NotifyContext(SIGINT, SIGTERM)   (cancels handler context)

WithGracefulShutdown[T]()
  → signal.NotifyContext(SIGINT, SIGTERM)   (cancels handler context)
  → defer scope.Shutdown(WithoutCancel(ctx)) (shuts down DI services)
```
