package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jean3690/ct-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Open a template directory for editing",
	Long: `Opens the template directory in your default editor.

The editor can be set via config (editor field) or the EDITOR environment variable.
On Windows, it defaults to opening in Explorer if no editor is configured.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := config.InitPaths()
		if err != nil {
			return err
		}

		templateDir := filepath.Join(paths.TemplatesDir, args[0])
		if _, err := os.Stat(templateDir); os.IsNotExist(err) {
			return fmt.Errorf("template '%s' not found", args[0])
		}

		editor := viper.GetString("editor")

		// If no editor configured, use OS-specific default.
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			// Fall back to showing the path.
			fmt.Printf("Template directory: %s\n", templateDir)
			fmt.Println("Set editor in config to open directly: ct config set editor <command>")
			return nil
		}

		c := exec.Command(editor, templateDir)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	},
}
