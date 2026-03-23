package v2

import (
	"reflect"
	"testing"
	"time"
)

func TestFlagTag_DefaultValue(t *testing.T) {
	tests := []struct {
		name     string
		tag      FlagTag
		expected any
	}{
		{
			name: "string default",
			tag: FlagTag{
				Type:    reflect.TypeFor[string](),
				Default: "hello",
			},
			expected: "hello",
		},
		{
			name: "bool true",
			tag: FlagTag{
				Type:    reflect.TypeFor[bool](),
				Default: "true",
			},
			expected: true,
		},
		{
			name: "bool false",
			tag: FlagTag{
				Type:    reflect.TypeFor[bool](),
				Default: "false",
			},
			expected: false,
		},
		{
			name: "int default",
			tag: FlagTag{
				Type:    reflect.TypeFor[int](),
				Default: "42",
			},
			expected: 42,
		},
		{
			name: "uint default",
			tag: FlagTag{
				Type:    reflect.TypeFor[uint](),
				Default: "42",
			},
			expected: uint(42),
		},
		{
			name: "uint32 default",
			tag: FlagTag{
				Type:    reflect.TypeFor[uint32](),
				Default: "4294967295",
			},
			expected: uint(4294967295),
		},
		{
			name: "float64 default",
			tag: FlagTag{
				Type:    reflect.TypeFor[float64](),
				Default: "3.14",
			},
			expected: 3.14,
		},
		{
			name: "empty default returns zero",
			tag: FlagTag{
				Type:    reflect.TypeFor[int](),
				Default: "",
			},
			expected: 0,
		},
		{
			name: "slice default",
			tag: FlagTag{
				Type:    reflect.TypeFor[[]string](),
				Default: "a,b,c",
			},
			expected: []string{"a", "b", "c"},
		},
		{
			name: "Duration default",
			tag: FlagTag{
				Type:    reflect.TypeFor[Duration](),
				Default: "5m",
			},
			expected: FromDuration(5 * time.Minute),
		},
		{
			name: "LogLevel default",
			tag: FlagTag{
				Type:    reflect.TypeFor[LogLevel](),
				Default: "info",
			},
			expected: "info", // Returns string for these types
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tag.DefaultValue()
			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestConfig_DefaultConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
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
