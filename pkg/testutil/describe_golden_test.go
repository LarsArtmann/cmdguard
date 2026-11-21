package testutil

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildDescribeFixture compiles the testdata fixture CLI and returns the
// binary path, so the capture test exercises a real process boundary.
func buildDescribeFixture(t *testing.T) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "fixtureapp")
	build := exec.CommandContext(t.Context(), "go", "build", "-o", bin, "./testdata/describeapp")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("building fixture: %v\n%s", err, out)
	}

	return bin
}

func TestDescribeContractFromBinary(t *testing.T) {
	t.Parallel()

	bin := buildDescribeFixture(t)

	contract := DescribeContractFromBinary(t, bin)

	var doc struct {
		SchemaVersion string `json:"schema_version"`
		Name          string `json:"name"`
		Flags         []struct {
			Name string `json:"name"`
			Env  string `json:"env"`
		} `json:"flags"`
		Commands []struct {
			Name  string `json:"name"`
			Flags []struct {
				Name string `json:"name"`
				Env  string `json:"env"`
			} `json:"flags"`
		} `json:"commands"`
	}
	if err := json.Unmarshal([]byte(contract), &doc); err != nil {
		t.Fatalf("describe output is not valid JSON: %v\n%s", err, contract)
	}

	if doc.SchemaVersion != "1" {
		t.Errorf("schema_version = %q, want %q", doc.SchemaVersion, "1")
	}

	if doc.Name != "fixtureapp" {
		t.Errorf("name = %q, want fixtureapp", doc.Name)
	}

	greetIdx := -1
	nameIdx := -1

	for i := range doc.Commands {
		if doc.Commands[i].Name == "greet" {
			greetIdx = i

			for j, f := range doc.Commands[i].Flags {
				if f.Name == "name" {
					nameIdx = j
				}
			}
		}
	}

	if greetIdx < 0 {
		t.Fatalf("describe output missing greet command:\n%s", contract)
	}

	if nameIdx < 0 {
		t.Errorf("greet flags = %+v, want a name flag", doc.Commands[greetIdx].Flags)
	}

	envSeen := false

	for _, f := range doc.Flags {
		if f.Env == "FIXTURE_APP_VERBOSE" {
			envSeen = true
		}
	}

	if !envSeen {
		t.Errorf("root flags missing FIXTURE_APP_VERBOSE env name (env prefix not applied):\n%s", contract)
	}
}
