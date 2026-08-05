const importPath = "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4";

export const heroCode = `func main() {
    cli, _ := v4.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{},
        v4.WithGracefulShutdown(), // SIGINT → reverse-order shutdown
    )

    // Register a service (lazy, lifecycle-managed)
    v4.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
        return &Database{DSN: "postgres://..."}, nil
    })

    cmd, _ := v4.NewCommand("query", v4.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, _ v4.NoFlags) error {
            db, _ := v4.Invoke[*Database](cli.Scope())
            return db.Query(ctx)
        },
        v4.WithShort("Query the database"),
    )
    v4.AddCommand(cli, cmd)

    cli.ExecuteAndExit(context.Background()) // zero panics, correct exit code
}`;
