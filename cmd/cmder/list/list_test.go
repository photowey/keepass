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

package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewRejectsMoreThanOneQueryArgument(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{"first", "second"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected argument validation error")
	}
}

func TestNewListsEntriesAsJSON(t *testing.T) {
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
		Tags:     []string{"code"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"--json", "--tag", "code"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"alias": "github"`) || strings.Contains(output, "secret-123") {
		t.Fatalf("unexpected list output %q", output)
	}
}

func TestNewFiltersEntriesByQueryAndRendersText(t *testing.T) {
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

	for _, input := range []manager.AddInput{
		{Alias: "github", Username: "octocat", Password: "secret-1", URI: "https://github.com", Tags: []string{"code"}},
		{Alias: "gitlab", Username: "gitlab-user", Password: "secret-2", URI: "https://gitlab.com", Tags: []string{"work"}},
	} {
		if _, _, err := mgr.Add("master", input); err != nil {
			t.Fatalf("Add(%s) error = %v", input.Alias, err)
		}
	}

	cmd := New()
	cmd.SetArgs([]string{"hub"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "github") || strings.Contains(output, "gitlab") || strings.Contains(output, "secret-1") {
		t.Fatalf("unexpected list text output %q", output)
	}
}

func TestNewReturnsEmptyJSONArrayWhenNoEntriesMatch(t *testing.T) {
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
		Tags:     []string{"code"},
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"--json", "--tag", "ops"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if strings.TrimSpace(out.String()) != "[]" {
		t.Fatalf("expected empty json array, got %q", out.String())
	}
}
