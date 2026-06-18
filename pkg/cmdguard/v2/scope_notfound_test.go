package v2

import (
	"errors"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestInvoke_MissingService_ReturnsErrServiceNotFound(t *testing.T) {
	t.Parallel()

	scope := NewScope("test")

	_, err := Invoke[string](scope)
	testutil.AssertExpectedError(t, err)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected errors.Is(err, ErrServiceNotFound), got: %v", err)
	}
}

func TestInvokeNamed_MissingService_ReturnsErrServiceNotFound(t *testing.T) {
	t.Parallel()

	scope := NewScope("test")

	_, err := InvokeNamed[string](scope, "missing")
	testutil.AssertExpectedError(t, err)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected errors.Is(err, ErrServiceNotFound), got: %v", err)
	}
}

func TestInvoke_ExistingService_NoErrServiceNotFound(t *testing.T) {
	t.Parallel()

	scope := NewScope("test")
	mustProvideValue(t, scope, "hello")

	v, err := Invoke[string](scope)
	testutil.AssertNoError(t, err)
	if v != "hello" {
		t.Errorf("expected 'hello', got %q", v)
	}
	if errors.Is(err, ErrServiceNotFound) {
		t.Error("did not expect ErrServiceNotFound for existing service")
	}
}
