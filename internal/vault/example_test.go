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

package vault

import (
	"fmt"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
)

func ExampleEncode() {
	env := home.Environment{
		RootDir:         "/tmp/example-home/.keepass",
		ConfigFile:      "/tmp/example-home/.keepass/keepass.config.json",
		DefaultVault:    "/tmp/example-home/.keepass/keepass.kp",
		ResolvedHomeDir: "/tmp/example-home",
	}

	cfg := configs.Default(env)
	cfg.Security.Argon2id.Time = 1
	cfg.Security.Argon2id.MemoryKiB = 8 * 1024
	cfg.Security.Argon2id.Threads = 1
	cfg.Security.Argon2id.KeyLength = 32

	data, err := Encode(1, []byte("secret payload"), "master-password", cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	plaintext, err := Decode(data, "master-password")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintext))

	// Output:
	// secret payload
}
