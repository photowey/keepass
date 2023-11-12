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

package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd `keepass config` cmd
//
// update an exits password node of $username with alias,
// and the alias must be unique under the username $username.
//
// pattern: $ keepass config -a $alias -u $username -p $password ...
//
// e.g.:
//
// $ keepass config photowey -a github -u photowey@github.com -p hello.github
//
// update a password node (-u photowey@github.com -p hello.github)
// with the username photowey@github.com and password hello.github under the namespace of username photowey, using the alias github.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Config keepass password node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hello keepass config~")
	},
}

func init() {

}
