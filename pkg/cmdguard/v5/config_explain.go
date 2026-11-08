package v5

import (
	"cmp"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	output "github.com/larsartmann/go-output"
	"github.com/spf13/cobra"
)

// RedactedValue replaces a secret's resolved value in every cmdguard output.
// Consumers that print resolved configuration themselves should reuse it so
// their logs match what `config explain` shows.
const RedactedValue = "***"

// configExplainCommandName is the parent command that groups config
// introspection. A consumer-defined command with this name wins.
const configExplainCommandName = "config"

// ValueSource identifies where a flag's effective value came from. The order of
// the constants matches precedence: flag beats environment beats config file
// beats the declared default.
type ValueSource string

const (
	// SourceFlag means an explicit command-line flag supplied the value.
	SourceFlag ValueSource = "flag"
	// SourceEnvironment means an environment variable supplied the value.
	SourceEnvironment ValueSource = "environment"
	// SourceConfigFile means a config file supplied the value.
	SourceConfigFile ValueSource = "config-file"
	// SourceDefault means the struct tag default (or the value passed to NewCLI).
	SourceDefault ValueSource = "default"
	// SourceUnset means nothing supplied a value and the tag declared no default.
	SourceUnset ValueSource = "unset"
)

// ResolvedFlag explains how one flag's effective value was decided.
//
// Value carries the resolved value and is redacted when Secret is true; a
// consumer that needs to know whether the secret was set must look at Source
// and Redacted, never at Value. Default is the value the struct tag declared,
// before any config file replaced it, so a reader can see both what the
// application ships with and what this invocation actually uses.
type ResolvedFlag struct {
	Name     string      `json:"name"`
	Key      string      `json:"key"`
	Source   ValueSource `json:"source"`
	Value    string      `json:"value"`
	Default  string      `json:"default,omitzero"`
	Env      string      `json:"env,omitzero"`
	Secret   bool        `json:"secret,omitzero"`
	Redacted bool        `json:"redacted,omitzero"`
}

// ConfigExplanation is the JSON contract of `config explain`. It reports the
// root config, because a config file can only feed root flags; command-level
// flags have no config-file source.
//
// SchemaVersion follows the same additive-only rule as [CLIDescription].
type ConfigExplanation struct {
	SchemaVersion string         `json:"schema_version"`
	ConfigFile    string         `json:"config_file,omitzero"`
	Flags         []ResolvedFlag `json:"flags"`
}

// ExplainConfig reports how every root-config flag's value was decided.
// Call it after Execute for a description of that run, or before Execute for a
// description of the constructor defaults plus the config file.
func (cli *CLI[T]) ExplainConfig() ConfigExplanation {
	tags := cli.registry.Tags()
	declared := cli.registry.DeclaredDefaults()

	field := reflect.Value{}

	if v, err := derefPointerToStruct(cli.config); err == nil {
		field = v
	}

	flags := make([]ResolvedFlag, 0, len(tags))

	for i, tag := range tags {
		flagDefault := ""
		if i < len(declared) {
			flagDefault = declared[i]
		}

		flags = append(flags, cli.resolveFlag(tag, flagDefault, field))
	}

	return ConfigExplanation{
		SchemaVersion: descriptionSchemaVersion,
		ConfigFile:    cli.configFilePath,
		Flags:         flags,
	}
}

// resolveFlag picks the highest-precedence source that supplied tag's value,
// mirroring the order [FlagRegistry.parseFlag] uses.
func (cli *CLI[T]) resolveFlag(
	tag FlagTag,
	declaredDefault string,
	cfgValue reflect.Value,
) ResolvedFlag {
	resolved := ResolvedFlag{
		Name:    tag.Name,
		Key:     tag.Key,
		Default: declaredDefault,
		Env:     cli.registry.envName(tag.Env),
		Secret:  tag.Secret,
	}

	effective := liveFieldValue(cfgValue, tag)

	if flag, ok := explicitFlagValue(cli.rootCmd, tag.Name); ok {
		resolved.Source = SourceFlag
		resolved.Value = flag
	} else if envValue, ok := lookupEnvValue(resolved.Env); ok {
		resolved.Source = SourceEnvironment
		resolved.Value = envValue
	} else if cli.configSuppliedField(tag.Field) {
		resolved.Source = SourceConfigFile
		resolved.Value = effective
	} else if declaredDefault != "" || effective != "" {
		resolved.Source = SourceDefault
		// The declared default is authoritative here: before Execute the struct
		// still holds the constructor's zero values, and after Execute the field
		// equals the default anyway. A non-zero constructor value is the fallback
		// for fields whose tag declares no default.
		resolved.Value = cmp.Or(declaredDefault, effective)
	} else {
		resolved.Source = SourceUnset
	}

	if tag.Secret {
		resolved.Value = RedactedValue
		resolved.Redacted = true
	}

	return resolved
}

// configSuppliedField reports whether the config file set field.
func (cli *CLI[T]) configSuppliedField(field string) bool {
	_, ok := cli.configSetFields[field]

	return ok
}

// envName applies the registry's prefix to an env tag, matching what flag
// parsing reads.
func (r *FlagRegistry) envName(env string) string {
	if env == "" || r.envPrefix == "" {
		return env
	}

	return r.envPrefix + env
}

// lookupEnvValue reads an environment variable, reporting whether it is set at
// all (an empty but set variable still counts as a source).
func lookupEnvValue(name string) (string, bool) {
	if name == "" {
		return "", false
	}

	return os.LookupEnv(name)
}

// explicitFlagValue returns a flag's value when the current invocation set it
// explicitly. pflag's Changed flag is the only reliable signal: comparing the
// value against the default cannot distinguish "--port 8080" from no flag.
func explicitFlagValue(cmd *cobra.Command, name string) (string, bool) {
	flag := lookupFlagInCommand(cmd, name)
	if flag == nil || !flag.Changed {
		return "", false
	}

	return flag.Value.String(), true
}

// liveFieldValue renders the field a tag addresses, falling back to the tag
// default when the struct value is unavailable.
func liveFieldValue(cfgValue reflect.Value, tag FlagTag) string {
	if !cfgValue.IsValid() {
		return tag.Default
	}

	field := fieldByTag(cfgValue, tag)
	if !field.IsValid() {
		return tag.Default
	}

	return formatFieldValue(field)
}

// ConfigExplainCommand builds the "config explain" command. It prints how each
// root-config value was decided: flag, environment, config file, or default.
//
// Register it with [AddCommand] or opt in with [WithConfigExplainCommand], which
// registers it at Execute time and yields to a consumer-defined "config" command.
func ConfigExplainCommand[T any](cli *CLI[T]) (Command[T, NoFlags], error) {
	appName := cli.spec.name

	explain, err := NewCommand(
		"explain",
		NoFlags{},
		func(_ context.Context, _ *T, _ NoFlags) error {
			return cli.writeConfigExplanation(cli.rootCmd.OutOrStdout())
		},
		WithShort("Explain where each configuration value came from"),
		WithLong(
			"Print every root configuration flag with the source that decided its effective value: "+
				"a command-line flag, an environment variable, a config file, or the declared default. "+
				"Values tagged secret:\"true\" are redacted.",
		),
		WithExample(appName+" config explain"),
		WithNoArgs(),
	)
	if err != nil {
		return Command[T, NoFlags]{}, fmt.Errorf("building %q command: %w", "config explain", err)
	}

	parent, err := NewParentCommand[T, NoFlags](
		configExplainCommandName,
		"Inspect how this CLI's configuration was resolved",
		NoFlags{},
		WithShort("Inspect configuration resolution"),
		WithSubcommands(explain),
	)
	if err != nil {
		return Command[T, NoFlags]{}, fmt.Errorf("building %q command: %w", configExplainCommandName, err)
	}

	return parent, nil
}

// ensureConfigExplainCommand registers the config command opted into by
// [WithConfigExplainCommand]. Registration is idempotent and yields to a
// consumer-defined command of the same name.
func (cli *CLI[T]) ensureConfigExplainCommand() error {
	if cli.configExplainWired || !cli.spec.configExplainCmd {
		return nil
	}

	if cli.hasCommand(configExplainCommandName) {
		return nil
	}

	cmd, err := ConfigExplainCommand(cli)
	if err != nil {
		return err
	}

	err = AddCommand(cli, cmd)
	if err != nil {
		return fmt.Errorf("registering %q command: %w", configExplainCommandName, err)
	}

	cli.configExplainWired = true

	return nil
}

// writeConfigExplanation renders the explanation as JSON for machine-readable
// output formats and as an aligned table otherwise.
func (cli *CLI[T]) writeConfigExplanation(w io.Writer) error {
	explanation := cli.ExplainConfig()

	if cli.OutputFormat() == output.FormatJSON {
		return writeConfigExplanationJSON(w, explanation)
	}

	return writeConfigExplanationTable(w, explanation)
}

// writeConfigExplanationJSON writes the versioned JSON contract.
func writeConfigExplanationJSON(w io.Writer, explanation ConfigExplanation) error {
	enc := jsontext.NewEncoder(w, jsontext.WithIndent("  "))

	err := json.MarshalEncode(enc, explanation)
	if err != nil {
		return fmt.Errorf("encoding config explanation: %w", err)
	}

	return nil
}

// writeConfigExplanationTable writes one aligned row per flag. Columns are
// ordered by what a reader asks first: which value won, then where it came from.
func writeConfigExplanationTable(w io.Writer, explanation ConfigExplanation) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	_, err := fmt.Fprintln(tw, "FLAG\tSOURCE\tVALUE\tDEFAULT\tENV")
	if err != nil {
		return fmt.Errorf("writing config explanation: %w", err)
	}

	for _, flag := range explanation.Flags {
		_, err = fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%s\t%s\n",
			flag.Name,
			flag.Source,
			displayValue(flag.Value),
			displayValue(flag.Default),
			displayValue(flag.Env),
		)
		if err != nil {
			return fmt.Errorf("writing config explanation: %w", err)
		}
	}

	err = tw.Flush()
	if err != nil {
		return fmt.Errorf("flushing config explanation: %w", err)
	}

	if explanation.ConfigFile != "" {
		_, err = fmt.Fprintf(w, "\nconfig file: %s\n", explanation.ConfigFile)
		if err != nil {
			return fmt.Errorf("writing config file path: %w", err)
		}
	}

	return nil
}

// displayValue renders an empty value as "-" so columns stay readable.
func displayValue(value string) string {
	if value == "" {
		return "-"
	}

	return strings.ReplaceAll(value, "\n", " ")
}
