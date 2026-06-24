package engine

// VariableType indicates the type of a template variable.
type VariableType string

const (
	VarString VariableType = "string"
	VarBool   VariableType = "bool"
	VarChoice VariableType = "choice"
)

// VariableDef defines a single template variable.
type VariableDef struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description,omitempty"`
	Prompt      string       `yaml:"prompt,omitempty"`
	Default     string       `yaml:"default,omitempty"`
	Type        VariableType `yaml:"type,omitempty"`
	Options     []string     `yaml:"options,omitempty"`
	Required    bool         `yaml:"required,omitempty"`
	Validate    string       `yaml:"validate,omitempty"`
}

// PromptText returns the prompt string for this variable, falling back to defaults.
func (v VariableDef) PromptText() string {
	if v.Prompt != "" {
		return v.Prompt
	}
	if v.Description != "" {
		return v.Description
	}
	return v.Name
}
