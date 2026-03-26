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

package doctor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewReportsMissingConfigAndVault(t *testing.T) {
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
	if !strings.Contains(output, `"overall_status": "warn"`) {
		t.Fatalf("expected warn status, got %s", output)
	}

	if !strings.Contains(output, `"present": false`) {
		t.Fatalf("expected missing resources in output, got %s", output)
	}
}

func TestNewPrintsHealthyTextReport(t *testing.T) {
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

	cmd := New()

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Overall: ok") || !strings.Contains(output, "Config present: true") || !strings.Contains(output, "Vault present: true") {
		t.Fatalf("unexpected healthy doctor output %q", output)
	}
}

func TestNewReportsRehashRecommendationInTextMode(t *testing.T) {
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

	cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save(updated cfg) error = %v", err)
	}

	cmd := New()

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Overall: warn") || !strings.Contains(output, "Recommendations:") || !strings.Contains(output, "keepass rehash") {
		t.Fatalf("unexpected warning doctor output %q", output)
	}
}

func TestNewRejectsUnexpectedArguments(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	cmd := New()
	cmd.SetArgs([]string{"unexpected"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected argument validation error")
	}
}

func TestNewJSONIncludesResolvedPaths(t *testing.T) {
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
	if !strings.Contains(output, filepath.Join(homeDir, "keepass.config.json")) || !strings.Contains(output, filepath.Join(homeDir, "keepass.kp")) {
		t.Fatalf("expected resolved paths in output, got %q", output)
	}

	if _, err := os.Stat(homeDir); err != nil {
		t.Fatalf("expected keepass home dir to remain accessible: %v", err)
	}
}
