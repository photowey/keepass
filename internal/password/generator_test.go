package password

import (
	"strings"
	"testing"
)

func TestGenerateUsesRequestedAlphabetAndLength(t *testing.T) {
	value, err := Generate(32, "abc123")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(value) != 32 {
		t.Fatalf("expected length 32, got %d", len(value))
	}

	for _, r := range value {
		if !strings.ContainsRune("abc123", r) {
			t.Fatalf("generated rune %q not found in alphabet", r)
		}
	}
}

func TestGenerateRejectsInvalidInput(t *testing.T) {
	if _, err := Generate(0, "abc"); err == nil {
		t.Fatal("expected error for zero length")
	}

	if _, err := Generate(8, ""); err == nil {
		t.Fatal("expected error for empty alphabet")
	}
}
