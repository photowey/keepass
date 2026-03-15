/*
 * Copyright © 2023 the original author or authors.
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

package config

import (
	"errors"

	"github.com/photowey/keepass/cmd/cmder/common"
	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show effective configuration and resolved paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := home.Detect()
			if err != nil {
				return err
			}

			cfg, err := configs.Load(env)
			initialized := true
			if err != nil {
				if !errors.Is(err, configs.ErrConfigNotFound) {
					return err
				}
				initialized = false
				cfg = configs.Default(env)
			}

			return common.PrintConfig(cmd.OutOrStdout(), env, cfg, initialized)
		},
	}
}
