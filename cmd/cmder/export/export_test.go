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

package export

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewRequiresPath(t *testing.T) {
	cmd := New()

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected usage error")
	}

	var cliErr common.CLIError
	if !errors.As(err, &cliErr) || cliErr.ExitCode != common.ExitCodeUsage {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestNewExportsEntriesToJSONFile(t *testing.T) {
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
		URI:      "https://github.com",
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	exportPath := filepath.Join(t.TempDir(), "entries.json")
	cmd := New()
	cmd.SetArgs([]string{"--path", exportPath})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if !strings.Contains(string(data), `"alias": "github"`) || !strings.Contains(out.String(), "Exported entries to ") {
		t.Fatalf("unexpected export result output=%q file=%q", out.String(), string(data))
	}
}
