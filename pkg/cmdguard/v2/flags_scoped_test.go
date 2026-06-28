package v2

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

// scopedConfig mirrors a real app layout: a few flags are global (inherited by
// every subcommand) and a few are local (only meaningful on the root run, plus
// any subcommand that explicitly opts in).
type scopedConfig struct {
	Global  string `default:"g-default" flag:"global"  help:"global flag"`
	Verbose bool   `default:"false"     flag:"verbose" help:"verbose"`
	Build   string `default:"full"      flag:"build"   help:"build mode"  local:"true"`
	Step    string `default:""          flag:"step"    help:"single step" local:"true"`
	Fix     bool   `default:"false"     flag:"fix"     help:"auto fix"    local:"true"`
}

// mergeInherited mirrors what cobra's execute() does before PersistentPreRunE:
// it merges parent persistent flags into the command so lookups see inherited
// flags. Calling InheritedFlags() has mergePersistentFlags() as a side effect.
func mergeInherited(cmd *cobra.Command) {
	_ = cmd.InheritedFlags()
}

func TestRegisterScopedFlags_RootGetsAll(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(scopedConfig{})
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}

	err = registry.RegisterScopedFlags(root)
	testutil.AssertNoError(t, err)
	mergeInherited(root)

	// Persistent flags land on root.
	testutil.AssertFlagRegistered(t, root, "global")
	testutil.AssertFlagRegistered(t, root, "verbose")
	// Local flags also land on root (the root run uses them).
	testutil.AssertFlagRegistered(t, root, "build")
	testutil.AssertFlagRegistered(t, root, "step")
	testutil.AssertFlagRegistered(t, root, "fix")
}

func TestLocalFlagsNotInheritedBySubcommand(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(scopedConfig{})
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}

	err = registry.RegisterScopedFlags(root)
	testutil.AssertNoError(t, err)

	sub := &cobra.Command{Use: "sub"}
	root.AddCommand(sub)
	mergeInherited(sub)

	// Persistent flags are inherited by the subcommand.
	testutil.AssertFlagRegistered(t, sub, "global")
	testutil.AssertFlagRegistered(t, sub, "verbose")
	// Local flags are NOT inherited — this is the whole point.
	testutil.AssertFlagNotRegistered(t, sub, "build")
	testutil.AssertFlagNotRegistered(t, sub, "step")
	testutil.AssertFlagNotRegistered(t, sub, "fix")
}

func TestRegisterLocalFlags_OnSubcommand(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(scopedConfig{})
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}
	_ = registry.RegisterScopedFlags(root)

	sub := &cobra.Command{Use: "sub"}
	root.AddCommand(sub)

	err = registry.RegisterLocalFlags(sub)
	testutil.AssertNoError(t, err)
	mergeInherited(sub)

	// Opting in makes the local group available on the subcommand.
	testutil.AssertFlagRegistered(t, sub, "build")
	testutil.AssertFlagRegistered(t, sub, "step")
	testutil.AssertFlagRegistered(t, sub, "fix")
	// Persistent flags are still inherited alongside the local group.
	testutil.AssertFlagRegistered(t, sub, "global")
}

// ParseFlags must not error when a local flag is absent on the running command.
// This is what lets utility subcommands run without the root's execution flags.
func TestParseFlags_SkipsAbsentLocalFlags(t *testing.T) {
	t.Parallel()

	cfg := scopedConfig{}
	registry, err := NewFlagRegistry(cfg)
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}
	_ = registry.RegisterScopedFlags(root)

	// A subcommand that does NOT opt into the local flag group.
	sub := &cobra.Command{Use: "sub"}
	root.AddCommand(sub)
	mergeInherited(sub)

	err = registry.ParseFlags(sub, &cfg)
	testutil.AssertNoError(t, err)

	// Persistent flag default is applied; local flags are untouched at default.
	testutil.AssertFieldEqString(t, cfg.Global, "g-default", "Global default")
	testutil.AssertFieldEqString(t, cfg.Build, "full", "Build default (untouched)")
}

// ParseFlags must read a local flag when the running command opted into the group.
func TestParseFlags_ReadsLocalFlagOnOptedInSubcommand(t *testing.T) {
	t.Parallel()

	cfg := scopedConfig{}
	registry, err := NewFlagRegistry(cfg)
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}
	_ = registry.RegisterScopedFlags(root)

	sub := &cobra.Command{Use: "sub"}
	root.AddCommand(sub)
	_ = registry.RegisterLocalFlags(sub)
	mergeInherited(sub)

	// Simulate the user passing --build fast on the subcommand.
	err = sub.Flags().Set("build", "fast")
	testutil.AssertNoError(t, err)

	err = registry.ParseFlags(sub, &cfg)
	testutil.AssertNoError(t, err)

	testutil.AssertFieldEqString(t, cfg.Build, "fast", "Build should be parsed from the local flag")
}

func TestRegisterLocalFlags_SkipsFlagsSubcommandAlreadyDefines(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(scopedConfig{})
	testutil.AssertNoError(t, err)

	root := &cobra.Command{Use: "root"}
	_ = registry.RegisterScopedFlags(root)

	// The subcommand defines its OWN build flag with a different default/help.
	sub := &cobra.Command{Use: "sub"}
	sub.Flags().String("build", "watch", "subcommand-specific build mode")
	root.AddCommand(sub)

	err = registry.RegisterLocalFlags(sub)
	testutil.AssertNoError(t, err)

	mergeInherited(sub)

	// The subcommand's own definition is preserved (not overwritten).
	build := sub.Flags().Lookup("build")
	testutil.AssertNotNil(t, build)

	if def, _ := build.Value.(interface{ String() string }); def.String() != "watch" {
		t.Errorf("expected subcommand's own build default %q, got %q", "watch", def.String())
	}

	// The rest of the local group is still registered.
	testutil.AssertFlagRegistered(t, sub, "step")
	testutil.AssertFlagRegistered(t, sub, "fix")
}

func TestParseLocalTag(t *testing.T) {
	t.Parallel()

	t.Run("explicit true", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x" local:"true"`
		}

		tags, err := ParseFlagTags(c{})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldLen(t, tags, 1, "tags")
		testutil.AssertBoolTrue(t, tags[0].Local, "tags[0].Local")
	})

	t.Run("absent defaults to false", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x"`
		}

		tags, err := ParseFlagTags(c{})
		testutil.AssertNoError(t, err)
		testutil.AssertBoolFalse(t, tags[0].Local, "tags[0].Local default")
	})

	t.Run("invalid value errors", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x" local:"maybe"`
		}

		_, err := ParseFlagTags(c{})
		testutil.AssertExpectedError(t, err)
		assertErrorContains(t, err, "invalid local tag")
	})
}
