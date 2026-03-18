package configs_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/testutil"
)

func TestLoadOrCreateCreatesValidatedConfig(t *testing.T) {
	env := testutil.NewEnvironment(t)

	cfg, err := configs.LoadOrCreate(env)
	if err != nil {
		t.Fatalf("configs.LoadOrCreate() error = %v", err)
	}

	if cfg.Version != configs.CurrentVersion {
		t.Fatalf("expected version %d, got %d", configs.CurrentVersion, cfg.Version)
	}

	if cfg.ResolveVaultPath(env) == "" {
		t.Fatal("expected resolved vault path")
	}

	if cfg.Security.Argon2id.MemoryKiB != 256*1024 {
		t.Fatalf("expected default argon2id memory 262144 KiB, got %d", cfg.Security.Argon2id.MemoryKiB)
	}

	if _, err := configs.Load(env); err != nil {
		t.Fatalf("configs.Load() error = %v", err)
	}
}

func TestLoadMissingConfigReturnsSentinel(t *testing.T) {
	env := testutil.NewEnvironment(t)

	_, err := configs.Load(env)
	if !errors.Is(err, configs.ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestValidateRejectsBrokenGeneratorConfig(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := configs.Default(env)
	cfg.PasswordGenerator.DefaultLength = 0

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for default_length")
	}

	cfg = configs.Default(env)
	cfg.PasswordGenerator.Preset = "unknown"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unsupported preset")
	}
}

func TestPasswordGeneratorEffectiveAlphabet(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := configs.Default(env)

	alphabet, err := cfg.PasswordGenerator.EffectiveAlphabet()
	if err != nil {
		t.Fatalf("EffectiveAlphabet() error = %v", err)
	}

	if alphabet == "" {
		t.Fatal("expected non-empty preset alphabet")
	}

	cfg = configs.Default(env)
	cfg.PasswordGenerator.Preset = "symbols"
	cfg.PasswordGenerator.Alphabet = "!custom!"

	alphabet, err = cfg.PasswordGenerator.EffectiveAlphabet()
	if err != nil {
		t.Fatalf("EffectiveAlphabet() custom error = %v", err)
	}

	if alphabet != "!custom!" {
		t.Fatalf("expected custom alphabet override, got %q", alphabet)
	}
}

func TestDefaultUsesEnvironmentVaultPathOutsideStandardHome(t *testing.T) {
	env := home.Environment{
		RootDir:         filepath.Join(t.TempDir(), "custom-home"),
		ConfigFile:      filepath.Join(t.TempDir(), "custom-home", home.ConfigFileName),
		DefaultVault:    filepath.Join(t.TempDir(), "custom-home", home.DefaultVaultName),
		ResolvedHomeDir: t.TempDir(),
	}

	cfg := configs.Default(env)
	if cfg.Vault.Path != env.DefaultVault {
		t.Fatalf("expected default vault path %q, got %q", env.DefaultVault, cfg.Vault.Path)
	}
}

func TestResolveVaultPathKeepsAbsolutePath(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := configs.Default(env)
	cfg.Vault.Path = "/var/lib/keepass/data.kp"

	if got := cfg.ResolveVaultPath(env); got != "/var/lib/keepass/data.kp" {
		t.Fatalf("expected absolute path, got %q", got)
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	env := testutil.NewEnvironment(t)
	if err := os.WriteFile(env.ConfigFile, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := configs.Load(env); err == nil {
		t.Fatal("expected invalid json error")
	}
}

func TestLoadOrCreateReturnsExistingLoadError(t *testing.T) {
	env := testutil.NewEnvironment(t)
	if err := os.WriteFile(env.ConfigFile, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := configs.LoadOrCreate(env); err == nil {
		t.Fatal("expected load error to be returned")
	}
}
