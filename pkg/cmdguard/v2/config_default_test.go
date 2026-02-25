package v2

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag.DefaultValue()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_DefaultConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := Config{}
		tags, err := ParseFlagTags(cfg)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(tags), 4)

		// Find LogLevel tag
		var logLevelTag *FlagTag

		for _, tag := range tags {
			if tag.Name == "log-level" {
				logLevelTag = &tag

				break
			}
		}

		require.NotNil(t, logLevelTag)
		assert.Equal(t, "info", logLevelTag.Default)
		assert.Equal(t, "l", logLevelTag.Short)
	})
}
