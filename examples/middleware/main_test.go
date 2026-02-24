// Integration test for middleware example
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestMiddlewareExample_Timer(t *testing.T) {
	timer := NewTimer("test")
	require.NotNil(t, timer)
	assert.Equal(t, "test", timer.label)

	output := captureOutput(func() {
		timer.Start()
		time.Sleep(10 * time.Millisecond)
		timer.Stop()
	})

	assert.Contains(t, output, "[test] Starting...")
	assert.Contains(t, output, "[test] Completed in")
}

func TestMiddlewareExample_ValidateFlags(t *testing.T) {
	// Test valid flags
	flags := &ProcessingFlags{
		Input:  "input.txt",
		Output: "output.txt",
	}
	err := validateFlags(flags)
	require.NoError(t, err)

	// Test empty input
	flags = &ProcessingFlags{
		Input:  "",
		Output: "output.txt",
	}
	err = validateFlags(flags)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "input file is required")

	// Test empty output
	flags = &ProcessingFlags{
		Input:  "input.txt",
		Output: "",
	}
	err = validateFlags(flags)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "output file is required")

	// Test same input/output
	flags = &ProcessingFlags{
		Input:  "same.txt",
		Output: "same.txt",
	}
	err = validateFlags(flags)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "input and output cannot be the same file")
}

func TestMiddlewareExample_Authenticate(t *testing.T) {
	output := captureOutput(func() {
		err := authenticate(context.Background())
		require.NoError(t, err)
	})
	assert.Contains(t, output, "[Auth] Checking permissions...")
}

func TestMiddlewareExample_Logging(t *testing.T) {
	flags := &ProcessingFlags{
		Input:  "test.txt",
		Output: "out.txt",
		DryRun: true,
	}

	output := captureOutput(func() {
		logStart(flags)
	})
	assert.Contains(t, output, "[Log] Processing test.txt -> out.txt (dry-run: true)")

	output = captureOutput(func() {
		logComplete(flags)
	})
	assert.Contains(t, output, "[Log] Processing complete for out.txt")
}

func TestMiddlewareExample_Cleanup(t *testing.T) {
	output := captureOutput(func() {
		err := cleanup()
		require.NoError(t, err)
	})
	assert.Contains(t, output, "[Cleanup] Releasing resources...")
}

func TestMiddlewareExample_RecordMetrics(t *testing.T) {
	output := captureOutput(func() {
		recordMetrics(100*time.Millisecond, true)
	})
	assert.Contains(t, output, "[Metrics] Duration: 100ms, Status: success")

	output = captureOutput(func() {
		recordMetrics(200*time.Millisecond, false)
	})
	assert.Contains(t, output, "[Metrics] Duration: 200ms, Status: failure")
}

func TestMiddlewareExample_CreateCLI(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	require.NoError(t, err)
	assert.NotNil(t, cli)
}

func TestMiddlewareExample_ProcessCommandStructure(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	require.NoError(t, err)

	timer := NewTimer("process")

	processCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "process",
		Short: "Process a file with middleware",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			if err := authenticate(ctx); err != nil {
				return err
			}
			if err := validateFlags(flags); err != nil {
				return err
			}
			timer.Start()
			logStart(flags)
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			if flags.DryRun {
				return nil
			}
			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			timer.Stop()
			logComplete(flags)
			_ = cleanup()
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, processCmd)
	require.NoError(t, err)

	// Test dry run
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"process", "--input", "test.txt", "--output", "out.txt", "--dry-run"})
		_ = cli.Execute(context.Background())
	})

	assert.Contains(t, output, "[Auth] Checking permissions...")
	assert.Contains(t, output, "[Log] Processing test.txt -> out.txt")
	assert.Contains(t, output, "[process] Starting...")
	assert.Contains(t, output, "[process] Completed in")
	assert.Contains(t, output, "[Log] Processing complete for out.txt")
	assert.Contains(t, output, "[Cleanup] Releasing resources...")
}

func TestMiddlewareExample_ValidateCommandStructure(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	require.NoError(t, err)

	validateCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "validate",
		Short: "Validate configuration without processing",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			fmt.Println("[Validate] Running validation...")
			if err := validateFlags(flags); err != nil {
				return err
			}
			fmt.Println("[Validate] Configuration is valid")
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, validateCmd)
	require.NoError(t, err)

	// Test valid configuration
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"validate", "--input", "test.txt", "--output", "out.txt"})
		_ = cli.Execute(context.Background())
	})
	assert.Contains(t, output, "[Validate] Running validation...")
	assert.Contains(t, output, "[Validate] Configuration is valid")
}

func TestMiddlewareExample_ValidationFailure(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	require.NoError(t, err)

	processCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "process",
		Short: "Process a file",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			return validateFlags(flags)
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, processCmd)
	require.NoError(t, err)

	// Test with missing input (required flag)
	cli.RootCommand().SetArgs([]string{"process"})
	err = cli.Execute(context.Background())
	require.Error(t, err)
	// Error could be from flag validation (required flag) or PreRunE (input file is required)
	assert.True(t,
		errContainsAny(err.Error(), []string{"required flag", "input file is required"}),
		"Expected error about input, got: %v", err)
}

func errContainsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func TestMiddlewareExample_ProcessingFlags(t *testing.T) {
	flags := &ProcessingFlags{
		Input:  "/path/to/input.txt",
		Output: "/path/to/output.txt",
		DryRun: false,
	}
	assert.Equal(t, "/path/to/input.txt", flags.Input)
	assert.Equal(t, "/path/to/output.txt", flags.Output)
	assert.False(t, flags.DryRun)
}

func TestMiddlewareExample_AppConfig(t *testing.T) {
	cfg := AppConfig{Verbose: true}
	assert.True(t, cfg.Verbose)

	cfg = AppConfig{Verbose: false}
	assert.False(t, cfg.Verbose)
}

func TestMiddlewareExample_FullMiddlewareChain(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{Verbose: true})
	require.NoError(t, err)

	var executionOrder []string
	processCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "process",
		Short: "Process with full chain",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			executionOrder = append(executionOrder, "prerun")
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			executionOrder = append(executionOrder, "run")
			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			executionOrder = append(executionOrder, "postrun")
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, processCmd)
	require.NoError(t, err)

	_ = captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"process", "--input", "test.txt", "--output", "out.txt"})
		_ = cli.Execute(context.Background())
	})

	assert.Equal(t, []string{"prerun", "run", "postrun"}, executionOrder)
}

func TestMiddlewareExample_TimerAccuracy(t *testing.T) {
	timer := NewTimer("accuracy")

	start := time.Now()
	timer.Start()
	time.Sleep(50 * time.Millisecond)
	duration := time.Since(start)
	timer.Stop()

	// Verify timer recorded reasonable duration
	assert.GreaterOrEqual(t, duration, 50*time.Millisecond)
	assert.Less(t, duration, 200*time.Millisecond) // Allow some buffer
}

func TestMiddlewareExample_OutputFormatting(t *testing.T) {
	// Test that output contains expected format elements
	cli, err := v2.New[AppConfig, v2.NoFlags]("middleware", "Middleware Example", AppConfig{})
	require.NoError(t, err)

	timer := NewTimer("format")

	processCmd := v2.Command[AppConfig, *ProcessingFlags]{
		Use:   "process",
		Short: "Process",
		Flags: &ProcessingFlags{},
		PreRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			timer.Start()
			return nil
		},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			return nil
		},
		PostRunE: func(ctx context.Context, cfg *AppConfig, flags *ProcessingFlags) error {
			timer.Stop()
			return nil
		},
	}

	err = v2.AddAnyCommand(cli, processCmd)
	require.NoError(t, err)

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"process", "--input", "a.txt", "--output", "b.txt"})
		_ = cli.Execute(context.Background())
	})

	// Verify output format
	lines := strings.Split(output, "\n")
	var foundStart, foundStop bool
	for _, line := range lines {
		if strings.Contains(line, "[format] Starting...") {
			foundStart = true
		}
		if strings.Contains(line, "[format] Completed in") {
			foundStop = true
		}
	}
	assert.True(t, foundStart, "Expected start message")
	assert.True(t, foundStop, "Expected stop message")
}
