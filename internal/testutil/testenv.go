package testutil

import (
	"path/filepath"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func NewEnvironment(t *testing.T) home.Environment {
	t.Helper()

	root := t.TempDir()
	return home.Environment{
		RootDir:         root,
		ConfigFile:      filepath.Join(root, home.ConfigFileName),
		DefaultVault:    filepath.Join(root, home.DefaultVaultName),
		ResolvedHomeDir: root,
	}
}

func TestConfig(env home.Environment) configs.Config {
	cfg := configs.Default(env)
	cfg.Vault.Path = env.DefaultVault
	cfg.Security.Argon2id.Time = 1
	cfg.Security.Argon2id.MemoryKiB = 8 * 1024
	cfg.Security.Argon2id.Threads = 1
	cfg.Security.Argon2id.KeyLength = 32
	return cfg
}
