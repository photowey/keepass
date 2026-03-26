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

package home

import (
	"fmt"
	"os"
)

func ExampleDetect() {
	originalValue, hadValue := os.LookupEnv(EnvKeepassHomePath)
	defer func() {
		if hadValue {
			_ = os.Setenv(EnvKeepassHomePath, originalValue)
			return
		}
		_ = os.Unsetenv(EnvKeepassHomePath)
	}()

	_ = os.Setenv(EnvKeepassHomePath, "/tmp/keepass-example")

	env, err := Detect()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(env.RootDir)
	fmt.Println(env.ConfigFile)
	fmt.Println(env.DefaultVault)

	// Output:
	// /tmp/keepass-example
	// /tmp/keepass-example/keepass.config.json
	// /tmp/keepass-example/keepass.kp
}
