package init

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRegistersForceFlag(t *testing.T) {
	cmd := New()

	if cmd.Flags().Lookup("force") == nil {
		t.Fatal("expected force flag to be registered")
	}
}

func TestNewInitializesVault(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	cmd := New()
	cmd.SetIn(strings.NewReader("master\nmaster\n"))

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Initialized vault at ") {
		t.Fatalf("unexpected init output %q", out.String())
	}
}

func TestNewForceReinitializesExistingVault(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("KEEPASS_HOME", homeDir)

	first := New()
	first.SetIn(strings.NewReader("master\nmaster\n"))
	if err := first.Execute(); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	second := New()
	second.SetArgs([]string{"--force"})
	second.SetIn(strings.NewReader("master2\nmaster2\n"))

	var out bytes.Buffer
	second.SetOut(&out)

	if err := second.Execute(); err != nil {
		t.Fatalf("forced Execute() error = %v", err)
	}

	if !strings.Contains(out.String(), "Initialized vault at ") {
		t.Fatalf("unexpected force init output %q", out.String())
	}
}
