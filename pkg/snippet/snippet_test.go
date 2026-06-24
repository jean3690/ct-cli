package snippet

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New("my-snippet", "a useful function", "go", "func main() {}", []string{"go", "example"})
	if s.Name != "my-snippet" {
		t.Errorf("expected name 'my-snippet', got %q", s.Name)
	}
	if s.Description != "a useful function" {
		t.Errorf("expected description, got %q", s.Description)
	}
	if s.Language != "go" {
		t.Errorf("expected language 'go', got %q", s.Language)
	}
	if s.Code != "func main() {}" {
		t.Errorf("expected code, got %q", s.Code)
	}
	if len(s.Tags) != 2 || s.Tags[0] != "go" || s.Tags[1] != "example" {
		t.Errorf("unexpected tags: %v", s.Tags)
	}
	if s.CreatedAt == "" || s.UpdatedAt == "" {
		t.Error("expected timestamps to be set")
	}
}

func TestMatchesQuery(t *testing.T) {
	s := New("HelloWorld", "a greeting snippet", "go", "fmt.Println(\"hello\")", []string{"greeting", "hello"})

	if !s.MatchesQuery("HelloWorld") {
		t.Error("expected exact name match")
	}
	if !s.MatchesQuery("hello") {
		t.Error("expected case-insensitive name match")
	}
	if !s.MatchesQuery("greeting") {
		t.Error("expected description match")
	}
	if !s.MatchesQuery("GREETING") {
		t.Error("expected case-insensitive description match")
	}
	if !s.MatchesQuery("hello") {
		t.Error("expected tag match")
	}
	if s.MatchesQuery("nonexistent") {
		t.Error("expected no match for nonexistent query")
	}
}
