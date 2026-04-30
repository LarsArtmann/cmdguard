package v2

// SuggestCommand returns the best matching command name for a potentially misspelled input.
// Returns the best match and true if a good match is found, or empty string and false otherwise.
func SuggestCommand(validCommands []string, input string) (string, bool) {
	if len(validCommands) == 0 {
		return "", false
	}

	return SuggestFlag(validCommands, input)
}
