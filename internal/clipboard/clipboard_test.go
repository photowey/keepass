package clipboard

import (
	"context"
	"testing"
	"time"
)

func TestClearAfterClearsClipboardAfterTimeout(t *testing.T) {
	t.Helper()

	var writes []string
	originalWriteAll := writeAll
	originalWaitForTimeout := waitForTimeout
	t.Cleanup(func() {
		writeAll = originalWriteAll
		waitForTimeout = originalWaitForTimeout
	})

	writeAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	waitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		if timeout != 5*time.Second {
			t.Fatalf("expected timeout 5s, got %s", timeout)
		}

		return false
	}

	interrupted, err := ClearAfter(context.Background(), 5*time.Second)
	if err != nil {
		t.Fatalf("ClearAfter() error = %v", err)
	}

	if interrupted {
		t.Fatal("expected timeout path, got interrupted")
	}

	if len(writes) != 1 || writes[0] != "" {
		t.Fatalf("expected clipboard clear write, got %#v", writes)
	}
}

func TestClearAfterClearsClipboardWhenInterrupted(t *testing.T) {
	t.Helper()

	var writes []string
	originalWriteAll := writeAll
	originalWaitForTimeout := waitForTimeout
	t.Cleanup(func() {
		writeAll = originalWriteAll
		waitForTimeout = originalWaitForTimeout
	})

	writeAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	waitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		return true
	}

	interrupted, err := ClearAfter(context.Background(), time.Second)
	if err != nil {
		t.Fatalf("ClearAfter() error = %v", err)
	}

	if !interrupted {
		t.Fatal("expected interrupted path")
	}

	if len(writes) != 1 || writes[0] != "" {
		t.Fatalf("expected clipboard clear write, got %#v", writes)
	}
}
