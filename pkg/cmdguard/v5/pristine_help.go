package v5

import (
	"encoding/json/v2"

	"github.com/spf13/cobra"
)

// pristineHelp is the annotation payload: a command's Long and Example text
// exactly as the author wrote them, captured before any help transform runs.
type pristineHelp struct {
	Long    string `json:"long"`
	Example string `json:"example"`
}

// stashPristineHelpText snapshots Long and Example of every command in the
// tree ahead of the help transforms, which overwrite the live cobra fields
// with terminal-styled bytes. The annotation is written only when absent, so
// repeated Execute calls re-stash nothing: after the first pass the live
// fields are already styled, and the snapshot must keep pointing at the
// original text. Encoding a plain string pair cannot fail, but a marshal
// error is tolerated as absent — every reader falls back to the live fields,
// which are identical when no transform ran.
func stashPristineHelpText(cmd *cobra.Command) {
	if _, ok := cmd.Annotations[pristineHelpAnnotation]; !ok {
		encoded, err := json.Marshal(pristineHelp{
			Long:    cmd.Long,
			Example: cmd.Example,
		})
		if err == nil {
			if cmd.Annotations == nil {
				cmd.Annotations = make(map[string]string)
			}

			cmd.Annotations[pristineHelpAnnotation] = string(encoded)
		}
	}

	for _, sub := range cmd.Commands() {
		stashPristineHelpText(sub)
	}
}

// pristineHelpFor returns the command's pre-transform Long and Example text.
// Absent or malformed annotations fall back to the live fields: annotations
// are internal, machine-written data, and neither Describe nor GenerateDocs
// may fail because of them.
func pristineHelpFor(cmd *cobra.Command) (string, string) {
	encoded, ok := cmd.Annotations[pristineHelpAnnotation]
	if ok {
		var help pristineHelp

		if json.Unmarshal([]byte(encoded), &help) == nil {
			return help.Long, help.Example
		}
	}

	return cmd.Long, cmd.Example
}
