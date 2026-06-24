package template

import (
	"fmt"

	tmplpkg "github.com/jean3690/ct-cli/pkg/template"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an installed template",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := tmplpkg.RemoveTemplate(args[0]); err != nil {
			return err
		}
		fmt.Printf("Template '%s' removed.\n", args[0])
		return nil
	},
}
