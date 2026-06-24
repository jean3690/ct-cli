package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initLang string

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Scaffold a new template (meta-template for authoring)",
	Long: `Creates a new template skeleton that you can fill in and install.

This is a "template for writing templates" — it generates:
  - manifest.yaml (with common variable prompts)
  - files/ directory for your .tmpl source files

After filling in your template files, install it with:
  ct template import ./<name>

Examples:
  ct template init rust-service --lang rust
  ct template init zig-cli --lang zig`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Don't overwrite existing directories.
		if _, err := os.Stat(name); err == nil {
			return fmt.Errorf("directory '%s' already exists", name)
		}

		lang := initLang
		langs := "[]"
		if lang != "" {
			langs = fmt.Sprintf("[\"%s\"]", lang)
		}

		// Create directory structure.
		dirs := []string{name, filepath.Join(name, "files")}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				return fmt.Errorf("creating %s: %w", d, err)
			}
		}

		// Write manifest.yaml skeleton.
		manifest := fmt.Sprintf(`# Template manifest for %s
# Fill in the details below, then install with: ct template import ./%s
apiVersion: v1
name: %s
description: "Describe your template here"
version: "0.1.0"
languages: %s

# Variables the user will be prompted for.
# Types: string (default), bool, choice
variables:
  - name: ProjectName
    prompt: "Project name?"
    type: string
    required: true
  # Add more variables as needed:
  # - name: UseDatabase
  #   prompt: "Include database support?"
  #   type: bool
  #   default: "false"

# Files to generate (relative to files/ directory).
# .tmpl files are rendered with Go text/template; others are copied as-is.
files:
  # - source: main.go.tmpl
  #   dest: "{{ .ProjectName }}/main.go"

# Optional: shell command to run after generation.
# Variables are substituted before execution.
# postGenerate: "go mod init {{ .ProjectName }} && go fmt ./..."
`, name, name, name, langs)

		manifestPath := filepath.Join(name, "manifest.yaml")
		if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
			return fmt.Errorf("writing manifest: %w", err)
		}

		fmt.Printf("Template '%s' created.\n\n", name)
		fmt.Println("  manifest.yaml  — template metadata and variables")
		fmt.Println("  files/          — put your .tmpl source files here")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Edit manifest.yaml to describe your template\n")
		fmt.Printf("  2. Add .tmpl files to files/\n")
		fmt.Printf("  3. Install: ct template import ./%s\n", name)
		fmt.Printf("  4. Test:    ct new %s -y ./test-output\n", name)

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initLang, "lang", "l", "", "primary language (e.g. go, rust, python)")
}
