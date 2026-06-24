package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type textEngine struct {
	funcs template.FuncMap
}

func (e *textEngine) AddFuncs(funcs template.FuncMap) {
	for k, v := range funcs {
		e.funcs[k] = v
	}
}

func (e *textEngine) registerDefaultFuncs() {
	e.funcs = template.FuncMap{
		"upper":  strings.ToUpper,
		"lower":  strings.ToLower,
		"title":  cases.Title(language.English).String,
		"snake":  toSnakeCase,
		"camel":  toCamelCase,
		"kebab":  toKebabCase,
		"pascal": toPascalCase,
		"now":    func(layout string) string { return time.Now().Format(layout) },
	}
}

func (e *textEngine) Render(tmplContent string, vars map[string]string) (string, error) {
	tmpl, err := template.New("inline").Funcs(e.funcs).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

func (e *textEngine) RenderFile(tmplPath, destPath string, vars map[string]string) error {
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("reading template file %s: %w", tmplPath, err)
	}

	rendered, err := e.Render(string(content), vars)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", tmplPath, err)
	}

	// Ensure destination directory exists.
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating dest dir %s: %w", destDir, err)
	}

	if err := os.WriteFile(destPath, []byte(rendered), 0644); err != nil {
		return fmt.Errorf("writing rendered file %s: %w", destPath, err)
	}

	return nil
}

// RenderString renders a template string (used for paths with variables).
func (e *textEngine) RenderString(tmplStr string, vars map[string]string) (string, error) {
	return e.Render(tmplStr, vars)
}

// --- case conversion helpers ---

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func toCamelCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(p)
		} else {
			parts[i] = cases.Title(language.English).String(p)
		}
	}
	return strings.Join(parts, "")
}

func toPascalCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.Title(p)
	}
	return strings.Join(parts, "")
}

func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

func splitWords(s string) []string {
	var parts []string
	var current []rune
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			if len(current) > 0 {
				parts = append(parts, string(current))
				current = nil
			}
		} else if unicode.IsUpper(r) && len(current) > 0 {
			parts = append(parts, string(current))
			current = []rune{unicode.ToLower(r)}
		} else {
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		parts = append(parts, string(current))
	}
	return parts
}
