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

import "fmt"

func ExampleSummary() {
	originalVersion := version
	originalCommit := commit
	originalBuildTime := buildTime
	defer func() {
		version = originalVersion
		commit = originalCommit
		buildTime = originalBuildTime
	}()

	version = "1.2.3"
	commit = "abc1234"
	buildTime = "2026-03-18T12:00:00Z"

	fmt.Println(Now())
	fmt.Println(Summary())

	// Output:
	// v1.2.3
	// v1.2.3 (commit abc1234, built 2026-03-18T12:00:00Z)
}
