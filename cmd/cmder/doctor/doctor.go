package doctor

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/audit"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Audit local vault and config health",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := audit.Collect()
			if err != nil {
				return err
			}

			if jsonOut {
				return common.PrintAuditJSON(cmd.OutOrStdout(), report)
			}

			_, err = common.PrintAuditText(cmd.OutOrStdout(), report)
			return err
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}
