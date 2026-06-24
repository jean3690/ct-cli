package snippet

import (
	"fmt"
	"strings"

	"github.com/jean/codeTemplateCli/pkg/snippet"
	"github.com/spf13/cobra"
)

var searchLimit int

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Aliases: []string{"find"},
	Short:   "Search snippets by name, description, or tags",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		results, err := store.Search(args[0])
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Printf("No snippets matching '%s'.\n", args[0])
			return nil
		}

		if searchLimit > 0 && len(results) > searchLimit {
			results = results[:searchLimit]
		}

		fmt.Printf("Found %d snippet(s) for '%s':\n\n", len(results), args[0])
		fmt.Printf("%-20s %-10s %s\n", "NAME", "LANGUAGE", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 80))
		for _, s := range results {
			desc := s.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			tags := ""
			if len(s.Tags) > 0 {
				tags = " [" + strings.Join(s.Tags, ", ") + "]"
			}
			fmt.Printf("%-20s %-10s %s%s\n", s.Name, s.Language, desc, tags)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 0, "max results to show")
}
