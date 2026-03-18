package audit

import (
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/testutil"
	"github.com/photowey/keepass/internal/vault"
)

func TestCollectReportsHealthyVault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	report, err := Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if report.Config.Present != true || report.Vault.Present != true {
		t.Fatalf("expected config and vault present, got %+v", report)
	}

	if report.RehashRecommended {
		t.Fatalf("expected no rehash recommendation, got %+v", report)
	}
}

func TestCollectRecommendsRehashWhenConfigDiffers(t *testing.T) {
	env := testutil.NewEnvironment(t)
	t.Setenv("KEEPASS_HOME", env.RootDir)

	cfg := testutil.TestConfig(env)
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	store := vault.NewStore(cfg.ResolveVaultPath(env), cfg)
	if err := store.Initialize("master-password", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	cfg.Security.Argon2id.MemoryKiB = 16 * 1024
	if err := configs.Save(env, cfg); err != nil {
		t.Fatalf("Save(updated cfg) error = %v", err)
	}

	report, err := Collect()
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	if !report.RehashRecommended {
		t.Fatalf("expected rehash recommendation, got %+v", report)
	}
}
