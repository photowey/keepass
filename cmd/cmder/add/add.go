/*
 * Copyright © 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package add

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var uri string
	var note string
	var tags []string
	var generate bool
	var revealGenerated bool

	cmd := &cobra.Command{
		Use:   "add [alias] [username]",
		Short: "Add a vault entry",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.OutOrStdout())

			alias := ""
			if len(args) > 0 {
				alias = args[0]
			} else {
				alias, err = prompter.Ask("Alias")
				if err != nil {
					return err
				}
			}

			username := ""
			if len(args) > 1 {
				username = args[1]
			} else {
				username, err = prompter.Ask("Username")
				if err != nil {
					return err
				}
			}

			if uri == "" {
				uri, err = prompter.AskOptional("URI (optional)")
				if err != nil {
					return err
				}
			}

			if note == "" {
				note, err = prompter.AskOptional("Note (optional)")
				if err != nil {
					return err
				}
			}

			if len(tags) == 0 {
				tagLine, err := prompter.AskOptional("Tags (comma-separated, optional)")
				if err != nil {
					return err
				}
				tags = common.ParseTags(tagLine)
			}

			accountPassword := ""
			if !generate {
				accountPassword, err = prompter.AskSecret("Account password (leave blank to generate)")
				if err != nil {
					return err
				}

				if accountPassword != "" {
					confirmed, err := prompter.AskSecret("Confirm account password")
					if err != nil {
						return err
					}

					if accountPassword != confirmed {
						return fmt.Errorf("account password does not match confirmation")
					}
				}
			}

			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			entry, generated, err := mgr.Add(masterPassword, manager.AddInput{
				Alias:            alias,
				Username:         username,
				Password:         accountPassword,
				URI:              uri,
				Note:             note,
				Tags:             tags,
				GeneratePassword: generate,
			})
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Added entry %s\n", entry.Alias)
			if err != nil {
				return err
			}

			if generated {
				if revealGenerated {
					_, err = fmt.Fprintf(cmd.OutOrStdout(), "Generated password: %s\n", entry.Password)
					return err
				}

				_, err = fmt.Fprintln(cmd.OutOrStdout(), "Password was generated and stored securely. Use `keepass get <alias> --reveal` to view it.")
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&uri, "uri", "", "account URI")
	cmd.Flags().StringVar(&note, "note", "", "account note")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "entry tag (repeatable)")
	cmd.Flags().BoolVar(&generate, "generate", false, "generate the account password")
	cmd.Flags().BoolVar(&revealGenerated, "reveal-generated", false, "print the generated password once after creation")

	return cmd
}
