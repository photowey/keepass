package rehash

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rehash",
		Short: "Rewrite the vault using the current security parameters",
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

			entryCount, err := mgr.Rehash(masterPassword)
			if err != nil {
				return common.MapError(err)
			}

			argon2id := mgr.Config().Security.Argon2id
			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"Rehashed vault at %s with Argon2id(time=%d, memory_kib=%d, threads=%d). Entries: %d\n",
				mgr.VaultPath(),
				argon2id.Time,
				argon2id.MemoryKiB,
				argon2id.Threads,
				entryCount,
			)
			return err
		},
	}

	return cmd
}
