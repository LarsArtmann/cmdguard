// Integration test for advanced-flags example
package main

import (
	"context"
	"testing"
	"time"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestAdvancedFlags_CreateCLI(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[GlobalConfig](
		"advflags",
		"Advanced Flags Example",
		GlobalConfig{},
	)
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	if cli == nil {
		t.Fatal("CLI is nil")
	}

	cfg := cli.Config()
	if cfg == nil {
		t.Fatal("Config is nil")
	}
}

func TestAdvancedFlags_ServerCommand(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[GlobalConfig](
		"advflags",
		"Advanced Flags Example",
		GlobalConfig{},
	)
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	serverCmd, err := v2.NewCommand[GlobalConfig, ServerFlags](
		"server",
		func(ctx context.Context, cfg *GlobalConfig, flags ServerFlags) error {
			return nil
		},
		v2.WithShort[GlobalConfig, ServerFlags]("Start the server"),
		v2.WithFlags[GlobalConfig, ServerFlags](ServerFlags{}),
	)
	if err != nil {
		t.Fatalf("Failed to create server command: %v", err)
	}

	err = v2.AddCommand(cli, serverCmd)
	if err != nil {
		t.Fatalf("Failed to add server command: %v", err)
	}

	// Verify command was added
	cmd := cli.RootCommand()
	if cmd == nil {
		t.Fatal("Root command is nil")
	}
}

func TestAdvancedFlags_ConfigCommand(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[GlobalConfig](
		"advflags",
		"Advanced Flags Example",
		GlobalConfig{},
	)
	if err != nil {
		t.Fatalf("Failed to create CLI: %v", err)
	}

	configCmd, err := v2.NewCommand[GlobalConfig, ConfigFlags](
		"config",
		func(ctx context.Context, cfg *GlobalConfig, flags ConfigFlags) error {
			return nil
		},
		v2.WithShort[GlobalConfig, ConfigFlags]("Configuration management"),
		v2.WithFlags[GlobalConfig, ConfigFlags](ConfigFlags{}),
	)
	if err != nil {
		t.Fatalf("Failed to create config command: %v", err)
	}

	err = v2.AddCommand(cli, configCmd)
	if err != nil {
		t.Fatalf("Failed to add config command: %v", err)
	}
}

func TestAdvancedFlags_EnumValidation(t *testing.T) {
	t.Parallel()
	// Test valid environment
	flags := EnumFlags{Environment: defaultEnvironment, Region: "us-west-2"}

	err := flags.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid environment, got: %v", err)
	}

	// Test invalid environment
	flags = EnumFlags{Environment: "invalid", Region: "us-west-2"}

	err = flags.Validate()
	if err == nil {
		t.Error("Expected error for invalid environment, got nil")
	}

	if err != nil && err.Error() == "" {
		t.Error("Expected error message for invalid environment")
	}
}

func TestAdvancedFlags_DurationParsing(t *testing.T) {
	t.Parallel()
	// Test valid duration parsing
	duration, err := v2.ParseDuration("1h30m")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}

	// Duration().Hours() returns float64
	hours := duration.Duration().Hours()
	if hours < 1.4 || hours > 1.6 {
		t.Errorf("Expected ~1.5 hours, got %f", hours)
	}

	// Test invalid duration
	_, err = v2.ParseDuration("not-a-duration")
	if err == nil {
		t.Error("Expected error for invalid duration, got nil")
	}
}

func TestAdvancedFlags_FormatSuggestion(t *testing.T) {
	t.Parallel()
	// Test suggestion for typo
	suggestion := suggestFormat("yam")
	if suggestion != formatYAML {
		t.Errorf("Expected suggestion 'yaml' for 'yam', got: %s", suggestion)
	}

	suggestion = suggestFormat("jsn")
	if suggestion != formatJSON {
		t.Errorf("Expected suggestion 'json' for 'jsn', got: %s", suggestion)
	}

	// Test when no close match
	suggestion = suggestFormat("xyz")
	if suggestion == "" {
		t.Error("Expected non-empty suggestion even for no close match")
	}
}

func TestGlobalConfig_LogLevel(t *testing.T) {
	t.Parallel()
	// Test that LogLevel type works correctly using ParseLogLevel
	logLevel, err := v2.ParseLogLevel("debug")
	if err != nil {
		t.Fatalf("Failed to parse log level: %v", err)
	}

	cfg := GlobalConfig{
		LogLevel:  logLevel,
		LogFormat: "json",
		Timeout:   v2.FromDuration(30 * time.Second),
	}

	if cfg.LogLevel.String() != "debug" {
		t.Errorf("Expected log level 'debug', got: %s", cfg.LogLevel.String())
	}

	if cfg.LogFormat != "json" {
		t.Errorf("Expected log format 'json', got: %s", cfg.LogFormat)
	}
}
