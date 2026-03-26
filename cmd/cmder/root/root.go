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

package root

import (
	"errors"
	"fmt"
	"os"

	"github.com/photowey/keepass/cmd/cmder/add"
	"github.com/photowey/keepass/cmd/cmder/audit"
	"github.com/photowey/keepass/cmd/cmder/backup"
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/cmd/cmder/completion"
	"github.com/photowey/keepass/cmd/cmder/config"
	"github.com/photowey/keepass/cmd/cmder/delete"
	"github.com/photowey/keepass/cmd/cmder/doctor"
	"github.com/photowey/keepass/cmd/cmder/export"
	"github.com/photowey/keepass/cmd/cmder/get"
	"github.com/photowey/keepass/cmd/cmder/importdata"
	initcmd "github.com/photowey/keepass/cmd/cmder/init"
	"github.com/photowey/keepass/cmd/cmder/list"
	"github.com/photowey/keepass/cmd/cmder/rehash"
	"github.com/photowey/keepass/cmd/cmder/restore"
	"github.com/photowey/keepass/cmd/cmder/rotate"
	"github.com/photowey/keepass/cmd/cmder/update"
	"github.com/photowey/keepass/internal/version"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var nonInteractive bool

	root := &cobra.Command{
		Use:           "keepass",
		Short:         "Secure local password manager",
		Long:          "A secure local password manager backed by a versioned encrypted vault.",
		Version:       version.Summary(),
		Aliases:       []string{"kee", "kps"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	root.AddCommand(
		initcmd.New(),
		add.New(),
		list.New(),
		get.New(),
		update.New(),
		audit.New(),
		rotate.New(),
		export.New(),
		importdata.New(),
		backup.New(),
		restore.New(),
		delete.New(),
		doctor.New(),
		rehash.New(),
		config.New(),
	)

	root.AddCommand(completion.New(root))

	root.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "disable interactive prompts and fail fast when required input is missing")

	return root
}

func Run() {
	if err := NewCommand().Execute(); err != nil {
		var cliErr common.CLIError
		if errors.As(err, &cliErr) && cliErr.ExitCode != 0 {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(cliErr.ExitCode)
		}

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(common.ExitCodeGeneric)
	}
}
