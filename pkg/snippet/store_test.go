package snippet

import (
	"os"
	"path/filepath"
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	return &Store{dir: dir}
}

func TestStoreSaveLoad(t *testing.T) {
	store := testStore(t)
	s := New("test-snip", "desc", "go", "package main", []string{"go"})

	if err := store.Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := store.Load("test-snip")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Name != s.Name {
		t.Errorf("expected name %q, got %q", s.Name, loaded.Name)
	}
	if loaded.Code != s.Code {
		t.Errorf("expected code %q, got %q", s.Code, loaded.Code)
	}
}

func TestStoreLoadNotFound(t *testing.T) {
	store := testStore(t)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Error("expected error for missing snippet")
	}
}

func TestStoreDelete(t *testing.T) {
	store := testStore(t)
	s := New("temp", "", "go", "code", nil)

	if err := store.Save(s); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete("temp"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load("temp"); err == nil {
		t.Error("expected error after delete")
	}
}

func TestStoreDeleteNotFound(t *testing.T) {
	store := testStore(t)
	if err := store.Delete("nonexistent"); err == nil {
		t.Error("expected error for missing snippet")
	}
}

func TestStoreList(t *testing.T) {
	store := testStore(t)

	store.Save(New("a", "", "go", "code", nil))
	store.Save(New("b", "", "py", "code", nil))
	store.Save(New("c", "", "go", "code", nil))

	// List all.
	all, err := store.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 snippets, got %d", len(all))
	}

	// Filter by language.
	goSnippets, err := store.List("go")
	if err != nil {
		t.Fatal(err)
	}
	if len(goSnippets) != 2 {
		t.Errorf("expected 2 go snippets, got %d", len(goSnippets))
	}
}

func TestStoreListEmpty(t *testing.T) {
	store := testStore(t)
	all, err := store.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 snippets from empty dir, got %d", len(all))
	}
}

func TestStoreListNonexistentDir(t *testing.T) {
	store := &Store{dir: filepath.Join(t.TempDir(), "nonexistent")}
	all, err := store.List("")
	if err != nil {
		t.Fatal(err)
	}
	if all != nil {
		t.Error("expected nil from missing dir")
	}
}

func TestStoreSearch(t *testing.T) {
	store := testStore(t)

	store.Save(New("hello-world", "a greeting program", "go", "code", []string{"greeting"}))
	store.Save(New("goodbye", "farewell program", "py", "code", []string{"farewell"}))
	store.Save(New("helloworld2", "another one", "js", "code", nil))

	results, err := store.Search("hello")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	// Exact match should sort first.
	if results[0].Name != "hello-world" {
		t.Errorf("expected 'hello-world' first (exact prefix), got %q", results[0].Name)
	}
}

func TestStoreSaveCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "sub", "snippets")
	store := &Store{dir: dir}
	s := New("test", "", "go", "code", nil)
	if err := store.Save(s); err != nil {
		t.Fatalf("save should create dirs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "test.yaml")); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
