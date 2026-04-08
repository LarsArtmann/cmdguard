// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// RootConfig is the application-level configuration for tests.
type RootConfig struct {
	Debug   bool   `default:"false" flag:"debug"   help:"Enable debug mode" short:"d"`
	Verbose bool   `                flag:"verbose" help:"Verbose output"    short:"v"`
	Level   string `                flag:"level"   help:"Log level"`
}

// GreetFlags are flags for the greet command.
type GreetFlags struct {
	Name  string `default:"World" flag:"name"  help:"Name to greet"      short:"n"`
	Shout bool   `default:"false" flag:"shout" help:"Shout the greeting" short:"s"`
}

// MathFlags are flags for the math command (different type than GreetFlags).
type MathFlags struct {
	X int `default:"0" flag:"x" help:"First operand"`
	Y int `default:"0" flag:"y" help:"Second operand"`
}

// ConfigFlags are flags for the config command (yet another type).
type ConfigFlags struct {
	File string `default:""      flag:"file" help:"Config file path" short:"f"`
	JSON bool   `default:"false" flag:"json" help:"Output as JSON"`
}

// DBFlags are flags for database commands.
type DBFlags struct {
	Host     string `default:"localhost" flag:"host"     help:"Database host"`
	Port     int    `default:"5432"      flag:"port"     help:"Database port"`
	Database string `default:""          flag:"database" help:"Database name"`
}

// MigrateFlags are flags for migration commands.
type MigrateFlags struct {
	Steps     int    `default:"0"  flag:"steps"     help:"Number of migrations"`
	Direction string `default:"up" flag:"direction" help:"Migration direction"`
}

func TestV2_MixedFlagTypes_BasicCommands(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.New[RootConfig, *GreetFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		greetCalled  bool
		greetFlags   *GreetFlags
		mathCalled   bool
		mathFlags    *MathFlags
		configCalled bool
		configFlags  *ConfigFlags
	)

	err = cli.AddCommand(v2.Command[RootConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{Name: "World", Shout: false},
		RunE: func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			greetCalled = true
			greetFlags = flags

			return nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *MathFlags]{
			Use:   "math",
			Short: "Do math",
			Flags: &MathFlags{X: 0, Y: 0},
			RunE: func(_ context.Context, _ *RootConfig, flags *MathFlags) error {
				mathCalled = true
				mathFlags = flags

				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *ConfigFlags]{
			Use:   "config",
			Short: "Manage config",
			Flags: &ConfigFlags{File: "", JSON: false},
			RunE: func(_ context.Context, _ *RootConfig, flags *ConfigFlags) error {
				configCalled = true
				configFlags = flags

				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=Alice", "--shout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !greetCalled {
		t.Error("greetCalled should be true")
	}

	if greetFlags.Name != "Alice" {
		t.Errorf("greetFlags.Name = %q, want %q", greetFlags.Name, "Alice")
	}

	if !greetFlags.Shout {
		t.Error("greetFlags.Shout should be true")
	}

	err = cli.ExecuteWithArgs(ctx, []string{"math", "--x=10", "--y=20"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mathCalled {
		t.Error("mathCalled should be true")
	}

	if mathFlags.X != 10 {
		t.Errorf("mathFlags.X = %d, want %d", mathFlags.X, 10)
	}

	if mathFlags.Y != 20 {
		t.Errorf("mathFlags.Y = %d, want %d", mathFlags.Y, 20)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"config", "--file=/etc/app.yaml", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !configCalled {
		t.Error("configCalled should be true")
	}

	if configFlags.File != "/etc/app.yaml" {
		t.Errorf("configFlags.File = %q, want %q", configFlags.File, "/etc/app.yaml")
	}

	if !configFlags.JSON {
		t.Error("configFlags.JSON should be true")
	}
}

func TestV2_MixedFlagTypes_NestedSubcommands(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		migrateCalled bool
		migrateFlags  *MigrateFlags
		statusCalled  bool
		statusFlags   *DBFlags
	)

	dbCmd := v2.Command[RootConfig, *DBFlags]{
		Use:   "db",
		Short: "Database commands",
		Flags: &DBFlags{Host: "localhost", Port: 5432, Database: ""},
		Commands: []v2.Command[RootConfig, *DBFlags]{
			{
				Use:   "status",
				Short: "Check database status",
				RunE: func(_ context.Context, _ *RootConfig, flags *DBFlags) error {
					statusCalled = true
					statusFlags = flags

					return nil
				},
			},
		},
	}

	err = v2.AddAnyCommand(cli, dbCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	migrateCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:   "migrate",
		Short: "Run migrations",
		Flags: &MigrateFlags{Steps: 0, Direction: "up"},
		RunE: func(_ context.Context, _ *RootConfig, flags *MigrateFlags) error {
			migrateCalled = true
			migrateFlags = flags

			return nil
		},
	}

	err = v2.AddAnyCommand(cli, migrateCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"db", "status", "--host=prod-db", "--port=3306"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !statusCalled {
		t.Error("statusCalled should be true")
	}

	if statusFlags.Host != "prod-db" {
		t.Errorf("statusFlags.Host = %q, want %q", statusFlags.Host, "prod-db")
	}

	if statusFlags.Port != 3306 {
		t.Errorf("statusFlags.Port = %d, want %d", statusFlags.Port, 3306)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"migrate", "--steps=5", "--direction=down"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !migrateCalled {
		t.Error("migrateCalled should be true")
	}

	if migrateFlags.Steps != 5 {
		t.Errorf("migrateFlags.Steps = %d, want %d", migrateFlags.Steps, 5)
	}

	if migrateFlags.Direction != "down" {
		t.Errorf("migrateFlags.Direction = %q, want %q", migrateFlags.Direction, "down")
	}
}
