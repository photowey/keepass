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
	"github.com/photowey/keepass/cmd/cmder/echo"
	"github.com/spf13/cobra"

	"github.com/photowey/keepass/internal/version"
)

var root = &cobra.Command{
	Use:     "keepass",
	Short:   "Password manager",
	Long:    "A cmd password manager implemented in Go.",
	Version: version.Now(),
	Aliases: []string{"kee", "kps", "keectl"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Welcome to keepass cmder %s~", version.Now())
	},
}

func init() {
	cobra.OnInitialize(onInit)
	root.AddCommand(add.Cmd, config.Cmd, echo.Cmd)
}

func Run() {
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
