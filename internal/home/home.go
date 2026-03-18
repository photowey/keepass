/*
 * Copyright © 2023 the original author or authors.
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

package home

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	RootDirName        = ".keepass"
	ConfigFileName     = "keepass.config.json"
	DefaultVaultName   = "keepass.kp"
	EnvKeepassHomePath = "KEEPASS_HOME"
)

type Environment struct {
	RootDir         string
	ConfigFile      string
	DefaultVault    string
	ResolvedHomeDir string
}

func Detect() (Environment, error) {
	if custom := os.Getenv(EnvKeepassHomePath); custom != "" {
		root := filepath.Clean(custom)
		return Environment{
			RootDir:         root,
			ConfigFile:      filepath.Join(root, ConfigFileName),
			DefaultVault:    filepath.Join(root, DefaultVaultName),
			ResolvedHomeDir: root,
		}, nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return Environment{}, fmt.Errorf("resolve user home: %w", err)
	}

	root := filepath.Join(userHome, RootDirName)

	return Environment{
		RootDir:         root,
		ConfigFile:      filepath.Join(root, ConfigFileName),
		DefaultVault:    filepath.Join(root, DefaultVaultName),
		ResolvedHomeDir: userHome,
	}, nil
}
