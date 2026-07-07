package v3

import (
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
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
		override := &TestConfig{Name: "override", Count: 0}

		result := MergeConfigs(base, override)
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.Name, "override", "Name")
		testutil.AssertFieldEq(t, result.Count, 0, "Count")
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
		testutil.AssertFieldEqString(t, result.Name, "", "Name")
	})

	t.Run("multiple configs — later overrides earlier", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			A string
			B string
			C string
		}

		first := &TestConfig{A: "a1", B: "b1", C: "c1"}
		second := &TestConfig{A: "a2", B: "b2", C: ""}
		third := &TestConfig{A: "a3", B: "", C: "c3"}

		result := MergeConfigs(first, second, third)
		testutil.AssertNotNil(t, result)

		testutil.AssertFieldEqString(t, result.A, "a3", "A")
		testutil.AssertFieldEqString(t, result.B, "", "B")
		testutil.AssertFieldEqString(t, result.C, "c3", "C")
	})

	t.Run("does not mutate input configs", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name  string
			Count int
		}

		base := &TestConfig{Name: "base", Count: 10}
		override := &TestConfig{Name: "override", Count: 20}

		result := MergeConfigs(base, override)

		testutil.AssertFieldEqString(t, base.Name, "base", "base.Name should not be mutated")
		testutil.AssertFieldEq(t, base.Count, 10, "base.Count should not be mutated")
		testutil.AssertFieldEqString(t, result.Name, "override", "result.Name")
		testutil.AssertFieldEq(t, result.Count, 20, "result.Count")
	})

	t.Run("returned config is independent", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string
		}

		base := &TestConfig{Name: "original"}

		result := MergeConfigs(base)
		result.Name = "modified"

		testutil.AssertFieldEqString(
			t,
			base.Name,
			"original",
			"base should not be affected by result mutation",
		)
	})

	t.Run("deep copies slices and maps", func(t *testing.T) {
		t.Parallel()

		type ComplexConfig struct {
			Tags   []string
			Labels map[string]string
		}

		base := &ComplexConfig{
			Tags:   []string{"a", "b", "c"},
			Labels: map[string]string{"env": "prod"},
		}

		result := MergeConfigs(base)

		result.Tags[0] = "modified"
		result.Labels["env"] = "dev"

		testutil.AssertFieldEqString(t, base.Tags[0], "a", "base slice should not be mutated")
		testutil.AssertFieldEqString(
			t,
			base.Labels["env"],
			"prod",
			"base map should not be mutated",
		)
	})

	t.Run("deep copies nested pointers", func(t *testing.T) {
		t.Parallel()

		type Inner struct {
			Value string
		}

		type Outer struct {
			Ref *Inner
		}

		base := &Outer{Ref: &Inner{Value: "original"}}

		result := MergeConfigs(base)
		result.Ref.Value = "modified"

		testutil.AssertFieldEqString(
			t,
			base.Ref.Value,
			"original",
			"base pointer target should not be mutated",
		)
	})
}

func TestMergeConfigs_ZeroValueOverride(t *testing.T) {
	t.Parallel()

	t.Run("false overrides true", func(t *testing.T) {
		t.Parallel()

		type Cfg struct {
			Enabled bool
		}

		base := &Cfg{Enabled: true}
		override := &Cfg{Enabled: false}

		result := MergeConfigs(base, override)

		testutil.AssertFieldEq(t, result.Enabled, false, "Enabled")
	})

	t.Run("0 overrides non-zero int", func(t *testing.T) {
		t.Parallel()

		type Cfg struct {
			Port int
		}

		base := &Cfg{Port: 8080}
		override := &Cfg{Port: 0}

		result := MergeConfigs(base, override)

		testutil.AssertFieldEq(t, result.Port, 0, "Port")
	})

	t.Run("empty string overrides non-empty", func(t *testing.T) {
		t.Parallel()

		type Cfg struct {
			Name string
		}

		base := &Cfg{Name: "production"}
		override := &Cfg{Name: ""}

		result := MergeConfigs(base, override)

		testutil.AssertFieldEqString(t, result.Name, "", "Name")
	})

	t.Run("zero float overrides non-zero", func(t *testing.T) {
		t.Parallel()

		type Cfg struct {
			Rate float64
		}

		base := &Cfg{Rate: 3.14}
		override := &Cfg{Rate: 0.0}

		result := MergeConfigs(base, override)

		testutil.AssertFieldEq(t, result.Rate, 0.0, "Rate")
	})
}
