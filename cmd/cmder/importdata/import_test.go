package importdata

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/transfer"
	"github.com/photowey/keepass/internal/vault"
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

func TestNewImportsEntriesFromJSONFile(t *testing.T) {
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

	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	doc := transfer.ExportDocument{
		Version:    transfer.ExportVersion,
		ExportedAt: now,
		Entries: []vault.Entry{
			{
				Alias:             "github",
				Username:          "octocat",
				Password:          "secret-123",
				URI:               "https://github.com",
				CreatedAt:         now,
				UpdatedAt:         now,
				PasswordUpdatedAt: now,
			},
		},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	importPath := filepath.Join(t.TempDir(), "entries.json")
	if err := os.WriteFile(importPath, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cmd := New()
	cmd.SetArgs([]string{"--path", importPath, "--conflict", "overwrite"})
	cmd.SetIn(strings.NewReader("master\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Added=1") {
		t.Fatalf("unexpected import output %q", out.String())
	}

	if _, err := mgr.Get("master", "github"); err != nil {
		t.Fatalf("expected imported entry, got error %v", err)
	}
}
