package v2

import (
	"testing"

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

	t.Run("help formatting without optional elements", func(t *testing.T) {
		type TestConfig struct {
			Name string `flag:"name" help:"The name"`
		}

		registry, err := NewFlagRegistry(TestConfig{})
		require.NoError(t, err)

		help := registry.GenerateHelp()
		assert.Contains(t, help, "--name")
		assert.NotContains(t, help, "-,", "should not show empty short flag")
		assert.NotContains(t, help, "default:", "should not show default when not set")
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
