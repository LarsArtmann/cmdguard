const importPath = "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3";

export const heroCode = `package main

import (
    "context"
    "fmt"
    "os"

    "github.com/samber/do/v2"
    "${importPath}"
)

type AppConfig struct {
    Verbose bool   \`flag:"verbose" short:"v" default:"false"\`
    Output  string \`flag:"output" short:"o" default:"text"\`
}

type Database struct{ DSN string }

func (db *Database) Shutdown() error { return db.Close() }

func main() {
    cli, err := v3.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{},
        v3.WithGracefulShutdown(), // SIGINT → reverse-order shutdown
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to create CLI: %v\\n", err)
        os.Exit(1)
    }

    // Register a service (lazy, lifecycle-managed)
    v3.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
        return &Database{DSN: "postgres://..."}, nil
    })

    cmd, err := v3.NewCommand("query", v3.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, _ v3.NoFlags) error {
            db, _ := v3.Invoke[*Database](cli.Scope())
            return db.Query(ctx)
        },
        v3.WithShort("Query the database"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to create command: %v\\n", err)
        os.Exit(1)
    }
    v3.AddCommand(cli, cmd)

    cli.ExecuteAndExit(context.Background()) // zero panics, correct exit code
}`;
