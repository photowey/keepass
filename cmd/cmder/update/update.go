package update

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var username string
	var uri string
	var note string
	var tags []string
	var clearTags bool
	var generate bool

	cmd := &cobra.Command{
		Use:   "update <alias>",
		Short: "Update an existing entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.OutOrStdout())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			current, err := mgr.Get(masterPassword, args[0])
			if err != nil {
				return err
			}

			var input manager.UpdateInput

			if username == "" {
				updated, err := prompter.AskDefault("Username", current.Username)
				if err != nil {
					return err
				}
				if updated != current.Username {
					username = updated
				}
			}
			if username != "" {
				input.Username = &username
			}

			if uri == "" {
				updated, err := prompter.AskDefault("URI", current.URI)
				if err != nil {
					return err
				}
				if updated != current.URI {
					uri = updated
				}
			}
			if uri != "" {
				input.URI = &uri
			}

			if note == "" {
				updated, err := prompter.AskDefault("Note", current.Note)
				if err != nil {
					return err
				}
				if updated != current.Note {
					note = updated
				}
			}
			if note != "" {
				input.Note = &note
			}

			if len(tags) == 0 && !clearTags {
				tagLine, err := prompter.AskDefault("Tags (comma-separated)", joinTags(current.Tags))
				if err != nil {
					return err
				}
				if tagLine != joinTags(current.Tags) {
					parsed := common.ParseTags(tagLine)
					tags = parsed
				}
			}
			if clearTags {
				empty := []string{}
				input.Tags = &empty
			} else if len(tags) > 0 {
				input.Tags = &tags
			}

			if !generate {
				changePassword, err := prompter.AskOptional("Password action (leave blank to keep, type `manual` to replace, type `generate` to generate)")
				if err != nil {
					return err
				}

				switch changePassword {
				case "":
				case "manual":
					updated, err := prompter.AskSecretWithConfirmation("New account password", "Confirm new account password")
					if err != nil {
						return err
					}
					input.Password = &updated
				case "generate":
					generate = true
				default:
					return fmt.Errorf("unsupported password action %q", changePassword)
				}
			}

			input.GeneratePassword = generate

			entry, generated, err := mgr.Update(masterPassword, args[0], input)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Updated entry %s\n", entry.Alias)
			if err != nil {
				return err
			}

			if generated {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), "Password was generated and stored securely. Use `keepass get <alias> --reveal` to view it.")
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "new username")
	cmd.Flags().StringVar(&uri, "uri", "", "new URI")
	cmd.Flags().StringVar(&note, "note", "", "new note")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "replace tags with the provided list")
	cmd.Flags().BoolVar(&clearTags, "clear-tags", false, "clear all tags")
	cmd.Flags().BoolVar(&generate, "generate", false, "generate a new password")

	return cmd
}

func joinTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	value := tags[0]
	for i := 1; i < len(tags); i++ {
		value += "," + tags[i]
	}

	return value
}
