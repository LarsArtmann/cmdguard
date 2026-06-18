package v2

import (
	"errors"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestResult_Ok(t *testing.T) {
	t.Parallel()

	r := Ok(42)
	testutil.AssertBoolTrue(t, r.IsOk(), "expected Ok to be IsOk")
	testutil.AssertBoolFalse(t, r.IsErr(), "expected Ok to not be IsErr")

	v, ok := r.Value()
	testutil.AssertBoolTrue(t, ok, "expected Value ok=true")
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestResult_Err(t *testing.T) {
	t.Parallel()

	errSentinel := errors.New("boom")
	r := Err[int](errSentinel)
	testutil.AssertBoolTrue(t, r.IsErr(), "expected Err to be IsErr")
	testutil.AssertBoolFalse(t, r.IsOk(), "expected Err to not be IsOk")

	if !errors.Is(r.ErrValue(), errSentinel) {
		t.Errorf("expected ErrValue to match sentinel")
	}
}

func TestResult_UnwrapOr(t *testing.T) {
	t.Parallel()

	if Ok(42).UnwrapOr(0) != 42 {
		t.Error("expected UnwrapOr to return value when Ok")
	}

	if Err[int](errors.New("x")).UnwrapOr(99) != 99 {
		t.Error("expected UnwrapOr to return fallback when Err")
	}
}

func TestResult_UnwrapOrElse(t *testing.T) {
	t.Parallel()

	if Ok(42).UnwrapOrElse(func(error) int { return 0 }) != 42 {
		t.Error("expected UnwrapOrElse to return value when Ok")
	}

	if Err[int](errors.New("x")).UnwrapOrElse(func(error) int { return 99 }) != 99 {
		t.Error("expected UnwrapOrElse to return fn result when Err")
	}
}

func TestResult_Map(t *testing.T) {
	t.Parallel()

	doubled := Ok(21).Map(func(v int) int { return v * 2 })
	if v, _ := doubled.Value(); v != 42 {
		t.Errorf("expected Map to double to 42, got %d", v)
	}

	errSentinel := errors.New("boom")
	mapped := Err[int](errSentinel).Map(func(v int) int { return v * 2 })
	testutil.AssertErrorIs(t, mapped.ErrValue(), errSentinel)
}

func TestResult_AndThen(t *testing.T) {
	t.Parallel()

	exclaim := Ok("hello").AndThen(func(s string) Result[string] {
		return Ok(s + "!")
	})

	if v, _ := exclaim.Value(); v != "hello!" {
		t.Errorf("expected AndThen to produce 'hello!', got %q", v)
	}

	errSentinel := errors.New("boom")
	chained := Err[string](errSentinel).AndThen(func(s string) Result[string] {
		return Ok(s + "!")
	})

	testutil.AssertErrorIs(t, chained.ErrValue(), errSentinel)
}

func TestValidated_Valid(t *testing.T) {
	t.Parallel()

	v := Valid("hello")
	testutil.AssertBoolTrue(t, v.IsValid(), "expected Valid to be IsValid")
	if val, ok := v.Value(); !ok || val != "hello" {
		t.Errorf("expected Value 'hello', got %q ok=%v", val, ok)
	}
}

func TestValidated_Invalid(t *testing.T) {
	t.Parallel()

	e1 := errors.New("too short")
	e2 := errors.New("missing uppercase")
	v := Invalid("ab", e1, e2)
	testutil.AssertBoolFalse(t, v.IsValid(), "expected Invalid to not be IsValid")

	errs := v.Errors()
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestValidated_AddErr(t *testing.T) {
	t.Parallel()

	v := Valid(42)
	testutil.AssertBoolTrue(t, v.IsValid(), "expected fresh Validated to be valid")

	v.AddErr(errors.New("oops"))
	testutil.AssertBoolFalse(t, v.IsValid(), "expected invalid after AddErr")
}

func TestValidated_Combine(t *testing.T) {
	t.Parallel()

	v1 := Invalid("x", errors.New("e1"))
	v2 := Invalid("x", errors.New("e2"), errors.New("e3"))

	v1.Combine(v2)
	if len(v1.Errors()) != 3 {
		t.Errorf("expected 3 combined errors, got %d", len(v1.Errors()))
	}
}

func TestValidated_ToResult(t *testing.T) {
	t.Parallel()

	okResult := Valid(42).ToResult()
	testutil.AssertBoolTrue(t, okResult.IsOk(), "expected ToResult to be Ok for valid")

	invalidResult := Invalid(0, errors.New("e1"), errors.New("e2")).ToResult()
	testutil.AssertBoolTrue(t, invalidResult.IsErr(), "expected ToResult to be Err for invalid")
}
