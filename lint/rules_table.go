package lint

import (
	"fmt"
	"strings"
)

// Markers bracketing the generated rule table in README.md. The table between
// them must equal RuleTableMarkdown() exactly; rules_table_test.go enforces
// the sync so adding a rule without refreshing the README fails CI.
const (
	readmeTableBegin = `<!-- BEGIN RULES TABLE: generated from allRuleDefs() — regenerate with ` + "`go run ./cmd/cmdguard-lint rules --markdown`" + ` -->`
	readmeTableEnd   = "<!-- END RULES TABLE -->"
)

// RuleTableMarkdown renders the full rule inventory as a Markdown table, one
// row per rule from allRuleDefs(). README.md embeds the output between the
// BEGIN/END markers; regenerate it with:
//
//	go run ./cmd/cmdguard-lint rules --markdown
//
// and replace the marked section. Pipes in descriptions are escaped so a
// future rule text cannot silently break the table.
func RuleTableMarkdown() string {
	var b strings.Builder

	b.WriteString("| ID | Severity | Rule | Description |\n")
	b.WriteString("| -- | -------- | ---- | ----------- |\n")

	for _, def := range allRuleDefs() {
		description := strings.ReplaceAll(def.meta.Description, "|", "\\|")
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", def.meta.ID, def.meta.Sev, def.meta.Name, description)
	}

	return b.String()
}
