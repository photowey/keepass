package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/audit"
	credentialaudit "github.com/photowey/keepass/internal/credentialaudit"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/internal/manager"
	"github.com/photowey/keepass/internal/prompt"
	"github.com/photowey/keepass/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewPrompter(in io.Reader, out io.Writer) *prompt.Prompter {
	return prompt.New(in, out)
}

func IsInteractive(in io.Reader) bool {
	file, ok := in.(*os.File)
	if !ok {
		// Best-effort: in tests we often pass non-*os.File readers to simulate
		// interactive input. In real CLI runs, stdin is an *os.File.
		return true
	}

	return term.IsTerminal(int(file.Fd()))
}

func IsNonInteractive(cmd *cobra.Command) bool {
	// Persistent flags are inherited by subcommands, so this works everywhere.
	value, err := cmd.Flags().GetBool("non-interactive")
	if err == nil {
		return value
	}

	value, _ = cmd.InheritedFlags().GetBool("non-interactive")
	return value
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

type entrySummary struct {
	Alias    string   `json:"alias"`
	Username string   `json:"username"`
	URI      string   `json:"uri,omitempty"`
	Note     string   `json:"note,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func PrintEntriesJSON(w io.Writer, entries []vault.Entry) error {
	out := make([]entrySummary, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entrySummary{
			Alias:    entry.Alias,
			Username: entry.Username,
			URI:      entry.URI,
			Note:     entry.Note,
			Tags:     entry.Tags,
		})
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
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

func PrintEntryJSON(w io.Writer, entry vault.Entry, reveal bool) error {
	if !reveal {
		entry.Password = ""
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
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

func PrintConfigText(w io.Writer, env home.Environment, cfg configs.Config, initialized bool) (int, error) {
	initValue := "false"
	if initialized {
		initValue = "true"
	}

	out := "" +
		fmt.Sprintf("Initialized: %s\n", initValue) +
		fmt.Sprintf("Root dir: %s\n", env.RootDir) +
		fmt.Sprintf("Config file: %s\n", env.ConfigFile) +
		fmt.Sprintf("Vault file: %s\n", cfg.ResolveVaultPath(env)) +
		"Tip: use `keepass config --json` for machine-readable output.\n"

	return fmt.Fprint(w, out)
}

func PrintAuditJSON(w io.Writer, report audit.Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

func PrintAuditText(w io.Writer, report audit.Report) (int, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Overall: %s\n", report.OverallStatus))
	builder.WriteString(fmt.Sprintf("Root dir: %s\n", report.RootDir))
	builder.WriteString(fmt.Sprintf("Config file: %s\n", report.ConfigFile))
	builder.WriteString(fmt.Sprintf("Vault file: %s\n", report.VaultFile))
	builder.WriteString(fmt.Sprintf("Config present: %t\n", report.Config.Present))
	builder.WriteString(fmt.Sprintf("Vault present: %t\n", report.Vault.Present))
	builder.WriteString(fmt.Sprintf("Configured Argon2id: time=%d memory_kib=%d threads=%d\n", report.Config.Argon2id.Time, report.Config.Argon2id.MemoryKiB, report.Config.Argon2id.Threads))
	builder.WriteString(fmt.Sprintf("Password preset: %s\n", report.Config.PasswordGenerator.Preset))
	builder.WriteString(fmt.Sprintf("Custom alphabet override: %t\n", report.Config.PasswordGenerator.UsesCustomAlphabet))

	if report.Vault.Argon2id != nil {
		builder.WriteString(fmt.Sprintf("Vault Argon2id: time=%d memory_kib=%d threads=%d\n", report.Vault.Argon2id.Time, report.Vault.Argon2id.MemoryKiB, report.Vault.Argon2id.Threads))
	}

	builder.WriteString("Checks:\n")
	for _, check := range report.Checks {
		builder.WriteString(fmt.Sprintf("- [%s] %s: %s\n", check.Status, check.Name, check.Message))
	}

	if len(report.Recommendations) > 0 {
		builder.WriteString("Recommendations:\n")
		for _, recommendation := range report.Recommendations {
			builder.WriteString(fmt.Sprintf("- %s\n", recommendation))
		}
	}

	builder.WriteString("Tip: use `keepass doctor --json` for machine-readable output.\n")
	return fmt.Fprint(w, builder.String())
}

func PrintCredentialAuditJSON(w io.Writer, report credentialaudit.Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

func PrintCredentialAuditText(w io.Writer, report credentialaudit.Report) (int, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Overall: %s\n", report.OverallStatus))
	builder.WriteString(fmt.Sprintf("Max password age days: %d\n", report.MaxAgeDays))
	builder.WriteString("Findings:\n")
	for _, finding := range report.Findings {
		builder.WriteString(fmt.Sprintf("- [%s] %s: %s\n", finding.Type, strings.Join(finding.Aliases, ", "), finding.Message))
	}
	builder.WriteString("Tip: use `keepass audit --json` for machine-readable output.\n")
	return fmt.Fprint(w, builder.String())
}
