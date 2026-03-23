package v2

import (
	"testing"
)

func TestMergeConfigs(t *testing.T) {
	t.Parallel()
	t.Run("empty configs", func(t *testing.T) {
		t.Parallel()
		result := MergeConfigs[int]()
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("single config", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			Name string
		}

		cfg := &TestConfig{Name: "test"}

		result := MergeConfigs(cfg)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.Name != "test" {
			t.Errorf("expected name 'test', got %q", result.Name)
		}
	})

	t.Run("merge two configs", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			Name  string
			Count int
		}

		base := &TestConfig{Name: "base", Count: 10}
		override := &TestConfig{Name: "override", Count: 0} // Zero value won't override

		result := MergeConfigs(base, override)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.Name != "override" { // Overridden
			t.Errorf("expected name 'override', got %q", result.Name)
		}

		if result.Count != 10 { // Not overridden (zero value)
			t.Errorf("expected count 10, got %d", result.Count)
		}
	})

	t.Run("nil base config", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			Name string
		}

		override := &TestConfig{Name: "override"}

		result := MergeConfigs(nil, override)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.Name != "override" {
			t.Errorf("expected name 'override', got %q", result.Name)
		}
	})

	t.Run("nil override config", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			Name string
		}

		base := &TestConfig{Name: "base"}

		result := MergeConfigs(base, nil)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.Name != "base" {
			t.Errorf("expected name 'base', got %q", result.Name)
		}
	})

	t.Run("nested struct merge", func(t *testing.T) {
		t.Parallel()
		type Inner struct {
			Value string
		}

		type Outer struct {
			Inner Inner
			Name  string
		}

		base := &Outer{Inner: Inner{Value: "base-inner"}, Name: "base"}
		override := &Outer{Inner: Inner{Value: "override-inner"}, Name: ""}

		result := MergeConfigs(base, override)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.Inner.Value != "override-inner" {
			t.Errorf("expected inner value 'override-inner', got %q", result.Inner.Value)
		}

		if result.Name != "base" { // Not overridden (empty)
			t.Errorf("expected name 'base', got %q", result.Name)
		}
	})

	t.Run("multiple configs", func(t *testing.T) {
		t.Parallel()
		type TestConfig struct {
			A string
			B string
			C string
		}

		first := &TestConfig{A: "a1"}
		second := &TestConfig{B: "b2"}
		third := &TestConfig{C: "c3"}

		result := MergeConfigs(first, second, third)
		if result == nil {
			t.Fatal("expected result to not be nil")
		}

		if result.A != "a1" {
			t.Errorf("expected A 'a1', got %q", result.A)
		}

		if result.B != "b2" {
			t.Errorf("expected B 'b2', got %q", result.B)
		}

		if result.C != "c3" {
			t.Errorf("expected C 'c3', got %q", result.C)
		}
	})
}
