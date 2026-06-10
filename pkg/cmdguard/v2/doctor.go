package v2

import (
	"context"
	"fmt"
	"sort"
)

// DoctorCheck is a named diagnostic check that can be added to the doctor command.
type DoctorCheck struct {
	Name string
	Run  func(ctx context.Context) error
}

// doctorConfig holds the configuration for the doctor command.
type doctorConfig[T any] struct {
	short   string
	long    string
	groupID string
	checks  []DoctorCheck
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
//	cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
//	    v2.WithCLIVersion[Config]("1.0.0"),
//	)
//	docCmd, err := v2.DoctorCommand[Config](cli)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	v2.AddCommand(cli, docCmd)
//
// With custom checks:
//
//	docCmd, err := v2.DoctorCommand[Config](cli,
//	    v2.WithDoctorCheck[Config]("database", func(ctx context.Context) error {
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

	appName := cli.name
	short := cfg.short
	long := cfg.long
	groupID := cfg.groupID
	customChecks := cfg.checks

	return NewCommand[T, NoFlags](
		"doctor",
		func(ctx context.Context, _ *T, _ NoFlags) error {
			w := cli.rootCmd.OutOrStdout()

			results := cli.HealthCheckResultsWithContext(ctx)

			for _, check := range customChecks {
				err := check.Run(ctx)
				results[check.Name] = err
			}

			names := make([]string, 0, len(results))

			for name := range results {
				names = append(names, name)
			}

			sort.Strings(names)

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
					appName,
					failed,
				)
				exitErr, _ := NewExitError(1, doctorErr)

				return exitErr
			}

			return nil
		},
		WithShort[T, NoFlags](short),
		WithLong[T, NoFlags](long),
		withDoctorGroupID[T](groupID),
	)
}

// withDoctorGroupID conditionally adds WithGroupID if non-empty.
func withDoctorGroupID[T any](groupID string) CommandOption[T, NoFlags] {
	return func(cmd *Command[T, NoFlags]) {
		if groupID != "" {
			WithGroupID[T, NoFlags](groupID)(cmd)
		}
	}
}
