package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredTemplate is a template found on disk with brief info.
type DiscoveredTemplate struct {
	Name        string
	Description string
	Version     string
	Languages   []string
	Path        string
}

// Discover scans a directory for template subdirectories containing manifest.yaml.
// If language is non-empty, only templates supporting that language are returned.
func Discover(templatesDir, language string) ([]DiscoveredTemplate, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading templates dir: %w", err)
	}

	var templates []DiscoveredTemplate
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(templatesDir, entry.Name(), "manifest.yaml")
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		m, err := LoadManifest(manifestPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping broken manifest %s: %v\n", manifestPath, err)
			continue
		}

		// Filter by language if specified.
		if language != "" && !m.supportsLanguage(language) {
			continue
		}

		templates = append(templates, DiscoveredTemplate{
			Name:        m.Name,
			Description: m.Description,
			Version:     m.Version,
			Languages:   m.Languages,
			Path:        filepath.Join(templatesDir, entry.Name()),
		})
	}

	return templates, nil
}

// supportsLanguage checks if the template supports a given language.
func (m *TemplateManifest) supportsLanguage(lang string) bool {
	if len(m.Languages) == 0 {
		return false // No languages declared = not language-specific, skip when filtering.
	}
	for _, l := range m.Languages {
		if strings.EqualFold(l, lang) {
			return true
		}
	}
	return false
}

// FindTemplate looks for a template by name in the templates directory.
func FindTemplate(templatesDir, name string) (*DiscoveredTemplate, error) {
	templates, err := Discover(templatesDir, "")
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		if t.Name == name {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("template not found: %s", name)
}
