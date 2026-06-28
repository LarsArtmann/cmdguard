package v2

import (
	"iter"
	"strings"
)

// maxEditDistance is the threshold for flag name suggestions.
const maxEditDistance = 3

// SuggestFlag returns the best matching flag name for a potentially misspelled input.
// Returns the best match and true if a good match is found, or empty string and false otherwise.
func SuggestFlag(validNames []string, input string) (string, bool) {
	if len(validNames) == 0 {
		return "", false
	}

	bestMatch := ""
	bestDist := maxEditDistance + 1

	for _, name := range validNames {
		dist := editDistance(input, name)
		if dist < bestDist {
			bestDist = dist
			bestMatch = name
		}
	}

	if bestDist <= maxEditDistance {
		return bestMatch, true
	}

	return "", false
}

// SuggestCommand returns the best matching command name for a potentially misspelled input.
// Returns the best match and true if a good match is found, or empty string and false otherwise.
func SuggestCommand(validCommands []string, input string) (string, bool) {
	if len(validCommands) == 0 {
		return "", false
	}

	return SuggestFlag(validCommands, input)
}

// editDistance computes the Levenshtein distance between two strings.
func editDistance(a, b string) int {
	aLen, bLen := len(a), len(b)
	if aLen == 0 {
		return bLen
	}

	if bLen == 0 {
		return aLen
	}

	// Use a single row for space optimization
	prev := make([]int, bLen+1, bLen+1)
	curr := make([]int, bLen+1, bLen+1)

	for j := 0; j <= bLen; j++ {
		prev[j] = j
	}

	for i := 1; i <= aLen; i++ {
		curr[0] = i

		for j := 1; j <= bLen; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}

			curr[j] = min(
				min(prev[j]+1, curr[j-1]+1),
				prev[j-1]+cost,
			)
		}

		prev, curr = curr, prev
	}

	return prev[bLen]
}

// GenerateHelp generates help text for all flags.
func (r *FlagRegistry) GenerateHelp() string {
	var lines []string

	for _, tag := range r.tags {
		line := "  --" + tag.Name
		if tag.Short != "" {
			line += ", -" + tag.Short
		}

		line += "\t" + tag.Help
		if tag.Default != "" {
			line += " (default: " + tag.Default + ")"
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// FlagNames returns all registered flag names for suggestion purposes.
// For zero-allocation iteration, prefer FlagNamesSeq().
func (r *FlagRegistry) FlagNames() []string {
	names := make([]string, len(r.tags), len(r.tags))
	for i, tag := range r.tags {
		names[i] = tag.Name
	}

	return names
}

// FlagNamesSeq returns an iterator over all registered flag names without allocating a slice.
// Prefer this over FlagNames() when you only need to range over the names.
func (r *FlagRegistry) FlagNamesSeq() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, tag := range r.tags {
			if !yield(tag.Name) {
				return
			}
		}
	}
}
