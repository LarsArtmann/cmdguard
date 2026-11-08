export const heroCode = `func main() {
    cli, _ := v5.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{}) // batteries ON by default

    // Register a service (lazy, lifecycle-managed)
    v5.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
        return &Database{DSN: "postgres://..."}, nil
    })

    cmd, _ := v5.NewCommand("query", v5.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, _ v5.NoFlags) error {
            db, _ := v5.Invoke[*Database](cli.Scope())
            return db.Query(ctx)
        },
        v5.WithShort("Query the database"),
    )
    v5.AddCommand(cli, cmd)

    cli.ExecuteAndExit(context.Background()) // zero panics, correct exit code
}`;
