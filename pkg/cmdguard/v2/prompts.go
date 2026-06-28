package v2

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

// PromptRunner defines the interface for interactive prompts.
// Tests can override the default implementation by replacing defaultPromptRunner.
type PromptRunner interface {
	PromptString(title, defaultValue string) (string, error)
	PromptSelect(title string, options []string) (string, error)
	PromptConfirm(title string) (bool, error)
}

type huhPromptRunner struct{}

func (h *huhPromptRunner) PromptString(title, defaultValue string) (string, error) {
	result := defaultValue

	err := huh.NewInput().
		Title(title).
		Value(&result).
		Run()
	if err != nil {
		return "", fmt.Errorf("title=%q, defaultValue=%q: running string prompt: %w", title, defaultValue, err)
	}

	return result, nil
}

func (h *huhPromptRunner) PromptSelect(title string, options []string) (string, error) {
	var result string

	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}

	err := huh.NewSelect[string]().
		Title(title).
		Options(opts...).
		Value(&result).
		Run()
	if err != nil {
		return "", fmt.Errorf("title=%q, options=%v: running select prompt: %w", title, options, err)
	}

	return result, nil
}

func (h *huhPromptRunner) PromptConfirm(title string) (bool, error) {
	var result bool

	err := huh.NewConfirm().
		Title(title).
		Value(&result).
		Run()
	if err != nil {
		return false, fmt.Errorf("title=%q: running confirm prompt: %w", title, err)
	}

	return result, nil
}

// defaultPromptRunner is the package-level prompt implementation.
// Override this in tests to avoid interactive terminal requirements.
//
//nolint:gochecknoglobals // package-level test hook
var defaultPromptRunner PromptRunner = &huhPromptRunner{}

// PromptString prompts the user for a string value.
func PromptString(title, defaultValue string) (string, error) {
	result, err := defaultPromptRunner.PromptString(title, defaultValue)
	if err != nil {
		return "", fmt.Errorf("title=%q, defaultValue=%q: prompting for string: %w", title, defaultValue, err)
	}

	return result, nil
}

// PromptSelect prompts the user to select from a list of options.
func PromptSelect(title string, options []string) (string, error) {
	result, err := defaultPromptRunner.PromptSelect(title, options)
	if err != nil {
		return "", fmt.Errorf("title=%q, options=%v: prompting for selection: %w", title, options, err)
	}

	return result, nil
}

// PromptConfirm prompts the user for a yes/no confirmation.
func PromptConfirm(title string) (bool, error) {
	result, err := defaultPromptRunner.PromptConfirm(title)
	if err != nil {
		return false, fmt.Errorf("title=%q: prompting for confirmation: %w", title, err)
	}

	return result, nil
}

// promptMissingCommandFlags interactively prompts for any command-level flags
// that have a prompt tag and were not explicitly provided via CLI arguments or
// environment variables. Prompted values are set on the cobra flag set so that
// subsequent flag parsing picks them up normally.
func promptMissingCommandFlags(c *cobra.Command, registry *FlagRegistry) error {
	if registry == nil {
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

		if err := c.Flags().Set(tag.Name, value); err != nil {
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
