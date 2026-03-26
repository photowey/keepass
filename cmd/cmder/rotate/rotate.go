/*
 * Copyright © 2023-present the keepass authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rotate

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
	var passwordValue string
	var generate bool
	var reveal bool
	var copyToClipboard bool
	var copyTimeoutSeconds int

	cmd := &cobra.Command{
		Use:   "rotate <alias>",
		Short: "Rotate the password for a single entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if generate && passwordValue != "" {
				return common.UsageError("cannot use --generate and --password together")
			}

			if !generate && passwordValue == "" {
				generate = true
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

			var passwordPtr *string
			if passwordValue != "" {
				passwordPtr = &passwordValue
			}

			entry, generated, err := mgr.Rotate(masterPassword, args[0], passwordPtr, generate)
			if err != nil {
				return common.MapError(err)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Rotated password for %s\n", entry.Alias)
			if err != nil {
				return err
			}

			if copyToClipboard {
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

			if reveal || generated {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "New password: %s\n", entry.Password)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&passwordValue, "password", "", "manual password value")
	cmd.Flags().BoolVar(&generate, "generate", false, "generate a new password")
	cmd.Flags().BoolVar(&reveal, "reveal", false, "print the new password once")
	cmd.Flags().BoolVar(&copyToClipboard, "copy", false, "copy the new password to clipboard")
	cmd.Flags().IntVar(&copyTimeoutSeconds, "copy-timeout", 15, "clipboard clear timeout in seconds (0 disables auto-clear)")
	return cmd
}
