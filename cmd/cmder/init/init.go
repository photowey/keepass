package init

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the encrypted vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadOrCreateManager()
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptNewMasterPassword(prompter)
			if err != nil {
				return err
			}

			if err := mgr.Initialize(masterPassword, force); err != nil {
				return err
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Initialized vault at %s\n",
				mgr.VaultPath(),
			)
			return err
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing vault")

	return cmd
}
