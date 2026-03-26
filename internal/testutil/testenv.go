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

package testutil

import (
	"path/filepath"
	"testing"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func NewEnvironment(t *testing.T) home.Environment {
	t.Helper()

	root := t.TempDir()
	return home.Environment{
		RootDir:         root,
		ConfigFile:      filepath.Join(root, home.ConfigFileName),
		DefaultVault:    filepath.Join(root, home.DefaultVaultName),
		ResolvedHomeDir: root,
	}
}

func TestConfig(env home.Environment) configs.Config {
	cfg := configs.Default(env)
	cfg.Vault.Path = env.DefaultVault
	cfg.Security.Argon2id.Time = 1
	cfg.Security.Argon2id.MemoryKiB = 8 * 1024
	cfg.Security.Argon2id.Threads = 1
	cfg.Security.Argon2id.KeyLength = 32
	return cfg
}
