package snippet

import (
	"strings"
	"time"
)

// Snippet represents a stored code snippet.
type Snippet struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Language    string   `yaml:"language"`
	Tags        []string `yaml:"tags,omitempty"`
	Code        string   `yaml:"code"`
	CreatedAt   string   `yaml:"createdAt"`
	UpdatedAt   string   `yaml:"updatedAt"`
}

// New creates a new Snippet with timestamps set.
func New(name, description, language, code string, tags []string) *Snippet {
	now := time.Now().Format(time.RFC3339)
	return &Snippet{
		Name:        name,
		Description: description,
		Language:    language,
		Tags:        tags,
		Code:        code,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MatchesQuery checks if the snippet matches a search query.
func (s *Snippet) MatchesQuery(query string) bool {
	ql := strings.ToLower(query)
	if strings.ToLower(s.Name) == ql {
		return true // Exact name match.
	}
	if strings.Contains(strings.ToLower(s.Name), ql) {
		return true
	}
	if strings.Contains(strings.ToLower(s.Description), ql) {
		return true
	}
	for _, t := range s.Tags {
		if strings.Contains(strings.ToLower(t), ql) {
			return true
		}
	}
	return false
}
