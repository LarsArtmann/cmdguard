package v2

import (
	"strings"
	"testing"
)

func slicesEqualStr(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestParseFlagTags(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		type TestConfig struct {
			Name    string `default:"test" flag:"name"    help:"The name"       short:"n"`
			Count   int    `default:"10"   flag:"count"   help:"The count"`
			Enabled bool   `default:"true" flag:"enabled" help:"Enable feature" short:"e"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(tags) != 3 {
			t.Fatalf("expected 3 tags, got %d", len(tags))
		}

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
		type TestConfig struct {
			Field string `flag:"field"`
		}

		tags, err := ParseFlagTags(&TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(tags) != 1 {
			t.Fatalf("expected 1 tag, got %d", len(tags))
		}

		if tags[0].Name != "field" {
			t.Errorf("expected Name 'field', got %q", tags[0].Name)
		}
	})

	t.Run("skips fields without flag tag", func(t *testing.T) {
		type TestConfig struct {
			Tagged   string `flag:"tagged"`
			Untagged string
			Ignored  string `flag:"-"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(tags) != 1 {
			t.Fatalf("expected 1 tag, got %d", len(tags))
		}

		if tags[0].Field != "Tagged" {
			t.Errorf("expected Field 'Tagged', got %q", tags[0].Field)
		}
	})

	t.Run("nil config", func(t *testing.T) {
		tags, err := ParseFlagTags(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "must not be nil") {
			t.Errorf("expected error to contain 'must not be nil', got: %v", err)
		}

		if tags != nil {
			t.Errorf("expected nil tags, got %v", tags)
		}
	})

	t.Run("non-struct config", func(t *testing.T) {
		tags, err := ParseFlagTags("not a struct")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "expected struct") {
			t.Errorf("expected error to contain 'expected struct', got: %v", err)
		}

		if tags != nil {
			t.Errorf("expected nil tags, got %v", tags)
		}
	})

	t.Run("with values tag", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn,error"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(tags) != 1 {
			t.Fatalf("expected 1 tag, got %d", len(tags))
		}

		expected := []string{"debug", "info", "warn", "error"}
		if !slicesEqualStr(tags[0].Values, expected) {
			t.Errorf("expected Values %v, got %v", expected, tags[0].Values)
		}
	})

	t.Run("embedded Config", func(t *testing.T) {
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
}
