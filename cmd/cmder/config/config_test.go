package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewPrintsDefaultConfigWhenUninitialized(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	cmd := New()
	cmd.SetArgs([]string{"--json"})

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, `"initialized": false`) {
		t.Fatalf("expected uninitialized config output, got %s", output)
	}

	if !strings.Contains(output, `"root_dir": `) {
		t.Fatalf("expected resolved root_dir in output, got %s", output)
	}
}
