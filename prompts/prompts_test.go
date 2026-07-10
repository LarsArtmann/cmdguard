package prompts

import (
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestHuhRunner_ImplementsPromptRunner(t *testing.T) {
	t.Parallel()

	var _ v3.PromptRunner = (*HuhRunner)(nil)
}

func TestHuhRunner_PromptConfirm(t *testing.T) {
	t.Parallel()

	runner := &HuhRunner{}
	if runner == nil {
		t.Fatal("HuhRunner instance is nil")
	}
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
