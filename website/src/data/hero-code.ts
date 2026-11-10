export const heroCode = `func main() {
    cli, err := v5.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{}) // batteries ON by default
    if err != nil {
        log.Fatal(err)
    }

    // Register a service (lazy, lifecycle-managed)
    err = v5.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
        return &Database{DSN: "postgres://..."}, nil
    })
    if err != nil {
        log.Fatal(err)
    }

    query, err := v5.NewCommand("query", v5.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, _ v5.NoFlags) error {
            db, invokeErr := v5.Invoke[*Database](cli.Scope())
            if invokeErr != nil {
                return invokeErr
            }

            return db.Query(ctx)
        },
        v5.WithShort("Query the database"),
    )

    b := v5.NewBuilder(cli)
    v5.Try(b, query, err) // collect constructor errors, no ladder
    if buildErr := b.Build(); buildErr != nil {
        log.Fatal(buildErr)
    }

    cli.ExecuteAndExit(context.Background()) // zero panics, correct exit code
}`;
