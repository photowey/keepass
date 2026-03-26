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

package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/vault"
)

type Status string

const (
	StatusOK    Status = "ok"
	StatusWarn  Status = "warn"
	StatusError Status = "error"
	StatusSkip  Status = "skip"
)

type Check struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message"`
}

type PasswordGeneratorInfo struct {
	DefaultLength           int    `json:"default_length"`
	Preset                  string `json:"preset"`
	UsesCustomAlphabet      bool   `json:"uses_custom_alphabet"`
	EffectiveAlphabetLength int    `json:"effective_alphabet_length"`
}

type ConfigInfo struct {
	Present           bool                   `json:"present"`
	Path              string                 `json:"path"`
	Argon2id          configs.Argon2idConfig `json:"argon2id"`
	PasswordGenerator PasswordGeneratorInfo  `json:"password_generator"`
}

type VaultInfo struct {
	Present       bool                    `json:"present"`
	Path          string                  `json:"path"`
	FormatVersion uint16                  `json:"format_version,omitempty"`
	KDF           string                  `json:"kdf,omitempty"`
	Cipher        string                  `json:"cipher,omitempty"`
	Argon2id      *configs.Argon2idConfig `json:"argon2id,omitempty"`
}

type Report struct {
	OverallStatus     Status     `json:"overall_status"`
	RootDir           string     `json:"root_dir"`
	ConfigFile        string     `json:"config_file"`
	VaultFile         string     `json:"vault_file"`
	RehashRecommended bool       `json:"rehash_recommended"`
	Recommendations   []string   `json:"recommendations,omitempty"`
	Checks            []Check    `json:"checks"`
	Config            ConfigInfo `json:"config"`
	Vault             VaultInfo  `json:"vault"`
}

func Collect() (Report, error) {
	env, err := home.Detect()
	if err != nil {
		return Report{}, err
	}

	report := Report{
		OverallStatus: StatusOK,
		RootDir:       env.RootDir,
		ConfigFile:    env.ConfigFile,
	}

	cfg, cfgPresent, err := loadEffectiveConfig(env)
	if err != nil {
		return Report{}, err
	}

	vaultPath := cfg.ResolveVaultPath(env)
	report.VaultFile = vaultPath
	report.Config = buildConfigInfo(cfgPresent, env.ConfigFile, cfg)
	report.Vault = VaultInfo{
		Path: vaultPath,
	}

	report.addCheck(configPresenceCheck(cfgPresent, env.ConfigFile))
	report.addPathPermissionChecks(env.RootDir, "root_dir", 0o700)
	if cfgPresent {
		report.addPathPermissionChecks(env.ConfigFile, "config_file", 0o600)
	}

	meta, vaultPresent, err := inspectVault(vaultPath)
	if err != nil {
		return Report{}, err
	}

	report.Vault.Present = vaultPresent
	report.addCheck(vaultPresenceCheck(vaultPresent, vaultPath))
	if vaultPresent {
		report.Vault.FormatVersion = meta.FormatVersion
		report.Vault.KDF = meta.KDF
		report.Vault.Cipher = meta.Cipher
		report.Vault.Argon2id = meta.Argon2id
		report.addPathPermissionChecks(vaultPath, "vault_file", 0o600)
		report.compareKDF(cfg.Security.Argon2id, meta.Argon2id)
	}

	return report, nil
}

func (r *Report) compareKDF(expected configs.Argon2idConfig, actual *configs.Argon2idConfig) {
	if actual == nil {
		return
	}

	if expected == *actual {
		r.addCheck(Check{
			Name:    "vault_kdf_alignment",
			Status:  StatusOK,
			Message: "Vault KDF parameters match the current config",
		})
		return
	}

	r.RehashRecommended = true
	r.Recommendations = append(r.Recommendations, "Run `keepass rehash` to rewrite the vault with the current Argon2 settings.")
	r.addCheck(Check{
		Name:    "vault_kdf_alignment",
		Status:  StatusWarn,
		Message: "Vault KDF parameters differ from the current config",
	})
}

func (r *Report) addPathPermissionChecks(path, name string, expected fs.FileMode) {
	if runtime.GOOS == "windows" {
		r.addCheck(Check{
			Name:    name + "_permissions",
			Status:  StatusSkip,
			Message: "Permission checks are skipped on this platform",
		})
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	actual := info.Mode().Perm()
	if actual == expected {
		r.addCheck(Check{
			Name:    name + "_permissions",
			Status:  StatusOK,
			Message: fmt.Sprintf("Permissions are %#o", actual),
		})
		return
	}

	r.addCheck(Check{
		Name:    name + "_permissions",
		Status:  StatusWarn,
		Message: fmt.Sprintf("Expected %#o, got %#o", expected, actual),
	})
}

func (r *Report) addCheck(check Check) {
	r.Checks = append(r.Checks, check)
	switch check.Status {
	case StatusError:
		r.OverallStatus = StatusError
	case StatusWarn:
		if r.OverallStatus != StatusError {
			r.OverallStatus = StatusWarn
		}
	}
}

func buildConfigInfo(present bool, path string, cfg configs.Config) ConfigInfo {
	alphabet, _ := cfg.PasswordGenerator.EffectiveAlphabet()
	preset := cfg.PasswordGenerator.Preset
	if preset == "" {
		preset = "compatible"
	}

	return ConfigInfo{
		Present:  present,
		Path:     path,
		Argon2id: cfg.Security.Argon2id,
		PasswordGenerator: PasswordGeneratorInfo{
			DefaultLength:           cfg.PasswordGenerator.DefaultLength,
			Preset:                  preset,
			UsesCustomAlphabet:      cfg.PasswordGenerator.Alphabet != "",
			EffectiveAlphabetLength: len(alphabet),
		},
	}
}

func loadEffectiveConfig(env home.Environment) (configs.Config, bool, error) {
	cfg, err := configs.Load(env)
	if err == nil {
		return cfg, true, nil
	}

	if errors.Is(err, configs.ErrConfigNotFound) {
		return configs.Default(env), false, nil
	}

	return configs.Config{}, false, err
}

func inspectVault(path string) (vault.Metadata, bool, error) {
	meta, err := vault.InspectFile(path)
	if err == nil {
		return meta, true, nil
	}

	if errors.Is(err, vault.ErrVaultNotInitialized) {
		return vault.Metadata{}, false, nil
	}

	return vault.Metadata{}, false, err
}

func configPresenceCheck(present bool, path string) Check {
	if present {
		return Check{
			Name:    "config_present",
			Status:  StatusOK,
			Message: fmt.Sprintf("Config found at %s", path),
		}
	}

	return Check{
		Name:    "config_present",
		Status:  StatusWarn,
		Message: fmt.Sprintf("Config not found at %s", path),
	}
}

func vaultPresenceCheck(present bool, path string) Check {
	if present {
		return Check{
			Name:    "vault_present",
			Status:  StatusOK,
			Message: fmt.Sprintf("Vault found at %s", path),
		}
	}

	return Check{
		Name:    "vault_present",
		Status:  StatusWarn,
		Message: fmt.Sprintf("Vault not found at %s", path),
	}
}

func (r Report) MarshalJSON() ([]byte, error) {
	type alias Report
	return json.Marshal(alias(r))
}
