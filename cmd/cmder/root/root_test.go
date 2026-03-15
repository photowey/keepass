package root_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/root"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/testutil"
)

func TestCommandFlow(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	runCommand(t, "master\nmaster\n", "init")
	runCommand(t, "secret123\nsecret123\nmaster\n", "add", "github", "abc", "--uri", "https://github.com", "--note", "personal", "--tag", "code")
	runCommand(t, "secret234\nsecret234\nmaster\n", "add", "gitea", "abc", "--uri", "https://gitea.example.com", "--note", "work", "--tag", "code")

	listOutput := runCommand(t, "master\n", "list", "--tag", "code")
	if !strings.Contains(listOutput, "github") || strings.Contains(listOutput, "secret123") {
		t.Fatalf("unexpected list output: %s", listOutput)
	}

	getOutput := runCommand(t, "master\n", "get", "gith")
	if !strings.Contains(getOutput, "Password: [hidden]") {
		t.Fatalf("unexpected get output: %s", getOutput)
	}

	revealOutput := runCommand(t, "master\n", "get", "gith", "--reveal")
	if !strings.Contains(revealOutput, "secret123") {
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
