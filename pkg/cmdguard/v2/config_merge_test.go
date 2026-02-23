package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeConfigs(t *testing.T) {
	t.Run("empty configs", func(t *testing.T) {
		result := MergeConfigs[int]()
		assert.Nil(t, result)
	})

	t.Run("single config", func(t *testing.T) {
		type TestConfig struct {
			Name string
		}
		cfg := &TestConfig{Name: "test"}
		result := MergeConfigs(cfg)
		require.NotNil(t, result)
		assert.Equal(t, "test", result.Name)
	})

	t.Run("merge two configs", func(t *testing.T) {
		type TestConfig struct {
			Name  string
			Count int
		}

		base := &TestConfig{Name: "base", Count: 10}
		override := &TestConfig{Name: "override", Count: 0} // Zero value won't override

		result := MergeConfigs(base, override)
		require.NotNil(t, result)
		assert.Equal(t, "override", result.Name) // Overridden
		assert.Equal(t, 10, result.Count)        // Not overridden (zero value)
	})

	t.Run("nil base config", func(t *testing.T) {
		type TestConfig struct {
			Name string
		}
		override := &TestConfig{Name: "override"}
		result := MergeConfigs[TestConfig](nil, override)
		require.NotNil(t, result)
		assert.Equal(t, "override", result.Name)
	})

	t.Run("nil override config", func(t *testing.T) {
		type TestConfig struct {
			Name string
		}
		base := &TestConfig{Name: "base"}
		result := MergeConfigs(base, nil)
		require.NotNil(t, result)
		assert.Equal(t, "base", result.Name)
	})

	t.Run("nested struct merge", func(t *testing.T) {
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
		require.NotNil(t, result)
		assert.Equal(t, "override-inner", result.Inner.Value)
		assert.Equal(t, "base", result.Name) // Not overridden (empty)
	})

	t.Run("multiple configs", func(t *testing.T) {
		type TestConfig struct {
			A string
			B string
			C string
		}

		first := &TestConfig{A: "a1"}
		second := &TestConfig{B: "b2"}
		third := &TestConfig{C: "c3"}

		result := MergeConfigs(first, second, third)
		require.NotNil(t, result)
		assert.Equal(t, "a1", result.A)
		assert.Equal(t, "b2", result.B)
		assert.Equal(t, "c3", result.C)
	})
}
