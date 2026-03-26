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

package audit

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var maxAgeDays int
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit credential hygiene inside the unlocked vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			report, err := mgr.AuditCredentials(masterPassword, maxAgeDays)
			if err != nil {
				return common.MapError(err)
			}

			if jsonOut {
				return common.PrintCredentialAuditJSON(cmd.OutOrStdout(), report)
			}

			_, err = common.PrintCredentialAuditText(cmd.OutOrStdout(), report)
			return err
		},
	}

	cmd.Flags().IntVar(&maxAgeDays, "max-password-age-days", 180, "mark passwords older than this number of days as stale")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}
