package vault

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"testing"

	"github.com/photowey/keepass/internal/testutil"
)

func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()

	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	store := NewStore(cfg.ResolveVaultPath(env), cfg)

	const masterPassword = "correct horse battery staple"
	if err := store.Initialize(masterPassword, false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	return store, masterPassword
}

func TestStoreRoundTripAndPermissions(t *testing.T) {
	store, masterPassword := newTestStore(t)

	document, err := store.Load(masterPassword)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	document.Entries = append(document.Entries, Entry{
		Alias:    "github",
		Username: "hellopass",
		Password: "github-secret-2024",
		Tags:     []string{"code"},
	})

	if err := store.Save(masterPassword, document); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	reloaded, err := store.Load(masterPassword)
	if err != nil {
		t.Fatalf("Load() after Save error = %v", err)
	}

	if len(reloaded.Entries) != 1 || reloaded.Entries[0].Alias != "github" {
		t.Fatalf("unexpected entries after reload: %+v", reloaded.Entries)
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(store.Path())
		if err != nil {
			t.Fatalf("Stat() error = %v", err)
		}

		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Fatalf("expected vault mode 0600, got %#o", perm)
		}
	}
}

func TestStoreRejectsWrongPassword(t *testing.T) {
	store, _ := newTestStore(t)

	if _, err := store.Load("wrong password"); !errors.Is(err, ErrDecryptFailed) {
		t.Fatalf("expected ErrDecryptFailed, got %v", err)
	}
}

func TestStoreLoadRejectsMissingVault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	store := NewStore(cfg.ResolveVaultPath(env), cfg)

	if _, err := store.Load("master"); !errors.Is(err, ErrVaultNotInitialized) {
		t.Fatalf("expected ErrVaultNotInitialized, got %v", err)
	}
}

func TestStoreInitializeRejectsBlankMasterAndExistingVault(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	store := NewStore(cfg.ResolveVaultPath(env), cfg)

	if err := store.Initialize("", false); err == nil {
		t.Fatal("expected blank master password error")
	}

	if err := store.Initialize("master", false); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	if err := store.Initialize("master", false); err == nil {
		t.Fatal("expected existing vault error")
	}
}

func TestStoreSaveRejectsNilDocument(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)
	store := NewStore(cfg.ResolveVaultPath(env), cfg)

	if err := store.Save("master", nil); err == nil {
		t.Fatal("expected nil document error")
	}
}

func TestDecodeRejectsTamperedCiphertext(t *testing.T) {
	store, masterPassword := newTestStore(t)

	data, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	data[len(data)-1] ^= 0xFF
	if err := os.WriteFile(store.Path(), data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.Load(masterPassword); !errors.Is(err, ErrDecryptFailed) {
		t.Fatalf("expected ErrDecryptFailed, got %v", err)
	}
}

func TestDecodeRejectsUnsupportedVersion(t *testing.T) {
	store, masterPassword := newTestStore(t)

	data, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	binary.BigEndian.PutUint16(data[4:6], 99)
	if _, err := Decode(data, masterPassword); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatalf("expected ErrUnsupportedVersion, got %v", err)
	}
}

func TestDecodeRejectsInvalidHeaderParameters(t *testing.T) {
	store, masterPassword := newTestStore(t)

	data, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	headerLength := int(binary.BigEndian.Uint32(data[6:10]))
	headerBytes := data[10 : 10+headerLength]

	var header v1Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	header.Argon2id.Time = 0
	mutatedHeader, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	mutated := make([]byte, 0, len(data)-len(headerBytes)+len(mutatedHeader))
	mutated = append(mutated, data[:10]...)
	binary.BigEndian.PutUint32(mutated[6:10], uint32(len(mutatedHeader)))
	mutated = append(mutated, mutatedHeader...)
	mutated = append(mutated, data[10+headerLength:]...)

	if _, err := Decode(mutated, masterPassword); err == nil {
		t.Fatal("expected invalid file error after header mutation")
	}
}

func TestEncodeRejectsUnsupportedVersionAndBlankMaster(t *testing.T) {
	env := testutil.NewEnvironment(t)
	cfg := testutil.TestConfig(env)

	if _, err := Encode(99, []byte("hello"), "master", cfg); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatalf("expected ErrUnsupportedVersion, got %v", err)
	}

	if _, err := Encode(1, []byte("hello"), "", cfg); err == nil {
		t.Fatal("expected blank master password error")
	}
}

func TestDecodeRejectsInvalidFileHeaders(t *testing.T) {
	if _, err := Decode([]byte("bad"), "master"); !errors.Is(err, ErrInvalidFile) {
		t.Fatalf("expected ErrInvalidFile for short input, got %v", err)
	}

	data := make([]byte, 10)
	copy(data[:4], []byte("BAD!"))
	if _, err := Decode(data, "master"); !errors.Is(err, ErrInvalidFile) {
		t.Fatalf("expected ErrInvalidFile for bad magic, got %v", err)
	}
}
