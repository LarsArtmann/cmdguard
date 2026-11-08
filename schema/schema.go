// Package schema turns a cmdguard config struct into a JSON Schema document.
//
// The point is to let editors, CI pipelines, and non-Go tools validate a config
// file without running the binary. The schema is derived from the same struct
// tags that drive flags and strict config loading ([v5.ConfigKeys]), so it
// cannot accept a key the loader rejects or miss one it accepts.
//
// This package deliberately lives outside the core module: it is a build-time
// code generator, not runtime machinery, and core has no reason to carry it.
//
// Validation at runtime is already covered by [cmdguard.WithStrictConfigKeys];
// this package is for everything that happens before the binary runs.
package schema

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

// Dialect is the JSON Schema dialect emitted. Draft 2020-12 is what editors and
// CI validators support today.
const Dialect = "https://json-schema.org/draft/2020-12/schema"

// ExtensionKey names the non-standard hints cmdguard adds to a property. It is
// namespaced so a validator ignores it and a tool can read it.
const ExtensionKey = "x-cmdguard"

// PropertyMap is a set of named properties. Plain Go maps marshal in map
// iteration order under encoding/json/v2, so this type sorts its keys:
// generated schemas are byte-for-byte reproducible.
type PropertyMap map[string]*Property

func (m PropertyMap) MarshalJSONTo(out *jsontext.Encoder) error {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}

	slices.Sort(names)

	if err := out.WriteToken(jsontext.BeginObject); err != nil {
		return fmt.Errorf("opening schema properties object: %w", err)
	}

	for _, name := range names {
		if err := out.WriteToken(jsontext.String(name)); err != nil {
			return fmt.Errorf("writing schema property name %q: %w", name, err)
		}

		if err := json.MarshalEncode(out, m[name]); err != nil {
			return fmt.Errorf("encoding schema property %q: %w", name, err)
		}
	}

	if err := out.WriteToken(jsontext.EndObject); err != nil {
		return fmt.Errorf("closing schema properties object: %w", err)
	}

	return nil
}

// Document is a JSON Schema describing a cmdguard config file.
type Document struct {
	Schema               string      `json:"$schema"`
	Title                string      `json:"title,omitempty"`
	Description          string      `json:"description,omitempty"`
	Type                 string      `json:"type"`
	Properties           PropertyMap `json:"properties"`
	Required             []string    `json:"required,omitempty"`
	AdditionalProperties bool        `json:"additionalProperties"`
}

// Property describes one accepted key. Standard fields carry the constraints a
// validator enforces; Hints carries what only cmdguard knows.
type Property struct {
	Type                 string      `json:"type,omitempty"`
	Description          string      `json:"description,omitempty"`
	Default              any         `json:"default,omitempty"`
	Enum                 []string    `json:"enum,omitempty"`
	Minimum              *float64    `json:"minimum,omitempty"`
	Maximum              *float64    `json:"maximum,omitempty"`
	MinLength            *int        `json:"minLength,omitempty"`
	MaxLength            *int        `json:"maxLength,omitempty"`
	Pattern              string      `json:"pattern,omitempty"`
	Format               string      `json:"format,omitempty"`
	Properties           PropertyMap `json:"properties,omitempty"`
	Required             []string    `json:"required,omitempty"`
	Items                *Property   `json:"items,omitempty"`
	AdditionalProperties *bool       `json:"additionalProperties,omitempty"`
	Hints                *Hints      `json:"x-cmdguard,omitempty"`
}

// Hints carries metadata JSON Schema has no home for.
type Hints struct {
	// Flag is the command-line flag that sets the same value.
	Flag string `json:"flag,omitempty"`
	// Env is the environment variable that sets the same value.
	Env string `json:"env,omitempty"`
	// Format is a cmdguard-specific type hint (duration, hostport, path) for
	// values JSON Schema only sees as strings.
	Format string `json:"format,omitempty"`
	// Secret marks a value that must never be logged or committed.
	Secret bool `json:"secret,omitempty"`
}

// FromConfig builds the schema for a cmdguard config struct.
//
// The result rejects unknown keys (additionalProperties: false) to match strict
// config loading: a typo caught by the schema is the same typo the binary would
// reject.
func FromConfig(cfg any) (*Document, error) {
	return fromConfigWithIdentity(cfg, "", "")
}

// FromCLI builds the schema for a CLI's config, using the CLI's name and short
// description as the document's title and description.
func FromCLI[T any](cli *v5.CLI[T]) (*Document, error) {
	if cli == nil {
		return nil, fmt.Errorf("building config schema: %w", v5.ErrConfigNil)
	}

	return fromConfigWithIdentity(cli.Config(), cli.Name(), cli.Short())
}

// Emit writes the schema as indented JSON to writer.
func Emit(writer io.Writer, cfg any) error {
	document, err := FromConfig(cfg)
	if err != nil {
		return err
	}

	return encode(writer, document)
}

// Marshal returns the schema as indented JSON.
func Marshal(cfg any) ([]byte, error) {
	document, err := FromConfig(cfg)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = encode(&buf, document)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// encode writes a document as indented JSON.
func encode(writer io.Writer, document *Document) error {
	enc := jsontext.NewEncoder(writer, jsontext.WithIndent("  "))

	err := json.MarshalEncode(enc, document)
	if err != nil {
		return fmt.Errorf("encoding config schema: %w", err)
	}

	return nil
}

// fromConfigWithIdentity builds a document and optionally names it.
func fromConfigWithIdentity(cfg any, title, description string) (*Document, error) {
	keys, err := v5.ConfigKeys(cfg)
	if err != nil {
		return nil, fmt.Errorf("reading config keys: %w", err)
	}

	properties, required := propertiesFor(keys)

	return &Document{
		Schema:               Dialect,
		Title:                title,
		Description:          description,
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: false,
	}, nil
}

// propertiesFor converts config keys into schema properties, collecting the
// required ones in declaration order.
func propertiesFor(keys []v5.ConfigKey) (PropertyMap, []string) {
	properties := make(PropertyMap, len(keys))
	required := make([]string, 0, len(keys))

	for _, key := range keys {
		properties[key.Key] = propertyFor(key)

		if key.Required {
			required = append(required, key.Key)
		}
	}

	return properties, required
}

// propertyFor converts one config key into a schema property.
func propertyFor(key v5.ConfigKey) *Property {
	property := &Property{
		Type:        key.Type,
		Description: key.Help,
		Enum:        key.Values,
		Hints:       hintsFor(key),
	}

	if key.Format != "" {
		if isStandardFormat(key.Format) {
			property.Format = key.Format
		}
	}

	applyValidate(property, key.Type, key.Validate)
	property.Default = defaultFor(key.Type, key.Default)

	if len(key.Properties) > 0 {
		property.Properties, property.Required = propertiesFor(key.Properties)

		closed := false
		property.AdditionalProperties = &closed
	}

	if key.Items != nil {
		property.Items = propertyFor(*key.Items)
	}

	return property
}

// hintsFor builds the cmdguard extension, returning nil when there is nothing
// to say so the emitted schema stays clean.
func hintsFor(key v5.ConfigKey) *Hints {
	hints := &Hints{Flag: key.Flag, Env: key.Env, Secret: key.Secret}

	if !isStandardFormat(key.Format) {
		hints.Format = key.Format
	}

	if *hints == (Hints{}) {
		return nil
	}

	return hints
}

// isStandardFormat reports whether JSON Schema defines the format name, so it
// belongs in "format" rather than the cmdguard extension.
func isStandardFormat(format string) bool {
	switch format {
	case "email", "uri", "date-time":
		return true
	default:
		return false
	}
}

// defaultFor converts a struct-tag default into the JSON value a validator and
// an editor expect. An unparseable default is emitted as the raw string rather
// than dropped: the tag is the authority, and a wrong-shaped default is a bug
// worth seeing.
func defaultFor(jsonType, raw string) any {
	if raw == "" {
		return nil
	}

	switch jsonType {
	case "boolean":
		if value, err := strconv.ParseBool(raw); err == nil {
			return value
		}
	case "integer":
		if value, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return value
		}
	case "number":
		if value, err := strconv.ParseFloat(raw, 64); err == nil {
			return value
		}
	case "array":
		return strings.Split(raw, ",")
	case "object", "string":
		return raw
	}

	return raw
}

// applyValidate translates cmdguard's validate tag into schema constraints.
// cmdguard's own semantics are preserved: on string fields min/max mean length,
// which is why the value kind decides the target.
func applyValidate(property *Property, jsonType, validate string) {
	if validate == "" {
		return
	}

	for rule := range strings.SplitSeq(validate, ",") {
		name, value, _ := strings.Cut(strings.TrimSpace(rule), "=")
		applyValidateRule(property, jsonType, name, value)
	}
}

// applyValidateRule applies one validate rule. Unknown rules are ignored: a
// consumer-registered validator cannot be expressed in JSON Schema, and failing
// the schema build for it would be worse than documenting the rest.
func applyValidateRule(property *Property, jsonType, name, value string) {
	switch name {
	case "min":
		applyBound(property, jsonType, value, true)
	case "max":
		applyBound(property, jsonType, value, false)
	case "minlen":
		if length, err := strconv.Atoi(value); err == nil {
			property.MinLength = &length
		}
	case "maxlen":
		if length, err := strconv.Atoi(value); err == nil {
			property.MaxLength = &length
		}
	case "nonempty":
		one := 1
		property.MinLength = &one
	case "regex":
		property.Pattern = value
	case "email":
		property.Format = "email"
	case "url":
		property.Format = "uri"
	}
}

// applyBound routes a min/max rule to length or numeric bounds, matching what
// cmdguard validation does.
func applyBound(property *Property, jsonType, value string, isMin bool) {
	if jsonType == "string" || jsonType == "array" {
		length, err := strconv.Atoi(value)
		if err != nil {
			return
		}

		if isMin {
			property.MinLength = &length

			return
		}

		property.MaxLength = &length

		return
	}

	bound, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}

	if isMin {
		property.Minimum = &bound

		return
	}

	property.Maximum = &bound
}

// CommandSchemas builds one flag schema per command path from the machine
// contract. Keys are command paths ("app", "app deploy"); each document
// describes the flags that command accepts, including inherited ones.
//
// This is for tooling that drives a binary: a config schema constrains a file,
// a flag schema constrains an invocation.
func CommandSchemas(description *v5.CLIDescription) map[string]*Document {
	if description == nil {
		return nil
	}

	schemas := make(map[string]*Document)

	schemas[description.Name] = commandSchema(description.Name, description.Short, description.Flags)

	collectCommandSchemas(description.Commands, schemas)

	return schemas
}

// commandSchema builds the flag schema for one command path.
func commandSchema(path, short string, flags []v5.FlagDescription) *Document {
	properties, required := flagProperties(flags)

	return &Document{
		Schema:               Dialect,
		Title:                path,
		Description:          short,
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: false,
	}
}

// collectCommandSchemas records every command in the tree, at any depth.
func collectCommandSchemas(commands []v5.CommandDescription, schemas map[string]*Document) {
	for _, command := range commands {
		schemas[command.Path] = commandSchema(command.Path, command.Short, command.Flags)
		collectCommandSchemas(command.Commands, schemas)
	}
}

// flagProperties converts flag descriptions into properties keyed by flag name.
func flagProperties(flags []v5.FlagDescription) (map[string]*Property, []string) {
	properties := make(map[string]*Property, len(flags))
	required := make([]string, 0, len(flags))

	for _, flag := range flags {
		jsonType := jsonTypeForFlag(flag.Type)

		properties[flag.Name] = &Property{
			Type:        jsonType,
			Description: flag.Usage,
			Default:     defaultFor(jsonType, flag.Default),
			Enum:        flag.Values,
			Hints:       &Hints{Env: flag.Env},
		}

		if flag.Required {
			required = append(required, flag.Name)
		}
	}

	return properties, required
}

// jsonTypeForFlag maps a pflag type name to a JSON type.
func jsonTypeForFlag(pflagType string) string {
	switch pflagType {
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "count":
		return "integer"
	case "float32", "float64":
		return "number"
	case "stringSlice", "stringArray", "intSlice", "int32Slice", "int64Slice",
		"uintSlice", "float32Slice", "float64Slice":
		return "array"
	default:
		return "string"
	}
}
