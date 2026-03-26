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
	"bytes"
	"strings"
	"testing"
)

func TestNewRegistersForceFlag(t *testing.T) {
	cmd := New()

	if cmd.Flags().Lookup("force") == nil {
		t.Fatal("expected force flag to be registered")
	}
}

func TestNewInitializesVault(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	cmd := New()
	cmd.SetIn(strings.NewReader("master\nmaster\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Initialized vault at ") {
		t.Fatalf("unexpected init output %q", out.String())
	}
}

func TestNewForceReinitializesExistingVault(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	first := New()
	first.SetIn(strings.NewReader("master\nmaster\n"))
	if err := first.Execute(); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	second := New()
	second.SetArgs([]string{"--force"})
	second.SetIn(strings.NewReader("master2\nmaster2\n"))

	var out bytes.Buffer
	second.SetOut(&out)

	if err := second.Execute(); err != nil {
		t.Fatalf("forced Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Initialized vault at ") {
		t.Fatalf("unexpected force init output %q", out.String())
	}
}
