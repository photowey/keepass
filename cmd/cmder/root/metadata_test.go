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

import "testing"

func TestNewCommandRegistersAliasesAndSubcommands(t *testing.T) {
	cmd := NewCommand()

	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "kee" || cmd.Aliases[1] != "kps" {
		t.Fatalf("unexpected aliases %#v", cmd.Aliases)
	}

	if cmd.PersistentFlags().Lookup("non-interactive") == nil {
		t.Fatal("expected non-interactive persistent flag")
	}

	expected := []string{
		"init", "add", "list", "get", "update", "audit", "rotate",
		"export", "import", "backup", "restore", "delete", "doctor",
		"rehash", "config", "completion",
	}

	commands := cmd.Commands()
	for _, use := range expected {
		found := false
		for _, sub := range commands {
			if sub.Name() == use {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected subcommand %q to be registered", use)
		}
	}
}
