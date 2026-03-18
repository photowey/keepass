package vault

import (
	"os"
	"testing"
)

func TestInspectFileReadsV1Metadata(t *testing.T) {
	store, _ := newTestStore(t)

	meta, err := InspectFile(store.Path())
	if err != nil {
		t.Fatalf("InspectFile() error = %v", err)
	}

	if meta.FormatVersion != 1 {
		t.Fatalf("expected format version 1, got %d", meta.FormatVersion)
	}

	if meta.KDF != "argon2id" || meta.Cipher != "xchacha20poly1305" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}

	if meta.Argon2id == nil || meta.Argon2id.MemoryKiB == 0 {
		t.Fatalf("expected argon2id metadata, got %+v", meta)
	}
}

func TestInspectFileRejectsMissingVault(t *testing.T) {
	if _, err := InspectFile("/tmp/does-not-exist-keepass.kp"); err == nil {
		t.Fatal("expected missing vault error")
	}
}

func TestInspectFileRejectsInvalidHeader(t *testing.T) {
	store, _ := newTestStore(t)

	if err := os.WriteFile(store.Path(), []byte("bad"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := InspectFile(store.Path()); err == nil {
		t.Fatal("expected inspect error for invalid vault")
	}
}
