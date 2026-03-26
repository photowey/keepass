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

package manager

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/transfer"
	"github.com/photowey/keepass/internal/vault"
)

func newTestManager(t *testing.T) (*Manager, string) {
	t.Helper()

	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	mgr := New(env, cfg)
	mgr.now = func() time.Time {
		return time.Unix(1_700_000_000, 0).UTC()
	}

	const masterPassword = "master-password"
	if err := mgr.Initialize(masterPassword, false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	return mgr, masterPassword
}

func TestAddGetAndUniquePrefixResolution(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	entry, generated, err := mgr.Add(masterPassword, AddInput{
		Alias:    "github",
		Username: "hellopass",
		Password: "github-secret-2024",
		Tags:     []string{"code"},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if generated {
		t.Fatal("expected manual password, got generated")
	}

	got, err := mgr.Get(masterPassword, "git")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Alias != entry.Alias {
		t.Fatalf("expected alias %s, got %s", entry.Alias, got.Alias)
	}
}

func TestGetRejectsAmbiguousPrefix(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "hellopass", Password: "github-secret-alpha"})
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "gitea", Username: "hellopass", Password: "gitea-secret-beta"})

	_, err := mgr.Get(masterPassword, "gi")
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous alias error, got %v", err)
	}
}

func TestListFiltersByTags(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "mysql-prod", Username: "service-admin", Password: "mysql-prod-secret", Tags: []string{"database", "prod"}})
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "blog", Username: "content-author", Password: "blog-editor-secret", Tags: []string{"blog"}})

	entries, err := mgr.List(masterPassword, ListFilter{Tags: []string{"database", "prod"}})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 1 || entries[0].Alias != "mysql-prod" {
		t.Fatalf("unexpected filtered entries: %+v", entries)
	}
}

func TestAddRejectsDuplicateAlias(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "hellopass", Password: "github-secret-alpha"})

	_, _, err := mgr.Add(masterPassword, AddInput{Alias: "github", Username: "backup-user", Password: "gitea-secret-beta"})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate alias error, got %v", err)
	}
}

func TestUpdateGenerateAndDelete(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "hellopass", Password: "github-secret-alpha"})

	note := "personal"
	uri := "https://github.com"
	tags := []string{"code"}
	entry, generated, err := mgr.Update(masterPassword, "github", UpdateInput{
		Note:             &note,
		URI:              &uri,
		Tags:             &tags,
		GeneratePassword: true,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if !generated {
		t.Fatal("expected generated password")
	}

	if len(entry.Password) != mgr.cfg.PasswordGenerator.DefaultLength {
		t.Fatalf("expected generated password length %d, got %d", mgr.cfg.PasswordGenerator.DefaultLength, len(entry.Password))
	}

	if entry.Note != note || entry.URI != uri || len(entry.Tags) != 1 || entry.Tags[0] != "code" {
		t.Fatalf("unexpected updated entry: %+v", entry)
	}

	if _, err := mgr.Delete(masterPassword, "github"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := mgr.Get(masterPassword, "github"); err == nil {
		t.Fatal("expected missing entry after delete")
	}
}

func TestRehashRewritesVaultWithCurrentSecurityParameters(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, err := mgr.Add(masterPassword, AddInput{
		Alias:    "github",
		Username: "hellopass",
		Password: "github-secret-2024",
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	mgr.cfg.Security.Argon2id.Time = 2
	mgr.cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	mgr.cfg.Security.Argon2id.Threads = 2
	mgr.cfg.Security.Argon2id.KeyLength = 32
	mgr.store = vault.NewStore(mgr.cfg.ResolveVaultPath(mgr.env), mgr.cfg)

	count, err := mgr.Rehash(masterPassword)
	if err != nil {
		t.Fatalf("Rehash() error = %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 entry after rehash, got %d", count)
	}

	header, err := readV1Header(mgr.VaultPath())
	if err != nil {
		t.Fatalf("readV1Header() error = %v", err)
	}

	if header.Argon2id.Time != mgr.cfg.Security.Argon2id.Time {
		t.Fatalf("expected time %d, got %d", mgr.cfg.Security.Argon2id.Time, header.Argon2id.Time)
	}

	if header.Argon2id.MemoryKiB != mgr.cfg.Security.Argon2id.MemoryKiB {
		t.Fatalf("expected memory_kib %d, got %d", mgr.cfg.Security.Argon2id.MemoryKiB, header.Argon2id.MemoryKiB)
	}

	entry, err := mgr.Get(masterPassword, "github")
	if err != nil {
		t.Fatalf("Get() after rehash error = %v", err)
	}

	if entry.Password != "github-secret-2024" {
		t.Fatalf("expected preserved password after rehash, got %q", entry.Password)
	}
}

func TestExportImportAndRotate(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, err := mgr.Add(masterPassword, AddInput{
		Alias:    "github",
		Username: "hellopass",
		Password: "github-secret-2024",
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	doc, err := mgr.Export(masterPassword)
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	rotated, _, err := mgr.Rotate(masterPassword, "github", stringPtr("github-rotated-2024"), false)
	if err != nil {
		t.Fatalf("Rotate() error = %v", err)
	}

	if rotated.Password != "github-rotated-2024" {
		t.Fatalf("expected rotated password, got %q", rotated.Password)
	}

	doc.Entries[0].Password = "github-imported-2024"
	result, err := mgr.Import(masterPassword, doc, "overwrite")
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	if result.Overwrote != 1 {
		t.Fatalf("expected 1 overwritten entry, got %+v", result)
	}

	entry, err := mgr.Get(masterPassword, "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if entry.Password != "github-imported-2024" {
		t.Fatalf("expected imported password, got %q", entry.Password)
	}

	report, err := mgr.AuditCredentials(masterPassword, 1)
	if err != nil {
		t.Fatalf("AuditCredentials() error = %v", err)
	}

	if report.OverallStatus == "" {
		t.Fatalf("expected non-empty audit report, got %+v", report)
	}
}

func TestAddRejectsInvalidInput(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, err := mgr.Add(masterPassword, AddInput{Alias: "GitHub", Username: "", Password: "secret"})
	if err == nil {
		t.Fatal("expected invalid input error")
	}
}

func TestUpdateRejectsBlankPasswordAndUsername(t *testing.T) {
	mgr, masterPassword := newTestManager(t)
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "hellopass", Password: "secret"})

	blank := "   "
	if _, _, err := mgr.Update(masterPassword, "github", UpdateInput{Username: &blank}); err == nil {
		t.Fatal("expected blank username error")
	}
	if _, _, err := mgr.Update(masterPassword, "github", UpdateInput{Password: &blank}); err == nil {
		t.Fatal("expected blank password error")
	}
}

func TestImportSupportsSkipAndFailStrategies(t *testing.T) {
	mgr, masterPassword := newTestManager(t)
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "hellopass", Password: "secret"})

	doc := transfer.ExportDocument{
		Version: transfer.ExportVersion,
		Entries: []vault.Entry{{
			Alias:    "github",
			Username: "updated",
			Password: "new-secret",
		}},
	}

	result, err := mgr.Import(masterPassword, doc, transfer.ConflictSkip)
	if err != nil {
		t.Fatalf("Import(skip) error = %v", err)
	}
	if result.Skipped != 1 {
		t.Fatalf("expected one skipped entry, got %+v", result)
	}

	if _, err := mgr.Import(masterPassword, doc, transfer.ConflictFail); err == nil {
		t.Fatal("expected conflict fail error")
	}
}

func TestLoadCurrentAndLoadOrCreateCurrent(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv(home.EnvKeepassHomePath, env.RootDir)

	if _, err := LoadCurrent(); err == nil {
		t.Fatal("expected LoadCurrent to fail before config exists")
	}

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := LoadCurrent(); err != nil {
		t.Fatalf("LoadCurrent() error = %v", err)
	}

	freshEnv := testutil.NewEnvironment(t)
	t.Setenv(home.EnvKeepassHomePath, freshEnv.RootDir)
	if _, err := LoadOrCreateCurrent(); err != nil {
		t.Fatalf("LoadOrCreateCurrent() error = %v", err)
	}
}

func stringPtr(value string) *string {
	return &value
}

type testV1Header struct {
	Argon2id struct {
		Time      uint32 `json:"time"`
		MemoryKiB uint32 `json:"memory_kib"`
		Threads   uint8  `json:"threads"`
		KeyLength uint32 `json:"key_length"`
	} `json:"argon2id"`
}

func readV1Header(path string) (testV1Header, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return testV1Header{}, err
	}

	headerLength := int(binary.BigEndian.Uint32(data[6:10]))
	headerBytes := data[10 : 10+headerLength]

	var header testV1Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return testV1Header{}, err
	}

	return header, nil
}
