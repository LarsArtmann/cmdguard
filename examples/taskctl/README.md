# taskctl — A Production-Grade Task Manager CLI

A complete, runnable example demonstrating **every major cmdguard feature** in a realistic application.

## Quick Start

```bash
# From project root:
go run examples/taskctl/main.go --help
go run examples/taskctl/main.go list
go run examples/taskctl/main.go list --format json --all
go run examples/taskctl/main.go add --title "Fix bug" --priority high
go run examples/taskctl/main.go done --id 1
go run examples/taskctl/main.go stats
go run examples/taskctl/main.go version
go run examples/taskctl/main.go doctor
go run examples/taskctl/main.go config show
go run examples/taskctl/main.go db migrate --env production
```

## Feature Matrix

Every cmdguard feature is demonstrated in this example:

| Feature                      | Where               | How                                                                                                                         |
| ---------------------------- | ------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **CLI Construction**         | main.go             | `NewCLI` with version, env prefix, config file, graceful shutdown, strict validation, config validation, middleware, groups |
| `WithCLIVersion`             | main.go             | `"1.0.0"`                                                                                                                   |
| `WithEnvPrefix`              | main.go             | `TASKCTL_`                                                                                                                  |
| `WithConfigFile`             | main.go             | JSON from `~/.config/taskctl/config.json`                                                                                   |
| `WithConfigValidation`       | main.go             | Validates DataDir not empty                                                                                                 |
| `WithGracefulShutdown`       | main.go             | Graceful DI shutdown on SIGINT/SIGTERM                                                                                      |
| `WithStrictValidation`       | main.go             | All commands have `WithShort`                                                                                               |
| `WithMiddleware`             | main.go             | `TimingMiddleware` + `RecoveryMiddleware`                                                                                   |
| `WithGroup`                  | main.go             | "tasks" + "system" groups                                                                                                   |
| `AddGlobalBoolFlag`          | main.go             | `--debug` global flag                                                                                                       |
| **env tags**                 | types.go            | `env:"TASK_LOG_LEVEL"`, `env:"TASK_DATA_DIR"`                                                                               |
| **count:"true"**             | types.go            | `-v`/`-vv`/`-vvv` verbosity                                                                                                 |
| **required:"true"**          | types.go            | `add --title`, `done --id`                                                                                                  |
| **prompt:"..."**             | types.go            | `prompt:"Task title?"` on add                                                                                               |
| **Duration type**            | types.go            | Config timeout                                                                                                              |
| **LogLevel type**            | types.go            | Config log-level                                                                                                            |
| **Port type**                | types.go            | Config port                                                                                                                 |
| `NewCommand`                 | commands.go         | list, add, done, stats, inspect                                                                                             |
| `DoctorCommand`              | commands.go         | Doctor command with DI health + custom checks                                                                               |
| `NewParentCommand`           | commands.go         | db, config command groups                                                                                                   |
| `WithFlags`                  | commands.go         | All task commands                                                                                                           |
| `WithShort`                  | commands.go         | All commands                                                                                                                |
| `WithLong`                   | main.go doc comment | Package docs                                                                                                                |
| `WithExample`                | commands.go         | list, add, done                                                                                                             |
| `WithAliases`                | commands.go         | `"ls"` alias for list                                                                                                       |
| `WithPreRunE`                | commands.go         | add validates title length                                                                                                  |
| `WithPostRunE`               | commands.go         | done triggers cleanup                                                                                                       |
| `WithHidden`                 | commands.go         | secret command                                                                                                              |
| `WithDeprecated`             | commands.go         | complete → done                                                                                                             |
| `WithExactArgs`              | commands.go         | inspect requires 1 arg                                                                                                      |
| `WithNoArgs`                 | commands.go         | list, add, done, stats                                                                                                      |
| `WithValidArgs`              | commands.go         | inspect static args                                                                                                         |
| `WithGroupID`                | commands.go         | tasks vs system grouping                                                                                                    |
| `WithPromptOnMissing`        | commands.go         | add prompts for missing title                                                                                               |
| `WithCompletion`             | commands.go         | done --id dynamic completion                                                                                                |
| `WithAuditLog`               | main.go             | DI audit logging via samber-do-auditlog; `AUDIT_LOG_FORMAT` picks export (html/json/ndjson/csv/tsv/mermaid/dot)             |
| `OutputTable`                | commands.go         | list, stats with --format                                                                                                   |
| `output.ParseFormat`         | commands.go         | list --format json/csv/yaml                                                                                                 |
| **DI: Provide/Invoke**       | types.go + main.go  | TaskStore with lifecycle                                                                                                    |
| **Shutdowner**               | types.go            | TaskStore.Shutdown                                                                                                          |
| **HealthcheckerWithContext** | types.go            | TaskStore.HealthCheck                                                                                                       |
| `NewScopeFromInjector`       | types.go            | In NewTaskStore                                                                                                             |
| `NewCommandError`            | commands.go         | Error wrapping                                                                                                              |
| `NewFlagError`               | commands.go         | Invalid format/priority                                                                                                     |
| `NewServiceError`            | commands.go         | Store resolution                                                                                                            |
| `NewExitError`               | commands.go         | done returns code 2 on not-found                                                                                            |
| `VersionCommand`             | commands.go         | Auto version subcommand                                                                                                     |
| **BranchingFlowContext**     | commands.go         | inspect reads flow path                                                                                                     |
| `EditInEditor`               | commands.go         | config edit opens $EDITOR                                                                                                   |
| `ValueOrDefault`             | main_test.go        | Helper usage in tests                                                                                                       |
| `EnsureValid`                | main_test.go        | Helper usage in tests                                                                                                       |
| `ParseDuration`              | main_test.go        | Duration parsing test                                                                                                       |
| **NoFlags type**             | commands.go         | doctor, version, hidden, deprecated                                                                                         |

## Architecture

```
examples/taskctl/
├── main.go        # CLI construction, DI setup, all CLI options
├── commands.go    # All command definitions with options
├── types.go       # Config, flags, domain types, TaskStore service
├── main_test.go   # Comprehensive integration tests (~66 tests)
└── README.md      # This file
```

## Running Tests

```bash
go test ./examples/taskctl/... -count=1 -race
```
