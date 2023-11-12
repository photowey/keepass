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

package filez

import (
	"os"
	"path/filepath"
)

func exists(names ...string) bool {
	_, err := os.Stat(filepath.Join(names...))

	return err == nil || os.IsExist(err)
}

func DirExists(path string) bool {
	return exists(path)
}

func FileExists(target, name string) bool {
	return exists(target, name)
}

func FileNotExists(target, name string) bool {
	return !FileExists(target, name)
}

func Write(filename, content string) {
	f, err := os.Create(filename)
	MustCheck(err)
	defer Close(f)
	_, err = f.WriteString(content)
	MustCheck(err)
}

func Close(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}
