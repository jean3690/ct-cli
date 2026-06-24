package snippet

import (
	"fmt"

	"github.com/jean3690/ct-cli/pkg/snippet"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm"},
	Short:   "Delete a code snippet",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		if err := store.Delete(args[0]); err != nil {
			return err
		}

		fmt.Printf("Snippet '%s' deleted.\n", args[0])
		return nil
	},
}
