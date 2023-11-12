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

package home

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/photowey/keepass/pkg/filez"
)

const (
	KeepassConfig = "keepass.json"
)

var (
	Home   = ".keepass"
	Usr, _ = user.Current()
	Dir    = filepath.Join(Usr.HomeDir, string(os.PathSeparator), Home)
)

func KeepassHome() {
	keepassHome := Dir
	if ok := filez.DirExists(keepassHome); !ok {
		if err := os.MkdirAll(keepassHome, os.ModePerm); err != nil {
			// panic(fmt.Sprintf("mkdir keepass home dir:%s error:%v", keepassHome, err))
		}
	}
}
