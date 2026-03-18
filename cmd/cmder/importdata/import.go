package importdata

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/transfer"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var path string
	var conflict string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import entries from a JSON export file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return common.UsageError("import requires --path")
			}

			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			doc, err := transfer.ReadExport(path)
			if err != nil {
				return err
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.ErrOrStderr())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			result, err := mgr.Import(masterPassword, doc, conflict)
			if err != nil {
				return common.MapError(err)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Import complete. Added=%d Overwrote=%d Skipped=%d\n", result.Added, result.Overwrote, result.Skipped)
			return err
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "source JSON file path")
	cmd.Flags().StringVar(&conflict, "conflict", transfer.ConflictFail, "conflict strategy: fail|skip|overwrite")
	return cmd
}
