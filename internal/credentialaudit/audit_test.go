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

package credentialaudit

import (
	"testing"
	"time"

	"github.com/photowey/keepass/internal/vault"
)

func TestAnalyzeFindsStaleDuplicateAndMissingMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0).UTC()
	entries := []vault.Entry{
		{
			Alias:             "github",
			Username:          "hellopass",
			Password:          "shared-secret",
			PasswordUpdatedAt: now.AddDate(0, 0, -365),
		},
		{
			Alias:             "gitea",
			Username:          "hellopass",
			Password:          "shared-secret",
			URI:               "https://gitea.example.com",
			PasswordUpdatedAt: now,
		},
	}

	report := Analyze(entries, 180, now)
	if report.OverallStatus != "warn" {
		t.Fatalf("expected warn status, got %+v", report)
	}

	if len(report.Findings) < 3 {
		t.Fatalf("expected multiple findings, got %+v", report)
	}
}
