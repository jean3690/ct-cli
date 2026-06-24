package snippet

import (
	"fmt"
	"strings"

	"github.com/jean/codeTemplateCli/pkg/snippet"
	"github.com/spf13/cobra"
)

var listLang string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List code snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		snippets, err := store.List(listLang)
		if err != nil {
			return err
		}

		if len(snippets) == 0 {
			filterMsg := ""
			if listLang != "" {
				filterMsg = " for language: " + listLang
			}
			fmt.Printf("No snippets found%s. Use 'ct snippet create <name>' to add one.\n", filterMsg)
			return nil
		}

		fmt.Printf("%-20s %-10s %s\n", "NAME", "LANGUAGE", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 80))
		for _, s := range snippets {
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
	listCmd.Flags().StringVarP(&listLang, "language", "l", "", "filter by language")
}
