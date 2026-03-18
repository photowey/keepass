package root_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/root"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/version"
)

func TestCommandFlow(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "github-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--uri", "https://github.com", "--note", "personal", "--tag", "code")
	runCommand(t, "gitea-secret-2024\ngitea-secret-2024\nmaster\n", "add", "gitea", "hellopass", "--uri", "https://gitea.example.com", "--note", "work", "--tag", "code")

	listOutput := runCommand(t, "master\n", "list", "--tag", "code")
	if !strings.Contains(listOutput, "github") || strings.Contains(listOutput, "github-secret-2024") {
		t.Fatalf("unexpected list output: %s", listOutput)
	}

	getOutput := runCommand(t, "master\n", "get", "gith")
	if !strings.Contains(getOutput, "Password: [hidden]") {
		t.Fatalf("unexpected get output: %s", getOutput)
	}

	revealOutput := runCommand(t, "master\n", "get", "gith", "--reveal")
	if !strings.Contains(revealOutput, "github-secret-2024") {
		t.Fatalf("unexpected reveal output: %s", revealOutput)
	}

	runCommand(t, "master\n\n\nnew note\ncode\n\n", "update", "gith")

	updatedOutput := runCommand(t, "master\n", "get", "github")
	if !strings.Contains(updatedOutput, "Note: new note") {
		t.Fatalf("unexpected updated output: %s", updatedOutput)
	}

	runCommand(t, "master\ny\n", "delete", "gith")

	finalList := runCommand(t, "master\n", "list")
	if strings.Contains(finalList, "github") {
		t.Fatalf("expected github to be deleted, got: %s", finalList)
	}
}

func TestAddNonInteractiveRequiresArgs(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	cmd := root.NewCommand()
	cmd.SetArgs([]string{"add"})
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	_ = w.Close()
	cmd.SetIn(r)
	t.Cleanup(func() { _ = r.Close() })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}

func TestJSONOutputDoesNotRevealPasswordByDefault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "\n\ngithub-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--tag", "code")

	listJSON := runCommand(t, "master\n", "list", "--json")
	if !strings.HasPrefix(strings.TrimSpace(listJSON), "[") || strings.Contains(listJSON, "github-secret-2024") {
		t.Fatalf("unexpected list json output: %s", listJSON)
	}

	getJSON := runCommand(t, "master\n", "get", "github", "--json")
	if !strings.Contains(getJSON, "\"alias\"") || strings.Contains(getJSON, "github-secret-2024") {
		t.Fatalf("unexpected get json output: %s", getJSON)
	}

	revealJSON := runCommand(t, "master\n", "get", "github", "--json", "--reveal")
	if !strings.Contains(revealJSON, "github-secret-2024") {
		t.Fatalf("unexpected reveal get json output: %s", revealJSON)
	}
}

func TestCompletionCommandOutputsScript(t *testing.T) {
	cmd := root.NewCommand()
	cmd.SetArgs([]string{"completion", "bash"})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(completion) error = %v, stderr = %s", err, stderr.String())
	}

	if strings.TrimSpace(stdout.String()) == "" {
		t.Fatalf("expected non-empty completion output")
	}
}

func TestVersionFlagOutputsBuildInfo(t *testing.T) {
	cmd := root.NewCommand()
	cmd.SetArgs([]string{"--version"})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(--version) error = %v, stderr = %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if !strings.Contains(output, version.Now()) {
		t.Fatalf("expected version output to contain %q, got %q", version.Now(), output)
	}

	if !strings.Contains(output, "commit") || !strings.Contains(output, "built") {
		t.Fatalf("expected version output to include build metadata, got %q", output)
	}
}

func TestNonInteractiveFlagForcesFastFail(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	cmd := root.NewCommand()
	cmd.SetArgs([]string{"add", "--non-interactive"})

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	_ = w.Close()
	cmd.SetIn(r)
	t.Cleanup(func() { _ = r.Close() })

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}

func TestUpdateNonInteractiveRequiresMutationFlags(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "\n\ngithub-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--tag", "code")

	cmd := root.NewCommand()
	cmd.SetArgs([]string{"update", "github", "--non-interactive"})
	cmd.SetIn(strings.NewReader("master\n"))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error, got nil, stdout=%s stderr=%s", stdout.String(), stderr.String())
	}
}

func TestUpdateNonInteractiveAppliesFlagChanges(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "github-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--uri", "https://github.com", "--note", "personal", "--tag", "code")

	runCommand(t, "master\n", "update", "github", "--non-interactive", "--password", "github-rotated-2024", "--clear-uri", "--clear-note", "--clear-tags")

	output := runCommand(t, "master\n", "get", "github", "--reveal")
	if !strings.Contains(output, "github-rotated-2024") || strings.Contains(output, "https://github.com") || strings.Contains(output, "personal") {
		t.Fatalf("unexpected updated output: %s", output)
	}
}

func TestCopyDoesNotRevealPassword(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "\n\ngithub-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--tag", "code")

	cmd := root.NewCommand()
	cmd.SetArgs([]string{"get", "github", "--copy", "--copy-timeout", "0"})
	cmd.SetIn(strings.NewReader("master\n"))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()
	if err != nil {
		// Clipboard can be unavailable in some environments (e.g., headless CI).
		t.Skipf("clipboard copy unavailable: %v", err)
	}

	if strings.Contains(stdout.String(), "github-secret-2024") {
		t.Fatalf("expected password not to appear in stdout, got: %s", stdout.String())
	}
}

func TestRehashCommandUsesCurrentConfig(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "\n\n\ngithub-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass")

	cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save(updated cfg) error = %v", err)
	}

	output := runCommand(t, "master\n", "rehash")
	if !strings.Contains(output, "Rehashed vault") || !strings.Contains(output, "memory_kib=16384") {
		t.Fatalf("unexpected rehash output: %s", output)
	}

	getOutput := runCommand(t, "master\n", "get", "github", "--reveal")
	if !strings.Contains(getOutput, "github-secret-2024") {
		t.Fatalf("expected rehash to preserve vault contents, got: %s", getOutput)
	}
}

func TestDoctorJSONReportsRehashRecommendation(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")

	cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save(updated cfg) error = %v", err)
	}

	output := runCommand(t, "", "doctor", "--json")
	if !strings.Contains(output, "\"rehash_recommended\": true") {
		t.Fatalf("expected rehash recommendation in doctor json output, got: %s", output)
	}

	if !strings.Contains(output, "\"overall_status\": \"warn\"") {
		t.Fatalf("expected warn overall status in doctor json output, got: %s", output)
	}
}

func TestExportImportAuditAndRotateCommands(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	exportPath := filepath.Join(t.TempDir(), "entries.json")
	backupPath := filepath.Join(t.TempDir(), "backup")

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "\n\ngithub-secret-2024\ngithub-secret-2024\nmaster\n", "add", "github", "hellopass", "--uri", "https://github.com")
	runCommand(t, "master\n", "export", "--path", exportPath)

	runCommand(t, "master\n", "rotate", "github", "--password", "github-rotated-2024", "--reveal")
	runCommand(t, "master\n", "import", "--path", exportPath, "--conflict", "overwrite")

	auditOutput := runCommand(t, "master\n", "audit", "--json")
	if !strings.Contains(auditOutput, "\"overall_status\"") {
		t.Fatalf("unexpected audit output: %s", auditOutput)
	}

	runCommand(t, "", "backup", "--path", backupPath)
	if _, err := os.Stat(filepath.Join(backupPath, "manifest.json")); err != nil {
		t.Fatalf("expected backup manifest: %v", err)
	}
}

func runCommand(t *testing.T, input string, args ...string) string {
	t.Helper()

	cmd := root.NewCommand()
	cmd.SetArgs(args)
	cmd.SetIn(strings.NewReader(input))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v, stderr = %s", args, err, stderr.String())
	}

	return stdout.String()
}
