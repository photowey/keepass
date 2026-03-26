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

package password_test

import (
	"fmt"
	"strings"

	"github.com/photowey/keepass/internal/password"
)

func ExampleGenerate() {
	value, err := password.Generate(8, "ab")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(value))
	fmt.Println(strings.Trim(value, "ab") == "")

	// Output:
	// 8
	// true
}

func ExampleAlphabetForPreset() {
	alphabet, err := password.AlphabetForPreset(password.PresetCompatible)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(strings.Contains(alphabet, "-"))
	fmt.Println(strings.Contains(alphabet, "0"))

	// Output:
	// true
	// false
}
