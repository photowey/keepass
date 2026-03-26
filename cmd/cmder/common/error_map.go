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

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/vault"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, configs.ErrConfigNotFound),
		errors.Is(err, vault.ErrVaultNotInitialized):
		return WithExitCode(ExitCodeNotInitialized, err)
	case errors.Is(err, vault.ErrDecryptFailed):
		return WithExitCode(ExitCodeUnlockFailed, err)
	default:
		return err
	}
}
