package export

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
)

func TestNewRequiresPath(t *testing.T) {
	cmd := New()

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected usage error")
	}

	var cliErr common.CLIError
	if !errors.As(err, &cliErr) || cliErr.ExitCode != common.ExitCodeUsage {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestNewExportsEntriesToJSONFile(t *testing.T) {
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

	exportPath := filepath.Join(t.TempDir(), "entries.json")
	cmd := New()
	cmd.SetArgs([]string{"--path", exportPath})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if !strings.Contains(string(data), `"alias": "github"`) || !strings.Contains(out.String(), "Exported entries to ") {
		t.Fatalf("unexpected export result output=%q file=%q", out.String(), string(data))
	}
}
