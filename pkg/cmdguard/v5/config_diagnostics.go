package v5

import (
	"errors"
	"fmt"
)

// ConfigDiagnostic names one configuration health check and its result. A nil
// Err means the check passed.
type ConfigDiagnostic struct {
	Name string
	Err  error
}

// Diagnostic names, stable because they appear in `app doctor` output.
const (
	// DiagnosticConfigPath reports whether the config file search found a file.
	DiagnosticConfigPath = "config-path"
	// DiagnosticConfigParse reports whether the found file parses into the
	// config struct.
	DiagnosticConfigParse = "config-parse"
	// DiagnosticConfigKeys reports config-file keys no field accepts, even when
	// tolerant mode would have ignored them at startup.
	DiagnosticConfigKeys = "config-keys"
	// DiagnosticConfigSchema reports a consumer-supplied schema check.
	DiagnosticConfigSchema = "config-schema"
)

// ConfigDiagnostics checks the CLI's config file. It returns nothing when the
// CLI has no config file configured, so doctor output stays free of noise for
// applications that do not use one.
//
// The checks re-read the file rather than trusting the startup load: startup
// succeeds by ignoring unknown keys in tolerant mode, which is exactly the
// failure this diagnoses. A parse failure at startup means NewCLI already
// returned an error, so in practice the parser check confirms what is already
// known; it is kept because doctor may also run on a CLI whose file changed on
// disk since construction.
func (cli *CLI[T]) ConfigDiagnostics() []ConfigDiagnostic {
	if cli.spec.configLoader == nil || len(cli.spec.configFilePaths) == 0 {
		return nil
	}

	paths := cli.configSearchPaths()

	diagnostics := make([]ConfigDiagnostic, 0, 2)
	diagnostics = append(diagnostics, ConfigDiagnostic{Name: DiagnosticConfigPath, Err: configPathDiagnostic(paths)})

	if firstExistingConfigPath(paths) == "" {
		return diagnostics
	}

	diagnostics = append(diagnostics,
		ConfigDiagnostic{Name: DiagnosticConfigParse, Err: cli.freshConfigLoadError(paths, false)},
		ConfigDiagnostic{Name: DiagnosticConfigKeys, Err: configKeyDiagnostic(cli.freshConfigLoadError(paths, true))},
	)

	return diagnostics
}

// configPathDiagnostic reports when no config file exists in the search order.
func configPathDiagnostic(paths []string) error {
	if firstExistingConfigPath(paths) != "" {
		return nil
	}

	return fmt.Errorf("%w: none of %v found", ErrConfigFileNotFound, paths)
}

// configKeyDiagnostic keeps only unknown-key failures: a plain parse error is
// already reported by the parse check and must not be counted twice.
func configKeyDiagnostic(err error) error {
	if err == nil || !errors.Is(err, ErrConfigUnknownKey) {
		return nil
	}

	return err
}

// freshConfigLoadError loads the config file into a zero config value. A fresh
// value guarantees the result reflects the file, not the application's current
// state.
func (cli *CLI[T]) freshConfigLoadError(paths []string, strict bool) error {
	var fresh T

	if _, isBuiltinLoader := cli.spec.configLoader.(*KoanfLoader); isBuiltinLoader {
		freshLoader := NewKoanfLoader(paths...)
		freshLoader.SetStrictKeys(strict)

		_, err := freshLoader.Load(nil, &fresh)

		return err
	}

	// A custom loader cannot be asked to reject unknown keys, and running it
	// twice would double any side effects. Report nothing rather than guess.
	if strict {
		return nil
	}

	_, err := loadConfigFile(paths, cli.spec.configLoader, &fresh)

	return err
}
