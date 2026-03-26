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

package restore

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var path string
	var force bool

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore config and vault from a backup bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return common.UsageError("restore requires --path")
			}

			if err := manager.RestoreCurrent(path, force); err != nil {
				return common.MapError(err)
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Restored backup bundle from %s\n", path)
			return err
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "source backup directory")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing config and vault files")
	return cmd
}
