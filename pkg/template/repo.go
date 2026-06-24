package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jean3690/ct-cli/pkg/config"
)

// ImportFromPath imports a template from a local directory into the templates store.
func ImportFromPath(srcPath, name string) error {
	paths, err := config.InitPaths()
	if err != nil {
		return err
	}

	// Validate source has a manifest.yaml and load it to get the canonical name.
	srcManifest := filepath.Join(srcPath, "manifest.yaml")
	m, err := LoadManifest(srcManifest)
	if err != nil {
		return fmt.Errorf("loading manifest from %s: %w", srcPath, err)
	}

	// If name not provided, use the manifest name.
	if name == "" {
		name = m.Name
	}

	destPath := filepath.Join(paths.TemplatesDir, name)

	// Check if already exists.
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("template %s already exists, remove it first", name)
	}

	// Copy directory.
	if err := copyDir(srcPath, destPath); err != nil {
		return fmt.Errorf("copying template: %w", err)
	}

	return nil
}

// ImportFromGit clones a git repository and imports it as a template.
// The repo must contain a manifest.yaml at its root.
func ImportFromGit(repoURL, name string) error {
	paths, err := config.InitPaths()
	if err != nil {
		return err
	}

	// Clone to a temp directory first.
	tmpDir, err := os.MkdirTemp("", "ct-import-")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Cloning %s...\n", repoURL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w\nMake sure git is installed and the URL is accessible", err)
	}

	// Validate the cloned repo has a manifest.yaml.
	srcManifest := filepath.Join(tmpDir, "manifest.yaml")
	m, err := LoadManifest(srcManifest)
	if err != nil {
		return fmt.Errorf("cloned repository does not have a valid manifest.yaml: %w", err)
	}

	if name == "" {
		name = m.Name
	}

	destPath := filepath.Join(paths.TemplatesDir, name)
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("template %s already exists, remove it first", name)
	}

	// Move from temp to templates dir.
	if err := copyDir(tmpDir, destPath); err != nil {
		return fmt.Errorf("installing template: %w", err)
	}

	return nil
}

// RemoveTemplate deletes a template from the store.
func RemoveTemplate(name string) error {
	paths, err := config.InitPaths()
	if err != nil {
		return err
	}

	target := filepath.Join(paths.TemplatesDir, name)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("template not found: %s", name)
	}

	return os.RemoveAll(target)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		return os.WriteFile(destPath, data, info.Mode())
	})
}
