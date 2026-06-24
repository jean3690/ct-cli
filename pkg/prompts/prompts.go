package prompts

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jean/codeTemplateCli/pkg/engine"
)

// CollectVariables prompts the user interactively for any missing required/optional variables.
func CollectVariables(defs []engine.VariableDef, provided map[string]string, interactive bool) (map[string]string, error) {
	result := make(map[string]string)
	for k, v := range provided {
		result[k] = v
	}

	for _, def := range defs {
		if _, ok := result[def.Name]; ok {
			continue // Already provided via flags.
		}

		if !interactive {
			if def.Required {
				return nil, fmt.Errorf("variable %q is required but not provided (use --var %s=<value> or interactive mode)", def.Name, def.Name)
			}
			if def.Default != "" {
				result[def.Name] = def.Default
			}
			continue
		}

		// Interactive prompt.
		val, err := promptVariable(def)
		if err != nil {
			return nil, fmt.Errorf("prompting for %s: %w", def.Name, err)
		}
		result[def.Name] = val
	}

	return result, nil
}

func promptVariable(def engine.VariableDef) (string, error) {
	promptText := def.PromptText()
	if def.Required && def.Default == "" {
		promptText += " (required)"
	}

	switch def.Type {
	case engine.VarBool:
		var answer bool
		defaultVal := def.Default == "true"
		if err := survey.AskOne(&survey.Confirm{
			Message: promptText,
			Default: defaultVal,
		}, &answer); err != nil {
			return "", err
		}
		if answer {
			return "true", nil
		}
		return "false", nil

	case engine.VarChoice:
		var answer string
		if err := survey.AskOne(&survey.Select{
			Message: promptText,
			Options: def.Options,
			Default: def.Default,
		}, &answer); err != nil {
			return "", err
		}
		return answer, nil

	default: // VarString
		var answer string
		opts := []survey.AskOpt{}
		if def.Required {
			opts = append(opts, survey.WithValidator(survey.Required))
		}
		if def.Validate != "" {
			re, err := regexp.Compile(def.Validate)
			if err != nil {
				return "", fmt.Errorf("invalid validation regex: %w", err)
			}
			opts = append(opts, survey.WithValidator(func(val interface{}) error {
				if str, ok := val.(string); ok {
					if !re.MatchString(str) {
						return fmt.Errorf("value must match pattern: %s", def.Validate)
					}
				}
				return nil
			}))
		}
		if err := survey.AskOne(&survey.Input{
			Message: promptText,
			Default: def.Default,
		}, &answer, opts...); err != nil {
			return "", err
		}
		return answer, nil
	}
}
