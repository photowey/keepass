package rotate

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
	_ "unsafe"
)

//go:linkname rotateClipboardWriteAll github.com/photowey/keepass/internal/clipboard.writeAll
var rotateClipboardWriteAll func(string) error

//go:linkname rotateClipboardWaitForTimeout github.com/photowey/keepass/internal/clipboard.waitForTimeout
var rotateClipboardWaitForTimeout func(context.Context, time.Duration) bool

func TestNewRejectsGenerateAndPasswordTogether(t *testing.T) {
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

func TestNewRotatesPasswordAndRevealsIt(t *testing.T) {
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
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--password", "new-secret", "--reveal"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "New password: new-secret") {
		t.Fatalf("unexpected rotate output %q", out.String())
	}
}

func TestNewGeneratesPasswordByDefault(t *testing.T) {
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
	if !strings.Contains(output, "Rotated password for github") || !strings.Contains(output, "New password: ") {
		t.Fatalf("unexpected rotate generate output %q", output)
	}
}

func TestNewManualRotateDoesNotRevealWithoutFlag(t *testing.T) {
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
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--password", "new-secret"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Rotated password for github") || strings.Contains(output, "New password: new-secret") {
		t.Fatalf("unexpected non-reveal rotate output %q", output)
	}
}

func TestNewCopiesRotatedPasswordWithoutPrintingPlaintext(t *testing.T) {
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
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	var writes []string
	originalWriteAll := rotateClipboardWriteAll
	originalWaitForTimeout := rotateClipboardWaitForTimeout
	t.Cleanup(func() {
		rotateClipboardWriteAll = originalWriteAll
		rotateClipboardWaitForTimeout = originalWaitForTimeout
	})

	rotateClipboardWriteAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	rotateClipboardWaitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		return false
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--password", "new-secret", "--copy", "--copy-timeout", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(writes) != 1 || writes[0] != "new-secret" {
		t.Fatalf("unexpected clipboard writes %#v", writes)
	}

	output := out.String()
	if !strings.Contains(output, "Copied password to clipboard.") || strings.Contains(output, "New password: new-secret") {
		t.Fatalf("unexpected rotate copy output %q", output)
	}
}

func TestNewCopiesAndClearsRotatedPasswordAfterTimeout(t *testing.T) {
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
	}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	var writes []string
	originalWriteAll := rotateClipboardWriteAll
	originalWaitForTimeout := rotateClipboardWaitForTimeout
	t.Cleanup(func() {
		rotateClipboardWriteAll = originalWriteAll
		rotateClipboardWaitForTimeout = originalWaitForTimeout
	})

	rotateClipboardWriteAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	rotateClipboardWaitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		if timeout != time.Second {
			t.Fatalf("expected one second timeout, got %s", timeout)
		}
		return false
	}

	cmd := New()
	cmd.SetArgs([]string{"github", "--password", "new-secret", "--copy", "--copy-timeout", "1"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if len(writes) != 2 || writes[0] != "new-secret" || writes[1] != "" {
		t.Fatalf("unexpected clipboard writes %#v", writes)
	}

	output := out.String()
	if !strings.Contains(output, "Waiting 1s before clearing") || !strings.Contains(output, "Clipboard cleared.") {
		t.Fatalf("unexpected rotate timed copy output %q", output)
	}
}
