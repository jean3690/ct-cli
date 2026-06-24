package snippet

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jean/codeTemplateCli/pkg/config"
	"github.com/jean/codeTemplateCli/pkg/snippet"
	"github.com/spf13/cobra"
)

var (
	snipName        string
	snipDesc        string
	snipLang        string
	snipCode        string
	snipTags        []string
)

var createCmd = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{"add"},
	Short:   "Create a new code snippet",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		name := args[0]
		interactive := config.IsInteractive()

		if snipCode != "" {
			// Non-interactive mode: use flags.
			s := snippet.New(name, snipDesc, snipLang, snipCode, snipTags)
			return store.Save(s)
		}

		if interactive {
			return createInteractive(store, name)
		}

		return fmt.Errorf("use --code to provide snippet content, or run interactively")
	},
}

func createInteractive(store *snippet.Store, name string) error {
	qs := []*survey.Question{
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Description:",
			},
		},
		{
			Name: "language",
			Prompt: &survey.Input{
				Message: "Language:",
				Help:    "e.g. go, python, javascript",
			},
		},
		{
			Name: "code",
			Prompt: &survey.Multiline{
				Message: "Paste your code (press <Esc> then <Enter> when done):",
			},
			Validate: survey.Required,
		},
		{
			Name: "tags",
			Prompt: &survey.Input{
				Message: "Tags (comma-separated):",
			},
		},
	}

	answers := struct {
		Description string
		Language    string
		Code        string
		Tags        string
	}{}

	if err := survey.Ask(qs, &answers); err != nil {
		return fmt.Errorf("prompt failed: %w", err)
	}

	var tags []string
	if answers.Tags != "" {
		for _, t := range strings.Split(answers.Tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	s := snippet.New(name, answers.Description, answers.Language, answers.Code, tags)
	if err := store.Save(s); err != nil {
		return err
	}

	fmt.Printf("Snippet '%s' created.\n", name)
	return nil
}

func init() {
	createCmd.Flags().StringVar(&snipDesc, "description", "", "snippet description")
	createCmd.Flags().StringVar(&snipLang, "language", "", "programming language (go, py, js, etc.)")
	createCmd.Flags().StringVar(&snipCode, "code", "", "snippet code content")
	createCmd.Flags().StringSliceVar(&snipTags, "tags", nil, "comma-separated tags")
}

