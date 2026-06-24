package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jean3690/ct-cli/pkg/engine"
	"github.com/jean3690/ct-cli/pkg/output"
	tmplpkg "github.com/jean3690/ct-cli/pkg/template"
)

// Scaffold generates a project from a template directory.
func Scaffold(templatePath, destDir string, vars map[string]string) error {
	manifestPath := filepath.Join(templatePath, "manifest.yaml")
	m, err := tmplpkg.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	eng := engine.New()

	renderedDir, err := eng.RenderString(destDir, vars)
	if err != nil {
		return fmt.Errorf("resolving destination dir: %w", err)
	}

	if m.HasScaffold() {
		if err := scaffoldDir(eng, templatePath, renderedDir, vars); err != nil {
			return err
		}
		return runPostGenerate(m, renderedDir, vars)
	}

	// File-based generation: process each FileMapping.
	filesDir := filepath.Join(templatePath, "files")
	for _, f := range m.Files {
		srcPath := filepath.Join(filesDir, f.Source)

		destRel, err := eng.RenderString(f.OutputDest(), vars)
		if err != nil {
			return fmt.Errorf("resolving dest for %s: %w", f.Source, err)
		}
		destPath := filepath.Join(renderedDir, destRel)

		if err := eng.RenderFile(srcPath, destPath, vars); err != nil {
			return err
		}
		output.FileCreated(destPath)
	}

	fmt.Printf("\nGenerated %d files.\n", len(m.Files))
	return runPostGenerate(m, renderedDir, vars)
}

func scaffoldDir(eng engine.TemplateEngine, templatePath, destDir string, vars map[string]string) error {
	filesDir := filepath.Join(templatePath, "files")

	count := 0
	err := filepath.Walk(filesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(filesDir, path)

		renderedRel, err := eng.RenderString(relPath, vars)
		if err != nil {
			return fmt.Errorf("resolving path %s: %w", relPath, err)
		}
		destPath := filepath.Join(destDir, renderedRel)

		if info.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("creating dir %s: %w", destPath, err)
			}
			output.DirCreated(destPath)
			return nil
		}

		if strings.HasSuffix(info.Name(), ".tmpl") {
			if err := eng.RenderFile(path, destPath, vars); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("creating dir for %s: %w", destPath, err)
			}
			if err := os.WriteFile(destPath, data, info.Mode()); err != nil {
				return fmt.Errorf("writing %s: %w", destPath, err)
			}
		}

		output.FileCreated(destPath)
		count++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nGenerated %d files.\n", count)
	return nil
}

// runPostGenerate executes the postGenerate script if defined.
func runPostGenerate(m *tmplpkg.TemplateManifest, destDir string, vars map[string]string) error {
	if m.PostGenerate == "" {
		return nil
	}

	eng := engine.New()
	script, err := eng.RenderString(m.PostGenerate, vars)
	if err != nil {
		return fmt.Errorf("rendering postGenerate: %w", err)
	}

	fmt.Printf("\nRunning post-generate hook...\n")
	fmt.Printf("  > %s\n", script)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", script)
	} else {
		cmd = exec.Command("sh", "-c", script)
	}
	cmd.Dir = destDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("postGenerate hook failed: %w", err)
	}

	return nil
}


