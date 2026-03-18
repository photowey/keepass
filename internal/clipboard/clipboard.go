package clipboard

import (
	"context"
	"time"

	sysclipboard "github.com/atotto/clipboard"
)

var writeAll = sysclipboard.WriteAll

var waitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return false
	case <-ctx.Done():
		return true
	}
}

func Copy(text string) error {
	return writeAll(text)
}

func Clear() error {
	return writeAll("")
}

func ClearAfter(ctx context.Context, timeout time.Duration) (bool, error) {
	if timeout <= 0 {
		return false, nil
	}

	interrupted := waitForTimeout(ctx, timeout)
	return interrupted, Clear()
}
