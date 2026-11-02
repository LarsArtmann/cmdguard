package schema_test

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/schema"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

var updateGolden = flag.Bool("update-golden", false, "rewrite golden schema fixture files")

// serverConfig is an array element: the schema must describe its properties
// under "items".
type serverConfig struct {
	Host string `json:"host" flag:"var-host" default:"localhost" help:"Server host"`
	Port int    `json:"port" flag:"var-port" default:"8080"      help:"Server port"`
}

// databaseConfig is a nested config group.
type databaseConfig struct {
	Host     string `json:"host"     flag:"db-host"     default:"localhost" env:"DB_HOST" help:"Database host"`
	Port     int    `json:"port"     flag:"db-port"     default:"5432"                    help:"Database port"     validate:"min=1,max=65535"`
	Password string `json:"password" flag:"db-password"                                   help:"Database password"                            secret:"true"`
}

// fixtureConfig exercises every mapping the emitter claims to support.
type fixtureConfig struct {
	LogLevel string         `json:"log_level" flag:"log-level" default:"info"  values:"debug,info,warn,error" help:"Log level"`
	Timeout  time.Duration  `json:"timeout"   flag:"timeout"   default:"30s"                                  help:"Request timeout"`
	Workers  int            `json:"workers"   flag:"workers"   default:"4"                                    help:"Worker count"     validate:"min=1,max=64"`
	Ratio    float64        `json:"ratio"     flag:"ratio"     default:"0.5"                                  help:"Sampling ratio"   validate:"min=0,max=1"`
	Debug    bool           `json:"debug"     flag:"debug"     default:"false"                                help:"Debug mode"`
	Name     string         `json:"name"      flag:"name"                                                     help:"Application name" validate:"minlen=3,maxlen=20,regex=^[a-z]+$" required:"true"`
	Tags     []string       `json:"tags"      flag:"tag"                                                      help:"Tags"`
	Database databaseConfig `json:"database"`
	Servers  []serverConfig `json:"servers"                                                                   help:"Servers"`
}

func TestFromConfig_Golden(t *testing.T) {
	t.Parallel()

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	got, err := json.Marshal(document, jsontext.WithIndent("  "))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	path := filepath.Join("testdata", "config.schema.json")

	if *updateGolden {
		if err := os.WriteFile(path, append(got, '\n'), 0o600); err != nil {
			t.Fatalf("writing golden file: %v", err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading golden file (run with -update-golden to create it): %v", err)
	}

	if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
		t.Errorf("schema drift; regenerate with -update-golden\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestFromConfig_RejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	if document.AdditionalProperties {
		t.Error("the schema must reject unknown keys, matching strict config loading")
	}

	if document.Schema != schema.Dialect {
		t.Errorf("$schema = %q, want %q", document.Schema, schema.Dialect)
	}
}

func TestFromConfig_TypesAndConstraints(t *testing.T) {
	t.Parallel()

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	tests := []struct {
		name  string
		check func(t *testing.T, p *schema.Property)
	}{
		{"log_level", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "string", "type")
			assertEq(t, len(p.Enum), 4, "enum size")
			assertEq[any](t, p.Default, "info", "default")
			assertEq(t, p.Hints.Flag, "log-level", "hint flag")
		}},
		{"timeout", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "string", "type")
			assertEq(t, p.Hints.Format, "duration", "duration hint")
		}},
		{"workers", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "integer", "type")
			assertEq(t, deref(t, p.Minimum), float64(1), "minimum")
			assertEq(t, deref(t, p.Maximum), float64(64), "maximum")
			assertEq[any](t, p.Default, int64(4), "default")
		}},
		{"ratio", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "number", "type")
			assertEq(t, deref(t, p.Minimum), float64(0), "minimum")
			assertEq[any](t, p.Default, 0.5, "default")
		}},
		{"debug", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "boolean", "type")
			assertEq[any](t, p.Default, false, "default")
		}},
		{"name", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, deref(t, p.MinLength), 3, "minLength")
			assertEq(t, deref(t, p.MaxLength), 20, "maxLength")
			assertEq(t, p.Pattern, "^[a-z]+$", "pattern")
		}},
		{"tags", func(t *testing.T, p *schema.Property) {
			t.Helper()
			assertEq(t, p.Type, "array", "type")
			assertEq(t, p.Items.Type, "string", "item type")
		}},
	}

	properties := document.Properties

	for _, test := range tests {
		property, ok := properties[test.name]
		if !ok {
			t.Errorf("property %q missing", test.name)

			continue
		}

		test.check(t, property)
	}

	if len(document.Required) != 1 || document.Required[0] != "name" {
		t.Errorf("required = %v, want [name]", document.Required)
	}
}

func TestFromConfig_NestedObject(t *testing.T) {
	t.Parallel()

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	database := document.Properties["database"]
	if database == nil || database.Type != "object" {
		t.Fatalf("database property = %+v, want an object", database)
	}

	if database.Properties["host"].Hints.Env != "DB_HOST" {
		t.Errorf("nested env hint lost: %+v", database.Properties["host"].Hints)
	}

	if !database.Properties["password"].Hints.Secret {
		t.Error("secret hint must survive nesting")
	}

	if database.Properties["port"].Maximum == nil {
		t.Error("nested validate constraint lost")
	}

	if database.AdditionalProperties == nil || *database.AdditionalProperties {
		t.Error("nested objects must reject unknown keys")
	}
}

func TestFromConfig_ArrayOfObjects(t *testing.T) {
	t.Parallel()

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	servers := document.Properties["servers"]
	if servers == nil || servers.Type != "array" {
		t.Fatalf("servers property = %+v, want an array", servers)
	}

	if servers.Items == nil || servers.Items.Type != "object" {
		t.Fatalf("servers items = %+v, want an object", servers.Items)
	}

	if len(servers.Items.Properties) != 2 {
		t.Errorf("servers item properties = %v, want host and port", servers.Items.Properties)
	}
}

func TestFromCLI_UsesTheCLIIdentity(t *testing.T) {
	t.Parallel()

	cli, err := v4.NewCLI("app", "My application", fixtureConfig{}, v4.WithFang(false))
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	document, err := schema.FromCLI(cli)
	if err != nil {
		t.Fatalf("FromCLI: %v", err)
	}

	assertEq(t, document.Title, "app", "title")
	assertEq(t, document.Description, "My application", "description")
}

func TestCommandSchemas_DescribesEveryCommand(t *testing.T) {
	t.Parallel()

	cli, err := v4.NewCLI("app", "My application", fixtureConfig{}, v4.WithFang(false))
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	deploy, err := v4.NewCommand(
		"deploy",
		struct {
			Region string `flag:"region" required:"true" help:"Region" values:"eu,us"`
		}{},
		func(_ context.Context, _ *fixtureConfig, _ struct {
			Region string `flag:"region" required:"true" help:"Region" values:"eu,us"`
		},
		) error {
			return nil
		},
		v4.WithShort("Deploy"),
	)
	if err != nil {
		t.Fatalf("NewCommand: %v", err)
	}

	if err := v4.AddCommand(cli, deploy); err != nil {
		t.Fatalf("AddCommand: %v", err)
	}

	schemas := schema.CommandSchemas(cli.Describe())

	root, ok := schemas["app"]
	if !ok {
		t.Fatalf("root schema missing, got keys %v", keysOf(schemas))
	}

	if _, ok := root.Properties["log-level"]; !ok {
		t.Error("root schema should describe the root's flags")
	}

	deploySchema, ok := schemas["app deploy"]
	if !ok {
		t.Fatalf("deploy schema missing, got keys %v", keysOf(schemas))
	}

	if len(deploySchema.Required) != 1 || deploySchema.Required[0] != "region" {
		t.Errorf("deploy required = %v, want [region]", deploySchema.Required)
	}

	if len(deploySchema.Properties["region"].Enum) != 2 {
		t.Errorf("deploy region enum = %v, want two values", deploySchema.Properties["region"].Enum)
	}
}

func TestEmit_WritesIndentedJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	if err := schema.Emit(&buf, fixtureConfig{}); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	if !strings.Contains(buf.String(), "\n  ") {
		t.Errorf("output is not indented:\n%s", buf.String())
	}

	if !strings.Contains(buf.String(), `"$schema": "`+schema.Dialect+`"`) {
		t.Errorf("output does not declare the dialect:\n%s", buf.String()[:60])
	}
}

func TestFromConfig_UnknownKeysMatchStrictLoading(t *testing.T) {
	t.Parallel()

	// The schema's property names are the canonical config keys. Anything else
	// must be rejected, which is exactly what additionalProperties: false says.
	keys, err := v4.ConfigKeys(fixtureConfig{})
	if err != nil {
		t.Fatalf("ConfigKeys: %v", err)
	}

	document, err := schema.FromConfig(fixtureConfig{})
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}

	if len(keys) != len(document.Properties) {
		t.Fatalf("schema has %d properties, config exposes %d keys", len(document.Properties), len(keys))
	}

	for _, key := range keys {
		if _, ok := document.Properties[key.Key]; !ok {
			t.Errorf("config key %q missing from the schema", key.Key)
		}
	}
}

func keysOf(schemas map[string]*schema.Document) []string {
	keys := make([]string, 0, len(schemas))

	for key := range schemas {
		keys = append(keys, key)
	}

	return keys
}

func deref[T any](t *testing.T, ptr *T) T {
	t.Helper()

	if ptr == nil {
		t.Fatalf("expected a value, got nil")
	}

	return *ptr
}

func assertEq[T comparable](t *testing.T, got, want T, label string) {
	t.Helper()

	if got != want {
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}
