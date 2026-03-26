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

package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewPrintsDefaultConfigWhenUninitialized(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	cmd := New()
	cmd.SetArgs([]string{"--json"})

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"initialized": false`) {
		t.Fatalf("expected uninitialized config output, got %s", output)
	}

	if !strings.Contains(output, `"root_dir": `) {
		t.Fatalf("expected resolved root_dir in output, got %s", output)
	}
}
