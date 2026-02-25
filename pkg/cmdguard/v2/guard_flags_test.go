package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneFlags(t *testing.T) {
	type TestFlags struct {
		Name  string
		Count int
	}

	t.Run("clones struct", func(t *testing.T) {
		original := TestFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		require.NotNil(t, cloned)
		// cloned is already TestFlags, no type assertion needed
		assert.Equal(t, original.Name, cloned.Name)
		assert.Equal(t, original.Count, cloned.Count)

		// Verify it's a copy (modifying clone doesn't affect original)
		cloned.Name = "modified"
		assert.Equal(t, "test", original.Name)
	})

	t.Run("clones pointer to struct", func(t *testing.T) {
		original := &TestFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		require.NotNil(t, cloned)
		// cloned is already *TestFlags, no type assertion needed
		assert.Equal(t, original.Name, cloned.Name)
		assert.Equal(t, original.Count, cloned.Count)

		// Verify it's a different pointer
		assert.NotSame(t, original, cloned)
	})

	t.Run("returns nil for nil pointer", func(t *testing.T) {
		var original *TestFlags // nil

		cloned := cloneFlags[*TestFlags](original)
		assert.Nil(t, cloned)
	})

	t.Run("returns as-is for non-struct", func(t *testing.T) {
		original := "string value"
		cloned := cloneFlags(original)
		assert.Equal(t, original, cloned)
	})
}

func TestFlagTypeConstraint(t *testing.T) {
	type testFlags struct {
		Name  string
		Count int
	}

	t.Run("accepts NoFlags (struct{})", func(t *testing.T) {
		err := FlagTypeConstraint[NoFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		err := FlagTypeConstraint[*testFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts empty struct", func(t *testing.T) {
		type EmptyFlags struct{}

		err := FlagTypeConstraint[EmptyFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts struct with fields", func(t *testing.T) {
		err := FlagTypeConstraint[testFlags]()
		assert.NoError(t, err)
	})

	t.Run("rejects pointer to non-struct", func(t *testing.T) {
		err := FlagTypeConstraint[*string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "*string")
	})

	t.Run("rejects int", func(t *testing.T) {
		err := FlagTypeConstraint[int]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "int")
	})

	t.Run("rejects string", func(t *testing.T) {
		err := FlagTypeConstraint[string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "string")
	})

	t.Run("rejects slice", func(t *testing.T) {
		err := FlagTypeConstraint[[]string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "[]string")
	})

	t.Run("rejects map", func(t *testing.T) {
		err := FlagTypeConstraint[map[string]string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "map[string]string")
	})
}
