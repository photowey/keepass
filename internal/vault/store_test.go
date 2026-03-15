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
		Username: "abc",
		Password: "s3cret",
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
