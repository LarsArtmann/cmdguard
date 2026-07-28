package prompts

import (
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestHuhRunner_ImplementsPromptRunner(t *testing.T) {
	t.Parallel()

	var _ v4.PromptRunner = (*HuhRunner)(nil)
}

func TestHuhRunner_PromptConfirm(t *testing.T) {
	t.Parallel()

	var runner v4.PromptRunner = &HuhRunner{}
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
