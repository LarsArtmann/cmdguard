// Package main demonstrates advanced flag usage with cmdguard v2.
//
// This example shows:
// - Custom flag types (Enum, Duration, LogLevel)
// - Flag validation
// - Required flags
// - Flag typo suggestions
// - Complex nested flags
//
// Usage:
//
//	go run examples/advanced-flags/main.go
//	go run examples/advanced-flags/main.go server --port 8080
//	go run examples/advanced-flags/main.go server --log-level debug
//	go run examples/advanced-flags/main.go config --required-flag value
package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	examplesinternal "github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

const (
	defaultEnvironment = "production"
	formatJSON         = "json"
	formatYAML         = "yaml"
	formatYML          = "yml"
)

// GlobalConfig contains application-wide settings.
type GlobalConfig struct {
	LogLevel  v2.LogLevel `default:"info" flag:"log-level"  help:"Log level (debug, info, warn, error)" short:"l"`
	LogFormat string      `default:"text" flag:"log-format" help:"Log format (text, json)"`
	Timeout   v2.Duration `default:"30s"  flag:"timeout"    help:"Default timeout"`
}

// ServerFlags for the server command.
type ServerFlags struct {
	Port     int         `default:"8080"      flag:"port"      help:"Server port"        short:"p"`
	Host     string      `default:"localhost" flag:"host"      help:"Server host"`
	LogLevel v2.LogLevel `default:""          flag:"log-level" help:"Override log level"`
}

// ConfigFlags for the config command.
type ConfigFlags struct {
	ConfigFile   string `default:""     flag:"config"        help:"Config file path"           short:"c"`
	RequiredFlag string `               flag:"required-flag" help:"Required flag demo"                   required:"true"`
	OutputFormat string `default:"yaml" flag:"output"        help:"Output format (json, yaml)"`
}

// EnumFlags demonstrates enum-like validation.
type EnumFlags struct {
	Environment string `default:"development" flag:"env"    help:"Environment (development, staging, production)"`
	Region      string `default:"us-east-1"   flag:"region" help:"AWS region"`
}

// Validate implements custom validation for EnumFlags.
func (f EnumFlags) Validate() error {
	validEnvs := []string{"development", "staging", defaultEnvironment}
	if !slices.Contains(validEnvs, f.Environment) {
		return fmt.Errorf("invalid environment: %q (must be one of: %v)", f.Environment, validEnvs)
	}

	return nil
}

// DurationFlags demonstrates duration parsing.
type DurationFlags struct {
	SessionTimeout v2.Duration `default:"1h" flag:"session-timeout" help:"Session timeout duration"`
	RetryDelay     v2.Duration `default:"5s" flag:"retry-delay"     help:"Retry delay between attempts"`
	MaxWait        v2.Duration `default:"0s" flag:"max-wait"        help:"Maximum wait time (0 = no limit)"`
}

func main() {
	ctx := context.Background()

	root, err := v2.NewCLI[GlobalConfig](
		"advflags",
		"Advanced Flags Example",
		GlobalConfig{},
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// Server command with custom flags
	serverCmd, err := v2.NewCommand[GlobalConfig, ServerFlags](
		"server",
		func(ctx context.Context, cfg *GlobalConfig, flags ServerFlags) error {
			logLevel := cfg.LogLevel
			if flags.LogLevel.String() != "" {
				logLevel = flags.LogLevel
			}

			fmt.Printf("Starting server on %s:%d\n", flags.Host, flags.Port)
			fmt.Printf("Log level: %s\n", logLevel)
			fmt.Printf("Timeout: %s\n", cfg.Timeout)

			return nil
		},
		v2.WithShort[GlobalConfig, ServerFlags]("Start the server"),
		v2.WithLong[GlobalConfig, ServerFlags](
			"Start the HTTP server with configurable host and port.",
		),
		v2.WithFlags[GlobalConfig, ServerFlags](ServerFlags{}),
	)
	if err != nil {
		examplesinternal.Fatalf("Error creating server command: %v\n", err)
	}

	if err := v2.AddCommand(root, serverCmd); err != nil {
		examplesinternal.Fatalf("Error adding server command: %v\n", err)
	}

	// Config command with required flag
	configCmd, err := v2.NewCommand[GlobalConfig, ConfigFlags](
		"config",
		func(ctx context.Context, cfg *GlobalConfig, flags ConfigFlags) error {
			fmt.Printf("Config file: %s\n", flags.ConfigFile)
			fmt.Printf("Required flag: %s\n", flags.RequiredFlag)
			fmt.Printf("Output format: %s\n", flags.OutputFormat)

			return nil
		},
		v2.WithShort[GlobalConfig, ConfigFlags]("Configuration management"),
		v2.WithFlags[GlobalConfig, ConfigFlags](ConfigFlags{}),
		v2.WithPreRunE[GlobalConfig, ConfigFlags](
			func(ctx context.Context, cfg *GlobalConfig, flags ConfigFlags) error {
				if flags.OutputFormat != "" {
					validFormats := []string{formatJSON, formatYAML, formatYML}
					if !slices.Contains(validFormats, flags.OutputFormat) {
						return fmt.Errorf("invalid output format: %q (suggestions: %s)",
							flags.OutputFormat,
							suggestFormat(flags.OutputFormat))
					}
				}

				return nil
			},
		),
	)
	if err != nil {
		examplesinternal.Fatalf("Error creating config command: %v\n", err)
	}

	if err := v2.AddCommand(root, configCmd); err != nil {
		examplesinternal.Fatalf("Error adding config command: %v\n", err)
	}

	// Enum demo command
	envCmd, err := v2.NewCommand[GlobalConfig, EnumFlags](
		"env",
		func(ctx context.Context, cfg *GlobalConfig, flags EnumFlags) error {
			fmt.Printf("Environment: %s\n", flags.Environment)
			fmt.Printf("Region: %s\n", flags.Region)

			if flags.Environment == defaultEnvironment {
				fmt.Println("⚠️  Running in PRODUCTION mode!")
			}

			return nil
		},
		v2.WithShort[GlobalConfig, EnumFlags]("Environment settings"),
		v2.WithFlags[GlobalConfig, EnumFlags](EnumFlags{}),
		v2.WithPreRunE[GlobalConfig, EnumFlags](
			func(ctx context.Context, cfg *GlobalConfig, flags EnumFlags) error {
				return flags.Validate()
			},
		),
	)
	if err != nil {
		examplesinternal.Fatalf("Error creating env command: %v\n", err)
	}

	if err := v2.AddCommand(root, envCmd); err != nil {
		examplesinternal.Fatalf("Error adding env command: %v\n", err)
	}

	// Duration demo command
	durationCmd, err := v2.NewCommand[GlobalConfig, DurationFlags]("duration",
		func(ctx context.Context, cfg *GlobalConfig, flags DurationFlags) error {
			fmt.Printf("Session timeout: %s (%g seconds)\n",
				flags.SessionTimeout,
				flags.SessionTimeout.Seconds())
			fmt.Printf("Retry delay: %s (%d milliseconds)\n",
				flags.RetryDelay,
				flags.RetryDelay.Milliseconds())
			fmt.Printf("Max wait: %s\n", flags.MaxWait)

			return nil
		},
		v2.WithShort[GlobalConfig, DurationFlags]("Duration settings demo"),
		v2.WithFlags[GlobalConfig, DurationFlags](DurationFlags{}),
	)
	if err != nil {
		examplesinternal.Fatalf("Error creating duration command: %v\n", err)
	}

	if err := v2.AddCommand(root, durationCmd); err != nil {
		examplesinternal.Fatalf("Error adding duration command: %v\n", err)
	}

	// Execute
	examplesinternal.Execute(ctx, root)
}

// suggestFormat returns a suggestion for invalid format.
func suggestFormat(input string) string {
	validFormats := []string{formatJSON, formatYAML, formatYML}

	var bestMatch string

	bestDistance := len(input)

	for _, format := range validFormats {
		dist := levenshteinDistance(input, format)
		if dist < bestDistance {
			bestDistance = dist
			bestMatch = format
		}
	}

	if bestMatch != "" {
		return bestMatch
	}

	return strings.Join(validFormats, ", ")
}

// levenshteinDistance calculates the edit distance between two strings.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}

	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			deletion := matrix[i-1][j] + 1
			insertion := matrix[i][j-1] + 1
			substitution := matrix[i-1][j-1] + cost

			matrix[i][j] = min(deletion, min(insertion, substitution))
		}
	}

	return matrix[len(s1)][len(s2)]
}
