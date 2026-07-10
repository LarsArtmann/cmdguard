package v3

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

type customColor string

type colorPlugin struct{}

func (colorPlugin) Name() string { return "color" }

func (colorPlugin) Register(r PluginRegistrar) error {
	r.TypeHandlerFunc(reflect.TypeFor[customColor](), TypeHandlerFunc{
		RegisterFunc: registerStringFlagFromTag,
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return customColor(value), nil
		},
		DefaultFunc: func(tag FlagTag) any {
			return customColor(tag.Default)
		},
	})

	r.Validator("nodash", func(value string) error {
		if strings.Contains(value, "-") {
			return errContainsDash
		}

		return nil
	})

	return nil
}

var errContainsDash = dashError("value must not contain dashes")

type dashError string

func (e dashError) Error() string { return string(e) }

func TestPlugin_RegisterTypeHandlerAndValidator(t *testing.T) {
	t.Parallel()

	err := RegisterPlugin(colorPlugin{})
	testutil.AssertNoError(t, err)

	h, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[customColor]())
	testutil.AssertBoolTrue(t, ok, "expected customColor handler to be registered")

	parsed, err := h.Parse("red", FlagTag{})
	testutil.AssertNoError(t, err)
	if parsed != customColor("red") {
		t.Errorf("expected customColor('red'), got %v", parsed)
	}

	validator, ok := lookupValidator("nodash")
	testutil.AssertBoolTrue(t, ok, "expected nodash validator to be registered")
	testutil.AssertNoError(t, validator("blue"))
	testutil.AssertExpectedError(t, validator("blue-red"))
}

func TestPlugin_NilReturnsError(t *testing.T) {
	t.Parallel()

	err := RegisterPlugin(nil)
	testutil.AssertExpectedError(t, err)
}

func TestPlugin_WithPluginOption(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Color customColor `flag:"color" default:"red"`
	}

	_ = RegisterPlugin(colorPlugin{})

	cli, err := NewCLI("plugintest", "1.0", cfg{})
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, c *cfg, _ NoFlags) error {
			if c.Color != "blue" {
				t.Errorf("expected 'blue', got %q", c.Color)
			}

			return nil
		},
		WithShort("run"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	execErr := cli.ExecuteWithArgs(context.Background(), []string{"run", "--color", "blue"})
	testutil.AssertNoError(t, execErr)
}

var errPluginFail = errors.New("plugin registration failed")

type failingPlugin struct{}

func (failingPlugin) Name() string { return "failing" }

func (failingPlugin) Register(_ PluginRegistrar) error {
	return errPluginFail
}

func TestPlugin_WithPluginReturnsErrorOnFailedRegister(t *testing.T) {
	t.Parallel()

	_, err := NewCLI("failingplugin", "1.0", struct{}{}, WithPlugin(failingPlugin{}))
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, errPluginFail) {
		t.Errorf("expected error to wrap errPluginFail, got: %v", err)
	}
}
