// Integration test for typed example
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// newGreetCmd creates a new greet command instance for testing.
// This helper ensures consistency across tests and avoids code duplication.
func newGreetCmd() v2.Command[AppConfig, *GreetFlags] {
	return v2.Command[AppConfig, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Flags: &GreetFlags{},
		RunE: func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("%s, %s%s", flags.Prefix, flags.Name, flags.Suffix)
			if flags.Shout {
				msg = strings.ToUpper(msg)
			}

			for range flags.Count {
				fmt.Println(msg)
			}

			return nil
		},
	}
}

// captureOutput captures stdout during the execution of f and returns it as a string.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, r) // Error intentionally ignored in test helper
	os.Stdout = old

	return buf.String()
}

func TestTypedExample_CreateCLI(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cli == nil {
		t.Fatal("cli is nil")
	}

	cfg := cli.Config()
	if cfg == nil {
		t.Fatal("config is nil")
	}
}

func TestTypedExample_VersionCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "A typed CLI application", AppConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cli.SetVersion("1.0.0")

	versionCmd := v2.Command[AppConfig, v2.NoFlags]{
		Use:   "version",
		Short: "Print version information",
		RunE: func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
			fmt.Println("myapp version 1.0.0")

			return nil
		},
	}

	err = cli.AddCommand(versionCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"version"})
		_ = cli.Execute(context.Background())
	})

	if !strings.Contains(output, "myapp version 1.0.0") {
		t.Errorf("output should contain %q, got %q", "myapp version 1.0.0", output)
	}
}

func TestTypedExample_GreetCommand(t *testing.T) {
	cli, err := v2.New[AppConfig, v2.NoFlags](
		"myapp",
		"A typed CLI application",
		AppConfig{Verbose: false},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd := newGreetCmd()

	err = v2.AddAnyCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test basic greeting
	output := captureOutput(func() {
		cli.RootCommand().SetArgs([]string{"greet"})
		_ = cli.Execute(context.Background())
	})
	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("output should contain %q, got %q", "Hello, World!", output)
	}
}
