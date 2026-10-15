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
//
// Columns are padded to the widest cell (dprint's Markdown table format), so
// the generated section is a fixed point for both this generator and the
// dprint formatter — rules_table_test.go pins it byte-for-byte, and running
// dprint over README.md never rewrites it.
func RuleTableMarkdown() string {
	header := []string{"ID", "Severity", "Rule", "Description"}
	sep := []string{"--", "--------", "----", "-----------"}

	rows := make([][]string, 0, len(allRuleDefs())+2)
	rows = append(rows, header, sep)

	for _, def := range allRuleDefs() {
		description := strings.ReplaceAll(def.meta.Description, "|", "\\|")
		rows = append(rows, []string{def.meta.ID, fmt.Sprint(def.meta.Sev), def.meta.Name, description})
	}

	width := make([]int, len(header))

	for _, row := range rows {
		if row[0] == "--" {
			continue
		}

		for i, cell := range row {
			if len(cell) > width[i] {
				width[i] = len(cell)
			}
		}
	}

	var b strings.Builder

	writeRow := func(cells []string, isSep bool) {
		parts := make([]string, len(cells))
		for i, cell := range cells {
			if isSep {
				parts[i] = " " + strings.Repeat("-", width[i]) + " "
			} else {
				parts[i] = " " + cell + strings.Repeat(" ", width[i]-len(cell)) + " "
			}
		}

		b.WriteString("|" + strings.Join(parts, "|") + "|\n")
	}

	for _, row := range rows {
		writeRow(row, row[0] == "--")
	}

	return b.String()
}
