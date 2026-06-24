package cmd

import (
	"fmt"
	"strings"

	"github.com/jean3690/ct-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the ~/.codeTemplate directory structure",
	Long: `Creates the default directory structure for ct under ~/.codeTemplate/,
including a config.yaml, templates/, and snippets/ directories.

Use CT_HOME environment variable to override the default location.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		if err := paths.EnsureDirs(); err != nil {
			return err
		}

		// Write default config if it doesn't exist.
		if _, err := config.ReadConfig(paths.ConfigFile); err != nil {
			// Use forward slashes to avoid YAML escape issues with Windows paths.
			usePath := func(p string) string {
				return strings.ReplaceAll(p, "\\", "/")
			}
			defaultCfg := fmt.Sprintf(`# ct configuration
storage:
  templatesDir: "%s"
  snippetsDir: "%s"

defaults:
  author: ""
  license: "MIT"
  interactive: true

editor: ""
`, usePath(paths.TemplatesDir), usePath(paths.SnippetsDir))

			if err := config.WriteConfig(paths.ConfigFile, []byte(defaultCfg)); err != nil {
				return fmt.Errorf("writing default config: %w", err)
			}
		}

		fmt.Printf("Initialized ct at %s\n", paths.ConfigDir)
		fmt.Println("  templates/  — store your code templates here")
		fmt.Println("  snippets/   — store your code snippets here")
		fmt.Println("  config.yaml — configuration file")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Bool("force", false, "overwrite existing config")
	viper.BindPFlag("init.force", initCmd.Flags().Lookup("force"))
}
