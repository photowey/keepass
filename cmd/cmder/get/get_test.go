package get

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/vault"
	_ "unsafe"
)

//go:linkname clipboardWriteAll github.com/photowey/keepass/internal/clipboard.writeAll
var clipboardWriteAll func(string) error

//go:linkname clipboardWaitForTimeout github.com/photowey/keepass/internal/clipboard.waitForTimeout
var clipboardWaitForTimeout func(context.Context, time.Duration) bool

func TestNewRegistersFlags(t *testing.T) {
	cmd := New()

	for _, flagName := range []string{"reveal", "json", "copy", "copy-timeout"} {
		if cmd.Flags().Lookup(flagName) == nil {
			t.Fatalf("expected flag %q to be registered", flagName)
		}
	}
}

func TestNewPrintsRevealedEntryAsJSON(t *testing.T) {
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

	cmd := New()
	cmd.SetArgs([]string{"github", "--json", "--reveal"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"alias": "github"`) || !strings.Contains(output, "secret-123") {
		t.Fatalf("unexpected get output %q", output)
	}
}

func TestNewPrintsHiddenEntryAsText(t *testing.T) {
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
	cmd.SetArgs([]string{"github"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Password: [hidden]") || strings.Contains(output, "secret-123") {
		t.Fatalf("unexpected get text output %q", output)
	}
}

func TestNewPrintsHiddenEntryAsJSON(t *testing.T) {
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
	cmd.SetArgs([]string{"github", "--json"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"password": ""`) || strings.Contains(output, "secret-123") {
		t.Fatalf("unexpected hidden get json output %q", output)
	}
}

func TestNewPrintsRevealedEntryAsText(t *testing.T) {
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
	cmd.SetArgs([]string{"github", "--reveal"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Password: secret-123") {
		t.Fatalf("unexpected revealed get text output %q", output)
	}
}

func TestNewCopiesPasswordWithoutRevealingIt(t *testing.T) {
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

	var writes []string
	originalWriteAll := clipboardWriteAll
	originalWaitForTimeout := clipboardWaitForTimeout
	t.Cleanup(func() {
		clipboardWriteAll = originalWriteAll
		clipboardWaitForTimeout = originalWaitForTimeout
	})

	clipboardWriteAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	clipboardWaitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		return false
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--copy", "--copy-timeout", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(writes) != 1 || writes[0] != "secret-123" {
		t.Fatalf("unexpected clipboard writes %#v", writes)
	}

	output := out.String()
	if !strings.Contains(output, "Copied password to clipboard.") || strings.Contains(output, "secret-123") {
		t.Fatalf("unexpected copy output %q", output)
	}
}

func TestNewCopiesAndClearsClipboardAfterTimeout(t *testing.T) {
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

	var writes []string
	originalWriteAll := clipboardWriteAll
	originalWaitForTimeout := clipboardWaitForTimeout
	t.Cleanup(func() {
		clipboardWriteAll = originalWriteAll
		clipboardWaitForTimeout = originalWaitForTimeout
	})

	clipboardWriteAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	clipboardWaitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		if timeout != time.Second {
			t.Fatalf("expected one second timeout, got %s", timeout)
		}
		return false
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--copy", "--copy-timeout", "1"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(writes) != 2 || writes[0] != "secret-123" || writes[1] != "" {
		t.Fatalf("unexpected clipboard writes %#v", writes)
	}

	output := out.String()
	if !strings.Contains(output, "Waiting 1s before clearing") || !strings.Contains(output, "Clipboard cleared.") {
		t.Fatalf("unexpected timed copy output %q", output)
	}
}

func TestNewReturnsErrorWhenClipboardCopyFails(t *testing.T) {
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

	originalWriteAll := clipboardWriteAll
	t.Cleanup(func() { clipboardWriteAll = originalWriteAll })
	clipboardWriteAll = func(text string) error {
		return errors.New("clipboard unavailable")
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--copy", "--copy-timeout", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "copy to clipboard") {
		t.Fatalf("expected clipboard copy error, got %v", err)
	}
}

func TestNewRejectsCopyWhenPasswordIsEmpty(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Save("master", &vault.Document{
		Version: 1,
		Entries: []vault.Entry{
			{
				Alias:    "github",
				Username: "octocat",
				Password: "",
			},
		},
	}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--copy", "--copy-timeout", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "empty password cannot be copied") {
		t.Fatalf("expected empty password error, got %v", err)
	}
}
