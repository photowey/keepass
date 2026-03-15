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

package root

import (
	"fmt"
	"os"

	"github.com/photowey/keepass/cmd/cmder/add"
	"github.com/photowey/keepass/cmd/cmder/config"
	"github.com/photowey/keepass/cmd/cmder/del"
	"github.com/photowey/keepass/cmd/cmder/get"
	"github.com/photowey/keepass/cmd/cmder/initcmd"
	"github.com/photowey/keepass/cmd/cmder/list"
	"github.com/photowey/keepass/cmd/cmder/update"
	"github.com/spf13/cobra"

	"github.com/photowey/keepass/internal/version"
)

func NewCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "keepass",
		Short:         "Secure local password manager",
		Long:          "A secure local password manager backed by a versioned encrypted vault.",
		Version:       version.Now(),
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
		del.New(),
		config.New(),
	)

	return root
}

func Run() {
	if err := NewCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
