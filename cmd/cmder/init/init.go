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

package init

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the encrypted vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadOrCreateManager()
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptNewMasterPassword(prompter)
			if err != nil {
				return err
			}

			if err := mgr.Initialize(masterPassword, force); err != nil {
				return err
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Initialized vault at %s\n",
				mgr.VaultPath(),
			)
			return err
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing vault")

	return cmd
}
