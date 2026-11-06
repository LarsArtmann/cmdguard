package v4

import (
	"os"
	"strings"
	"testing"
)

// errorReferencePath is the committed doc the generated table must match.
const errorReferencePath = "../../../docs/ERROR_REFERENCE.md"

func TestErrorReferenceDoc(t *testing.T) { //nolint:paralleltest // rewrites a repo file under -update-golden
	doc, err := os.ReadFile(errorReferencePath)
	if err != nil {
		t.Fatalf("reading error reference doc (run from pkg/cmdguard/v4): %v", err)
	}

	content := string(doc)

	generated := GenerateErrorReference()

	if *updateGolden {
		spliced := content

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
		t.Fatalf("expected exactly one %s / %s marker pair in the doc, got %d/%d", errorRefBegin, errorRefEnd, begin, end)
	}

	start := strings.Index(content, errorRefBegin)
	stop := strings.Index(content, errorRefEnd) + len(errorRefEnd)
	if stop < start {
		t.Fatal("error reference markers are out of order")
	}

	current := content[start:stop]
	if current != generated {
		t.Fatalf(`docs/ERROR_REFERENCE.md drifted from the sentinel fact table.

Regenerate with:
  go test ./pkg/cmdguard/v4 -run TestErrorReferenceDoc -update-golden`)
	}
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
