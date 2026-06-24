package snippet

import "github.com/spf13/cobra"

// Cmd is the parent snippet command.
var Cmd = &cobra.Command{
	Use:     "snippet",
	Aliases: []string{"snip"},
	Short:   "Manage code snippets",
	Long:    `Manage code snippets stored under ~/.codeTemplate/snippets/.`,
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(searchCmd)
}
