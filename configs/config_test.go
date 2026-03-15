package configs_test

import (
	"errors"
	"testing"

	"github.com/photowey/keepass/configs"
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
	cfg.PasswordGenerator.Alphabet = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for alphabet")
	}
}
