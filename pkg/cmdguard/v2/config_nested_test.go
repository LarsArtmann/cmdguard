package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type nestedDBConfig struct {
	Host string `flag:"host" default:"localhost" help:"database host"`
	Port int    `flag:"port" default:"5432"      help:"database port"`
}

type nestedRootConfig struct {
	Debug bool           `flag:"debug" default:"false" help:"debug mode"`
	DB    nestedDBConfig // no flag tag → recurse into it
}

func TestParseFlagTags_RecursesIntoNestedStructs(t *testing.T) {
	t.Parallel()

	tags, err := ParseFlagTags(&nestedRootConfig{})
	testutil.AssertNoError(t, err)

	names := make(map[string]bool)
	for _, tag := range tags {
		names[tag.Name] = true
	}

	for _, want := range []string{"debug", "host", "port"} {
		if !names[want] {
			t.Errorf("expected flag %q from nested struct, tags found: %v", want, names)
		}
	}
}

func TestParseFlagTags_NestedFieldIndex(t *testing.T) {
	t.Parallel()

	tags, err := ParseFlagTags(&nestedRootConfig{})
	testutil.AssertNoError(t, err)

	found := false

	for _, tag := range tags {
		if tag.Name == "host" {
			found = true
			if len(tag.Index) < 2 {
				t.Errorf("expected nested Index path for 'host', got %v", tag.Index)
			}
		}
	}

	if !found {
		t.Error("did not find 'host' tag")
	}
}

func TestNestedConfig_SetFieldByIndex(t *testing.T) {
	t.Parallel()

	cfg := &nestedRootConfig{}
	tags, err := ParseFlagTags(cfg)
	testutil.AssertNoError(t, err)

	found := false

	for _, tag := range tags {
		if tag.Name != "host" {
			continue
		}

		found = true
		err := setFieldByTag(cfg, tag, "db.example.com", globalTypeRegistry)
		testutil.AssertNoError(t, err)
		if cfg.DB.Host != "db.example.com" {
			t.Errorf("expected DB.Host='db.example.com', got %q", cfg.DB.Host)
		}
	}

	if !found {
		t.Error("did not find 'host' tag")
	}
}

func TestNestedConfig_DoesNotRecurseIntoDuration(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Timeout Duration `flag:"timeout" default:"5s"`
	}

	tags, err := ParseFlagTags(&cfg{})
	testutil.AssertNoError(t, err)

	if len(tags) != 1 {
		t.Fatalf("expected 1 tag (timeout), got %d: %v", len(tags), tags)
	}

	if tags[0].Name != "timeout" {
		t.Errorf("expected 'timeout', got %q", tags[0].Name)
	}
}

func TestNestedConfig_JSONFile(t *testing.T) {
	t.Parallel()

	jsonData := `{"debug": true, "DB": {"host": "prod.db", "port": 9999}}`

	cfg := &nestedRootConfig{}
	loader := &jsonLoader{}
	setFields, err := loader.Load([]byte(jsonData), cfg)
	testutil.AssertNoError(t, err)

	if !cfg.Debug {
		t.Error("expected Debug=true")
	}

	if cfg.DB.Host != "prod.db" {
		t.Errorf("expected DB.Host='prod.db', got %q", cfg.DB.Host)
	}

	if cfg.DB.Port != 9999 {
		t.Errorf("expected DB.Port=9999, got %d", cfg.DB.Port)
	}

	foundHost := false

	for _, f := range setFields {
		if f == "Host" {
			foundHost = true
		}
	}

	if !foundHost {
		t.Errorf("expected 'Host' in setFields, got %v", setFields)
	}
}

func TestNestedConfig_CLIIntegration(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("test", "1.0", nestedRootConfig{})
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"run",
		NoFlags{},
		func(_ context.Context, _ *nestedRootConfig, _ NoFlags) error {
			return nil
		},
		WithShort("run"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	execErr := cli.ExecuteWithArgs(context.Background(), []string{"run", "--host", "cli.db", "--port", "8080"})
	testutil.AssertNoError(t, execErr)
}
