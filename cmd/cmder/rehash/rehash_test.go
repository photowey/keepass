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
	"bytes"
	"strings"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewRejectsUnexpectedArguments(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"unexpected"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected argument validation error")
	}
}

func TestNewRehashesVault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	mgr := manager.New(env, cfg)
	if err := mgr.Initialize("master", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if _, _, err := mgr.Add("master", manager.AddInput{
		Alias:    "github",
		Username: "octocat",
		Password: "secret-123",
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Rehashed vault at ") || !strings.Contains(out.String(), "Entries: 1") {
		t.Fatalf("unexpected rehash output %q", out.String())
	}
}
