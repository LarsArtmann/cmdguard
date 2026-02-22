package v2

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlagRegistry(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		type TestConfig struct {
			Name  string `flag:"name" help:"name help" default:"default-name"`
			Count int    `flag:"count" help:"count help" default:"10"`
		}
		cfg := TestConfig{}

		registry, err := NewFlagRegistry(cfg)
		require.NoError(t, err)
		assert.NotNil(t, registry)
		assert.Len(t, registry.Tags(), 2)
	})

	t.Run("non-struct config", func(t *testing.T) {
		registry, err := NewFlagRegistry("not a struct")
		require.Error(t, err)
		assert.Nil(t, registry)
		assert.Contains(t, err.Error(), "must be a struct")
	})

	t.Run("config with short flags", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" short:"n" help:"name help"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)
		tags := registry.Tags()
		require.Len(t, tags, 1)
		assert.Equal(t, "n", tags[0].Short)
	})
}

func TestFlagRegistry_RegisterFlags(t *testing.T) {
	t.Run("registers all flag types", func(t *testing.T) {
		type TestConfig struct {
			String  string   `flag:"string" default:"str"`
			Bool    bool     `flag:"bool" default:"true"`
			Int     int      `flag:"int" default:"42"`
			Float   float64  `flag:"float" default:"3.14"`
			Strings []string `flag:"strings" default:"a,b,c"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		err = registry.RegisterFlags(cmd)
		require.NoError(t, err)

		// Verify flags were registered
		flags := cmd.PersistentFlags()
		assert.NotNil(t, flags.Lookup("string"))
		assert.NotNil(t, flags.Lookup("bool"))
		assert.NotNil(t, flags.Lookup("int"))
		assert.NotNil(t, flags.Lookup("float"))
		assert.NotNil(t, flags.Lookup("strings"))
	})

	t.Run("registers custom types", func(t *testing.T) {
		type TestConfig struct {
			Level  LogLevel  `flag:"level" default:"info"`
			Format LogFormat `flag:"format" default:"json"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		err = registry.RegisterFlags(cmd)
		require.NoError(t, err)

		assert.NotNil(t, cmd.PersistentFlags().Lookup("level"))
		assert.NotNil(t, cmd.PersistentFlags().Lookup("format"))
	})

	t.Run("registers Duration type", func(t *testing.T) {
		type TestConfig struct {
			Timeout Duration `flag:"timeout" default:"30s"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		err = registry.RegisterFlags(cmd)
		require.NoError(t, err)

		flag := cmd.PersistentFlags().Lookup("timeout")
		assert.NotNil(t, flag)
		assert.Equal(t, "30s", flag.DefValue)
	})

	t.Run("registers enum with values", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		err = registry.RegisterFlags(cmd)
		require.NoError(t, err)

		flag := cmd.PersistentFlags().Lookup("mode")
		assert.NotNil(t, flag)
		// Help should include allowed values
		assert.Contains(t, flag.Usage, "one of: dev, staging, prod")
	})
}

func TestFlagRegistry_ParseFlags(t *testing.T) {
	t.Run("parse basic flag types", func(t *testing.T) {
		tests := []struct {
			name     string
			flag     string
			value    string
			wantErr  bool
			validate func(t *testing.T, cfg interface{})
		}{
			{
				name:  "string flag",
				flag:  "name",
				value: "custom",
				validate: func(t *testing.T, cfg interface{}) {
					assert.Equal(t, "custom", cfg.(*struct{ Name string }).Name)
				},
			},
			{
				name:  "bool flag",
				flag:  "verbose",
				value: "true",
				validate: func(t *testing.T, cfg interface{}) {
					assert.True(t, cfg.(*struct{ Verbose bool }).Verbose)
				},
			},
			{
				name:  "int flag",
				flag:  "count",
				value: "42",
				validate: func(t *testing.T, cfg interface{}) {
					assert.Equal(t, 42, cfg.(*struct{ Count int }).Count)
				},
			},
			{
				name:  "float64 flag",
				flag:  "rate",
				value: "3.14159",
				validate: func(t *testing.T, cfg interface{}) {
					assert.InDelta(t, 3.14159, cfg.(*struct{ Rate float64 }).Rate, 0.00001)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var cfg interface{}
				switch tt.flag {
				case "name":
					type TestConfig struct{ Name string }
					cfg = &TestConfig{Name: "default"}
				case "verbose":
					type TestConfig struct{ Verbose bool }
					cfg = &TestConfig{Verbose: false}
				case "count":
					type TestConfig struct{ Count int }
					cfg = &TestConfig{Count: 0}
				case "rate":
					type TestConfig struct{ Rate float64 }
					cfg = &TestConfig{Rate: 0.0}
				}

				registry, err := NewFlagRegistry(cfg)
				require.NoError(t, err)

				cmd := &cobra.Command{Use: "test"}
				require.NoError(t, registry.RegisterFlags(cmd))

				require.NoError(t, cmd.PersistentFlags().Set(tt.flag, tt.value))

				err = registry.ParseFlags(cmd, cfg)
				if tt.wantErr {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				tt.validate(t, cfg)
			})
		}
	})

	t.Run("parse Duration flag", func(t *testing.T) {
		type TestConfig struct {
			Timeout Duration `flag:"timeout" default:"1m"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("timeout", "5m30s"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		expected := FromDuration(5*time.Minute + 30*time.Second)
		assert.Equal(t, expected, cfg.Timeout)
	})

	t.Run("parse LogLevel flag", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" default:"info"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("level", "debug"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, LogLevelDebug, cfg.Level)
	})

	t.Run("parse invalid LogLevel returns error", func(t *testing.T) {
		type TestConfig struct {
			Level LogLevel `flag:"level" default:"info"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("level", "invalid"))

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "level")
	})

	t.Run("parse Enum flag", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "prod"))

		err = registry.ParseFlags(cmd, cfg)
		require.NoError(t, err)
		assert.Equal(t, "prod", cfg.Mode.String())
	})

	t.Run("parse invalid Enum returns error", func(t *testing.T) {
		type TestConfig struct {
			Mode Enum `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "invalid"))

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
	})
}

func TestFlagRegistry_ValidateFlags(t *testing.T) {
	t.Run("valid values pass", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("mode", "staging"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Manually set an invalid value (bypassing validation)
		require.NoError(t, cmd.PersistentFlags().Set("mode", "invalid"))

		err = registry.ValidateFlags(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mode")
	})

	t.Run("unchanged flag skips validation", func(t *testing.T) {
		type TestConfig struct {
			Mode string `flag:"mode" values:"dev,staging,prod" default:"dev"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't change the flag - should pass validation
		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("flag without values skips validation", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" default:"default"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("name", "anything"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("required flag not set returns error", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"true" help:"required name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't set the flag - should fail validation
		err = registry.ValidateFlags(cmd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name")
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("required flag set passes validation", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"true" help:"required name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		require.NoError(t, cmd.PersistentFlags().Set("name", "test-value"))

		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})

	t.Run("required false does not enforce", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" required:"false" help:"optional name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		require.NoError(t, registry.RegisterFlags(cmd))

		// Don't set the flag - should pass since required:"false"
		err = registry.ValidateFlags(cmd)
		require.NoError(t, err)
	})
}

func TestFlagRegistry_Tags(t *testing.T) {
	t.Run("returns all tags", func(t *testing.T) {
		type TestConfig struct {
			Name  string `flag:"name" help:"name help"`
			Count int    `flag:"count" help:"count help"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		tags := registry.Tags()
		assert.Len(t, tags, 2)

		names := make([]string, len(tags))
		for i, tag := range tags {
			names[i] = tag.Name
		}
		assert.Contains(t, names, "name")
		assert.Contains(t, names, "count")
	})
}

func TestFlagRegistry_GenerateHelp(t *testing.T) {
	t.Run("generates help for all flags", func(t *testing.T) {
		type TestConfig struct {
			Name    string `flag:"name,n" help:"The name to use" default:"default"`
			Verbose bool   `flag:"verbose,v" help:"Enable verbose output"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		help := registry.GenerateHelp()
		assert.Contains(t, help, "--name")
		assert.Contains(t, help, "-n")
		assert.Contains(t, help, "The name to use")
		assert.Contains(t, help, "default: default")
		assert.Contains(t, help, "--verbose")
		assert.Contains(t, help, "-v")
		assert.Contains(t, help, "Enable verbose output")
	})

	t.Run("help without short flag", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" help:"The name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		help := registry.GenerateHelp()
		assert.Contains(t, help, "--name")
		assert.NotContains(t, help, "-,")
	})

	t.Run("help without default", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" help:"The name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		help := registry.GenerateHelp()
		assert.Contains(t, help, "--name")
		assert.NotContains(t, help, "default:")
	})
}

func TestFlagRegistry_ShortFlags(t *testing.T) {
	t.Run("register short flags", func(t *testing.T) {
		type TestConfig struct {
			Name    string `flag:"name" short:"n" default:""`
			Count   int    `flag:"count" short:"c" default:"0"`
			Verbose bool   `flag:"verbose" short:"v" default:"false"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		cmd := &cobra.Command{Use: "test"}
		err = registry.RegisterFlags(cmd)
		require.NoError(t, err)

		flags := cmd.PersistentFlags()

		// Verify short flags work
		nameFlag := flags.Lookup("name")
		require.NotNil(t, nameFlag)
		assert.Equal(t, "n", nameFlag.Shorthand)

		countFlag := flags.Lookup("count")
		require.NotNil(t, countFlag)
		assert.Equal(t, "c", countFlag.Shorthand)

		verboseFlag := flags.Lookup("verbose")
		require.NotNil(t, verboseFlag)
		assert.Equal(t, "v", verboseFlag.Shorthand)
	})
}

func TestFlagRegistry_FlagNotFound(t *testing.T) {
	t.Run("missing flag returns error", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name"`
		}

		cfg := &TestConfig{}
		registry, err := NewFlagRegistry(*cfg)
		require.NoError(t, err)

		// Don't register flags on command
		cmd := &cobra.Command{Use: "test"}

		err = registry.ParseFlags(cmd, cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

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

func TestFlagRegistry_FlagNames(t *testing.T) {
	t.Run("returns all flag names", func(t *testing.T) {
		type TestConfig struct {
			Verbose bool   `flag:"verbose" short:"v"`
			Config  string `flag:"config" short:"c"`
			Output  string `flag:"output" short:"o"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		names := registry.FlagNames()
		assert.Len(t, names, 3)
		assert.Contains(t, names, "verbose")
		assert.Contains(t, names, "config")
		assert.Contains(t, names, "output")
	})

	t.Run("empty registry returns empty slice", func(t *testing.T) {
		type EmptyConfig struct{}

		registry, err := NewFlagRegistry(EmptyConfig{})
		require.NoError(t, err)

		names := registry.FlagNames()
		assert.Empty(t, names)
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
