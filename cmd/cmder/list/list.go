/*
 * Copyright © 2023-present the keepass authors. All rights reserved.
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

package list

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var tags []string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "list [query]",
		Short: "List entries",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			query := ""
			if len(args) == 1 {
				query = args[0]
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			entries, err := mgr.List(masterPassword, manager.ListFilter{
				Query: query,
				Tags:  tags,
			})
			if err != nil {
				return common.MapError(err)
			}

			if jsonOut {
				return common.PrintEntriesJSON(cmd.OutOrStdout(), entries)
			}

			common.PrintEntries(cmd.OutOrStdout(), entries)
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&tags, "tag", nil, "filter by tag (repeatable, all tags must match)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	return cmd
}
