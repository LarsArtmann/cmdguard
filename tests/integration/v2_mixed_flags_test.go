// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	"github.com/spf13/cobra"

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
	ctx := context.Background()

	// Create CLI with RootConfig and GreetFlags
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

	// Add greet command with same flag type as CLI root
	err = cli.AddCommand(v2.Command[RootConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{Name: "World", Shout: false},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
			greetCalled = true
			greetFlags = flags
			return nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add math command with different flag type
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *MathFlags](
		cli,
		v2.Command[RootConfig, *MathFlags]{
			Use:   "math",
			Short: "Do math",
			Flags: &MathFlags{X: 0, Y: 0},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *MathFlags) error {
				mathCalled = true
				mathFlags = flags
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add config command with yet another flag type
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *ConfigFlags](
		cli,
		v2.Command[RootConfig, *ConfigFlags]{
			Use:   "config",
			Short: "Manage config",
			Flags: &ConfigFlags{File: "", JSON: false},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *ConfigFlags) error {
				configCalled = true
				configFlags = flags
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test greet command with flags
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

	// Test math command with different flags
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

	// Test config command with yet more flags
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
	ctx := context.Background()

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

	// Create db parent command with DBFlags
	dbCmd := v2.Command[RootConfig, *DBFlags]{
		Use:   "db",
		Short: "Database commands",
		Flags: &DBFlags{Host: "localhost", Port: 5432, Database: ""},
		Commands: []v2.Command[RootConfig, *DBFlags]{
			{
				Use:   "status",
				Short: "Check database status",
				RunE: func(ctx context.Context, cfg *RootConfig, flags *DBFlags) error {
					statusCalled = true
					statusFlags = flags
					return nil
				},
			},
		},
	}

	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *DBFlags](cli, dbCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add migrate command with different flag type at same level as db
	migrateCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:   "migrate",
		Short: "Run migrations",
		Flags: &MigrateFlags{Steps: 0, Direction: "up"},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *MigrateFlags) error {
			migrateCalled = true
			migrateFlags = flags
			return nil
		},
	}

	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *MigrateFlags](cli, migrateCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test db status with DBFlags
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

	// Test migrate with MigrateFlags
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

// assertCommandExecution runs a command multiple times and verifies the execution state.
func assertCommandExecution[
	T any,
	F any,
](
	t *testing.T,
	ctx context.Context,
	cli *v2.GuardedCommand[T, F],
	args []string,
	wantExecuted string,
	assertFlags func(t *testing.T, flags any),
) {
	t.Helper()

	for range 3 {
		err := cli.ExecuteWithArgs(ctx, args)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i+1, err)
		}
		if wantExecuted != lastExecuted {
			t.Errorf("iteration %d: lastExecuted = %q, want %q", i, lastExecuted, wantExecuted)
		}
		assertFlags(t, lastFlags)
	}
}

var (
	lastExecuted string
	lastFlags    any
)

func TestV2_MixedFlagTypes_NoInterference(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, *GreetFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add command A with GreetFlags
	err = cli.AddCommand(v2.Command[RootConfig, *GreetFlags]{
		Use:   "cmd-a",
		Short: "Command A",
		Flags: &GreetFlags{Name: "default", Shout: false},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
			lastExecuted = "A"
			lastFlags = flags
			return nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add command B with MathFlags
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *MathFlags](
		cli,
		v2.Command[RootConfig, *MathFlags]{
			Use:   "cmd-b",
			Short: "Command B",
			Flags: &MathFlags{X: 0, Y: 0},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *MathFlags) error {
				lastExecuted = "B"
				lastFlags = flags
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Execute A multiple times
	assertCommandExecution(
		t,
		ctx,
		cli,
		[]string{"cmd-a", "--name=test"},
		"A",
		func(t *testing.T, flags any) {
			gf, ok := flags.(*GreetFlags)
			if !ok {
				t.Fatalf("expected *GreetFlags, got %T", flags)
			}
			if gf.Name != "test" {
				t.Errorf("gf.Name = %q, want %q", gf.Name, "test")
			}
		},
	)

	// Execute B multiple times
	assertCommandExecution(
		t,
		ctx,
		cli,
		[]string{"cmd-b", "--x=42"},
		"B",
		func(t *testing.T, flags any) {
			mf, ok := flags.(*MathFlags)
			if !ok {
				t.Fatalf("expected *MathFlags, got %T", flags)
			}
			if mf.X != 42 {
				t.Errorf("mf.X = %d, want %d", mf.X, 42)
			}
		},
	)

	// Interleave executions
	for range 5 {
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-a", "--shout"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lastExecuted != "A" {
			t.Errorf("lastExecuted = %q, want %q", lastExecuted, "A")
		}

		err = cli.ExecuteWithArgs(ctx, []string{"cmd-b", "--y=99"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lastExecuted != "B" {
			t.Errorf("lastExecuted = %q, want %q", lastExecuted, "B")
		}
	}
}

func TestV2_MixedFlagTypes_WithNoFlags(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executed bool

	// Add command with NoFlags
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, v2.NoFlags](
		cli,
		v2.Command[RootConfig, v2.NoFlags]{
			Use:   "simple",
			Short: "Simple command",
			RunE: func(ctx context.Context, cfg *RootConfig, flags v2.NoFlags) error {
				executed = true
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Add command with actual flags
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](
		cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "greet",
			Short: "Greet command",
			Flags: &GreetFlags{Name: "World", Shout: false},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
				executed = true
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test simple command (NoFlags)
	executed = false
	err = cli.ExecuteWithArgs(ctx, []string{"simple"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !executed {
		t.Error("executed should be true")
	}

	// Test greet command with flags
	executed = false
	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !executed {
		t.Error("executed should be true")
	}
}

func TestV2_MixedFlagTypes_WithLifecycleHooks(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		preRunCalled  bool
		runCalled     bool
		postRunCalled bool
		receivedFlags *GreetFlags
	)

	// Add command with lifecycle hooks
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](
		cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "greet",
			Short: "Greet with lifecycle",
			Flags: &GreetFlags{Name: "World", Shout: false},
			PreRunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
				preRunCalled = true
				return nil
			},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
				runCalled = true
				receivedFlags = flags
				return nil
			},
			PostRunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
				postRunCalled = true
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=TestUser", "--shout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !preRunCalled {
		t.Error("preRunCalled should be true")
	}
	if !runCalled {
		t.Error("runCalled should be true")
	}
	if !postRunCalled {
		t.Error("postRunCalled should be true")
	}
	if receivedFlags.Name != "TestUser" {
		t.Errorf("receivedFlags.Name = %q, want %q", receivedFlags.Name, "TestUser")
	}
	if !receivedFlags.Shout {
		t.Error("receivedFlags.Shout should be true")
	}
}

func TestV2_MixedFlagTypes_ValidationErrors(t *testing.T) {
	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Command without Use should fail
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](
		cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "",
			Short: "Invalid command",
			RunE:  func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error { return nil },
		},
	)
	if err == nil {
		t.Error("expected error for empty Use field")
	}

	// Command without RunE and no subcommands should fail
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](
		cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "invalid",
			Short: "No handler",
		},
	)
	if err == nil {
		t.Error("expected error for missing RunE")
	}
}

func TestV2_MixedFlagTypes_ConfigAccess(t *testing.T) {
	ctx := context.Background()

	defaultConfig := RootConfig{
		Debug:   false,
		Verbose: true,
		Level:   "debug",
	}

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", defaultConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var receivedConfig *RootConfig

	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](
		cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "check",
			Short: "Check config access",
			Flags: &GreetFlags{},
			RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
				receivedConfig = cfg
				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"check"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedConfig == nil {
		t.Fatal("receivedConfig is nil")
	}
	if !receivedConfig.Verbose {
		t.Error("receivedConfig.Verbose should be true")
	}
	if receivedConfig.Level != "debug" {
		t.Errorf("receivedConfig.Level = %q, want %q", receivedConfig.Level, "debug")
	}
}

func TestV2_MixedFlagTypes_DeeplyNested(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executedFlags *MigrateFlags

	// Create a deep command tree: migrate -> up
	// This tests that nested commands work with AddAnyCommand
	migrateUpCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:   "up",
		Short: "Run up migrations",
		RunE: func(ctx context.Context, cfg *RootConfig, flags *MigrateFlags) error {
			executedFlags = flags
			return nil
		},
	}

	migrateCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:      "migrate",
		Short:    "Migration commands",
		Flags:    &MigrateFlags{Steps: 0, Direction: "up"},
		Commands: []v2.Command[RootConfig, *MigrateFlags]{migrateUpCmd},
	}

	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *MigrateFlags](cli, migrateCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"migrate", "up", "--steps=3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if executedFlags.Steps != 3 {
		t.Errorf("executedFlags.Steps = %d, want %d", executedFlags.Steps, 3)
	}
}
