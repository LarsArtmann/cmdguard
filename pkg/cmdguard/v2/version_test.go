package v2

import (
	"bytes"
	"testing"
)

func TestVersionCommandSuccess(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI[Config](
		"myapp", "My app", Config{},
		WithCLIVersion("2.0.0"),
	)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err := VersionCommand[Config](cli)
	if err != nil {
		t.Fatal(err)
	}

	if cmd.Use() != "version" {
		t.Errorf("got Use %q, want %q", cmd.Use(), "version")
	}
}

func TestGenerateVersionCommand(t *testing.T) {
	t.Parallel()

	t.Run("writes version to custom writer", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[Config](
			"myapp", "My app", Config{},
			WithCLIVersion("3.0.0"),
		)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer

		cmd, err := GenerateVersionCommand[Config](cli, &buf)
		if err != nil {
			t.Fatal(err)
		}

		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := "myapp 3.0.0\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns error when no version set", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[Config]("myapp", "My app", Config{})
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer

		_, err = GenerateVersionCommand[Config](cli, &buf)
		if err == nil {
			t.Fatal("expected error when no version set")
		}
	})
}
