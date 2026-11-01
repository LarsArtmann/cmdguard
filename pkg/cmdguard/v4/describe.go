package v4

import (
	"context"
	"cmp"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// descriptionSchemaVersion is the version of the JSON contract emitted by
// Describe and WriteDescription. Within a major version only additive changes
// are allowed: new fields may appear, existing field names and meanings may
// not change. Consumers (including AI agents) pin this value.
const descriptionSchemaVersion = "1"

// flagMetadataAnnotation stores struct-tag metadata that pflag does not model
// (env, required, enum values, count, prompt, validate) on the cobra command
// that owns the flags. Describe reads it back so the machine-readable contract
// carries everything a caller needs to construct a valid invocation.
const flagMetadataAnnotation = "cmdguard.io/flag-metadata"

// describeCommandName is the reserved command name for the auto-registered
// JSON description command. A consumer-defined command with this name wins.
const describeCommandName = "describe"

// helpCommandName is cobra's generated help command, which carries no behavior
// of its own and is therefore not part of the description contract.
const helpCommandName = "help"

// CLIDescription is the machine-readable contract of a cmdguard CLI: every
// command, every accepted flag, and the metadata needed to drive the binary
// without reading its source or scraping its prose help.
//
// SchemaVersion is stable within a major: fields may be added, never renamed
// or repurposed. Treat unknown fields as additive.
type CLIDescription struct {
	SchemaVersion string               `json:"schema_version"`
	Name          string               `json:"name"`
	Short         string               `json:"short"`
	Long          string               `json:"long,omitempty"`
	Version       string               `json:"version,omitempty"`
	Flags         []FlagDescription    `json:"flags,omitempty"`
	Commands      []CommandDescription `json:"commands,omitempty"`
}

// CommandDescription describes one command in the CLI tree.
type CommandDescription struct {
	Name       string               `json:"name"`
	Path       string               `json:"path"`
	Short      string               `json:"short"`
	Long       string               `json:"long,omitempty"`
	Example    string               `json:"example,omitempty"`
	Aliases    []string             `json:"aliases,omitempty"`
	Group      string               `json:"group,omitempty"`
	Hidden     bool                 `json:"hidden,omitempty"`
	Deprecated string               `json:"deprecated,omitempty"`
	Flags      []FlagDescription    `json:"flags,omitempty"`
	Commands   []CommandDescription `json:"commands,omitempty"`
}

// FlagDescription describes one accepted flag, including the struct-tag
// metadata (env, required, enum values) that determines a valid invocation.
type FlagDescription struct {
	Name       string   `json:"name"`
	Shorthand  string   `json:"shorthand,omitempty"`
	Type       string   `json:"type"`
	Default    string   `json:"default,omitempty"`
	Usage      string   `json:"usage"`
	Env        string   `json:"env,omitempty"`
	Values     []string `json:"values,omitempty"`
	Required   bool     `json:"required,omitempty"`
	Count      bool     `json:"count,omitempty"`
	Prompt     string   `json:"prompt,omitempty"`
	Validate   string   `json:"validate,omitempty"`
	Hidden     bool     `json:"hidden,omitempty"`
	Deprecated string   `json:"deprecated,omitempty"`
	Inherited  bool     `json:"inherited,omitempty"`
}

// flagMetadata is the internal, annotation-encoded subset of FlagTag that
// Describe needs. It is deliberately decoupled from FlagTag so the persisted
// annotation shape is an explicit contract rather than an accident of the
// struct-tag parser's field set.
type flagMetadata struct {
	Name     string   `json:"name"`
	Env      string   `json:"env,omitempty"`
	Values   []string `json:"values,omitempty"`
	Required bool     `json:"required,omitempty"`
	Count    bool     `json:"count,omitempty"`
	Prompt   string   `json:"prompt,omitempty"`
	Validate string   `json:"validate,omitempty"`
}

// Describe returns a machine-readable description of the CLI command tree.
//
// The result is stable JSON-shaped data that a caller (an AI agent, a test, a
// documentation generator) can use to drive the binary without scraping help
// text. Describe is safe to call before Execute; it initializes cobra's
// default help and version flags so the contract lists them.
func (cli *CLI[T]) Describe() *CLIDescription {
	cli.rootCmd.InitDefaultHelpFlag()
	cli.rootCmd.InitDefaultVersionFlag()

	return &CLIDescription{
		SchemaVersion: descriptionSchemaVersion,
		Name:          cli.spec.name,
		Short:         cli.spec.short,
		Long:          cli.spec.long,
		Version:       cli.spec.version,
		Flags:         flagsFor(cli.rootCmd),
		Commands:      subcommandDescriptions(cli.rootCmd),
	}
}

// WriteDescription writes the CLI description as indented JSON to w.
func (cli *CLI[T]) WriteDescription(w io.Writer) error {
	enc := jsontext.NewEncoder(w, jsontext.WithIndent("  "))

	err := json.MarshalEncode(enc, cli.Describe())
	if err != nil {
		return fmt.Errorf("encoding CLI description: %w", err)
	}

	return nil
}

// ensureDescribeCommand registers the hidden "describe" command that prints the
// machine-readable contract. Registration is idempotent, skipped when disabled,
// and yields to a consumer-defined command of the same name.
func (cli *CLI[T]) ensureDescribeCommand() error {
	if cli.describeWired || cli.spec.noDescribeCmd {
		return nil
	}

	for _, existing := range cli.rootCmd.Commands() {
		if existing.Name() == describeCommandName {
			return nil
		}
	}

	cmd, err := NewCommand(
		describeCommandName,
		NoFlags{},
		func(_ context.Context, _ *T, _ NoFlags) error {
			return cli.WriteDescription(cli.rootCmd.OutOrStdout())
		},
		WithShort("Print a machine-readable JSON description of this CLI"),
		WithLong(
			"Print the full command and flag tree as JSON. "+
				"Designed for AI agents and tooling that need to drive this CLI without parsing help text.",
		),
		WithExample(cli.spec.name+" "+describeCommandName),
		WithHidden(true),
		WithNoArgs(),
	)
	if err != nil {
		return fmt.Errorf("building %q command: %w", describeCommandName, err)
	}

	err = AddCommand(cli, cmd)
	if err != nil {
		return fmt.Errorf("registering %q command: %w", describeCommandName, err)
	}

	cli.describeWired = true

	return nil
}

// attachFlagMetadata encodes registry tag metadata onto cmd's annotations so
// Describe can report env/required/values without retaining the registry.
func attachFlagMetadata(cmd *cobra.Command, registry *FlagRegistry) error {
	if registry == nil {
		return nil
	}

	tags := registry.Tags()
	if len(tags) == 0 {
		return nil
	}

	metas := make([]flagMetadata, 0, len(tags))

	for _, tag := range tags {
		env := tag.Env
		if env != "" && registry.envPrefix != "" {
			env = registry.envPrefix + env
		}

		metas = append(metas, flagMetadata{
			Name:     tag.Name,
			Env:      env,
			Values:   slices.Clone(tag.Values),
			Required: tag.Required,
			Count:    tag.Count,
			Prompt:   tag.Prompt,
			Validate: tag.Validate,
		})
	}

	encoded, err := json.Marshal(metas)
	if err != nil {
		return fmt.Errorf("encoding flag metadata for command %q: %w", cmd.Name(), err)
	}

	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	cmd.Annotations[flagMetadataAnnotation] = string(encoded)

	return nil
}

// flagMetadataFor decodes cmd's annotation-encoded tag metadata, keyed by flag
// name. Malformed metadata is treated as absent: annotations are internal,
// machine-written data, and Describe must never fail because of them.
func flagMetadataFor(cmd *cobra.Command) map[string]flagMetadata {
	encoded, ok := cmd.Annotations[flagMetadataAnnotation]
	if !ok {
		return nil
	}

	var metas []flagMetadata
	if err := json.Unmarshal([]byte(encoded), &metas); err != nil {
		return nil
	}

	out := make(map[string]flagMetadata, len(metas))
	for _, meta := range metas {
		out[meta.Name] = meta
	}

	return out
}

// subcommandDescriptions returns descriptions of cmd's subcommands, excluding
// cobra's generated help command (which carries no behavior of its own).
func subcommandDescriptions(cmd *cobra.Command) []CommandDescription {
	subs := cmd.Commands()
	out := make([]CommandDescription, 0, len(subs))

	for _, sub := range subs {
		if sub.Name() == helpCommandName {
			continue
		}

		out = append(out, commandDescription(sub))
	}

	return out
}

// commandDescription describes one command and, recursively, its subtree.
func commandDescription(cmd *cobra.Command) CommandDescription {
	cmd.InitDefaultHelpFlag()

	desc := CommandDescription{
		Name:       cmd.Name(),
		Path:       cmd.CommandPath(),
		Short:      cmd.Short,
		Long:       cmd.Long,
		Example:    cmd.Example,
		Aliases:    slices.Clone(cmd.Aliases),
		Group:      cmd.GroupID,
		Hidden:     cmd.Hidden,
		Deprecated: cmd.Deprecated,
		Flags:      flagsFor(cmd),
		Commands:   subcommandDescriptions(cmd),
	}

	return desc
}

// flagsFor returns every flag accepted by cmd, local and inherited, sorted by
// name for a deterministic contract.
func flagsFor(cmd *cobra.Command) []FlagDescription {
	meta := flagMetadataFor(cmd)
	seen := make(map[string]struct{})
	out := make([]FlagDescription, 0)

	appendSet := func(flags *pflag.FlagSet, inherited bool) {
		flags.VisitAll(func(flag *pflag.Flag) {
			if _, ok := seen[flag.Name]; ok {
				return
			}

			seen[flag.Name] = struct{}{}

			out = append(out, flagDescription(flag, meta[flag.Name], inherited))
		})
	}

	appendSet(cmd.LocalFlags(), false)
	appendSet(cmd.InheritedFlags(), true)

	slices.SortFunc(out, func(a, b FlagDescription) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return out
}

// flagDescription converts one pflag plus its optional tag metadata into the
// contract shape.
func flagDescription(flag *pflag.Flag, meta flagMetadata, inherited bool) FlagDescription {
	return FlagDescription{
		Name:       flag.Name,
		Shorthand:  flag.Shorthand,
		Type:       flag.Value.Type(),
		Default:    flag.DefValue,
		Usage:      flag.Usage,
		Env:        meta.Env,
		Values:     meta.Values,
		Required:   meta.Required,
		Count:      meta.Count,
		Prompt:     meta.Prompt,
		Validate:   meta.Validate,
		Hidden:     flag.Hidden,
		Deprecated: flag.Deprecated,
		Inherited:  inherited,
	}
}
