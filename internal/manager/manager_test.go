package manager

import (
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/internal/testutil"
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
		Username: "abc",
		Password: "secret123",
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

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "abc", Password: "one"})
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "gitea", Username: "abc", Password: "two"})

	_, err := mgr.Get(masterPassword, "gi")
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous alias error, got %v", err)
	}
}

func TestListFiltersByTags(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "mysql-prod", Username: "root", Password: "one", Tags: []string{"database", "prod"}})
	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "blog", Username: "writer", Password: "two", Tags: []string{"blog"}})

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

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "abc", Password: "one"})

	_, _, err := mgr.Add(masterPassword, AddInput{Alias: "github", Username: "other", Password: "two"})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate alias error, got %v", err)
	}
}

func TestUpdateGenerateAndDelete(t *testing.T) {
	mgr, masterPassword := newTestManager(t)

	_, _, _ = mgr.Add(masterPassword, AddInput{Alias: "github", Username: "abc", Password: "one"})

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
