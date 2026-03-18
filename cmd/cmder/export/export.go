package export

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/transfer"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export unlocked entries to a JSON file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return common.UsageError("export requires --path")
			}

			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			doc, err := mgr.Export(masterPassword)
			if err != nil {
				return common.MapError(err)
			}

			if err := transfer.WriteExport(path, doc); err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte("Exported entries to " + path + "\n"))
			return err
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "destination JSON file path")
	return cmd
}
