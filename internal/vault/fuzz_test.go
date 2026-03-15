package vault

import (
	"path/filepath"
	"testing"

	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/testutil"
)

func FuzzDecode(f *testing.F) {
	root := f.TempDir()
	env := home.Environment{
		RootDir:         root,
		ConfigFile:      filepath.Join(root, home.ConfigFileName),
		DefaultVault:    filepath.Join(root, home.DefaultVaultName),
		ResolvedHomeDir: root,
	}
	cfg := testutil.TestConfig(env)
	payload := []byte(`{"version":1,"entries":[]}`)
	data, err := Encode(1, payload, "master", cfg)
	if err == nil {
		f.Add(data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data, "master")
	})
}
