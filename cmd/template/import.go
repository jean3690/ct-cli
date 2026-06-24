package template

import (
	"fmt"

	tmplpkg "github.com/jean3690/ct-cli/pkg/template"
	"github.com/spf13/cobra"
)

var (
	importName string
	importGit  bool
)

var importCmd = &cobra.Command{
	Use:   "import <source> [name]",
	Short: "Import a template from a local directory or git repository",
	Long: `Import a template from a local directory or git repository.

The source must contain a manifest.yaml at its root.

Examples:
  ct template import ./my-template
  ct template import ./my-template my-name
  ct template import --git https://github.com/user/go-service-template
  ct template import --git https://github.com/user/tmpl my-name`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		name := importName
		if name == "" && len(args) > 1 {
			name = args[1]
		}

		if importGit {
			if err := tmplpkg.ImportFromGit(src, name); err != nil {
				return err
			}
		} else {
			if err := tmplpkg.ImportFromPath(src, name); err != nil {
				return err
			}
		}

		fmt.Printf("Template imported successfully.\n")
		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&importName, "name", "n", "", "template name (defaults to manifest name)")
	importCmd.Flags().BoolVarP(&importGit, "git", "g", false, "import from a git repository URL")
}
