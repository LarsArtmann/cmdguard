package prompts

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestHuhPrompter_ImplementsPromptRunner(t *testing.T) {
	t.Parallel()

	var _ v5.PromptRunner = (*HuhPrompter)(nil)
}

func TestHuhPrompter_PromptConfirm(t *testing.T) {
	t.Parallel()

	var runner v5.PromptRunner = &HuhPrompter{}
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
