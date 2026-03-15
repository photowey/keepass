package get

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var reveal bool

	cmd := &cobra.Command{
		Use:   "get <alias>",
		Short: "Show a single entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.OutOrStdout())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			entry, err := mgr.Get(masterPassword, args[0])
			if err != nil {
				return err
			}

			common.PrintEntry(cmd.OutOrStdout(), entry, reveal)
			return nil
		},
	}

	cmd.Flags().BoolVar(&reveal, "reveal", false, "reveal the plaintext password")

	return cmd
}
