package v2

import (
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestMergeConfigs(t *testing.T) {
	t.Parallel()
	t.Run("empty configs", func(t *testing.T) {
		t.Parallel()

		result := MergeConfigs[int]()
		testutil.AssertNil(t, result)
	})

	t.Run("single config", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string
		}

		cfg := &TestConfig{Name: "test"}

		result := MergeConfigs(cfg)
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Name, "test", "Name")
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
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Name, "override", "Name")
		testutil.AssertFieldEq(t, result.Count, 10, "Count")
	})

	t.Run("nil base config", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string
		}

		override := &TestConfig{Name: "override"}

		result := MergeConfigs(nil, override)
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Name, "override", "Name")
	})

	t.Run("nil override config", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string
		}

		base := &TestConfig{Name: "base"}

		result := MergeConfigs(base, nil)
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Name, "base", "Name")
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
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Inner.Value, "override-inner", "Inner.Value")
		testutil.AssertFieldEqString(t, result.Name, "base", "Name")
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
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.A, "a1", "A")
		testutil.AssertFieldEqString(t, result.B, "b2", "B")
		testutil.AssertFieldEqString(t, result.C, "c3", "C")
	})
}
