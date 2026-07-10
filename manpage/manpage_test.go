package manpage

import (
	"strings"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type testConfig struct {
	Name string `flag:"name" short:"n" default:"" help:"Your name"`
}

func TestGenerate_NonEmpty(t *testing.T) {
	t.Parallel()

	cli, err := v3.NewCLI[testConfig]("testapp", "Test App", testConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	content, err := Generate(cli, 1)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generate returned empty content")
	}

	if !strings.Contains(content, "testapp") {
		t.Fatalf("Generate should contain app name 'testapp', got:\n%s", content)
	}
}

func TestWrite_WritesToWriter(t *testing.T) {
	t.Parallel()

	cli, err := v3.NewCLI[testConfig]("testapp", "Test App", testConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	var buf strings.Builder

	err = Write(cli, &buf, 1)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("Write produced no output")
	}
}

func TestGenerateCommand_ReturnsCommand(t *testing.T) {
	t.Parallel()

	cli, err := v3.NewCLI[testConfig]("myapp", "My Application", testConfig{})
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd, err := GenerateCommand(cli)
	if err != nil {
		t.Fatalf("GenerateCommand failed: %v", err)
	}

	if cmd == nil {
		t.Fatal("GenerateCommand returned nil command")
	}

	if cmd.RunE == nil && cmd.Run == nil {
		t.Fatal("GenerateCommand should return a runnable command")
	}
}
