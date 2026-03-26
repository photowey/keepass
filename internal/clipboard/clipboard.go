/*
 * Copyright © 2023-present the keepass authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
