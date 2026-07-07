package v3

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestNewFlagRegistry(t *testing.T) {
	t.Parallel()
	t.Run("valid config", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name  string `default:"default-name" flag:"name"  help:"name help"`
			Count int    `default:"10"           flag:"count" help:"count help"`
		}

		cfg := testConfig{}

		registry, err := NewFlagRegistry(cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertNotNil(t, registry)

		testutil.AssertFieldLen(t, registry.Tags(), 2, "Tags")
	})

	t.Run("non-struct config", func(t *testing.T) {
		t.Parallel()

		registry, err := NewFlagRegistry("not a struct")
		testutil.AssertExpectedError(t, err)

		testutil.AssertNil(t, registry)

		assertErrorContains(t, err, "expected struct")
	})

	t.Run("config with short flags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name string `flag:"name" help:"name help" short:"n"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		testutil.AssertNoError(t, err)

		tags := registry.Tags()
		testutil.AssertFieldLen(t, tags, 1, "tags")

		if len(tags) > 0 {
			testutil.AssertFieldEqString(t, tags[0].Short, "n", "tags[0].Short")
		}
	})
}

func TestFlagRegistry_RegisterFlags(t *testing.T) {
	t.Parallel()
	t.Run("registers all flag types", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			String  string   `default:"str"   flag:"string"`
			Bool    bool     `default:"true"  flag:"bool"`
			Int     int      `default:"42"    flag:"int"`
			Uint    uint     `default:"10"    flag:"uint"`
			Uint64  uint64   `default:"1024"  flag:"uint64"`
			Float   float64  `default:"3.14"  flag:"float"`
			Strings []string `default:"a,b,c" flag:"strings"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		testutil.AssertNoError(t, err)

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		testutil.AssertNoError(t, err)

		testutil.AssertFlagRegistered(t, cmd, "string")
		testutil.AssertFlagRegistered(t, cmd, "bool")
		testutil.AssertFlagRegistered(t, cmd, "int")
		testutil.AssertFlagRegistered(t, cmd, "uint")
		testutil.AssertFlagRegistered(t, cmd, "uint64")
		testutil.AssertFlagRegistered(t, cmd, "float")
		testutil.AssertFlagRegistered(t, cmd, "strings")
	})

	t.Run("registers custom types", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Level  LogLevel  `default:"info" flag:"level"`
			Format LogFormat `default:"json" flag:"format"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		testutil.AssertNoError(t, err)

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		testutil.AssertNoError(t, err)

		testutil.AssertFlagRegistered(t, cmd, "level")
		testutil.AssertFlagRegistered(t, cmd, "format")
	})

	t.Run("registers Duration type", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Timeout Duration `default:"30s" flag:"timeout"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		testutil.AssertNoError(t, err)

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		testutil.AssertNoError(t, err)

		flag := cmd.Flags().Lookup("timeout")
		if flag == nil {
			t.Fatal("expected 'timeout' flag to be registered")
		}

		testutil.AssertFieldEqString(t, flag.DefValue, "30s", "DefValue")
	})

	t.Run("registers enum with values", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Mode Enum `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		flag := cmd.Flags().Lookup("mode")
		if flag == nil {
			t.Fatal("expected 'mode' flag to be registered")
		}

		testutil.AssertStringFieldContains(
			t,
			flag.Usage,
			"one of: dev, staging, prod",
			"flag.Usage",
		)
	})
}

func TestFlagRegistry_Tags(t *testing.T) {
	t.Parallel()
	t.Run("returns all tags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name  string `flag:"name"  help:"name help"`
			Count int    `flag:"count" help:"count help"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tags := registry.Tags()
		testutil.AssertFieldLen(t, tags, 2, "tags")

		names := make([]string, len(tags))
		for i, tag := range tags {
			names[i] = tag.Name
		}

		testutil.AssertContainsString(t, names, "name")
		testutil.AssertContainsString(t, names, "count")
	})
}

func TestFlagRegistry_FlagNames(t *testing.T) {
	t.Parallel()
	t.Run("returns all flag names", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Verbose bool   `flag:"verbose" short:"v"`
			Config  string `flag:"config"  short:"c"`
			Output  string `flag:"output"  short:"o"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		names := registry.FlagNames()
		testutil.AssertFieldLen(t, names, 3, "names")

		testutil.AssertContainsString(t, names, "verbose")
		testutil.AssertContainsString(t, names, "config")
		testutil.AssertContainsString(t, names, "output")
	})

	t.Run("empty registry returns empty slice", func(t *testing.T) {
		t.Parallel()

		type emptyConfig struct{}

		registry, err := NewFlagRegistry(emptyConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		names := registry.FlagNames()
		testutil.AssertFieldLen(t, names, 0, "names")
	})
}
