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
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/photowey/keepass/pkg/filez"
)

const (
	Home            = ".keepass"
	KeepassConfig   = "keepass.json"
	KeepassDatabase = "database"
)

var (
	Usr, _   = user.Current()
	Dir      = filepath.Join(Usr.HomeDir, string(os.PathSeparator), Home)
	Database = filepath.Join(Dir, string(os.PathSeparator), KeepassDatabase)
)

func KeepassHome() {
	keepassHome := Dir
	keepassDatabase := Database

	if ok := filez.DirExists(keepassHome); !ok {
		if err := os.MkdirAll(keepassHome, os.ModePerm); err != nil {
			panic(fmt.Sprintf("mkdir keepass home dir:%s error:%v", keepassHome, err))
		}
	}

	if ok := filez.DirExists(keepassDatabase); !ok {
		if err := os.MkdirAll(keepassDatabase, os.ModePerm); err != nil {
			panic(fmt.Sprintf("mkdir keepass database home dir:%s error:%v", keepassDatabase, err))
		}
	}

	if filez.FileNotExists(keepassHome, KeepassConfig) {
		buf := bytes.NewBufferString(keepassConfigData)
		keepassConfigFile := filepath.Join(keepassHome, strings.ToLower(KeepassConfig))
		if err := os.WriteFile(keepassConfigFile, buf.Bytes(), 0o644); err != nil {
			panic(fmt.Sprintf("writing file %s: %v", keepassConfigFile, err))
		}
	}
}
