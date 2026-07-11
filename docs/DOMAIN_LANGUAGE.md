# Domain Language

A **Unified Language** for `cmdguard` — shared across contributors and AI assistants.

Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

## Glossary

| Term          | Definition                                                          | Context                              |
| ------------- | ------------------------------------------------------------------- | ------------------------------------ |
| CLI           | A `CLI[T]` instance — the root application container                | Entry point for all cmdguard apps    |
| Command       | A `Command[T,F]` — a typed subcommand with config T and flags F     | Building block for command trees     |
| Config        | The root type parameter `T` on `CLI[T]` — typed application config  | Drives flag registration             |
| Flags         | The per-command type parameter `F` on `Command[T,F]`                | Typed flag struct with struct tags   |
| FlagRegistry  | Holds parsed flag tags and type/validator registries                | Copy-on-write, per-CLI instance      |
| FlagTag       | Parsed metadata from struct tags (`flag`, `default`, `help`, etc.)  | Drives flag registration and parsing |
| Scope         | DI container wrapping `samber/do/v2` injector                       | Manages service lifecycle            |
| Plugin        | Bundles custom type handlers + validators for one-step registration | Applied globally or per-instance     |
| TypeHandler   | Interface for registering, parsing, and defaulting a flag type      | Extensible via `RegisterTypeHandler` |
| FlagValidator | Named validation function applied to flag values                    | Registered via `RegisterValidator`   |
| FlowContext   | Tracks command execution path and branch state                      | Used by middleware and hooks         |
| SilenceUsage  | Suppresses cobra's usage-on-error footgun                           | True by default                      |

## Entities

Objects with identity and lifecycle.

| Term         | Definition                                                | Context                               |
| ------------ | --------------------------------------------------------- | ------------------------------------- |
| CLI[T]       | Root application — owns scope, root command, registry     | Created once per process              |
| Command[T,F] | Typed subcommand with RunE handler                        | Registered via `AddCommand`           |
| Scope        | DI container — manages service registration and lifecycle | Created by CLI or provided externally |

## Value Objects

Immutable objects defined by attributes.

| Term         | Definition                                         | Context                                 |
| ------------ | -------------------------------------------------- | --------------------------------------- |
| FlagTag      | Parsed struct tag metadata for a single flag field | Built by `ParseFlagTags`                |
| OutputFormat | Type-safe enum for output formats (table, json, …) | Aliased from go-output                  |
| Duration     | Validated duration string (e.g. "5s", "1h30m")     | Custom type with Parse/Default handlers |
| Port         | Validated network port (1-65535)                   | Custom type with Parse/Default handlers |
| Email        | RFC 5322 validated email address                   | Custom type with Parse/Default handlers |
| URL          | Validated URL string                               | Custom type with Parse/Default handlers |
| Enum         | Enumerated string value from a fixed set           | Custom type with Parse/Default handlers |

## Commands

Actions the system can perform.

| Term            | Definition                                       | Context                             |
| --------------- | ------------------------------------------------ | ----------------------------------- |
| NewCLI          | Creates a CLI application with typed config      | Entry point                         |
| AddCommand      | Registers a typed subcommand on the CLI          | Supports per-command flag types     |
| Execute         | Runs the CLI with os.Args                        | Returns error for exit-code mapping |
| ExecuteWithArgs | Runs the CLI with explicit args (testing)        | Used in tests                       |
| ExecuteAndExit  | Runs CLI and calls os.Exit with mapped exit code | Blessed entry point                 |

---

## Bounded Contexts

cmdguard has distinct areas of responsibility. Each context has its own vocabulary and rules.

### 1. Command Construction

Building typed command trees from Go structs.

| Term          | Definition                                                              |
| ------------- | ----------------------------------------------------------------------- |
| Use           | The command name as typed by the user (e.g. `"deploy"`)                 |
| Short         | One-line description shown in help and command listings                 |
| Long          | Multi-line description shown in detailed help                           |
| Example       | Usage example string shown in help                                      |
| RunE          | The handler function executed when the command runs                     |
| ParentCommand | A command with subcommands but no own RunE                              |
| NoFlags       | Sentinel type for commands that take no flags (`type NoFlags struct{}`) |

### 2. Flag System

Parsing, validating, and defaulting typed flags from struct tags.

| Term          | Definition                                                                                                                                 |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| FlagTag       | Parsed metadata from struct tags (`flag`, `default`, `help`, `short`, `env`, `validate`, `count`, `required`, `prompt`, `local`, `hidden`) |
| FlagRegistry  | Per-CLI collection of flag tags + type/validator registries (copy-on-write)                                                                |
| TypeHandler   | Interface for registering, parsing, and defaulting a specific flag type                                                                    |
| FlagValidator | Named validation function applied to flag values at parse time                                                                             |
| DefaultValue  | The initial value for a flag before env/flag overrides                                                                                     |
| Precedence    | Resolution order: explicit flag > env var > config file > tag default                                                                      |
| Counting Flag | A flag using `count:"true"` that increments on each occurrence (e.g. `-vvv`)                                                               |
| Scoped Flag   | A flag with `local:"true"` that is NOT inherited by subcommands                                                                            |
| Hidden Flag   | A flag with `hidden:"true"` that is excluded from help but fully functional                                                                |

### 3. Dependency Injection

Service lifecycle and injection via samber/do/v2.

| Term             | Definition                                                        |
| ---------------- | ----------------------------------------------------------------- |
| Scope            | DI container wrapping `do.Injector`                               |
| Injector         | The underlying samber/do/v2 service container                     |
| Override         | Replace a service for testing (clone scope, override, invoke)     |
| CloneScope       | Copy registrations without invoked state                          |
| Shutdown         | DI service shutdown in reverse invocation order                   |
| GracefulShutdown | SIGINT/SIGTERM-triggered DI shutdown via `WithGracefulShutdown()` |

### 4. Output and Formatting

Rendering command output in multiple formats.

| Term         | Definition                                                            |
| ------------ | --------------------------------------------------------------------- |
| OutputFormat | Type-safe enum for output formats (table, json, yaml, markdown, etc.) |
| OutputTable  | Render tabular data in the selected output format                     |
| OutputResult | Render structured data in the selected output format                  |
| AuditLog     | DI service invocation audit trail (11 export formats)                 |

### 5. Extension Points

Interfaces and hooks for customization without modifying core.

| Term          | Definition                                                                   |
| ------------- | ---------------------------------------------------------------------------- |
| Plugin        | Bundles custom type handlers + validators for one-step registration          |
| PromptRunner  | Interface for interactive flag prompting (huh/v2 impl in prompts sub-module) |
| HelpTransform | Function hook for transforming command help text before display              |
| Middleware    | Generic chain type for wrapping command execution                            |

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
