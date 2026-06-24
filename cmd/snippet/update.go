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
	updDesc string
	updLang string
	updCode string
	updTags []string
)

var updateCmd = &cobra.Command{
	Use:     "update <name>",
	Aliases: []string{"edit"},
	Short:   "Update an existing code snippet",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snippet.NewStore()
		if err != nil {
			return err
		}

		name := args[0]
		snip, err := store.Load(name)
		if err != nil {
			return err
		}

		interactive := config.IsInteractive()

		if updCode != "" {
			if updDesc != "" {
				snip.Description = updDesc
			}
			if updLang != "" {
				snip.Language = updLang
			}
			if cmd.Flags().Changed("tags") {
				snip.Tags = updTags
			}
			snip.Code = updCode
			return store.Save(snip)
		}

		if interactive {
			return updateInteractive(store, snip)
		}

		return fmt.Errorf("use --code to provide new snippet content, or run interactively")
	},
}

func updateInteractive(store *snippet.Store, snip *snippet.Snippet) error {
	qs := []*survey.Question{
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Description:",
				Default: snip.Description,
			},
		},
		{
			Name: "language",
			Prompt: &survey.Input{
				Message: "Language:",
				Default: snip.Language,
				Help:    "e.g. go, python, javascript",
			},
		},
		{
			Name: "code",
			Prompt: &survey.Multiline{
				Message: "Edit code:",
				Default: snip.Code,
			},
			Validate: survey.Required,
		},
		{
			Name: "tags",
			Prompt: &survey.Input{
				Message: "Tags (comma-separated):",
				Default: strings.Join(snip.Tags, ", "),
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

	snip.Description = answers.Description
	snip.Language = answers.Language
	snip.Code = answers.Code

	var tags []string
	if answers.Tags != "" {
		for _, t := range strings.Split(answers.Tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}
	snip.Tags = tags

	if err := store.Save(snip); err != nil {
		return err
	}

	fmt.Printf("Snippet '%s' updated.\n", snip.Name)
	return nil
}

func init() {
	updateCmd.Flags().StringVar(&updDesc, "description", "", "new snippet description")
	updateCmd.Flags().StringVar(&updLang, "language", "", "new programming language")
	updateCmd.Flags().StringVar(&updCode, "code", "", "new snippet code content")
	updateCmd.Flags().StringSliceVar(&updTags, "tags", nil, "new comma-separated tags")
}
