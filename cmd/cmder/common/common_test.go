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
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/photowey/keepass/configs"
	auditreport "github.com/photowey/keepass/internal/audit"
	credentialaudit "github.com/photowey/keepass/internal/credentialaudit"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/vault"
	"github.com/spf13/cobra"
)

func TestIsInteractiveReturnsTrueForNonFileReader(t *testing.T) {
	if !IsInteractive(strings.NewReader("input")) {
		t.Fatal("expected non-file reader to be treated as interactive")
	}
}

func TestIsNonInteractiveReadsInheritedFlag(t *testing.T) {
	root := &cobra.Command{Use: "keepass"}
	root.PersistentFlags().Bool("non-interactive", false, "")

	cmd := &cobra.Command{Use: "child"}
	root.AddCommand(cmd)

	if err := root.PersistentFlags().Set("non-interactive", "true"); err != nil {
		t.Fatalf("Set(non-interactive) error = %v", err)
	}

	if !IsNonInteractive(cmd) {
		t.Fatal("expected inherited non-interactive flag to be true")
	}
}

func TestWithExitCodeAndUsageError(t *testing.T) {
	if got := WithExitCode(ExitCodeUsage, nil); got != nil {
		t.Fatalf("expected nil passthrough, got %v", got)
	}

	err := UsageError("bad input")
	var cliErr CLIError
	if !errors.As(err, &cliErr) {
		t.Fatalf("expected CLIError, got %v", err)
	}

	if cliErr.ExitCode != ExitCodeUsage {
		t.Fatalf("expected usage exit code, got %d", cliErr.ExitCode)
	}

	if cliErr.Error() != "bad input" {
		t.Fatalf("unexpected cli error message %q", cliErr.Error())
	}
}

func TestCLIErrorMessageWithoutWrappedError(t *testing.T) {
	err := CLIError{ExitCode: ExitCodeUnlockFailed}

	if got := err.Error(); got != "exit code 4" {
		t.Fatalf("unexpected error string %q", got)
	}
}

func TestMapError(t *testing.T) {
	t.Run("config missing maps to not initialized", func(t *testing.T) {
		err := MapError(configs.ErrConfigNotFound)
		assertCLIErrorCode(t, err, ExitCodeNotInitialized)
	})

	t.Run("vault missing maps to not initialized", func(t *testing.T) {
		err := MapError(vault.ErrVaultNotInitialized)
		assertCLIErrorCode(t, err, ExitCodeNotInitialized)
	})

	t.Run("decrypt failed maps to unlock failed", func(t *testing.T) {
		err := MapError(vault.ErrDecryptFailed)
		assertCLIErrorCode(t, err, ExitCodeUnlockFailed)
	})

	t.Run("unknown error passes through", func(t *testing.T) {
		original := errors.New("boom")
		err := MapError(original)
		if !errors.Is(err, original) {
			t.Fatalf("expected original error passthrough, got %v", err)
		}
	})
}

func TestPrintEntriesAndEntryOutputs(t *testing.T) {
	entry := sampleEntry()

	var out bytes.Buffer
	PrintEntries(&out, []vault.Entry{entry})
	text := out.String()
	if !strings.Contains(text, "ALIAS") || !strings.Contains(text, "github") {
		t.Fatalf("unexpected entries output %q", text)
	}

	out.Reset()
	PrintEntry(&out, entry, false)
	if strings.Contains(out.String(), entry.Password) || !strings.Contains(out.String(), "[hidden]") {
		t.Fatalf("expected hidden password output, got %q", out.String())
	}

	out.Reset()
	PrintEntry(&out, entry, true)
	if !strings.Contains(out.String(), entry.Password) {
		t.Fatalf("expected revealed password output, got %q", out.String())
	}
}

func TestPrintEntriesJSONAndEntryJSON(t *testing.T) {
	entry := sampleEntry()

	var out bytes.Buffer
	if err := PrintEntriesJSON(&out, []vault.Entry{entry}); err != nil {
		t.Fatalf("PrintEntriesJSON() error = %v", err)
	}
	if strings.Contains(out.String(), entry.Password) || !strings.Contains(out.String(), `"alias": "github"`) {
		t.Fatalf("unexpected entries json output %q", out.String())
	}

	out.Reset()
	if err := PrintEntryJSON(&out, entry, false); err != nil {
		t.Fatalf("PrintEntryJSON() error = %v", err)
	}
	if strings.Contains(out.String(), entry.Password) {
		t.Fatalf("expected password to be hidden in json output, got %q", out.String())
	}

	out.Reset()
	if err := PrintEntryJSON(&out, entry, true); err != nil {
		t.Fatalf("PrintEntryJSON(reveal) error = %v", err)
	}
	if !strings.Contains(out.String(), entry.Password) {
		t.Fatalf("expected revealed password in json output, got %q", out.String())
	}
}

func TestPrintConfigOutputs(t *testing.T) {
	env := home.Environment{
		RootDir:         "/tmp/keepass-home",
		ConfigFile:      "/tmp/keepass-home/keepass.config.json",
		DefaultVault:    "/tmp/keepass-home/keepass.kp",
		ResolvedHomeDir: "/tmp",
	}
	cfg := configs.Default(env)

	var out bytes.Buffer
	if err := PrintConfig(&out, env, cfg, false); err != nil {
		t.Fatalf("PrintConfig() error = %v", err)
	}
	if !strings.Contains(out.String(), `"initialized": false`) || !strings.Contains(out.String(), `"vault_file"`) {
		t.Fatalf("unexpected config json output %q", out.String())
	}

	out.Reset()
	if _, err := PrintConfigText(&out, env, cfg, true); err != nil {
		t.Fatalf("PrintConfigText() error = %v", err)
	}
	if !strings.Contains(out.String(), "Initialized: true") || !strings.Contains(out.String(), "Tip: use `keepass config --json`") {
		t.Fatalf("unexpected config text output %q", out.String())
	}
}

func TestPrintAuditOutputs(t *testing.T) {
	report := auditreport.Report{
		OverallStatus:     auditreport.StatusWarn,
		RootDir:           "/tmp/keepass-home",
		ConfigFile:        "/tmp/keepass-home/keepass.config.json",
		VaultFile:         "/tmp/keepass-home/keepass.kp",
		RehashRecommended: true,
		Recommendations:   []string{"Run `keepass rehash`."},
		Checks: []auditreport.Check{
			{Name: "config_present", Status: auditreport.StatusOK, Message: "Config found"},
			{Name: "vault_kdf_alignment", Status: auditreport.StatusWarn, Message: "KDF mismatch"},
		},
		Config: auditreport.ConfigInfo{
			Present:  true,
			Path:     "/tmp/keepass-home/keepass.config.json",
			Argon2id: configs.Argon2idConfig{Time: 1, MemoryKiB: 8192, Threads: 1, KeyLength: 32},
			PasswordGenerator: auditreport.PasswordGeneratorInfo{
				DefaultLength:           21,
				Preset:                  "compatible",
				UsesCustomAlphabet:      false,
				EffectiveAlphabetLength: 64,
			},
		},
		Vault: auditreport.VaultInfo{
			Present:  true,
			Path:     "/tmp/keepass-home/keepass.kp",
			Argon2id: &configs.Argon2idConfig{Time: 1, MemoryKiB: 8192, Threads: 1, KeyLength: 32},
		},
	}

	var out bytes.Buffer
	if err := PrintAuditJSON(&out, report); err != nil {
		t.Fatalf("PrintAuditJSON() error = %v", err)
	}
	if !strings.Contains(out.String(), `"rehash_recommended": true`) {
		t.Fatalf("unexpected audit json output %q", out.String())
	}

	out.Reset()
	if _, err := PrintAuditText(&out, report); err != nil {
		t.Fatalf("PrintAuditText() error = %v", err)
	}
	if !strings.Contains(out.String(), "Recommendations:") || !strings.Contains(out.String(), "[warn] vault_kdf_alignment") {
		t.Fatalf("unexpected audit text output %q", out.String())
	}
}

func TestPrintCredentialAuditOutputs(t *testing.T) {
	report := credentialaudit.Report{
		OverallStatus: "warn",
		MaxAgeDays:    90,
		Findings: []credentialaudit.Finding{
			{
				Type:    credentialaudit.FindingDuplicatePassword,
				Aliases: []string{"github", "gitlab"},
				Message: "Entries share the same password",
			},
		},
	}

	var out bytes.Buffer
	if err := PrintCredentialAuditJSON(&out, report); err != nil {
		t.Fatalf("PrintCredentialAuditJSON() error = %v", err)
	}
	if !strings.Contains(out.String(), `"overall_status": "warn"`) {
		t.Fatalf("unexpected credential audit json output %q", out.String())
	}

	out.Reset()
	if _, err := PrintCredentialAuditText(&out, report); err != nil {
		t.Fatalf("PrintCredentialAuditText() error = %v", err)
	}
	if !strings.Contains(out.String(), "Findings:") || !strings.Contains(out.String(), "github, gitlab") {
		t.Fatalf("unexpected credential audit text output %q", out.String())
	}
}

func TestIsInteractiveOnPipeReturnsFalse(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	defer func() {
		_ = r.Close()
		_ = w.Close()
	}()

	if IsInteractive(r) {
		t.Fatal("expected pipe input to be non-interactive")
	}
}

func assertCLIErrorCode(t *testing.T, err error, code int) {
	t.Helper()

	var cliErr CLIError
	if !errors.As(err, &cliErr) {
		t.Fatalf("expected CLIError, got %v", err)
	}
	if cliErr.ExitCode != code {
		t.Fatalf("expected exit code %d, got %d", code, cliErr.ExitCode)
	}
}

func sampleEntry() vault.Entry {
	now := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	return vault.Entry{
		Alias:             "github",
		Username:          "octocat",
		Password:          "secret-123",
		URI:               "https://github.com",
		Note:              "personal",
		Tags:              []string{"code", "personal"},
		CreatedAt:         now,
		UpdatedAt:         now,
		PasswordUpdatedAt: now,
	}
}
