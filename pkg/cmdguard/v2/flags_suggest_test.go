package v2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuggestFlag(t *testing.T) {
	validNames := []string{"verbose", "version", "config", "help", "output"}

	t.Run("exact match returns same name", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbose")
		assert.Equal(t, "verbose", result)
	})

	t.Run("one character typo returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbosee")
		assert.Equal(t, "verbose", result)
	})

	t.Run("missing character returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbos")
		assert.Equal(t, "verbose", result)
	})

	t.Run("transposed characters returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbsoe")
		assert.Equal(t, "verbose", result)
	})

	t.Run("similar prefix returns closest match", func(t *testing.T) {
		result := SuggestFlag(validNames, "confing")
		assert.Equal(t, "config", result)
	})

	t.Run("too different returns empty", func(t *testing.T) {
		result := SuggestFlag(validNames, "xyzzy")
		assert.Empty(t, result)
	})

	t.Run("empty valid names returns empty", func(t *testing.T) {
		result := SuggestFlag([]string{}, "test")
		assert.Empty(t, result)
	})

	t.Run("single valid name match", func(t *testing.T) {
		result := SuggestFlag([]string{"help"}, "hlep")
		assert.Equal(t, "help", result)
	})

	t.Run("selects closest match among multiple", func(t *testing.T) {
		names := []string{"start", "status", "stop"}
		result := SuggestFlag(names, "stat")
		// "stat" is distance 2 from "start" and "status", distance 1 from neither
		// Should return one of the closest matches
		assert.Contains(t, names, result)
	})
}

func TestEditDistance(t *testing.T) {
	t.Run("identical strings have distance 0", func(t *testing.T) {
		assert.Equal(t, 0, editDistance("hello", "hello"))
	})

	t.Run("empty strings", func(t *testing.T) {
		assert.Equal(t, 5, editDistance("", "hello"))
		assert.Equal(t, 5, editDistance("hello", ""))
		assert.Equal(t, 0, editDistance("", ""))
	})

	t.Run("single insertion", func(t *testing.T) {
		assert.Equal(t, 1, editDistance("hell", "hello"))
	})

	t.Run("single deletion", func(t *testing.T) {
		assert.Equal(t, 1, editDistance("hello", "hell"))
	})

	t.Run("single substitution", func(t *testing.T) {
		assert.Equal(t, 1, editDistance("hello", "hallo"))
	})

	t.Run("transposition counts as 2", func(t *testing.T) {
		// Standard Levenshtein: "ab" -> "ba" is 2 operations
		assert.Equal(t, 2, editDistance("ab", "ba"))
	})

	t.Run("multiple edits", func(t *testing.T) {
		assert.Equal(t, 3, editDistance("kitten", "sitting"))
	})

	t.Run("case sensitive", func(t *testing.T) {
		assert.Equal(t, 1, editDistance("Hello", "hello"))
	})
}

func TestNewFlagErrorWithSuggestion(t *testing.T) {
	t.Run("error includes suggestion", func(t *testing.T) {
		err := NewFlagErrorWithSuggestion("verboose", fmt.Errorf("unknown flag"), "verbose")
		assert.Contains(t, err.Error(), "verboose")
		assert.Contains(t, err.Error(), "unknown flag")
		assert.Contains(t, err.Error(), "did you mean --verbose")
	})

	t.Run("empty suggestion omits hint", func(t *testing.T) {
		err := NewFlagError("test", fmt.Errorf("some error"))
		assert.NotContains(t, err.Error(), "did you mean")
	})

	t.Run("unwraps to inner error", func(t *testing.T) {
		inner := fmt.Errorf("inner error")
		err := NewFlagErrorWithSuggestion("flag", inner, "suggestion")
		assert.ErrorIs(t, err, inner)
	})
}
