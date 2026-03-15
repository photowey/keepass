package common

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/prompt"
	"github.com/photowey/keepass/internal/vault"
)

func NewPrompter(in io.Reader, out io.Writer) *prompt.Prompter {
	return prompt.New(in, out)
}

func LoadManager() (*manager.Manager, error) {
	return manager.LoadCurrent()
}

func LoadOrCreateManager() (*manager.Manager, error) {
	return manager.LoadOrCreateCurrent()
}

func PromptMasterPassword(p *prompt.Prompter) (string, error) {
	return p.AskSecret("Master password")
}

func PromptNewMasterPassword(p *prompt.Prompter) (string, error) {
	return p.AskSecretWithConfirmation("Master password", "Confirm master password")
}

func ParseTags(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}

	items := strings.Split(input, ",")
	tags := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	return tags
}

func PrintEntries(w io.Writer, entries []vault.Entry) {
	tab := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(tab, "ALIAS\tUSERNAME\tTAGS\tURI\tNOTE")
	for _, entry := range entries {
		_, _ = fmt.Fprintf(
			tab,
			"%s\t%s\t%s\t%s\t%s\n",
			entry.Alias,
			entry.Username,
			strings.Join(entry.Tags, ","),
			entry.URI,
			entry.Note,
		)
	}
	_ = tab.Flush()
}

func PrintEntry(w io.Writer, entry vault.Entry, reveal bool) {
	passwordValue := "[hidden]"
	if reveal {
		passwordValue = entry.Password
	}

	_, _ = fmt.Fprintf(w, "Alias: %s\n", entry.Alias)
	_, _ = fmt.Fprintf(w, "Username: %s\n", entry.Username)
	_, _ = fmt.Fprintf(w, "Password: %s\n", passwordValue)
	_, _ = fmt.Fprintf(w, "URI: %s\n", entry.URI)
	_, _ = fmt.Fprintf(w, "Note: %s\n", entry.Note)
	_, _ = fmt.Fprintf(w, "Tags: %s\n", strings.Join(entry.Tags, ", "))
}

func PrintConfig(w io.Writer, env home.Environment, cfg configs.Config, initialized bool) error {
	payload := map[string]any{
		"initialized": initialized,
		"root_dir":    env.RootDir,
		"config_file": env.ConfigFile,
		"vault_file":  cfg.ResolveVaultPath(env),
		"config":      cfg,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}
