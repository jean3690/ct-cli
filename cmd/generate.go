package cmd

import (
	"fmt"

	"github.com/jean3690/ct-cli/pkg/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genVarFlags []string

var generateCmd = &cobra.Command{
	Use:     "generate [template]",
	Aliases: []string{"gen"},
	Short:   "Generate files from a template into the current directory",
	Long: `Generates files from a template into the current working directory.

Unlike 'ct new', this command generates files into an existing project
without creating a new project directory structure.

If no template is specified, the default template from config is used.
Set it with: ct config set defaults.template <name>`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := ""
		if len(args) > 0 {
			templateName = args[0]
		} else {
			templateName = viper.GetString("defaults.template")
			if templateName == "" {
				return fmt.Errorf("no template specified and no default configured. Set one with: ct config set defaults.template <name>")
			}
		}

		t, m, vars, err := loadTemplateAndVars(templateName, genVarFlags)
		if err != nil {
			return err
		}

		if dryRun {
			dryRunPreview(m, vars, ".")
			return nil
		}

		return scaffold.Scaffold(t.Path, ".", vars)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringArrayVarP(&genVarFlags, "var", "", nil, "template variable (e.g. --var Name=value)")
	generateCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing files")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview what would be generated")
}
