package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestOk(t *testing.T) {
	t.Parallel()

	r := Ok(42)

	if !r.IsOk() {
		t.Error("expected IsOk() to be true")
	}

	if r.IsErr() {
		t.Error("expected IsErr() to be false")
	}

	val, err := r.Get()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestErr(t *testing.T) {
	t.Parallel()

	r := Err[int](errors.New("failed"))

	if r.IsOk() {
		t.Error("expected IsOk() to be false")
	}

	if !r.IsErr() {
		t.Error("expected IsErr() to be true")
	}

	val, err := r.Get()
	if err == nil {
		t.Error("expected error")
	}

	if val != 0 {
		t.Errorf("expected zero value, got %d", val)
	}
}

func TestResult_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		r := Ok("hello")
		if r.Unwrap() != "hello" {
			t.Error("expected hello")
		}
	})

	t.Run("err panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
		}()

		Err[int](errors.New("boom")).Unwrap()
	})
}

func TestResult_UnwrapOr(t *testing.T) {
	t.Parallel()

	t.Run("ok returns value", func(t *testing.T) {
		t.Parallel()

		r := Ok(42)
		if r.UnwrapOr(-1) != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("err returns default", func(t *testing.T) {
		t.Parallel()

		r := Err[int](errors.New("fail"))
		if r.UnwrapOr(-1) != -1 {
			t.Error("expected -1")
		}
	})
}

func TestResult_UnwrapOrElse(t *testing.T) {
	t.Parallel()

	t.Run("ok returns value", func(t *testing.T) {
		t.Parallel()

		r := Ok(10)
		if r.UnwrapOrElse(func(_ error) int { return 99 }) != 10 {
			t.Error("expected 10")
		}
	})

	t.Run("err calls function", func(t *testing.T) {
		t.Parallel()

		r := Err[int](errors.New("fail"))
		if r.UnwrapOrElse(func(err error) int {
			return len(err.Error())
		}) != 4 {
			t.Error("expected 4")
		}
	})
}

func TestResult_UnwrapErr(t *testing.T) {
	t.Parallel()

	t.Run("err returns error", func(t *testing.T) {
		t.Parallel()

		e := errors.New("fail")

		r := Err[int](e)
		if r.UnwrapErr() != e {
			t.Error("expected same error")
		}
	})

	t.Run("ok panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
		}()

		Ok(42).UnwrapErr()
	})
}

func TestResult_Expect(t *testing.T) {
	t.Parallel()

	t.Run("ok returns value", func(t *testing.T) {
		t.Parallel()

		r := Ok(42)
		if r.Expect("should be ok") != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("err panics with message", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}

			s := fmt.Sprintf("%v", r)
			if s == "" {
				t.Error("expected panic message")
			}
		}()

		Err[int](errors.New("boom")).Expect("something went wrong")
	})
}

func TestResult_ExpectErr(t *testing.T) {
	t.Parallel()

	t.Run("err returns error", func(t *testing.T) {
		t.Parallel()

		e := errors.New("fail")

		r := Err[int](e)
		if r.ExpectErr("should be err") != e {
			t.Error("expected same error")
		}
	})

	t.Run("ok panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic")
			}
		}()

		Ok(42).ExpectErr("should not be ok")
	})
}

func TestResult_Map(t *testing.T) {
	t.Parallel()

	t.Run("ok maps value", func(t *testing.T) {
		t.Parallel()

		r := Ok(10).Map(func(v int) int { return v * 2 })
		if r.Unwrap() != 20 {
			t.Error("expected 20")
		}
	})

	t.Run("err stays err", func(t *testing.T) {
		t.Parallel()

		e := errors.New("fail")

		r := Err[int](e).Map(func(v int) int { return v * 2 })
		if r.IsOk() {
			t.Error("expected Err")
		}

		if r.UnwrapErr() != e {
			t.Error("expected same error")
		}
	})
}

func TestResult_MapErr(t *testing.T) {
	t.Parallel()

	t.Run("ok stays ok", func(t *testing.T) {
		t.Parallel()

		r := Ok(10).MapErr(func(err error) error { return fmt.Errorf("wrapped: %w", err) })
		if !r.IsOk() {
			t.Error("expected Ok")
		}
	})

	t.Run("err maps error", func(t *testing.T) {
		t.Parallel()

		wrapped := fmt.Errorf("wrapped: %w", errors.New("fail"))

		r := Err[int](errors.New("fail")).MapErr(func(_ error) error {
			return wrapped
		})
		if r.IsOk() {
			t.Error("expected Err")
		}

		if r.UnwrapErr() != wrapped {
			t.Error("expected wrapped error")
		}
	})
}

func TestResult_MapOr(t *testing.T) {
	t.Parallel()

	t.Run("ok maps value", func(t *testing.T) {
		t.Parallel()

		v := Ok(10).MapOr(-1, func(v int) int { return v * 2 })
		if v != 20 {
			t.Error("expected 20")
		}
	})

	t.Run("err returns default", func(t *testing.T) {
		t.Parallel()

		v := Err[int](errors.New("fail")).MapOr(-1, func(v int) int { return v * 2 })
		if v != -1 {
			t.Error("expected -1")
		}
	})
}

func TestResult_And(t *testing.T) {
	t.Parallel()

	t.Run("ok and ok", func(t *testing.T) {
		t.Parallel()

		r := Ok(1).And(Ok(2))
		if r.Unwrap() != 2 {
			t.Error("expected 2")
		}
	})

	t.Run("ok and err", func(t *testing.T) {
		t.Parallel()

		r := Ok(1).And(Err[int](errors.New("fail")))
		if r.IsOk() {
			t.Error("expected Err")
		}
	})

	t.Run("err and ok", func(t *testing.T) {
		t.Parallel()

		e := errors.New("first")

		r := Err[int](e).And(Ok(2))
		if r.UnwrapErr() != e {
			t.Error("expected first error")
		}
	})
}

func TestResult_Or(t *testing.T) {
	t.Parallel()

	t.Run("ok or other", func(t *testing.T) {
		t.Parallel()

		r := Ok(1).Or(Ok(2))
		if r.Unwrap() != 1 {
			t.Error("expected 1")
		}
	})

	t.Run("err or ok", func(t *testing.T) {
		t.Parallel()

		r := Err[int](errors.New("fail")).Or(Ok(2))
		if r.Unwrap() != 2 {
			t.Error("expected 2")
		}
	})

	t.Run("err or err", func(t *testing.T) {
		t.Parallel()

		r := Err[int](errors.New("first")).Or(Err[int](errors.New("second")))
		if r.IsOk() {
			t.Error("expected Err")
		}
	})
}

func TestResult_IfOk(t *testing.T) {
	t.Parallel()

	t.Run("ok calls function", func(t *testing.T) {
		t.Parallel()

		called := false

		Ok(42).IfOk(func(v int) { called = true })

		if !called {
			t.Error("expected function to be called")
		}
	})

	t.Run("err does not call function", func(t *testing.T) {
		t.Parallel()

		called := false

		Err[int](errors.New("fail")).IfOk(func(_ int) { called = true })

		if called {
			t.Error("expected function not to be called")
		}
	})
}

func TestResult_IfErr(t *testing.T) {
	t.Parallel()

	t.Run("err calls function", func(t *testing.T) {
		t.Parallel()

		called := false

		Err[int](errors.New("fail")).IfErr(func(_ error) { called = true })

		if !called {
			t.Error("expected function to be called")
		}
	})

	t.Run("ok does not call function", func(t *testing.T) {
		t.Parallel()

		called := false

		Ok(42).IfErr(func(_ error) { called = true })

		if called {
			t.Error("expected function not to be called")
		}
	})
}

func TestResult_ToOption(t *testing.T) {
	t.Parallel()

	t.Run("ok to some", func(t *testing.T) {
		t.Parallel()

		opt := Ok(42).ToOption()
		if !opt.IsSome() {
			t.Error("expected Some")
		}

		if v, ok := opt.Get(); !ok || v != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("err to none", func(t *testing.T) {
		t.Parallel()

		opt := Err[int](errors.New("fail")).ToOption()
		if !opt.IsNone() {
			t.Error("expected None")
		}
	})
}

func TestResult_String(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		if s := Ok(42).String(); s != "Ok(42)" {
			t.Errorf("expected Ok(42), got %s", s)
		}
	})

	t.Run("err", func(t *testing.T) {
		t.Parallel()

		s := Err[int](errors.New("fail")).String()
		if s != "Err(fail)" {
			t.Errorf("expected Err(fail), got %s", s)
		}
	})
}

func TestResult_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(Ok(42))
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != "42" {
			t.Errorf("expected 42, got %s", data)
		}
	})

	t.Run("ok with string", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(Ok("hello"))
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != `"hello"` {
			t.Errorf("expected \"hello\", got %s", data)
		}
	})

	t.Run("err", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(Err[int](errors.New("fail")))
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != `{"error":"fail"}` {
			t.Errorf("expected {\"error\":\"fail\"}, got %s", data)
		}
	})
}

func TestResultFrom(t *testing.T) {
	t.Parallel()

	t.Run("value no error", func(t *testing.T) {
		t.Parallel()

		r := ResultFrom(42, nil)
		if !r.IsOk() {
			t.Error("expected Ok")
		}

		if r.Unwrap() != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		r := ResultFrom(0, errors.New("fail"))
		if !r.IsErr() {
			t.Error("expected Err")
		}
	})
}

func TestResult_ToPair(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		val, err := Ok(42).ToPair()
		if err != nil {
			t.Error("expected no error")
		}

		if val != 42 {
			t.Error("expected 42")
		}
	})

	t.Run("err", func(t *testing.T) {
		t.Parallel()

		val, err := Err[int](errors.New("fail")).ToPair()
		if err == nil {
			t.Error("expected error")
		}

		if val != 0 {
			t.Errorf("expected zero, got %d", val)
		}
	})
}

func TestResult_Chain(t *testing.T) {
	t.Parallel()

	t.Run("successful pipeline", func(t *testing.T) {
		t.Parallel()

		r := Ok(10).
			Map(func(v int) int { return v * 2 }).
			Map(func(v int) int { return v + 5 })

		if r.Unwrap() != 25 {
			t.Error("expected 25")
		}
	})

	t.Run("error short-circuits", func(t *testing.T) {
		t.Parallel()

		called := false
		r := Err[int](errors.New("fail")).
			Map(func(v int) int {
				called = true

				return v
			})

		if called {
			t.Error("expected Map not to be called")
		}

		if r.IsOk() {
			t.Error("expected Err")
		}
	})
}

func TestResult_PracticalUsage(t *testing.T) {
	t.Parallel()

	t.Run("parse port", func(t *testing.T) {
		t.Parallel()

		parse := func(s string) Result[Port] {
			return ResultFrom(ParsePort(s))
		}

		r := parse("8080")
		if !r.IsOk() {
			t.Error("expected Ok")
		}

		if r.Unwrap().Int() != 8080 {
			t.Error("expected 8080")
		}

		r = parse("invalid")
		if !r.IsErr() {
			t.Error("expected Err")
		}
	})

	t.Run("compose with map", func(t *testing.T) {
		t.Parallel()

		double := func(r Result[int]) Result[int] {
			return r.Map(func(v int) int { return v * 2 })
		}

		result := double(Ok(21))
		if result.Unwrap() != 42 {
			t.Error("expected 42")
		}

		result = double(Err[int](errors.New("fail")))
		if result.IsOk() {
			t.Error("expected Err")
		}
	})
}
