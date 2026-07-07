package v3

import (
	"slices"
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestParseFlagTags(t *testing.T) {
	t.Parallel()
	t.Run("valid struct", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name    string `default:"test" flag:"name"    help:"The name"       short:"n"`
			Count   int    `default:"10"   flag:"count"   help:"The count"`
			Enabled bool   `default:"true" flag:"enabled" help:"Enable feature" short:"e"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		testutil.AssertFieldLen(t, tags, 3, "tags")

		// Check first field
		if tags[0].Field != "Name" {
			t.Errorf("expected Field 'Name', got %q", tags[0].Field)
		}

		if tags[0].Name != "name" {
			t.Errorf("expected Name 'name', got %q", tags[0].Name)
		}

		if tags[0].Short != "n" {
			t.Errorf("expected Short 'n', got %q", tags[0].Short)
		}

		if tags[0].Default != "test" {
			t.Errorf("expected Default 'test', got %q", tags[0].Default)
		}

		if tags[0].Help != "The name" {
			t.Errorf("expected Help 'The name', got %q", tags[0].Help)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Field string `flag:"field"`
		}

		tags, err := ParseFlagTags(&TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		testutil.AssertFieldLen(t, tags, 1, "tags")

		if tags[0].Name != "field" {
			t.Errorf("expected Name 'field', got %q", tags[0].Name)
		}
	})

	t.Run("skips fields without flag tag", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Tagged   string `flag:"tagged"`
			Untagged string
			Ignored  string `flag:"-"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		testutil.AssertFieldLen(t, tags, 1, "tags")

		if tags[0].Field != "Tagged" {
			t.Errorf("expected Field 'Tagged', got %q", tags[0].Field)
		}
	})

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()

		tags, err := ParseFlagTags(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "must not be nil")

		if tags != nil {
			t.Errorf("expected nil tags, got %v", tags)
		}
	})

	t.Run("non-struct config", func(t *testing.T) {
		t.Parallel()

		tags, err := ParseFlagTags("not a struct")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "expected struct")

		if tags != nil {
			t.Errorf("expected nil tags, got %v", tags)
		}
	})

	t.Run("with values tag", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn,error"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		testutil.AssertFieldLen(t, tags, 1, "tags")

		expected := []string{"debug", "info", "warn", "error"}
		if !slices.Equal(tags[0].Values, expected) {
			t.Errorf("expected Values %v, got %v", expected, tags[0].Values)
		}
	})

	t.Run("embedded Config", func(t *testing.T) {
		t.Parallel()

		type AppConfig struct {
			Config

			AppName string `default:"myapp" flag:"app-name"`
		}

		tags, err := ParseFlagTags(AppConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		// Config has 4 fields + AppName = 5
		if len(tags) < 1 {
			t.Errorf("expected at least 1 tag, got %d", len(tags))
		}
	})

	t.Run("invalid required tag returns error", func(t *testing.T) {
		t.Parallel()

		type BadConfig struct {
			Name string `flag:"name" required:"tru"`
		}

		_, err := ParseFlagTags(BadConfig{})
		if err == nil {
			t.Fatal("expected error for invalid required tag, got nil")
		}

		assertErrorContains(t, err, "Name", "required", "tru")
	})
}
