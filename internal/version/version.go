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

package version

import (
	"fmt"
	"strings"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func Now() string {
	value := strings.TrimSpace(version)
	switch {
	case value == "", value == "unknown":
		return "dev"
	case value == "dev":
		return value
	case strings.HasPrefix(value, "v"):
		return value
	default:
		return "v" + value
	}
}

func Commit() string {
	return normalizeValue(commit, "unknown")
}

func BuildTime() string {
	return normalizeValue(buildTime, "unknown")
}

func Summary() string {
	return fmt.Sprintf("%s (commit %s, built %s)", Now(), Commit(), BuildTime())
}

func normalizeValue(value, fallback string) string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return fallback
	}

	return normalized
}
