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

package root_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExitCodeNotInitialized(t *testing.T) {
	t.Parallel()

	bin := buildTestBinary(t)

	cmd := exec.Command(bin, "list", "--json")
	cmd.Env = append(cmd.Environ(), "KEEPASS_HOME="+t.TempDir())
	_ = cmd.Run()

	if cmd.ProcessState == nil {
		t.Fatalf("missing process state")
	}

	if code := cmd.ProcessState.ExitCode(); code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}
}

func TestExitCodeUsageForNonInteractiveDeleteWithoutYes(t *testing.T) {
	t.Parallel()

	bin := buildTestBinary(t)
	homeDir := t.TempDir()

	initCmd := exec.Command(bin, "init")
	initCmd.Env = append(initCmd.Environ(), "KEEPASS_HOME="+homeDir)
	initCmd.Stdin = strings.NewReader("master\nmaster\n")
	if out, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("init error = %v, out=%s", err, string(out))
	}

	addCmd := exec.Command(bin, "add", "github", "hellopass")
	addCmd.Env = append(addCmd.Environ(), "KEEPASS_HOME="+homeDir)
	addCmd.Stdin = strings.NewReader("master\n")
	if out, err := addCmd.CombinedOutput(); err != nil {
		t.Fatalf("add error = %v, out=%s", err, string(out))
	}

	cmd := exec.Command(bin, "delete", "github", "--non-interactive")
	cmd.Env = append(cmd.Environ(), "KEEPASS_HOME="+homeDir)
	cmd.Stdin = strings.NewReader("master\n")
	_ = cmd.Run()

	if cmd.ProcessState == nil {
		t.Fatalf("missing process state")
	}

	if code := cmd.ProcessState.ExitCode(); code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
}

func TestExitCodeUsageForNonInteractiveUpdateWithoutFlags(t *testing.T) {
	t.Parallel()

	bin := buildTestBinary(t)
	homeDir := t.TempDir()

	initCmd := exec.Command(bin, "init")
	initCmd.Env = append(initCmd.Environ(), "KEEPASS_HOME="+homeDir)
	initCmd.Stdin = strings.NewReader("master\nmaster\n")
	if out, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("init error = %v, out=%s", err, string(out))
	}

	addCmd := exec.Command(bin, "add", "github", "hellopass")
	addCmd.Env = append(addCmd.Environ(), "KEEPASS_HOME="+homeDir)
	addCmd.Stdin = strings.NewReader("master\n")
	if out, err := addCmd.CombinedOutput(); err != nil {
		t.Fatalf("add error = %v, out=%s", err, string(out))
	}

	cmd := exec.Command(bin, "update", "github", "--non-interactive")
	cmd.Env = append(cmd.Environ(), "KEEPASS_HOME="+homeDir)
	cmd.Stdin = strings.NewReader("master\n")
	_ = cmd.Run()

	if cmd.ProcessState == nil {
		t.Fatalf("missing process state")
	}

	if code := cmd.ProcessState.ExitCode(); code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
}

func buildTestBinary(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))

	bin := filepath.Join(t.TempDir(), "keepass-test-bin")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build error = %v, out=%s", err, string(out))
	}

	return bin
}
