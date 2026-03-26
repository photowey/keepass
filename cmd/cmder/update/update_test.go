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

package update

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestHasUpdateIntent(t *testing.T) {
	if hasUpdateIntent(manager.UpdateInput{}) {
		t.Fatal("expected empty input to have no update intent")
	}

	username := "alice"
	if !hasUpdateIntent(manager.UpdateInput{Username: &username}) {
		t.Fatal("expected username mutation to count as update intent")
	}
}

func TestJoinTags(t *testing.T) {
	if got := joinTags(nil); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}

	if got := joinTags([]string{"ops", "prod"}); got != "ops,prod" {
		t.Fatalf("unexpected joined tags %q", got)
	}
}

func TestNewRejectsConflictingURIFlags(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"github", "--uri", "https://example.com", "--clear-uri"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected conflict error")
	}

	var cliErr common.CLIError
	if !errors.As(err, &cliErr) || cliErr.ExitCode != common.ExitCodeUsage {
		t.Fatalf("expected usage error, got %v", err)
	}

	if !strings.Contains(err.Error(), "--uri and --clear-uri") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestNewRejectsConflictingPasswordFlags(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"github", "--generate", "--password", "secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected conflict error")
	}

	var cliErr common.CLIError
	if !errors.As(err, &cliErr) || cliErr.ExitCode != common.ExitCodeUsage {
		t.Fatalf("expected usage error, got %v", err)
	}

	if !strings.Contains(err.Error(), "--generate and --password") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestNewUpdatesEntryInNonInteractiveMode(t *testing.T) {
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
	cmd.Flags().Bool("non-interactive", false, "")
	if err := cmd.Flags().Set("non-interactive", "true"); err != nil {
		t.Fatalf("Set(non-interactive) error = %v", err)
	}
	cmd.SetArgs([]string{"github", "--username", "hubot", "--password", "new-secret"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated entry github") {
		t.Fatalf("unexpected update output %q", out.String())
	}

	entry, err := mgr.Get("master", "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if entry.Username != "hubot" || entry.Password != "new-secret" {
		t.Fatalf("unexpected stored entry %#v", entry)
	}
}

func TestNewClearsFieldsInNonInteractiveMode(t *testing.T) {
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
		Note:     "personal",
		Tags:     []string{"code"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.Flags().Bool("non-interactive", false, "")
	if err := cmd.Flags().Set("non-interactive", "true"); err != nil {
		t.Fatalf("Set(non-interactive) error = %v", err)
	}
	cmd.SetArgs([]string{"github", "--clear-uri", "--clear-note", "--clear-tags"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	entry, err := mgr.Get("master", "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if entry.URI != "" || entry.Note != "" || len(entry.Tags) != 0 {
		t.Fatalf("expected cleared fields, got %#v", entry)
	}
}

func TestNewRejectsConflictingNoteFlags(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"github", "--note", "personal", "--clear-note"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected conflict error")
	}

	var cliErr common.CLIError
	if !errors.As(err, &cliErr) || cliErr.ExitCode != common.ExitCodeUsage {
		t.Fatalf("expected usage error, got %v", err)
	}

	if !strings.Contains(err.Error(), "--note and --clear-note") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestNewSupportsInteractiveGenerateAction(t *testing.T) {
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
		Password: "old-secret",
		URI:      "https://github.com",
		Note:     "personal",
		Tags:     []string{"code"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"github"})
	cmd.SetIn(strings.NewReader("master\n\n\n\n\ngenerate\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Password was generated and stored securely.") {
		t.Fatalf("unexpected interactive generate output %q", out.String())
	}

	entry, err := mgr.Get("master", "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if entry.Password == "old-secret" {
		t.Fatalf("expected password to change after interactive generate, got %#v", entry)
	}
}
