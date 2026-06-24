package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jean3690/ct-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and change ct configuration settings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		fmt.Printf("Config file: %s\n\n", paths.ConfigFile)
		fmt.Printf("%-25s %s\n", "KEY", "VALUE")
		fmt.Println("--------------------------------------------------")
		for _, k := range viper.AllKeys() {
			fmt.Printf("%-25s %v\n", k, viper.Get(k))
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value and persist it to config.yaml.

Common keys:
  defaults.template    - default template for ct new/generate
  defaults.author      - default author name
  defaults.interactive - enable/disable interactive mode (true/false)
  defaults.license     - default license (MIT, Apache-2.0, etc.)
  editor              - editor command for ct template edit`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		viper.Set(key, value)

		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		// Ensure config file exists.
		if err := os.MkdirAll(filepath.Dir(paths.ConfigFile), 0755); err != nil {
			return err
		}

		if err := viper.WriteConfigAs(paths.ConfigFile); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}
