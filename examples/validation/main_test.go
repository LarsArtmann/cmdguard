//nolint:paralleltest
package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid name", "Alice", false, ""},
		{"empty name", "", true, "name cannot be empty"},
		{"single char", "A", true, "name must be at least 2 characters"},
		{"valid short name", "Bo", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)

				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateName(%q) error = %q, want %q", tt.input, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"valid count 1", 1, false},
		{"valid count 5", 5, false},
		{"valid count 10", 10, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"over limit", 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCount(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty (optional)", "", false},
		{"valid email", "alice@example.com", false},
		{"missing @", "aliceexample.com", true},
		{"missing local part", "@example.com", true},
		{"missing domain", "alice@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	t.Run("all valid", func(t *testing.T) {
		errs := ValidateFlags("Alice", 3, "alice@example.com")
		if len(errs) > 0 {
			t.Errorf("ValidateFlags() returned errors: %v", errs)
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		errs := ValidateFlags("", 0, "invalid")
		if len(errs) != 3 {
			t.Errorf("ValidateFlags() returned %d errors, want 3", len(errs))
		}
	})
}

func TestValidationExample_CLI(t *testing.T) {
	//nolint:paralleltest
	t.Run("greet with empty name fails", func(t *testing.T) {
		cmd := exec.CommandContext(
			context.Background(), "go", "run", ".", "greet", "--name=", "--count=1",
		)
		out, _ := cmd.CombinedOutput()

		if cmd.ProcessState.ExitCode() == 0 {
			t.Fatalf("expected non-zero exit, got 0; output: %s", out)
		}

		if !strings.Contains(string(out), "name cannot be empty") {
			t.Errorf("expected 'name cannot be empty' in output, got: %s", out)
		}
	})

	//nolint:paralleltest
	t.Run("greet with valid args succeeds", func(t *testing.T) {
		cmd := exec.CommandContext(
			context.Background(),
			"go",
			"run",
			".",
			"greet",
			"--name=Alice",
			"--count=1",
		)
		out, _ := cmd.CombinedOutput()

		if cmd.ProcessState.ExitCode() != 0 {
			t.Fatalf("expected zero exit, got %d; output: %s", cmd.ProcessState.ExitCode(), out)
		}

		if !strings.Contains(string(out), "Hello, Alice!") {
			t.Errorf("expected 'Hello, Alice!' in output, got: %s", out)
		}
	})

	//nolint:paralleltest
	t.Run("process without input fails", func(t *testing.T) {
		cmd := exec.CommandContext(context.Background(), "go", "run", ".", "process", "--workers=1")
		out, _ := cmd.CombinedOutput()

		if cmd.ProcessState.ExitCode() == 0 {
			t.Fatalf("expected non-zero exit, got 0; output: %s", out)
		}

		if !strings.Contains(string(out), "input") {
			t.Errorf("expected 'input' error in output, got: %s", out)
		}
	})
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "name", Message: "cannot be empty"}
	if err.Error() != "name: cannot be empty" {
		t.Errorf("ValidationError.Error() = %q, want %q", err.Error(), "name: cannot be empty")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
