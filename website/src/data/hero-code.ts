const importPath = "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3";

export const heroCode = `func main() {
    cli, _ := v3.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{},
        v3.WithGracefulShutdown(), // SIGINT → reverse-order shutdown
    )

    // Register a service (lazy, lifecycle-managed)
    v3.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
        return &Database{DSN: "postgres://..."}, nil
    })

    cmd, _ := v3.NewCommand("query", v3.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, _ v3.NoFlags) error {
            db, _ := v3.Invoke[*Database](cli.Scope())
            return db.Query(ctx)
        },
        v3.WithShort("Query the database"),
    )
    v3.AddCommand(cli, cmd)

    cli.ExecuteAndExit(context.Background()) // zero panics, correct exit code
}`;
