package get

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/internal/clipboard"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var reveal bool
	var jsonOut bool
	var copyToClipboard bool
	var copyTimeoutSeconds int

	cmd := &cobra.Command{
		Use:   "get <alias>",
		Short: "Show a single entry",
		Args:  cobra.ExactArgs(1),
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

			entry, err := mgr.Get(masterPassword, args[0])
			if err != nil {
				return common.MapError(err)
			}

			if copyToClipboard {
				if entry.Password == "" {
					return fmt.Errorf("empty password cannot be copied")
				}

				if err := clipboard.Copy(entry.Password); err != nil {
					return fmt.Errorf("copy to clipboard: %w", err)
				}

				if copyTimeoutSeconds > 0 {
					_, err = fmt.Fprintf(cmd.OutOrStdout(), "Copied password to clipboard. Waiting %ds before clearing...\n", copyTimeoutSeconds)
					if err != nil {
						return err
					}

					ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
					defer stop()

					_, err = clipboard.ClearAfter(ctx, time.Duration(copyTimeoutSeconds)*time.Second)
					if err != nil {
						return fmt.Errorf("clear clipboard: %w", err)
					}

					_, err = fmt.Fprintln(cmd.OutOrStdout(), "Clipboard cleared.")
					return err
				}

				_, err = fmt.Fprintln(cmd.OutOrStdout(), "Copied password to clipboard.")
				return err
			}

			if jsonOut {
				return common.PrintEntryJSON(cmd.OutOrStdout(), entry, reveal)
			}

			common.PrintEntry(cmd.OutOrStdout(), entry, reveal)
			return nil
		},
	}

	cmd.Flags().BoolVar(&reveal, "reveal", false, "reveal the plaintext password")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	cmd.Flags().BoolVar(&copyToClipboard, "copy", false, "copy password to clipboard (does not print plaintext)")
	cmd.Flags().IntVar(&copyTimeoutSeconds, "copy-timeout", 15, "clipboard clear timeout in seconds (0 disables auto-clear)")

	return cmd
}
