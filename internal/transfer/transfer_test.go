package transfer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/vault"
)

func TestWriteReadExportRoundTrip(t *testing.T) {
	doc := ExportDocument{
		Version:    ExportVersion,
		ExportedAt: time.Unix(1_700_000_000, 0).UTC(),
		Entries: []vault.Entry{{
			Alias:    "github",
			Username: "hellopass",
			Password: "github-secret-2024",
		}},
	}

	path := filepath.Join(t.TempDir(), "export.json")
	if err := WriteExport(path, doc); err != nil {
		t.Fatalf("WriteExport() error = %v", err)
	}

	loaded, err := ReadExport(path)
	if err != nil {
		t.Fatalf("ReadExport() error = %v", err)
	}

	if len(loaded.Entries) != 1 || loaded.Entries[0].Alias != "github" {
		t.Fatalf("unexpected export content: %+v", loaded)
	}
}

func TestBackupAndRestoreBundle(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	backupDir := filepath.Join(t.TempDir(), "bundle")
	if _, err := CreateBackupBundle(backupDir, env, cfg, false, time.Unix(1_700_000_000, 0)); err != nil {
		t.Fatalf("CreateBackupBundle() error = %v", err)
	}

	restoreRoot := filepath.Join(t.TempDir(), "restore-home")
	restoreEnv := home.Environment{
		RootDir:         restoreRoot,
		ConfigFile:      filepath.Join(restoreRoot, home.ConfigFileName),
		DefaultVault:    filepath.Join(restoreRoot, home.DefaultVaultName),
		ResolvedHomeDir: restoreRoot,
	}

	if err := RestoreBackupBundle(backupDir, restoreEnv, false); err != nil {
		t.Fatalf("RestoreBackupBundle() error = %v", err)
	}

	if _, err := os.Stat(restoreEnv.ConfigFile); err != nil {
		t.Fatalf("expected restored config file: %v", err)
	}

	if _, err := os.Stat(restoreEnv.DefaultVault); err != nil {
		t.Fatalf("expected restored vault file: %v", err)
	}
}

func TestNormalizeConflictStrategy(t *testing.T) {
	if got, err := NormalizeConflictStrategy(""); err != nil || got != ConflictFail {
		t.Fatalf("expected default fail strategy, got %q err=%v", got, err)
	}

	if _, err := NormalizeConflictStrategy("invalid"); err == nil {
		t.Fatal("expected invalid strategy error")
	}
}

func TestReadExportRejectsUnsupportedVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "export.json")
	data, err := json.Marshal(ExportDocument{Version: 99})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := ReadExport(path); err == nil {
		t.Fatal("expected unsupported version error")
	}
}

func TestCreateBackupBundleRejectsNonEmptyDirWithoutForce(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	backupDir := filepath.Join(t.TempDir(), "bundle")
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(backupDir, "existing.txt"), []byte("x"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := CreateBackupBundle(backupDir, env, cfg, false, time.Now()); err == nil {
		t.Fatal("expected non-empty backup dir error")
	}
}
