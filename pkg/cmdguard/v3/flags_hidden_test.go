package v3

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestHiddenFlag_NotShownInHelp(t *testing.T) {
	t.Parallel()

	type config struct {
		Public string `flag:"public" help:"visible flag" local:"true"`
		Secret string `flag:"secret" help:"hidden flag"  local:"true" hidden:"true"`
	}

	registry, err := NewFlagRegistry(config{})
	testutil.AssertNoError(t, err)

	cmd := &cobra.Command{Use: "test"}
	err = registry.RegisterScopedFlags(cmd)
	testutil.AssertNoError(t, err)

	publicFlag := cmd.Flags().Lookup("public")
	if publicFlag == nil {
		t.Fatal("public flag not registered")
	}

	if publicFlag.Hidden {
		t.Error("public flag should NOT be hidden")
	}

	secretFlag := cmd.Flags().Lookup("secret")
	if secretFlag == nil {
		t.Fatal("secret flag not registered")
	}

	if !secretFlag.Hidden {
		t.Error("secret flag should be hidden")
	}
}

func TestHiddenFlag_StillFunctional(t *testing.T) {
	t.Parallel()

	type config struct {
		Debug string `flag:"debug" help:"debug" hidden:"true" local:"true"`
	}

	registry, err := NewFlagRegistry(config{})
	testutil.AssertNoError(t, err)

	cmd := &cobra.Command{Use: "test"}
	err = registry.RegisterScopedFlags(cmd)
	testutil.AssertNoError(t, err)

	err = cmd.Flags().Set("debug", "value")
	testutil.AssertNoError(t, err)
}

func TestParseHiddenTag(t *testing.T) {
	t.Parallel()

	t.Run("explicit true", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x" hidden:"true"`
		}

		tags, err := ParseFlagTags(c{})
		testutil.AssertNoError(t, err)
		testutil.AssertFieldLen(t, tags, 1, "tags")
		testutil.AssertBoolTrue(t, tags[0].Hidden, "tags[0].Hidden")
	})

	t.Run("absent defaults to false", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x"`
		}

		tags, err := ParseFlagTags(c{})
		testutil.AssertNoError(t, err)
		testutil.AssertBoolFalse(t, tags[0].Hidden, "tags[0].Hidden default")
	})

	t.Run("invalid value errors", func(t *testing.T) {
		t.Parallel()

		type c struct {
			X string `flag:"x" hidden:"maybe"`
		}

		_, err := ParseFlagTags(c{})
		testutil.AssertExpectedError(t, err)
	})
}
