package v4

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
)

// PromptRunner defines the interface for interactive prompts.
// The optional prompts module (github.com/larsartmann/cmdguard/prompts)
// provides a huh/v2 implementation. Without it, prompting returns an error.
type PromptRunner interface {
	PromptString(title, defaultValue string) (string, error)
	PromptSelect(title string, options []string) (string, error)
	PromptConfirm(title string) (bool, error)
}

// defaultPromptRunner is the package-level prompt implementation.
// nil means prompting is disabled (prompts module not registered).
//
//nolint:gochecknoglobals // package-level hook for optional prompts module
var defaultPromptRunner PromptRunner

// SetPromptRunner registers a prompt implementation. The optional prompts
// module calls this during Register(). When nil (the default), prompt-related
// features return ErrPromptNotRegistered.
func SetPromptRunner(r PromptRunner) {
	defaultPromptRunner = r
}

// ErrPromptNotRegistered is returned when prompt features are used without
// the prompts module registered.
var ErrPromptNotRegistered = errors.New(
	"prompts module not registered: import github.com/larsartmann/cmdguard/prompts and call Register()",
)

// PromptString prompts the user for a string value.
func PromptString(title, defaultValue string) (string, error) {
	if defaultPromptRunner == nil {
		return "", fmt.Errorf("title=%q: %w", title, ErrPromptNotRegistered)
	}

	result, err := defaultPromptRunner.PromptString(title, defaultValue)
	if err != nil {
		return "", fmt.Errorf("title=%q, defaultValue=%q: prompting for string: %w", title, defaultValue, err)
	}

	return result, nil
}

// PromptSelect prompts the user to select from a list of options.
func PromptSelect(title string, options []string) (string, error) {
	if defaultPromptRunner == nil {
		return "", fmt.Errorf("title=%q: %w", title, ErrPromptNotRegistered)
	}

	result, err := defaultPromptRunner.PromptSelect(title, options)
	if err != nil {
		return "", fmt.Errorf("title=%q, options=%v: prompting for selection: %w", title, options, err)
	}

	return result, nil
}

// PromptConfirm prompts the user for a yes/no confirmation.
func PromptConfirm(title string) (bool, error) {
	if defaultPromptRunner == nil {
		return false, fmt.Errorf("title=%q: %w", title, ErrPromptNotRegistered)
	}

	result, err := defaultPromptRunner.PromptConfirm(title)
	if err != nil {
		return false, fmt.Errorf("title=%q: prompting for confirmation: %w", title, err)
	}

	return result, nil
}

// promptMissingCommandFlags interactively prompts for any command-level flags
// that have a prompt tag and were not explicitly provided via CLI arguments or
// environment variables.
func promptMissingCommandFlags(c *cobra.Command, registry *FlagRegistry) error {
	if registry == nil || defaultPromptRunner == nil {
		return nil
	}

	for _, tag := range registry.tags {
		if tag.Prompt == "" {
			continue
		}

		flag := c.Flags().Lookup(tag.Name)
		if flag == nil {
			continue
		}

		if flag.Changed {
			continue
		}

		if isEnvSet(registry, tag) {
			continue
		}

		var value string

		switch {
		case len(tag.Values) > 0:
			selected, err := PromptSelect(tag.Prompt, tag.Values)
			if err != nil {
				return fmt.Errorf("flag=%q, value=%q: prompting for flag %q: %w", tag.Name, tag.Default, tag.Name, err)
			}

			value = selected
		case tag.Type.Kind() == reflect.Bool:
			confirmed, err := PromptConfirm(tag.Prompt)
			if err != nil {
				return fmt.Errorf("flag=%q, value=%q: prompting for flag %q: %w", tag.Name, tag.Default, tag.Name, err)
			}

			value = strconv.FormatBool(confirmed)
		default:
			result, err := PromptString(tag.Prompt, tag.Default)
			if err != nil {
				return fmt.Errorf("flag=%q, value=%q: prompting for flag %q: %w", tag.Name, tag.Default, tag.Name, err)
			}

			value = result
		}

		err := c.Flags().Set(tag.Name, value)
		if err != nil {
			return fmt.Errorf(
				"flag=%q, value=%q: setting prompted value for flag %q: %w",
				tag.Name,
				value,
				tag.Name,
				err,
			)
		}
	}

	return nil
}

func isEnvSet(registry *FlagRegistry, tag FlagTag) bool {
	if tag.Env == "" {
		return false
	}

	envName := tag.Env
	if registry.envPrefix != "" {
		envName = registry.envPrefix + envName
	}

	_, ok := os.LookupEnv(envName)

	return ok
}
