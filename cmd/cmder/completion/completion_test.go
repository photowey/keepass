package completion

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRejectsUnsupportedShell(t *testing.T) {
	root := &cobra.Command{Use: "keepass"}
	cmd := New(root)
	cmd.SetArgs([]string{"invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected unsupported shell error")
	}

	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestNewGeneratesBashCompletion(t *testing.T) {
	root := &cobra.Command{Use: "keepass"}
	cmd := New(root)
	cmd.SetArgs([]string{"bash"})

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if strings.TrimSpace(out.String()) == "" {
		t.Fatal("expected completion script output")
	}
}

func TestNewGeneratesOtherShellCompletions(t *testing.T) {
	for _, shell := range []string{"zsh", "fish", "powershell"} {
		t.Run(shell, func(t *testing.T) {
			root := &cobra.Command{Use: "keepass"}
			cmd := New(root)
			cmd.SetArgs([]string{shell})

			var out bytes.Buffer
			cmd.SetOut(&out)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute(%s) error = %v", shell, err)
			}

			if strings.TrimSpace(out.String()) == "" {
				t.Fatalf("expected non-empty %s completion script", shell)
			}
		})
	}
}
