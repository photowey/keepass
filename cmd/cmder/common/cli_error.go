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

package common

import (
	"errors"
	"fmt"
)

const (
	ExitCodeGeneric        = 1
	ExitCodeUsage          = 2
	ExitCodeNotInitialized = 3
	ExitCodeUnlockFailed   = 4
)

type CLIError struct {
	ExitCode int
	Err      error
}

func (e CLIError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.ExitCode)
	}
	return e.Err.Error()
}

func (e CLIError) Unwrap() error { return e.Err }

func WithExitCode(code int, err error) error {
	if err == nil {
		return nil
	}
	return CLIError{ExitCode: code, Err: err}
}

func UsageError(message string) error {
	return WithExitCode(ExitCodeUsage, errors.New(message))
}
