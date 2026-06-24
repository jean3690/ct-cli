package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderSimple(t *testing.T) {
	eng := New()
	result, err := eng.Render("Hello {{ .Name }}!", map[string]string{"Name": "World"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Hello World!" {
		t.Errorf("expected 'Hello World!', got %q", result)
	}
}

func TestRenderWithFuncs(t *testing.T) {
	eng := New()
	result, err := eng.Render("{{ .Name | upper }}", map[string]string{"Name": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", result)
	}
}

func TestRenderFile(t *testing.T) {
	eng := New()

	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "test.go.tmpl")
	destPath := filepath.Join(dir, "output/test.go")

	if err := os.WriteFile(tmplPath, []byte("package {{ .Pkg }}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := eng.RenderFile(tmplPath, destPath, map[string]string{"Pkg": "main"}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "package main\n" {
		t.Errorf("expected 'package main\\n', got %q", string(data))
	}
}

func TestRenderString(t *testing.T) {
	eng := New()
	result, err := eng.RenderString("{{ .Dir }}/{{ .File }}", map[string]string{"Dir": "src", "File": "main.go"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "src/main.go" {
		t.Errorf("expected 'src/main.go', got %q", result)
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"ProjectName", "project_name"},
		{"HTTPServer", "http_server"},
		{"simple", "simple"},
		{"XMLParser", "xml_parser"},
	}
	for _, tt := range tests {
		got := toSnakeCase(tt.in)
		if got != tt.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct{ in, want string }{
		{"hello_world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
		{"hello-world", "helloWorld"},
	}
	for _, tt := range tests {
		got := toCamelCase(tt.in)
		if got != tt.want {
			t.Errorf("toCamelCase(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestPascalCase(t *testing.T) {
	got := toPascalCase("hello_world")
	if got != "HelloWorld" {
		t.Errorf("expected 'HelloWorld', got %q", got)
	}
}

func TestVariableDefPromptText(t *testing.T) {
	v := VariableDef{Name: "Foo", Prompt: "What is foo?"}
	if v.PromptText() != "What is foo?" {
		t.Error("expected explicit prompt")
	}

	v2 := VariableDef{Name: "Bar", Description: "Bar description"}
	if v2.PromptText() != "Bar description" {
		t.Error("expected description fallback")
	}

	v3 := VariableDef{Name: "Baz"}
	if v3.PromptText() != "Baz" {
		t.Error("expected name fallback")
	}
}

func TestRenderPathWithVars(t *testing.T) {
	eng := New()
	result, err := eng.RenderString("{{ .ProjectName }}/src/main.go", map[string]string{"ProjectName": "myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "myapp") {
		t.Errorf("expected path to contain 'myapp', got %q", result)
	}
}
