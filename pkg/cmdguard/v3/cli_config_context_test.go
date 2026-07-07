package v3

import (
	"context"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

type ctxTestConfig struct {
	Name string `flag:"name" default:"world" help:"name"`
}

func TestConfigFromContext_StoredByPreRunE(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI[ctxTestConfig]("test", "Test", ctxTestConfig{})
	testutil.AssertNoError(t, err)

	var gotCfg *ctxTestConfig

	subCmd := &cobra.Command{
		Use: "sub",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, _ := ConfigFromContext[ctxTestConfig](cmd.Context())
			gotCfg = cfg

			return nil
		},
	}

	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertNoError(t, err)

	if gotCfg == nil {
		t.Fatal("ConfigFromContext returned nil after PersistentPreRunE ran")
	}

	if gotCfg.Name != "world" {
		t.Errorf("config Name = %q, want %q", gotCfg.Name, "world")
	}
}

func TestConfigFromContext_NilContext(t *testing.T) {
	t.Parallel()

	cfg, ok := ConfigFromContext[ctxTestConfig](context.TODO())
	if ok {
		t.Error("expected ok=false for nil context")
	}

	if cfg != nil {
		t.Error("expected nil config for nil context")
	}
}

func TestConfigFromContext_NotStored(t *testing.T) {
	t.Parallel()

	_, ok := ConfigFromContext[ctxTestConfig](context.Background())
	if ok {
		t.Error("expected ok=false for context without config")
	}
}

func TestConfigFromContext_WrongType(t *testing.T) {
	t.Parallel()

	type other struct{}

	ctx := context.WithValue(context.Background(), configKey, &other{})
	_, ok := ConfigFromContext[ctxTestConfig](ctx)
	if ok {
		t.Error("expected ok=false when config is wrong type")
	}
}
