package snippet

import (
	"fmt"

	"github.com/jean/codeTemplateCli/pkg/snippet"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get <name>",
	Aliases: []string{"show", "cat"},
	Short:   "Display snippet content",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		snip, err := store.Load(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("# %s (%s)\n", snip.Name, snip.Language)
		if snip.Description != "" {
			fmt.Printf("# %s\n", snip.Description)
		}
		if len(snip.Tags) > 0 {
			fmt.Printf("# tags: ")
			for i, t := range snip.Tags {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(t)
			}
			fmt.Println()
		}
		fmt.Println()
		fmt.Println(snip.Code)
		return nil
	},
}
