package restore

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var path string
	var force bool

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore config and vault from a backup bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return common.UsageError("restore requires --path")
			}

			if err := manager.RestoreCurrent(path, force); err != nil {
				return common.MapError(err)
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Restored backup bundle from %s\n", path)
			return err
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "source backup directory")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing config and vault files")
	return cmd
}
