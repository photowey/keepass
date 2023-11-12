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

package stringz

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func String(source any) string {
	return fmt.Sprintf("%v", source)
}

func ReplaceTemplate(template string, args ...any) string {
	return fmt.Sprintf(template, args...)
}

func IsBlankString(str string) bool {
	return "" == str
}

func IsNotBlankString(str string) bool {
	return !IsBlankString(str)
}

func IsEmptyStringSlice(target []string) bool {
	return len(target) == 0
}

func IsNotEmptyStringSlice(target []string) bool {
	return !IsEmptyStringSlice(target)
}

// ----------------------------------------------------------------

func Tail(content, separator string) string {
	if IsBlankString(content) {
		return ""
	}
	items := strings.Split(content, separator)

	return items[len(items)-1]
}

// ----------------------------------------------------------------

func Pascal(content string) string {
	if IsBlankString(content) {
		return ""
	}
	items := strings.ToUpper(content[0:1]) + content[1:]

	return items
}

// ----------------------------------------------------------------

func ToPath(content, separator string) string {
	if IsBlankString(content) {
		return ""
	}

	pathSeparator := string(os.PathSeparator)
	items := strings.ReplaceAll(content, separator, pathSeparator)

	return items
}

// ----------------------------------------------------------------

func ArrayContains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.EqualFold(a, needle) {
			return true
		}
	}

	return false
}

func ArrayNotContains(haystack []string, needle string) bool {
	return !ArrayContains(haystack, needle)
}

// ----------------------------------------------------------------

func ToProjectPath(content string) string {
	if IsBlankString(content) {
		return ""
	}

	return ToPath(content, "-")
}

func isWindows() bool {
	goos := runtime.GOOS
	if strings.HasPrefix(strings.ToLower(goos), "windows") {
		return true
	}
	return false
}
