package restore

import (
	"bytes"
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

func TestNewRestoresBackupBundle(t *testing.T) {
	sourceEnv := testutil.NewEnvironment(t)
	sourceCfg := testutil.TestConfig(sourceEnv)
	if err := configs.Save(sourceEnv, sourceCfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	mgr := manager.New(sourceEnv, sourceCfg)
	if err := mgr.Initialize("master", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	backupDir := filepath.Join(t.TempDir(), "bundle")
	if _, err := transfer.CreateBackupBundle(backupDir, sourceEnv, sourceCfg, false, time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("CreateBackupBundle() error = %v", err)
	}

	restoreHome := t.TempDir()
	t.Setenv("KEEPASS_HOME", restoreHome)

	cmd := New()
	cmd.SetArgs([]string{"--path", backupDir})

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Restored backup bundle from ") {
		t.Fatalf("unexpected restore output %q", out.String())
	}

	if _, err := os.Stat(filepath.Join(restoreHome, "keepass.config.json")); err != nil {
		t.Fatalf("expected restored config file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(restoreHome, "keepass.kp")); err != nil {
		t.Fatalf("expected restored vault file: %v", err)
	}
}
