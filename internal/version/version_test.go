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

import "testing"

func TestNowPrefixesReleaseVersion(t *testing.T) {
	originalVersion := version
	t.Cleanup(func() {
		version = originalVersion
	})

	version = "1.2.3"

	if got := Now(); got != "v1.2.3" {
		t.Fatalf("expected prefixed release version, got %q", got)
	}
}

func TestNowKeepsDevVersion(t *testing.T) {
	originalVersion := version
	t.Cleanup(func() {
		version = originalVersion
	})

	version = "dev"

	if got := Now(); got != "dev" {
		t.Fatalf("expected dev version, got %q", got)
	}
}

func TestSummaryIncludesMetadata(t *testing.T) {
	originalVersion := version
	originalCommit := commit
	originalBuildTime := buildTime
	t.Cleanup(func() {
		version = originalVersion
		commit = originalCommit
		buildTime = originalBuildTime
	})

	version = "v2.0.0"
	commit = "abc1234"
	buildTime = "2026-03-18T12:00:00Z"

	got := Summary()
	want := "v2.0.0 (commit abc1234, built 2026-03-18T12:00:00Z)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
