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

package rehash

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rehash",
		Short: "Rewrite the vault using the current security parameters",
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

			entryCount, err := mgr.Rehash(masterPassword)
			if err != nil {
				return common.MapError(err)
			}

			argon2id := mgr.Config().Security.Argon2id
			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Rehashed vault at %s with Argon2id(time=%d, memory_kib=%d, threads=%d). Entries: %d\n",
				mgr.VaultPath(),
				argon2id.Time,
				argon2id.MemoryKiB,
				argon2id.Threads,
				entryCount,
			)
			return err
		},
	}

	return cmd
}
