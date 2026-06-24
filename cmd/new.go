package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jean3690/ct-cli/pkg/config"
	"github.com/jean3690/ct-cli/pkg/scaffold"
	tmplpkg "github.com/jean3690/ct-cli/pkg/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	varFlags []string
	force    bool
	dryRun   bool
)

var newCmd = &cobra.Command{
	Use:   "new [template] [directory]",
	Short: "Scaffold a new project from a template",
	Long: `Generates a new project directory from an installed template.

If template is not specified, the default template from config is used.
If directory is not specified, the current directory is used.
Variables are prompted interactively (unless --yes), or passed via --var.

Set default template: ct config set defaults.template <name>`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := ""
		destDir := "."
		if len(args) >= 1 {
			templateName = args[0]
		}
		if len(args) >= 2 {
			destDir = args[1]
		}

		// If first arg isn't a known template and we have a default, treat it as directory.
		if templateName != "" {
			paths, err := config.InitPaths()
			if err != nil {
				return err
			}
			if _, err := tmplpkg.FindTemplate(paths.TemplatesDir, templateName); err != nil {
				defaultTmpl := viper.GetString("defaults.template")
				if defaultTmpl != "" {
					destDir = templateName
					templateName = defaultTmpl
				}
			}
		}

		if templateName == "" {
			templateName = viper.GetString("defaults.template")
			if templateName == "" {
				return fmt.Errorf("no template specified and no default configured; set one with: ct config set defaults.template <name>")
			}
		}

		t, m, vars, err := loadTemplateAndVars(templateName, varFlags)
		if err != nil {
			return err
		}

		if dryRun {
			dryRunPreview(m, vars, destDir)
			return nil
		}

		if !force {
			absDest, _ := filepath.Abs(destDir)
			if entries, _ := os.ReadDir(absDest); len(entries) > 0 {
				return fmt.Errorf("destination %s is not empty. Use --force to overwrite", destDir)
			}
		}

		return scaffold.Scaffold(t.Path, destDir, vars)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringArrayVarP(&varFlags, "var", "", nil, "template variable (e.g. --var Name=value)")
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing files")
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview what would be generated")
}

func parseVarFlags(flags []string) map[string]string {
	result := make(map[string]string)
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
