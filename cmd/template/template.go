package template

import "github.com/spf13/cobra"

// Cmd is the parent template command.
var Cmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"tpl", "tmpl"},
	Short:   "Manage code templates",
	Long:    `Manage code templates stored under ~/.codeTemplate/templates/.`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(infoCmd)
	Cmd.AddCommand(importCmd)
	Cmd.AddCommand(removeCmd)
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(editCmd)
}
