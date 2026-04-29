package v2

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ParseFlags populates a config struct from parsed flags.
func (r *FlagRegistry) ParseFlags(cmd *cobra.Command, cfg any) error {
	for _, tag := range r.tags {
		err := r.parseFlag(cmd, cfg, tag)
		if err != nil {
			return fmt.Errorf("parsing flags for command %q: %w", cmd.Use, err)
		}
	}

	return nil
}

// parseFlag reads a flag value and sets it on the config struct.
// Priority: explicit flag > environment variable > default value.
func (r *FlagRegistry) parseFlag(cmd *cobra.Command, cfg any, tag FlagTag) error {
	flag, err := r.lookupFlag(cmd, tag)
	if err != nil {
		return fmt.Errorf("looking up flag %q on command %q: %w", tag.Name, cmd.Use, err)
	}

	var value string
	hasValue := false

	// Priority 1: explicit flag value
	if flag.Changed {
		value = flag.Value.String()
		hasValue = true
	}

	// Priority 2: environment variable
	if !hasValue && tag.Env != "" {
		envName := tag.Env
		if r.envPrefix != "" {
			envName = r.envPrefix + envName
		}

		if envVal, ok := os.LookupEnv(envName); ok {
			value = envVal
			hasValue = true
		}
	}

	// Priority 3: default value
	if !hasValue {
		if tag.Default == "" {
			return nil
		}

		value = tag.Default
	}

	return r.parseAndSetValue(cfg, tag, value)
}

// lookupFlag finds a flag in the command.
func (r *FlagRegistry) lookupFlag(cmd *cobra.Command, tag FlagTag) (*pflag.Flag, error) {
	flag := cmd.Flags().Lookup(tag.Name)
	if flag == nil {
		// Try persistent flags
		flag = cmd.PersistentFlags().Lookup(tag.Name)
	}

	if flag == nil {
		return nil, fmt.Errorf(
			"flag %q not found in command %q: %w",
			tag.Name,
			cmd.Use,
			ErrFlagNotFound,
		)
	}

	return flag, nil
}

// parseAndSetValue parses the flag value via the TypeHandler registry and sets it on config.
func (r *FlagRegistry) parseAndSetValue(cfg any, tag FlagTag, value string) error {
	parsed, err := dispatchParse(value, tag)
	if err != nil {
		return fmt.Errorf("parsing flag %q with value %q: %w", tag.Name, value, err)
	}

	// Handle []string from slice parsing
	if sl, ok := parsed.([]string); ok {
		return SetField(cfg, tag.Field, sl)
	}

	// Handle type conversion via reflect
	parsedVal := reflect.ValueOf(parsed)
	fieldVal := reflect.ValueOf(cfg)
	if fieldVal.Kind() == reflect.Pointer {
		fieldVal = fieldVal.Elem()
	}
	field := fieldVal.FieldByName(tag.Field)
	if !field.IsValid() {
		return fmt.Errorf("field %q not found in %T: %w", tag.Field, cfg, ErrFieldNotFound)
	}

	// Use ConvertibleTo for numeric narrowing
	if parsedVal.Type().ConvertibleTo(field.Type()) {
		field.Set(parsedVal.Convert(field.Type()))
		return nil
	}

	return SetField(cfg, tag.Field, parsed)
}
