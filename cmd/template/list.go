package template

import (
	"fmt"
	"strings"

	"github.com/jean/codeTemplateCli/pkg/config"
	"github.com/jean/codeTemplateCli/pkg/template"
	"github.com/spf13/cobra"
)

var listLanguage string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		templates, err := template.Discover(paths.TemplatesDir, listLanguage)
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			if listLanguage != "" {
				fmt.Printf("No templates found for language '%s'. Use 'ct template import' to add one.\n", listLanguage)
			} else {
				fmt.Println("No templates installed. Use 'ct template import <path>' to add one.")
			}
			return nil
		}

		fmt.Printf("%-20s %-10s %-20s %s\n", "NAME", "VERSION", "LANGUAGES", "DESCRIPTION")
		fmt.Println(strings.Repeat("-", 80))
		for _, t := range templates {
			desc := t.Description
			if len(desc) > 40 {
				desc = desc[:37] + "..."
			}
			langs := strings.Join(t.Languages, ", ")
			if langs == "" {
				langs = "*"
			}
			fmt.Printf("%-20s %-10s %-20s %s\n", t.Name, t.Version, langs, desc)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&listLanguage, "language", "l", "", "filter by programming language (e.g. go, python, rust)")
}
