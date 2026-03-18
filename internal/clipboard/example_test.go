package clipboard

import (
	"context"
	"fmt"
	"time"
)

func ExampleClearAfter() {
	originalWriteAll := writeAll
	originalWaitForTimeout := waitForTimeout
	defer func() {
		writeAll = originalWriteAll
		waitForTimeout = originalWaitForTimeout
	}()

	writes := make([]string, 0, 1)
	writeAll = func(text string) error {
		writes = append(writes, text)
		return nil
	}
	waitForTimeout = func(ctx context.Context, timeout time.Duration) bool {
		return true
	}

	interrupted, err := ClearAfter(context.Background(), time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(interrupted)
	fmt.Println(len(writes))
	fmt.Println(writes[0] == "")

	// Output:
	// true
	// 1
	// true
}
