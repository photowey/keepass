package add

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewExposesExpectedFlags(t *testing.T) {
	cmd := New()

	for _, flagName := range []string{"uri", "note", "tag", "generate", "reveal-generated"} {
		if cmd.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected flag %q to be registered", flagName)
		}
	}
}

func TestNewRequiresAliasInNonInteractiveMode(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	cmd := New()
	cmd.Flags().Bool("non-interactive", false, "")
	if err := cmd.Flags().Set("non-interactive", "true"); err != nil {
		t.Fatalf("Set(non-interactive) error = %v", err)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	_ = w.Close()
	defer func() { _ = r.Close() }()

	cmd.SetIn(r)

	err = cmd.Execute()
	if err == nil {
		t.Fatal("expected alias validation error")
	}

	var cliErr common.CLIError
	if errors.As(err, &cliErr) {
		t.Fatalf("expected plain validation error, got CLIError %v", err)
	}

	if !strings.Contains(err.Error(), "alias is required in non-interactive mode") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestNewAddsEntryWithGeneratedPassword(t *testing.T) {
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
	cmd.SetArgs([]string{"github", "octocat", "--generate", "--reveal-generated", "--uri", "https://github.com", "--note", "personal", "--tag", "code"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Added entry github") || !strings.Contains(output, "Generated password: ") {
		t.Fatalf("unexpected add output %q", output)
	}

	entry, err := mgr.Get("master", "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if entry.Username != "octocat" || entry.URI != "https://github.com" {
		t.Fatalf("unexpected stored entry %#v", entry)
	}
}

func TestNewAddsEntryThroughInteractivePrompts(t *testing.T) {
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
	cmd.SetIn(strings.NewReader("github\noctocat\nhttps://github.com\npersonal\ncode,git\naccount-secret\naccount-secret\nmaster\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	entry, err := mgr.Get("master", "github")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if entry.Password != "account-secret" || entry.Note != "personal" || len(entry.Tags) != 2 {
		t.Fatalf("unexpected stored entry %#v", entry)
	}
}
