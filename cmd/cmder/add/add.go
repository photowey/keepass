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
	"github.com/photowey/keepass/internal/action"
	"github.com/photowey/keepass/internal/types"
	"github.com/photowey/keepass/pkg/stringz"
	"github.com/spf13/cobra"
)

// Cmd `keepass add` cmd
//
// Add a password node into the namespace(master username) of $username.
//
// Pattern: $ keepass add -n $namespace -a $alias -u $username -p $password ...
//
// or
//
// Pattern: $ keepass add $namespace -a $alias -u $username -p $password ...
//
// e.g.:
//
// $ keepass add -n photowey -a github -u photowey@github.com -p hello.github -i https://github.com -t "the username and password of github.com website".
//
// or
//
// $ keepass add photowey -a github -u photowey@github.com -p hello.github -i https://github.com -t "the username and password of github.com website".
//
// or:
//
// $ keepass add < testdata.json
//
// Json example:
//
//	{
//	 "namespace": "photowey",
//	 "nodes": [
//	   {
//	     "alias": "github",
//	     "username": "photowey@github.com",
//	     "password": "hello.github",
//	     "uri": "https://github.com",
//	     "note": "the username and password of github.com website"
//	   }
//	 ]
//	}
//
// -n: namespace
//
// -a: alias: the alias must be unique under the namespace of $username.
//
// -p: password
//
// -i: uri
//
// -t: note
//
// -
//
// Add a password node (-u photowey@github.com -p hello.github)
// with the username photowey@github.com and password hello.github under the namespace of username photowey, using the alias github.
var Cmd = &cobra.Command{
	Use:   "add",
	Short: "Add keepass password node",
	Run: func(cmd *cobra.Command, args []string) {
		ns := namespace
		if stringz.IsBlankString(ns) {
			if len(args) > 0 {
				ns = args[0]
			}
		}

		if stringz.IsBlankString(ns) {
			panic("The namespace(master username) parameter cannot be blank. Either of the two commands, " +
				"`$ keepass add $username` or `$ keepass add -n $username` can be used.")
		}

		action.OnAddEvent(&types.AddEvent{
			Namespace: ns,
			Nodes: []*types.Node{{
				Alias:    alias,
				Username: username,
				Password: password,
				Uri:      uri,
				Note:     note,
			}},
		})
	},
}

var (
	namespace string
	alias     string
	username  string
	password  string
	uri       string
	note      string
)

func init() {
	Cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "the namespace of local username")
	Cmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "the alias of password node")
	Cmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username")
	Cmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
	Cmd.PersistentFlags().StringVarP(&uri, "uri", "i", "", "uri")
	Cmd.PersistentFlags().StringVarP(&note, "note", "t", "", "the note of password node(`pn`)")
}
