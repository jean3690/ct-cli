package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/jean/codeTemplateCli/pkg/config"
	"github.com/jean/codeTemplateCli/pkg/engine"
	"github.com/jean/codeTemplateCli/pkg/prompts"
	tmplpkg "github.com/jean/codeTemplateCli/pkg/template"
)

// loadTemplateAndVars resolves a template by name, loads its manifest,
// collects variables, and returns everything needed for scaffolding.
func loadTemplateAndVars(templateName string, varFlags []string) (*tmplpkg.DiscoveredTemplate, *tmplpkg.TemplateManifest, map[string]string, error) {
	paths, err := config.InitPaths()
	if err != nil {
		return nil, nil, nil, err
	}

	t, err := tmplpkg.FindTemplate(paths.TemplatesDir, templateName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("template %s not found. Use 'ct template list' to see available templates", templateName)
	}

	manifestPath := filepath.Join(t.Path, "manifest.yaml")
	m, err := tmplpkg.LoadManifest(manifestPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading manifest: %w", err)
	}

	provided := parseVarFlags(varFlags)
	vars, err := prompts.CollectVariables(m.Variables, provided, config.IsInteractive())
	if err != nil {
		return nil, nil, nil, err
	}
	if vars == nil {
		vars = make(map[string]string)
	}
	vars["TemplateName"] = m.Name

	return t, m, vars, nil
}

// dryRunPreview prints what would be generated without actually creating files.
func dryRunPreview(m *tmplpkg.TemplateManifest, vars map[string]string, destDir string) {
	fmt.Printf("Template: %s (%s)\n", m.Name, m.Description)
	fmt.Println("Variables:")
	for k, v := range vars {
		fmt.Printf("  %s = %s\n", k, v)
	}
	fmt.Println("\nWould generate:")
	eng := engine.New()
	for _, f := range m.Files {
		dest, _ := eng.RenderString(f.OutputDest(), vars)
		fmt.Printf("  %s -> %s\n", f.Source, filepath.Join(destDir, dest))
	}
}
