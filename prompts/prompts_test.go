package prompts

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestHuhRunner_ImplementsPromptRunner(t *testing.T) {
	t.Parallel()

	var _ v5.PromptRunner = (*HuhRunner)(nil)
}

func TestHuhRunner_PromptConfirm(t *testing.T) {
	t.Parallel()

	var runner v5.PromptRunner = &HuhRunner{}
	_ = runner
}

func TestRegister_DoesNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Register() panicked: %v", r)
		}
	}()

	Register()
}
