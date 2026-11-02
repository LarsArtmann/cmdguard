package v4

import (
	"context"
	"fmt"
	"slices"
)

// DoctorCheck is a named diagnostic check that can be added to the doctor command.
type DoctorCheck struct {
	Name string
	Run  func(ctx context.Context) error
}

// doctorConfig holds the configuration for the doctor command.
type doctorConfig[T any] struct {
	short           string
	long            string
	groupID         string
	checks          []DoctorCheck
	schemaValidator func(configPath string) error
}

// DoctorOption configures the doctor command.
type DoctorOption[T any] func(*doctorConfig[T])

// WithDoctorShort sets a custom short description for the doctor command.
func WithDoctorShort[T any](short string) DoctorOption[T] {
	return func(cfg *doctorConfig[T]) {
		cfg.short = short
	}
}

// WithDoctorLong sets a custom long description for the doctor command.
func WithDoctorLong[T any](long string) DoctorOption[T] {
	return func(cfg *doctorConfig[T]) {
		cfg.long = long
	}
}

// WithDoctorGroupID sets the command group ID for the doctor command.
func WithDoctorGroupID[T any](groupID string) DoctorOption[T] {
	return func(cfg *doctorConfig[T]) {
		cfg.groupID = groupID
	}
}

// WithDoctorSchemaValidator adds a config-schema conformance check. It runs
// only when a config file exists, and receives the path that was found. The
// schema emitter lives in its own module (it needs reflection over the config
// struct), so cmdguard core cannot import it; wire it here instead:
//
//	cmdguard.WithDoctorSchemaValidator[Config](func(path string) error {
//	    return schema.ValidateFile(path, Config{})
//	})
func WithDoctorSchemaValidator[T any](validate func(configPath string) error) DoctorOption[T] {
	return func(cfg *doctorConfig[T]) {
		cfg.schemaValidator = validate
	}
}

// WithDoctorCheck adds a named diagnostic check to the doctor command.
// Checks run after DI health checks and can verify external dependencies
// like database connections, file system state, or API reachability.
func WithDoctorCheck[T any](name string, run func(ctx context.Context) error) DoctorOption[T] {
	return func(cfg *doctorConfig[T]) {
		cfg.checks = append(cfg.checks, DoctorCheck{Name: name, Run: run})
	}
}

// DoctorCommand creates a typed "doctor" subcommand that runs health checks
// on all DI services and any custom diagnostic checks, reporting per-check status.
//
// Usage:
//
//	cli, _ := v4.NewCLI[Config]("myapp", "My app", Config{},
//	    v4.WithCLIVersion[Config]("1.0.0"),
//	)
//	docCmd, err := v4.DoctorCommand[Config](cli)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	v4.AddCommand(cli, docCmd)
//
// With custom checks:
//
//	docCmd, err := v4.DoctorCommand[Config](cli,
//	    v4.WithDoctorCheck[Config]("database", func(ctx context.Context) error {
//	        return db.Ping(ctx)
//	    }),
//	)
func DoctorCommand[T any](cli *CLI[T], opts ...DoctorOption[T]) (Command[T, NoFlags], error) {
	cfg := doctorConfig[T]{
		short: "Check system health",
		long:  "Run health checks on all services and report diagnostic status.",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return NewCommand(
		"doctor",
		NoFlags{},
		doctorRun(cli, cfg),
		WithShort(cfg.short),
		WithLong(cfg.long),
		WithExample(cli.spec.name+" doctor"),
		withDoctorGroupID(cfg.groupID),
	)
}

// doctorRun builds the doctor handler: DI health, config diagnostics, the
// optional schema validator, and custom checks — one report line each.
func doctorRun[T any](cli *CLI[T], cfg doctorConfig[T]) func(context.Context, *T, NoFlags) error {
	return func(ctx context.Context, _ *T, _ NoFlags) error {
		w := cli.rootCmd.OutOrStdout()

		results := cli.HealthCheckResultsWithContext(ctx)

		for _, diagnostic := range cli.ConfigDiagnostics() {
			results[diagnostic.Name] = diagnostic.Err
		}

		if cfg.schemaValidator != nil {
			results[DiagnosticConfigSchema] = runSchemaValidator(cfg.schemaValidator, cli.configFilePath)
		}

		for _, check := range cfg.checks {
			err := check.Run(ctx)
			results[check.Name] = err
		}

		names := make([]string, 0, len(results))

		for name := range results {
			names = append(names, name)
		}

		slices.Sort(names)

		passed := 0
		failed := 0

		for _, name := range names {
			err := results[name]
			if err != nil {
				fmt.Fprintf(w, "✗ %s\t%s\n", name, err)

				failed++
			} else {
				fmt.Fprintf(w, "✓ %s\tok\n", name)

				passed++
			}
		}

		fmt.Fprintf(w, "\n%d passed, %d failed\n", passed, failed)

		if failed > 0 {
			doctorErr := fmt.Errorf(
				"%w: %s: %d check(s) failed",
				ErrDoctorFailed,
				cli.spec.name,
				failed,
			)
			exitErr, _ := NewExitError(1, doctorErr)

			return exitErr
		}

		return nil
	}
}

// withDoctorGroupID conditionally adds WithGroupID if non-empty.
func withDoctorGroupID(groupID string) CommandOption {
	return func(s *commandSpec) {
		if groupID != "" {
			s.group = groupID
		}
	}
}

// runSchemaValidator runs a consumer-supplied schema check. No config file
// means there is nothing to validate, so the check passes rather than failing
// for a reason the user cannot act on.
func runSchemaValidator(validate func(configPath string) error, configPath string) error {
	if configPath == "" {
		return nil
	}

	return validate(configPath)
}
