package template

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jean3690/ct-cli/pkg/engine"
	"gopkg.in/yaml.v3"
)

// ManifestAPIVersion is the current API version.
const ManifestAPIVersion = "v1"

// TemplateManifest describes a single template.
type TemplateManifest struct {
	APIVersion   string               `yaml:"apiVersion"`
	Name         string               `yaml:"name"`
	Description  string               `yaml:"description"`
	Version      string               `yaml:"version"`
	Author       string               `yaml:"author,omitempty"`
	Languages    []string             `yaml:"languages,omitempty"`
	Variables    []engine.VariableDef `yaml:"variables"`
	Files        []FileMapping        `yaml:"files"`
	ScaffoldDir  string               `yaml:"scaffoldDir,omitempty"`
	PostGenerate string               `yaml:"postGenerate,omitempty"`
}

// FileMapping maps a template source file to its output destination.
type FileMapping struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest,omitempty"`
}

// LoadManifest reads and parses a manifest.yaml file.
func LoadManifest(manifestPath string) (*TemplateManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}

	var m TemplateManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}

	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("validating manifest: %w", err)
	}

	return &m, nil
}

// SaveManifest writes a manifest to a YAML file.
func SaveManifest(m *TemplateManifest, path string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating manifest dir: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Validate checks that the manifest has required fields.
func (m *TemplateManifest) Validate() error {
	if m.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}

	seen := make(map[string]bool)
	for i, v := range m.Variables {
		if v.Name == "" {
			return fmt.Errorf("variable[%d]: name is required", i)
		}
		if seen[v.Name] {
			return fmt.Errorf("duplicate variable name: %s", v.Name)
		}
		seen[v.Name] = true

		if v.Type == "" {
			m.Variables[i].Type = engine.VarString
		}
		if v.Type == engine.VarChoice && len(v.Options) == 0 {
			return fmt.Errorf("variable %s: type 'choice' requires options", v.Name)
		}
	}

	for i, f := range m.Files {
		if f.Source == "" {
			return fmt.Errorf("file[%d]: source is required", i)
		}
	}

	return nil
}

// HasScaffold returns true if this is a scaffold (multi-file project) template.
func (m *TemplateManifest) HasScaffold() bool {
	return m.ScaffoldDir != ""
}

// OutputDest returns the destination path for a file mapping, or a default derived from source.
func (f FileMapping) OutputDest() string {
	if f.Dest != "" {
		return f.Dest
	}
	// Strip .tmpl extension if present.
	return f.Source[:len(f.Source)-len(".tmpl")]
}
