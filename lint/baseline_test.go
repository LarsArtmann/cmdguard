package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/go-finding"
)

func findingAt(rule, file string, line int) finding.Finding {
	return finding.Finding{
		Rule: finding.RuleName(rule),
		Position: finding.Position{
			File: finding.FilePath(file),
			Line: line,
		},
	}
}

func TestApplyBaseline(t *testing.T) {
	t.Parallel()

	t.Run("nil baseline marks everything fresh", func(t *testing.T) {
		t.Parallel()

		findings := []finding.Finding{findingAt("CG004", "a.go", 10)}

		fresh, stale := ApplyBaseline(findings, nil)

		if len(fresh) != 1 || len(stale) != 0 {
			t.Errorf("fresh=%d stale=%d, want 1/0", len(fresh), len(stale))
		}
	})

	t.Run("exact match is consumed", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG004", File: "a.go", Line: 10},
		}}
		findings := []finding.Finding{findingAt("CG004", "a.go", 10)}

		fresh, stale := ApplyBaseline(findings, baseline)

		if len(fresh) != 0 {
			t.Errorf("fresh=%d, want 0 (finding is baselined)", len(fresh))
		}

		if len(stale) != 0 {
			t.Errorf("stale=%d, want 0 (entry was consumed)", len(stale))
		}
	})

	t.Run("line drift within same rule+file is absorbed", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG006", File: "a.go", Line: 10},
		}}
		findings := []finding.Finding{findingAt("CG006", "a.go", 42)}

		fresh, _ := ApplyBaseline(findings, baseline)

		if len(fresh) != 0 {
			t.Errorf("fresh=%d, want 0 (edit above the finding shifted its line)", len(fresh))
		}
	})

	t.Run("surplus findings for a rule+file are fresh", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG006", File: "a.go", Line: 10},
			{Rule: "CG006", File: "a.go", Line: 20},
		}}
		findings := []finding.Finding{
			findingAt("CG006", "a.go", 10),
			findingAt("CG006", "a.go", 20),
			findingAt("CG006", "a.go", 30),
		}

		fresh, _ := ApplyBaseline(findings, baseline)

		if len(fresh) != 1 {
			t.Errorf("fresh=%d, want 1 (the new third finding)", len(fresh))
		}

		if fresh[0].Position.Line != 30 {
			t.Errorf("fresh finding line=%d, want 30", fresh[0].Position.Line)
		}
	})

	t.Run("rules and files do not cross-match", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG004", File: "a.go", Line: 10},
			{Rule: "CG006", File: "b.go", Line: 10},
		}}
		findings := []finding.Finding{
			findingAt("CG006", "a.go", 10),
			findingAt("CG004", "b.go", 10),
		}

		fresh, _ := ApplyBaseline(findings, baseline)

		if len(fresh) != 2 {
			t.Errorf("fresh=%d, want 2 (neither rule nor file matched)", len(fresh))
		}
	})

	t.Run("same-named files in different dirs do not merge", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG006", File: "cmd/util.go", Line: 5},
		}}
		findings := []finding.Finding{findingAt("CG006", "internal/util.go", 5)}

		fresh, _ := ApplyBaseline(findings, baseline)

		if len(fresh) != 1 {
			t.Errorf("fresh=%d, want 1 (file paths compare exactly)", len(fresh))
		}
	})

	t.Run("fixed findings surface as stale entries", func(t *testing.T) {
		t.Parallel()

		baseline := &Baseline{Version: baselineVersion, Entries: []BaselineEntry{
			{Rule: "CG004", File: "a.go", Line: 10},
			{Rule: "CG006", File: "b.go", Line: 20},
		}}
		findings := []finding.Finding{findingAt("CG006", "b.go", 20)}

		_, stale := ApplyBaseline(findings, baseline)

		if len(stale) != 1 || stale[0].Rule != "CG004" || stale[0].File != "a.go" {
			t.Errorf("stale=%v, want the unconsumed CG004@a.go entry", stale)
		}
	})
}

func TestBaselineRoundTrip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), DefaultBaselinePath)
	findings := []finding.Finding{
		findingAt("CG001", "main.go", 7),
		findingAt("CG006", "cmd/run.go", 91),
	}

	if err := WriteBaseline(findings, path); err != nil {
		t.Fatalf("WriteBaseline: %v", err)
	}

	loaded, found, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if !found {
		t.Fatal("LoadBaseline reported not-found for an existing file")
	}

	if loaded.Version != baselineVersion {
		t.Errorf("Version=%d, want %d", loaded.Version, baselineVersion)
	}

	if len(loaded.Entries) != len(findings) {
		t.Fatalf("entries=%d, want %d", len(loaded.Entries), len(findings))
	}

	fresh, stale := ApplyBaseline(findings, loaded)

	if len(fresh) != 0 || len(stale) != 0 {
		t.Errorf("round-trip fresh=%d stale=%d, want 0/0", len(fresh), len(stale))
	}
}

func TestLoadBaseline(t *testing.T) {
	t.Parallel()

	t.Run("missing file is nil without error", func(t *testing.T) {
		t.Parallel()

		baseline, found, err := LoadBaseline(filepath.Join(t.TempDir(), "absent.json"))

		if baseline != nil || found || err != nil {
			t.Errorf("baseline=%v found=%v err=%v, want nil/false/nil", baseline, found, err)
		}
	})

	t.Run("wrong version is rejected loudly", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), DefaultBaselinePath)
		if err := os.WriteFile(path, []byte(`{"version":99,"entries":[]}`), 0o644); err != nil {
			t.Fatalf("writing fixture: %v", err)
		}

		if _, _, err := LoadBaseline(path); err == nil {
			t.Error("expected an error for a version mismatch")
		}
	})

	t.Run("invalid JSON is rejected", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), DefaultBaselinePath)
		if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
			t.Fatalf("writing fixture: %v", err)
		}

		if _, _, err := LoadBaseline(path); err == nil {
			t.Error("expected an error for invalid JSON")
		}
	})
}
