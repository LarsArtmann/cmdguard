// Package integration provides end-to-end tests for cmdguard v2 API.
package integration

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RootConfig is the application-level configuration for tests.
type RootConfig struct {
	Debug   bool   `flag:"debug" short:"d" default:"false" help:"Enable debug mode"`
	Verbose bool   `flag:"verbose" short:"v" help:"Verbose output"`
	Level   string `flag:"level" help:"Log level"`
}

// GreetFlags are flags for the greet command.
type GreetFlags struct {
	Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
	Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}

// MathFlags are flags for the math command (different type than GreetFlags).
type MathFlags struct {
	X int `flag:"x" default:"0" help:"First operand"`
	Y int `flag:"y" default:"0" help:"Second operand"`
}

// ConfigFlags are flags for the config command (yet another type).
type ConfigFlags struct {
	File string `flag:"file" short:"f" default:"" help:"Config file path"`
	Json bool   `flag:"json" default:"false" help:"Output as JSON"`
}

// DBFlags are flags for database commands.
type DBFlags struct {
	Host     string `flag:"host" default:"localhost" help:"Database host"`
	Port     int    `flag:"port" default:"5432" help:"Database port"`
	Database string `flag:"database" default:"" help:"Database name"`
}

// MigrateFlags are flags for migration commands.
type MigrateFlags struct {
	Steps     int    `flag:"steps" default:"0" help:"Number of migrations"`
	Direction string `flag:"direction" default:"up" help:"Migration direction"`
}

func TestV2_MixedFlagTypes_BasicCommands(t *testing.T) {
	ctx := context.Background()

	// Create CLI with RootConfig and GreetFlags
	cli, err := v2.New[RootConfig, *GreetFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

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
	require.NoError(t, err)

	// Add math command with different flag type
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *MathFlags](cli, v2.Command[RootConfig, *MathFlags]{
		Use:   "math",
		Short: "Do math",
		Flags: &MathFlags{X: 0, Y: 0},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *MathFlags) error {
			mathCalled = true
			mathFlags = flags
			return nil
		},
	})
	require.NoError(t, err)

	// Add config command with yet another flag type
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *ConfigFlags](cli, v2.Command[RootConfig, *ConfigFlags]{
		Use:   "config",
		Short: "Manage config",
		Flags: &ConfigFlags{File: "", Json: false},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *ConfigFlags) error {
			configCalled = true
			configFlags = flags
			return nil
		},
	})
	require.NoError(t, err)

	// Test greet command with flags
	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=Alice", "--shout"})
	require.NoError(t, err)
	assert.True(t, greetCalled)
	assert.Equal(t, "Alice", greetFlags.Name)
	assert.True(t, greetFlags.Shout)

	// Test math command with different flags
	err = cli.ExecuteWithArgs(ctx, []string{"math", "--x=10", "--y=20"})
	require.NoError(t, err)
	assert.True(t, mathCalled)
	assert.Equal(t, 10, mathFlags.X)
	assert.Equal(t, 20, mathFlags.Y)

	// Test config command with yet more flags
	err = cli.ExecuteWithArgs(ctx, []string{"config", "--file=/etc/app.yaml", "--json"})
	require.NoError(t, err)
	assert.True(t, configCalled)
	assert.Equal(t, "/etc/app.yaml", configFlags.File)
	assert.True(t, configFlags.Json)
}

func TestV2_MixedFlagTypes_NestedSubcommands(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

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
	require.NoError(t, err)

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
	require.NoError(t, err)

	// Test db status with DBFlags
	err = cli.ExecuteWithArgs(ctx, []string{"db", "status", "--host=prod-db", "--port=3306"})
	require.NoError(t, err)
	assert.True(t, statusCalled)
	assert.Equal(t, "prod-db", statusFlags.Host)
	assert.Equal(t, 3306, statusFlags.Port)

	// Test migrate with MigrateFlags
	err = cli.ExecuteWithArgs(ctx, []string{"migrate", "--steps=5", "--direction=down"})
	require.NoError(t, err)
	assert.True(t, migrateCalled)
	assert.Equal(t, 5, migrateFlags.Steps)
	assert.Equal(t, "down", migrateFlags.Direction)
}

func TestV2_MixedFlagTypes_NoInterference(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, *GreetFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

	var lastExecuted string
	var lastFlags any

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
	require.NoError(t, err)

	// Add command B with MathFlags
	err = v2.AddAnyCommand[RootConfig, *GreetFlags, *MathFlags](cli, v2.Command[RootConfig, *MathFlags]{
		Use:   "cmd-b",
		Short: "Command B",
		Flags: &MathFlags{X: 0, Y: 0},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *MathFlags) error {
			lastExecuted = "B"
			lastFlags = flags
			return nil
		},
	})
	require.NoError(t, err)

	// Execute A multiple times
	for i := 0; i < 3; i++ {
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-a", "--name=test"})
		require.NoError(t, err)
		assert.Equal(t, "A", lastExecuted)
		gf, ok := lastFlags.(*GreetFlags)
		assert.True(t, ok)
		assert.Equal(t, "test", gf.Name)
	}

	// Execute B multiple times
	for i := 0; i < 3; i++ {
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-b", "--x=42"})
		require.NoError(t, err)
		assert.Equal(t, "B", lastExecuted)
		mf, ok := lastFlags.(*MathFlags)
		assert.True(t, ok)
		assert.Equal(t, 42, mf.X)
	}

	// Interleave executions
	for i := 0; i < 5; i++ {
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-a", "--shout"})
		require.NoError(t, err)
		assert.Equal(t, "A", lastExecuted)

		err = cli.ExecuteWithArgs(ctx, []string{"cmd-b", "--y=99"})
		require.NoError(t, err)
		assert.Equal(t, "B", lastExecuted)
	}
}

func TestV2_MixedFlagTypes_WithNoFlags(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

	var executed bool

	// Add command with NoFlags
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, v2.NoFlags](cli, v2.Command[RootConfig, v2.NoFlags]{
		Use:   "simple",
		Short: "Simple command",
		RunE: func(ctx context.Context, cfg *RootConfig, flags v2.NoFlags) error {
			executed = true
			return nil
		},
	})
	require.NoError(t, err)

	// Add command with actual flags
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](cli, v2.Command[RootConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet command",
		Flags: &GreetFlags{Name: "World", Shout: false},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
			executed = true
			return nil
		},
	})
	require.NoError(t, err)

	// Test simple command (NoFlags)
	executed = false
	err = cli.ExecuteWithArgs(ctx, []string{"simple"})
	require.NoError(t, err)
	assert.True(t, executed)

	// Test greet command with flags
	executed = false
	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=Bob"})
	require.NoError(t, err)
	assert.True(t, executed)
}

func TestV2_MixedFlagTypes_WithLifecycleHooks(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

	var (
		preRunCalled  bool
		runCalled     bool
		postRunCalled bool
		receivedFlags *GreetFlags
	)

	// Add command with lifecycle hooks
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](cli, v2.Command[RootConfig, *GreetFlags]{
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
	})
	require.NoError(t, err)

	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=TestUser", "--shout"})
	require.NoError(t, err)

	assert.True(t, preRunCalled)
	assert.True(t, runCalled)
	assert.True(t, postRunCalled)
	assert.Equal(t, "TestUser", receivedFlags.Name)
	assert.True(t, receivedFlags.Shout)
}

func TestV2_MixedFlagTypes_ValidationErrors(t *testing.T) {
	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

	// Command without Use should fail
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](cli, v2.Command[RootConfig, *GreetFlags]{
		Use:   "",
		Short: "Invalid command",
		RunE:  func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error { return nil },
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no Use field")

	// Command without RunE and no subcommands should fail
	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](cli, v2.Command[RootConfig, *GreetFlags]{
		Use:   "invalid",
		Short: "No handler",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no RunE")
}

func TestV2_MixedFlagTypes_ConfigAccess(t *testing.T) {
	ctx := context.Background()

	defaultConfig := RootConfig{
		Debug:   false,
		Verbose: true,
		Level:   "debug",
	}

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", defaultConfig)
	require.NoError(t, err)

	var receivedConfig *RootConfig

	err = v2.AddAnyCommand[RootConfig, v2.NoFlags, *GreetFlags](cli, v2.Command[RootConfig, *GreetFlags]{
		Use:   "check",
		Short: "Check config access",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *RootConfig, flags *GreetFlags) error {
			receivedConfig = cfg
			return nil
		},
	})
	require.NoError(t, err)

	err = cli.ExecuteWithArgs(ctx, []string{"check"})
	require.NoError(t, err)

	require.NotNil(t, receivedConfig)
	assert.True(t, receivedConfig.Verbose)
	assert.Equal(t, "debug", receivedConfig.Level)
}

func TestV2_MixedFlagTypes_DeeplyNested(t *testing.T) {
	ctx := context.Background()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	require.NoError(t, err)

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
	require.NoError(t, err)

	err = cli.ExecuteWithArgs(ctx, []string{"migrate", "up", "--steps=3"})
	require.NoError(t, err)
	assert.Equal(t, 3, executedFlags.Steps)
}
