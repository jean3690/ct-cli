package cmd

import (
	"fmt"
	"os"

	"github.com/jean/codeTemplateCli/cmd/snippet"
	"github.com/jean/codeTemplateCli/cmd/template"
	"github.com/jean/codeTemplateCli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.1.0"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "ct-cli",
	Version: version,
	Short:   "ct-cli — a code template and scaffolding CLI",
	Long: `ct-cli is a CLI tool for managing code templates, scaffolding projects,
and organizing code snippets.

Store templates and snippets under ~/.codeTemplate/ and generate
projects or files with interactive prompts or pure CLI flags.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.InitViper(cfgFile)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(template.Cmd)
	rootCmd.AddCommand(snippet.Cmd)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ~/.codeTemplate/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "skip interactive confirmations")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))
}
