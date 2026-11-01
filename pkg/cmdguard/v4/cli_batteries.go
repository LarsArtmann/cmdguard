package v4

import "fmt"

// NewApp constructs a CLI and lets build register its commands before
// returning it. It is the one-call entry point for the common shape:
//
//	cli, err := cmdguard.NewApp("app", "My app", AppConfig{},
//		func(cli *cmdguard.CLI[AppConfig]) error {
//			return cmdguard.AddCommands(cli, listCmd, addCmd)
//		},
//		cmdguard.WithBatteries(),
//	)
//	if err != nil {
//		return err
//	}
//	cli.ExecuteAndExit(ctx)
//
// A nil build is valid (root-only CLI). A build error is returned wrapped and
// the partially built CLI is discarded, so callers never execute a half-wired
// application.
func NewApp[T any](
	name, short string,
	defaults T,
	build func(*CLI[T]) error,
	opts ...CLIOption,
) (*CLI[T], error) {
	cli, err := NewCLI(name, short, defaults, opts...)
	if err != nil {
		return nil, err
	}

	if build != nil {
		err = build(cli)
		if err != nil {
			return nil, fmt.Errorf("building CLI %q: %w", name, err)
		}
	}

	return cli, nil
}

// ensureBatteryCommands registers the commands implied by WithBatteries:
// version (when a version is set) and doctor. Registration is idempotent and
// yields to consumer-defined commands of the same name.
func (cli *CLI[T]) ensureBatteryCommands() error {
	if !cli.spec.batteries {
		return nil
	}

	if cli.spec.version != "" && !cli.hasCommand("version") {
		cmd, err := VersionCommand(cli)
		if err != nil {
			return fmt.Errorf("building version command: %w", err)
		}

		err = AddCommand(cli, cmd)
		if err != nil {
			return fmt.Errorf("registering version command: %w", err)
		}
	}

	if !cli.hasCommand("doctor") {
		cmd, err := DoctorCommand(cli)
		if err != nil {
			return fmt.Errorf("building doctor command: %w", err)
		}

		err = AddCommand(cli, cmd)
		if err != nil {
			return fmt.Errorf("registering doctor command: %w", err)
		}
	}

	return nil
}

// hasCommand reports whether a direct subcommand with the given name exists.
func (cli *CLI[T]) hasCommand(name string) bool {
	for _, cmd := range cli.rootCmd.Commands() {
		if cmd.Name() == name {
			return true
		}
	}

	return false
}
