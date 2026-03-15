package list

import (
	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/manager"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var tags []string

	cmd := &cobra.Command{
		Use:   "list [query]",
		Short: "List entries",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := common.LoadManager()
			if err != nil {
				return err
			}

			query := ""
			if len(args) == 1 {
				query = args[0]
			}

			prompter := common.NewPrompter(cmd.InOrStdin(), cmd.OutOrStdout())
			masterPassword, err := common.PromptMasterPassword(prompter)
			if err != nil {
				return err
			}

			entries, err := mgr.List(masterPassword, manager.ListFilter{
				Query: query,
				Tags:  tags,
			})
			if err != nil {
				return err
			}

			common.PrintEntries(cmd.OutOrStdout(), entries)
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&tags, "tag", nil, "filter by tag (repeatable, all tags must match)")

	return cmd
}
