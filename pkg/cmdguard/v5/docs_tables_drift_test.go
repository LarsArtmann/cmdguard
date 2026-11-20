package v5

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Drift tripwires (SU25/SU26/SU27, 2026-11-20): the hand-kept doc tables that
// describe code facts are pinned to the code itself. When a fact changes in
// source, these tests fail until every mirrored table is updated — the
// docgen_drift pattern applied to README + website + SKILL.md.

// documentedTagFields maps every doc-facing FlagTag field to the struct-tag
// spelling the docs use. Field/Type/Index/Key are plumbing, not user-facing
// tags, and are deliberately absent: add a user-facing tag here AND to the
// mirrored tables when FlagTag grows one.
var documentedTagFields = map[string]string{
	"Name":     "flag",
	"Short":    "short",
	"Default":  "default",
	"Help":     "help",
	"Values":   "values",
	"Required": "required",
	"Validate": "validate",
	"Env":      "env",
	"Count":    "count",
	"Prompt":   "prompt",
	"Local":    "local",
	"Hidden":   "hidden",
	"Secret":   "secret",
}

// plumbingTagFields are FlagTag fields that carry no struct tag of their own.
var plumbingTagFields = map[string]bool{
	"Field": true,
	"Type":  true,
	"Index": true,
	"Key":   true,
}

func readRepoDoc(t *testing.T, rel string) string {
	t.Helper()
	path := filepath.Join("..", "..", "..", rel)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return string(content)
}

// rowHasRule reports whether a markdown validate table contains a row whose
// first cell is exactly the rule name (with or without a `=param` suffix).
func rowHasRule(table, rule string) bool {
	return strings.Contains(table, "| `"+rule+"`") || strings.Contains(table, "| `"+rule+"=")
}

func TestFlagTagFieldsTrackedInDocs(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeFor[FlagTag]()
	seen := map[string]bool{}
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		seen[name] = true
		if plumbingTagFields[name] {
			continue
		}
		if _, ok := documentedTagFields[name]; !ok {
			t.Errorf("FlagTag gained a doc-facing field %q: add its struct-tag spelling to documentedTagFields and to the README/SKILL tag tables", name)
		}
	}
	for name := range documentedTagFields {
		if !seen[name] {
			t.Errorf("documentedTagFields lists %q but FlagTag has no such field: remove it from the map and the doc tables", name)
		}
	}

	readme := readRepoDoc(t, "README.md")
	for _, tag := range documentedTagFields {
		if !strings.Contains(readme, "| `"+tag+"`") {
			t.Errorf("README 13-tag table lost a row for tag %q (source of truth: FlagTag in config.go)", tag)
		}
	}

	skill := readRepoDoc(t, ".agents/skills/cmdguard/SKILL.md")
	skillSet := "flag, short, default, help, values, required, count, local, hidden, validate, env, prompt, secret"
	if !strings.Contains(skill, skillSet) {
		t.Errorf("SKILL.md full-tag-set line drifted from FlagTag: expected %q", skillSet)
	}
}

func TestValidateRulesTrackedInDocs(t *testing.T) {
	t.Parallel()

	registry := newValidatorRegistry()
	names := make([]string, 0, len(registry.validators))
	for name := range registry.validators {
		names = append(names, name)
	}
	sort.Strings(names)

	const wantCount = 10
	if len(names) != wantCount {
		t.Errorf("builtin validator count is %d, docs say %d: update the README + website validate tables (and this constant) together", len(names), wantCount)
	}

	readme := readRepoDoc(t, "README.md")
	site := readRepoDoc(t, "website/src/content/docs/guides/type-safe-flags.mdx")
	for _, rule := range names {
		if !rowHasRule(readme, rule) {
			t.Errorf("README validate table lost the %q row (source of truth: registerBuiltins in flags_validate.go)", rule)
		}
		if !rowHasRule(site, rule) {
			t.Errorf("website type-safe-flags.mdx validate table lost the %q row — the site copy must track the same registry", rule)
		}
	}
}
