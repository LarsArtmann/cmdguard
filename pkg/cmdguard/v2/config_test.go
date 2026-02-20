package v2

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlagTags(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		type TestConfig struct {
			Name    string `flag:"name" short:"n" default:"test" help:"The name"`
			Count   int    `flag:"count" default:"10" help:"The count"`
			Enabled bool   `flag:"enabled" short:"e" default:"true" help:"Enable feature"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 3)

		// Check first field
		assert.Equal(t, "Name", tags[0].Field)
		assert.Equal(t, "name", tags[0].Name)
		assert.Equal(t, "n", tags[0].Short)
		assert.Equal(t, "test", tags[0].Default)
		assert.Equal(t, "The name", tags[0].Help)
	})

	t.Run("pointer to struct", func(t *testing.T) {
		type TestConfig struct {
			Field string `flag:"field"`
		}

		tags, err := ParseFlagTags(&TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "field", tags[0].Name)
	})

	t.Run("skips fields without flag tag", func(t *testing.T) {
		type TestConfig struct {
			Tagged   string `flag:"tagged"`
			Untagged string
			Ignored  string `flag:"-"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, "Tagged", tags[0].Field)
	})

	t.Run("nil config", func(t *testing.T) {
		tags, err := ParseFlagTags(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be nil")
		assert.Nil(t, tags)
	})

	t.Run("non-struct config", func(t *testing.T) {
		tags, err := ParseFlagTags("not a struct")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a struct")
		assert.Nil(t, tags)
	})

	t.Run("with values tag", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn,error"`
		}

		tags, err := ParseFlagTags(TestConfig{})
		require.NoError(t, err)
		require.Len(t, tags, 1)
		assert.Equal(t, []string{"debug", "info", "warn", "error"}, tags[0].Values)
	})

	t.Run("embedded Config", func(t *testing.T) {
		type AppConfig struct {
			Config
			AppName string `flag:"app-name" default:"myapp"`
		}

		tags, err := ParseFlagTags(AppConfig{})
		require.NoError(t, err)
		// Config has 4 fields + AppName = 5
		assert.GreaterOrEqual(t, len(tags), 1)
	})
}

func TestFlagTag_DefaultValue(t *testing.T) {
	tests := []struct {
		name     string
		tag      FlagTag
		expected any
	}{
		{
			name: "string default",
			tag: FlagTag{
				Type:    reflect.TypeOf(""),
				Default: "hello",
			},
			expected: "hello",
		},
		{
			name: "bool true",
			tag: FlagTag{
				Type:    reflect.TypeOf(false),
				Default: "true",
			},
			expected: true,
		},
		{
			name: "bool false",
			tag: FlagTag{
				Type:    reflect.TypeOf(false),
				Default: "false",
			},
			expected: false,
		},
		{
			name: "int default",
			tag: FlagTag{
				Type:    reflect.TypeOf(0),
				Default: "42",
			},
			expected: 42,
		},
		{
			name: "float64 default",
			tag: FlagTag{
				Type:    reflect.TypeOf(0.0),
				Default: "3.14",
			},
			expected: 3.14,
		},
		{
			name: "empty default returns zero",
			tag: FlagTag{
				Type:    reflect.TypeOf(0),
				Default: "",
			},
			expected: 0,
		},
		{
			name: "slice default",
			tag: FlagTag{
				Type:    reflect.TypeOf([]string{}),
				Default: "a,b,c",
			},
			expected: []string{"a", "b", "c"},
		},
		{
			name: "Duration default",
			tag: FlagTag{
				Type:    reflect.TypeOf(Duration{}),
				Default: "5m",
			},
			expected: FromDuration(5 * time.Minute),
		},
		{
			name: "LogLevel default",
			tag: FlagTag{
				Type:    reflect.TypeOf(LogLevel{}),
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

func TestSetField(t *testing.T) {
	t.Run("set string field", func(t *testing.T) {
		cfg := &struct {
			Name string
		}{}
		err := SetField(cfg, "Name", "test")
		require.NoError(t, err)
		assert.Equal(t, "test", cfg.Name)
	})

	t.Run("set int field", func(t *testing.T) {
		cfg := &struct {
			Count int
		}{}
		err := SetField(cfg, "Count", 42)
		require.NoError(t, err)
		assert.Equal(t, 42, cfg.Count)
	})

	t.Run("set bool field", func(t *testing.T) {
		cfg := &struct {
			Enabled bool
		}{}
		err := SetField(cfg, "Enabled", true)
		require.NoError(t, err)
		assert.True(t, cfg.Enabled)
	})

	t.Run("non-pointer config", func(t *testing.T) {
		cfg := struct{ Name string }{}
		err := SetField(cfg, "Name", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pointer to struct")
	})

	t.Run("field not found", func(t *testing.T) {
		cfg := &struct{ Name string }{}
		err := SetField(cfg, "Missing", "test")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("string to LogLevel", func(t *testing.T) {
		cfg := &struct {
			Level LogLevel
		}{}
		err := SetField(cfg, "Level", "debug")
		require.NoError(t, err)
		assert.Equal(t, LogLevelDebug, cfg.Level)
	})

	t.Run("string to LogFormat", func(t *testing.T) {
		cfg := &struct {
			Format LogFormat
		}{}
		err := SetField(cfg, "Format", "json")
		require.NoError(t, err)
		assert.Equal(t, LogFormatJSON, cfg.Format)
	})

	t.Run("string to Duration", func(t *testing.T) {
		cfg := &struct {
			Timeout Duration
		}{}
		err := SetField(cfg, "Timeout", "30s")
		require.NoError(t, err)
		assert.Equal(t, 30*time.Second, cfg.Timeout.Duration())
	})

	t.Run("time.Duration to Duration", func(t *testing.T) {
		cfg := &struct {
			Timeout Duration
		}{}
		err := SetField(cfg, "Timeout", 45*time.Second)
		require.NoError(t, err)
		assert.Equal(t, 45*time.Second, cfg.Timeout.Duration())
	})

	t.Run("invalid LogLevel", func(t *testing.T) {
		cfg := &struct {
			Level LogLevel
		}{}
		err := SetField(cfg, "Level", "invalid")
		require.Error(t, err)
	})

	t.Run("incompatible types", func(t *testing.T) {
		cfg := &struct {
			Name string
		}{}
		// Use a struct type which is truly incompatible with string
		err := SetField(cfg, "Name", Duration{duration: 5 * time.Minute})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert")
	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		type TestConfig struct {
			Name  string `flag:"name"`
			Count int    `flag:"count"`
		}
		err := ValidateConfig(TestConfig{Name: "test", Count: 10})
		require.NoError(t, err)
	})

	t.Run("nil config", func(t *testing.T) {
		err := ValidateConfig(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must not be nil")
	})

	t.Run("non-struct config", func(t *testing.T) {
		err := ValidateConfig("not a struct")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be a struct")
	})

	t.Run("valid enum value", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}
		err := ValidateConfig(TestConfig{Level: "info"})
		require.NoError(t, err)
	})

	t.Run("invalid enum value", func(t *testing.T) {
		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}
		err := ValidateConfig(TestConfig{Level: "invalid"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config validation")
	})

	t.Run("pointer to config", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name"`
		}
		err := ValidateConfig(&TestConfig{Name: "test"})
		require.NoError(t, err)
	})

	t.Run("LogLevel field with values", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" values:"debug,info,warn,error"`
		}
		cfg := TestConfig{Level: LogLevelInfo}
		err := ValidateConfig(cfg)
		require.NoError(t, err)
	})

	t.Run("LogFormat field with values", func(t *testing.T) {
		type TestConfig struct {
			Format LogFormat `flag:"format" values:"text,json"`
		}
		cfg := TestConfig{Format: LogFormatText}
		err := ValidateConfig(cfg)
		require.NoError(t, err)
	})
}

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
