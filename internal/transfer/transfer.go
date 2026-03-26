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

package transfer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/vault"
	"github.com/photowey/keepass/pkg/files"
)

const (
	ExportVersion   = 1
	ManifestVersion = 1

	ConflictFail      = "fail"
	ConflictSkip      = "skip"
	ConflictOverwrite = "overwrite"
)

type ExportDocument struct {
	Version    int           `json:"version"`
	ExportedAt time.Time     `json:"exported_at"`
	Entries    []vault.Entry `json:"entries"`
}

type ImportResult struct {
	Added     int `json:"added"`
	Overwrote int `json:"overwrote"`
	Skipped   int `json:"skipped"`
}

type BackupManifest struct {
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	ConfigFile string    `json:"config_file"`
	VaultFile  string    `json:"vault_file"`
}

func WriteExport(path string, doc ExportDocument) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("encode export: %w", err)
	}

	return files.WriteFileAtomic(path, append(data, '\n'), 0o600)
}

func ReadExport(path string) (ExportDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExportDocument{}, fmt.Errorf("read export: %w", err)
	}

	var doc ExportDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return ExportDocument{}, fmt.Errorf("decode export: %w", err)
	}

	if doc.Version != ExportVersion {
		return ExportDocument{}, fmt.Errorf("unsupported export version %d", doc.Version)
	}

	return doc, nil
}

func NormalizeConflictStrategy(strategy string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(strategy))
	if normalized == "" {
		normalized = ConflictFail
	}

	switch normalized {
	case ConflictFail, ConflictSkip, ConflictOverwrite:
		return normalized, nil
	default:
		return "", fmt.Errorf("unsupported conflict strategy %q", strategy)
	}
}

func CreateBackupBundle(path string, env home.Environment, cfg configs.Config, force bool, now time.Time) (string, error) {
	if path == "" {
		return "", errors.New("backup path cannot be blank")
	}

	targetDir := filepath.Clean(path)
	if info, err := os.Stat(targetDir); err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("backup path %s is not a directory", targetDir)
		}
		entries, err := os.ReadDir(targetDir)
		if err != nil {
			return "", fmt.Errorf("read backup dir: %w", err)
		}
		if len(entries) > 0 && !force {
			return "", fmt.Errorf("backup directory %s is not empty", targetDir)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat backup dir: %w", err)
	}

	if err := files.EnsureDir(targetDir, 0o700); err != nil {
		return "", fmt.Errorf("ensure backup dir: %w", err)
	}

	configBase := filepath.Base(env.ConfigFile)
	vaultBase := filepath.Base(cfg.ResolveVaultPath(env))
	manifest := BackupManifest{
		Version:    ManifestVersion,
		CreatedAt:  now.UTC(),
		ConfigFile: configBase,
		VaultFile:  vaultBase,
	}

	if err := copyFile(env.ConfigFile, filepath.Join(targetDir, configBase), 0o600); err != nil {
		return "", err
	}

	if err := copyFile(cfg.ResolveVaultPath(env), filepath.Join(targetDir, vaultBase), 0o600); err != nil {
		return "", err
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encode manifest: %w", err)
	}

	if err := files.WriteFileAtomic(filepath.Join(targetDir, "manifest.json"), append(manifestBytes, '\n'), 0o600); err != nil {
		return "", fmt.Errorf("write manifest: %w", err)
	}

	return targetDir, nil
}

func RestoreBackupBundle(path string, env home.Environment, force bool) error {
	targetDir := filepath.Clean(path)
	data, err := os.ReadFile(filepath.Join(targetDir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	var manifest BackupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}

	if manifest.Version != ManifestVersion {
		return fmt.Errorf("unsupported manifest version %d", manifest.Version)
	}

	configTarget := env.ConfigFile
	configSource := filepath.Join(targetDir, manifest.ConfigFile)
	configBytes, err := os.ReadFile(configSource)
	if err != nil {
		return fmt.Errorf("read backup config: %w", err)
	}

	var cfg configs.Config
	if err := json.Unmarshal(configBytes, &cfg); err != nil {
		return fmt.Errorf("decode backup config: %w", err)
	}

	vaultTarget, restoredConfig, err := restoredVaultTarget(env, manifest, cfg)
	if err != nil {
		return err
	}
	if !force {
		if files.Exists(configTarget) || files.Exists(vaultTarget) {
			return errors.New("restore target already contains config or vault, use --force to overwrite")
		}
	}

	if err := files.EnsureDir(env.RootDir, 0o700); err != nil {
		return fmt.Errorf("ensure keepass home: %w", err)
	}

	restoredConfigBytes, err := json.MarshalIndent(restoredConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("encode restored config: %w", err)
	}

	if err := files.WriteFileAtomic(configTarget, append(restoredConfigBytes, '\n'), 0o600); err != nil {
		return fmt.Errorf("write restored config: %w", err)
	}

	if err := copyFile(filepath.Join(targetDir, manifest.VaultFile), vaultTarget, 0o600); err != nil {
		return err
	}

	return nil
}

func restoredVaultTarget(env home.Environment, manifest BackupManifest, cfg configs.Config) (string, configs.Config, error) {
	restored := cfg

	switch {
	case strings.HasPrefix(restored.Vault.Path, "~/"):
		return restored.ResolveVaultPath(env), restored, nil
	case filepath.IsAbs(restored.Vault.Path):
		target := filepath.Join(env.RootDir, filepath.Base(manifest.VaultFile))
		restored.Vault.Path = target
		return target, restored, nil
	default:
		target := filepath.Join(env.RootDir, filepath.Clean(restored.Vault.Path))
		restored.Vault.Path = target
		return target, restored, nil
	}
}

func copyFile(src, dst string, perm fs.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}

	if err := files.WriteFileAtomic(dst, data, perm); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}

	return nil
}

func SortEntries(entries []vault.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Alias < entries[j].Alias
	})
}
