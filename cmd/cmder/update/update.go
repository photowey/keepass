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
	var passwordValue string
	var tags []string
	var clearTags bool
	var clearURI bool
	var clearNote bool
	var generate bool

	cmd := &cobra.Command{
		Use:   "update <alias>",
		Short: "Update an existing entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if clearURI && uri != "" {
				return common.UsageError("cannot use --uri and --clear-uri together")
			}

			if clearNote && note != "" {
				return common.UsageError("cannot use --note and --clear-note together")
			}

			if generate && passwordValue != "" {
				return common.UsageError("cannot use --generate and --password together")
			}

			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			in := cmd.InOrStdin()
			prompter := common.NewPrompter(in, cmd.ErrOrStderr())
			interactive := common.IsInteractive(in) && !common.IsNonInteractive(cmd)
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			current, err := mgr.Get(masterPassword, args[0])
			if err != nil {
				return common.MapError(err)
			}

			var input manager.UpdateInput

			if username == "" && interactive {
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

			if clearURI {
				uri = ""
				input.URI = &uri
			} else if uri == "" && interactive {
				updated, err := prompter.AskDefault("URI", current.URI)
				if err != nil {
					return err
				}
				if updated != current.URI {
					uri = updated
				}
			}
			if !clearURI && uri != "" {
				input.URI = &uri
			}

			if clearNote {
				note = ""
				input.Note = &note
			} else if note == "" && interactive {
				updated, err := prompter.AskDefault("Note", current.Note)
				if err != nil {
					return err
				}
				if updated != current.Note {
					note = updated
				}
			}
			if !clearNote && note != "" {
				input.Note = &note
			}

			if len(tags) == 0 && !clearTags && interactive {
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

			if passwordValue != "" {
				input.Password = &passwordValue
			} else if !generate && interactive {
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

			if !interactive && !hasUpdateIntent(input) {
				return common.UsageError("update requires at least one mutation flag in non-interactive mode")
			}

			entry, generated, err := mgr.Update(masterPassword, args[0], input)
			if err != nil {
				return common.MapError(err)
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
	cmd.Flags().StringVar(&passwordValue, "password", "", "new account password")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "replace tags with the provided list")
	cmd.Flags().BoolVar(&clearTags, "clear-tags", false, "clear all tags")
	cmd.Flags().BoolVar(&clearURI, "clear-uri", false, "clear the URI")
	cmd.Flags().BoolVar(&clearNote, "clear-note", false, "clear the note")
	cmd.Flags().BoolVar(&generate, "generate", false, "generate a new password")

	return cmd
}

func hasUpdateIntent(input manager.UpdateInput) bool {
	return input.Username != nil ||
		input.Password != nil ||
		input.URI != nil ||
		input.Note != nil ||
		input.Tags != nil ||
		input.GeneratePassword
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
