package delete

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
				return common.MapError(err)
			}

			in := cmd.InOrStdin()
			prompter := common.NewPrompter(in, cmd.ErrOrStderr())
			interactive := common.IsInteractive(in) && !common.IsNonInteractive(cmd)
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			entry, err := mgr.Get(masterPassword, args[0])
			if err != nil {
				return common.MapError(err)
			}

			if !yes {
				if !interactive {
					return common.UsageError("delete requires --yes in non-interactive mode")
				}

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
				return common.MapError(err)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Deleted entry %s\n", removed.Alias)
			return err
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "skip deletion confirmation")

	return cmd
}
