package root

import "testing"

func TestNewCommandRegistersAliasesAndSubcommands(t *testing.T) {
	cmd := NewCommand()

	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "kee" || cmd.Aliases[1] != "kps" {
		t.Fatalf("unexpected aliases %#v", cmd.Aliases)
	}

	if cmd.PersistentFlags().Lookup("non-interactive") == nil {
		t.Fatal("expected non-interactive persistent flag")
	}

	expected := []string{
		"init", "add", "list", "get", "update", "audit", "rotate",
		"export", "import", "backup", "restore", "delete", "doctor",
		"rehash", "config", "completion",
	}

	commands := cmd.Commands()
	for _, use := range expected {
		found := false
		for _, sub := range commands {
			if sub.Name() == use {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected subcommand %q to be registered", use)
		}
	}
}
