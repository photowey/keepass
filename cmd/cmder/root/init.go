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

package root

import (
	"path/filepath"
	"strings"

	"github.com/photowey/keepass/configs"
	"github.com/photowey/keepass/internal/home"
	"github.com/photowey/keepass/pkg/filez"
)

func onInit() {
	home.KeepassHome()
	configLoad()
}

func configLoad() {
	keepassHome := home.Dir
	keepassConfigFile := filepath.Join(keepassHome, strings.ToLower(home.KeepassConfig))
	if filez.FileExists(keepassHome, home.KeepassConfig) {
		configs.Init(keepassConfigFile)

		return
	}

	// fmt.Printf("keepass config file not exists\n")
}
