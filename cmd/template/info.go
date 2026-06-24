package template

import (
	"fmt"
	"path/filepath"

	"github.com/jean3690/ct-cli/pkg/config"
	tmplpkg "github.com/jean3690/ct-cli/pkg/template"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show template details including variables",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		manifestPath := filepath.Join(paths.TemplatesDir, args[0], "manifest.yaml")
		m, err := tmplpkg.LoadManifest(manifestPath)
		if err != nil {
			return fmt.Errorf("loading template %s: %w", args[0], err)
		}

		fmt.Printf("Name:        %s\n", m.Name)
		fmt.Printf("Description: %s\n", m.Description)
		fmt.Printf("Version:     %s\n", m.Version)
		if m.Author != "" {
			fmt.Printf("Author:      %s\n", m.Author)
		}
		if len(m.Languages) > 0 {
			fmt.Print("Languages:   ")
			for i, l := range m.Languages {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Print(l)
			}
			fmt.Println()
		}
		fmt.Printf("Type:        ")
		if m.HasScaffold() {
			fmt.Println("project scaffold")
		} else {
			fmt.Println("file generator")
		}

		if len(m.Variables) > 0 {
			fmt.Println("\nVariables:")
			fmt.Printf("  %-16s %-8s %-10s %s\n", "NAME", "TYPE", "DEFAULT", "PROMPT")
			fmt.Println("  --------------------------------------------------------")
			for _, v := range m.Variables {
				req := ""
				if v.Required {
					req = "*"
				}
				fmt.Printf("  %-16s %-8s %-10s %s%s\n", v.Name, v.Type, v.Default, v.PromptText(), req)
			}
			fmt.Println("  * = required")
		}

		if len(m.Files) > 0 {
			fmt.Println("\nFiles:")
			for _, f := range m.Files {
				fmt.Printf("  %s -> %s\n", f.Source, f.OutputDest())
			}
		}

		return nil
	},
}
