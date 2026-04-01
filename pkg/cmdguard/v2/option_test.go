package v2

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSome(t *testing.T) {
	t.Parallel()
	opt := Some(42)
	if !opt.IsSome() {
		t.Error("Some(42).IsSome() should be true")
	}
	if opt.IsNone() {
		t.Error("Some(42).IsNone() should be false")
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()
	if opt.IsSome() {
		t.Error("None().IsSome() should be false")
	}
	if !opt.IsNone() {
		t.Error("None().IsNone() should be true")
	}
}

func TestOption_Get(t *testing.T) {
	// Test Some
	someOpt := Some("hello")
	val, ok := someOpt.Get()
	if !ok {
		t.Error("Some().Get() should return ok=true")
	}
	if val != "hello" {
		t.Errorf("Some().Get() should return value, got %v", val)
	}

	// Test None
	noneOpt := None[string]()
	val, ok = noneOpt.Get()
	if ok {
		t.Error("None().Get() should return ok=false")
	}
	if val != "" {
		t.Errorf("None().Get() should return zero value, got %v", val)
	}
}

func TestOption_Unwrap(t *testing.T) {
	opt := Some(42)
	if opt.Unwrap() != 42 {
		t.Error("Some(42).Unwrap() should return 42")
	}

	// Test panic on None
	defer func() {
		if r := recover(); r == nil {
			t.Error("None().Unwrap() should panic")
		}
	}()
	None[int]().Unwrap()
}

func TestOption_UnwrapOr(t *testing.T) {
	someOpt := Some(42)
	if someOpt.UnwrapOr(0) != 42 {
		t.Error("Some(42).UnwrapOr(0) should return 42")
	}

	noneOpt := None[int]()
	if noneOpt.UnwrapOr(100) != 100 {
		t.Error("None().UnwrapOr(100) should return 100")
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	someOpt := Some(42)
	result := someOpt.UnwrapOrElse(func() int { return 100 })
	if result != 42 {
		t.Error("Some(42).UnwrapOrElse should return 42")
	}

	noneOpt := None[int]()
	result = noneOpt.UnwrapOrElse(func() int { return 100 })
	if result != 100 {
		t.Error("None().UnwrapOrElse should return computed value")
	}
}

func TestOption_UnwrapOrError(t *testing.T) {
	someOpt := Some(42)
	val, err := someOpt.UnwrapOrError(errors.New("not found"))
	if err != nil {
		t.Error("Some().UnwrapOrError should not return error")
	}
	if val != 42 {
		t.Error("Some(42).UnwrapOrError should return 42")
	}

	noneOpt := None[int]()
	_, err = noneOpt.UnwrapOrError(errors.New("not found"))
	if err == nil {
		t.Error("None().UnwrapOrError should return error")
	}
}

func TestOption_Expect(t *testing.T) {
	opt := Some(42)
	if opt.Expect("should not panic") != 42 {
		t.Error("Expect should return value")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("None().Expect should panic")
		}
	}()
	None[int]().Expect("custom error message")
}

func TestOption_Map(t *testing.T) {
	someOpt := Some(21)
	result := someOpt.Map(func(v int) int { return v * 2 })
	if result.UnwrapOr(0) != 42 {
		t.Error("Map should apply function to Some value")
	}

	noneOpt := None[int]()
	mapped := noneOpt.Map(func(v int) int { return v * 2 })
	if mapped.IsSome() {
		t.Error("Map on None should return None")
	}
}

func TestOption_MapOr(t *testing.T) {
	someOpt := Some(21)
	result := someOpt.MapOr(0, func(v int) int { return v * 2 })
	if result != 42 {
		t.Error("MapOr should apply function to Some value")
	}

	noneOpt := None[int]()
	result = noneOpt.MapOr(100, func(v int) int { return v * 2 })
	if result != 100 {
		t.Error("MapOr should return default for None")
	}
}

func TestOption_And(t *testing.T) {
	x := Some(42)
	y := Some(100)
	result := x.And(y)
	if result.UnwrapOr(0) != 100 {
		t.Error("Some.And(Some) should return second Some")
	}

	x = None[int]()
	result = x.And(Some(100))
	if result.IsSome() {
		t.Error("None.And(Some) should return None")
	}
}

func TestOption_Or(t *testing.T) {
	x := Some(42)
	y := Some(100)
	result := x.Or(y)
	if result.UnwrapOr(0) != 42 {
		t.Error("Some.Or(Some) should return first Some")
	}

	x = None[int]()
	result = x.Or(Some(100))
	if result.UnwrapOr(0) != 100 {
		t.Error("None.Or(Some) should return second Some")
	}
}

func TestOption_Filter(t *testing.T) {
	someOpt := Some(42)
	result := someOpt.Filter(func(v int) bool { return v > 40 })
	if result.IsNone() {
		t.Error("Filter with true predicate should return Some")
	}

	result = someOpt.Filter(func(v int) bool { return v > 50 })
	if result.IsSome() {
		t.Error("Filter with false predicate should return None")
	}

	noneOpt := None[int]()
	result = noneOpt.Filter(func(v int) bool { return true })
	if result.IsSome() {
		t.Error("Filter on None should return None")
	}
}

func TestOption_IfSome(t *testing.T) {
	called := false
	Some(42).IfSome(func(v int) { called = true })
	if !called {
		t.Error("IfSome should call function for Some")
	}

	called = false
	None[int]().IfSome(func(v int) { called = true })
	if called {
		t.Error("IfSome should not call function for None")
	}
}

func TestOption_IfNone(t *testing.T) {
	called := false
	Some(42).IfNone(func() { called = true })
	if called {
		t.Error("IfNone should not call function for Some")
	}

	called = false
	None[int]().IfNone(func() { called = true })
	if !called {
		t.Error("IfNone should call function for None")
	}
}

func TestOption_MarshalJSON(t *testing.T) {
	someOpt := Some(42)
	data, err := json.Marshal(someOpt)
	if err != nil {
		t.Fatalf("Marshal Some failed: %v", err)
	}
	if string(data) != "42" {
		t.Errorf("Marshal Some should return value, got %s", string(data))
	}

	noneOpt := None[int]()
	data, err = json.Marshal(noneOpt)
	if err != nil {
		t.Fatalf("Marshal None failed: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Marshal None should return null, got %s", string(data))
	}
}

func TestOption_UnmarshalJSON(t *testing.T) {
	var opt Option[int]
	err := json.Unmarshal([]byte("42"), &opt)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if !opt.IsSome() || opt.Unwrap() != 42 {
		t.Error("Unmarshal 42 should return Some(42)")
	}

	err = json.Unmarshal([]byte("null"), &opt)
	if err != nil {
		t.Fatalf("Unmarshal null failed: %v", err)
	}
	if opt.IsSome() {
		t.Error("Unmarshal null should return None")
	}
}

func TestOption_String(t *testing.T) {
	someOpt := Some(42)
	if someOpt.String() != "Some(42)" {
		t.Errorf("Some.String() should be Some(42), got %s", someOpt.String())
	}

	noneOpt := None[int]()
	if noneOpt.String() != "None" {
		t.Errorf("None.String() should be None, got %s", noneOpt.String())
	}
}

func TestOption_WithStrings(t *testing.T) {
	opt := Some("hello")
	if opt.UnwrapOr("") != "hello" {
		t.Error("Option[string] should work correctly")
	}
}

func TestOption_WithStructs(t *testing.T) {
	type Config struct {
		Name  string
		Value int
	}

	cfg := Config{Name: "test", Value: 42}
	opt := Some(cfg)
	if opt.Unwrap().Name != "test" || opt.Unwrap().Value != 42 {
		t.Error("Option[struct] should work correctly")
	}
}

func TestOption_MarshalText(t *testing.T) {
	// Test Some with string
	someOpt := Some("hello")
	text, err := someOpt.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	if string(text) != "hello" {
		t.Errorf("Some.MarshalText() should return 'hello', got %s", string(text))
	}

	// Test Some with int
	someInt := Some(42)
	text, err = someInt.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	if string(text) != "42" {
		t.Errorf("Some(42).MarshalText() should return '42', got %s", string(text))
	}

	// Test None
	noneOpt := None[string]()
	text, err = noneOpt.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	if string(text) != "" {
		t.Errorf("None.MarshalText() should return '', got %s", string(text))
	}
}
