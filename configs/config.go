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

package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/pkg/filez"
)

const CurrentVersion = 1

var ErrConfigNotFound = errors.New("keepass config not found, run `keepass init` first")

type Config struct {
	Version           int               `json:"version"`
	Vault             VaultConfig       `json:"vault"`
	Security          SecurityConfig    `json:"security"`
	PasswordGenerator PasswordGenerator `json:"password_generator"`
}

type VaultConfig struct {
	Path          string `json:"path"`
	FormatVersion uint16 `json:"format_version"`
}

type SecurityConfig struct {
	Argon2id Argon2idConfig `json:"argon2id"`
}

type Argon2idConfig struct {
	Time      uint32 `json:"time"`
	MemoryKiB uint32 `json:"memory_kib"`
	Threads   uint8  `json:"threads"`
	KeyLength uint32 `json:"key_length"`
}

type PasswordGenerator struct {
	DefaultLength int    `json:"default_length"`
	Alphabet      string `json:"alphabet"`
}

func Default(env home.Environment) Config {
	vaultPath := "~/" + filepath.ToSlash(filepath.Join(home.RootDirName, home.DefaultVaultName))
	if env.RootDir != filepath.Join(env.ResolvedHomeDir, home.RootDirName) {
		vaultPath = env.DefaultVault
	}

	return Config{
		Version: CurrentVersion,
		Vault: VaultConfig{
			Path:          vaultPath,
			FormatVersion: 1,
		},
		Security: SecurityConfig{
			Argon2id: Argon2idConfig{
				Time:      3,
				MemoryKiB: 64 * 1024,
				Threads:   4,
				KeyLength: 32,
			},
		},
		PasswordGenerator: PasswordGenerator{
			DefaultLength: 21,
			Alphabet:      "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_",
		},
	}
}

func Load(env home.Environment) (Config, error) {
	data, err := os.ReadFile(env.ConfigFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, ErrConfigNotFound
		}

		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func LoadOrCreate(env home.Environment) (Config, error) {
	cfg, err := Load(env)
	if err == nil {
		return cfg, nil
	}

	if !errors.Is(err, ErrConfigNotFound) {
		return Config{}, err
	}

	cfg = Default(env)
	if err := Save(env, cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Save(env home.Environment, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := filez.EnsureDir(env.RootDir, 0o700); err != nil {
		return fmt.Errorf("ensure keepass home: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	if err := filez.WriteFileAtomic(env.ConfigFile, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

func (c Config) Validate() error {
	if c.Version != CurrentVersion {
		return fmt.Errorf("unsupported config version %d", c.Version)
	}

	if c.Vault.FormatVersion == 0 {
		return errors.New("vault format version must be positive")
	}

	if strings.TrimSpace(c.Vault.Path) == "" {
		return errors.New("vault path cannot be blank")
	}

	if c.Security.Argon2id.Time == 0 {
		return errors.New("argon2id time must be positive")
	}

	if c.Security.Argon2id.MemoryKiB < 8*1024 {
		return errors.New("argon2id memory_kib must be at least 8192")
	}

	if c.Security.Argon2id.Threads == 0 {
		return errors.New("argon2id threads must be positive")
	}

	if c.Security.Argon2id.KeyLength < 32 {
		return errors.New("argon2id key_length must be at least 32")
	}

	if c.PasswordGenerator.DefaultLength <= 0 {
		return errors.New("password_generator.default_length must be positive")
	}

	if strings.TrimSpace(c.PasswordGenerator.Alphabet) == "" {
		return errors.New("password_generator.alphabet cannot be blank")
	}

	return nil
}

func (c Config) ResolveVaultPath(env home.Environment) string {
	if strings.HasPrefix(c.Vault.Path, "~/") {
		return filepath.Join(env.ResolvedHomeDir, filepath.FromSlash(strings.TrimPrefix(c.Vault.Path, "~/")))
	}

	return filepath.Clean(c.Vault.Path)
}
