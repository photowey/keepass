package backup

import (
	"fmt"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var path string
	var force bool

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Create an encrypted local backup bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				return common.UsageError("backup requires --path")
			}

			mgr, err := common.LoadManager()
			if err != nil {
				return common.MapError(err)
			}

			outPath, err := mgr.CreateBackup(path, force)
			if err != nil {
				return common.MapError(err)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created backup bundle at %s\n", outPath)
			return err
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "destination backup directory")
	cmd.Flags().BoolVar(&force, "force", false, "allow writing into an existing non-empty backup directory")
	return cmd
}
