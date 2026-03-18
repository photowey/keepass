package audit

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var maxAgeDays int
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit credential hygiene inside the unlocked vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			report, err := mgr.AuditCredentials(masterPassword, maxAgeDays)
			if err != nil {
				return common.MapError(err)
			}

			if jsonOut {
				return common.PrintCredentialAuditJSON(cmd.OutOrStdout(), report)
			}

			_, err = common.PrintCredentialAuditText(cmd.OutOrStdout(), report)
			return err
		},
	}

	cmd.Flags().IntVar(&maxAgeDays, "max-password-age-days", 180, "mark passwords older than this number of days as stale")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}
