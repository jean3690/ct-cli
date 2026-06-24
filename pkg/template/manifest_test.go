package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "manifest.yaml")
	data := `apiVersion: v1
name: test-tmpl
version: "1.0"
description: A test template
variables:
  - name: ProjectName
    type: string
    required: true
  - name: UseDB
    type: bool
    default: "false"
files:
  - source: main.go.tmpl
    dest: "{{ .ProjectName }}/main.go"
`
	if err := saveManifestContent(manifestPath, []byte(data)); err != nil {
		t.Fatal(err)
	}

	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}

	if m.Name != "test-tmpl" {
		t.Errorf("expected name 'test-tmpl', got %q", m.Name)
	}
	if m.Version != "1.0" {
		t.Errorf("expected version '1.0', got %q", m.Version)
	}
	if len(m.Variables) != 2 {
		t.Errorf("expected 2 variables, got %d", len(m.Variables))
	}
	if len(m.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(m.Files))
	}
	if m.Files[0].OutputDest() != "{{ .ProjectName }}/main.go" {
		t.Errorf("unexpected dest: %s", m.Files[0].OutputDest())
	}
}

func TestValidateMissingName(t *testing.T) {
	m := &TemplateManifest{APIVersion: "v1", Version: "1.0"}
	if err := m.Validate(); err == nil {
		t.Error("expected validation error for missing name")
	}
}

func TestHasScaffold(t *testing.T) {
	m := &TemplateManifest{ScaffoldDir: "files/skeleton"}
	if !m.HasScaffold() {
		t.Error("expected HasScaffold to return true")
	}

	m2 := &TemplateManifest{}
	if m2.HasScaffold() {
		t.Error("expected HasScaffold to return false")
	}
}

func saveManifestContent(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
