package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDir(t *testing.T) {
	dir, err := DefaultDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir == "" {
		t.Error("expected non-empty default dir")
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got %q", dir)
	}
}

func TestInitPaths(t *testing.T) {
	paths, err := InitPaths()
	if err != nil {
		t.Fatal(err)
	}
	if paths.TemplatesDir == "" || paths.SnippetsDir == "" || paths.ConfigDir == "" {
		t.Error("expected all paths to be set")
	}
	if filepath.Base(paths.ConfigFile) != "config.yaml" {
		t.Errorf("expected config.yaml, got %q", filepath.Base(paths.ConfigFile))
	}
}

func TestInitPathsCTHome(t *testing.T) {
	old := os.Getenv("CT_HOME")
	defer os.Setenv("CT_HOME", old)

	customDir := t.TempDir()
	os.Setenv("CT_HOME", customDir)

	paths, err := InitPaths()
	if err != nil {
		t.Fatal(err)
	}
	if paths.ConfigDir != customDir {
		t.Errorf("expected ConfigDir %q, got %q", customDir, paths.ConfigDir)
	}
}

func TestEnsureDirs(t *testing.T) {
	dir := t.TempDir()
	paths := &Paths{
		ConfigDir:    filepath.Join(dir, "config"),
		TemplatesDir: filepath.Join(dir, "templates"),
		SnippetsDir:  filepath.Join(dir, "snippets"),
	}

	if err := paths.EnsureDirs(); err != nil {
		t.Fatal(err)
	}

	for _, d := range []string{paths.ConfigDir, paths.TemplatesDir, paths.SnippetsDir} {
		info, err := os.Stat(d)
		if err != nil {
			t.Errorf("expected %s to exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", d)
		}
	}
}

func TestIsInteractive(t *testing.T) {
	// By default it should respect viper settings; hard to test without
	// setting up viper, but at minimum this function should not panic.
	_ = IsInteractive()
}

func TestReadWriteConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test-config.yaml")
	data := []byte("key: value\n")

	if err := WriteConfig(path, data); err != nil {
		t.Fatal(err)
	}

	read, err := ReadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(read) != string(data) {
		t.Errorf("expected %q, got %q", string(data), string(read))
	}
}
