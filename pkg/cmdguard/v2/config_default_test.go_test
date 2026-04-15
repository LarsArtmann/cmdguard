package v2

import (
	"reflect"
	"testing"
	"time"
)

// newFlagTag creates a FlagTag with the given type and default value.
func newFlagTag[T any](defaultVal string) FlagTag {
	return FlagTag{
		Type:    reflect.TypeFor[T](),
		Default: defaultVal,
	}
}

func TestFlagTag_DefaultValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tag      FlagTag
		expected any
	}{
		{name: "string default", tag: newFlagTag[string]("hello"), expected: "hello"},
		{name: "bool true", tag: newFlagTag[bool]("true"), expected: true},
		{name: "bool false", tag: newFlagTag[bool]("false"), expected: false},
		{name: "int default", tag: newFlagTag[int]("42"), expected: 42},
		{name: "uint default", tag: newFlagTag[uint]("42"), expected: uint(42)},
		{name: "uint32 default", tag: newFlagTag[uint32]("4294967295"), expected: uint(4294967295)},
		{name: "float64 default", tag: newFlagTag[float64]("3.14"), expected: 3.14},
		{name: "empty default returns zero", tag: newFlagTag[int](""), expected: 0},
		{
			name:     "slice default",
			tag:      newFlagTag[[]string]("a,b,c"),
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Duration default",
			tag:      newFlagTag[Duration]("5m"),
			expected: FromDuration(5 * time.Minute),
		},
		{name: "LogLevel default", tag: newFlagTag[LogLevel]("info"), expected: "info"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := tc.tag.DefaultValue()
			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestConfig_DefaultConfig(t *testing.T) {
	t.Parallel()
	t.Run("default values", func(t *testing.T) {
		t.Parallel()

		cfg := Config{}

		tags, err := ParseFlagTags(cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(tags) < 4 {
			t.Fatalf("expected at least 4 tags, got %d", len(tags))
		}

		// Find LogLevel tag
		var logLevelTag *FlagTag

		for _, tag := range tags {
			if tag.Name == "log-level" {
				logLevelTag = &tag

				break
			}
		}

		if logLevelTag == nil {
			t.Fatal("expected logLevelTag to not be nil")
		}

		if logLevelTag.Default != "info" {
			t.Errorf("expected Default 'info', got %q", logLevelTag.Default)
		}

		if logLevelTag.Short != "l" {
			t.Errorf("expected Short 'l', got %q", logLevelTag.Short)
		}
	})
}
