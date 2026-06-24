package snippet

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jean3690/ct-cli/pkg/config"
	"gopkg.in/yaml.v3"
)

// Store manages snippet persistence.
type Store struct {
	dir string
}

// NewStore creates a Store using the configured snippets directory.
func NewStore() (*Store, error) {
	paths, err := config.InitPaths()
	if err != nil {
		return nil, err
	}
	return &Store{dir: paths.SnippetsDir}, nil
}

// Save writes a snippet to disk.
func (s *Store) Save(snip *Snippet) error {
	path := filepath.Join(s.dir, snip.Name+".yaml")
	data, err := yaml.Marshal(snip)
	if err != nil {
		return fmt.Errorf("marshaling snippet: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snippet from disk by name.
func (s *Store) Load(name string) (*Snippet, error) {
	path := filepath.Join(s.dir, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snippet not found: %s", name)
		}
		return nil, fmt.Errorf("reading snippet: %w", err)
	}
	var snip Snippet
	if err := yaml.Unmarshal(data, &snip); err != nil {
		return nil, fmt.Errorf("parsing snippet: %w", err)
	}
	return &snip, nil
}

// Delete removes a snippet from disk.
func (s *Store) Delete(name string) error {
	path := filepath.Join(s.dir, name+".yaml")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snippet not found: %s", name)
		}
		return err
	}
	return nil
}

// List returns all snippets, optionally filtered by language.
func (s *Store) List(language string) ([]*Snippet, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var snippets []*Snippet
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		snip, err := s.Load(strings.TrimSuffix(entry.Name(), ".yaml"))
		if err != nil {
			continue
		}
		if language != "" && language != snip.Language {
			continue
		}
		snippets = append(snippets, snip)
	}

	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].Name < snippets[j].Name
	})

	return snippets, nil
}

// Search finds snippets matching a query.
func (s *Store) Search(query string) ([]*Snippet, error) {
	all, err := s.List("")
	if err != nil {
		return nil, err
	}

	var results []*Snippet
	for _, snip := range all {
		if snip.MatchesQuery(query) {
			results = append(results, snip)
		}
	}

	// Sort: exact name match first, then prefix, then substring.
	sort.Slice(results, func(i, j int) bool {
		qi := matchScore(results[i].Name, query)
		qj := matchScore(results[j].Name, query)
		return qi > qj
	})

	return results, nil
}

func matchScore(name, query string) int {
	ql := strings.ToLower(query)
	nl := strings.ToLower(name)
	if nl == ql {
		return 3
	}
	if strings.HasPrefix(nl, ql) {
		return 2
	}
	if strings.Contains(nl, ql) {
		return 1
	}
	return 0
}
