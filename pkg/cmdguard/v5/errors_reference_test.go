package v5

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// errorReferencePath is the committed doc the generated table must match.
const errorReferencePath = "../../../docs/ERROR_REFERENCE.md"

func TestErrorReferenceDoc(t *testing.T) { //nolint:paralleltest // rewrites a repo file under -update-golden
	doc, err := os.ReadFile(errorReferencePath)
	if err != nil {
		t.Fatalf("reading error reference doc (run from pkg/cmdguard/v5): %v", err)
	}

	content := string(doc)

	generated := GenerateErrorReference()

	if *updateGolden {
		var spliced string

		start := strings.Index(content, errorRefBegin)
		stop := strings.Index(content, errorRefEnd)

		if start >= 0 && stop >= start {
			spliced = content[:start] + generated + content[stop+len(errorRefEnd):]
		} else {
			spliced = content + "\n" + generated + "\n"
		}

		if writeErr := os.WriteFile(errorReferencePath, []byte(spliced), 0o644); writeErr != nil {
			t.Fatalf("rewriting error reference doc: %v", writeErr)
		}

		return
	}

	begin := strings.Count(content, errorRefBegin)
	end := strings.Count(content, errorRefEnd)
	if begin != 1 || end != 1 {
		t.Fatalf(
			"expected exactly one %s / %s marker pair in the doc, got %d/%d",
			errorRefBegin,
			errorRefEnd,
			begin,
			end,
		)
	}

	start := strings.Index(content, errorRefBegin)
	stop := strings.Index(content, errorRefEnd) + len(errorRefEnd)
	if stop < start {
		t.Fatal("error reference markers are out of order")
	}

	current := content[start:stop]
	if normalizeTableSpacing(current) != normalizeTableSpacing(generated) {
		t.Fatalf(`docs/ERROR_REFERENCE.md drifted from the sentinel fact table.

Regenerate with:
  go test ./pkg/cmdguard/v5 -run TestErrorReferenceDoc -update-golden`)
	}
}

func normalizeTableSpacing(s string) string {
	pipeSpace := regexp.MustCompile(` *\| *`)
	dashes := regexp.MustCompile(`-{3,}`)
	separator := regexp.MustCompile(`^[|: -]+$`)

	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if separator.MatchString(line) {
			line = dashes.ReplaceAllString(line, "---")
		}

		lines[i] = pipeSpace.ReplaceAllString(line, "|")
	}

	return strings.Join(lines, "\n")
}

func TestNormalizeTableSpacing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string is unchanged",
			input: "",
			want:  "",
		},
		{
			name: "aligned table collapses to unpadded cells and short separators",
			input: "| Go name      | Code             |\n" +
				"| ------------- | ---------------- |\n" +
				"| `ErrRequired` | `required_flag`  |",
			want: "|Go name|Code|\n|---|---|\n|`ErrRequired`|`required_flag`|",
		},
		{
			name:  "unaligned table normalizes to the same shape as aligned",
			input: "| Go name | Code |\n| --- | --- |\n| `ErrRequired` | `required_flag` |",
			want:  "|Go name|Code|\n|---|---|\n|`ErrRequired`|`required_flag`|",
		},
		{
			name:  "colon alignment markers survive dash collapsing",
			input: "| left | center |\n| :--- | :---: |\n| a    | b      |",
			want:  "|left|center|\n|:---|:---:|\n|a|b|",
		},
		{
			name:  "spaces inside cells are preserved",
			input: "| two words | none |\n| --- | --- |",
			want:  "|two words|none|\n|---|---|",
		},
		{
			name:  "prose without pipes is untouched",
			input: "Some prose line\nwith a hyphenated-word and - dashes",
			want:  "Some prose line\nwith a hyphenated-word and - dashes",
		},
		{
			name:  "long plain-dash rule collapses like a separator row",
			input: "intro\n----------\noutro",
			want:  "intro\n---\noutro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeTableSpacing(tt.input)
			if got != tt.want {
				t.Fatalf("normalizeTableSpacing(%q) = %q, want %q", tt.input, got, tt.want)
			}

			if again := normalizeTableSpacing(got); again != got {
				t.Fatalf("normalizeTableSpacing is not idempotent: %q -> %q", got, again)
			}
		})
	}

	t.Run("aligned and unaligned renderings compare equal", func(t *testing.T) {
		t.Parallel()

		alignedInput := "| Go name      | Code             |\n" +
			"| ------------- | ---------------- |\n" +
			"| `ErrRequired` | `required_flag`  |"
		unalignedInput := "| Go name | Code |\n| --- | --- |\n| `ErrRequired` | `required_flag` |"

		aligned := normalizeTableSpacing(alignedInput)
		unaligned := normalizeTableSpacing(unalignedInput)

		if aligned != unaligned {
			t.Fatalf("aligned vs unaligned normalized forms differ:\n%q\n%q", aligned, unaligned)
		}
	})
}

func TestGenerateErrorReference(t *testing.T) {
	t.Parallel()

	ref := GenerateErrorReference()

	if !strings.Contains(ref, errorRefBegin) {
		t.Fatal("missing begin marker")
	}

	if !strings.Contains(ref, errorRefEnd) {
		t.Fatal("missing end marker")
	}

	for _, fact := range sentinelFacts() {
		needle := "| `" + fact.code + "` | `" + fact.name + "` |"
		if !strings.Contains(ref, needle) {
			t.Errorf("table row missing for %s", fact.name)
		}
	}

	if !strings.Contains(ref, "### flag errors") {
		t.Fatal("expected class sections in the generated reference")
	}
}

// skillMdPath is the consumer-facing skill source (not part of the
// published module zip, so consumer checkouts skip the count tripwire).
const skillMdPath = "../../../.agents/skills/cmdguard/SKILL.md"

// TestSkillSentinelCountMatchesReference fails when the consumer SKILL.md's
// stated sentinel count drifts from the generated ERROR_REFERENCE count.
// The SKILL.md count is hand-derived via an rg pattern and its stored
// derivation rotted once already (the pattern matched nothing after a
// formatter re-wrapped the table, caught 2026-11-22); this tripwire turns
// that cross-file agreement into a machine check: every sentinel count
// stated in the skill must equal the live generator's count, and at least
// one must be present.
func TestSkillSentinelCountMatchesReference(t *testing.T) {
	t.Parallel()

	ref := GenerateErrorReference()
	wantCounts := regexp.MustCompile(`(\d+) sentinels`).FindAllStringSubmatch(ref, -1)
	if len(wantCounts) == 0 {
		t.Fatal("generated reference carries no sentinel count; the count-tripwire parse is broken")
	}
	want := wantCounts[0][1]
	for _, m := range wantCounts[1:] {
		if m[1] != want {
			t.Fatalf("generated reference states conflicting sentinel counts: %s vs %s", want, m[1])
		}
	}

	skill, err := os.ReadFile(skillMdPath)
	if os.IsNotExist(err) {
		t.Skipf("consumer checkout: %s not present outside this repo", skillMdPath)
	}
	if err != nil {
		t.Fatalf("reading skill md: %v", err)
	}

	skillCounts := regexp.MustCompile(`(\d+) (?:documented sentinel|sentinels)`).
		FindAllStringSubmatch(string(skill), -1)
	if len(skillCounts) == 0 {
		t.Fatal(
			"SKILL.md states no sentinel count; the derivation pattern in the skill has rotted — re-derive and update it",
		)
	}
	seen := map[string]bool{}
	for _, m := range skillCounts {
		if m[1] != want {
			t.Errorf("SKILL.md sentinel count %s does not match the generated reference count %s", m[1], want)
		}
		seen[m[1]] = true
	}
	if len(seen) > 1 {
		keys := make([]string, 0, len(seen))
		for k := range seen {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		t.Errorf("SKILL.md states conflicting sentinel counts %v; unify them on %s", keys, want)
	}
}
