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

package configs_test

import (
	"fmt"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func ExampleConfig_ResolveVaultPath() {
	env := home.Environment{
		RootDir:         "/tmp/example-home/.keepass",
		ConfigFile:      "/tmp/example-home/.keepass/keepass.config.json",
		DefaultVault:    "/tmp/example-home/.keepass/keepass.kp",
		ResolvedHomeDir: "/tmp/example-home",
	}

	cfg := configs.Default(env)
	fmt.Println(cfg.ResolveVaultPath(env))

	// Output:
	// /tmp/example-home/.keepass/keepass.kp
}

func ExamplePasswordGenerator_EffectiveAlphabet() {
	generator := configs.PasswordGenerator{Alphabet: "abc123"}
	alphabet, err := generator.EffectiveAlphabet()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(alphabet)

	// Output:
	// abc123
}
