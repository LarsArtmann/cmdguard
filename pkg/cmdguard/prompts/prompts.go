// Package prompts provides interactive terminal prompting via huh/v2.
// It is an optional module — import it only when you need interactive
// prompts, to avoid pulling in bubbles, bubbletea, and the full TUI framework.
//
// Usage:
//
//	import (
//	    "github.com/larsartmann/cmdguard/prompts"
//	    v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
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

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

// HuhRunner implements v2.PromptRunner using the huh/v2 TUI library.
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
	v2.SetPromptRunner(&HuhRunner{})
}
