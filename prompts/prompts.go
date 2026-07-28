// Package prompts provides interactive terminal prompting via huh/v2.
// It is an optional module — import it only when you need interactive
// prompts, to avoid pulling in bubbles, bubbletea, and the full TUI framework.
//
// Usage:
//
//	import (
//	    "github.com/larsartmann/cmdguard/prompts"
//	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
//	)
//
//	func main() {
//	    prompts.Register() // wire huh-based prompting into cmdguard
//	    // ... build CLI as normal
//	}
package prompts

import (
	"fmt"

	"charm.land/huh/v2"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// HuhRunner implements v4.PromptRunner using the huh/v2 TUI library.
// TODO(v5): "Runner" suffix is generic; consider HuhPrompter (see naming-review 2026-07-18).
type HuhRunner struct{}

func (h *HuhRunner) PromptString(title, defaultValue string) (string, error) {
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

func (h *HuhRunner) PromptSelect(title string, options []string) (string, error) {
	var result string

	opts := make([]huh.Option[string], 0, len(options))
	for _, o := range options {
		opts = append(opts, huh.NewOption(o, o))
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

func (h *HuhRunner) PromptConfirm(title string) (bool, error) {
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

// Register wires huh-based prompting into cmdguard. Call this once at startup
// to enable interactive prompts for missing flags with the `prompt:"..."` tag.
func Register() {
	v4.SetPromptRunner(&HuhRunner{})
}
