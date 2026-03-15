package del

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <alias>",
		Short: "Delete an entry",
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

			if !yes {
				confirmed, err := prompter.Confirm(fmt.Sprintf("Delete entry %s", entry.Alias), false)
				if err != nil {
					return err
				}
				if !confirmed {
					_, err = fmt.Fprintln(cmd.OutOrStdout(), "Deletion cancelled.")
					return err
				}
			}

			removed, err := mgr.Delete(masterPassword, args[0])
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Deleted entry %s\n", removed.Alias)
			return err
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "skip deletion confirmation")

	return cmd
}
